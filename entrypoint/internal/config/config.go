package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	PBCONFIG_PATH                 = "config.pbtxt"
	ENV_PREFIX                    = "TRTLLM"
	ENV_SPLITTER                  = "__"
	ENV_BLS_IDENTIFIER            = "BLS"
	ENV_ENGINE_IDENTIFIER         = "MODEL"
	ENV_ENSEMBLE_IDENTIFIER       = "ENSEMBLE"
	ENV_POSTPROCESSING_IDENTIFIER = "POSTPROCESSING"
	ENV_PREPROCESSING_IDENTIFIER  = "PREPROCESSING"
	DIR_BLS_IDENTIFIER            = "tensorrt_llm_bls"
	DIR_ENGINE_IDENTIFIER         = "tensorrt_llm"
	DIR_ENSEMBLE_IDENTIFIER       = "ensemble"
	DIR_POSTPROCESSING_IDENTIFIER = "postprocessing"
	DIR_PREPROCESSING_IDENTIFIER  = "preprocessing"
)

type Opts struct {
	GrpcPort    int
	HttpPort    int
	MetricsPort int
	ModelRepo   string
	Verbosity   int
	WorldSize   int
}

type inferenceConfig struct {
	bls            []string
	engine         []string
	ensemble       []string
	postprocessing []string
	preprocessing  []string
}

func replace(path string, args []string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	replacer := strings.NewReplacer(args...)
	newContent := replacer.Replace(string(content))
	return os.WriteFile(path, []byte(newContent), info.Mode().Perm())
}

func FillConfigTemplatesFromEnv(modelRepo string) error {
	var cfg inferenceConfig
	for _, entry := range os.Environ() {
		k, v, _ := strings.Cut(entry, "=")
		kparts := strings.SplitN(k, ENV_SPLITTER, 3)
		if len(kparts) != 3 || kparts[0] != ENV_PREFIX {
			continue
		}
		toReplace := fmt.Sprintf("${%s}", strings.ToLower(kparts[2]))
		switch kparts[1] {
		case ENV_BLS_IDENTIFIER:
			cfg.bls = append(cfg.bls, toReplace, v)
		case ENV_ENGINE_IDENTIFIER:
			cfg.engine = append(cfg.engine, toReplace, v)
		case ENV_ENSEMBLE_IDENTIFIER:
			cfg.ensemble = append(cfg.ensemble, toReplace, v)
		case ENV_POSTPROCESSING_IDENTIFIER:
			cfg.postprocessing = append(cfg.postprocessing, toReplace, v)
		case ENV_PREPROCESSING_IDENTIFIER:
			cfg.preprocessing = append(cfg.preprocessing, toReplace, v)
		}
	}

	for _, elem := range []struct {
		path string
		args []string
	}{
		{path: filepath.Join(modelRepo, DIR_BLS_IDENTIFIER, PBCONFIG_PATH), args: cfg.bls},
		{path: filepath.Join(modelRepo, DIR_ENSEMBLE_IDENTIFIER, PBCONFIG_PATH), args: cfg.ensemble},
		{path: filepath.Join(modelRepo, DIR_ENGINE_IDENTIFIER, PBCONFIG_PATH), args: append([]string{"${engine_dir}", filepath.Join(modelRepo, DIR_ENGINE_IDENTIFIER, "1")}, cfg.engine...)},
		{path: filepath.Join(modelRepo, DIR_POSTPROCESSING_IDENTIFIER, PBCONFIG_PATH), args: cfg.postprocessing},
		{path: filepath.Join(modelRepo, DIR_PREPROCESSING_IDENTIFIER, PBCONFIG_PATH), args: cfg.preprocessing},
	} {
		if err := replace(elem.path, elem.args); err != nil {
			return fmt.Errorf("failed to update configuration %s: %v", elem.path, err)
		}
	}

	return nil
}

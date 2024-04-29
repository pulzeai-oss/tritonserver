# Triton Inference Server + TensorRT-LLM

A thin wrapper for simplifying the deployment of [Triton Inference Server](https://github.com/triton-inference-server/server) with the [TensorRT-LLM](https://github.com/triton-inference-server/tensorrtllm_backend) backend.

Instructions for deployment of TensorRT-LLM engines leave much to be desired, for the following reasons:
- Requires arbitrary scripts:
    - Bootstrapping model repository from configuration templates
    - Rendering configuration templates with user-specified values
- No clear compatibility matrix between backend and server versions
- Python entrypoint fails to forward signals, preventing graceful termination

We simplify this by:
- Baking the necessary configuration templates into the docker image
- Using an entrypoint that:
    - Renders the configuration templates on the fly from environment variables
    - Is replaced by `mpirun` via an `exec` call that correctly forwards signals to the respective workers, allowing in-flight requests to be drained before termination
- Bundling the correct revision of the TensorRT-LLM backend for the base [Triton NGC image](https://catalog.ngc.nvidia.com/orgs/nvidia/containers/tritonserver)

## Configuration

[Model configuration](https://github.com/triton-inference-server/tensorrtllm_backend/tree/41fe3a6a9daa12c64403e084298c6169b07d489d?tab=readme-ov-file#modify-the-model-configuration) is set via environment variables. For instance, setting `TRTLLM__ENGINE__TRITON_MAX_BATCH_SIZE=1` replaces the `${triton_max_batch_size}` template variable in `tensorrt_llm/config.pbtxt`. Similarly, setting `TRTLLM__BLS__ACCUMULATE_TOKENS=false` replaces the `${accumulate_tokens}` template variable in `tensorrt_llm_bls/config.pbtxt`. You get the idea!

The pre-built TensorRT-LLM engine should be mounted into `/srv/run/repo/tensorrt_llm/1`.

See [examples](./examples) for instructions on deploying on [Kubernetes](./examples/k8s) or locally using [docker-compose](./examples/docker-compose).

## TODOs
- [ ] Helm chart

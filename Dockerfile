FROM --platform=${BUILDPLATFORM} golang:1.22-alpine AS entrypoint

ARG TARGETARCH
ENV GOARCH "${TARGETARCH}"
ENV GOOS linux

WORKDIR /build

COPY entrypoint/go.mod entrypoint/go.sum ./
RUN go mod download
COPY entrypoint ./
RUN go build -ldflags="-w -s" -o dist/entrypoint main.go


FROM nvcr.io/nvidia/tritonserver:24.03-trtllm-python-py3

# Install tini
RUN apt-get update && apt-get install --yes --no-install-recommends tini && rm -rf /var/lib/apt/lists/*

# Create nonroot account
RUN groupadd --gid 65532 nonroot \
    && useradd --no-log-init --create-home \
        --uid 65532 \
        --gid 65532 \
        --shell /sbin/nologin \
        nonroot

ENV TRTLLM__MODEL_REPO_DIR /srv/run/repo
WORKDIR ${TRTLLM__MODEL_REPO_DIR}

# Copy shell model repository
COPY tensorrtllm_backend/all_models/inflight_batcher_llm ${TRTLLM__MODEL_REPO_DIR}

COPY --from=entrypoint /build/dist/entrypoint /entrypoint
ENTRYPOINT [ "/usr/bin/tini", "-g", "--", "/entrypoint"]

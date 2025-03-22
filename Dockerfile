# Build the manager binary
FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.24@sha256:52ff1b35ff8de185bf9fd26c70077190cd0bed1e9f16a2d498ce907e5c421268 AS builder

ARG TARGETOS
ARG TARGETARCH
ENV CGO_ENABLED=0

WORKDIR /workspace
COPY . .

# Build greenhouse operator and tooling.
RUN --mount=type=cache,target=/go/pkg/mod \
	--mount=type=cache,target=/root/.cache/go-build \
	make action-build CGO_ENABLED=${CGO_ENABLED} GOOS=${TARGETOS} GOARCH=${TARGETARCH}

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM --platform=${BUILDPLATFORM:-linux/amd64} gcr.io/distroless/static:nonroot@sha256:b35229a3a6398fe8f86138c74c611e386f128c20378354fc5442811700d5600d
LABEL source_repository="https://github.com/cloudoperators/greenhouse"
WORKDIR /
COPY --from=builder /workspace/bin/* .
USER 65532:65532

RUN ["/greenhouse", "--version"]
ENTRYPOINT ["/greenhouse"]

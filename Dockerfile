# Build the manager binary
FROM golang:1.21 as builder

WORKDIR /workspace

COPY Makefile .
RUN make controller-gen envtest
COPY . .

# Build greenhouse operator and tooling.
RUN --mount=type=cache,target=/go/pkg/mod \
	--mount=type=cache,target=/root/.cache/go-build \
	make test build CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
LABEL source_repository="https://github.com/cloudoperators/greenhouse"
WORKDIR /
COPY --from=builder /workspace/bin/* .
USER 65532:65532

RUN ["/greenhouse", "--version"]
ENTRYPOINT ["/greenhouse"]

# syntax=docker/dockerfile:1
FROM public.ecr.aws/docker/library/golang:1.24.0 AS builder

WORKDIR /src
COPY . .

RUN --mount=type=cache,target=/go/pkg/mod <<EOT
    go mod download
EOT

RUN <<EOT
    make install-build-deps
    make build
    chown 65532 eksctl
EOT

FROM public.ecr.aws/eks-distro/kubernetes/go-runner:v0.16.4-eks-1-32-6 AS go-runner
COPY --from=builder /src/eksctl /eksctl
ENTRYPOINT ["/eksctl"]

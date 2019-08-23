# Make sure to run the following commands after changes to this file are made:
# `make update-build-image-manifest && make push-build-image`

ARG BASE_BUILD_IMAGE_ID
FROM ${BASE_BUILD_IMAGE_ID}

WORKDIR /src
ENV CGO_ENABLED=0
COPY install-build-deps.sh go.mod go.sum /src/

# Install all go dependencies and remove the go caches in a single step to reduce the image footprint
# (caches won't be used later on, we overwrite them by volume-mounting)
RUN ./install-build-deps.sh && \
    go install github.com/goreleaser/goreleaser && \
    go build github.com/kubernetes-sigs/aws-iam-authenticator/cmd/aws-iam-authenticator && \
    rm -rf /root/.cache/go-build /go/pkg/mod



# Runtime dependencies. Build the root filesystem of the eksctl image at /out
RUN mkdir -p /out/etc/apk && cp -r /etc/apk/* /out/etc/apk/
RUN apk add --no-cache --initdb --root /out \
    alpine-baselayout \
    busybox \
    ca-certificates \
    coreutils \
    git \
    libc6-compat \
    openssh \
    && true

RUN mv ./aws-iam-authenticator /out/usr/local/bin/aws-iam-authenticator

ENV KUBECTL_VERSION v1.11.5
RUN curl --silent --location "https://dl.k8s.io/${KUBECTL_VERSION}/bin/linux/amd64/kubectl" --output /out/usr/local/bin/kubectl \
    && chmod +x /out/usr/local/bin/kubectl

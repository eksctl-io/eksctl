# Make sure to run the following commands after changes to this file are made:
# `make -f Makefile.docker update-build-image-tag && make -f Makefile.docker push-build-image`

# This digest corresponds to golang:1.14.2-alpine3.11
FROM golang@sha256:62cd35bbeb9aadff6764dd8809c788267d72b12066bb40c080431510bbe81e36 AS base

# Build-time dependencies
RUN apk add --no-cache \
    bash \
    curl \
    docker-cli \
    g++ \
    gcc \
    git \
    libsass-dev \
    make \
    musl-dev \
    jq \
    && true

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

ENV KUBECTL_VERSION v1.15.10
RUN curl --silent --location "https://dl.k8s.io/${KUBECTL_VERSION}/bin/linux/amd64/kubectl" --output /out/usr/local/bin/kubectl \
    && chmod +x /out/usr/local/bin/kubectl

# Remaining dependencies are controlled by go.mod
WORKDIR /src
ENV CGO_ENABLED=0 GOPROXY=https://proxy.golang.org,direct

RUN git config --global url."git@github.com:".insteadOf "https://github.com/"

COPY .requirements install-build-deps.sh go.mod go.sum /src/

# Install all build tools dependencies
RUN ./install-build-deps.sh

# The authenticator is a runtime dependency, so it needs to be in /out
RUN go install sigs.k8s.io/aws-iam-authenticator/cmd/aws-iam-authenticator \
    && mv $GOPATH/bin/aws-iam-authenticator /out/usr/local/bin/aws-iam-authenticator

# Add kubectl and aws-iam-authenticator to the PATH
ENV PATH="${PATH}:/out/usr/bin:/out/usr/local/bin"

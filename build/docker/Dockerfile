# Make sure to run the following commands after changes to this file are made:
# `make -f Makefile.docker update-build-image-tag && make -f Makefile.docker push-build-image`

FROM golang:1.17.0-alpine3.14@sha256:4ace2a06b957dec72202e03e71ccb661f082cc9d9a1849c0efef41ddc246bdba AS base

# Add kubectl and aws-iam-authenticator to the PATH
ENV PATH="${PATH}:/out/usr/bin:/out/usr/local/bin"

# Build-time dependencies
RUN apk add --no-cache \
    bash \
    curl \
    docker-cli \
    g++ \
    gcc \
    git \
    httpie \
    libsass-dev \
    make \
    musl-dev \
    jq \
    py3-setuptools \
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

RUN curl --silent --location "https://dl.k8s.io/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" --output /out/usr/local/bin/kubectl \
    && chmod +x /out/usr/local/bin/kubectl

ENV GITHUB_CLI_VERSION 1.14.0
RUN curl --silent --location "https://github.com/cli/cli/releases/download/v${GITHUB_CLI_VERSION}/gh_${GITHUB_CLI_VERSION}_linux_amd64.tar.gz" \
    | tar xvz -C /out/usr/local/bin --strip-components=2 gh_${GITHUB_CLI_VERSION}_linux_amd64/bin/gh && gh --version

# Remaining dependencies are controlled by go.mod
WORKDIR /src
ENV CGO_ENABLED=0 GOPROXY=https://proxy.golang.org,direct

# The authenticator is a runtime dependency, so it needs to be in /out
ENV AWS_IAM_AUTH_VERSION 0.5.2
RUN curl --silent --location "https://github.com/kubernetes-sigs/aws-iam-authenticator/releases/download/v${AWS_IAM_AUTH_VERSION}/aws-iam-authenticator_${AWS_IAM_AUTH_VERSION}_linux_amd64" --output /out/usr/local/bin/aws-iam-authenticator \
    && chmod +x /out/usr/local/bin/aws-iam-authenticator
COPY .requirements install-build-deps.sh go.mod go.sum /src/

# Install all build tools dependencies
RUN ./install-build-deps.sh

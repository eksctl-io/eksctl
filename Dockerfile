ARG EKSCTL_BUILD_IMAGE
FROM $EKSCTL_BUILD_IMAGE AS build

LABEL eksctl.builder=true

RUN apk add --no-cache \
      py-pip \
      python \
      python-dev \
      && true

RUN mkdir -p /out/etc/apk && cp -r /etc/apk/* /out/etc/apk/
RUN apk add --no-cache --initdb --root /out \
    alpine-baselayout \
    busybox \
    ca-certificates \
    coreutils \
    git \
    libc6-compat \
    libgcc \
    libstdc++ \
    python \
    && true

ENV EKSCTL $GOPATH/src/github.com/weaveworks/eksctl
RUN mkdir -p "$(dirname ${EKSCTL})"
COPY . $EKSCTL

ARG COVERALLS_TOKEN
ENV COVERALLS_TOKEN $COVERALLS_TOKEN

ARG TEST_TARGET

ENV JUNIT_REPORT_DIR $GOPATH/src/github.com/weaveworks/eksctl/test-results/ginkgo
RUN mkdir -p "${JUNIT_REPORT_DIR}"

WORKDIR $EKSCTL
RUN make $TEST_TARGET
RUN make lint
RUN make build \
    && cp ./eksctl /out/usr/local/bin/eksctl

RUN go build ./vendor/github.com/kubernetes-sigs/aws-iam-authenticator/cmd/aws-iam-authenticator \
    && cp ./aws-iam-authenticator /out/usr/local/bin/aws-iam-authenticator

RUN pip install --root=/out aws-mfa==0.0.12 awscli==1.16.34

WORKDIR /out

ENV KUBECTL_VERSION v1.11.5
RUN curl --silent --location "https://dl.k8s.io/${KUBECTL_VERSION}/bin/linux/amd64/kubectl" --output usr/local/bin/kubectl \
    && chmod +x usr/local/bin/kubectl

FROM scratch
CMD eksctl
COPY --from=build  /out /

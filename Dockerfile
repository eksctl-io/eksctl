ARG EKSCTL_BUILD_IMAGE
FROM $EKSCTL_BUILD_IMAGE AS build

LABEL eksctl.builder=true

RUN mkdir -p /out/etc/apk && cp -r /etc/apk/* /out/etc/apk/
RUN apk add --no-cache --initdb --root /out \
    alpine-baselayout \
    busybox \
    ca-certificates \
    coreutils \
    git \
    libc6-compat \
    && true

ENV KUBECTL_VERSION v1.11.5
RUN curl --silent --location "https://dl.k8s.io/${KUBECTL_VERSION}/bin/linux/amd64/kubectl" --output /out/usr/local/bin/kubectl \
    && chmod +x /out/usr/local/bin/kubectl

COPY . /src

ARG COVERALLS_TOKEN
ENV COVERALLS_TOKEN $COVERALLS_TOKEN

ARG TEST_TARGET
ARG GO_BUILD_TAGS

ENV JUNIT_REPORT_DIR /src/test-results/ginkgo
RUN mkdir -p "${JUNIT_REPORT_DIR}"

WORKDIR /src
RUN make $TEST_TARGET
RUN make build GO_BUILD_TAGS="${GO_BUILD_TAGS}" \
    && cp ./eksctl /out/usr/local/bin/eksctl
RUN make build-integration-test \
    && mkdir -p /out/usr/local/share/eksctl \
    && cp integration/*.yaml /out/usr/local/share/eksctl \
    && cp ./eksctl-integration-test /out/usr/local/bin/eksctl-integration-test

RUN go build github.com/kubernetes-sigs/aws-iam-authenticator/cmd/aws-iam-authenticator \
    && cp ./aws-iam-authenticator /out/usr/local/bin/aws-iam-authenticator

FROM scratch
CMD eksctl
COPY --from=build  /out /

ARG EKSCTL_BUILD_IMAGE
FROM $EKSCTL_BUILD_IMAGE AS build

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

ARG JUNIT_REPORT_FOLDER
ENV JUNIT_REPORT_FOLDER $JUNIT_REPORT_FOLDER

WORKDIR $EKSCTL
RUN make lint
RUN make test
RUN make build \
    && cp ./eksctl /out/usr/local/bin/eksctl

RUN go build ./vendor/github.com/kubernetes-sigs/aws-iam-authenticator/cmd/aws-iam-authenticator \
    && cp ./aws-iam-authenticator /out/usr/local/bin/aws-iam-authenticator

RUN pip install --root=/out aws-mfa==0.0.12 awscli==1.16.34

WORKDIR /out

ENV KUBECTL_VERSION v1.10.3
RUN curl --silent --location "https://dl.k8s.io/${KUBECTL_VERSION}/bin/linux/amd64/kubectl" --output usr/local/bin/kubectl \
    && chmod +x usr/local/bin/kubectl

FROM scratch
CMD eksctl
COPY --from=build  /out /

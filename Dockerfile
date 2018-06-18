ARG EKSCTL_BUILD_IMAGE
FROM $EKSCTL_BUILD_IMAGE AS build

RUN apk add --update \
      py-pip \
      python \
      python-dev \
      && true

RUN mkdir /out
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

RUN go get github.com/jteeuwen/go-bindata/go-bindata

ENV EKSCTL $GOPATH/src/github.com/weaveworks/eksctl
RUN mkdir -p "$(dirname ${EKSCTL})"
COPY . $EKSCTL

WORKDIR $EKSCTL
RUN make \
    && cp ./eksctl /out/usr/local/bin/eksctl

RUN go build ./vendor/github.com/heptio/authenticator/cmd/heptio-authenticator-aws \
    && cp ./heptio-authenticator-aws /out/usr/local/bin/heptio-authenticator-aws

RUN pip install --root=/out aws-mfa==0.0.12 awscli==1.15.40

WORKDIR /out

ENV KUBECTL_VERSION v1.10.3
RUN curl --silent --location "https://dl.k8s.io/${KUBECTL_VERSION}/bin/linux/amd64/kubectl" --output usr/local/bin/kubectl \
    && chmod +x usr/local/bin/kubectl

FROM scratch
CMD eksctl
COPY --from=build  /out /

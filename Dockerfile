ARG BUILD_IMAGE=public.ecr.aws/eksctl/eksctl-build:7b4e42f61d4655dcfac75894275b2a0c4d7bfe2c
FROM $BUILD_IMAGE as build

WORKDIR /src

COPY . /src

RUN make test
RUN make build \
    && cp ./eksctl /out/usr/local/bin/eksctl
RUN make build-integration-test \
    && mkdir -p /out/usr/local/share/eksctl \
    && cp -r integration/data/*.yaml integration/scripts /out/usr/local/share/eksctl \
    && cp ./eksctl-integration-test /out/usr/local/bin/eksctl-integration-test

FROM scratch
COPY --from=build /out /
ENTRYPOINT ["eksctl"]

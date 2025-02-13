ARG BUILD_IMAGE=public.ecr.aws/eksctl/eksctl-build:9b9ac54ee74f36fa22be8c7ec24f695dbd4c0d72
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

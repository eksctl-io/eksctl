ARG BUILD_IMAGE=public.ecr.aws/eksctl/eksctl-build:741e7e49004ca5aabe086bf130aaca1cdcbc5c44
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

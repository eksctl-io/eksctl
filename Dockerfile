ARG BUILD_IMAGE=/weaveworks/eksctl-build:5705efb32cdbf4531b282dba24dae03a872f5b31
FROM $BUILD_IMAGE as build

WORKDIR /src
ENV JUNIT_REPORT_DIR="${JUNIT_REPORT_DIR:-/src/test-results/ginkgo}"
RUN mkdir -p "${JUNIT_REPORT_DIR}"

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

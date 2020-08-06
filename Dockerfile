ARG BUILD_IMAGE=/weaveworks/eksctl-build:0849e149eca6668d7e1a5b400dbbbb5cac650e52
FROM $BUILD_IMAGE as build

WORKDIR /src
ENV JUNIT_REPORT_DIR="${JUNIT_REPORT_DIR:-/src/test-results/ginkgo}"

COPY . /src

RUN mkdir -p "${JUNIT_REPORT_DIR}"

RUN make test
RUN make build \
    && cp ./eksctl /out/usr/local/bin/eksctl
RUN make build-integration-test \
    && mkdir -p /out/usr/local/share/eksctl \
    && cp -r integration/data/*.yaml integration/scripts /out/usr/local/share/eksctl \
    && cp ./eksctl-integration-test /out/usr/local/bin/eksctl-integration-test

FROM scratch
CMD eksctl
COPY --from=build /out /

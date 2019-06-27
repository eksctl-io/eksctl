FROM scratch
CMD eksctl
COPY --from=weaveworks/eksctl-builder:latest /out /

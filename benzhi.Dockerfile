FROM scratch

ARG TARGETARCH
COPY .docker-bin/pagestore-${TARGETARCH}.test /pagestore.test

ENTRYPOINT ["/pagestore.test"]

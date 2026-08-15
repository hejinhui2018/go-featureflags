FROM scratch

ARG TARGETARCH
COPY .docker-bin/revisionstore-${TARGETARCH}.test /revisionstore.test

ENTRYPOINT ["/revisionstore.test"]

FROM scratch

ARG TARGETARCH
COPY .docker-bin/ttlstore-${TARGETARCH}.test /ttlstore.test

ENTRYPOINT ["/ttlstore.test"]

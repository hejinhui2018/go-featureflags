FROM scratch

ARG TARGETARCH
COPY .docker-bin/cursorstore-${TARGETARCH}.test /cursorstore.test

ENTRYPOINT ["/cursorstore.test"]

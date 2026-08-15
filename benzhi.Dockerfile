ARG BASE_IMAGE=golang:1.22
FROM ${BASE_IMAGE}

WORKDIR /src
COPY go.mod ./
COPY README.md BUG_REPRO.md ./
COPY cmd ./cmd
COPY batchstore ./batchstore

CMD ["go", "test", "./..."]

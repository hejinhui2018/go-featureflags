FROM golang:1.22
WORKDIR /src
COPY go.mod bodylimit.go bodylimit_test.go BUG_REPRO.md ./
RUN go test ./...
CMD ["go", "test", "./...", "-count=1"]

FROM golang:1.22
WORKDIR /src
COPY go.mod cancel.go cancel_test.go BUG_REPRO.md ./
RUN go test ./...
CMD ["go", "test", "./...", "-count=1"]

FROM golang:1.22

WORKDIR /src
COPY go.mod ./
COPY README.md BUG_REPRO.md ./
COPY cmd ./cmd
COPY internal ./internal

RUN go test ./... && go vet ./... && go build ./...

CMD ["go", "test", "./..."]

FROM golang:1.23-alpine

WORKDIR /src
COPY . .
RUN go test ./... \
    && go vet ./... \
    && go build ./...

CMD ["go", "run", "./cmd/intervalmerge"]

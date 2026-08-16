ARG BASE_IMAGE=golang:1.23-alpine
FROM ${BASE_IMAGE}
WORKDIR /src
COPY . .
RUN go test ./... && go vet ./... && go build ./...
CMD ["go", "test", "./..."]

FROM golang:1.22

WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . .
CMD ["go", "test", "./...", "-count=1"]

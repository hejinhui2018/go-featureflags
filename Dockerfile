FROM golang:1.22
WORKDIR /src
COPY . .
ENV GOTOOLCHAIN=local
CMD ["go", "test", "./...", "-count=1"]

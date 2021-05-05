FROM golang:latest
WORKDIR /usr/app
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
RUN mkdir data
COPY main.go main.go
COPY server server
RUN go build -o sis .

CMD ["./sis"]

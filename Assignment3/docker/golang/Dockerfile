FROM golang:1.15

WORKDIR /go/src/github.com/postgres-go
COPY . .
RUN go get -u github.com/lib/pq

WORKDIR /go/src/yuta/go
ADD ./go .
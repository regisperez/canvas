FROM golang:1.18

ENV GO111MODULE=on

COPY ./ /go/src/canvas

WORKDIR /go/src/canvas

COPY go.mod ./
COPY go.sum ./

RUN go mod download

RUN go build -o canvas

CMD canvas

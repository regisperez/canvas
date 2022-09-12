FROM golang:1.18

ENV GO111MODULE=on

COPY ./ /go/src/canvas

WORKDIR /go/src/canvas

COPY go.mod ./
COPY go.sum ./

RUN go mod download

RUN go get github.com/githubnemo/CompileDaemon

RUN go install github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon --build="go build -o canvas" --command=./canvas

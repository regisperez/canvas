FROM golang:1.18

ENV GO111MODULE=on

COPY ./ /go/src/canvas

WORKDIR /go/src/canvas/code

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . . ./

RUN go build -o /canvas

EXPOSE 8080

CMD [ "/canvas" ]


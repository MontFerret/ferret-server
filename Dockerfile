FROM golang:1.11-alpine

WORKDIR /go/src/github.com/MontFerret/ferret-server

COPY . /go/src/github.com/MontFerret/ferret-server

RUN go build && mv ferret-server /go/bin/

CMD ["ferret-server"]
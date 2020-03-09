FROM golang:1.13

RUN mkdir -p /go/src/github.com/nitrajka/paymentsFutured
WORKDIR /go/src/github.com/nitrajka/paymentsFutured

ADD . /go/src/github.com/nitrajka/paymentsFutured

RUN go get -v
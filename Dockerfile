FROM golang:1.15

WORKDIR /go/src/github.com/factorim/make-env-file

COPY . .

RUN go install
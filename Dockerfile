FROM golang:latest
RUN mkdir /go/src/ingrid
ADD . /go/src/ingrid/
WORKDIR /go/src/ingrid
RUN go build cmd/main.go
CMD ["./main"]
# TODO: Make this example with our best pratices

FROM golang:1.14.4-alpine3.12 AS builder
WORKDIR /go/src/github.com/alexellis/href-counter/
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY main.go aws_storage.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/alexellis/href-counter/app .
CMD ["./app"]  
# TODO: Make this example with our best pratices

FROM golang:1.14.4-alpine3.12 AS builder
WORKDIR /go/src/github.com/dfds/azure-devops-exporter/
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY main.go aws_storage.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/dfds/azure-devops-exporter/app .

COPY ./entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh
ENTRYPOINT /app/entrypoint.sh
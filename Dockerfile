FROM golang:1.24-alpine AS builder
WORKDIR /opt
COPY . .
RUN go build -o /main main.go

FROM alpine:3.17
WORKDIR /opt
COPY --from=builder /main /opt/main
COPY input_files /opt/input_files
RUN mkdir -p /opt/logs
ENTRYPOINT ["/opt/main"]
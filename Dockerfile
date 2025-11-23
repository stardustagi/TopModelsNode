# syntax=docker/dockerfile:1

FROM golang:1.25.0 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o topmodelslogin main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app
COPY --from=builder /app/topmodelslogin /app/topModelsNode
COPY --from=builder /app/docs /app/docs
EXPOSE 8080
CMD ["/app/topModelsNode"]


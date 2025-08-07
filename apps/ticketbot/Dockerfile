FROM golang:1.24.5-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o ticketbot ./apps/ticketbot

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/ticketbot .
CMD ["./ticketbot"]
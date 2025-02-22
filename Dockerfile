FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -o kargo-promotion-check-ext-argo main.go

# Multiphase build
FROM alpine:3.21.3

WORKDIR /app
COPY --from=builder /app/kargo-promotion-check-ext-argo /app/kpcea

CMD ["/app/kpcea"]
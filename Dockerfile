FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -o kargo-promotion-check-ext-argo main.go

# Multiphase build
FROM alpine:3.21.3

ARG ARGOCD_VERSION="v2.14.2"
WORKDIR /app

#RUN apk add --no-cache curl && \
#    curl -sSL -o /usr/local/bin/argocd https://github.com/argoproj/argo-cd/releases/${ARGOCD_VERSION}/download/argocd-linux-amd64 && \
#    chmod +x /usr/local/bin/argocd
COPY --from=builder /app/kargo-promotion-check-ext-argo /app/kpcea

CMD ["/app/kpcea"]
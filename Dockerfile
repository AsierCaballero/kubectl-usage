FROM golang:1.22-alpine AS builder

ARG TARGETOS=linux
ARG TARGETARCH=amd64

WORKDIR /workspace
RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags="-s -w" -o kubectl-usage .

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata
RUN adduser -D -u 1000 appuser

COPY --from=builder /workspace/kubectl-usage /usr/local/bin/kubectl-usage

USER appuser
ENTRYPOINT ["kubectl-usage"]

# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /src

# Cache module downloads
COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION=dev
ARG BUILDDATE=unknown

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -ldflags "-s -w -X main.Version=${VERSION} -X main.BuildDate=${BUILDDATE}" \
    -o /nagiosql .

# Runtime stage — minimal image
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /nagiosql /usr/local/bin/nagiosql

# Configuration is mounted via volume or environment variables.
EXPOSE 8081

ENTRYPOINT ["nagiosql"]
CMD ["serve"]

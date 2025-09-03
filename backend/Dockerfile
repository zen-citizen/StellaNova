# build
FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG COMMIT_SHA=unknown

RUN BUILD_TIME=$(date '+%Y-%m-%d %H:%M:%S UTC') && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags "-w -s -X 'backend/internal/version.commitSHA=${COMMIT_SHA}' -X 'backend/internal/version.buildTime=${BUILD_TIME}'" \
    -o bin/backend \
    ./cmd/server

# run
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

RUN adduser -D -s /bin/sh appuser

WORKDIR /app

COPY --from=builder /app/bin/backend .
COPY --from=builder /app/assets ./assets

RUN chown -R appuser:appuser /app

USER appuser

EXPOSE 8080

CMD ["./backend"]

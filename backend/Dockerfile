# Stage 1: Dependencies (shared)
FROM golang:1.25.6-alpine AS dependencies

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Stage 2: Development
FROM dependencies AS development

RUN go install github.com/air-verse/air@latest

COPY . .

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]

# Stage 3: Builder (for production)
FROM dependencies AS builder

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main .

# Stage 4: Production
FROM alpine:latest AS production

RUN apk --no-cache add ca-certificates

WORKDIR /app

RUN addgroup -S appuser && adduser -S appuser -G appuser

COPY --from=builder /app/main .

RUN chown appuser:appuser /app/main

USER appuser

EXPOSE 8080

CMD ["./main"]

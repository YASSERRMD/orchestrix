FROM golang:1.21-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o /orchestrix ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /orchestrix .
COPY --from=builder /app/migrations ./migrations

ENV SERVER_PORT=8080
ENV DB_PATH=/app/data/orchestrix.db
ENV JWT_SECRET=orchestrix-secret-key-change-in-production

RUN mkdir -p /app/data

EXPOSE 8080

CMD ["./orchestrix"]

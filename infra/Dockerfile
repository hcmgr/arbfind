# --- Build stage ---
FROM golang:1.23.0 AS builder
WORKDIR /app
COPY ./src ./src
COPY go.mod go.sum ./
RUN mkdir build
RUN go build -o build/arb ./src

# --- Runtime stage ---
FROM debian:bookworm-slim
WORKDIR /app
COPY config.json ./
COPY --from=builder /app/build/arb .
ENTRYPOINT ["./arb", "--docker"]


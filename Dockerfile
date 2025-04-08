# FROM golang:1.24.1
# WORKDIR /app
# COPY ./src ./src
# COPY go.mod go.sum ./
# COPY config.json ./
# RUN mkdir build && go build -o build/arb ./src
# WORKDIR /app/build
# ENTRYPOINT ["./arb"]

# --- Build stage ---
FROM golang:1.23.0 AS builder
WORKDIR /app
COPY ./src ./src
COPY go.mod go.sum ./
COPY config.json ./
RUN mkdir build
RUN go build -o build/arb ./src

# --- Runtime stage ---
FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/build/arb .
ENTRYPOINT ["./arb"]


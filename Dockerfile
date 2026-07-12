FROM golang:1.26.5-alpine AS builder
ARG CGO_ENABLED=0
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o discord-bot

FROM alpine:3.23.5
COPY --from=builder /app/discord-bot /discord-bot
ENTRYPOINT ["/discord-bot"]

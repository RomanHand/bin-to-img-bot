FROM golang:1.23.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o tg-bti-bot cmd/btiBot/main.go

RUN ls -l tg-bti-bot

FROM debian:12.7

RUN apt update && apt install ca-certificates nmap -y

RUN mkdir /app/

COPY --from=builder /app/tg-bti-bot /usr/local/bin/tg-bti-bot
COPY --from=builder /app/config.yml /etc/tg-bti-bot/config.yml

RUN mkdir /app/var/ /app/var/tg-bins/ /app/var/tg-imgs/ 
RUN chmod +x /usr/local/bin/tg-bti-bot



CMD ["tg-bti-bot"]
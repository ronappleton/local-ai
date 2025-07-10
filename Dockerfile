# syntax=docker/dockerfile:1

FROM golang:1.24 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o /usr/local/bin/local-ai ./

FROM ubuntu:22.04
RUN apt-get update && apt-get install -y ca-certificates sqlite3 && rm -rf /var/lib/apt/lists/*
WORKDIR /data
COPY --from=build /usr/local/bin/local-ai /usr/local/bin/local-ai
EXPOSE 8081
CMD ["local-ai", "serve"]
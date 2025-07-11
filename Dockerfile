# syntax=docker/dockerfile:1

FROM golang:1.24 AS go-build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o /usr/local/bin/codex ./

FROM node:20 AS client-build
WORKDIR /app
COPY src/client/package*.json ./
RUN npm install
COPY src/client .
RUN npm run build

FROM ubuntu:22.04
RUN apt-get update && apt-get install -y ca-certificates sqlite3 && rm -rf /var/lib/apt/lists/*
WORKDIR /data
COPY --from=go-build /usr/local/bin/codex /usr/local/bin/codex
COPY --from=client-build /app/dist /client
EXPOSE 8081
CMD ["codex", "serve"]
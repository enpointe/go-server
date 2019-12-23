# Dockerfile Web Server
FROM golang:latest as builder

ENV GO111MODULE=on

WORKDIR /config
WORKDIR /public
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make server

# final stage
FROM scratch
COPY --from=builder /app/cmd/server/server /app/
COPY --from=builder /app/cmd/server/public/index.html /public/
COPY --from=builder /app/Configuration.json /config/
EXPOSE 8080
ENTRYPOINT ["/app/server", "-config", "/config/Configuration.json"]
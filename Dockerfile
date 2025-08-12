FROM golang:1.24.2-alpine as build

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -a -o /handle ./cmd/main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=build /handle /handle
COPY .env .env
EXPOSE 8080
ENTRYPOINT ["/handle"]
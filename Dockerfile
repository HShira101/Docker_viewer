# Desarrollo: Air hot reload (usar con docker-compose.dev.yml)
FROM golang:1.25-alpine AS dev
RUN go install github.com/air-verse/air@latest
WORKDIR /app/Front
EXPOSE 10000
CMD ["air"]

# Build producción
FROM golang:1.25-alpine AS build
WORKDIR /app/Front
COPY Front/ .
RUN go build -o /server .

# Producción
FROM alpine:latest AS prod
WORKDIR /app
COPY --from=build /server .
EXPOSE 10000
CMD ["./server"]

# ─── Front desarrollo ────────────────────────────────────────────────────────
FROM golang:1.25-alpine AS front-dev
RUN go install github.com/air-verse/air@latest
WORKDIR /app/Front
COPY Front/go.mod ./
RUN go mod download
EXPOSE 10000
CMD ["air"]

# ─── Back desarrollo ─────────────────────────────────────────────────────────
FROM golang:1.25-alpine AS back-dev
RUN go install github.com/air-verse/air@latest
WORKDIR /app/Back
COPY Back/go.mod ./
EXPOSE 10001
CMD ["sh", "-c", "go mod tidy && air"]

# ─── Build front producción ───────────────────────────────────────────────────
FROM golang:1.25-alpine AS front-build
WORKDIR /app/Front
COPY Front/ .
RUN go build -o /server .

# ─── Build back producción ────────────────────────────────────────────────────
FROM golang:1.25-alpine AS back-build
WORKDIR /app/Back
COPY Back/ .
RUN go mod tidy && go build -o /server .

# ─── Front producción ─────────────────────────────────────────────────────────
FROM alpine:latest AS front-prod
WORKDIR /app
COPY --from=front-build /server .
EXPOSE 10000
CMD ["./server"]

# ─── Back producción ──────────────────────────────────────────────────────────
FROM alpine:latest AS back-prod
WORKDIR /app
COPY --from=back-build /server .
EXPOSE 10001
CMD ["./server"]

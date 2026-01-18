FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o /app/backend ./cmd/backend

FROM scratch AS prod

COPY --from=builder /app/backend /backend

ENTRYPOINT ["/backend"]

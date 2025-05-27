FROM golang:1.23-alpine AS builder

WORKDIR /src
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/people ./cmd/server

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=builder /bin/people .
COPY configs ./configs
COPY .env .

EXPOSE 8000
ENTRYPOINT ["./people"]

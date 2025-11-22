FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod ./

COPY api ./api
COPY cmd ./cmd

RUN go mod download && go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8000

CMD ["./main"]
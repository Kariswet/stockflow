
FROM golang:1.27-alpine AS builder


RUN apk add --no-cache git ca-certificates

WORKDIR /app


COPY source/go.mod source/go.sum ./


ENV GOPROXY=https://proxy.golang.org,direct


RUN go mod download

COPY source/ .

RUN CGO_ENABLED=0 GOOS=linux go build -o stockflow-api ./main.go

FROM alpine:latest

RUN apk add --no-cache tzdata ca-certificates

WORKDIR /root/

COPY --from=builder /app/stockflow-api .

EXPOSE 8080

CMD ["./stockflow-api"]
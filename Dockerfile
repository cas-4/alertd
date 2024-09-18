FROM golang:alpine AS builder

RUN apk --update add ca-certificates git

WORKDIR /app

COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build cmd/main.go

# Run the exe file
FROM scratch

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app .


EXPOSE 8080

ENV BACKEND_URL="http://0.0.0.0:8000/graphql"
ENV ADDRESS="0.0.0.0:8080"
ENV REDIS="redis://default:@localhost:6379/0"
ENV DEBUG=0

CMD ["./main"]

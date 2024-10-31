FROM golang:1.23.1-alpine3.20 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/main /usr/local/bin/
# CMD ["main"]
CMD ["sh", "-c", "sleep 10 && main"]
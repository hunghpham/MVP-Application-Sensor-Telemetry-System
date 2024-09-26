# Build the Go binary, simply just building the Go image
FROM golang:1.17 as builder
WORKDIR /app
COPY . .
RUN go get -u ./...
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /go-webserver

# Deploy the Go binary, copying image to the server location to run it
FROM scratch
COPY --from=builder /go-webserver /go-webserver
CMD ["/go-webserver"]
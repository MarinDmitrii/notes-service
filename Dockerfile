FROM golang:1.21 as builder
WORKDIR /app
COPY . ./
# RUN make build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o notes-service cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/notes-service ./
COPY .env ./
EXPOSE 9090
CMD ["./notes-service"]
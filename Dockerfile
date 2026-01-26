FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o go_app ./cmd

FROM scratch
WORKDIR /app
COPY --from=builder /app/go_app .
EXPOSE 8080
CMD ["./go_app"]

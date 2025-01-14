# Build stage
FROM golang:1.23.3-alpine3.20 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -o app cmd/main.go

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/app /usr/bin/
COPY app.env .
COPY db/migration ./db/migration
COPY templates/ ./templates


EXPOSE 8080

ENTRYPOINT ["app"]
CMD ["api"]
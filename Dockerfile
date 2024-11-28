# Build stage
FROM golang:1.23.3-alpine3.20 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./db/migration

RUN chmod +x start.sh wait-for.sh

EXPOSE 8080
CMD [ "/app/main", "top_phrases/cmd", "api" ]
ENTRYPOINT [ "/app/start.sh" ]
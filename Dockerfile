FROM golang:1.24-alpine

WORKDIR /app

COPY . .

RUN go mod tidy && go build -o ai-reviewer ./webhook-service/cmd/

RUN chmod +x /app/entrypoint.sh

ENTRYPOINT ["/app/entrypoint.sh"]
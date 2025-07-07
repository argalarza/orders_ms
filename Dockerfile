FROM golang:1.23-alpine

WORKDIR /app

# Instala git y otras herramientas necesarias
RUN apk add --no-cache git
COPY go.mod ./
COPY go.sum ./
COPY . ./


RUN go mod download
RUN go build -o main cmd/main.go

EXPOSE 5001

CMD ["/app/main"]

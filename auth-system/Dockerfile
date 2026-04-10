# Dockerfile para aplicacao Go
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Instalar dependencias do sistema
RUN apk add --no-cache git

# Copiar go.mod e go.sum
COPY go.mod go.sum* ./

# Baixar dependencias
RUN go mod download

# Copiar codigo fonte
COPY . .

# Compilar aplicacao
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Instalar certificados SSL
RUN apk --no-cache add ca-certificates tzdata

# Copiar binario compilado
COPY --from=builder /app/main .

# Copiar templates e arquivos estaticos
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

# Expor porta
EXPOSE 8080

# Comando para executar
CMD ["./main"]

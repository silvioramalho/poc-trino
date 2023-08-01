# Etapa de compilação
FROM golang:1.17 AS build

WORKDIR /app

# Copie os arquivos de código fonte para a imagem Docker
COPY . .

# Compile o aplicativo Go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/httpserver/main.go

# Etapa de execução
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copie o binário compilado para a imagem final
COPY --from=build /app/main .

# Exponha a porta na qual seu aplicativo estará ouvindo
EXPOSE 8080

# Execute o binário
CMD ["./main"]
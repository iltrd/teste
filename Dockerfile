# Imagem base
FROM golang:1.17.2-alpine3.14

# Instala dependências
RUN apk add --no-cache git postgresql-client

# Diretório de trabalho
WORKDIR /app

# Copia os arquivos da aplicação para o container
COPY . .

# Compila a aplicação
RUN go build -o app .

# Expõe a porta utilizada pela aplicação
EXPOSE 8080

# Configurações do banco de dados
ENV DB_HOST=localhost
ENV DB_PORT=5432
ENV DB_USER=postgres
ENV DB_PASSWORD=postgres
ENV DB_NAME=postgres

# Comando para iniciar a aplicação
CMD ["./app"]

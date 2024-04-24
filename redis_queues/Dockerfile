# Derleme aşaması
FROM golang:latest AS builder

WORKDIR /app

COPY . .

# Modül gereksinimlerini yükle ve bağımlılıkları indir
RUN go mod download

# Uygulamayı derle
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Çalışma aşaması
FROM alpine:latest

WORKDIR /root/

# Önceki aşamadan derlenen uygulamayı kopyala
COPY --from=builder /app/main .

# Uygulamayı çalıştır
CMD ["./main"]

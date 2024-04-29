# Go'nun resmi base image'ını kullanarak başlayın
FROM golang:1.18 AS builder

# Çalışma dizinini ayarlayın
WORKDIR /app

# Go mod dosyalarını kopyalayın ve bağımlılıkları indirin
COPY go.mod go.sum ./
RUN go mod download

# Kaynak kodunu kopyalayın
COPY . .

# Uygulamayı build edin
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# İkinci stage, uygulamayı çalıştırmak için hafif bir base image kullanın
FROM alpine:latest  
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# İlk stage'den build edilen executable'ı kopyalayın
COPY --from=builder /app/main .

# Uygulamanın çalışacağı portu belirtin
EXPOSE 3000

# Uygulamayı çalıştırın
CMD ["./main"]

# Base image olarak alpine Linux kullanıyoruz
FROM golang:alpine AS builder

# Gerekli paketleri yükle
RUN apk update && apk add --no-cache git

# Çalışma dizinini oluştur
WORKDIR /app

# Modülleri kopyala
COPY go.mod ./
COPY go.sum ./

# Modülleri indir
RUN go mod download

# Kodları kopyala
COPY . .

# Kodu derle
RUN go build -o main .

# Eğer prod ortamında kullanıyorsanız aşağıdaki satırı kullanarak builder adımını atlayabilirsiniz
FROM alpine:latest

# Çalıştırılabilir dosyayı kopyala
 COPY --from=builder /app/main /usr/local/bin/main

# Çalışma dizinini belirt
WORKDIR /usr/local/bin/

# Uygulamayı başlat
CMD ["./main"]

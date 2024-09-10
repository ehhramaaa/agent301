# Gunakan image Go berbasis Alpine
FROM golang:latest AS builder

# Set working directory
WORKDIR /app

# Salin file go mod dan sum
COPY go.mod go.sum ./

# Download dependensi
RUN go mod download

# Salin kode sumber
COPY . .

# Build aplikasi
RUN go build -o main .

# Gunakan image Alpine sebagai base untuk runtime
FROM alpine:latest

# Install dependensi runtime (jika ada)
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Salin binary dari stage builder
COPY --from=builder /app/main .

# Jalankan aplikasi
CMD ["./main"]
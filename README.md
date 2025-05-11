# Botopia - Bot WhatsApp berbasis Command

Botopia adalah bot WhatsApp modular yang berbasis perintah (command) dan dibangun menggunakan Go dengan library whatsmeow.

## Fitur Utama

- Sistem command yang modular dan mudah diperluas
- Sistem routing perintah yang fleksibel dengan dukungan grouping
- Logging yang informatif dalam bahasa Indonesia
- Hot-reload untuk pengembangan yang lebih cepat
- Koneksi WhatsApp yang robust dengan dukungan QR Code

## Persyaratan

- Go 1.18+
- Git

## Instalasi

1. Clone repositori
   ```
   git clone https://github.com/gwenziro/botopia.git
   cd botopia
   ```

2. Install dependensi development
   ```
   make install-dev
   ```

3. Build aplikasi
   ```
   make build
   ```

## Menjalankan Bot

### Mode Produksi


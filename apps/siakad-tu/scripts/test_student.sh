#!/bin/bash

# Navigasi dari folder /scripts ke folder /apps/siakad-tu/src
# Menggunakan dirname $0 agar skrip aman dijalankan dari direktori mana pun
PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "🚀 Menjalankan uji coba student.test..."

# Jalankan pengujian langsung menembak ke folder __tests__
node "$PROJECT_DIR/src/domain/student/__tests__/student.test.js"

# Menangkap status sukses/gagal dari Node.js
if [ $? -eq 0 ]; then
    echo "🎉 Pengujian Berhasil!"
else
    echo "❌ Pengujian Gagal, silakan periksa kode Anda."
    exit 1
fi

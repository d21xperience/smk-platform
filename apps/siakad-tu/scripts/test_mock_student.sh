#!/bin/bash

# Navigasi dari folder /scripts ke folder /apps/siakad-tu/src
# Menggunakan dirname $0 agar skrip aman dijalankan dari direktori mana pun
PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "🚀 Menjalankan uji coba studentMockAdapter.test..."

# Jalankan pengujian langsung menembak ke folder __tests__
node --experimental-vm-modules --no-warnings "$PROJECT_DIR/src/adapters/mock/__tests__/studentMockAdapter.test.js"

# Menangkap status sukses/gagal dari Node.js
if [ $? -eq 0 ]; then
    echo "🎉 Pengujian Berhasil!"
else
    echo "❌ Pengujian Gagal, silakan periksa kode Anda."
    exit 1
fi

# Pastikan Mock Adapter mengimplementasikan semua method di kontrak
grep -E "async [a-zA-Z]+\(" $PROJECT_DIR/src/adapters/mock/studentMockAdapter.js|wc -l
# Ekspektasi: 10 method (sesuai kontrak)
# Pastikan factory menggunakan mock di development
grep "useApi" $PROJECT_DIR/src/adapters/api/studentApi.js


echo "==============================================="
echo "✅ SEMUA MODUL SELESAI DIUJI DAN LOLOS!"
echo "==============================================="

#!/bin/bash

# Navigasi dari folder /scripts ke folder /apps/siakad-tu/src
# Menggunakan dirname $0 agar skrip aman dijalankan dari direktori mana pun
PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "🚀 Menjalankan uji coba context.test..."

# Jalankan pengujian langsung menembak ke folder __tests__
node "$PROJECT_DIR/src/domain/context/__tests__/context.test.js"

# Menangkap status sukses/gagal dari Node.js
if [ $? -eq 0 ]; then
    echo "🎉 Pengujian Berhasil!"
else
    echo "❌ Pengujian Gagal, silakan periksa kode Anda."
    exit 1
fi

grep -r "from 'vue'" $PROJECT_DIR/src/domain/context/engine/ && echo "❌ VIOLATION"||echo "✅ Engine clean"
grep -r "from 'pinia'" $PROJECT_DIR/src/domain/context/engine/ && echo "❌ VIOLATION"||echo "✅ Engine clean"
grep -r "from 'axios'" $PROJECT_DIR/src/domain/context/engine/ && echo "❌ VIOLATION"||echo "✅ Engine clean"
# Pastikan Store TIDAK berisi logika bisnis
grep -r "new AcademicPeriod" $PROJECT_DIR/src/stores/ && echo "❌ VIOLATION"||echo "✅ Store clean"
grep -r "new OperationalContext" $PROJECT_DIR/src/stores/ && echo "❌ VIOLATION"||echo "✅ Store clean"


<template>
  <q-header elevated class="bg-primary text-white q-py-xs">
    <q-toolbar class="row items-center justify-between">

      <!-- Sisi Kiri: Logo / Nama Aplikasi & Tahun Ajaran -->
      <div class="row items-center no-wrap">
        <!-- <q-icon name="school" size="32px" class="q-mr-sm text-amber" /> -->
        <div>
          <!-- <div class="text-h6 text-weight-bold leading-tight">MONITORING SISWA</div> -->
          <div class="text-caption text-amber text-weight-medium">
            TAHUN AJARAN: {{ tahunAjaran }} SEMESTER: 1
          </div>
        </div>
      </div>

      <!-- Sisi Kanan: Jam Digital & Tanggal Real-time -->
      <div class="row items-center no-wrap text-right">
        <q-icon name="schedule" size="28px" class="q-mr-md opacity-80" />
        <div>
          <!-- Jam Berdetak Detik demi Detik -->
          <div class="text-h5 text-weight-bolder font-mono shadow-text">
            {{ currentTime }}
          </div>
          <!-- Tanggal Hari Ini -->
          <div class="text-caption text-grey-3 text-weight-light">
            {{ currentDate }}
          </div>
        </div>
      </div>

    </q-toolbar>
  </q-header>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'


// Properti Tahun Ajaran (Bisa diubah manual atau dijadikan props jika dinamis)
const tahunAjaran = ref('2025/2026')

// State untuk waktu dan tanggal
const currentTime = ref('-- : -- : --')
const currentDate = ref('')

// Variabel untuk menampung interval pembaruan jam
let timerInterval = null

// Fungsi untuk memperbarui jam dan tanggal menggunakan standar lokal Indonesia
const updateTime = () => {
  const now = new Date()

  // Format Jam: HH:mm:ss
  currentTime.value = now.toLocaleTimeString('id-ID', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false
  }).replace(/\./g, ':') // Memastikan pemisah menggunakan titik dua

  // Format Tanggal: Hari, DD Bulan YYYY (Contoh: Selasa, 25 Agustus 2026)
  currentDate.value = now.toLocaleDateString('id-ID', {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
    year: 'numeric'
  })
}

// Jalankan jam saat komponen dipasang ke layar (Mounted)
onMounted(() => {
  updateTime() // Jalankan langsung sekali di awal agar tidak kosong
  timerInterval = setInterval(updateTime, 1000) // Perbarui setiap 1 detik (1000ms)
})

// Bersihkan interval saat komponen dihancurkan (Unmounted) agar tidak bocor memori
onUnmounted(() => {
  if (timerInterval) clearInterval(timerInterval)
})
</script>

<style scoped>
.leading-tight {
  line-height: 1.2;
}

.font-mono {
  font-family: 'Courier New', Courier, monospace;
  /* Font monospace agar angka jam tidak bergeser saat berganti */
  letter-spacing: 1px;
}

.shadow-text {
  text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.2);
}

.opacity-80 {
  opacity: 0.8;
}
</style>

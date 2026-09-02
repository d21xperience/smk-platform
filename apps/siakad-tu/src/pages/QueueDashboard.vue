<template>
  <!-- Kontainer Utama Mengunci Tinggi Sesuai Layar (Viewport) -->
  <q-page class="bg-black text-white column no-wrap window-height q-pa-md">

    <!-- Konten Utama: Dibagi Menjadi Dua Kolom Besar (Kiri & Kanan) -->
    <div class="row col no-wrap q-col-gutter-md min-height-0">

      <!-- ================= SISI KIRI (Video & Nomor Aktif) ================= -->
      <div class="col-8 column no-wrap justify-between">

        <!-- Sisi Kiri Atas: Pemutar Video / Live TV (Ambil 70% Ruang Kiri) -->
        <div class="col-8 bg-grey-10 rounded-borders overflow-hidden relative-position shadow-2">
          <!-- Anda bisa mengganti tag ini dengan <video> atau iframe live streaming -->
          <div class="absolute-full flex flex-center text-grey-6 text-subtitle1">
            <q-icon name="live_tv" size="64px" class="q-mb-sm full-width text-center" />
            [ AREA VIDEO / LIVE STREAMING / PROMOSI ]
          </div>
        </div>

        <!-- Sisi Kiri Bawah: Display Nomor Antrean Utama (Ambil Sisa Ruang Kiri) -->
        <div class="col-4 row q-col-gutter-md q-pt-md">
          <!-- Antrean Customer Service (CS) -->
          <div class="col-6 column">
            <q-card flat bordered class="col bg-grey-9 text-white column justify-center items-center">
              <div class="text-subtitle2 text-amber text-weight-bold">CUSTOMER SERVICE 3</div>
              <div class="text-h2 font-mono text-weight-bolder text-green-4 q-my-xs">824</div>
            </q-card>
          </div>

          <!-- Antrean Teller -->
          <div class="col-6 column">
            <q-card flat bordered class="col bg-grey-9 text-white column justify-center items-center">
              <div class="text-subtitle2 text-amber text-weight-bold">TELLER 1</div>
              <div class="text-h2 font-mono text-weight-bolder text-yellow-4 q-my-xs">054</div>
            </q-card>
          </div>
        </div>

      </div>

      <!-- ================= SISI KANAN (Daftar Antrean Lain) ================= -->
      <div class="col-4 column no-wrap">
        <q-card flat bordered class="col bg-grey-9 column no-wrap">
          <!-- Header Sisi Kanan -->
          <q-card-section class="bg-amber-9 text-dark text-weight-bold text-center q-py-sm">
            <div class="text-subtitle1 text-weight-bolder">DAFTAR ANTREAN / KAUNTER</div>
          </q-card-section>

          <!-- List Riwayat Panggilan Antrean (Dibuat Pas Sesuai Sisa Tinggi) -->
          <q-card-section class="col overflow-hidden q-pa-none">
            <q-list separator class="full-height">
              <q-item v-for="item in daftarAntrean" :key="item.loket" class="q-py-md items-center">
                <q-item-section>
                  <q-item-label class="text-subtitle1 text-weight-bold text-grey-4">{{ item.loket }}</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-item-label class="text-h5 font-mono text-weight-bolder text-amber-4">{{ item.nomor
                    }}</q-item-label>
                </q-item-section>
              </q-item>
            </q-list>
          </q-card-section>
        </q-card>
      </div>

    </div>

    <!-- Running Text di Bagian Paling Bawah (Sama Seperti Struktur Sebelumnya) -->
    <div class="flex-shrink-0 q-mt-md" style="height: 36px;">
      <!-- Anda bisa memanggil komponen RunningTextTicker yang kita buat sebelumnya di sini -->
      <div class="bg-amber-9 text-dark full-height flex items-center q-px-md text-weight-bold rounded-borders">
        <q-icon name="volume_up" size="20px" class="q-mr-sm" />
        PENGUMUMAN: UTAMAKAN ANTRI DENGAN TERTIB. BUDAYAKAN ANTRI UNTUK KENYAMANAN BERSAMA.
      </div>
    </div>

  </q-page>
</template>

<script setup>
import { ref } from 'vue'

// Data simulasi untuk mengisi list antrean di sebelah kanan layar
const daftarAntrean = ref([
  { loket: 'Customer Service 1', nomor: '821' },
  { loket: 'Customer Service 2', nomor: '823' },
  { loket: 'Teller 2', nomor: '051' },
  { loket: 'Teller 3', nomor: '052' },
  { loket: 'Teller 4', nomor: '053' },
  { loket: 'Loket Pengaduan', nomor: 'A-12' }
])
</script>

<style scoped>
.window-height {
  height: 100vh;
  /* Kunci tinggi tepat 100% tinggi layar monitor */
  overflow: hidden;
}

.min-height-0 {
  min-height: 0;
  /* Mencegah elemen flexbox meluap ke bawah */
}

.font-mono {
  font-family: 'Courier New', Courier, monospace;
}
</style>

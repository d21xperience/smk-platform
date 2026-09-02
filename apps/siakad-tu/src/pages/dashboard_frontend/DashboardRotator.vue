<template>
  <!-- Kunci tinggi halaman agar pas satu layar penuh tanpa scroll -->
  <q-page class="column no-wrap bg-grey-2 text-dark">

    <!-- ================= TAMPILAN 1: DASBOR STATISTIK SEKOLAH ================= -->
    <div v-show="activeLayout === 1" class=" animated fadeIn">
      <ApexBarChart title="Jumlah Siswa Masuk & Lulus per Tahun Ajaran" :labels="tahunList" :lulusData="dataLulus"
        :masukData="dataMasuk"/>
    </div>


    <!-- ================= TAMPILAN 2: LAYOUT ALA PERBANKAN (MEDIA & INFORMASI) ================= -->
    <!-- <div v-show="activeLayout === 2"
      class="col row no-wrap q-col-gutter-md q-pa-md bg-black text-white animated fadeIn min-height-0"> -->
    <div v-show="activeLayout === 2" class=" animated fadeIn">
      <LandingKohor1 />
      <!-- Sisi Kiri: Pemutar Video Sekolah / Live TV (Dominan 75%) -->
      <!-- <div class="col-9 column no-wrap">
        <div class="col bg-grey-10 rounded-borders overflow-hidden relative-position shadow-2 video-container">

          <TiktokPlayer :video-ids="['7655271631866088712', '7655267655376522504']" :autoplay="true" :muted="false"
            :estimated-duration="25000" />
        </div>
      </div> -->

      <!-- Sisi Kanan: Pengumuman Penting / Agenda Hari Ini -->
      <!-- <div class=" col-3 column no-wrap">
        <q-card flat bordered class="col bg-grey-9 column no-wrap border-amber">
          <q-card-section class="bg-amber-9 text-dark text-weight-bold text-center q-py-sm">
            <div class="text-subtitle1 text-weight-bolder">AGENDA & INFO PENTING</div>
          </q-card-section>

          <q-card-section class="col overflow-auto q-pa-md">
            <div class="q-mb-md text-amber text-weight-bold">► Ujian Tengah Semester</div>
            <p class="text-caption text-grey-4 q-mb-lg">Pelaksanaan UTS gelombang pertama dimulai senin depan pukul
              07.30 WIB.</p>

            <div class="q-mb-md text-amber text-weight-bold">► Kunjungan Industri</div>
            <p class="text-caption text-grey-4 q-mb-lg">Siswa kelas XI TKJT akan melaksanakan kunjungan industri ke
              PT
              Telekomunikasi pada tanggal 4.</p>
          </q-card-section>
        </q-card>
      </div> -->

    </div>


    <!-- ================= TAMPILAN 3: STRUKTUR ORGANISASI SEKOLAH ================= -->
    <div v-show="activeLayout === 3" class="col column no-wrap q-pa-md animated fadeIn bg-blue-grey-1 min-height-0">
      <!-- Header Internal Layout 3 -->
      <div class="text-h5 text-weight-bold text-center text-primary q-mb-md flex-shrink-0">
        STRUKTUR ORGANISASI MANAGEMENT SEKOLAH
      </div>

      <!-- Wadah Utama Bagan Struktur (Membagi Ruang Secara Fleksibel) -->
      <DashboardStrukturOrganisasi />
    </div>
    <div v-show="activeLayout === 4"
      class="col row no-wrap q-col-gutter-lg q-pa-md animated fadeIn bg-grey-1 min-height-0">
      <DashboardProfilSekolah />
    </div>
  </q-page>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
// import TiktokPlayer from '@/components/TiktokPlayer.vue'
import DashboardProfilSekolah from './DashboardProfilSekolah.vue'
import DashboardStrukturOrganisasi from './DashboardStrukturOrganisasi.vue'
import LandingKohor1 from './LandingKohor1.vue'
import ApexBarChart from '@/components/ApexBarChart.vue'
const tahunList = ref(['2021/2022', '2022/2023', '2023/2024', '2024/2025', '2025/2026', '2026/2027'])
const dataLulus = ref([224, 200, 205, 205, 200, 250])
const dataMasuk = ref([205, 205, 350, 200, 250, 350])
// import 'tiktok-video-element'
// State untuk melacak layout mana yang sedang aktif (1 = Statistik, 2 = Perbankan/Media)
const activeLayout = ref(1)

// Konfigurasi waktu jeda rotasi (Contoh: Berubah setiap 30 detik / 30000 milidetik)
const intervalWaktu = 30000
let rotationTimer = null

// Fungsi untuk mengganti layout secara bergantian
// Perbaikan Logika Rotasi untuk 3 Layout (1 -> 2 -> 3 -> 1)
const startRotation = () => {
  rotationTimer = setInterval(() => {
    if (activeLayout.value === 1) {
      activeLayout.value = 2
    } else if (activeLayout.value === 2) {
      activeLayout.value = 3
    } else if (activeLayout.value === 3) {
      activeLayout.value = 4
    } else {
      activeLayout.value = 1 // Kembali ke layout awal (Statistik)
    }
  }, intervalWaktu)
}

onMounted(() => {
  startRotation()
})

onUnmounted(() => {
  if (rotationTimer) clearInterval(rotationTimer)
})

// Data Struktur Wakasek
// const wakasekList = ref([
//   { nama: 'Ahmad Rifai, S.Pd', jabatan: 'WAKASEK KURIKULUM' },
//   { nama: 'Siti Aminah, M.Pd', jabatan: 'WAKASEK KESISWAAN' },
//   { nama: 'Dedi Sutisna, S.T', jabatan: 'WAKASEK SARPRAS' },
//   { nama: 'Rian Hidayat, M.Kom', jabatan: 'WAKASEK HUBIN / HUMAS' }
// ])
</script>

<style scoped>
.min-height-0 {
  min-height: 0;
}

.flex-shrink-0 {
  flex-shrink: 0;
}

.border-amber {
  border: 1px solid var(--q-warning) !important;
}

/* Efek transisi halus (Fade In) saat layout berganti */
.animated {
  animation-duration: 0.8s;
  animation-fill-mode: both;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: scale(0.99);
  }

  to {
    opacity: 1;
    transform: scale(1);
  }
}

.fadeIn {
  animation-name: fadeIn;
}

.tiktok-video {
  --controls-color: #ff0050;
  /* Contoh: mengubah warna kontrol */
  --background-color: #1a1a2e;
  /* Contoh: mengubah latar belakang */
  --border-radius: 16px;
  /* Contoh: mengubah sudut */
}

.video-player-container {
  border: 2px solid #ff0050;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.video-player {
  border-radius: 16px;
}

.video-container {
  width: 100%;
  height: 100vh;
  /* atau 100% jika di dalam parent yang sudah punya tinggi */
}

.tiktok-player {
  width: 100% !important;
  height: 100% !important;
  display: block !important;
}

/* Styling Khusus Bagan Struktur Organisasi */
.org-container {
  max-width: 1000px;
  margin: 0 auto;
}

.org-card {
  width: 210px;
  display: inline-block;
  border-radius: 8px;
  transition: transform 0.3s;
}

.org-card:hover {
  transform: translateY(-5px);
}
</style>

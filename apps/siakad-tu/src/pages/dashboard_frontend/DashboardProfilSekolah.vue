<template>
  <!-- SISI KIRI: Profil Singkat & Visi Misi Sekolah (Lebar 4 Kolom) -->
  <div class="col-4 column no-wrap justify-between">
    <q-card flat bordered class="col column no-wrap bg-primary text-white shadow-2">
      <q-card-section class="text-center q-pb-none">
        <q-icon name="school" size="56px" class="text-amber q-mb-xs" />
        <div class="text-h5 text-weight-bolder">PROFIL SEKOLAH</div>
        <div class="text-caption text-grey-3">Mewujudkan Generasi Unggul & Berkarakter</div>
      </q-card-section>

      <q-separator dark inset class="q-my-md" />

      <q-card-section class="col overflow-auto q-pt-none text-body2 text-justify">
        <div class="text-subtitle2 text-weight-bold text-amber q-mb-xs">VISI:</div>
        <p class="q-mb-md text-grey-2">Menjadi lembaga pendidikan kejuruan yang menghasilkan lulusan beriman, bertakwa,
          kompeten, dan siap bersaing di dunia industri global.</p>

        <div class="text-subtitle2 text-weight-bold text-amber q-mb-xs">MISI UTAMA:</div>
        <ul class="q-pl-md text-grey-2" style="margin: 0; padding-left: 20px;">
          <li class="q-mb-xs">Menyelenggarakan pembelajaran berbasis industri dan teknologi.</li>
          <li class="q-mb-xs">Mengembangkan kemitraan strategis dengan dunia usaha/kerja (DUDI).</li>
          <li class="q-mb-xs">Menanamkan nilai karakter, disiplin, dan jiwa wirausaha.</li>
        </ul>
      </q-card-section>
    </q-card>
  </div>

  <!-- SISI KANAN: Daftar 5 Kompetensi Keahlian / Jurusan (Lebar 8 Kolom) -->
  <div class="col-8 column no-wrap">
    <div class="text-h5 text-weight-bold text-primary q-mb-md flex-shrink-0 row items-center">
      <q-icon name="stars" class="q-mr-sm" />
      KOMPETENSI KEAHLIAN (JURUSAN)
    </div>

    <!-- Grid 5 Jurusan Flexbox Berjejer Rapi -->
    <div class="col row q-col-gutter-md min-height-0 overflow-hidden">
      <div class="col-4 column" v-for="jurusan in daftarJurusan" :key="jurusan.kode">
        <q-card flat bordered class="col column no-wrap jurusan-card shadow-1">
          <!-- Header Kartu Warna Jurusan -->
          <div :class="`bg-${jurusan.color} text-white q-pa-md text-center column items-center`">
            <q-avatar size="44px" class="bg-white-top q-mb-xs">
              <q-icon :name="jurusan.icon" :class="`text-${jurusan.color}`" size="24px" />
            </q-avatar>
            <div class="text-h6 text-weight-bolder leading-none">{{ jurusan.kode }}</div>
          </div>

          <!-- Deskripsi Kompetensi Jurusan -->
          <q-card-section class="col overflow-auto q-pa-sm text-center flex flex-center">
            <div>
              <div class="text-subtitle2 text-weight-bold text-grey-9 q-mb-xs">{{ jurusan.nama }}</div>
              <div class="text-caption text-grey-6 text-weight-light">{{ jurusan.deskripsi }}</div>
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

// const activeLayout = ref(1)
// const intervalWaktu = 30000 // Jeda pergantian layout (30 detik)
// let rotationTimer = null

// Data 5 Jurusan Sekolah Sesuai Request Anda
const daftarJurusan = ref([
  {
    kode: 'TKR',
    nama: 'Teknik Kendaraan Ringan',
    icon: 'directions_car',
    color: 'red-8',
    deskripsi: 'Fokus pada perawatan, perbaikan, dan diagnosis teknologi mesin otomotif mobil.'
  },
  {
    kode: 'TSM',
    nama: 'Teknik Sepeda Motor',
    icon: 'two_wheeler',
    color: 'orange-9',
    deskripsi: 'Ahli dalam pemeliharaan mekanis, kelistrikan, dan sistem injeksi sepeda motor roda dua.'
  },
  {
    kode: 'TKJ',
    nama: 'Teknik Komputer & Jaringan',
    icon: 'lan',
    color: 'indigo-8',
    deskripsi: 'Menguasai infrastruktur jaringan internet, perakitan komputer, server, dan cyber security.'
  },
  {
    kode: 'OTKP',
    nama: 'Otomatisasi & Tata Kelola Perkantoran',
    icon: 'business_center',
    color: 'teal-8',
    deskripsi: 'Kompeten dalam manajemen dokumen, korespondensi bisnis, kearsipan, dan administrasi digital.'
  },
  {
    kode: 'AKL',
    nama: 'Akuntansi & Keuangan Lembaga',
    icon: 'payments',
    color: 'green-8',
    deskripsi: 'Mahir dalam mengelola siklus keuangan pembukuan, perpajakan, dan audit perbankan/perusahaan.'
  }
])

// Pembaruan Logika Rotasi Mutakhir untuk 4 Layout Berantai
// const startRotation = () => {
//   rotationTimer = setInterval(() => {
//     if (activeLayout.value === 1) {
//       activeLayout.value = 2
//     } else if (activeLayout.value === 2) {
//       activeLayout.value = 3
//     } else if (activeLayout.value === 3) {
//       activeLayout.value = 4
//     } else {
//       activeLayout.value = 1 // Kembali ke layout awal (Statistik)
//     }
//   }, intervalWaktu)
// }

// onMounted(() => {
//   startRotation()
// })

// onUnmounted(() => {
//   if (rotationTimer) clearInterval(rotationTimer)
// })
</script>

<style scoped>
.min-height-0 {
  min-height: 0;
}

.flex-shrink-0 {
  flex-shrink: 0;
}

.leading-none {
  line-height: 1;
}

.bg-white-top {
  background: rgba(255, 255, 255, 0.95);
}

/* Animasi Hover Halus pada Kartu Jurusan */
.jurusan-card {
  height: 100%;
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
}

.jurusan-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.15) !important;
}

/* Animasi FadeIn Transisi */
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
</style>

<template>
  <q-page class="bg-grey-1 q-py-xl">
    <div class="container q-px-md max-width-center" v-if="jurusan">

      <!-- NAVIGATION BREADCRUMB -->
      <q-breadcrumbs class="text-grey q-mb-lg text-caption text-sm-body2">
        <q-breadcrumbs-el label="Beranda" icon="home" to="/" />
        <q-breadcrumbs-el label="Jurusan" icon="widgets" />
        <q-breadcrumbs-el :label="jurusan.title" />
      </q-breadcrumbs>

      <div class="row q-col-gutter-lg">

        <!-- KOLOM KIRI: VISUAL & INFORMASI INTI (7/12 KELOMPOK DESKTOP) -->
        <div class="col-12 col-md-8">
          <q-card flat bordered class="bg-white rounded-borders shadow-1">
            <!-- Gambar Utama / Cover Banner Jurusan -->
            <q-img :src="jurusan.cover" :ratio="16 / 9" class="rounded-top">
              <div class="absolute-bottom-left q-ma-md q-pa-none bg-transparent">
                <q-badge :color="jurusan.themeColor" class="text-weight-bold q-px-md q-py-xs text-subtitle2">
                  <q-icon :name="jurusan.icon" class="q-mr-xs" /> {{ jurusan.badge }}
                </q-badge>
              </div>
            </q-img>

            <q-card-section class="q-pa-lg">
              <h1 class="text-h4 text-weight-bold text-grey-9 q-my-none text-sm-h3 leading-tight">
                {{ jurusantitleLong || jurusan.title }}
              </h1>

              <div class="row items-center q-mt-md text-grey-6 text-caption text-sm-body2">
                <q-icon name="verified" :color="jurusan.themeColor" size="xs" class="q-mr-xs" />
                <span class="text-weight-bold">Status Akreditasi: </span>
                <q-badge color="green-1" text-color="green-7" class="q-ml-sm text-weight-bold">A (Sangat Baik)</q-badge>
              </div>

              <q-separator class="q-my-lg" />

              <!-- DESKRIPSI JURUSAN -->
              <div class="text-h6 text-weight-bold text-primary q-mb-sm">Tentang Kompetensi Keahlian</div>
              <p class="text-body1 text-grey-8 line-height-relaxed" v-html="jurusan.description"></p>

              <!-- MATERI / KURIKULUM YANG DIPELAJARI -->
              <div class="text-h6 text-weight-bold text-primary q-mt-xl q-mb-md">Kompetensi Inti Materi</div>
              <div class="row q-col-gutter-sm">
                <div v-for="(materi, index) in jurusan.curriculum" :key="index" class="col-12 col-sm-6">
                  <q-item dense class="bg-grey-2 rounded-borders q-py-sm border-light">
                    <q-item-section avatar min-width="auto" class="q-pr-sm">
                      <q-icon name="check" :color="jurusan.themeColor" size="xs" />
                    </q-item-section>
                    <q-item-section class="text-body2 text-grey-9 text-weight-medium">
                      {{ materi }}
                    </q-item-section>
                  </q-item>
                </div>
              </div>
            </q-card-section>
          </q-card>
        </div>

        <!-- KOLOM KANAN: MITRA INDUSTRI & PROSPEK (4/12 KELOMPOK DESKTOP) -->
        <div class="col-12 col-md-4">
          <!-- CARD MITRA & KERJASAMA -->
          <q-card flat bordered class="bg-white rounded-borders shadow-1 q-mb-lg">
            <q-card-section class="bg-primary text-white q-py-sm text-weight-bold text-subtitle1 row items-center">
              <q-icon name="handshake" class="q-mr-sm" /> Sinkronisasi Industri
            </q-card-section>
            <q-card-section class="q-pa-md text-center">
              <div class="text-caption text-grey-7 q-mb-md">Lulusan jurusan ini terhubung dan diserap langsung oleh
                jaringan mitra:</div>
              <div class="row q-col-gutter-md justify-center items-center">
                <div v-for="(partner, i) in jurusan.partners" :key="i" class="col-6">
                  <div
                    class="partner-box border-light rounded-borders q-pa-sm text-weight-bold text-grey-8 text-caption">
                    {{ partner }}
                  </div>
                </div>
              </div>
            </q-card-section>
          </q-card>

          <!-- CARD PELUANG KARIER LULUSAN -->
          <q-card flat bordered class="bg-white rounded-borders shadow-1 q-mb-lg">
            <q-card-section class="bg-dark text-white q-py-sm text-weight-bold text-subtitle1 row items-center">
              <q-icon name="work" class="q-mr-sm" /> Prospek Kerja Lulusan
            </q-card-section>
            <q-list bordered separator class="rounded-borders">
              <q-item v-for="(job, idx) in jurusan.careerProspects" :key="idx">
                <q-item-section avatar min-width="auto" class="q-pr-sm">
                  <q-icon name="arrow_right" color="grey-7" />
                </q-item-section>
                <q-item-section class="text-body2 text-grey-9 text-weight-medium">{{ job }}</q-item-section>
              </q-item>
            </q-list>
          </q-card>

          <!-- ADMISSION CTA DI HP / PC -->
          <q-card
            :class="`bg-${jurusan.themeColor} text-white rounded-borders shadow-3 text-center q-pa-md shadow - 3`">
            <div class="text-subtitle1 text-weight-bold">Tertarik dengan Jurusan Ini?</div>
            <div class="text-caption text-white-8 q-mt-xs">Kuota terbatas untuk tiap gelombang pendaftaran PPDB.</div>
            <q-btn unelevated color="white" :text-color="jurusan.themeColor" label="Daftar PPDB Sekarang"
              class="full-width text-weight-bold q-mt-md" icon="assignment" @click="bukaPendaftaran" />
          </q-card>
        </div>

      </div>
    </div>

    <!-- TAMPILAN JIKA PARAMETER JURUSAN TIDAK VALID -->
    <div v-else class="text-center q-py-xl text-grey-7">
      <q-icon name="warning" size="xl" color="negative" />
      <div class="text-h6 q-mt-sm">Halaman Jurusan Tidak Ditemukan</div>
      <q-btn color="primary" label="Kembali ke Beranda" class="q-mt-md" to="/" />
    </div>
  </q-page>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const jurusan = ref(null)

// DATABASE LOKAL KHUSUS DATA DETAIL JURUSAN (TBSM & TKJ)
const masterJurusanData = {
  tbsm: {
    title: 'TBSM',
    titleLong: 'Teknik & Bisnis Sepeda Motor (TBSM)',
    badge: 'Kombinasi Industri Astra',
    icon: 'motorcycle',
    themeColor: 'red-7',
    cover: 'https://unsplash.com',
    description: 'Program keahlian Teknik & Bisnis Sepeda Motor (TBSM) membekali siswa dengan keterampilan teknis mutakhir dalam perawatan, perbaikan, diagnostik malfungsi mekanikal maupun sistem elektrikal sepeda motor. Bekerja sama erat dengan <b>PT Astra Honda Motor (AHM)</b>, kurikulum kami diselaraskan penuh dengan standar kerja teknisi bengkel resmi AHASS.',
    curriculum: [
      'Perawatan Berkala Mesin Motor',
      'Sistem Injeksi PGM-FI Honda',
      'Sasis & Sistem Suspensi Modern',
      'Sistem Kelistrikan & Penerangan',
      'Manajemen Bisnis Bengkel',
      'Diagnosis Kerusakan Berbasis Scanner komputer'
    ],
    partners: ['PT Astra Honda Motor', 'Jaringan AHASS', 'Federal Oil', 'Astra Otoparts'],
    careerProspects: [
      'Mekanik Profesional Bengkel Resmi',
      'Service Advisor (SA)',
      'Wirausaha Bengkel Mandiri',
      'Teknisi Assembly Pabrik Otomotif',
      'Staff Gudang Suku Cadang'
    ]
  },
  tkj: {
    title: 'TKJ',
    titleLong: 'Teknik Komputer & Jaringan (TKJ)',
    badge: 'Cisco & MikroTik Academy',
    icon: 'computer',
    themeColor: 'blue-7',
    cover: 'https://unsplash.com',
    description: 'Jurusan Teknik Komputer & Jaringan (TKJ) memfokuskan siswa pada penguasaan instalasi perangkat keras, sistem operasi komputer, perancangan jaringan kabel dan nirkabel, administrasi server Linux/Windows, serta konsep dasar keamanan siber siber berskala enterprise.',
    curriculum: [
      'Arsitektur Jaringan (LAN/WAN)',
      'Konfigurasi Router MikroTik & Cisco',
      'Administrasi Sistem Server Cloud',
      'Cybersecurity & Network Monitoring',
      'Perakitan & Troubleshooting PC',
      'Manajemen Fiber Optic (FO)'
    ],
    partners: ['MikroTik Academy', 'PT Telkom Indonesia', 'Biznet Networks', 'Lokal ISP'],
    careerProspects: [
      'Network Administrator',
      'IT Support & Helpdesk',
      'System Administrator (SysAdmin)',
      'Teknisi Infrastruktur Jaringan',
      'Wirausaha Jaringan (RT/RW Net)'
    ]
  }
}

// Fungsi memuat data berdasarkan id parameter router (:id)
const loadJurusanData = () => {
  const jurusanId = route.params.id ? route.params.id.toLowerCase() : ''
  jurusan.value = masterJurusanData[jurusanId] || null
}

onMounted(() => {
  loadJurusanData()
})

// Deteksi perubahan parameter URL seandainya pindah jurusan langsung via link
watch(() => route.params.id, () => {
  loadJurusanData()
})

function bukaPendaftaran() {
  window.open('https://sekolah.sch.id', '_blank')
}
</script>

<style scoped>
.max-width-center {
  max-width: 1100px;
  margin: 0 auto;
}

.line-height-relaxed {
  line-height: 1.7;
}

.border-light {
  border: 1px solid #e0e0e0;
}

.partner-box {
  background-color: #f9f9f9;
  min-height: 55px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.white-8 {
  color: rgba(255, 255, 255, 0.8) !important;
}

.rounded-top {
  border-top-left-radius: 8px;
  border-top-right-radius: 8px;
}
</style>

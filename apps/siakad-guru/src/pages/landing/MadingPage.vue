<template>
  <q-page class="bg-grey-1 q-py-xl">
    <div class="container q-px-md max-width-center">
      <!-- 1. HEADER HALAMAN MADING -->
      <div class="text-center q-mb-xl">
        <div class="text-h4 text-sm-h3 text-weight-bold text-primary">
          Mading Digital <span class="text-amber">Vokasi</span>
        </div>
        <p class="text-body1 text-grey-7 q-mt-sm">
          Pusat informasi, kegiatan OSIS, kreativitas siswa, dan pengumuman resmi jurusan.
        </p>
      </div>

      <!-- 2. FILTER KATEGORI (RESPONSIF) -->
      <!-- Menggunakan QTabs dengan scrollable agar pas di layar HP tanpa merusak layout -->
      <div class="row justify-center q-mb-xl">
        <q-tabs
          v-model="selectedCategory"
          dense
          class="text-grey-7 bg-white shadow-1 rounded-borders q-pa-xs"
          active-color="primary"
          indicator-color="primary"
          align="center"
          narrow-indicator
        >
          <q-tab name="all" icon="dashboard" label="Semua Info" />
          <q-tab name="osis" icon="emoji_events" label="Kegiatan OSIS" />
          <q-tab name="tbsm" icon="motorcycle" label="Seputar TBSM" />
          <q-tab name="ppdb" icon="campaign" label="Info PPDB" />
        </q-tabs>
      </div>

      <!-- 3. GRID UTAMA POSTER MADING -->
      <!-- Animasi transisi halus saat kategori poster difilter -->
      <div class="row q-col-gutter-lg justify-center">
        <div
          v-for="(poster, index) in filteredPosters"
          :key="index"
          class="col-12 col-sm-6 col-md-4"
        >
          <q-card class="mading-card shadow-3 flex flex-column no-wrap" flat bordered>
            <!-- Komponen Gambar dengan Fitur Zoom & Lazy Load -->
            <q-img
              :src="poster.image"
              :ratio="3 / 4"
              class="cursor-pointer"
              @click="bukaPratinjauPoster(poster)"
            >
              <!-- Badge Kategori Pojok Kiri Atas -->
              <div class="absolute-top-left q-ma-sm q-pa-none bg-transparent">
                <q-badge
                  :color="getBadgeColor(poster.category)"
                  class="text-weight-bold q-px-sm q-py-xs"
                >
                  {{ poster.categoryName }}
                </q-badge>
              </div>

              <!-- Efek Hover Text saat kursor berada di atas poster -->
              <div class="absolute-full flex flex-center poster-overlay text-center q-px-md">
                <q-btn round color="white" text-color="black" icon="zoom_in" size="md" />
                <div class="text-subtitle2 q-mt-sm text-white text-weight-medium">
                  Klik Untuk Memperbesar
                </div>
              </div>
            </q-img>

            <!-- Deskripsi & Info Pembuat Poster -->
            <q-card-section class="col">
              <div class="text-subtitle1 text-weight-bold text-grey-9 lines-2">
                {{ poster.title }}
              </div>
              <div class="row items-center justify-between q-mt-md text-caption text-grey-6">
                <div class="row items-center">
                  <q-icon name="person" class="q-mr-xs" />
                  <span>Oleh: {{ poster.author }}</span>
                </div>
                <div class="row items-center">
                  <q-icon name="event" class="q-mr-xs" />
                  <span>{{ poster.date }}</span>
                </div>
              </div>
            </q-card-section>

            <q-separator />

            <!-- 4. BUTTON CTA INTERAKTIF (KELEBIHAN WEB DIBANDING MEDIA SOSIAL) -->
            <q-card-actions class="bg-grey-5">
              <q-btn
                flat
                class="full-width text-weight-bold"
                :color="getBadgeColor(poster.category)"
                :label="poster.actionLabel"
                :icon="poster.actionIcon"
                @click="eksekusiAksiPoster(poster)"
              />
            </q-card-actions>
          </q-card>
        </div>

        <!-- TAMPILAN JIKA DATA KATEGORI KOSONG -->
        <div v-if="filteredPosters.length === 0" class="col-12 text-center q-py-xl text-grey-6">
          <q-icon name="layers_clear" size="xl" />
          <div class="text-h6 q-mt-sm">Belum ada poster di kategori ini.</div>
        </div>
      </div>
    </div>

    <!-- 5. DIALOG POP-UP UNTUK MEMPERBESAR GAMBAR POSTER (LIGHTBOX) -->
    <q-dialog v-model="dialogOpen">
      <q-card style="width: 90vw; max-width: 500px" class="bg-transparent no-shadow">
        <q-card-section class="row items-center justify-end q-pa-none q-mb-xs">
          <q-btn icon="close" flat round dense v-close-popup color="white" />
        </q-card-section>
        <q-img :src="activePoster.image" class="rounded-borders" />
        <div class="bg-white q-pa-md rounded-borders q-mt-sm">
          <div class="text-h6 text-weight-bold">{{ activePoster.title }}</div>
          <div class="text-body2 text-grey-7 q-mt-xs">
            Diterbitkan oleh {{ activePoster.author }}
          </div>
        </div>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const selectedCategory = ref('all')
const dialogOpen = ref(false)
const activePoster = ref({})

// Data Dummy Poster Mading Kreatif
const posters = ref([
  {
    title: 'Penerimaan Peserta Didik Baru (PPDB) Kompetensi Keahlian TBSM 2026/2027',
    image: 'https://unsplash.com',
    category: 'ppdb',
    categoryName: 'PPDB 2026',
    author: 'Panitia PPDB',
    date: '28 Mei 2026',
    actionLabel: 'Daftar TBSM Sekarang',
    actionIcon: 'assignment_turned_in',
    actionType: 'link_eksternal',
    actionUrl: 'https://sekolah.sch.id',
  },
  {
    title: 'Servis Sepeda Motor Gratis & Ganti Oli Menyambut Hari Pendidikan Nasional',
    image: 'https://unsplash.com',
    category: 'tbsm',
    categoryName: 'Seputar TBSM',
    author: 'Kaprog TBSM',
    date: '20 Mei 2026',
    actionLabel: 'Isi Formulir Booking Servis',
    actionIcon: 'build',
    actionType: 'internal_route',
    actionUrl: '/faq', // Bisa disesuaikan
  },
  {
    title: 'Turnamen Futsal Antar Kelas "OSIS CUP VI" - Junjung Tinggi Sportivitas!',
    image: 'https://unsplash.com',
    category: 'osis',
    categoryName: 'Kegiatan OSIS',
    author: 'Pengurus OSIS',
    date: '15 Mei 2026',
    actionLabel: 'Lihat Jadwal & Bagan Pertandingan',
    actionIcon: 'emoji_events',
    actionType: 'dialog_info',
    actionUrl: '',
  },
])

// Filter Logika Poster berdasarkan Tab yang dipilih
const filteredPosters = computed(() => {
  if (selectedCategory.value === 'all') return posters.value
  return posters.value.filter((p) => p.category === selectedCategory.value)
})

// Fungsi Penentuan Warna Identitas Badge Kontributor
const getBadgeColor = (category) => {
  if (category === 'osis') return 'blue-7'
  if (category === 'tbsm') return 'red-7'
  if (category === 'ppdb') return 'amber-9'
  return 'grey-7'
}

// Membuka Fitur Zoom Gambar Poster
function bukaPratinjauPoster(poster) {
  activePoster.value = poster
  dialogOpen.value = true
}

// Logika Penanganan Klik Tombol CTA Poster yang Dinamis
function eksekusiAksiPoster(poster) {
  if (poster.actionType === 'link_eksternal') {
    window.open(poster.actionUrl, '_blank')
  } else if (poster.actionType === 'internal_route') {
    router.push(poster.actionUrl)
  } else {
    bukaPratinjauPoster(poster)
  }
}
</script>

<style scoped>
.max-width-center {
  max-width: 1100px;
  margin: 0 auto;
}

.mading-card {
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
  border-radius: 12px;
  overflow: hidden;
  height: 100%;
}

.mading-card:hover {
  transform: translateY(-6px);
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.15) !important;
}

.lines-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* Efek Overlay Zoom yang Menawan saat Poster di-Hover */
.poster-overlay {
  background: rgba(0, 0, 0, 0.6);
  opacity: 0;
  transition: opacity 0.3s ease;
  flex-direction: column;
}

.q-img:hover .poster-overlay {
  opacity: 1;
}
</style>

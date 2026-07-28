<template>
  <q-page class="bg-grey-1 q-py-xl">
    <div class="container q-px-md max-width-center">

      <!-- 1. HEADER HALAMAN -->
      <div class="text-center q-mb-xl">
        <div class="text-h4 text-sm-h3 text-weight-bold text-primary">
          Berita & <span class="text-amber">Prestasi</span> Terbaru
        </div>
        <p class="text-body1 text-grey-7 q-mt-sm">
          Ikuti rekam jejak aktivitas, prestasi membanggakan, dan agenda resmi SMK Pasundan Jatinangor.
        </p>
      </div>

      <!-- 2. KONTROL PENCARIAN & FILTER KATEGORI (RESPONSIF) -->
      <div class="row q-col-gutter-md justify-between items-center q-mb-xl">
        <!-- Filter Kategori dengan QTabs (Bisa Digeser di HP) -->
        <div class="col-12 col-md-7 row justify-center justify-md-start">
          <q-tabs v-model="selectedFilter" dense class="text-grey-7 bg-white shadow-1 rounded-borders q-pa-xs"
            active-color="primary" indicator-color="primary" narrow-indicator>
            <q-tab name="all" icon="newspaper" label="Semua Kabar" />
            <q-tab name="Kegiatan" icon="rocket_launch" label="Kegiatan Sekolah" />
            <q-tab name="Prestasi" icon="emoji_events" label="Prestasi Siswa" />
          </q-tabs>
        </div>

        <!-- Input Pencarian Berita -->
        <div class="col-12 col-md-4">
          <q-input outlined dense v-model="searchQuery" placeholder="Cari berita atau prestasi..."
            class="bg-white shadow-1 rounded-borders" clearable>
            <template v-slot:prepend>
              <q-icon name="search" color="primary" />
            </template>
          </q-input>
        </div>
      </div>

      <!-- 3. GRID UTAMA DAFTAR BERITA -->
      <div class="row q-col-gutter-lg">
        <div v-for="(news, index) in filteredNews" :key="index" class="col-12 col-sm-6 col-md-4">
          <q-card class="news-public-card shadow-1 flat bordered flex flex-column justify-between bg-white">
            <div>
              <!-- Foto Cover dengan Badge Kategori -->
              <q-img :src="news.cover" :ratio="16 / 10">
                <div class="absolute-bottom-left q-ma-sm q-pa-none bg-transparent">
                  <q-badge :color="news.type === 'Prestasi' ? 'positive' : 'info'"
                    class="text-weight-bold q-px-sm q-py-xs">
                    {{ news.type }}
                  </q-badge>
                </div>
              </q-img>

              <!-- Konten Teks -->
              <q-card-section class="q-pa-md">
                <div class="text-caption text-grey-6 row items-center q-mb-xs">
                  <q-icon name="event" size="xs" class="q-mr-xs" /> {{ news.date }}
                  <q-space />
                  <q-icon name="person" size="xs" class="q-mr-xs" /> {{ news.author }}
                </div>

                <div class="text-subtitle1 text-weight-bold text-grey-9 lines-2 title-link q-mt-xs"
                  @click="bukaDetailBerita(news)">
                  {{ news.title }}
                </div>

                <p class="text-body2 text-grey-7 q-mt-sm lines-3 line-height-relaxed">
                  {{ news.excerpt }}
                </p>
              </q-card-section>
            </div>

            <!-- Tombol Aksi -->
            <div>
              <q-separator />
              <q-card-actions align="right" class="q-px-md q-py-sm bg-grey-1">
                <q-btn flat color="primary" label="Baca Selengkapnya" icon-right="chevron_right"
                  class="text-weight-bold text-caption" @click="bukaDetailBerita(news)" />
              </q-card-actions>
            </div>
          </q-card>
        </div>

        <!-- TAMPILAN JIKA BERITA TIDAK DITEMUKAN -->
        <div v-if="filteredNews.length === 0" class="col-12 text-center q-py-xl text-grey-6">
          <q-icon name="find_in_page" size="xl" color="grey-4" />
          <div class="text-h6 q-mt-sm">Kabar berita tidak ditemukan</div>
          <p class="text-caption">Coba gunakan kata kunci pencarian atau filter kategori lainnya.</p>
        </div>
      </div>

    </div>

    <!-- 4. DIALOG POP-UP UNTUK MEMBACA DETAIL BERITA (MODAL READ) -->
    <q-dialog v-model="detailDialogOpen" scrollable max-width="700px">
      <q-card style="width: 100%; max-width: 700px;" class="bg-white rounded-borders">
        <q-card-section class="row items-center justify-between bg-primary text-white q-py-sm">
          <div class="text-subtitle2 text-weight-bold uppercase-text">Detail Publikasi</div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-none">
          <q-img :src="activeNews.cover" :ratio="16 / 9" />
        </q-card-section>

        <q-card-section class="q-pa-lg scroll" style="max-height: 50vh">
          <q-badge :color="activeNews.type === 'Prestasi' ? 'positive' : 'info'" class="text-weight-bold q-mb-sm">
            {{ activeNews.type }}
          </q-badge>

          <h2 class="text-h5 text-weight-bold text-grey-9 q-my-none leading-tight">
            {{ activeNews.title }}
          </h2>

          <div class="text-caption text-grey-6 row items-center q-mt-md q-mb-lg">
            <q-icon name="event" size="xs" class="q-mr-xs" /> Dipublikasikan pada: {{ activeNews.date }}
            <q-space />
            <q-icon name="person" size="xs" class="q-mr-xs" /> Kontributor: {{ activeNews.author }}
          </div>

          <q-separator class="q-my-md" />

          <!-- Isi Konten Artikel Lengkap -->
          <div class="text-body1 text-grey-8 line-height-relaxed" v-html="activeNews.content"></div>
        </q-card-section>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed } from 'vue'

const searchQuery = ref('')
const selectedFilter = ref('all')
const detailDialogOpen = ref(false)
const activeNews = ref({})

// Data Master Berita Sekolah (Sinkron dengan inputan Panel Admin Guru)
const newsData = ref([
  {
    title: 'Siswa TBSM SMK Pasundan Jatinangor Raih Juara 1 Kompetensi Mekanik Astra Honda Tingkat Provinsi',
    excerpt: 'Prestasi gemilang kembali ditorehkan oleh siswa jurusan TBSM dalam ajang Honda Skill Contest 2026. Kompetisi ketat menguji ketepatan troubleshooting mesin injeksi terbaru.',
    content: 'Prestasi gemilang kembali ditorehkan oleh siswa jurusan TBSM dalam ajang Honda Skill Contest 2026. Kompetisi ketat menguji ketepatan troubleshooting mesin injeksi terbaru.<br><br>Kepala Program TBSM menyampaikan rasa bangganya atas kedisiplinan siswa dalam berlatih di bengkel standar AHASS milik sekolah. Juara ini sekaligus mengamankan tiket menuju kompetisi mekanik tingkat nasional mewakili provinsi.',
    cover: 'https://unsplash.com',
    type: 'Prestasi',
    date: '30 Mei 2026',
    author: 'Kaprog TBSM'
  },
  {
    title: 'Pelaksanaan Servis Motor Gratis oleh Komunitas Siswa TBSM bagi Warga Sekitar Sekolah',
    excerpt: 'Sebagai bentuk bakti sosial, jurusan TBSM mendirikan posko servis gratis. Kegiatan ini melatih kematangan mental siswa dalam menghadapi konsumen riil di lapangan.',
    content: 'Sebagai bentuk bakti sosial, jurusan TBSM mendirikan posko servis gratis. Kegiatan ini melatih kematangan mental siswa dalam menghadapi konsumen riil di lapangan.<br><br>Lebih dari 50 motor warga sekitar berhasil diservis oleh para siswa kelas XII yang didampingi oleh instruktur mekanik berpengalaman. Kegiatan ini direncanakan menjadi agenda rutin tahunan sekolah.',
    cover: 'https://unsplash.com',
    type: 'Kegiatan',
    date: '25 Mei 2026',
    author: 'Humas Hub'
  },
  {
    title: 'Kunjungan Industri Semester Genap Kelas XI ke Plant Perakitan Utama PT Astra Honda Motor',
    excerpt: 'Siswa diajak melihat langsung alur produksi perakitan motor berskala robotik modern, guna menyelaraskan wawasan budaya kerja industri langsung sejak bangku sekolah.',
    content: 'Siswa diajak melihat langsung alur produksi perakitan motor berskala robotik modern, guna menyelaraskan wawasan budaya kerja industri langsung sejak bangku sekolah.<br><br>Melalui kunjungan ini, diharapkan siswa memiliki gambaran nyata mengenai kedisiplinan, keselamatan kerja (K3), serta etos kerja tinggi yang diterapkan di dunia industri berskala internasional.',
    cover: 'https://unsplash.com',
    type: 'Kegiatan',
    date: '18 Mei 2026',
    author: 'Kesiswaan'
  }
])

// Logika Pemrosesan Filter Tab dan Pencarian Teks
const filteredNews = computed(() => {
  return newsData.value.filter(news => {
    // 1. Cek Pencarian Teks (Judul / Excerpt)
    const matchesSearch = !searchQuery.value ||
      news.title.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      news.excerpt.toLowerCase().includes(searchQuery.value.toLowerCase())

    // 2. Cek Filter Kategori Tab
    const matchesFilter = selectedFilter.value === 'all' ||
      news.type.includes(selectedFilter.value)

    return matchesSearch && matchesFilter
  })
})

function bukaDetailBerita(news) {
  activeNews.value = news
  detailDialogOpen.value = true
}
</script>

<style scoped>
/* .max-width-center {
  max-width: 1100px;
  margin: 0 auto;
}

.news-public-card {
  transition: transform 0.3s, box-shadow 0.3s;
  border-radius: 12px;
  overflow: hidden;
  height: 100%;
}

.news-public-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.12) !important;
}

.title-link {
  cursor: pointer;
  transition: color 0.2s;
}

.title-link:hover {
  color: var(--q-primary);
}

.lines-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.lines-3 {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.line-height-relaxed {
  line-height: 1.6;
}

.uppercase-text {
  text-transform: uppercase;
} */
</style>

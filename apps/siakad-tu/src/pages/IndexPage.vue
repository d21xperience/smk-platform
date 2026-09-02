<template>
  <q-page>
    <!-- 1. HERO SECTION -->
    <section class="hero-section relative-position">
      <q-carousel v-model="slide" animated arrows infinite autoplay transition-prev="fade" transition-next="fade"
        height="65vh" class="bg-dark">
        <q-carousel-slide :name="1" img-src="https://picsum.photos/1600/900?random=1" class="hero-slide">
          <div class="absolute-full flex flex-center hero-overlay">
            <div class="text-center q-px-md hero-content">
              <div class="text-overline text-weight-medium text-secondary letter-spacing-2">
                LEMBAGA PENDIDIKAN VOKASI TERAKREDITASI
              </div>
              <div class="text-h3 text-sm-h2 text-weight-bold text-white q-mt-sm">
                SMK Maju Berkarya
              </div>
              <div class="text-subtitle1 text-sm-h6 text-grey-3 q-mt-md hero-subtitle">
                Membentuk Generasi Vokasi yang Kompeten, Berdaya Saing, dan Siap Kerja
              </div>
              <div class="q-mt-lg">
                <q-btn unelevated color="secondary" text-color="dark" label="Profil Sekolah"
                  class="text-weight-bold q-px-lg" size="md" no-caps />
                <q-btn outline color="white" label="Pendaftaran Siswa Baru" class="text-weight-bold q-px-lg q-ml-sm"
                  size="md" no-caps />
              </div>
            </div>
          </div>
        </q-carousel-slide>

        <q-carousel-slide :name="2" img-src="https://picsum.photos/1600/900?random=2" class="hero-slide">
          <div class="absolute-full flex flex-center hero-overlay">
            <div class="text-center q-px-md hero-content">
              <div class="text-overline text-weight-medium text-secondary letter-spacing-2">
                KEMITRAAN INDUSTRI
              </div>
              <div class="text-h3 text-sm-h2 text-weight-bold text-white q-mt-sm">
                Fasilitas Berstandar Industri
              </div>
              <div class="text-subtitle1 text-sm-h6 text-grey-3 q-mt-md hero-subtitle">
                Praktik Kerja Nyata Bersama Mitra Perusahaan Nasional & Internasional
              </div>
              <div class="q-mt-lg">
                <q-btn unelevated color="secondary" text-color="dark" label="Lihat Fasilitas"
                  class="text-weight-bold q-px-lg" size="md" no-caps />
              </div>
            </div>
          </div>
        </q-carousel-slide>
      </q-carousel>
    </section>

    <!-- 2. STATISTIK / KEPERCAYAAN INSTITUSIONAL -->
    <section class="stats-section bg-primary text-white">
      <div class="container q-px-md max-width-center">
        <div class="row q-col-gutter-lg text-center">
          <div class="col-6 col-md-3" v-for="stat in stats" :key="stat.label">
            <div class="text-h4 text-weight-bold text-secondary">{{ stat.value }}</div>
            <div class="text-caption text-grey-4 q-mt-xs text-uppercase letter-spacing-1">
              {{ stat.label }}
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- 3. SAMBUTAN KEPALA SEKOLAH -->
    <section class="q-py-xl bg-body">
      <div class="container q-px-md max-width-center">
        <div class="row q-col-gutter-xl items-center">
          <div class="col-12 col-md-4 text-center">
            <q-avatar size="180px" class="shadow-2 principal-avatar">
              <img src="https://pravatar.cc/300" alt="Kepala Sekolah SMK Maju Berkarya" />
            </q-avatar>
          </div>
          <div class="col-12 col-md-8">
            <div class="text-overline text-secondary text-weight-medium letter-spacing-2">
              SAMBUTAN
            </div>
            <div class="text-h5 text-weight-bold text-primary q-mb-md">
              Kepala Sekolah
            </div>
            <q-icon name="format_quote" size="32px" color="grey-4" class="q-mb-sm" />
            <p class="text-body1 text-grey-8 quote-text">
              Selamat datang di SMK Maju Berkarya. Kami berkomitmen menyediakan pendidikan kejuruan
              yang relevan dengan perkembangan industri global. Bersama staf pengajar yang kompeten
              dan fasilitas modern, kami siap mengantarkan putra-putri Anda menuju gerbang kesuksesan
              karier.
            </p>
            <q-separator class="q-my-md" style="width: 60px" size="3px" color="secondary" />
            <div class="text-subtitle1 text-weight-bold text-primary">Dr. H. Ahmad Subarjo, M.Pd.</div>
            <div class="text-caption text-grey">Kepala SMK Maju Berkarya</div>
          </div>
        </div>
      </div>
    </section>

    <!-- 4. KOMPETENSI KEAHLIAN (JURUSAN) -->
    <section class="q-py-xl bg-grey-1">
      <div class="container q-px-md max-width-center">
        <div class="text-center q-mb-xl">
          <div class="text-overline text-secondary text-weight-medium letter-spacing-2">
            PROGRAM STUDI
          </div>
          <div class="text-h5 text-weight-bold text-primary q-mt-xs">
            Kompetensi Keahlian
          </div>
        </div>

        <div class="row q-col-gutter-lg justify-center">
          <div class="col-12 col-sm-6 col-md-4" v-for="jurusan in jurusanList" :key="jurusan.title">
            <q-card class="program-card" flat bordered>
              <q-img :src="jurusan.image" :ratio="16 / 9">
                <div v-if="jurusan.badge" class="absolute-top-left badge-accredited q-pa-xs">
                  <q-icon name="verified" size="14px" class="q-mr-xs" />
                  <span class="text-caption text-weight-bold">{{ jurusan.badge }}</span>
                </div>
              </q-img>
              <q-card-section>
                <div class="text-subtitle2 text-primary text-weight-medium">
                  {{ jurusan.partner }}
                </div>
                <div class="text-h6 text-weight-bold q-mt-xs">{{ jurusan.title }}</div>
                <p class="text-body2 q-mt-sm text-grey-8">
                  {{ jurusan.description }}
                </p>
              </q-card-section>
              <q-separator />
              <q-card-actions align="right">
                <q-btn flat color="primary" label="Detail Jurusan" :icon="jurusan.icon" no-caps
                  class="text-weight-medium" />
              </q-card-actions>
            </q-card>
          </div>
        </div>
      </div>
    </section>

    <!-- 5. BERITA & PRESTASI -->
    <section class="q-py-xl bg-body">
      <div class="container q-px-md max-width-center">
        <div class="row items-center justify-between q-mb-xl">
          <div>
            <div class="text-overline text-secondary text-weight-medium letter-spacing-2">
              INFORMASI TERKINI
            </div>
            <div class="text-h5 text-weight-bold text-primary q-mt-xs">
              Berita &amp; Prestasi
            </div>
          </div>
          <q-btn flat color="primary" label="Lihat Semua" no-caps icon-right="arrow_forward" />
        </div>

        <div class="row q-col-gutter-lg">
          <div class="col-12 col-sm-6 col-md-4" v-for="news in newsList" :key="news.title">
            <q-card flat bordered class="news-card">
              <q-img :src="news.image" ratio="16/9" />
              <q-card-section>
                <div class="text-caption text-grey q-mb-xs">
                  {{ news.date }} • {{ news.author }}
                </div>
                <div class="text-subtitle1 text-weight-bold news-title">
                  {{ news.title }}
                </div>
              </q-card-section>
            </q-card>
          </div>
        </div>
      </div>
    </section>
  </q-page>
</template>

<script setup>
import { ref } from 'vue'

const slide = ref(1)

const stats = ref([
  { value: '25+', label: 'Tahun Berdiri' },
  { value: '1.200+', label: 'Alumni Bekerja' },
  { value: '15', label: 'Mitra Industri' },
  { value: 'A', label: 'Akreditasi' },
])

const jurusanList = ref([
  {
    title: 'Teknik & Bisnis Sepeda Motor (TBSM)',
    partner: 'Kerjasama Resmi Astra Honda Motor (AHM)',
    description:
      'Mempelajari teknik perawatan mekanis, kelistrikan, hingga manajemen bisnis bengkel sepeda motor modern.',
    image: 'https://picsum.photos/600/340?random=11',
    icon: 'motorcycle',
    badge: 'Terakreditasi A',
  },
  {
    title: 'Teknik Komputer & Jaringan',
    partner: 'MikroTik & Cisco Academy',
    description:
      'Fokus pada administrasi jaringan komputer, keamanan siber, cloud computing, dan perakitan infrastruktur IT.',
    image: 'https://picsum.photos/600/340?random=12',
    icon: 'dns',
    badge: null,
  },
  {
    title: 'Akuntansi & Keuangan Lembaga',
    partner: 'Bersertifikat Kompetensi BNSP',
    description:
      'Membekali siswa dengan kemampuan pembukuan, perpajakan, dan pengelolaan keuangan berbasis sistem digital.',
    image: 'https://picsum.photos/600/340?random=13',
    icon: 'calculate',
    badge: null,
  },
])

const newsList = ref([
  {
    title: 'Siswa TBSM Raih Juara 1 Mekanik Edukasi Tingkat Provinsi 2026',
    date: '30 Mei 2026',
    author: 'Kesiswaan',
    image: 'https://picsum.photos/600/340?random=21',
  },
  {
    title: 'Penandatanganan MoU Baru dengan Mitra Industri Otomotif',
    date: '18 Mei 2026',
    author: 'Humas',
    image: 'https://picsum.photos/600/340?random=22',
  },
  {
    title: 'Uji Kompetensi Keahlian (UKK) Tahun Ajaran 2025/2026 Resmi Dibuka',
    date: '05 Mei 2026',
    author: 'Kurikulum',
    image: 'https://picsum.photos/600/340?random=23',
  },
])
</script>

<style scoped>
.max-width-center {
  max-width: 1140px;
  margin: 0 auto;
}

.letter-spacing-1 {
  letter-spacing: 0.06em;
}

.letter-spacing-2 {
  letter-spacing: 0.12em;
}

/* Hero */
.hero-overlay {
  background: linear-gradient(180deg, rgba(10, 25, 47, 0.35) 0%, rgba(10, 25, 47, 0.75) 100%);
}

.hero-content {
  max-width: 720px;
}

.hero-subtitle {
  max-width: 560px;
  margin-left: auto;
  margin-right: auto;
}

/* Stats */
.stats-section {
  padding: 40px 0;
}

/* Principal */
.principal-avatar {
  border: 4px solid white;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.12);
}

.quote-text {
  font-size: 1.05rem;
  line-height: 1.75;
}

/* Program cards */
.program-card {
  border-radius: 6px;
  transition: box-shadow 0.2s ease, transform 0.2s ease;
  overflow: hidden;
}

.program-card:hover {
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
  transform: translateY(-2px);
}

.badge-accredited {
  background-color: rgba(10, 25, 47, 0.85);
  color: white;
  border-radius: 4px;
  margin: 8px;
  display: inline-flex;
  align-items: center;
}

/* News cards */
.news-card {
  border-radius: 6px;
  transition: box-shadow 0.2s ease;
}

.news-card:hover {
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
}

.news-title {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  cursor: pointer;
  transition: color 0.2s ease;
}

.news-title:hover {
  color: var(--q-primary);
}
</style>

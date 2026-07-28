<template>
  <q-page class="bg-grey-1 q-py-xl">
    <div class="container q-px-md max-width-center">

      <!-- 1. HEADER HALAMAN -->
      <div class="text-center q-mb-xl">
        <div class="text-h4 text-sm-h3 text-weight-bold text-primary">
          Fasilitas & <span class="text-amber">Infrastruktur</span>
        </div>
        <p class="text-body1 text-grey-7 q-mt-sm">
          Lingkungan belajar modern dengan sarana dan laboratorium berstandar industri nasional dan internasional.
        </p>
      </div>

      <!-- 2. FILTER KATEGORI FASILITAS -->
      <div class="row justify-center q-mb-xl">
        <q-tabs v-model="selectedCategory" dense class="text-grey-7 bg-white shadow-1 rounded-borders q-pa-xs"
          active-color="primary" indicator-color="primary" align="center" narrow-indicator>
          <q-tab name="all" icon="business" label="Semua Fasilitas" />
          <q-tab name="tbsm" icon="motorcycle" label="Lab / Bengkel TBSM" />
          <q-tab name="tkj" icon="computer" label="Lab Komputer TKJ" />
          <q-tab name="umum" icon="school" label="Fasilitas Umum" />
        </q-tabs>
      </div>

      <!-- 3. GRID VISUAL FASILITAS -->
      <div class="row q-col-gutter-lg justify-center">
        <div v-for="(facility, index) in filteredFacilities" :key="index" class="col-12 col-sm-6 col-md-4">
          <q-card class="facility-card shadow-1 flat bordered flex flex-column no-wrap bg-white">

            <!-- Gambar Fasilitas dengan Efek Hover Zoom -->
            <q-img :src="facility.image" :ratio="16 / 10" class="cursor-pointer" @click="bukaLightbox(facility)">
              <div class="absolute-bottom text-subtitle2 bg-black-5 text-weight-bold row items-center justify-between">
                <span>{{ facility.name }}</span>
                <q-icon name="zoom_in" size="xs" />
              </div>
            </q-img>

            <!-- Deskripsi Singkat Alat Praktek / Sarana -->
            <q-card-section class="col q-pa-md">
              <p class="text-body2 text-grey-8 q-mb-none line-height-relaxed">
                {{ facility.description }}
              </p>

              <!-- Tag Sub-Fasilitas Kecil -->
              <div class="q-mt-sm q-gutter-xs">
                <q-badge v-for="(tag, i) in facility.tags" :key="i" outline color="primary"
                  class="text-weight-medium text-caption">
                  {{ tag }}
                </q-badge>
              </div>
            </q-card-section>
          </q-card>
        </div>
      </div>

    </div>

    <!-- 4. DIALOG POP-UP INTERAKTIF (LIGHTBOX) -->
    <q-dialog v-model="lightboxOpen">
      <q-card style="width: 90vw; max-width: 650px;" class="bg-transparent no-shadow">
        <q-card-section class="row items-center justify-end q-pa-none q-mb-xs">
          <q-btn icon="close" flat round dense v-close-popup color="white" />
        </q-card-section>

        <q-img :src="activeFacility.image" class="rounded-borders shadow-5" />

        <div class="bg-white q-pa-md rounded-borders q-mt-sm shadow-5">
          <div class="text-h6 text-weight-bold text-grey-9">{{ activeFacility.name }}</div>
          <p class="text-body2 text-grey-7 q-mt-xs q-mb-none">{{ activeFacility.description }}</p>
        </div>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed } from 'vue'

const selectedCategory = ref('all')
const lightboxOpen = ref(false)
const activeFacility = ref({})

// Data Master Fasilitas Sekolah (TBSM, TKJ, dan Umum)
const facilities = ref([
  {
    name: 'Bengkel Utama Standar AHASS',
    category: 'tbsm',
    image: 'https://unsplash.com',
    description: 'Laboratorium praktik utama TBSM yang didesain menyerupai layout bengkel resmi AHASS Honda, lengkap dengan bike lift pneumatik.',
    tags: ['Bike Lift', 'Honda Standard', 'Pneumatic Tools']
  },
  {
    name: 'Ruang Teori & Simulasi Injeksi FI',
    category: 'tbsm',
    image: 'https://unsplash.com',
    description: 'Ruang kelas khusus untuk simulasi troubleshooting kelistrikan sistem PGM-FI menggunakan papan sirkuit simulasi digital.',
    tags: ['PGM-FI Trainer', 'Scanner HIDS', 'Kelistrikan']
  },
  {
    name: 'Laboratorium Komputer & Server TKJ',
    category: 'tkj',
    image: 'https://unsplash.com',
    description: 'Lab jaringan komputer yang dilengkapi rak server fisik, puluhan router MikroTik, dan switch Cisco untuk simulasi topologi jaringan berskala besar.',
    tags: ['MikroTik Router', 'Cisco Switch', 'Linux Server']
  },
  {
    name: 'Laboratorium Cloud & Cyber Security',
    category: 'tkj',
    image: 'https://unsplash.com',
    description: 'Ruangan khusus ber-AC berkecepatan internet tinggi untuk melatih siswa melakukan monitoring keamanan siber dan konfigurasi cloud computing.',
    tags: ['Fiber Optic', 'Cyber Defense', 'Cloud Server']
  },
  {
    name: 'Perpustakaan Digital (E-Library)',
    category: 'umum',
    image: 'https://unsplash.com',
    description: 'Ruang baca tenang yang dilengkapi puluhan tablet komputer untuk mempermudah siswa mengakses ribuan e-book materi pelajaran dan jurnal industri.',
    tags: ['E-Book', 'Ruang AC', 'Wi-Fi Area']
  },
  {
    name: 'Gedung Olahraga & Seni (Indoor)',
    category: 'umum',
    image: 'https://unsplash.com',
    description: 'Fasilitas lapangan olahraga tertutup multifungsi untuk kegiatan turnamen futsal, basket, bulutangkis, serta acara panggung pameran seni OSIS.',
    tags: ['Futsal Indoor', 'Stage Mandiri', 'Basketball']
  }
])

// Logika Penyaringan Kategori
const filteredFacilities = computed(() => {
  if (selectedCategory.value === 'all') return facilities.value
  return facilities.value.filter(f => f.category === selectedCategory.value)
})

function bukaLightbox(facility) {
  activeFacility.value = facility
  lightboxOpen.value = true
}
</script>

<style scoped>
.max-width-center {
  max-width: 1100px;
  margin: 0 auto;
}

.facility-card {
  transition: transform 0.3s, box-shadow 0.3s;
  border-radius: 12px;
  overflow: hidden;
  height: 100%;
}

.facility-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.12) !important;
}

.bg-black-5 {
  background: rgba(0, 0, 0, 0.55) !important;
}

.line-height-relaxed {
  line-height: 1.6;
}

.rounded-top {
  border-top-left-radius: 12px;
  border-top-right-radius: 12px;
}
</style>

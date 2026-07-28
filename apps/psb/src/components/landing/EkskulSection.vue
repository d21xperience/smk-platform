<template>
  <!-- ====================================================================== -->
  <!-- NEW SECTION: EKSTRAKURIKULER (DEVELOPING TALENTS & SKILLS)            -->
  <!-- ====================================================================== -->
  <section class="q-py-xl bg-grey-1 border-top-light">
    <div class="container q-px-md max-width-center">
      <!-- HEADER SECTION -->
      <div class="text-center q-mb-xl">
        <div class="text-h5 text-sm-h4 text-weight-bold text-primary">Kegiatan Ekstrakurikuler</div>
        <p class="text-grey-7 q-mt-xs q-mb-none text-caption text-sm-body2">
          Wadah pengembangan minat, bakat, kepemimpinan, dan kreativitas siswa di luar jam akademis.
        </p>
      </div>

      <!-- GRID DAFAR EKSTRAKURIKULER -->
      <div class="row q-col-gutter-lg justify-center">
        <div v-for="(ekskul, index) in ekstraList" :key="index" class="col-12 col-sm-6 col-md-3">
          <q-card class="ekskul-card shadow-1 flat bordered bg-white overflow-hidden">
            <!-- Foto Dokumentasi Ekskul -->
            <q-img :src="ekskul.image" :ratio="4 / 3">
              <div class="absolute-top-right q-ma-sm q-pa-none bg-transparent">
                <q-badge
                  :color="getEkskulBadgeColor(ekskul.category)"
                  class="text-weight-bold q-px-sm q-py-xs"
                >
                  {{ ekskul.category }}
                </q-badge>
              </div>
            </q-img>

            <!-- Deskripsi & Judul Ekskul -->
            <q-card-section class="q-pa-md">
              <div class="text-subtitle1 text-weight-bold text-grey-9 row items-center no-wrap">
                <q-icon :name="ekskul.icon" color="primary" size="xs" class="q-mr-xs" />
                <span class="ellipsis">{{ ekskul.name }}</span>
              </div>
              <p class="text-caption text-grey-7 q-mt-sm q-mb-none lines-2 line-height-relaxed">
                {{ ekskul.description }}
              </p>
            </q-card-section>

            <q-separator />

            <!-- Info Jadwal Latihan -->
            <q-card-actions
              class="bg-grey-1 q-px-md q-py-xs justify-between items-center text-caption text-grey-6"
            >
              <div class="row items-center">
                <q-icon name="schedule" size="xs" class="q-mr-xs" />
                <span>{{ ekskul.schedule }}</span>
              </div>
              <div
                class="text-weight-bold text-primary cursor-pointer"
                @click="gabungEkskul(ekskul.name)"
              >
                Gabung <q-icon name="chevron_right" />
              </div>
            </q-card-actions>
          </q-card>
        </div>
      </div>

      <!-- VIEW MORE LINK -->
      <div class="row justify-center q-mt-xl">
        <q-btn
          flat
          color="primary"
          label="Lihat Galeri Prestasi & Liputan Ekskul"
          icon="collections_bookmark"
          to="/mading"
          class="text-weight-bold"
        />
      </div>
    </div>
  </section>
</template>

<script setup>
import { ref } from 'vue'
import { useQuasar } from 'quasar'
const $q = useQuasar()
// Data Master Kegiatan Ekstrakurikuler Sekolah
const ekstraList = ref([
  {
    name: 'Futsal & Sepak Bola',
    category: 'Olahraga',
    icon: 'sports_soccer',
    schedule: 'Rabu & Jumat (15.30)',
    description:
      'Mengembangkan bakat kerja sama tim, taktik, Fisik, dan teknik futsal siswa untuk kompetisi antar sekolah.',
    image: 'https://unsplash.com',
  },
  {
    name: 'Road Race & Racing Club',
    category: 'Teknologi',
    icon: 'motorcycle',
    schedule: 'Sabtu (09.00)',
    description:
      'Wadah khusus siswa TBSM mendalami teknik modifikasi mesin balap, manajemen mekanik pit-stop, dan keselamatan berkendara.',
    image: 'https://unsplash.com',
  },
  {
    name: 'Cyber Security & Robotik',
    category: 'Teknologi',
    icon: 'smart_toy',
    schedule: 'Kamis (15.30)',
    description:
      'Eksplorasi pembuatan robotika berbasis IoT, pemrograman mikrokontroler, dan latihan dasar pertahanan siber (CTF).',
    image: 'https://unsplash.com',
  },
  {
    name: 'Pramuka Garuda',
    category: 'Kerohanian',
    icon: 'military_tech',
    schedule: 'Sabtu (13.00)',
    description:
      'Pembentukan mental kepemimpinan mandiri, kedisiplinan tinggi, ketangkasan bertahan hidup, dan pengabdian masyarakat.',
    image: 'https://unsplash.com',
  },
])

function getEkskulBadgeColor(category) {
  if (category === 'Olahraga') return 'orange-8'
  if (category === 'Teknologi') return 'blue-7'
  return 'green-7'
}

function gabungEkskul(namaEkskul) {
  $q.notify({
    type: 'info',
    message: `Pendaftaran ekskul ${namaEkskul} dapat dilakukan secara mandiri oleh siswa aktif melalui portal internal SIAKAD setelah masuk tahun ajaran baru.`,
    position: 'bottom',
    timeout: 3500,
  })
}
</script>

<style scoped>
.ekskul-card {
  border-radius: 14px;
  transition:
    transform 0.3s ease,
    box-shadow 0.3s ease;
  height: 100%;
}

.ekskul-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.1) !important;
}

.lines-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>

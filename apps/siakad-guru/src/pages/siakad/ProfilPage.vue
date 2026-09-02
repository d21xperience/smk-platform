<template>
  <q-page class="profile-page q-pa-md">
    <!-- HEADER PROFIL -->
    <q-card flat bordered class="profile-banner q-mb-md overflow-hidden">
      <div class="profile-banner__bg"></div>
      <q-card-section class="profile-banner__content">
        <div class="row items-end q-col-gutter-md">
          <div class="col-auto">
            <q-avatar size="96px" class="profile-banner__avatar shadow-3">
              <img v-if="profile.foto" :src="profile.foto" :alt="profile.nama" />
              <q-icon v-else name="person" size="56px" color="grey-6" />
            </q-avatar>
          </div>
          <div class="col">
            <div class="text-h6 text-weight-bold text-white">{{ profile.nama }}</div>
            <div class="text-body2 text-grey-3">
              {{ profile.jabatan }} &middot; {{ profile.unitKerja }}
            </div>
            <div class="row q-gutter-xs q-mt-xs">
              <q-badge color="amber-8" text-color="white" class="q-px-sm">
                NIP {{ profile.nip }}
              </q-badge>
              <q-badge outline color="white" text-color="white" class="q-px-sm">
                {{ profile.statusKepegawaian }}
              </q-badge>
            </div>
          </div>
          <div class="col-auto gt-xs">
            <q-btn
              no-caps
              outline
              color="white"
              icon="edit"
              label="Edit Profil"
              @click="showEditNotice"
            />
          </div>
        </div>
      </q-card-section>
    </q-card>

    <div class="row q-col-gutter-md">
      <!-- KOLOM KIRI: DATA DIRI -->
      <div class="col-12 col-md-4">
        <q-card flat bordered class="info-card q-mb-md">
          <q-card-section>
            <div class="text-subtitle1 text-weight-bold q-mb-sm">Data Pribadi</div>
            <q-list dense>
              <q-item v-for="item in dataPribadi" :key="item.label" class="q-px-none">
                <q-item-section>
                  <q-item-label caption>{{ item.label }}</q-item-label>
                  <q-item-label>{{ item.value || '-' }}</q-item-label>
                </q-item-section>
              </q-item>
            </q-list>
          </q-card-section>
        </q-card>

        <q-card flat bordered class="info-card">
          <q-card-section>
            <div class="text-subtitle1 text-weight-bold q-mb-sm">Data Kepegawaian</div>
            <q-list dense>
              <q-item v-for="item in dataKepegawaian" :key="item.label" class="q-px-none">
                <q-item-section>
                  <q-item-label caption>{{ item.label }}</q-item-label>
                  <q-item-label>{{ item.value || '-' }}</q-item-label>
                </q-item-section>
              </q-item>
            </q-list>
          </q-card-section>
        </q-card>

        <!-- ALERT BERKAS -->
        <div class="q-mt-md" v-if="documentAlerts.length">
          <div class="text-subtitle1 text-weight-bold q-mb-sm">Perlu Perhatian</div>
          <q-card
            flat
            bordered
            class="alert-card q-mb-sm"
            v-for="alert in documentAlerts"
            :key="alert.id"
          >
            <q-card-section class="row items-center no-wrap q-gutter-sm">
              <q-icon name="warning" color="red" size="24px" />
              <div>
                <div class="text-weight-medium text-body2">{{ alert.title }}</div>
                <div class="text-caption text-grey-7">{{ alert.description }}</div>
              </div>
            </q-card-section>
          </q-card>
        </div>
      </div>

      <!-- KOLOM KANAN: ARSIP BERKAS -->
      <div class="col-12 col-md-8">
        <q-card flat bordered>
          <q-card-section>
            <div class="row items-center justify-between q-mb-md wrap">
              <div class="text-subtitle1 text-weight-bold">Arsip Berkas</div>
              <div class="row q-gutter-sm items-center">
                <q-btn
                  no-caps
                  unelevated
                  color="primary"
                  icon="upload_file"
                  label="Unggah Berkas"
                  @click="showUploadNotice"
                />
              </div>
            </div>

            <!-- FILTER & PENCARIAN -->
            <div class="row q-col-gutter-sm q-mb-md">
              <div class="col-12 col-sm-6">
                <q-input
                  dense
                  outlined
                  v-model="searchQuery"
                  placeholder="Cari nama berkas..."
                  clearable
                >
                  <template v-slot:prepend>
                    <q-icon name="search" />
                  </template>
                </q-input>
              </div>
              <div class="col-12 col-sm-6">
                <q-select
                  dense
                  outlined
                  v-model="filterType"
                  :options="typeFilterOptions"
                  emit-value
                  map-options
                  label="Jenis Berkas"
                />
              </div>
            </div>

            <!-- DAFTAR BERKAS PER TAHUN -->
            <q-list v-if="filteredYears.length" bordered class="rounded-borders">
              <q-expansion-item
                v-for="group in filteredYears"
                :key="group.year"
                :label="`Tahun ${group.year}`"
                :caption="`${group.files.length} berkas`"
                header-class="year-header text-weight-bold"
                default-opened
                icon="folder"
              >
                <q-separator />
                <q-list separator>
                  <q-item v-for="file in group.files" :key="file.id">
                    <q-item-section avatar>
                      <q-icon
                        :name="fileIcon(file.ext)"
                        :color="fileIconColor(file.ext)"
                        size="32px"
                      />
                    </q-item-section>

                    <q-item-section>
                      <q-item-label class="text-weight-medium">{{ file.name }}</q-item-label>
                      <q-item-label caption>
                        {{ file.type }} &middot; {{ file.size }} &middot; Diunggah
                        {{ file.uploadedAt }}
                      </q-item-label>
                    </q-item-section>

                    <q-item-section side v-if="file.status !== 'valid'">
                      <q-badge :color="file.status === 'expired' ? 'red' : 'orange'" outline>
                        {{ file.status === 'expired' ? 'Kadaluarsa' : 'Perlu Diperbarui' }}
                      </q-badge>
                    </q-item-section>

                    <q-item-section side>
                      <q-btn
                        flat
                        round
                        dense
                        icon="download"
                        color="primary"
                        @click="handleDownload(file)"
                      >
                        <q-tooltip>Unduh Berkas</q-tooltip>
                      </q-btn>
                    </q-item-section>
                  </q-item>
                </q-list>
              </q-expansion-item>
            </q-list>

            <!-- STATE KOSONG -->
            <div v-else class="text-center text-grey-6 q-pa-xl">
              <q-icon name="folder_off" size="40px" class="q-mb-sm" />
              <div>Tidak ada berkas yang cocok dengan pencarian</div>
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()

// TODO: ganti dengan data dari endpoint GET /api/guru/profile (Go backend)
const profile = reactive({
  nama: 'Deden Moh Jaenudin, S.H.',
  jabatan: 'Guru Mata Pelajaran Pendidikan Pancasila',
  unitKerja: 'SMK Pasundan Jatinangor',
  nip: '198705122011011008',
  statusKepegawaian: 'PNS',
  foto: '', // TODO: isi dengan URL foto profil dari server, kosong = tampil ikon default
})

// TODO: ganti dengan data dari endpoint GET /api/guru/profile (Go backend)
const dataPribadi = reactive([
  { label: 'NUPTK', value: '1234567890123456' },
  { label: 'Tempat, Tanggal Lahir', value: 'Sumedang, 21 September 1987' },
  { label: 'Jenis Kelamin', value: 'Laki-laki' },
  { label: 'Alamat', value: 'Jl. Raya Jatinangor No. 45, Sumedang' },
  { label: 'No. HP', value: '0812-3456-7890' },
  { label: 'Email', value: 'd21xperience@gmail.com' },
])

// TODO: ganti dengan data dari endpoint GET /api/guru/profile (Go backend)
const dataKepegawaian = reactive([
  { label: 'TMT (Tanggal Mulai Tugas)', value: '1 Januari 2011' },
  { label: 'Pangkat/Golongan', value: 'Penata Tk. I / III-d' },
  { label: 'Jurusan/Program Keahlian', value: 'Rekayasa Perangkat Lunak' },
  { label: 'Status Sertifikasi', value: 'Sudah Sertifikasi (2016)' },
])

// TODO: ganti dengan data dari endpoint GET /api/guru/berkas (Go backend)
// Struktur: dikelompokkan per tahun terbit/unggah berkas
const documentsByYear = reactive([
  {
    year: 2024,
    files: [
      {
        id: 'f1',
        name: 'SK Mengajar Tahun Ajaran 2024-2025.pdf',
        type: 'SK Mengajar',
        ext: 'pdf',
        size: '1.2 MB',
        uploadedAt: '15 Jul 2024',
        status: 'valid',
        url: '/files/sk-mengajar-2024.pdf',
      },
    ],
  },
  {
    year: 2016,
    files: [
      {
        id: 'f2',
        name: 'Sertifikat Pendidik Profesional.pdf',
        type: 'Sertifikat Pendidik',
        ext: 'pdf',
        size: '890 KB',
        uploadedAt: '2 Mar 2016',
        status: 'valid',
        url: '/files/sertifikat-pendidik.pdf',
      },
    ],
  },
  {
    year: 2011,
    files: [
      {
        id: 'f3',
        name: 'SK Pengangkatan CPNS.pdf',
        type: 'SK Pengangkatan',
        ext: 'pdf',
        size: '650 KB',
        uploadedAt: '5 Jan 2011',
        status: 'valid',
        url: '/files/sk-pengangkatan-2011.pdf',
      },
      {
        id: 'f4',
        name: 'Ijazah S1 Ilmu Hukum.jpg',
        type: 'Ijazah',
        ext: 'jpg',
        size: '2.1 MB',
        uploadedAt: '5 Jan 2011',
        status: 'perlu_perbarui',
        url: '/files/ijazah-s1.jpg',
      },
    ],
  },
])

const searchQuery = ref('')
const filterType = ref('semua')

const typeFilterOptions = computed(() => {
  const types = new Set()
  documentsByYear.forEach((group) => group.files.forEach((f) => types.add(f.type)))
  return [
    { label: 'Semua Jenis', value: 'semua' },
    ...Array.from(types).map((t) => ({ label: t, value: t })),
  ]
})

// Filter berkas berdasarkan pencarian nama & jenis, tanpa mengubah data asli
const filteredYears = computed(() => {
  const query = (searchQuery.value || '').toLowerCase().trim()

  return documentsByYear
    .map((group) => {
      const files = group.files.filter((file) => {
        const matchQuery = !query || file.name.toLowerCase().includes(query)
        const matchType = filterType.value === 'semua' || file.type === filterType.value
        return matchQuery && matchType
      })
      return { year: group.year, files }
    })
    .filter((group) => group.files.length > 0)
    .sort((a, b) => b.year - a.year)
})

// Peringatan otomatis untuk berkas kadaluarsa / perlu diperbarui / belum diunggah
const documentAlerts = computed(() => {
  const alerts = []

  documentsByYear.forEach((group) => {
    group.files.forEach((file) => {
      if (file.status === 'expired') {
        alerts.push({
          id: `alert-${file.id}`,
          title: `${file.type} sudah kadaluarsa`,
          description: file.name,
        })
      } else if (file.status === 'perlu_perbarui') {
        alerts.push({
          id: `alert-${file.id}`,
          title: `${file.type} perlu diperbarui`,
          description: file.name,
        })
      }
    })
  })

  // TODO: tambahkan pengecekan berkas wajib yang belum diunggah sama sekali
  // (mis. Ijazah terbaru, Sertifikat Pendidik) berdasarkan aturan dari backend

  return alerts
})

function fileIcon(ext) {
  if (ext === 'pdf') return 'picture_as_pdf'
  if (['jpg', 'jpeg', 'png'].includes(ext)) return 'image'
  return 'description'
}

function fileIconColor(ext) {
  if (ext === 'pdf') return 'red'
  if (['jpg', 'jpeg', 'png'].includes(ext)) return 'blue'
  return 'grey-7'
}

function handleDownload(file) {
  // TODO: ganti dengan permintaan unduh sesungguhnya ke Go backend,
  // misalnya GET /api/guru/berkas/:id/download yang mengembalikan signed URL atau blob
  $q.notify({
    type: 'positive',
    message: `Mengunduh "${file.name}"...`,
    position: 'top',
    timeout: 1500,
  })
  // window.open(file.url, '_blank')
}

function showUploadNotice() {
  // TODO: buka dialog unggah berkas baru, lalu POST ke /api/guru/berkas (Go backend)
  $q.notify({
    type: 'info',
    message: 'Fitur unggah berkas akan segera tersedia.',
    position: 'top',
    timeout: 2000,
  })
}

function showEditNotice() {
  // TODO: arahkan ke halaman/dialog edit profil, lalu PUT ke /api/guru/profile (Go backend)
  $q.notify({
    type: 'info',
    message: 'Fitur edit profil akan segera tersedia.',
    position: 'top',
    timeout: 2000,
  })
}
</script>

<style scoped>
.profile-page {
  max-width: 1200px;
  margin: 0 auto;
}

.profile-banner {
  border-radius: 8px;
  position: relative;
}

.profile-banner__bg {
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, #0d2b4e 0%, #1a4a7a 100%);
}

.profile-banner__content {
  position: relative;
  z-index: 1;
}

.profile-banner__avatar {
  border: 3px solid white;
  background: white;
}

.info-card,
.alert-card {
  border-radius: 8px;
}

.alert-card {
  border-left: 4px solid var(--q-red, #c10015);
}

.year-header {
  font-size: 0.95rem;
}
</style>

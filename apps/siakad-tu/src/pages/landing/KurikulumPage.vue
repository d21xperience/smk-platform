<template>
  <q-page class="dashboard-page bg-grey-1">
    <!-- HEADER -->
    <div class="dashboard-header q-px-md q-py-lg">
      <div class="container max-width-center row items-center justify-between">
        <div>
          <div class="text-overline text-secondary text-weight-medium letter-spacing-2">
            SISTEM INFORMASI AKADEMIK
          </div>
          <div class="text-h5 text-weight-bold text-white q-mt-xs">
            Dashboard Kurikulum
          </div>
          <div class="text-caption text-grey-4 q-mt-xs">
            {{ todayLabel }} • Semester Genap 2025/2026
          </div>
        </div>
        <div class="row items-center q-gutter-sm">
          <q-btn unelevated color="secondary" text-color="dark" label="Susun Jadwal" no-caps icon="event" />
          <q-btn round flat color="white" icon="notifications">
            <q-badge color="secondary" text-color="dark" floating rounded>{{ notifCount }}</q-badge>
          </q-btn>
        </div>
      </div>
    </div>

    <div class="container max-width-center q-px-md q-py-lg">
      <!-- KPI CARDS -->
      <div class="row q-col-gutter-md q-mb-lg">
        <div class="col-6 col-md-3" v-for="kpi in kpis" :key="kpi.label">
          <q-card flat bordered class="kpi-card">
            <q-card-section class="row items-center no-wrap">
              <q-avatar :color="kpi.color" text-color="white" :icon="kpi.icon" size="46px" class="q-mr-md" />
              <div>
                <div class="text-h6 text-weight-bold text-primary">{{ kpi.value }}</div>
                <div class="text-caption text-grey-7">{{ kpi.label }}</div>
              </div>
            </q-card-section>
          </q-card>
        </div>
      </div>

      <div class="row q-col-gutter-lg">
        <!-- LEFT COLUMN -->
        <div class="col-12 col-md-8">
          <!-- MAPEL PERLU PERHATIAN -->
          <q-card flat bordered class="q-mb-lg content-card border-alert" v-if="mapelPerhatian.length">
            <q-card-section>
              <div class="row items-center">
                <q-icon name="report" color="red" size="22px" class="q-mr-sm" />
                <div class="text-subtitle1 text-weight-bold text-red-9">
                  Mata Pelajaran Tertinggal dari Target
                </div>
              </div>
            </q-card-section>
            <q-separator />
            <q-list separator>
              <q-item v-for="mapel in mapelPerhatian" :key="mapel.nama">
                <q-item-section avatar>
                  <q-avatar color="red-1" text-color="red-9" icon="menu_book" size="38px" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-medium">{{ mapel.nama }}</q-item-label>
                  <q-item-label caption>{{ mapel.kelas }} • Guru: {{ mapel.guru }}</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-badge color="red" text-color="white" class="q-px-sm">
                    Progress {{ mapel.progress }}%
                  </q-badge>
                </q-item-section>
              </q-item>
            </q-list>
          </q-card>

          <!-- PROGRESS PENYELESAIAN MATERI -->
          <q-card flat bordered class="q-mb-lg content-card">
            <q-card-section>
              <div class="row items-center justify-between">
                <div class="text-subtitle1 text-weight-bold text-primary">
                  Progress Penyelesaian Materi per Kelompok Mapel
                </div>
                <q-badge color="grey-3" text-color="grey-8" class="q-px-sm">
                  Minggu ke-{{ mingguKe }} semester
                </q-badge>
              </div>
            </q-card-section>
            <q-separator />
            <q-card-section>
              <div v-for="kelompok in progressMapel" :key="kelompok.nama" class="q-mb-md">
                <div class="row items-center justify-between q-mb-xs">
                  <div class="text-body2 text-weight-medium">{{ kelompok.nama }}</div>
                  <div class="text-caption text-grey-7">{{ kelompok.progress }}% dari target</div>
                </div>
                <div class="bar-track">
                  <div
                    class="bar-fill"
                    :style="{ width: kelompok.progress + '%', backgroundColor: progressWarna(kelompok.progress) }"
                  />
                </div>
              </div>
            </q-card-section>
          </q-card>

          <!-- STATUS PENGUMPULAN PERANGKAT AJAR -->
          <q-card flat bordered class="content-card">
            <q-card-section>
              <div class="row items-center justify-between">
                <div class="text-subtitle1 text-weight-bold text-primary">
                  Status Pengumpulan RPP &amp; Perangkat Ajar
                </div>
                <q-btn flat dense color="primary" label="Lihat Semua" no-caps icon-right="arrow_forward" />
              </div>
            </q-card-section>
            <q-separator />
            <q-list separator>
              <q-item v-for="rpp in statusRpp" :key="rpp.id">
                <q-item-section avatar>
                  <q-avatar :color="statusWarna(rpp.status)" text-color="white" icon="description" size="38px" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-medium">{{ rpp.mapel }}</q-item-label>
                  <q-item-label caption>{{ rpp.guru }} • {{ rpp.kelas }} • Tenggat {{ rpp.tenggat }}</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-chip dense :color="statusWarna(rpp.status)" text-color="white" class="text-caption">
                    {{ rpp.status }}
                  </q-chip>
                </q-item-section>
              </q-item>
            </q-list>
          </q-card>
        </div>

        <!-- RIGHT COLUMN -->
        <div class="col-12 col-md-4">
          <!-- BEBAN MENGAJAR GURU -->
          <q-card flat bordered class="q-mb-lg content-card">
            <q-card-section>
              <div class="text-subtitle1 text-weight-bold text-primary">
                Beban Mengajar Guru (JP/Minggu)
              </div>
            </q-card-section>
            <q-separator />
            <q-list separator>
              <q-item v-for="guru in bebanMengajar" :key="guru.nama">
                <q-item-section>
                  <q-item-label class="text-weight-medium">{{ guru.nama }}</q-item-label>
                  <q-item-label caption>{{ guru.mapel }}</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-badge :color="guru.jp > 24 ? 'orange' : 'primary'" text-color="white">
                    {{ guru.jp }} JP
                  </q-badge>
                </q-item-section>
              </q-item>
            </q-list>
          </q-card>

          <!-- JADWAL ASESMEN -->
          <q-card flat bordered class="q-mb-lg content-card">
            <q-card-section>
              <div class="text-subtitle1 text-weight-bold text-primary">
                Jadwal Asesmen Mendatang
              </div>
            </q-card-section>
            <q-separator />
            <q-list separator>
              <q-item v-for="asesmen in jadwalAsesmen" :key="asesmen.judul">
                <q-item-section avatar>
                  <div class="agenda-date text-center">
                    <div class="text-weight-bold text-primary">{{ asesmen.tanggal }}</div>
                    <div class="text-caption text-grey-7">{{ asesmen.bulan }}</div>
                  </div>
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-medium">{{ asesmen.judul }}</q-item-label>
                  <q-item-label caption>{{ asesmen.cakupan }}</q-item-label>
                </q-item-section>
              </q-item>
            </q-list>
          </q-card>

          <!-- KKM & KELULUSAN -->
          <q-card flat bordered class="content-card">
            <q-card-section>
              <div class="text-subtitle1 text-weight-bold text-primary">
                Ketercapaian KKM per Mapel Inti
              </div>
            </q-card-section>
            <q-separator />
            <q-list separator>
              <q-item v-for="kkm in ketercapaianKkm" :key="kkm.mapel">
                <q-item-section>
                  <q-item-label class="text-weight-medium">{{ kkm.mapel }}</q-item-label>
                  <q-item-label caption>KKM {{ kkm.nilaiKkm }}</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <div class="text-weight-bold" :style="{ color: kkm.persenTuntas >= 85 ? '#2e7d32' : '#c62828' }">
                    {{ kkm.persenTuntas }}% tuntas
                  </div>
                </q-item-section>
              </q-item>
            </q-list>
          </q-card>
        </div>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, computed } from 'vue'

const todayLabel = computed(() => {
  const today = new Date()
  return today.toLocaleDateString('id-ID', {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  })
})

const notifCount = ref(5)
const mingguKe = 14

const kpis = ref([
  { label: 'Total Mata Pelajaran', value: '38', icon: 'menu_book', color: 'primary' },
  { label: 'Guru Mengajar Aktif', value: '76', icon: 'badge', color: 'secondary' },
  { label: 'Progress Kurikulum Semester', value: '68%', icon: 'trending_up', color: 'teal' },
  { label: 'RPP Belum Dikumpulkan', value: '11', icon: 'pending_actions', color: 'orange-8' },
])

const mapelPerhatian = ref([
  { nama: 'Matematika', kelas: 'X TKJ 2', guru: 'Ibu Sri Wahyuni', progress: 42 },
  { nama: 'Produktif TBSM - Sistem Kelistrikan', kelas: 'XI TBSM 1', guru: 'Bpk. Slamet Riyadi', progress: 38 },
])

const progressMapel = ref([
  { nama: 'Mata Pelajaran Umum (Adaptif & Normatif)', progress: 72 },
  { nama: 'Produktif TBSM', progress: 65 },
  { nama: 'Produktif TKJ', progress: 70 },
  { nama: 'Produktif Akuntansi', progress: 58 },
  { nama: 'Muatan Lokal & Proyek P5', progress: 80 },
])

const progressWarna = (progress) => {
  if (progress >= 70) return '#2e7d32'
  if (progress >= 50) return '#f9a825'
  return '#c62828'
}

const statusRpp = ref([
  {
    id: 1,
    mapel: 'Matematika',
    guru: 'Ibu Sri Wahyuni',
    kelas: 'X TKJ 2',
    tenggat: '12 Juli 2026',
    status: 'Terlambat',
  },
  {
    id: 2,
    mapel: 'Produktif TBSM - Sistem Kelistrikan',
    guru: 'Bpk. Slamet Riyadi',
    kelas: 'XI TBSM 1',
    tenggat: '14 Juli 2026',
    status: 'Menunggu',
  },
  {
    id: 3,
    mapel: 'Bahasa Indonesia',
    guru: 'Ibu Ratna Sari',
    kelas: 'XII Akuntansi',
    tenggat: '10 Juli 2026',
    status: 'Terkumpul',
  },
  {
    id: 4,
    mapel: 'Produktif TKJ - Administrasi Jaringan',
    guru: 'Bpk. Yusuf Hidayat',
    kelas: 'XI TKJ 1',
    tenggat: '10 Juli 2026',
    status: 'Terkumpul',
  },
])

const statusWarna = (status) => {
  if (status === 'Terlambat') return 'red-7'
  if (status === 'Menunggu') return 'orange-7'
  return 'teal-7'
}

const bebanMengajar = ref([
  { nama: 'Bpk. Slamet Riyadi', mapel: 'Produktif TBSM', jp: 28 },
  { nama: 'Ibu Sri Wahyuni', mapel: 'Matematika', jp: 24 },
  { nama: 'Bpk. Yusuf Hidayat', mapel: 'Produktif TKJ', jp: 26 },
  { nama: 'Ibu Ratna Sari', mapel: 'Bahasa Indonesia', jp: 22 },
])

const jadwalAsesmen = ref([
  { tanggal: '16', bulan: 'Jul', judul: 'Penilaian Tengah Semester Ganjil', cakupan: 'Seluruh Mapel, Kelas X & XI' },
  { tanggal: '17', bulan: 'Jul', judul: 'Ujian Praktik Kejuruan (UPK) TBSM', cakupan: 'Kelas XII TBSM' },
  { tanggal: '24', bulan: 'Jul', judul: 'Asesmen Diagnostik Awal Semester', cakupan: 'Kelas X (Seluruh Jurusan)' },
])

const ketercapaianKkm = ref([
  { mapel: 'Matematika', nilaiKkm: 75, persenTuntas: 78 },
  { mapel: 'Bahasa Indonesia', nilaiKkm: 75, persenTuntas: 91 },
  { mapel: 'Produktif TBSM', nilaiKkm: 78, persenTuntas: 82 },
  { mapel: 'Produktif TKJ', nilaiKkm: 78, persenTuntas: 74 },
])
</script>

<style scoped>
.max-width-center {
  max-width: 1200px;
  margin: 0 auto;
}

.letter-spacing-2 {
  letter-spacing: 0.12em;
}

.dashboard-header {
  background: linear-gradient(135deg, #0a192f 0%, #123a63 100%);
}

.kpi-card {
  border-radius: 8px;
  height: 100%;
}

.content-card {
  border-radius: 8px;
}

.border-alert {
  border-left: 4px solid #c62828;
}

.bar-track {
  width: 100%;
  height: 10px;
  border-radius: 6px;
  background-color: #eceff1;
  overflow: hidden;
}

.bar-fill {
  height: 100%;
  border-radius: 6px;
  transition: width 0.4s ease;
}

.agenda-date {
  min-width: 44px;
}
</style>

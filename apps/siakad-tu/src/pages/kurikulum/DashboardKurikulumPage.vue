<template>
  <q-page class="q-pa-md bg-grey-2">
    <!-- 1. HEADER HALAMAN -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h5 class="text-h5 text-weight-bold q-my-none">Dashboard Kurikulum</h5>
        <div class="text-caption text-grey-7">Selamat datang kembali, Admin TU Kurikulum.</div>
      </div>
      <q-btn color="primary" icon="refresh" label="Segarkan Data" @click="refreshDashboard" :loading="isRefreshing"
        dense class="q-px-sm" />
    </div>

    <!-- 2. BANNER STATUS SINKRONISASI DAPODIK -->
    <q-banner inline-actions rounded
      :class="dapodikSync.isSynced ? 'bg-positive text-white q-mb-lg' : 'bg-warning text-white q-mb-lg'">
      <template v-slot:avatar>
        <q-icon :name="dapodikSync.isSynced ? 'cloud_done' : 'cloud_sync'" size="sm" />
      </template>
      <strong>Status Dapodik:</strong> {{ dapodikSync.message }} (Terakhir sinkron: {{ dapodikSync.lastSync }})
      <template v-slot:action>
        <q-btn flat label="Sinkronkan Sekarang" color="white" @click="syncWithDapodik" :loading="isSyncing" />
      </template>
    </q-banner>

    <!-- 3. KARTU RINGKASAN STATISTIK (GRID) -->
    <div class="row q-col-gutter-md q-mb-lg">
      <div v-for="stat in statsData" :key="stat.title" class="col-12 col-sm-6 col-md-3">
        <q-card flat bordered class="my-stat-card">
          <q-card-section class="row items-center justify-between no-wrap">
            <div>
              <div class="text-overline text-grey-7 text-weight-medium">{{ stat.title }}</div>
              <div class="text-h4 text-weight-bold text-dark q-mt-xs">{{ stat.value }}</div>
            </div>
            <q-avatar :icon="stat.icon" :bg-color="stat.bgColor" :text-color="stat.iconColor" size="lg" />
          </q-card-section>
          <q-separator />
          <q-card-section class="q-py-xs text-caption text-grey-7">
            {{ stat.desc }}
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- 4. BAGIAN UTAMA (AGENDA & LOG AKTIVITAS) -->
    <div class="row q-col-gutter-md">
      <!-- Kolon Kiri: Agenda Akademik Terdekat -->
      <div class="col-12 col-md-6">
        <q-card flat bordered class="full-height">
          <q-card-section class="bg-primary text-white row items-center">
            <q-icon name="event" size="sm" class="q-mr-sm" />
            <div class="text-subtitle1 text-weight-bold">Agenda Akademik Terdekat</div>
          </q-card-section>

          <q-card-section>
            <q-timeline color="primary" dense>
              <q-timeline-entry v-for="agenda in agendaList" :key="agenda.id" :title="agenda.title"
                :subtitle="agenda.date" :side="right" :icon="agenda.icon" :color="agenda.color">
                <div class="text-caption text-grey-8">{{ agenda.desc }}</div>
              </q-timeline-entry>
            </q-timeline>
          </q-card-section>
        </q-card>
      </div>

      <!-- Kolom Kanan: Log Aktivitas Terbaru -->
      <div class="col-12 col-md-6">
        <q-card flat bordered class="full-height">
          <q-card-section class="bg-dark text-white row items-center">
            <q-icon name="history" size="sm" class="q-mr-sm" />
            <div class="text-subtitle1 text-weight-bold">Log Aktivitas Sistem</div>
          </q-card-section>

          <q-card-section class="q-pa-none">
            <q-list separator>
              <q-item v-for="log in logList" :key="log.id" clickable v-ripple>
                <q-item-section avatar>
                  <q-icon :name="log.icon" :color="log.iconColor" />
                </q-item-section>

                <q-item-section>
                  <q-item-label class="text-weight-medium text-body2">{{ log.activity }}</q-item-label>
                  <q-item-label caption>{{ log.user }} &bull; {{ log.time }}</q-item-label>
                </q-item-section>
              </q-item>
            </q-list>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()

// State untuk efek loading tombol
const isRefreshing = ref(false)
const isSyncing = ref(false)

// 1. Data Mock Banner Dapodik
const dapodikSync = ref({
  isSynced: true,
  message: 'Data lokal sinkron dengan server pusat.',
  lastSync: 'Hari ini, 08:30 WIB'
})

// 2. Data Mock Kartu Statistik
const statsData = ref([
  {
    title: 'Total Siswa Aktif',
    value: '842',
    icon: 'groups',
    bgColor: 'blue-1',
    iconColor: 'blue-8',
    desc: 'Terbagi dalam 24 rombel (X, XI, XII)'
  },
  {
    title: 'Guru Pengajar',
    value: '54',
    icon: 'person',
    bgColor: 'green-1',
    iconColor: 'green-8',
    desc: '42 Guru Tetap, 12 Guru Honorer'
  },
  {
    title: 'Modul Ajar Masuk',
    value: '48 / 54',
    icon: 'auto_stories',
    bgColor: 'orange-1',
    iconColor: 'orange-8',
    desc: '6 Guru belum mengumpulkan perangkat'
  },
  {
    title: 'Belum Tuntas Nilai',
    value: '18',
    icon: 'gavel',
    bgColor: 'red-1',
    iconColor: 'red-8',
    desc: 'Siswa dalam proses perbaikan/remedial'
  }
])

// 3. Data Mock Agenda Akademik
const agendaList = ref([
  {
    id: 1,
    title: 'Batas Akhir Penyerahan Modul Ajar',
    date: '15 Juni 2026',
    desc: 'Pengumpulan berkas perangkat ajar Semester Ganjil melalui portal TU.',
    icon: 'assignment_turned_in',
    color: 'warning'
  },
  {
    id: 2,
    title: 'Uji Kompetensi Keahlian (UKK) SMK',
    date: '22 - 26 Juni 2026',
    desc: 'Pelaksanaan ujian praktik produktif menggandeng DUDI (Dunia Usaha Dunia Industri).',
    icon: 'engineering',
    color: 'primary'
  },
  {
    id: 3,
    title: 'Rapat Pleno Kenaikan Kelas',
    date: '30 Juni 2026',
    desc: 'Evaluasi akhir rekapitulasi leger nilai bersama seluruh wali kelas dan Kepala Sekolah.',
    icon: 'gavel',
    color: 'negative'
  }
])

// 4. Data Mock Log Aktivitas Terbaru
const logList = ref([
  {
    id: 1,
    activity: 'Mengubah jadwal pelajaran Kelas XI TKJ 2',
    user: 'Andi (TU Kurikulum)',
    time: '10 menit yang lalu',
    icon: 'edit_calendar',
    iconColor: 'blue'
  },
  {
    id: 2,
    activity: 'Mengunggah RPP Informatika Semester Ganjil',
    user: 'Drs. Supardi (Guru)',
    time: '45 menit yang lalu',
    icon: 'upload_file',
    iconColor: 'green'
  },
  {
    id: 3,
    activity: 'Sinkronisasi profil siswa baru ke Dapodik sukses',
    user: 'Sistem Otomatis',
    time: '2 jam yang lalu',
    icon: 'cloud_done',
    iconColor: 'teal'
  },
  {
    id: 4,
    activity: 'Mencetak Lembar Berita Acara Tryout',
    user: 'Rina (Staf TU)',
    time: 'Kemarin, 14:15 WIB',
    icon: 'print',
    iconColor: 'grey-7'
  }
])

// Fungsi Segarkan Halaman
const refreshDashboard = () => {
  isRefreshing.value = true
  setTimeout(() => {
    isRefreshing.value = false
    $q.notify({
      type: 'positive',
      message: 'Data dashboard berhasil diperbarui!',
      position: 'top-right'
    })
  }, 1000)
}

// Fungsi Simulasi Sinkronisasi Dapodik
const syncWithDapodik = () => {
  isSyncing.value = true
  setTimeout(() => {
    isSyncing.value = false
    dapodikSync.value.isSynced = true
    dapodikSync.value.lastSync = 'Baru saja'
    $q.notify({
      type: 'positive',
      message: 'Sinkronisasi server Dapodik sukses dilakukan.',
      icon: 'cloud_done'
    })
  }, 2000)
}
</script>

<style scoped>
.my-stat-card {
  transition: transform 0.2s;
  border-radius: 8px;
}

.my-stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1) !important;
}
</style>

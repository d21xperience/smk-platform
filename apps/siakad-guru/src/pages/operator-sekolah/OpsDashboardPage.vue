<template>
  <q-page padding>
    <div class="row q-col-gutter-md">
      <!-- Metric Cards -->
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="bg-warning text-white">
          <q-card-section>
            <div class="text-subtitle2">Tugas Pending</div>
            <div class="text-h4">15</div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="bg-positive text-white">
          <q-card-section>
            <div class="text-subtitle2">Selesai Hari Ini</div>
            <div class="text-h4">7</div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="bg-primary text-white">
          <q-card-section>
            <div class="text-subtitle2">Total Siswa</div>
            <div class="text-h4">1.230</div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="bg-secondary text-white">
          <q-card-section>
            <div class="text-subtitle2">GTK</div>
            <div class="text-h4">89</div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Status Koneksi Multi-Sistem -->
    <q-card class="q-mt-md">
      <q-card-section>
        <div class="text-h6">Status Koneksi Sistem</div>
        <div class="row q-col-gutter-sm">
          <div class="col-6 col-md-3">
            <q-chip icon="cloud_done" color="green" text-color="white">Dapodik</q-chip>
          </div>
          <div class="col-6 col-md-3">
            <q-chip icon="cloud_off" color="red" text-color="white">SIAKAD</q-chip>
          </div>
          <div class="col-6 col-md-3">
            <q-chip icon="cloud_done" color="green" text-color="white">e-Rapor</q-chip>
          </div>
          <div class="col-6 col-md-3">
            <q-chip icon="cloud_done" color="green" text-color="white">Backup</q-chip>
          </div>
        </div>
      </q-card-section>
    </q-card>

    <!-- Bar Chart Aktivitas 7 Hari -->
    <q-card class="q-mt-md">
      <q-card-section>
        <div class="text-h6">Aktivitas 7 Hari Terakhir</div>
        <apexchart type="bar" height="250" :options="chartOptions" :series="chartSeries" />
      </q-card-section>
    </q-card>

    <!-- Log Aktivitas Terbaru -->
    <q-card class="q-mt-md">
      <q-card-section>
        <div class="text-h6">Aktivitas Terbaru</div>
        <q-list dense separator>
          <q-item v-for="log in latestLogs" :key="log.id">
            <q-item-section avatar>
              <q-icon :name="log.icon" :color="log.color" />
            </q-item-section>
            <q-item-section>
              <q-item-label>{{ log.deskripsi }}</q-item-label>
              <q-item-label caption>{{ log.waktu }}</q-item-label>
            </q-item-section>
          </q-item>
        </q-list>
      </q-card-section>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref } from 'vue'

// Dummy chart
const chartOptions = {
  chart: { id: 'aktivitas-7hari', toolbar: { show: false } },
  xaxis: { categories: ['Sen', 'Sel', 'Rab', 'Kam', 'Jum', 'Sab', 'Min'] },
  colors: ['#1976d2']
}
const chartSeries = [{
  name: 'Jumlah Tugas',
  data: [12, 8, 15, 10, 9, 4, 6]
}]

// Dummy logs
const latestLogs = ref([
  { id: 1, icon: 'check_circle', color: 'green', deskripsi: 'Tugas mutasi siswa A diselesaikan', waktu: '10:30' },
  { id: 2, icon: 'sync', color: 'blue', deskripsi: 'Sinkronisasi Dapodik berhasil', waktu: '09:15' },
  { id: 3, icon: 'block', color: 'red', deskripsi: 'Pengajuan biodata B ditolak', waktu: '08:00' }
])
</script>

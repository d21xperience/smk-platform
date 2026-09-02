<template>
  <q-page class="q-pa-md bg-grey-2">
    <q-inner-loading :showing="loading" />

    <!-- Header Dashboard-->
    <div class="row q-col-gutter-sm items-center q-mb-md">
      <div class="col-12 col-md-6">
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Dashboard Kesiswaan</h1>
        <p class="text-caption text-grey-7 q-mb-none">Selamat datang di Panel Administrasi TU Kesiswaan SMK</p>
      </div>
      <div class="col-12 col-md-6 text-right text-grey-6 text-caption">
        <q-icon name="calendar_today" class="q-mr-xs" />{{ currentDate }}
      </div>
    </div>

    <!-- 1. Row Metrik Utama -->
    <div class="row q-col-gutter-md q-mb-lg">
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="bg-white text-dark shadow-1 rounded-borders">
          <q-card-section class="flex justify-between items-center">
            <div>
              <div class="text-caption text-grey-7 text-uppercase text-weight-medium">Total Siswa Aktif</div>
              <div class="text-h4 text-weight-bold text-primary q-my-xs">{{ metrics.totalSiswa || 0 }}</div>
            </div>
            <q-avatar color="blue-1" text-color="primary" icon="groups" size="50px" />
          </q-card-section>
        </q-card>
      </div>
      <!-- (Card lainnya mengikuti pola yang sama, disederhanakan untuk ringkasan) -->
    </div>

    <!-- 2. Row Grafik Analisis -->
    <div class="row q-col-gutter-md q-mb-lg">
      <div class="col-12 col-md-7">
        <q-card class="shadow-1 fit">
          <q-card-section>
            <div class="text-subtitle1 text-weight-bold text-grey-9">Distribusi Siswa Per Kompetensi Keahlian</div>
          </q-card-section>
          <q-card-section>
            <apexchart type="bar" height="280" :options="chartOptions" :series="chartSeries"></apexchart>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-md-5">
        <q-card class="shadow-1 fit">
          <q-card-section class="bg-amber-1 text-amber-10 flex items-center">
            <q-icon name="warning" size="sm" class="q-mr-sm" />
            <div>
              <div class="text-weight-bold">Perhatian Kesiswaan (Jatuh Tempo)</div>
            </div>
          </q-card-section>
          <q-separator />
          <q-card-section class="q-pa-none">
            <q-list separator>
              <q-item v-for="(reminder, idx) in reminders" :key="idx" class="q-py-md">
                <q-item-section avatar>
                  <q-avatar color="grey-2" icon="corporate_fare" text-color="grey-8" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-medium">{{ reminder.perusahaan }}</q-item-label>
                  <q-item-label caption>{{ reminder.jurusan }} • {{ reminder.jumlahSiswa }} Siswa</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-badge color="red" :label="reminder.sisaHari + ' Hari Lagi'" />
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
import { ref } from 'vue';
import apexchart from 'vue3-apexcharts';
import { useDashboard } from '@/composables/student/useDashboard';

const currentDate = ref(new Date().toLocaleDateString('id-ID', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' }));

const { metrics, chartSeries, reminders, loading } = useDashboard();

const chartOptions = ref({
  chart: { id: 'siswa-jurusan-chart', toolbar: { show: false } },
  colors: ['#027BE3'],
  plotOptions: { bar: { borderRadius: 4, horizontal: false, columnWidth: '45%' } },
  dataLabels: { enabled: false },
  xaxis: { categories: ['RPL', 'TKRO', 'AKL', 'ULP', 'DPIB'] }
});

// const tableColumns = [
//   { name: 'nama', align: 'left', label: 'Nama Siswa', field: 'nama', sortable: true },
//   { name: 'kelas', align: 'left', label: 'Kelas', field: 'kelas' },
//   { name: 'keperluan', align: 'left', label: 'Keperluan', field: 'keperluan' },
//   { name: 'status', align: 'center', label: 'Status', field: 'status' },
//   { name: 'aksi', align: 'center', label: 'Aksi' }
// ];
</script>

<style scoped>
.text-h4 { line-height: 1.1; }
</style>

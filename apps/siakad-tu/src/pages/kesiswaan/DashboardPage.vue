<template>
  <q-page class="q-pa-md bg-grey-2">
    <q-inner-loading :showing="loading" />

    <!-- Header Dashboard -->
    <div class="row q-col-gutter-sm items-center q-mb-md">
      <div class="col-12">
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Dashboard Kesiswaan</h1>
        <p class="text-caption text-grey-7 q-mb-none">Ringkasan data kesiswaan semester aktif</p>
      </div>
    </div>

    <!-- 1. Row Metrik Utama (4 Kartu) -->
    <div class="row q-col-gutter-md q-mb-lg">
      <!-- Total Siswa -->
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="bg-white text-dark shadow-1 rounded-borders">
          <q-card-section class="flex justify-between items-center">
            <div>
              <div class="text-caption text-grey-7 text-uppercase text-weight-medium">Total Siswa</div>
              <div class="text-h4 text-weight-bold text-primary q-my-xs">{{ metrics.totalSiswa || 0 }}</div>
            </div>
            <q-avatar color="blue-1" text-color="primary" icon="groups" size="50px" />
          </q-card-section>
        </q-card>
      </div>

      <!-- Total Guru -->
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="bg-white text-dark shadow-1 rounded-borders">
          <q-card-section class="flex justify-between items-center">
            <div>
              <div class="text-caption text-grey-7 text-uppercase text-weight-medium">Total Guru & Staf</div>
              <div class="text-h4 text-weight-bold text-teal q-my-xs">{{ metrics.totalGuru || 0 }}</div>
            </div>
            <q-avatar color="teal-1" text-color="teal" icon="badge" size="50px" />
          </q-card-section>
        </q-card>
      </div>

      <!-- Total Rombel -->
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="bg-white text-dark shadow-1 rounded-borders">
          <q-card-section class="flex justify-between items-center">
            <div>
              <div class="text-caption text-grey-7 text-uppercase text-weight-medium">Total Rombel</div>
              <div class="text-h4 text-weight-bold text-indigo q-my-xs">{{ metrics.totalRombel || 0 }}</div>
            </div>
            <q-avatar color="indigo-1" text-color="indigo" icon="meeting_room" size="50px" />
          </q-card-section>
        </q-card>
      </div>

      <!-- Total Jurusan -->
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="bg-white text-dark shadow-1 rounded-borders">
          <q-card-section class="flex justify-between items-center">
            <div>
              <div class="text-caption text-grey-7 text-uppercase text-weight-medium">Kompetensi Keahlian</div>
              <div class="text-h4 text-weight-bold text-orange q-my-xs">{{ metrics.totalJurusan || 0 }}</div>
            </div>
            <q-avatar color="orange-1" text-color="orange" icon="school" size="50px" />
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- 2. Row Grafik Analisis -->
    <div class="row q-col-gutter-md q-mb-lg">
      <div class="col-12 col-md-8">
        <q-card class="shadow-1 fit">
          <q-card-section>
            <div class="text-subtitle1 text-weight-bold text-grey-9">Distribusi Siswa Per Kompetensi Keahlian</div>
          </q-card-section>
          <q-card-section>
            <apexchart type="bar" height="300" :options="chartOptions" :series="chartSeries"></apexchart>
          </q-card-section>
        </q-card>
      </div>
      
      <div class="col-12 col-md-4">
        <!-- Placeholder untuk Widget Lain -->
        <q-card class="shadow-1 fit bg-blue-1 text-blue-9 flex flex-center" style="min-height: 300px">
           <div class="text-center">
             <q-icon name="insights" size="48px" class="q-mb-sm" />
             <div class="text-h6">Statistik Lanjutan</div>
             <div class="text-caption">Segera Hadir</div>
           </div>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import apexchart from 'vue3-apexcharts'
import { useDashboard } from '@/composables/student/useDashboard'
import { useStudentAffairsStore } from '@/stores/kesiswaan/studentAffairsStore'

const { metrics, chartSeries, loading } = useDashboard()
const store = useStudentAffairsStore()

// Computed untuk opsi ApexCharts agar dinamis
const chartOptions = computed(() => ({
  chart: { id: 'siswa-jurusan-chart', toolbar: { show: false } },
  colors: ['#027BE3'],
  plotOptions: { bar: { borderRadius: 4, horizontal: false, columnWidth: '45%' } },
  dataLabels: { enabled: true },
  xaxis: { 
    categories: store.dashboardData?.chartCategories || ['Belum Ada Data'] 
  }
}))
</script>

<style scoped>
.text-h4 { line-height: 1.1; }
</style>

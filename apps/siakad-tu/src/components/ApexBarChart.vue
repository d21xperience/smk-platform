<template>
  <q-card flat bordered class="full-height">
    <q-card-section class="q-pb-none text-center flex-shrink-0">
      {{ title }}
    </q-card-section>
    <!-- Komponen Utama ApexCharts -->
    <apexchart type="bar" :options="chartOptions" :series="chartSeries" height="500px" />
  </q-card>
</template>

<script setup>
import { computed } from 'vue'
import apexchart from 'vue3-apexcharts'
// Definisikan Props agar data bersifat dinamis
const props = defineProps({
  title: { type: String, default: 'CHART TITLE' },
  labels: { type: Array, required: true },     // Contoh: ['2021', '2022', ...]
  lulusData: { type: Array, required: true },  // Contoh: [224, 200, ...]
  masukData: { type: Array, required: true }   // Contoh: [205, 205, ...]
})

// Format data untuk series ApexCharts
const chartSeries = computed(() => [
  { name: 'LULUS', data: props.lulusData },
  { name: 'MASUK', data: props.masukData }
])

// Konfigurasi Chart disesuaikan mirip dengan gambar sampel
const chartOptions = computed(() => ({
  chart: {
    type: 'bar',
    toolbar: { show: false } // Menghilangkan menu download bawaan ApexCharts
  },
  colors: ['#3b6bc6', '#e87e25'], // Warna Biru dan Oranye kustom sesuai gambar
  plotOptions: {
    bar: {
      horizontal: false,
      columnWidth: '70%',      // Mengatur lebar batang utama
      dataLabels: {
        position: 'top'        // Memunculkan angka angka di atas batang
      }
    }
  },
  dataLabels: {
    enabled: true,
    offsetX: 0,
    offsetY: -20,              // Mengatur posisi angka agar menggantung di atas batang
    style: {
      fontSize: '11px',
      colors: ['#666666']      // Warna teks angka label
    }
  },
  stroke: {
    show: true,
    width: 0,                  // Menghilangkan border putih di sekitar batang agar rapat
    colors: ['transparent']
  },
  legend: {
    position: 'top',
    horizontalAlign: 'center',
    markers: { shape: 'square', radius: 0 }, // Indikator legend berbentuk kotak murni
    itemMargin: { horizontal: 10, vertical: 5 },
    style: { fontWeight: 'bold' }
  },
  grid: {
    show: false               // Menghilangkan garis grid latar belakang horizontal/vertikal
  },
  xaxis: {
    categories: props.labels,
    axisTicks: { show: false }, // Menghilangkan centang kecil pada sumbu X
    axisBorder: { color: '#cccccc' }
  },
  yaxis: {
    min: 0,
    max: 400,
    tickAmount: 8,            // Rentang kelipatan 50 dari 0 hingga 400 (400 / 8 = 50)
    axisBorder: { show: true, color: '#cccccc' }
  }
}))
</script>

<style scoped>
.full-height {
  margin-top: 20px;
  /* height: 100vh; */
  overflow: hidden;
  border: 1px solid blue;
}
</style>

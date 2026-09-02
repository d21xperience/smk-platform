<template>
  <q-card flat bordered class="q-mb-lg">
    <q-card-section>
      <!-- Judul Grafik -->
      <div class="text-h6 text-weight-bold text-center text-grey-8 q-mb-sm">
        {{ title }}
      </div>

      <!-- Komponen Utama ApexCharts -->
      <apexchart type="bar" height="100%" :options="chartOptions" :series="chartSeries" />
    </q-card-section>
  </q-card>
</template>

<script setup>
import { computed } from 'vue'
import apexchart from 'vue3-apexcharts'

const props = defineProps({
  title: { type: String, default: 'GRAFIK KELAS' },
  labels: { type: Array, required: true },    // Data Nama Kelas
  lakiData: { type: Array, required: true },  // Data Jumlah Laki-laki
  perempuanData: { type: Array, required: true } // Data Jumlah Perempuan
})

// Format data untuk series ApexCharts
const chartSeries = computed(() => [
  { name: 'LAKI-LAKI', data: props.lakiData },
  { name: 'PEREMPUAN', data: props.perempuanData }
])

const chartOptions = computed(() => ({
  chart: {
    type: 'bar',
    stacked: true, // Menumpuk batang laki-laki dan perempuan
    toolbar: { show: false }
  },
  colors: ['#3b6bc6', '#e87e25'], // Warna Biru dan Oranye
  plotOptions: {
    bar: {
      horizontal: false,
      columnWidth: '60%',
      dataLabels: {
        total: {
          enabled: true, // Menampilkan angka TOTAL di paling atas batang ditumpuk
          style: {
            fontSize: '12px',
            fontWeight: 'bold',
            colors: ['#000000'] // Warna teks TOTAL hitam tegas
          }
        }
      }
    }
  },
  dataLabels: {
    enabled: true, // Menampilkan angka detail di dalam/atas masing-masing potongan gender
    style: {
      fontSize: '11px',
      colors: ['#fff'] // Warna teks angka gender putih agar kontras di dalam batang
    }
  },
  stroke: {
    show: true,
    width: 1,
    colors: ['#fff']
  },
  legend: {
    position: 'top',
    horizontalAlign: 'center',
    markers: { shape: 'square', radius: 0 }
  },
  grid: {
    show: false // Menghilangkan garis latar belakang sesuai sampel sebelumnya
  },
  xaxis: {
    categories: props.labels,
    axisTicks: { show: false },
    axisBorder: { color: '#cccccc' }
  },
  yaxis: {
    min: 0,
    max: 50, // Disesuaikan dengan total maksimum kelas Anda (nilai tertinggi 41)
    tickAmount: 5,
    axisBorder: { show: true, color: '#cccccc' }
  }
}))
</script>

<style scoped>
.chart-card {
  width: 100%;
  max-width: 850px;
  margin: 0 auto;
}
</style>

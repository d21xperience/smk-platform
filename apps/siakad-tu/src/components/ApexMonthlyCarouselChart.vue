<template>
  <q-card flat bordered class="chart-card">
    <!-- Judul Utama -->
    <q-card-section class="q-pb-none text-center">
      <div class="text-h6 text-weight-bold">{{ title }}</div>
    </q-card-section>

    <!-- Carousel Per Tingkat Kelas -->
    <q-carousel v-model="activeSlide" transition-prev="slide-right" transition-next="slide-left" swipeable animated
      control-color="primary"  padding height="480px" class="bg-transparent" :autoplay=5000 infinite>
      <q-carousel-slide v-for="tingkat in daftarTingkat" :key="tingkat" :name="tingkat"
        class="column no-wrap flex-center q-pa-md">
        <!-- Sub-Judul Slide -->
        <!-- <div class="text-subtitle1 text-weight-bold text-primary q-mb-md">
          TREN JUMLAH SISWA TINGKAT {{ tingkat }} (SEMESTER 1)
        </div> -->

        <!-- Grafik Garis ApexCharts -->
        <div class="full-width">
          <apexchart type="line" height="360" :options="getChartOptions(tingkat)" :series="getChartSeries(tingkat)" />
        </div>
      </q-carousel-slide>
    </q-carousel>
  </q-card>
</template>

<script setup>
import { ref, computed } from 'vue'
import apexchart from 'vue3-apexcharts'

const props = defineProps({
  title: { type: String, default: '' },
  rawMonthlyData: { type: Array, required: true }
})

const activeSlide = ref('X')

// List bulan semester 1 sebagai sumbu X
const listBulan = ['JULI', 'AGUSTUS', 'SEPTEMBER', 'OKTOBER', 'NOVEMBER', 'DESEMBER']

const daftarTingkat = computed(() => {
  return [...new Set(props.rawMonthlyData.map(item => item.tingkat))]
})

// Fungsi memformat series data per kelas di tingkat tertentu
const getChartSeries = (tingkat) => {
  const filtered = props.rawMonthlyData.filter(item => item.tingkat === tingkat)

  // Mengubah data baris kelas menjadi baris tren garis garis
  return filtered.map(item => ({
    name: item.nama,
    data: [item.juli, item.agustus, item.september, item.oktober, item.november, item.desember]
  }))
}

const getChartOptions = () => {
  return {
    chart: {
      type: 'line',
      toolbar: { show: false },
      dropShadow: {
        enabled: true,
        color: '#000',
        top: 3,
        left: 3,
        blur: 3,
        opacity: 0.1
      }
    },
    // Ketebalan garis dan kehalusan lekukan (smooth)
    stroke: {
      chart: 'smooth',
      width: 3
    },
    // Penanda bulatan angka di setiap bulan
    markers: {
      size: 4,
      hover: { size: 6 }
    },
    dataLabels: {
      enabled: true, // Menampilkan angka siswa langsung di atas garis bulanan
      style: { fontSize: '10px' }
    },
    grid: {
      borderColor: '#e7e7e7',
      row: {
        colors: ['#f3f3f3', 'transparent'], // Baris latar selang-seling agar mudah dilihat
        opacity: 0.5
      }
    },
    xaxis: {
      categories: listBulan,
      // title: { text: 'Bulan Semester 1', style: { color: '#666' } }
    },
    yaxis: {
      min: 0,
      max: 50, // Nilai disesuaikan data tabel Anda (Maksimal di angka 45)
      title: { text: 'Jumlah Siswa' }
    },
    legend: {
      position: 'bottom',
      horizontalAlign: 'center'
    }
  }
}
</script>

<style scoped>
.chart-card {
  width: 100%;
  max-width: 900px;
  margin: 0 auto;
}
</style>

<template>
  <q-page class="q-pa-md bg-grey-2">
    <div class="row justify-between items-center q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Statistik Keberadaan GTK</h1>
        <p class="text-caption text-grey-7 q-mb-none">Pusat Data & Analitik Kepegawaian</p>
      </div>
      <q-btn color="primary" icon="refresh" label="Refresh Data" @click="fetchData" :loading="isLoading" />
    </div>

    <q-linear-progress v-if="isLoading" indeterminate color="primary" class="q-mb-md" />

    <div class="row q-col-gutter-md">
      <!-- KPI Cards -->
      <div class="col-12 col-sm-6 col-md-4 col-lg-2" v-for="(kpi, index) in kpiList" :key="index">
        <q-card class="shadow-1 h-full text-center q-pt-md" :class="`bg-${kpi.color}-1`">
          <q-icon :name="kpi.icon" :color="kpi.color" size="32px" />
          <div class="text-h4 text-weight-bold q-mt-sm" :class="`text-${kpi.color}-9`">{{ kpi.value }}</div>
          <div class="text-caption text-grey-8 q-mb-md">{{ kpi.title }}</div>
        </q-card>
      </div>

      <!-- Charts Row -->
      <div class="col-12 col-md-6">
        <q-card class="shadow-1 q-pa-md">
          <div class="text-subtitle1 text-weight-bold text-primary q-mb-md">Komposisi Status Kepegawaian</div>
          <VueApexCharts type="donut" height="300" :options="donutOptions" :series="donutSeries"></VueApexCharts>
        </q-card>
      </div>

      <div class="col-12 col-md-6">
        <q-card class="shadow-1 q-pa-md">
          <div class="text-subtitle1 text-weight-bold text-primary q-mb-md">Distribusi Unit Kerja</div>
          <VueApexCharts type="bar" height="300" :options="barOptions" :series="barSeries"></VueApexCharts>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useQuasar } from 'quasar'
import { api } from '@/boot/axios'
import VueApexCharts from 'vue3-apexcharts'

const $q = useQuasar()
const isLoading = ref(false)
const stats = ref({
  kpi: {
    totalPegawai: 0,
    totalHadir: 0,
    totalCuti: 0,
    totalAlpha: 0,
    avgSkp: 0,
    guruBebanKurang: 0
  },
  charts: {
    komposisi: [],
    distribusiUnit: []
  }
})

const kpiList = computed(() => [
  { title: 'Total GTK', value: stats.value.kpi.totalPegawai, icon: 'groups', color: 'blue' },
  { title: 'Hadir (Bulan Ini)', value: stats.value.kpi.totalHadir, icon: 'how_to_reg', color: 'green' },
  { title: 'Cuti Disetujui', value: stats.value.kpi.totalCuti, icon: 'event_busy', color: 'orange' },
  { title: 'Alpha (Bulan Ini)', value: stats.value.kpi.totalAlpha, icon: 'person_off', color: 'red' },
  { title: 'Rata-rata SKP', value: (stats.value.kpi.avgSkp || 0).toFixed(2), icon: 'trending_up', color: 'teal' },
  { title: 'Beban < 24 JTM', value: stats.value.kpi.guruBebanKurang, icon: 'warning', color: 'deep-orange' }
])

// ApexCharts Options
const donutOptions = ref({
  labels: [],
  responsive: [{ breakpoint: 480, options: { chart: { width: 200 } } }],
  legend: { position: 'bottom' }
})
const donutSeries = ref([])

const barOptions = ref({
  chart: { type: 'bar' },
  plotOptions: { bar: { horizontal: false, columnWidth: '55%', borderRadius: 4 } },
  dataLabels: { enabled: false },
  xaxis: { categories: [] },
  colors: ['#1976D2']
})
const barSeries = ref([])

const fetchData = async () => {
  isLoading.value = true
  try {
    const res = await api.get('/kepegawaian/statistik')
    stats.value = res.data

    // Map Data Donut Chart
    donutOptions.value = { ...donutOptions.value, labels: res.data.charts.komposisi.map(i => i.label) }
    donutSeries.value = res.data.charts.komposisi.map(i => i.value)

    // Map Data Bar Chart
    barOptions.value = { ...barOptions.value, xaxis: { categories: res.data.charts.distribusiUnit.map(i => i.label) } }
    barSeries.value = [{ name: 'Jumlah GTK', data: res.data.charts.distribusiUnit.map(i => i.value) }]

  } catch (err) {
    $q.notify({ type: 'negative', message: 'Gagal memuat data statistik' })
    console.error(err)
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.h-full {
  height: 100%;
}
</style>

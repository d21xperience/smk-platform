<template>
  <q-page class="q-pa-md bg-grey-2">
    <div class="row justify-between items-center q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Statistik Keberadaan GTK</h1>
        <p class="text-caption text-grey-7 q-mb-none">Pusat Data & Analitik Kepegawaian</p>
      </div>
      <q-btn color="primary" icon="refresh" label="Terapkan Filter" @click="fetchData" :loading="isLoading" />
    </div>

    <!-- Filter Bar -->
    <q-card class="shadow-1 q-mb-md">
      <q-card-section class="row items-center q-col-gutter-md">
        <div class="col-12 col-md-3">
          <q-input v-model="filter.startDate" type="date" label="Tanggal Mulai" outlined dense />
        </div>
        <div class="col-12 col-md-3">
          <q-input v-model="filter.endDate" type="date" label="Tanggal Akhir" outlined dense />
        </div>
        <div class="col-12 col-md-6 text-right">
          <q-btn color="secondary" icon="download" label="Export PDF Snapshot" outline />
        </div>
      </q-card-section>
    </q-card>

    <!-- KPI Cards -->
    <div class="row q-col-gutter-md q-mb-lg">
      <!-- Total Pegawai -->
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="bg-primary text-white shadow-2">
          <q-card-section>
            <div class="row items-center justify-between">
              <div>
                <div class="text-subtitle2 text-grey-3">Total Pegawai</div>
                <div class="text-h3 text-weight-bold">{{ stats.kpis?.totalPegawai || 0 }}</div>
              </div>
              <q-icon name="groups" size="48px" class="text-white opacity-50" />
            </div>
          </q-card-section>
        </q-card>
      </div>

      <!-- Rata Kehadiran -->
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="bg-green text-white shadow-2">
          <q-card-section>
            <div class="row items-center justify-between">
              <div>
                <div class="text-subtitle2 text-grey-3">Rata-rata Kehadiran</div>
                <div class="text-h3 text-weight-bold">{{ (stats.kpis?.rataKehadiran || 0).toFixed(2) }}%</div>
              </div>
              <q-icon name="fact_check" size="48px" class="text-white opacity-50" />
            </div>
          </q-card-section>
        </q-card>
      </div>

      <!-- Cuti Pending -->
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="bg-orange text-white shadow-2">
          <q-card-section>
            <div class="row items-center justify-between">
              <div>
                <div class="text-subtitle2 text-grey-3">Cuti Pending</div>
                <div class="text-h3 text-weight-bold">{{ stats.kpis?.cutiPending || 0 }}</div>
              </div>
              <q-icon name="pending_actions" size="48px" class="text-white opacity-50" />
            </div>
          </q-card-section>
        </q-card>
      </div>

      <!-- Total Alpha -->
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="bg-red text-white shadow-2">
          <q-card-section>
            <div class="row items-center justify-between">
              <div>
                <div class="text-subtitle2 text-grey-3">Total Alpha</div>
                <div class="text-h3 text-weight-bold">{{ stats.kpis?.totalAlpha || 0 }}</div>
              </div>
              <q-icon name="warning" size="48px" class="text-white opacity-50" />
            </div>
          </q-card-section>
        </q-card>
      </div>

      <!-- Rata SKP -->
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="bg-purple text-white shadow-2">
          <q-card-section>
            <div class="row items-center justify-between">
              <div>
                <div class="text-subtitle2 text-grey-3">Rata-rata SKP</div>
                <div class="text-h3 text-weight-bold">{{ (stats.kpis?.rataSkp || 0).toFixed(2) }}</div>
              </div>
              <q-icon name="trending_up" size="48px" class="text-white opacity-50" />
            </div>
          </q-card-section>
        </q-card>
      </div>

      <!-- Dokumen Kadaluarsa -->
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="bg-brown text-white shadow-2">
          <q-card-section>
            <div class="row items-center justify-between">
              <div>
                <div class="text-subtitle2 text-grey-3">Dokumen Kadaluarsa</div>
                <div class="text-h3 text-weight-bold">{{ stats.kpis?.dokumenKadaluarsa || 0 }}</div>
              </div>
              <q-icon name="description" size="48px" class="text-white opacity-50" />
            </div>
          </q-card-section>
        </q-card>
      </div>

      <!-- Beban Kurang 24 JTM -->
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="bg-teal text-white shadow-2">
          <q-card-section>
            <div class="row items-center justify-between">
              <div>
                <div class="text-subtitle2 text-grey-3">Guru Beban &lt; 24 JTM</div>
                <div class="text-h3 text-weight-bold">{{ stats.kpis?.bebanKurang24 || 0 }}</div>
              </div>
              <q-icon name="schedule" size="48px" class="text-white opacity-50" />
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Charts Section -->
    <div class="row q-col-gutter-md">
      <!-- Pie Chart: Komposisi Pegawai -->
      <div class="col-12 col-md-6">
        <q-card class="shadow-1">
          <q-card-section>
            <div class="text-h6 text-primary q-mb-md">Komposisi Status Pegawai</div>
            <VueApexCharts type="pie" height="300" :options="pieOptions" :series="pieSeries"></VueApexCharts>
          </q-card-section>
        </q-card>
      </div>

      <!-- Bar Chart: Distribusi Unit Kerja -->
      <div class="col-12 col-md-6">
        <q-card class="shadow-1">
          <q-card-section>
            <div class="text-h6 text-primary q-mb-md">Distribusi Unit Kerja</div>
            <VueApexCharts type="bar" height="300" :options="barOptions" :series="barSeries"></VueApexCharts>
          </q-card-section>
        </q-card>
      </div>

      <!-- Horizontal Bar: Sebaran Pangkat -->
      <div class="col-12">
        <q-card class="shadow-1">
          <q-card-section>
            <div class="text-h6 text-primary q-mb-md">Sebaran Pangkat / Golongan</div>
            <VueApexCharts type="bar" height="300" :options="horizontalBarOptions" :series="horizontalBarSeries">
            </VueApexCharts>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useQuasar } from 'quasar'
import { api } from '@/boot/axios'
import VueApexCharts from 'vue3-apexcharts' // <-- Import ini sekarang akan terpakai di template

const $q = useQuasar()
const isLoading = ref(false)

const filter = ref({
  startDate: new Date(new Date().getFullYear(), new Date().getMonth(), 1).toISOString().split('T')[0], // Awal bulan
  endDate: new Date().toISOString().split('T')[0] // Hari ini
})

const stats = ref({
  kpis: {},
  charts: {
    komposisiPegawai: [],
    distribusiUnit: [],
    sebaranPangkat: []
  }
})

// === PIE CHART CONFIG ===
const pieSeries = computed(() => stats.value.charts?.komposisiPegawai?.map(item => item.value) || [])
const pieOptions = computed(() => ({
  labels: stats.value.charts?.komposisiPegawai?.map(item => item.label) || [],
  legend: { position: 'bottom' },
  responsive: [{ breakpoint: 480, options: { chart: { width: 200 } } }]
}))

// === BAR CHART CONFIG ===
const barSeries = computed(() => [{
  name: 'Jumlah Pegawai',
  data: stats.value.charts?.distribusiUnit?.map(item => item.value) || []
}])
const barOptions = computed(() => ({
  xaxis: { categories: stats.value.charts?.distribusiUnit?.map(item => item.label) || [] },
  colors: ['#1976D2'],
  plotOptions: { bar: { borderRadius: 4, horizontal: false } }
}))

// === HORIZONTAL BAR CONFIG ===
const horizontalBarSeries = computed(() => [{
  name: 'Jumlah',
  data: stats.value.charts?.sebaranPangkat?.map(item => item.value) || []
}])
const horizontalBarOptions = computed(() => ({
  xaxis: { categories: stats.value.charts?.sebaranPangkat?.map(item => item.label) || [] },
  colors: ['#388E3C'],
  plotOptions: { bar: { borderRadius: 4, horizontal: true } }
}))

const fetchData = async () => {
  isLoading.value = true
  try {
    const res = await api.get('/kepegawaian/statistics/dashboard', {
      params: {
        start_date: filter.value.startDate,
        end_date: filter.value.endDate
      }
    })
    stats.value = res.data
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
.opacity-50 {
  opacity: 0.5;
}
</style>

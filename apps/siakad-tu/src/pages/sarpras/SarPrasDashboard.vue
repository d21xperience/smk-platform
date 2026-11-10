<template>
  <q-page class="q-pa-md">
    <div class="row q-mb-md items-center justify-between">
      <div class="text-h6 text-weight-bold">Dashboard Sarpras</div>
    </div>

    <!-- METRIK CARDS -->
    <div class="row q-col-gutter-md q-mb-md">
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="bg-blue-9 text-white">
          <q-card-section>
            <div class="text-caption opacity-80">Total Aset Aktif</div>
            <div class="text-h4 text-weight-bold">{{ stats.total_assets }} <span class="text-subtitle2">Unit</span>
            </div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="bg-teal-8 text-white">
          <q-card-section>
            <div class="text-caption opacity-80">Total Nilai Aset</div>
            <div class="text-h5 text-weight-bold">Rp {{ (stats.total_value / 1000000).toFixed(1) }} <span
                class="text-subtitle2">Juta</span></div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="bg-orange-8 text-white">
          <q-card-section>
            <div class="text-caption opacity-80">Peminjaman Aktif</div>
            <div class="text-h4 text-weight-bold">{{ stats.active_borrowings }} <span
                class="text-subtitle2">Transaksi</span></div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="bg-red-8 text-white">
          <q-card-section>
            <div class="text-caption opacity-80">Usulan Dihapus</div>
            <div class="text-h4 text-weight-bold">{{ stats.pending_disposals }} <span class="text-subtitle2">Aset</span>
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- CHARTS & LAPORAN -->
    <div class="row q-col-gutter-md">
      <!-- KONDISI ASET (PIE CHART) -->
      <div class="col-12 col-md-6">
        <q-card flat bordered>
          <q-card-section class="text-subtitle1 text-weight-bold">Kondisi Aset</q-card-section>
          <q-card-section>
            <!-- Perbaikan: Gunakan komponen VueApexCharts -->
            <apexchart type="pie" height="300" :options="chartOptions" :series="chartSeries" />
          </q-card-section>
        </q-card>
      </div>

      <!-- LAPORAN KOMPREHENSIF -->
      <div class="col-12 col-md-6">
        <q-card flat bordered class="column full-height">
          <q-card-section class="text-subtitle1 text-weight-bold bg-grey-2">Generate Laporan</q-card-section>
          <q-card-section class="col q-gutter-y-md">
            <q-banner class="bg-blue-1 text-blue-9 rounded-borders">
              Pilih rentang tanggal untuk mengunduh Laporan Komprehensif (Penerimaan, Peminjaman, Ringkasan Aset) dalam
              format PDF.
            </q-banner>
            <div class="row q-col-gutter-sm">
              <div class="col-6">
                <q-input outlined dense v-model="reportForm.startDate" type="date" label="Dari Tanggal" />
              </div>
              <div class="col-6">
                <q-input outlined dense v-model="reportForm.endDate" type="date" label="Sampai Tanggal" />
              </div>
            </div>
            <q-btn color="primary" icon="picture_as_pdf" label="Unduh Laporan PDF" @click="unduhLaporan"
              :loading="isDownloading" class="full-width" />
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useQuasar } from 'quasar'
import { useSarprasDashboardStore } from '@/stores/sarpras/sarprasDashboard'
import { useContextStore } from '@/stores/contextStore'
// PERBAIKAN: Import dengan nama 'apexchart' agar cocok dengan tag <apexchart> di template
import apexchart from 'vue3-apexcharts'

const $q = useQuasar()
const store = useSarprasDashboardStore()
const contextStore = useContextStore()

// ApexCharts setup
const stats = computed(() => store.stats)
const chartSeries = computed(() => [stats.value.condition_baik, stats.value.condition_rr, stats.value.condition_rb])
const chartOptions = {
  labels: ['Baik', 'Rusak Ringan', 'Rusak Berat'],
  colors: ['#21ba45', '#f2c037', '#c10015'],
  legend: { position: 'bottom' }
}

// Laporan Form
const isDownloading = ref(false)
const reportForm = ref({
  startDate: new Date(new Date().getFullYear(), new Date().getMonth(), 1).toISOString().split('T')[0],
  endDate: new Date().toISOString().split('T')[0]
})

onMounted(() => {
  store.fetchStats()
})

watch(() => contextStore.selectedSemesterId, (newVal) => {
  if (newVal) store.fetchStats()
})

async function unduhLaporan() {
  if (!reportForm.value.startDate || !reportForm.value.endDate) {
    $q.notify({ type: 'negative', message: 'Tanggal wajib diisi' })
    return
  }
  isDownloading.value = true
  try {
    await store.downloadComprehensiveReport(reportForm.value.startDate, reportForm.value.endDate)
    $q.notify({ type: 'positive', message: 'Laporan berhasil diunduh' })
  } catch (e) {
    $q.notify({ type: 'negative', message: `Gagal mengunduh laporan ${e}` })
  } finally {
    isDownloading.value = false
  }
}
</script>

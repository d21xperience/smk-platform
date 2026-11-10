<template>
  <q-page class="q-pa-md bg-grey-2">
    <div class="row justify-between items-center q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Buku Kontrol Kenaikan</h1>
        <p class="text-caption text-grey-7 q-mb-none">Monitoring Kenaikan Gaji, Pangkat, dan Pensiun Pegawai</p>
      </div>
      <div class="row q-gutter-sm">
        <q-btn color="secondary" icon="autorenew" label="Inisialisasi Baru" @click="initMassal" :loading="isInitializing" outline />
        <q-btn color="primary" icon="refresh" label="Refresh" @click="fetchData" :loading="isLoading" />
      </div>
    </div>

    <q-linear-progress v-if="isLoading" indeterminate color="primary" class="q-mb-md" />

    <!-- KPI Cards -->
    <div class="row q-col-gutter-md q-mb-md">
      <div class="col-6 col-md-2" v-for="kpi in kpiList" :key="kpi.label">
        <q-card class="text-center q-pa-sm" :class="`bg-${kpi.color}-1`">
          <div class="text-h5 text-weight-bold" :class="`text-${kpi.color}-9`">{{ kpi.value }}</div>
          <div class="text-caption text-grey-8">{{ kpi.label }}</div>
        </q-card>
      </div>
    </div>

    <!-- Table -->
    <q-card class="shadow-1">
      <q-card-section>
        <q-table 
          :rows="bukuKontrolList" 
          :columns="columns" 
          row-key="ptkId" 
          flat 
          bordered 
          dense 
          :loading="isLoading"
        >
          <template v-slot:body-cell-namaPegawai="props">
            <q-td :props="props">
              <div class="text-weight-medium">{{ props.row.namaPegawai }}</div>
              <div class="text-caption text-grey-7">{{ props.row.nip }}</div>
            </q-td>
          </template>
          <template v-slot:body-cell-pangkat="props">
            <q-td :props="props" class="text-center">
              <q-badge color="blue-6" class="q-pa-xs">{{ props.row.kodeGolongan }}</q-badge>
              <div class="text-caption q-mt-xs">{{ props.row.namaPangkat }}</div>
            </q-td>
          </template>
          <template v-slot:body-cell-tanggalKgb="props">
            <q-td :props="props" class="text-center">
              {{ formatDate(props.row.tanggalKgbBerikut) }}
            </q-td>
          </template>
          <template v-slot:body-cell-tanggalPangkat="props">
            <q-td :props="props" class="text-center">
              {{ formatDate(props.row.tanggalPangkatBerikut) }}
            </q-td>
          </template>
          <template v-slot:body-cell-tanggalPensiun="props">
            <q-td :props="props" class="text-center">
              {{ formatDate(props.row.tanggalPensiun) }}
            </q-td>
          </template>
          <template v-slot:body-cell-status="props">
            <q-td :props="props" class="text-center">
              <q-badge :color="getStatusColor(props.row.statusMonitoring)" class="q-pa-sm">
                {{ props.row.statusMonitoring }}
              </q-badge>
              <div class="text-caption q-mt-xs" :class="props.row.selisihHari < 0 ? 'text-negative' : 'text-grey-7'">
                {{ props.row.selisihHari }} hari
              </div>
            </q-td>
          </template>
        </q-table>
      </q-card-section>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useQuasar } from 'quasar'
import { api } from '@/boot/axios'

const $q = useQuasar()
const isLoading = ref(false)
const isInitializing = ref(false)
const bukuKontrolList = ref([])
const kpiData = ref({ total: 0, aman: 0, perluPersiapan: 0, segera: 0, urgent: 0, overdue: 0 })

const columns = [
  { name: 'namaPegawai', label: 'Nama / NIP', align: 'left', sortable: true },
  { name: 'pangkat', label: 'Golongan / Pangkat', align: 'center' },
  { name: 'tanggalKgb', label: 'KGB Berikutnya', align: 'center' },
  { name: 'tanggalPangkat', label: 'Pangkat Berikutnya', align: 'center' },
  { name: 'tanggalPensiun', label: 'Pensiun', align: 'center' },
  { name: 'status', label: 'Status Monitoring', align: 'center' }
]

const kpiList = computed(() => [
  { label: 'Total Pegawai', value: kpiData.value.total, color: 'blue' },
  { label: 'Aman', value: kpiData.value.aman, color: 'green' },
  { label: 'Persiapan', value: kpiData.value.perluPersiapan, color: 'yellow' },
  { label: 'Segera', value: kpiData.value.segera, color: 'orange' },
  { label: 'Urgent', value: kpiData.value.urgent, color: 'red' },
  { label: 'Overdue', value: kpiData.value.overdue, color: 'deep-orange' }
])

const formatDate = (val) => {
  if (!val) return '-'
  return new Date(val).toLocaleDateString('id-ID', { year: 'numeric', month: 'short', day: 'numeric' })
}

const getStatusColor = (status) => {
  switch (status) {
    case 'AMAN': return 'green-6'
    case 'PERLU_PERSIAPAN': return 'yellow-8'
    case 'SEGERA': return 'orange-6'
    case 'URGENT': return 'red-6'
    case 'OVERDUE': return 'red-9'
    default: return 'grey-6'
  }
}

const fetchData = async () => {
  isLoading.value = true
  try {
    const res = await api.get('/kepegawaian/buku-kontrol/dashboard')
    bukuKontrolList.value = res.data.table || []
    kpiData.value = res.data.kpi || { total: 0, aman: 0, perluPersiapan: 0, segera: 0, urgent: 0, overdue: 0 }
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Gagal memuat data buku kontrol' })
    console.error(err)
  } finally {
    isLoading.value = false
  }
}

const initMassal = async () => {
  isInitializing.value = true
  try {
    const res = await api.post('/kepegawaian/buku-kontrol/initialize')
    $q.notify({ type: 'positive', message: `Inisialisasi selesai. Total: ${res.data.total}` })
    await fetchData()
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Gagal melakukan inisialisasi' })
    console.error(err)
  } finally {
    isInitializing.value = false
  }
}

onMounted(fetchData)
</script>

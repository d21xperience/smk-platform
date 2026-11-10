<template>
  <q-page class="q-pa-md">
    <!-- Header -->
    <div class="row q-mb-md items-center justify-between">
      <div class="text-h6 text-weight-bold">Buku Perawatan & Perbaikan Aset</div>
    </div>

    <!-- Filter & Tombol Aksi -->
    <q-card flat bordered class="q-mb-md bg-blue-grey-1">
      <q-card-section class="row items-center q-col-gutter-md">
        <div class="col-12 col-sm-auto">
          <q-input outlined dense bg-color="white" v-model="filter.startDate" type="date" label="Dari Tanggal" />
        </div>
        <div class="col-12 col-sm-auto">
          <q-input outlined dense bg-color="white" v-model="filter.endDate" type="date" label="Sampai Tanggal" />
        </div>
        <div class="col-12 col-sm-auto">
          <q-btn color="primary" icon="search" label="Filter" @click="loadData" :loading="store.isLoading" />
        </div>
        <q-space />
        <div class="col-12 col-sm-auto">
          <q-btn color="negative" icon="picture_as_pdf" label="Unduh Buku Perawatan (PDF)" @click="unduhPDF"
            :loading="isDownloading" />
        </div>
      </q-card-section>
    </q-card>

    <!-- Tabel Riwayat -->
    <q-table :rows="store.logs" :columns="columns" row-key="id" flat bordered :loading="store.isLoading"
      :rows-per-page-options="[10, 25, 50, 0]">
      <template v-slot:body-cell-service_date="props">
        <q-td :props="props" class="text-center">{{ new Date(props.value).toLocaleDateString('id-ID') }}</q-td>
      </template>
      <template v-slot:body-cell-cost="props">
        <q-td :props="props" class="text-right">
          Rp {{ props.value.toLocaleString('id-ID') }}
        </q-td>
      </template>
      <template v-slot:body-cell-vendor="props">
        <q-td :props="props">
          <div>{{ props.row.vendor || '-' }}</div>
          <div class="text-caption text-grey">{{ props.row.technician || '' }}</div>
        </q-td>
      </template>
      <template v-slot:body-cell-asset="props">
        <q-td :props="props">
          <div class="text-weight-medium">{{ props.row.asset_name }}</div>
          <div class="text-caption text-grey">{{ props.row.asset_code }}</div>
        </q-td>
      </template>
    </q-table>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useMaintenanceStore } from '@/stores/sarpras/maintenance'

const $q = useQuasar()
const store = useMaintenanceStore()

const isDownloading = ref(false)
const filter = ref({
  startDate: new Date(new Date().getFullYear(), new Date().getMonth(), 1).toISOString().split('T')[0],
  endDate: new Date().toISOString().split('T')[0]
})

const columns = [
  { name: 'service_date', label: 'Tanggal', field: 'service_date', align: 'center' },
  { name: 'asset', label: 'Aset', align: 'left' },
  { name: 'vendor', label: 'Vendor / Teknisi', align: 'left' },
  { name: 'cost', label: 'Biaya', field: 'cost', align: 'right' },
  { name: 'notes', label: 'Catatan', field: 'notes', align: 'left' },
]

onMounted(() => {
  loadData()
})

function loadData() {
  store.fetchMaintenanceLogs(filter.value.startDate, filter.value.endDate)
}

async function unduhPDF() {
  isDownloading.value = true
  try {
    await store.downloadMaintenanceBookPDF(filter.value.startDate, filter.value.endDate)
    $q.notify({ type: 'positive', message: 'Buku Perawatan berhasil diunduh' })
  } catch (e) {
    $q.notify({ type: 'negative', message: `Gagal mengunduh PDF ${e}` })
  } finally {
    isDownloading.value = false
  }
}
</script>

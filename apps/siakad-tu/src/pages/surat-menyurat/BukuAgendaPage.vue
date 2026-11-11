<template>
  <q-page class="q-pa-md">
    <div class="row items-center q-mb-md">
      <div class="text-h5 text-weight-bold">Buku Agenda Surat</div>
      <q-space />
      <q-btn color="positive" icon="download" label="Export Excel" :loading="isExporting" @click="exportExcel" />
    </div>

    <q-card flat bordered class="q-mb-md bg-grey-1">
      <q-card-section class="row q-col-gutter-sm items-center">
        <div class="col-12 col-md-3">
          <q-select v-model="filter.jenis" :options="[{label: 'Semua (IN/OUT)', value: 'BOTH'}, {label: 'Surat Masuk', value: 'IN'}, {label: 'Surat Keluar', value: 'OUT'}]" label="Jenis Buku" outlined dense emit-value map-options @update:model-value="fetchData" />
        </div>
        <div class="col-6 col-md-2">
          <q-select v-model="filter.bulan" :options="bulanOptions" label="Bulan" outlined dense emit-value map-options @update:model-value="fetchData" />
        </div>
        <div class="col-6 col-md-2">
          <q-input v-model.number="filter.tahun" label="Tahun" type="number" outlined dense @update:model-value="fetchData" />
        </div>
      </q-card-section>
    </q-card>

    <q-table :rows="suratStore.agendaList" :columns="columns" row-key="nomor_agenda" flat bordered :loading="suratStore.isLoadingAgenda">
      <template v-slot:body-cell-jenis="props">
        <q-td :props="props">
          <q-badge :color="props.row.jenis === 'IN' ? 'blue-6' : 'orange-6'">{{ props.row.jenis === 'IN' ? 'MASUK' : 'KELUAR' }}</q-badge>
        </q-td>
      </template>
      <template v-slot:body-cell-nomor_agenda="props">
        <q-td :props="props">
          <span class="text-weight-bold text-primary">{{ props.row.nomor_agenda }}</span>
        </q-td>
      </template>
      <template v-slot:body-cell-pihak="props">
        <q-td :props="props">
          {{ props.row.from_party }} <br> <span class="text-caption text-grey">ke: {{ props.row.to_party }}</span>
        </q-td>
      </template>
    </q-table>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Notify } from 'quasar'
import { useSuratStore } from '@/stores/surat/surat'

const suratStore = useSuratStore()

const isExporting = ref(false)
const filter = ref({ jenis: 'BOTH', bulan: new Date().getMonth() + 1, tahun: new Date().getFullYear() })

const bulanOptions = [
  { label: 'Januari', value: 1 }, { label: 'Februari', value: 2 }, { label: 'Maret', value: 3 },
  { label: 'April', value: 4 }, { label: 'Mei', value: 5 }, { label: 'Juni', value: 6 },
  { label: 'Juli', value: 7 }, { label: 'Agustus', value: 8 }, { label: 'September', value: 9 },
  { label: 'Oktober', value: 10 }, { label: 'November', value: 11 }, { label: 'Desember', value: 12 }
]

const columns = [
  { name: 'nomor_agenda', label: 'No. Agenda', align: 'left', sortable: true },
  { name: 'tanggal_catat', label: 'Tgl Catat', field: 'tanggal_catat', align: 'left' },
  { name: 'jenis', label: 'Jenis', align: 'center' },
  { name: 'number', label: 'Nomor Surat', field: 'number', align: 'left' },
  { name: 'date', label: 'Tgl Surat', field: 'date', align: 'left' },
  { name: 'subject', label: 'Perihal', field: 'subject', align: 'left' },
  { name: 'pihak', label: 'Pengirim / Penerima', align: 'left' },
  { name: 'kode_klasifikasi', label: 'Klasifikasi', field: 'kode_klasifikasi', align: 'center' }
]

async function fetchData() {
  try {
    await suratStore.fetchAgenda(filter.value)
  } catch {
    Notify.create({ type: 'negative', message: 'Gagal memuat data agenda' })
  }
}

async function exportExcel() {
  isExporting.value = true
  try {
    const fileUrl = await suratStore.exportAgendaExcel(filter.value)
    if (fileUrl) {
      const fullUrl = `http://localhost:8081${fileUrl}`
      window.open(fullUrl, '_blank')
      Notify.create({ type: 'positive', message: 'File Excel berhasil dibuat' })
    }
  } catch {
    Notify.create({ type: 'negative', message: 'Gagal membuat file Excel' })
  } finally {
    isExporting.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>

<template>
  <q-page class="q-pa-md bg-grey-2">
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Rekapitulasi & Pelaporan Berkala</h1>
        <p class="text-caption text-grey-7 q-mb-none">Penyusunan statistik bulanan siswa aktif untuk Dinas Pendidikan / Yayasan</p>
      </div>
      <div>
        <q-btn-dropdown color="primary" icon="download" label="Unduh Laporan Resmi" :loading="downloading">
          <q-list>
            <q-item clickable v-close-popup @click="downloadPDF">
              <q-item-section avatar><q-icon name="picture_as_pdf" color="red" /></q-item-section>
              <q-item-section><q-item-label>Laporan Bulanan (.pdf)</q-item-label></q-item-section>
            </q-item>
            <q-item clickable v-close-popup @click="downloadExcel">
              <q-item-section avatar><q-icon name="table_view" color="green-8" /></q-item-section>
              <q-item-section><q-item-label>Rekap Statistik Siswa (.xlsx)</q-item-label></q-item-section>
            </q-item>
          </q-list>
        </q-btn-dropdown>
      </div>
    </div>

    <q-banner v-if="store.error" class="bg-red-1 text-red-9 q-mb-md" rounded>
      <q-icon name="error" class="q-mr-sm" /> {{ store.error }}
    </q-banner>

    <q-card class="shadow-1 q-mb-md">
      <q-card-section class="flex justify-between items-center q-py-sm">
        <div class="text-subtitle2 text-weight-bold text-grey-9">
          Rekapitulasi Keadaan Siswa per Jurusan
        </div>
        <q-btn icon="refresh" flat color="primary" :loading="store.loading" @click="loadData">
          <q-tooltip>Muat Ulang</q-tooltip>
        </q-btn>
      </q-card-section>

      <q-card-section class="q-pa-none">
        <q-table flat :rows="store.stats" :columns="columnsStatistik" row-key="jurusan" :loading="store.loading"
          :pagination="{ rowsPerPage: 10 }" hide-bottom>
          <template v-slot:body-cell-jurusan="props">
            <q-td :props="props" class="text-weight-bold text-primary">
              {{ props.row.jurusan }}
            </q-td>
          </template>
          <template v-slot:body-cell-akhir="props">
            <q-td :props="props" class="text-weight-bold">
              {{ props.row.akhir }}
            </q-td>
          </template>
        </q-table>
      </q-card-section>

      <q-card-section class="bg-blue-1 text-primary text-caption">
        <q-icon name="info" /> Rumus Baku Keadaan: <strong>Siswa Akhir = (Siswa Awal + Mutasi Masuk) - Mutasi Keluar</strong>.
        Angka ini ditarik secara *real-time* dari database Dapodik lokal.
      </q-card-section>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useQuasar } from 'quasar'
import { useKesiswaanReportStore } from '@/stores/kesiswaan/kesiswaanReportStore.js'
import { useContextStore } from '@/stores/contextStore.js'

const $q = useQuasar()
const store = useKesiswaanReportStore()
const contextStore = useContextStore()
const downloading = ref(false)

const columnsStatistik = [
  { name: 'jurusan', align: 'left', label: 'Kompetensi Keahlian (SMK)', field: 'jurusan' },
  { name: 'awal', align: 'center', label: 'Siswa Awal Bulan', field: 'awal', sortable: true },
  { name: 'masuk', align: 'center', label: 'Mutasi Masuk (+)', field: 'masuk' },
  { name: 'keluar', align: 'center', label: 'Mutasi Keluar (-)', field: 'keluar' },
  { name: 'akhir', align: 'center', label: 'Siswa Akhir Bulan', field: 'akhir', sortable: true }
]

onMounted(loadData)

watch(() => contextStore.selectedSemesterId, () => {
  loadData()
})

async function loadData() {
  await store.fetchStatistik()
}

async function downloadPDF() {
  downloading.value = true
  try {
    const result = await store.downloadPDF()
    if (!result.success) {
      $q.notify({ color: 'negative', message: result.error || 'Gagal mengunduh PDF', icon: 'error' })
      return
    }
    const url = URL.createObjectURL(result.data)
    const link = document.createElement('a')
    link.href = url
    link.download = `Rekap_Kesiswaan.pdf`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)
  } finally {
    downloading.value = false
  }
}

async function downloadExcel() {
  downloading.value = true
  try {
    const result = await store.downloadExcel()
    if (!result.success) {
      $q.notify({ color: 'negative', message: result.error || 'Gagal mengunduh Excel', icon: 'error' })
      return
    }
    const url = URL.createObjectURL(result.data)
    const link = document.createElement('a')
    link.href = url
    link.download = `Rekap_Kesiswaan.xlsx`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)
  } finally {
    downloading.value = false
  }
}
</script>

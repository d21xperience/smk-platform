<template>
  <q-page class="q-pa-md bg-grey-2">
    <!-- Header Halaman -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Direktori Siswa</h1>
        <p class="text-caption text-grey-7 q-mb-none">Manajemen data komputerisasi Buku Klaper dan data aktif siswa</p>
      </div>
      <div class="q-gutter-sm">
        <q-btn color="green-7" icon="table_view" label="Import XLS" @click="showImportDialog = true" />
        <q-btn color="deep-orange-7" icon="print" label="Cetak PDF" @click="exportToExcel" />
        <q-btn color="primary" icon="add" label="Tambah Siswa" @click="showAddStudents = true" />
      </div>
    </div>

    <!-- Panel Filter & Pencarian -->
    <q-card class="shadow-1 q-mb-md">
      <q-card-section class="q-py-sm">
        <div class="row q-col-gutter-sm items-center">
          <!-- Input Pencarian -->
          <div class="col-12 col-md-4">
            <q-input v-model="localFilters.search" dense outlined placeholder="Cari Nama, NIS, NISN..." clearable
              @update:model-value="applySearch">
              <template v-slot:prepend><q-icon name="search" /></template>
            </q-input>
          </div>

          <!-- Filter Tingkat (Mengirim Nilai Integer) -->
          <div class="col-6 col-md-2">
            <q-select v-model="localFilters.tingkat" 
              :options="[{label: 'Kelas X', value: 10}, {label: 'Kelas XI', value: 11}, {label: 'Kelas XII', value: 12}]" 
              label="Tingkat" dense outlined emit-value map-options clearable @update:model-value="applyFilters" />
          </div>
          
          <!-- Filter Jurusan -->
          <div class="col-6 col-md-2">
            <q-select v-model="localFilters.jurusan" 
              :options="['RPL', 'TKJ', 'AKL', 'OTKP', 'TKRO']" 
              label="Jurusan" dense outlined clearable @update:model-value="applyFilters" />
          </div>

          <!-- Filter Kelas (Dinamis dari API) -->
          <div class="col-6 col-md-2">
            <q-select v-model="localFilters.kelas" 
              :options="classOptions" 
              label="Kelas" dense outlined emit-value map-options clearable 
              :loading="loadingClasses" @update:model-value="applyFilters" />
          </div>

          <!-- Tombol Reset Filter -->
          <div class="col-6 col-md-2 flex justify-end">
            <q-btn flat color="grey-7" label="Reset" icon="refresh" @click="resetFilters" class="full-width" />
          </div>
        </div>
      </q-card-section>
    </q-card>

    <!-- Tabel Utama Data Siswa -->
    <q-card class="shadow-1">
      <q-card-section class="q-pa-none">
        <q-table flat :rows="students" :columns="columns" row-key="id" :loading="isLoading"
          :pagination="initialPagination" @request="onRequest">
          
          <!-- Status Badge -->
          <template v-slot:body-cell-status="props">
            <q-td :props="props">
              <q-badge :color="getStatusColor(props.row.status)" :label="getStatusLabel(props.row.status)" />
            </q-td>
          </template>

          <!-- Action Buttons -->
          <template v-slot:body-cell-aksi="props">
            <q-td :props="props" class="q-gutter-xs text-center">
              <q-btn size="sm" color="blue" icon="visibility" round flat @click="viewDetail(props.row.id)">
                <q-tooltip>Buku Induk Lengkap</q-tooltip>
              </q-btn>
              <q-btn size="sm" color="amber-9" icon="edit" round flat @click="editSiswa(props.row)">
                <q-tooltip>Ubah Data</q-tooltip>
              </q-btn>
              <q-btn size="sm" color="red" icon="person_remove" round flat @click="openDeleteDialog(props.row)">
                <q-tooltip>Proses Keluar / Mutasi</q-tooltip>
              </q-btn>
            </q-td>
          </template>
        </q-table>
      </q-card-section>
    </q-card>

    <!-- Komponen Import Dinamis -->
    <BaseImportDialog v-model="showImportDialog" title="Data Siswa" :loading="isImporting"
      @download-template="downloadTemplate" @submit="processImport" />

    <!-- Modal Pilihan Tambah Siswa -->
    <q-dialog v-model="showAddStudents">
      <q-card style="min-width: 400px; max-width: 500px; border-radius: 12px;">
        <q-card-section class="row items-center q-pb-none">
          <div class="text-h6 text-weight-bold text-primary">Metode Tambah Siswa</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>
        <q-card-section class="q-pa-md">
          <q-list bordered separator class="rounded-borders text-grey-9">
            <q-item clickable v-close-popup @click="goToPage('/psb/siswa/baru')">
              <q-item-section avatar><q-icon name="fiber_new" color="green" size="md" /></q-item-section>
              <q-item-section>
                <q-item-label class="text-weight-bold">Siswa Baru</q-item-label>
                <q-item-label caption>Input calon siswa baru melalui menu PSB</q-item-label>
              </q-item-section>
            </q-item>
            <q-item clickable v-close-popup @click="goToPage('/kesiswaan/mutasi-siswa')">
              <q-item-section avatar><q-icon name="sync_alt" color="orange" size="md" /></q-item-section>
              <q-item-section>
                <q-item-label class="text-weight-bold">Siswa Mutasi</q-item-label>
                <q-item-label caption>Input siswa pindahan (Masuk/Keluar)</q-item-label>
              </q-item-section>
            </q-item>
          </q-list>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Dialog Hapus & Notifikasi -->
    <q-dialog v-model="showDeleteDialog" persistent>
      <q-card style="width: 400px">
        <q-card-section>
          <div class="text-h6">Konfirmasi Hapus</div>
          <div class="q-mt-md">Apakah Anda yakin ingin menghapus siswa <strong>{{ selectedStudent?.fullName }}</strong>?</div>
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" color="grey" @click="showDeleteDialog = false" />
          <q-btn flat label="Hapus" color="negative" @click="confirmDelete" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <q-dialog v-model="showNotification">
      <q-card>
        <q-card-section :class="notificationType === 'success' ? 'bg-positive text-white' : 'bg-negative text-white'">
          <div class="text-h6">{{ notificationType === 'success' ? 'Berhasil' : 'Gagal' }}</div>
          <div>{{ notificationMessage }}</div>
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="OK" color="primary" @click="showNotification = false" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useQuasar } from 'quasar'
import * as XLSX from 'xlsx'
import { useStudentList } from '@/composables/student/useStudentList'
import { studentApi } from '@/adapters/api/studentApi.js'
import { useContextStore } from '@/stores/contextStore.js'
import BaseImportDialog from '@/components/ui/BaseImportDialog.vue'

const router = useRouter()
const $q = useQuasar()
const contextStore = useContextStore()

const showAddStudents = ref(false)
const showImportDialog = ref(false)
const isImporting = ref(false)

const {
  students, isLoading, localFilters,
  classOptions, loadingClasses,
  showDeleteDialog, showNotification, notificationMessage, notificationType,
  selectedStudent, getStatusColor, getStatusLabel,
  onRequest, viewDetail, openDeleteDialog, confirmDelete,
  loadData, applyFilters, applySearch, resetFilters
} = useStudentList()

const goToPage = (path) => router.push(path)

// Kolom Tabel
const columns = [
  { name: 'nis', align: 'left', label: 'NIS', field: 'nis', sortable: true },
  { name: 'nisn', align: 'left', label: 'NISN', field: 'nisn', sortable: true },
  { name: 'nama', align: 'left', label: 'Nama Lengkap', field: 'fullName', sortable: true },
  { name: 'jk', label: 'JK', field: 'gender', align: 'center' },
  { name: 'kelas', align: 'left', label: 'Kelas', field: 'className', sortable: true },
  { name: 'status', align: 'center', label: 'Status', field: 'status' },
  { name: 'aksi', label: 'Aksi', field: 'actions', align: 'center' },
]

const initialPagination = { sortBy: 'nama', descending: false, page: 1, rowsPerPage: 10 }

onMounted(() => {
  loadData()
})

const editSiswa = (siswa) => {
  router.push(`/kesiswaan/siswa/edit/${siswa.id}`)
}

const exportToExcel = () => {
  $q.notify({ color: 'green-8', message: 'Mengekspor data ke format .xlsx...', icon: 'file_download' })
}

// LOGIKA IMPORT EXCEL
const downloadTemplate = () => {
  const templateData = [
    { NISN: '0088779136', NIS: '20251788779136', Nama: 'Ahmad Fauzi', 'L/P': 'L', Kelas: 'X RPL 1' }
  ]

  const worksheet = XLSX.utils.json_to_sheet(templateData)
  const workbook = XLSX.utils.book_new()
  XLSX.utils.book_append_sheet(workbook, worksheet, 'Template Siswa')
  worksheet['!cols'] = [{ wch: 15 }, { wch: 20 }, { wch: 30 }, { wch: 5 }, { wch: 15 }]

  XLSX.writeFile(workbook, 'Template_Import_Siswa.xlsx')
  $q.notify({ type: 'info', message: 'Template berhasil diunduh' })
}

const processImport = (file) => {
  isImporting.value = true
  const reader = new FileReader()

  reader.onload = async (e) => {
    try {
      const data = new Uint8Array(e.target.result)
      const workbook = XLSX.read(data, { type: 'array' })
      const firstSheetName = workbook.SheetNames[0]
      const worksheet = workbook.Sheets[firstSheetName]
      const jsonData = XLSX.utils.sheet_to_json(worksheet)

      if (jsonData.length === 0) {
        throw new Error('File Excel kosong atau format tidak sesuai.')
      }

      let successCount = 0
      let failedCount = 0
      const context = contextStore.currentContext

      for (const row of jsonData) {
        const payload = {
          nisn: row.NISN || '',
          nis: row.NIS || '',
          nama: row.Nama || '',
          jenisKelamin: row['L/P'] || 'L',
          kelas: row.Kelas || 'X'
        }

        const res = await studentApi.registerStudent(payload, context)

        if (res.success) {
          successCount++
        } else {
          failedCount++
          console.warn(`Gagal import NIS ${row.NIS}:`, res.error?.message)
        }
      }

      isImporting.value = false
      showImportDialog.value = false

      $q.notify({
        type: failedCount === 0 ? 'positive' : 'warning',
        message: `Import selesai: ${successCount} berhasil, ${failedCount} gagal.`,
        caption: 'Tabel akan diperbarui otomatis.'
      })

      loadData()

    } catch (err) {
      isImporting.value = false
      $q.notify({ type: 'negative', message: 'Gagal memproses file: ' + err.message })
    }
  }

  reader.readAsArrayBuffer(file)
}
</script>

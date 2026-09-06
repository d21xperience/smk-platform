<template>
  <q-btn color="green-7" icon="table_view" label="Export Excel" @click="openExportExcel" />

  <!-- ==================== DIALOG EXPORT EXCEL ==================== -->
  <q-dialog v-model="exportExcelDialog" persistent>
    <q-card style="min-width: 500px">
      <q-card-section class="bg-green-7 text-white">
        <div class="text-h6"><q-icon name="table_view" class="q-mr-sm" />Export Data ke Excel</div>
      </q-card-section>
      <q-card-section>
        <div class="text-subtitle2 q-mb-sm">Pilih Kolom yang Ingin Diekspor:</div>
        <div class="row q-col-gutter-sm">
          <div class="col-6" v-for="col in excelColumnOptions" :key="col.field">
            <q-checkbox v-model="selectedExcelColumns" :val="col.field" :label="col.label" dense />
          </div>
        </div>
        <q-separator class="q-my-md" />
        <div class="row q-col-gutter-sm">
          <div class="col-12">
            <q-input v-model="excelFileName" label="Nama File Excel" outlined dense hint="Tanpa ekstensi .xlsx" />
          </div>
        </div>
      </q-card-section>
      <q-card-actions align="right">
        <q-btn flat label="Batal" color="negative" v-close-popup />
        <q-btn label="Export" icon="download" color="green-7" :disable="selectedExcelColumns.length === 0"
          @click="doExportExcel" />
      </q-card-actions>
    </q-card>
  </q-dialog>


</template>


<script setup>
import { ref } from 'vue'
import * as XLSX from 'xlsx'
import { useQuasar } from 'quasar'
const $q = useQuasar()

const props = defineProps({
  data: { type: Array, required: true }
})
// ==================== STATE EXPORT EXCEL ====================
const exportExcelDialog = ref(false)
const excelFileName = ref('Data_Siswa')

// Semua kolom yang tersedia untuk export
const excelColumnOptions = [
  { field: 'nis', label: 'NIS' },
  { field: 'nama', label: 'Nama Lengkap' },
  { field: 'kelas_nama', label: 'Kelas' },
  { field: 'tempat_lahir', label: 'Tempat Lahir' },
  { field: 'tanggal_lahir', label: 'Tanggal Lahir' },
  { field: 'jenis_kelamin', label: 'Jenis Kelamin' },
  { field: 'alamat', label: 'Alamat' },
  { field: 'no_hp', label: 'No. HP Siswa' },
  { field: 'status', label: 'Status' },
  { field: 'orang_tua.ayah', label: 'Nama Ayah' },
  { field: 'orang_tua.ibu', label: 'Nama Ibu' },
  { field: 'orang_tua.pekerjaan_ayah', label: 'Pekerjaan Ayah' },
  { field: 'orang_tua.pekerjaan_ibu', label: 'Pekerjaan Ibu' },
  { field: 'orang_tua.no_hp_ortu', label: 'No. HP Orang Tua' },
]

// Kolom default yang sudah tercentang
const selectedExcelColumns = ref(['nis', 'nama', 'kelas_nama', 'jenis_kelamin', 'status'])
/**
 * Ambil nilai nested dari objek menggunakan dot-notation
 * Contoh: getNestedValue(row, 'orang_tua.ayah') → row.orang_tua.ayah
 */
function getNestedValue(obj, path) {
  return path.split('.').reduce((acc, key) => acc?.[key] ?? '', obj)
}

// ==================== FUNGSI EXPORT EXCEL ====================
function openExportExcel() {
  exportExcelDialog.value = true
}
function doExportExcel() {
  const activeColumns = excelColumnOptions.filter(c => selectedExcelColumns.value.includes(c.field))

  // Buat array of objects sesuai kolom terpilih
  const exportData = props.data.list.map((row, index) => {
    const entry = { 'No.': index + 1 }
    activeColumns.forEach(col => {
      entry[col.label] = getNestedValue(row, col.field)
    })
    return entry
  })

  const ws = XLSX.utils.json_to_sheet(exportData)
  const wb = XLSX.utils.book_new()
  XLSX.utils.book_append_sheet(wb, ws, 'Data Siswa')

  // Auto - width kolom berdasarkan konten
  const colWidths = Object.keys(exportData[0] || {}).map(key => ({
    wch: Math.max(key.length, ...exportData.map(r => String(r[key] ?? '').length)) + 2
  }))
  ws['!cols'] = colWidths

  XLSX.writeFile(wb, `${excelFileName.value || 'Data_Siswa'}.xlsx`)
  exportExcelDialog.value = false
  $q.notify({ type: 'positive', message: 'File Excel berhasil diunduh', icon: 'check_circle' })
}

</script>

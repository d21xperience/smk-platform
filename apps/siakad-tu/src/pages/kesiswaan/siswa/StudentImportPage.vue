<template>
  <BaseImportData title="Siswa" :columns="siswaColumns" :previewRows="siswaPreviewData"
    pastePlaceholderHint="NISN | NIS | Nama Lengkap | L/P | Kelas"
    examplePasteText="0088779136 | 20251788 | Ahmad Fauzi | L | X RPL 1&#10;0088779019 | 20251789 | Siti Nurhaliza | P | X TKJ 2"
    @download-template="handleDownloadTemplate" @file-loaded="handleFileParsing" @submit-upload="saveSiswaData"
    @submit-paste="handlePaste" @clear-file="siswaPreviewData = []" />
</template>

<script setup>
import { ref } from 'vue'
import BaseImportData from '@/components/ui/BaseImportData.vue'
import { studentApi } from '@/adapters/api/studentApi.js'
import { useContextStore } from '@/stores/contextStore.js'
import { useQuasar } from 'quasar'
import * as XLSX from 'xlsx' // Pastikan: npm install xlsx

const $q = useQuasar()
const contextStore = useContextStore()
const siswaPreviewData = ref([])

// Kolom spesifik data Siswa
const siswaColumns = [
  { name: 'nisn', label: 'NISN', field: 'nisn', align: 'left' },
  { name: 'nis', label: 'NIS', field: 'nis', align: 'left' },
  { name: 'nama', label: 'Nama Lengkap', field: 'nama', align: 'left' },
  { name: 'jk', label: 'L/P', field: 'jk', align: 'center' },
  { name: 'kelas', label: 'Kelas', field: 'kelas', align: 'left' }
]

const handleDownloadTemplate = () => {
  // Logika download template statis
  window.open('/templates/template_siswa.xlsx')
}

const handleFileParsing = (file) => {
  const reader = new FileReader()
  reader.onload = (e) => {
    const data = new Uint8Array(e.target.result)
    const workbook = XLSX.read(data, { type: 'array' })
    const firstSheetName = workbook.SheetNames[0]
    const worksheet = workbook.Sheets[firstSheetName]

    // Parse ke JSON
    const jsonData = XLSX.utils.sheet_to_json(worksheet, { header: 1 })

    // Asumsi Baris 1 Header, ambil dari baris ke-2 (index 1)
    const rows = jsonData.slice(1).map((row, index) => ({
      id: index,
      nisn: row[0] || '',
      nis: row[1] || '',
      nama: row[2] || '',
      jk: row[3] || '',
      kelas: row[4] || ''
    }))

    siswaPreviewData.value = rows
  }
  reader.readAsArrayBuffer(file)
}

const handlePaste = (textData) => {
  const rows = textData.trim().split('\n').map((line, index) => {
    const cols = line.split('|').map(c => c.trim())
    return {
      id: index,
      nisn: cols[0] || '',
      nis: cols[1] || '',
      nama: cols[2] || '',
      jk: cols[3] || '',
      kelas: cols[4] || ''
    }
  })
  siswaPreviewData.value = rows
}

const saveSiswaData = async () => {
  if (siswaPreviewData.value.length === 0) return

  $q.loading.show({ message: 'Menyimpan data siswa...' })
  const context = contextStore.currentContext
  let successCount = 0

  try {
    // Loop simpan ke Backend
    for (const row of siswaPreviewData.value) {
      const payload = {
        nis: row.nis,
        nisn: row.nisn,
        nama: row.nama,
        jenisKelamin: row.jk,
        kelas: row.kelas,
        status: 'ACTIVE'
      }
      const res = await studentApi.registerStudent(payload, context)
      if (res.success) successCount++
    }

    $q.notify({ type: 'positive', message: `${successCount} data siswa berhasil diimport!` })
    siswaPreviewData.value = []
  } catch (err) {
    $q.notify({ type: 'negative', message: `Terjadi kesalahan saat import data. ${err}` })
  } finally {
    $q.loading.hide()
  }
}
</script>

<template>
  <BaseImportData title="Guru" :columns="guruColumns" :previewRows="guruPreviewData"
    pastePlaceholderHint="NUPTK | Nama Lengkap | Jenis Kelamin | Mata Pelajaran"
    examplePasteText="987654321	Ahmad Sobari	L	Matematika&#10;123456789	Siti Rahma	P	Bahasa Inggris"
    @download-template="handleDownloadTemplate" @file-loaded="handleFileParsing" @submit-upload="saveGuruData"
    @submit-paste="saveGuruDataPaste" @clear-file="guruPreviewData = []" />
</template>

<script setup>
import { ref } from 'vue'
import BaseImportData from '@/components/ui/BaseImportData.vue'

const guruPreviewData = ref([])

// Kolom spesifik data Guru
const guruColumns = [
  { name: 'nuptk', label: 'NUPTK', field: 'nuptk', align: 'left' },
  { name: 'nama', label: 'Nama Guru', field: 'nama', align: 'left' },
  { name: 'jk', label: 'L/P', field: 'jk', align: 'center' },
  { name: 'mapel', label: 'Mata Pelajaran', field: 'mapel', align: 'left' }
]

const handleDownloadTemplate = () => {
  console.log('Download template Excel Guru...')
  // window.open('/templates/template_guru.xlsx')
}

const handleFileParsing = (file) => {
  console.log('Membaca file guru:', file.name)
  // Dummy data simulasi pratinjau guru
  guruPreviewData.value = [
    { id: 1, nuptk: '2044761', nama: 'Hendra Wijaya, M.Pd', jk: 'L', mapel: 'Fisika' },
    { id: 2, nuptk: '2044762', nama: 'Rina Melati, S.Pd', jk: 'P', mapel: 'Kimia' }
  ]
}

const saveGuruData = (file) => {
  alert('Data Guru dari file berhasil disimpan ke database!', file)
}

const saveGuruDataPaste = (textData) => {
  alert('Data Guru dari teks salinan berhasil diproses!', textData)
}
</script>

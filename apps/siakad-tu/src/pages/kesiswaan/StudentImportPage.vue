<template>
  <BaseImportData title="Siswa" :columns="siswaColumns" :previewRows="siswaPreviewData"
    pastePlaceholderHint="NUPTK | Nama Lengkap | Jenis Kelamin | Mata Pelajaran"
    examplePasteText="987654321	Ahmad Sobari	L	Matematika&#10;123456789	Siti Rahma	P	Bahasa Inggris"
    @download-template="handleDownloadTemplate" @file-loaded="handleFileParsing" @submit-upload="saveGuruData"
    @submit-paste="saveGuruDataPaste" @clear-file="siswaPreviewData = []" />
</template>

<script setup>
import { ref } from 'vue'
import BaseImportData from '@/components/ui/BaseImportData.vue'

const siswaPreviewData = ref([])

// Kolom spesifik data Siswa
const siswaColumns = [
  { name: 'nuptk', label: 'NUPTK', field: 'nuptk', align: 'left' },
  { name: 'nama', label: 'Nama Siswa', field: 'nama', align: 'left' },
  { name: 'jk', label: 'L/P', field: 'jk', align: 'center' },
  { name: 'mapel', label: 'Mata Pelajaran', field: 'mapel', align: 'left' }
]

const handleDownloadTemplate = () => {
  console.log('Download template Excel Siswa...')
  // window.open('/templates/template_siswa.xlsx')
}

const handleFileParsing = (file) => {
  console.log('Membaca file Siswa:', file.name)
  // Dummy data simulasi pratinjau Siswa
  siswaPreviewData.value = [
    { id: 1, nuptk: '2044761', nama: 'Hendra Wijaya, M.Pd', jk: 'L', mapel: 'Fisika' },
    { id: 2, nuptk: '2044762', nama: 'Rina Melati, S.Pd', jk: 'P', mapel: 'Kimia' }
  ]
}

const saveGuruData = (file) => {
  alert('Data Siswa dari file berhasil disimpan ke database!', file)
}

const saveGuruDataPaste = (textData) => {
  alert('Data Siswa dari teks salinan berhasil diproses!', textData)
}
</script>

<template>
  <q-page class="q-pa-md bg-grey-1">
    <!-- Judul Halaman Dinamis via Props -->
    <div class="row items-center q-mb-md">
      <q-btn flat round icon="arrow_back" color="primary" @click="$router.back()" />
      <h5 class="text-h5 text-weight-bold q-my-none q-ml-sm text-primary">Import Data {{ title }}</h5>
    </div>

    <q-card flat bordered class="rounded-borders shadow-1">
      <q-tabs v-model="activeTab" dense class="text-grey" active-color="primary" indicator-color="primary"
        align="justify" narrow-indicator>
        <q-tab name="upload" icon="upload_file" :label="'Upload File ' + title" />
        <q-tab name="copypaste" icon="content_paste" label="Copy Paste Data" />
      </q-tabs>

      <q-separator />

      <q-tab-panels v-model="activeTab" animated class="q-pa-lg">

        <!-- ==================== TAB A: UPLOAD ==================== -->
        <q-tab-panel name="upload" class="q-pa-none">

          <!-- Download Template Dinamis -->
          <div class="bg-blue-1 border-blue q-pa-md rounded-borders row items-center justify-between q-mb-lg">
            <div class="row items-center">
              <q-icon name="info" color="primary" size="sm" class="q-mr-sm" />
              <div>
                <div class="text-weight-bold text-primary">Gunakan Template Standar {{ title }}</div>
                <div class="text-caption text-grey-8">Pastikan susunan kolom sesuai format template agar tidak terjadi
                  eror.</div>
              </div>
            </div>
            <q-btn color="primary" icon="download" label="Unduh Format Excel" outline
              @click="$emit('download-template')" />
          </div>

          <!-- Drag and Drop Area -->
          <div class="dropzone-area flex flex-center column q-pa-xl text-center q-mb-lg cursor-pointer"
            :class="{ 'dropzone-active': isDragging }" @dragover.prevent="isDragging = true"
            @dragleave.prevent="isDragging = false" @drop.prevent="handleFileDrop" @click="triggerFileInput">
            <input type="file" ref="fileInputRef" class="hidden" accept=".xlsx, .xls" @change="handleFileSelect" />
            <q-icon :name="selectedFile ? 'task_alt' : 'cloud_upload'" :color="selectedFile ? 'positive' : 'primary'"
              size="56px" class="q-mb-md" />

            <div v-if="!selectedFile">
              <div class="text-subtitle1 text-weight-bold text-grey-9">
                Tarik & Lepaskan file Excel di sini atau <span class="text-primary text-underline">Pilih File</span>
              </div>
              <div class="text-caption text-grey-6 q-mt-xs">Format dokumen hanya .xlsx dan .xls</div>
            </div>
            <div v-else>
              <div class="text-subtitle1 text-weight-bold text-positive">{{ selectedFile.name }}</div>
              <q-btn flat color="negative" label="Hapus File" icon="delete" size="sm" class="q-mt-sm"
                @click.stop="clearFile" />
            </div>
          </div>

          <!-- Preview Tabel Dinamis via Props Columns dan Rows -->
          <div v-if="previewRows.length > 0" class="q-mb-lg animate-fade">
            <div class="text-subtitle1 text-weight-bold text-grey-9 q-mb-sm">Pratinjau Data ({{ previewRows.length }}
              Baris)</div>
            <q-table flat bordered dense :rows="previewRows" :columns="columns" row-key="id"
              :pagination="{ rowsPerPage: 5 }" class="bg-white shadow-1" />
          </div>

          <div class="row justify-end q-gutter-sm">
            <q-btn label="Batal" color="grey" flat @click="clearFile" :disabled="!selectedFile" />
            <q-btn label="Proses Impor" color="primary" icon="bolt" :disabled="!selectedFile"
              @click="$emit('submit-upload', selectedFile)" />
          </div>
        </q-tab-panel>

        <!-- ==================== TAB B: COPY PASTE ==================== -->
        <q-tab-panel name="copypaste" class="q-pa-none">
          <div class="text-subtitle1 text-weight-bold text-grey-9 q-mb-xs">Tempel Data Massal</div>
          <div class="text-caption text-grey-6 q-mb-md">
            Format kolom yang disarankan: <span class="text-weight-bold text-primary">{{ pastePlaceholderHint }}</span>
          </div>

          <q-input v-model="pasteText" type="textarea" filled rows="10" :placeholder="examplePasteText"
            class="text-body2 monospace-font q-mb-lg" />

          <div class="row justify-end q-gutter-sm">
            <q-btn label="Bersihkan" color="grey" flat @click="pasteText = ''" :disabled="!pasteText" />
            <q-btn label="Proses Salinan" color="primary" icon="content_paste_go" :disabled="!pasteText"
              @click="$emit('submit-paste', pasteText)" />
          </div>
        </q-tab-panel>

      </q-tab-panels>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, defineProps, defineEmits } from 'vue'

// 1. Mendefinisikan PROPS agar komponen flexibel
defineProps({
  title: String,               // Contoh: 'Guru' atau 'Pegawai'
  columns: Array,             // Struktur kolom pratinjau tabel
  previewRows: Array,         // Data baris hasil baca excel untuk pratinjau
  pastePlaceholderHint: String, // Petunjuk teks format kolom copy-paste
  examplePasteText: String     // Contoh isi data teks di dalam textarea
})

// 2. Mendefinisikan EMIT untuk mengirim data keluar
const emit = defineEmits(['download-template', 'file-loaded', 'submit-upload', 'submit-paste', 'clear-file'])

const activeTab = ref('upload')
const isDragging = ref(false)
const selectedFile = ref(null)
const fileInputRef = ref(null)
const pasteText = ref('')

const triggerFileInput = () => fileInputRef.value.click()
const handleFileSelect = (e) => { if (e.target.files.length > 0) onLoadFile(e.target.files[0]) }

const handleFileDrop = (e) => {
  isDragging.value = false
  if (e.dataTransfer.files.length > 0) {
    const file = e.dataTransfer.files[0]
    if (file.name.endsWith('.xlsx') || file.name.endsWith('.xls')) onLoadFile(file)
  }
}

const onLoadFile = (file) => {
  selectedFile.value = file
  emit('file-loaded', file) // Kirim file ke Parent untuk di-parse menjadi tabel pratinjau
}

const clearFile = () => {
  selectedFile.value = null
  if (fileInputRef.value) fileInputRef.value.value = ''
  emit('clear-file')
}
</script>

<style scoped>
.dropzone-area {
  border: 2px dashed #9e9e9e;
  border-radius: 8px;
  background-color: #fafafa;
  transition: all 0.3s ease;
}

.dropzone-area:hover {
  border-color: #1976d2;
  background-color: #e3f2fd;
}

.dropzone-active {
  border-color: #2e7d32;
  background-color: #e8f5e9;
}

.border-blue {
  border: 1px solid #bbdefb;
}

.monospace-font {
  font-family: 'Courier New', Courier, monospace;
}
</style>

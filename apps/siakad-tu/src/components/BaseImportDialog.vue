<template>
  <q-dialog v-model="dialogVisible" persistent>
    <q-card style="min-width: 400px; max-width: 500px; border-radius: 12px;">
      <!-- Header -->
      <q-card-section class="row items-center q-pb-none">
        <div class="text-h6 text-weight-bold text-primary">{{ title }}</div>
        <q-space />
        <q-btn icon="close" flat round dense v-close-popup />
      </q-card-section>

      <!-- Progress Bar (Sesuai Preferensi: Bukan Spinner) -->
      <q-linear-progress v-if="loading" indeterminate color="primary" class="q-mt-sm" />

      <q-card-section class="q-pa-md">
        <div class="text-caption text-grey-7 q-mb-sm">
          Format yang didukung: .xlsx, .xls. Pastikan data sesuai dengan template.
        </div>

        <!-- File Picker -->
        <q-file v-model="selectedFile" label="Pilih File Excel" outlined accept=".xlsx, .xls" :disable="loading"
          @update:model-value="handleFileSelected">
          <template v-slot:prepend>
            <q-icon name="attach_file" />
          </template>
        </q-file>

        <!-- Preview Data (Opsional, untuk UX yang lebih baik) -->
        <div v-if="parsedPreview.length > 0" class="q-mt-md">
          <div class="text-caption text-weight-bold text-grey-8">
            Pratinjau: {{ parsedPreview.length }} baris data ditemukan.
          </div>
        </div>
      </q-card-section>

      <!-- Actions -->
      <q-card-actions align="right" class="q-pa-md">
        <q-btn flat color="primary" label="Unduh Template" icon="download" :disable="loading"
          @click="$emit('download-template')" />
        <q-space />
        <q-btn flat label="Batal" color="grey" :disable="loading" v-close-popup />
        <q-btn color="primary" label="Proses Import" :disable="!selectedFile || loading" :loading="loading"
          @click="submitData" />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { ref, computed } from 'vue'
import * as XLSX from 'xlsx'

const props = defineProps({
  modelValue: { type: Boolean, required: true },
  title: { type: String, default: 'Import Data' },
  loading: { type: Boolean, default: false }
})

const emit = defineEmits(['update:modelValue', 'download-template', 'submit'])

const selectedFile = ref(null)
const parsedPreview = ref([])

const dialogVisible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const handleFileSelected = (file) => {
  if (!file) {
    parsedPreview.value = []
    return
  }

  const reader = new FileReader()
  reader.onload = (e) => {
    try {
      const data = new Uint8Array(e.target.result)
      const workbook = XLSX.read(data, { type: 'array' })
      const firstSheetName = workbook.SheetNames[0]
      const worksheet = workbook.Sheets[firstSheetName]

      // Parse ke JSON dan simpan sebagai preview
      parsedPreview.value = XLSX.utils.sheet_to_json(worksheet)
    } catch (err) {
      console.error('Gagal membaca file Excel:', err)
      parsedPreview.value = []
    }
  }
  reader.readAsArrayBuffer(file)
}

const submitData = () => {
  if (parsedPreview.value.length === 0) {
    // Fallback jika parsing gagal tapi file ada, atau validasi dasar
    alert('File kosong atau format tidak valid.')
    return
  }
  // Kirim data yang sudah diparsing ke parent, bukan file mentah
  emit('submit', parsedPreview.value)
}
</script>

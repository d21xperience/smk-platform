<template>
  <q-dialog v-model="formDialogOpen" persistent>
    <q-card style="min-width: 500px; max-width: 90vw">
      <q-card-section>
        <div class="text-h6">{{ isEditing ? 'Edit Dokumen' : 'Upload Dokumen Baru' }}</div>
      </q-card-section>

      <q-card-section class="q-gutter-md">
        <q-select v-model="form.documentType" :options="documentTypeOptions" label="Jenis Dokumen" outlined dense
          emit-value map-options :disable="isEditing" />

        <q-select v-model="form.documentScope" :options="documentScopeOptions" label="Scope Dokumen" outlined dense
          emit-value map-options :disable="isEditing" />

        <q-select v-if="form.documentScope === 'academic_year' || form.documentScope === 'semester'"
          v-model="form.academicYearId" :options="academicYearOptions" label="Tahun Pelajaran" outlined dense emit-value
          map-options :disable="isEditing" @update:model-value="handleAcademicYearChange" />

        <q-select v-if="form.documentScope === 'semester'" v-model="form.semesterId" :options="semesterOptions"
          label="Semester" outlined dense emit-value map-options :disable="isEditing"
          @update:model-value="handleSemesterChange" />

        <q-file v-model="selectedFile" label="Pilih File" outlined dense accept=".pdf,.jpg,.jpeg,.png,.doc,.docx"
          :max-file-size="maxFileSize" @update:model-value="handleFileSelect">
          <template #prepend>
            <q-icon name="attach_file" />
          </template>
        </q-file>

        <div v-if="form.fileName" class="text-caption text-grey">
          File: {{ form.fileName }} ({{ formatFileSize(form.fileSize) }})
        </div>
      </q-card-section>

      <q-card-actions align="right">
        <q-btn flat label="Batal" color="primary" @click="closeFormDialog" />
        <q-btn flat label="Simpan" color="positive" :loading="isSaving" :disable="!isFormValid" @click="handleSave" />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useQuasar } from 'quasar'
import { useTeacherDocument } from '../../composables/useTeacherDocument.js'
import { getAllDocumentTypes, getMaxFileSize } from '../../domain/teacher-document/models/DocumentType.js'
import { getAllDocumentScopes } from '../../domain/teacher-document/models/DocumentScope.js'

const $q = useQuasar()
const {
  formDialogOpen,
  editingDocument,
  isEditing,
  isSaving,
  saveDocument,
  closeFormDialog
} = useTeacherDocument()

const documentTypeOptions = getAllDocumentTypes()
const documentScopeOptions = getAllDocumentScopes()
const maxFileSize = getMaxFileSize()

const academicYearOptions = [
  { label: '2026/2027', value: 'AY-2026', name: '2026/2027' },
  { label: '2025/2026', value: 'AY-2025', name: '2025/2026' }
]

const semesterOptions = [
  { label: 'Semester 1 (20261)', value: '20261', name: 'Semester 1' },
  { label: 'Semester 2 (20262)', value: '20262', name: 'Semester 2' }
]

const form = ref({
  id: null,
  documentType: '',
  documentScope: '',
  academicYearId: null,
  academicYearName: null,
  semesterId: null,
  semesterName: null,
  fileName: '',
  fileSize: 0,
  mimeType: ''
})

const selectedFile = ref(null)

const isFormValid = computed(() => {
  if (!form.value.documentType) return false
  if (!form.value.documentScope) return false
  if (!form.value.fileName) return false

  if (form.value.documentScope === 'academic_year' && !form.value.academicYearId) {
    return false
  }

  if (form.value.documentScope === 'semester') {
    if (!form.value.academicYearId || !form.value.semesterId) {
      return false
    }
  }

  return true
})
const resetForm = () => {
  form.value = {
    id: null,
    documentType: '',
    documentScope: '',
    academicYearId: null,
    academicYearName: null,
    semesterId: null,
    semesterName: null,
    fileName: '',
    fileSize: 0,
    mimeType: ''
  }
  selectedFile.value = null
}
watch(editingDocument, (newVal) => {
  if (newVal) {
    form.value = {
      id: newVal.id,
      documentType: newVal.documentType,
      documentScope: newVal.documentScope,
      academicYearId: newVal.academicYearId,
      academicYearName: newVal.academicYearName,
      semesterId: newVal.semesterId,
      semesterName: newVal.semesterName,
      fileName: newVal.fileName,
      fileSize: newVal.fileSize,
      mimeType: newVal.mimeType
    }
    selectedFile.value = null
  } else {
    resetForm()
  }
}, { immediate: true })



const handleFileSelect = (file) => {
  if (!file) {
    form.value.fileName = ''
    form.value.fileSize = 0
    form.value.mimeType = ''
    return
  }

  form.value.fileName = file.name
  form.value.fileSize = file.size
  form.value.mimeType = file.type
}

const handleAcademicYearChange = (value) => {
  const selected = academicYearOptions.find(opt => opt.value === value)
  form.value.academicYearName = selected ? selected.name : null
  form.value.semesterId = null
  form.value.semesterName = null
}

const handleSemesterChange = (value) => {
  const selected = semesterOptions.find(opt => opt.value === value)
  form.value.semesterName = selected ? selected.name : null
}

const handleSave = async () => {
  try {
    await saveDocument({ ...form.value })
    closeFormDialog()
    resetForm()
    $q.notify({
      type: 'positive',
      message: isEditing.value ? 'Dokumen berhasil diperbarui' : 'Dokumen berhasil diupload'
    })
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: err.message
    })
  }
}

const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}
</script>

<template>
  <q-page padding>
    <div class="row items-center justify-between q-mb-md">
      <div class="text-h5">Teacher Documents</div>
      <q-btn color="primary" label="Upload Dokumen" icon="upload" @click="openFormDialog(null)" />
    </div>

    <q-banner v-if="error" class="bg-negative text-white q-mb-md">
      {{ error }}
    </q-banner>

    <q-card class="q-mb-md">
      <q-card-section>
        <div class="row q-gutter-md">
          <q-select v-model="filterScope" :options="scopeOptions" label="Filter Scope" outlined dense emit-value
            map-options clearable style="min-width: 200px" @update:model-value="handleFilterChange" />

          <q-select v-if="filterScope === 'academic_year' || filterScope === 'semester'" v-model="filterAcademicYear"
            :options="academicYearOptions" label="Tahun Pelajaran" outlined dense emit-value map-options clearable
            style="min-width: 200px" @update:model-value="handleFilterChange" />

          <q-select v-if="filterScope === 'semester'" v-model="filterSemester" :options="semesterOptions"
            label="Semester" outlined dense emit-value map-options clearable style="min-width: 200px"
            @update:model-value="handleFilterChange" />
        </div>
      </q-card-section>
    </q-card>

    <q-card v-if="isLoading" class="q-mb-md">
      <q-card-section class="text-center">
        <q-spinner color="primary" size="2em" />
        <div class="q-mt-sm">Loading...</div>
      </q-card-section>
    </q-card>

    <template v-else>
      <q-card v-if="documents.length === 0">
        <q-card-section class="text-center text-grey">
          Tidak ada dokumen. Klik "Upload Dokumen" untuk menambahkan.
        </q-card-section>
      </q-card>

      <q-table v-else :rows="documents" :columns="columns" row-key="id" :pagination="{ rowsPerPage: 10 }" flat bordered>
        <template #body-cell-documentType="props">
          <q-td :props="props">
            <q-badge color="info">
              {{ getDocumentTypeLabel(props.value) }}
            </q-badge>
          </q-td>
        </template>

        <template #body-cell-documentScope="props">
          <q-td :props="props">
            <q-badge :color="getScopeColor(props.value)">
              {{ getDocumentScopeLabel(props.value) }}
            </q-badge>
          </q-td>
        </template>

        <template #body-cell-period="props">
          <q-td :props="props">
            <span v-if="props.row.documentScope === 'global'" class="text-grey">-</span>
            <span v-else-if="props.row.documentScope === 'academic_year'">
              {{ props.row.academicYearName }}
            </span>
            <span v-else>
              {{ props.row.academicYearName }} / {{ props.row.semesterName }}
            </span>
          </q-td>
        </template>

        <template #body-cell-fileSize="props">
          <q-td :props="props">
            {{ formatFileSize(props.value) }}
          </q-td>
        </template>

        <template #body-cell-actions="props">
          <q-td :props="props">
            <q-btn flat dense round icon="download" color="primary" @click="handleDownload(props.row)">
              <q-tooltip>Download</q-tooltip>
            </q-btn>
            <q-btn flat dense round icon="edit" color="warning" @click="openFormDialog(props.row)">
              <q-tooltip>Edit</q-tooltip>
            </q-btn>
            <q-btn flat dense round icon="delete" color="negative" @click="handleDelete(props.row)">
              <q-tooltip>Delete</q-tooltip>
            </q-btn>
          </q-td>
        </template>
      </q-table>
    </template>

    <TeacherDocumentFormDialog />

    <q-dialog v-model="confirmDeleteDialog">
      <q-card style="min-width: 350px">
        <q-card-section>
          <div class="text-h6">Konfirmasi Hapus</div>
        </q-card-section>

        <q-card-section class="q-pt-none">
          Yakin ingin menghapus dokumen "{{ documentToDelete?.fileName }}"?
          Aksi ini tidak dapat dibatalkan.
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="Batal" color="primary" v-close-popup />
          <q-btn flat label="Hapus" color="negative" :loading="isDeleting" @click="confirmDelete" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useTeacherDocument } from '../../composables/useTeacherDocument.js'
import { getDocumentTypeLabel } from '../../domain/teacher-document/models/DocumentType.js'
import { getDocumentScopeLabel, getAllDocumentScopes } from '../../domain/teacher-document/models/DocumentScope.js'
import TeacherDocumentFormDialog from './TeacherDocumentFormDialog.vue'

const $q = useQuasar()
const {
  documents,
  isLoading,
  isDeleting,
  error,
  listDocuments,
  deleteDocument,
  downloadDocument,
  openFormDialog
} = useTeacherDocument()

const filterScope = ref(null)
const filterAcademicYear = ref(null)
const filterSemester = ref(null)
const confirmDeleteDialog = ref(false)
const documentToDelete = ref(null)

const scopeOptions = getAllDocumentScopes()

const academicYearOptions = [
  { label: '2026/2027', value: 'AY-2026' },
  { label: '2025/2026', value: 'AY-2025' }
]

const semesterOptions = [
  { label: 'Semester 1 (20261)', value: '20261' },
  { label: 'Semester 2 (20262)', value: '20262' }
]

const columns = [
  { name: 'fileName', label: 'Nama File', field: 'fileName', align: 'left', sortable: true },
  { name: 'documentType', label: 'Jenis', field: 'documentType', align: 'left', sortable: true },
  { name: 'documentScope', label: 'Scope', field: 'documentScope', align: 'left', sortable: true },
  { name: 'period', label: 'Periode', field: 'documentScope', align: 'left' },
  { name: 'fileSize', label: 'Ukuran', field: 'fileSize', align: 'right', sortable: true },
  { name: 'createdAt', label: 'Tanggal Upload', field: 'createdAt', align: 'left', sortable: true, format: val => new Date(val).toLocaleDateString('id-ID') },
  { name: 'actions', label: 'Aksi', field: 'id', align: 'center' }
]

onMounted(async () => {
  await listDocuments()
})

const handleFilterChange = async () => {
  await listDocuments({
    documentScope: filterScope.value,
    academicYearId: filterAcademicYear.value,
    semesterId: filterSemester.value
  })
}

const handleDownload = async (document) => {
  try {
    const downloadInfo = await downloadDocument(document.id)
    $q.notify({
      type: 'positive',
      message: `Download URL generated: ${downloadInfo.fileName}`,
      caption: `Expires at: ${new Date(downloadInfo.expiresAt).toLocaleString('id-ID')}`
    })
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: err.message
    })
  }
}

const handleDelete = (document) => {
  documentToDelete.value = document
  confirmDeleteDialog.value = true
}

const confirmDelete = async () => {
  try {
    await deleteDocument(documentToDelete.value.id)
    confirmDeleteDialog.value = false
    documentToDelete.value = null
    $q.notify({
      type: 'positive',
      message: 'Dokumen berhasil dihapus'
    })
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: err.message
    })
  }
}

const getScopeColor = (scope) => {
  const colorMap = {
    global: 'grey',
    academic_year: 'blue',
    semester: 'green'
  }
  return colorMap[scope] || 'grey'
}

const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}
</script>

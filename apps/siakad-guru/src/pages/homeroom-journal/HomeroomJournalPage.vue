<!-- FILE: src/pages/homeroom-journal/HomeroomJournalPage.vue -->
<!-- STATUS: NEW -->
<!-- STATUS IMPLEMENTASI: COMPLETE -->

<template>
  <q-page class="q-pa-md">
    <div class="row items-center justify-between q-mb-md">
      <div class="row items-center">
        <div class="text-h5">Catatan Siswa</div>
        <q-badge v-if="className" color="primary" class="q-ml-md" style="font-size: 14px; padding: 4px 12px;">
          Kelas {{ className }}
        </q-badge>
      </div>
      <q-btn color="primary" icon="add" label="Buat Catatan" @click="openCreateDialog" />
    </div>

    <q-banner v-if="error" class="bg-negative text-white q-mb-md">
      {{ error }}
    </q-banner>

    <div v-if="loading" class="flex flex-center" style="min-height: 200px;">
      <q-spinner-dots size="60px" color="primary" />
    </div>

    <template v-if="!loading">
      <div v-if="summary" class="row q-col-gutter-md q-mb-md">
        <div class="col-6 col-sm-3">
          <q-card flat bordered>
            <q-card-section class="text-center q-py-sm">
              <div class="text-h4 text-primary">{{ summary.totalNotes }}</div>
              <div class="text-caption text-grey-7">Total Catatan</div>
            </q-card-section>
          </q-card>
        </div>
        <div class="col-6 col-sm-3">
          <q-card flat bordered>
            <q-card-section class="text-center q-py-sm">
              <div class="text-h4 text-warning">{{ summary.draftCount }}</div>
              <div class="text-caption text-grey-7">Draft</div>
            </q-card-section>
          </q-card>
        </div>
        <div class="col-6 col-sm-3">
          <q-card flat bordered>
            <q-card-section class="text-center q-py-sm">
              <div class="text-h4 text-positive">{{ summary.submittedCount }}</div>
              <div class="text-caption text-grey-7">Submitted</div>
            </q-card-section>
          </q-card>
        </div>
        <div class="col-6 col-sm-3">
          <q-card flat bordered>
            <q-card-section class="text-center q-py-sm">
              <div class="text-h4 text-negative">{{ summary.urgentCount }}</div>
              <div class="text-caption text-grey-7">Urgent</div>
            </q-card-section>
          </q-card>
        </div>
      </div>

      <q-table
        flat
        bordered
        :rows="notes"
        :columns="columns"
        row-key="noteId"
        :pagination="{ rowsPerPage: 10, sortBy: 'noteDate', ascending: false }"
      >
        <template #body-cell-category="props">
          <q-td :props="props" class="text-center">
            <q-badge :color="getCategoryColor(props.row.category)">
              {{ getCategoryLabel(props.row.category) }}
            </q-badge>
          </q-td>
        </template>

        <template #body-cell-priority="props">
          <q-td :props="props" class="text-center">
            <q-badge :color="getPriorityColor(props.row.priority)" outline>
              {{ props.row.priority }}
            </q-badge>
          </q-td>
        </template>

        <template #body-cell-status="props">
          <q-td :props="props" class="text-center">
            <q-badge :color="props.row.status === 'SUBMITTED' ? 'positive' : 'warning'">
              {{ props.row.status === 'SUBMITTED' ? 'Submitted' : 'Draft' }}
            </q-badge>
          </q-td>
        </template>

        <template #body-cell-actions="props">
          <q-td :props="props" class="text-center">
            <q-btn
              v-if="props.row.status === 'DRAFT'"
              flat
              dense
              icon="send"
              color="primary"
              size="sm"
              title="Submit catatan"
              :loading="submitting"
              @click="confirmSubmit(props.row)"
            />
            <q-btn
              v-if="props.row.status === 'DRAFT'"
              flat
              dense
              icon="delete"
              color="negative"
              size="sm"
              title="Hapus draft"
              :loading="deleting"
              @click="confirmDelete(props.row)"
            />
            <q-btn
              flat
              dense
              icon="visibility"
              color="grey"
              size="sm"
              title="Lihat detail"
              @click="viewNote(props.row)"
            />
          </q-td>
        </template>
      </q-table>
    </template>

    <q-dialog v-model="showCreateDialog" persistent>
      <q-card style="min-width: 500px;">
        <q-card-section>
          <div class="text-h6">Buat Catatan Baru</div>
        </q-card-section>

        <q-card-section class="q-pt-none">
          <q-select
            v-model="formData.studentId"
            :options="studentOptions"
            label="Siswa"
            outlined
            dense
            emit-value
            map-options
            class="q-mb-md"
            @update:model-value="onStudentSelected"
          />

          <div class="row q-col-gutter-md q-mb-md">
            <div class="col-6">
              <q-select
                v-model="formData.category"
                :options="categoryOptions"
                label="Kategori"
                outlined
                dense
                emit-value
                map-options
              />
            </div>
            <div class="col-6">
              <q-select
                v-model="formData.priority"
                :options="priorityOptions"
                label="Prioritas"
                outlined
                dense
                emit-value
                map-options
              />
            </div>
          </div>

          <q-input
            v-model="formData.title"
            label="Judul Catatan"
            outlined
            dense
            class="q-mb-md"
            :maxlength="200"
            counter
          />

          <q-input
            v-model="formData.description"
            label="Deskripsi"
            type="textarea"
            outlined
            dense
            rows="4"
            class="q-mb-md"
            :maxlength="1000"
            counter
          />

          <q-input
            v-model="formData.noteDate"
            label="Tanggal Catatan"
            outlined
            dense
            mask="####-##-##"
            hint="Format: YYYY-MM-DD"
          >
            <template #append>
              <q-icon name="event" class="cursor-pointer">
                <q-popup-proxy cover transition-show="scale" transition-hide="scale">
                  <q-date v-model="formData.noteDate" mask="YYYY-MM-DD">
                    <div class="row items-center justify-end">
                      <q-btn v-close-popup label="Tutup" color="primary" flat />
                    </div>
                  </q-date>
                </q-popup-proxy>
              </q-icon>
            </template>
          </q-input>
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="Batal" color="grey" @click="showCreateDialog = false" />
          <q-btn
            label="Simpan Draft"
            color="primary"
            :loading="creating"
            @click="handleCreateNote"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <q-dialog v-model="showDetailDialog">
      <q-card v-if="selectedNote" style="min-width: 450px;">
        <q-card-section>
          <div class="text-h6">{{ selectedNote.title }}</div>
          <div class="text-caption text-grey">
            {{ selectedNote.studentName }} — {{ selectedNote.noteDate }}
          </div>
        </q-card-section>

        <q-card-section class="q-pt-none">
          <div class="row q-col-gutter-sm q-mb-md">
            <div class="col-4">
              <q-badge :color="getCategoryColor(selectedNote.category)">
                {{ getCategoryLabel(selectedNote.category) }}
              </q-badge>
            </div>
            <div class="col-4">
              <q-badge :color="getPriorityColor(selectedNote.priority)" outline>
                {{ selectedNote.priority }}
              </q-badge>
            </div>
            <div class="col-4">
              <q-badge :color="selectedNote.status === 'SUBMITTED' ? 'positive' : 'warning'">
                {{ selectedNote.status }}
              </q-badge>
            </div>
          </div>
          <p>{{ selectedNote.description }}</p>
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="Tutup" color="grey" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { onMounted, ref, computed } from 'vue'
import { useQuasar } from 'quasar'
import { useHomeroomJournal } from '@/composables/useHomeroomJournal.js'
import { useHomeroomStudent } from '@/composables/useHomeroomStudent.js'

const $q = useQuasar()

const {
  notes,
  summary,
  loading,
  creating,
  submitting,
  deleting,
  error,
  loadNotes,
  createNote,
  submitNote,
  deleteNote,
} = useHomeroomJournal()

const {
  students,
  className,
  loadStudents,
} = useHomeroomStudent()

const showCreateDialog = ref(false)
const showDetailDialog = ref(false)
const selectedNote = ref(null)

const formData = ref({
  studentId: null,
  studentName: '',
  category: 'NEUTRAL',
  priority: 'MEDIUM',
  title: '',
  description: '',
  noteDate: '',
})

const columns = ref([
  { name: 'noteDate', label: 'Tanggal', field: 'noteDate', align: 'center', sortable: true },
  { name: 'studentName', label: 'Siswa', field: 'studentName', align: 'left', sortable: true },
  { name: 'category', label: 'Kategori', field: 'category', align: 'center' },
  { name: 'priority', label: 'Prioritas', field: 'priority', align: 'center' },
  { name: 'title', label: 'Judul', field: 'title', align: 'left' },
  { name: 'status', label: 'Status', field: 'status', align: 'center' },
  { name: 'actions', label: 'Aksi', field: 'actions', align: 'center' },
])

const categoryOptions = ref([
  { label: 'Positif', value: 'POSITIVE' },
  { label: 'Negatif', value: 'NEGATIVE' },
  { label: 'Netral', value: 'NEUTRAL' },
  { label: 'Insiden', value: 'INCIDENT' },
])

const priorityOptions = ref([
  { label: 'Low', value: 'LOW' },
  { label: 'Medium', value: 'MEDIUM' },
  { label: 'High', value: 'HIGH' },
  { label: 'Urgent', value: 'URGENT' },
])

const studentOptions = computed(() => {
  return students.value.map(s => ({
    label: `${s.seatNumber}. ${s.fullName}`,
    value: s.id,
  }))
})

function onStudentSelected(studentId) {
  const found = students.value.find(s => s.id === studentId)
  if (found) {
    formData.value.studentName = found.fullName
  }
}

function openCreateDialog() {
  formData.value = {
    studentId: null,
    studentName: '',
    category: 'NEUTRAL',
    priority: 'MEDIUM',
    title: '',
    description: '',
    noteDate: new Date().toISOString().split('T')[0],
  }
  showCreateDialog.value = true
}

async function handleCreateNote() {
  try {
    await createNote({ ...formData.value })
    showCreateDialog.value = false
    $q.notify({
      type: 'positive',
      message: 'Catatan berhasil dibuat sebagai draft.',
      position: 'top',
      timeout: 3000,
    })
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: err.message || 'Gagal membuat catatan.',
      position: 'top',
      timeout: 3000,
    })
  }
}

function confirmSubmit(note) {
  $q.dialog({
    title: 'Konfirmasi Submit',
    message: `Submit catatan "${note.title}" untuk ${note.studentName}? Setelah di-submit, catatan tidak bisa diubah.`,
    cancel: true,
    persistent: true,
  }).onOk(async () => {
    try {
      await submitNote(note)
      $q.notify({
        type: 'positive',
        message: 'Catatan berhasil di-submit.',
        position: 'top',
        timeout: 3000,
      })
    } catch (err) {
      $q.notify({
        type: 'negative',
        message: err.message || 'Gagal submit catatan.',
        position: 'top',
        timeout: 3000,
      })
    }
  })
}

function confirmDelete(note) {
  $q.dialog({
    title: 'Konfirmasi Hapus',
    message: `Hapus draft catatan "${note.title}"?`,
    cancel: true,
    persistent: true,
  }).onOk(async () => {
    try {
      await deleteNote(note)
      $q.notify({
        type: 'positive',
        message: 'Draft catatan berhasil dihapus.',
        position: 'top',
        timeout: 3000,
      })
    } catch (err) {
      $q.notify({
        type: 'negative',
        message: err.message || 'Gagal menghapus catatan.',
        position: 'top',
        timeout: 3000,
      })
    }
  })
}

function viewNote(note) {
  selectedNote.value = note
  showDetailDialog.value = true
}

function getCategoryColor(category) {
  switch (category) {
    case 'POSITIVE': return 'positive'
    case 'NEGATIVE': return 'negative'
    case 'NEUTRAL': return 'grey'
    case 'INCIDENT': return 'deep-orange'
    default: return 'grey'
  }
}

function getCategoryLabel(category) {
  switch (category) {
    case 'POSITIVE': return 'Positif'
    case 'NEGATIVE': return 'Negatif'
    case 'NEUTRAL': return 'Netral'
    case 'INCIDENT': return 'Insiden'
    default: return category
  }
}

function getPriorityColor(priority) {
  switch (priority) {
    case 'LOW': return 'grey'
    case 'MEDIUM': return 'primary'
    case 'HIGH': return 'warning'
    case 'URGENT': return 'negative'
    default: return 'grey'
  }
}

onMounted(() => {
  loadNotes()
  loadStudents()
})
</script>

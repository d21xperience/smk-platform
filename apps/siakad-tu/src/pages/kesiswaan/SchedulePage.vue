<template>
  <q-page class="q-pa-md">
    <div class="row items-center q-mb-md">
      <div class="text-h5 text-weight-bold">Jadwal Pelajaran</div>
      <q-space />
      <q-btn color="secondary" label="Impor Excel" icon="upload_file" class="q-mr-sm" @click="triggerImport" />
      <q-btn color="primary" label="Tambah Jadwal" icon="add" @click="openDialog" />
      <input type="file" ref="fileInput" style="display: none" accept=".xlsx, .xls" @change="handleImportExcel" />
    </div>

    <q-table
      :rows="schedules"
      :columns="columns"
      row-key="id"
      flat bordered
      :loading="isLoading"
    >
      <template v-slot:body-cell-day="props">
        <q-td :props="props">
          <q-badge color="blue-6">{{ props.row.day }}</q-badge>
        </q-td>
      </template>
      <template v-slot:body-cell-time="props">
        <q-td :props="props">
          {{ props.row.start_time }} - {{ props.row.end_time }}
        </q-td>
      </template>
    </q-table>

    <!-- Modal Tambah Jadwal -->
    <q-dialog v-model="showDialog" persistent>
      <q-card style="min-width: 500px">
        <q-card-section class="row items-center">
          <div class="text-h6">Tambah Jadwal Pelajaran</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pt-none q-gutter-y-md">
          <q-select v-model="form.class_id" :options="classOptions" label="Pilih Kelas" outlined dense emit-value map-options />
          <q-select v-model="form.subject_id" :options="subjectOptions" label="Pilih Mata Pelajaran" outlined dense emit-value map-options />
          <q-select v-model="form.teacher_id" :options="teacherOptions" label="Pilih Guru Pengampu" outlined dense emit-value map-options />
          
          <div class="row q-col-gutter-sm">
            <div class="col-4">
              <q-select v-model="form.day" :options="dayOptions" label="Hari" outlined dense />
            </div>
            <div class="col-4">
              <q-input v-model="form.start_time" label="Jam Mulai" outlined dense mask="##:##" fill-mask="0" />
            </div>
            <div class="col-4">
              <q-input v-model="form.end_time" label="Jam Selesai" outlined dense mask="##:##" fill-mask="0" />
            </div>
          </div>
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn color="primary" label="Simpan" :loading="isSubmitting" @click="submitSchedule" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { Notify } from 'quasar'
import { api } from '@/boot/axios'
import * as XLSX from 'xlsx'

const schedules = ref([])
const isLoading = ref(false)
const isSubmitting = ref(false)
const showDialog = ref(false)
const fileInput = ref(null)

const classes = ref([])
const subjects = ref([])
const teachers = ref([])

const form = ref({
  class_id: '',
  subject_id: null,
  teacher_id: '',
  day: 'Senin',
  start_time: '07:00',
  end_time: '08:30'
})

const dayOptions = ['Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat', 'Sabtu']

const classOptions = computed(() => classes.value.map(c => ({ label: c.nm_kelas || c.Nama, value: c.rombongan_belajar_id || c.ID })))
const subjectOptions = computed(() => subjects.value.map(s => ({ label: s.name, value: s.id })))
const teacherOptions = computed(() => teachers.value.map(t => ({ label: t.name, value: t.id })))

const columns = [
  { name: 'day', label: 'Hari', field: 'day', align: 'left', sortable: true },
  { name: 'time', label: 'Jam', align: 'left' },
  { name: 'class_name', label: 'Kelas', field: 'class_name', align: 'left' },
  { name: 'subject_name', label: 'Mata Pelajaran', field: 'subject_name', align: 'left' },
  { name: 'teacher_name', label: 'Guru Pengampu', field: 'teacher_name', align: 'left' }
]

async function fetchData() {
  isLoading.value = true
  try {
    // PERBAIKAN: Endpoint harusnya /schedules, pastikan backend sudah direstart
    const resSchedules = await api.get('/schedules')
    schedules.value = resSchedules.data || []

    const resOptions = await api.get('/schedules/options')
    classes.value = resOptions.data.classes || []
    subjects.value = resOptions.data.subjects || []
    teachers.value = resOptions.data.teachers || []
  } catch (err) {
    console.error('Fetch error:', err)
    Notify.create({ type: 'negative', message: 'Gagal memuat data. Pastikan backend sudah direstart.' })
  } finally {
    isLoading.value = false
  }
}

function openDialog() {
  form.value = { class_id: '', subject_id: null, teacher_id: '', day: 'Senin', start_time: '07:00', end_time: '08:30' }
  showDialog.value = true
}

async function submitSchedule() {
  if (!form.value.class_id || !form.value.subject_id || !form.value.teacher_id) {
    Notify.create({ type: 'warning', message: 'Lengkapi semua isian terlebih dahulu' })
    return
  }
  isSubmitting.value = true
  try {
    await api.post('/schedules', form.value)
    Notify.create({ type: 'positive', message: 'Jadwal berhasil ditambahkan' })
    showDialog.value = false
    await fetchData()
  } catch (err) {
    Notify.create({ type: 'negative', message: err.response?.data?.error || 'Gagal menyimpan jadwal' })
  } finally {
    isSubmitting.value = false
  }
}

function triggerImport() {
  fileInput.value.click()
}

function handleImportExcel(event) {
  const file = event.target.files[0]
  if (!file) return

  const reader = new FileReader()
  reader.onload = async (e) => {
    try {
      const data = new Uint8Array(e.target.result)
      const workbook = XLSX.read(data, { type: 'array' })
      const worksheet = workbook.Sheets[workbook.SheetNames[0]]
      const json = XLSX.utils.sheet_to_json(worksheet)
      
      if (json.length === 0) throw new Error("File Excel kosong")
      
      isSubmitting.value = true
      let successCount = 0
      
      for (const row of json) {
        try {
          await api.post('/schedules', {
            class_id: row['ID Kelas'],
            subject_id: parseInt(row['ID Mapel']),
            teacher_id: row['ID Guru'],
            day: row['Hari'],
            start_time: row['Mulai'],
            end_time: row['Selesai']
          })
          successCount++
        } catch (err) {
          console.error("Gagal import baris:", row, err)
        }
      }
      
      Notify.create({ type: 'positive', message: `${successCount} Jadwal berhasil diimpor` })
      await fetchData()
    } catch (err) {
      Notify.create({ type: 'negative', message: 'Gagal membaca file Excel: ' + err.message })
    } finally {
      isSubmitting.value = false
      event.target.value = ''
    }
  }
  reader.readAsArrayBuffer(file)
}

onMounted(fetchData)
</script>

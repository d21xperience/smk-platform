<template>
  <q-page class="q-pa-md">
    <div class="text-h5 q-mb-md">Kelola Ujian</div>

    <q-tabs v-model="tab" dense class="q-mb-md" align="left">
      <q-tab name="create" label="Buat Ujian" />
      <q-tab name="list" label="Daftar Ujian" />
      <q-tab name="schedule" label="Jadwal Ujian" />
    </q-tabs>

    <q-tab-panels v-model="tab" animated>
      <!-- Tab Buat Ujian -->
      <q-tab-panel name="create" class="q-pa-none">
        <q-card flat bordered>
          <q-card-section>
            <div class="text-h6">{{ isEditExam ? 'Edit Ujian' : 'Buat Ujian Baru' }}</div>
          </q-card-section>
          <q-card-section>
            <q-form @submit="submitExam" class="q-gutter-md">
              <q-input v-model="examForm.nama" label="Nama Ujian" outlined
                :rules="[val => !!val || 'Nama ujian harus diisi']" />
              <div class="row q-col-gutter-md">
                <div class="col-12 col-md-6">
                  <q-select v-model="examForm.kelas_id" :options="kelasOptions" label="Kelas" option-value="id"
                    option-label="nama" outlined :rules="[val => !!val]" />
                </div>
                <div class="col-12 col-md-6">
                  <q-select v-model="examForm.mapel_id" :options="mapelOptions" label="Mata Pelajaran" option-value="id"
                    option-label="nama" outlined :rules="[val => !!val]" />
                </div>
              </div>
              <div class="row q-col-gutter-md">
                <div class="col-12 col-md-6">
                  <q-input v-model="examForm.tanggal_mulai" label="Tanggal & Waktu Mulai" type="datetime-local" outlined
                    :rules="[val => !!val]" />
                </div>
                <div class="col-12 col-md-6">
                  <q-input v-model="examForm.tanggal_selesai" label="Tanggal & Waktu Selesai" type="datetime-local"
                    outlined :rules="[val => !!val]" />
                </div>
              </div>
              <div class="row q-col-gutter-md">
                <div class="col-12 col-md-4">
                  <q-input v-model.number="examForm.durasi" label="Durasi (menit)" type="number" outlined
                    :rules="[val => val > 0 || 'Durasi harus > 0']" />
                </div>
                <div class="col-12 col-md-4">
                  <q-select v-model="examForm.status" :options="['aktif', 'nonaktif']" label="Status" outlined />
                </div>
                <div class="col-12 col-md-4">
                  <q-input v-model.number="examForm.jumlah_soal" label="Jumlah Soal yang Ditampilkan" type="number"
                    outlined hint="Kosongkan = semua soal dari bank" />
                </div>
              </div>
              <div class="row q-col-gutter-md">
                <div class="col-12">
                  <q-select v-model="examForm.soal_ids" :options="soalOptions" label="Pilih Soal dari Bank" multiple
                    outlined use-chips option-value="id" option-label="pertanyaan_preview" :loading="loadingSoal"
                    @focus="loadSoalList" />
                </div>
              </div>
              <div class="q-gutter-sm">
                <q-btn label="Simpan" type="submit" color="primary" />
                <q-btn label="Reset" flat color="negative" @click="resetExamForm" />
              </div>
            </q-form>
          </q-card-section>
        </q-card>
      </q-tab-panel>

      <!-- Tab Daftar Ujian -->
      <q-tab-panel name="list" class="q-pa-none">
        <q-table :rows="examList" :columns="examColumns" row-key="id" flat bordered dense :loading="loading">
          <template v-slot:body-cell-status="props">
            <q-td :props="props">
              <q-badge :color="props.row.status === 'aktif' ? 'positive' : 'grey'" :label="props.row.status" />
            </q-td>
          </template>
          <template v-slot:body-cell-aksi="props">
            <q-td :props="props">
              <q-btn flat round dense icon="edit" color="info" @click="editExam(props.row)" class="q-mr-sm" />
              <q-btn flat round dense icon="delete" color="negative" @click="confirmDeleteExam(props.row)" />
            </q-td>
          </template>
        </q-table>
      </q-tab-panel>

      <!-- Tab Jadwal Ujian -->
      <q-tab-panel name="schedule" class="q-pa-none">
        <q-calendar v-model="selectedDate" view="month" :events="calendarEvents" @click:event="onEventClick"
          class="q-mb-md" style="height: 600px;" />
        <!-- Atau bisa gunakan tabel sederhana jika tidak ingin pakai q-calendar -->
        <div class="q-mt-md">
          <q-table :rows="upcomingExams" :columns="scheduleColumns" row-key="id" flat bordered dense
            title="Ujian Mendatang" />
        </div>
      </q-tab-panel>
    </q-tab-panels>

    <!-- Dialog Konfirmasi Hapus -->
    <q-dialog v-model="deleteDialog" persistent>
      <q-card>
        <q-card-section>
          <div class="text-h6">Konfirmasi Hapus</div>
        </q-card-section>
        <q-card-section>
          Apakah Anda yakin ingin menghapus ujian <strong>{{ deleteItem?.nama }}</strong>?
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" color="negative" v-close-popup />
          <q-btn flat label="Hapus" color="primary" @click="executeDelete" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
// Jika ingin menggunakan q-calendar, install: npm install @quasar/quasar-ui-qcalendar
// Untuk contoh ini, kita hanya akan menggunakan tabel sederhana dulu.

const $q = useQuasar()

// Data dummy (integrasi dari modul sebelumnya)
const kelasList = ref([
  { id: 1, nama: '10 IPA 1' }, { id: 2, nama: '10 IPA 2' }, { id: 3, nama: '11 IPA 1' }
])
const mapelList = ref([
  { id: 1, nama: 'Matematika' }, { id: 2, nama: 'Fisika' }, { id: 3, nama: 'Bahasa Indonesia' }
])
const soalList = ref([
  { id: 1, pertanyaan: 'Hasil 2+2?', kategori_nama: 'Matematika' },
  { id: 2, pertanyaan: 'Satuan gaya?', kategori_nama: 'Fisika' },
  { id: 3, pertanyaan: 'Apa ibu kota Indonesia?', kategori_nama: 'Pengetahuan Umum' }
])

// Data ujian dummy
const examList = ref([
  {
    id: 1,
    nama: 'UTS Matematika Ganjil',
    kelas_id: 1,
    kelas_nama: '10 IPA 1',
    mapel_id: 1,
    mapel_nama: 'Matematika',
    tanggal_mulai: '2025-06-10T08:00',
    tanggal_selesai: '2025-06-10T10:00',
    durasi: 120,
    status: 'aktif',
    jumlah_soal: null,
    soal_ids: [1, 2]
  },
  {
    id: 2,
    nama: 'UAS Fisika',
    kelas_id: 2,
    kelas_nama: '10 IPA 2',
    mapel_id: 2,
    mapel_nama: 'Fisika',
    tanggal_mulai: '2025-06-15T09:00',
    tanggal_selesai: '2025-06-15T11:00',
    durasi: 120,
    status: 'aktif',
    jumlah_soal: 20,
    soal_ids: []
  }
])

const tab = ref('create')
const loading = ref(false)
const loadingSoal = ref(false)

// Form state
const isEditExam = ref(false)
const editExamId = ref(null)
const examForm = ref({
  nama: '',
  kelas_id: null,
  mapel_id: null,
  tanggal_mulai: '',
  tanggal_selesai: '',
  durasi: 60,
  status: 'aktif',
  jumlah_soal: null,
  soal_ids: []
})

// Options
const kelasOptions = computed(() => kelasList.value)
const mapelOptions = computed(() => mapelList.value)
const soalOptions = ref([])

function loadSoalList() {
  if (soalOptions.value.length) return
  loadingSoal.value = true
  setTimeout(() => {
    soalOptions.value = soalList.value.map(s => ({
      id: s.id,
      pertanyaan_preview: `${s.id}. ${s.pertanyaan.substring(0, 50)}... (${s.kategori_nama})`
    }))
    loadingSoal.value = false
  }, 500)
}

// Exam CRUD
function resetExamForm() {
  isEditExam.value = false
  editExamId.value = null
  examForm.value = {
    nama: '',
    kelas_id: null,
    mapel_id: null,
    tanggal_mulai: '',
    tanggal_selesai: '',
    durasi: 60,
    status: 'aktif',
    jumlah_soal: null,
    soal_ids: []
  }
  tab.value = 'list' // arahkan ke daftar setelah reset
}

function submitExam() {
  // Validasi dasar
  if (!examForm.value.nama || !examForm.value.kelas_id || !examForm.value.mapel_id ||
    !examForm.value.tanggal_mulai || !examForm.value.tanggal_selesai || !examForm.value.durasi) {
    $q.notify({ type: 'negative', message: 'Lengkapi semua field yang wajib' })
    return
  }
  const kelas = kelasOptions.value.find(k => k.id === examForm.value.kelas_id)
  const mapel = mapelOptions.value.find(m => m.id === examForm.value.mapel_id)
  const examData = {
    ...examForm.value,
    kelas_nama: kelas?.nama,
    mapel_nama: mapel?.nama,
    id: isEditExam.value ? editExamId.value : Date.now()
  }
  if (isEditExam.value) {
    const idx = examList.value.findIndex(e => e.id === editExamId.value)
    if (idx !== -1) examList.value[idx] = examData
    $q.notify({ type: 'positive', message: 'Ujian berhasil diupdate' })
  } else {
    examList.value.push(examData)
    $q.notify({ type: 'positive', message: 'Ujian berhasil ditambahkan' })
  }
  resetExamForm()
}

function editExam(exam) {
  isEditExam.value = true
  editExamId.value = exam.id
  examForm.value = {
    nama: exam.nama,
    kelas_id: exam.kelas_id,
    mapel_id: exam.mapel_id,
    tanggal_mulai: exam.tanggal_mulai,
    tanggal_selesai: exam.tanggal_selesai,
    durasi: exam.durasi,
    status: exam.status,
    jumlah_soal: exam.jumlah_soal,
    soal_ids: exam.soal_ids || []
  }
  tab.value = 'create'
  if (exam.soal_ids && exam.soal_ids.length) loadSoalList() // preload
}

function confirmDeleteExam(exam) {
  deleteItem.value = exam
  deleteDialog.value = true
}

const deleteDialog = ref(false)
const deleteItem = ref(null)
function executeDelete() {
  const idx = examList.value.findIndex(e => e.id === deleteItem.value.id)
  if (idx !== -1) examList.value.splice(idx, 1)
  $q.notify({ type: 'positive', message: 'Ujian dihapus' })
  deleteDialog.value = false
  deleteItem.value = null
}

// Tabel daftar ujian
const examColumns = [
  { name: 'nama', label: 'Nama Ujian', field: 'nama', align: 'left' },
  { name: 'kelas_nama', label: 'Kelas', field: 'kelas_nama', align: 'left' },
  { name: 'mapel_nama', label: 'Mata Pelajaran', field: 'mapel_nama', align: 'left' },
  { name: 'tanggal_mulai', label: 'Mulai', field: 'tanggal_mulai', align: 'left' },
  { name: 'durasi', label: 'Durasi (menit)', field: 'durasi', align: 'center' },
  { name: 'status', label: 'Status', field: 'status', align: 'center' },
  { name: 'aksi', label: 'Aksi', field: 'aksi', align: 'center' }
]

// Jadwal ujian (upcoming)
const upcomingExams = computed(() => {
  const now = new Date().toISOString()
  return examList.value.filter(e => e.tanggal_mulai > now).sort((a, b) => a.tanggal_mulai.localeCompare(b.tanggal_mulai))
})
const scheduleColumns = [
  { name: 'nama', label: 'Ujian', field: 'nama' },
  { name: 'kelas_nama', label: 'Kelas', field: 'kelas_nama' },
  { name: 'mapel_nama', label: 'Mapel', field: 'mapel_nama' },
  { name: 'tanggal_mulai', label: 'Mulai', field: 'tanggal_mulai' },
  { name: 'durasi', label: 'Durasi', field: 'durasi' }
]

// Untuk calendar events (opsional, butuh quasar-calendar)
const selectedDate = ref(new Date().toISOString().substr(0, 10))
const calendarEvents = computed(() => {
  return examList.value.map(e => ({
    id: e.id,
    title: e.nama,
    start: e.tanggal_mulai,
    end: e.tanggal_selesai,
    color: e.status === 'aktif' ? 'green' : 'grey'
  }))
})
function onEventClick(event) {
  $q.dialog({ message: `Detail ujian: ${event.title}`, ok: 'Tutup' })
}

onMounted(() => {
  // bisa preload
})
</script>

<style scoped>
/* optional */
</style>

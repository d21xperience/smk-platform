<template>
  <q-page class="q-pa-md">
    <!-- HEADER & FILTER (Sama seperti sebelumnya) -->
    <div class="row items-center justify-between q-mb-md">
      <div class="col">
        <div class="row items-center q-gutter-sm">
          <q-btn v-if="cameFromSchedule" flat dense round icon="arrow_back" color="grey-7" @click="router.back()" />
          <div>
            <div class="text-h5 text-weight-bold text-grey-9">Absensi Siswa</div>
            <div class="text-caption text-grey-7">
              <!-- Periode: <span class="text-weight-bold text-primary">{{ academicStore.activePeriodCode }}</span> -->
            </div>
          </div>
        </div>
      </div>
    </div>

    <q-card flat bordered class="q-mb-md">
      <q-card-section class="row q-col-gutter-md">
        <div class="col-12 col-md-6">
          <q-input v-model="selectedDate" outlined dense label="Tanggal" mask="####-##-##">
            <template v-slot:append>
              <q-icon name="event" class="cursor-pointer">
                <q-popup-proxy cover transition-show="scale" transition-hide="scale">
                  <q-date v-model="selectedDate" mask="YYYY-MM-DD" @update:model-value="onDateChange" />
                </q-popup-proxy>
              </q-icon>
            </template>
          </q-input>
        </div>
        <div class="col-12 col-md-6">
          <q-select v-model="selectedClass" :options="classOptions" option-label="name" option-value="id" outlined dense
            label="Pilih Kelas" emit-value map-options @update:model-value="onClassChange" />
        </div>
      </q-card-section>
    </q-card>

    <!-- JIKA KELAS SUDAH DIPILIH -->
    <template v-if="studentList.length > 0">
      <!-- SUMMARY CARDS -->
      <div class="row q-col-gutter-sm q-mb-md">
        <div class="col-3" v-for="(count, label) in summaryLabels" :key="label">
          <q-card :class="`bg-${count.color}-1 text-${count.color}-10`" flat bordered>
            <q-card-section class="q-pa-sm text-center">
              <div class="text-h6 text-weight-bold">{{ attendanceSummary[label] }}</div>
              <div class="text-caption">{{ count.text }}</div>
            </q-card-section>
          </q-card>
        </div>
      </div>

      <!-- FITUR UTAMA: DAFTAR KETIDAKHADIRAN -->
      <q-card flat bordered class="q-mb-md">
        <q-card-section class="row items-center justify-between">
          <div class="text-subtitle1 text-weight-bold text-grey-9">
            <q-icon name="person_off" class="q-mr-xs" /> Catatan Ketidakhadiran
          </div>
          <q-btn unelevated no-caps color="primary" icon="add" label="Catat Tidak Hadir" size="md"
            @click="openAbsenceDialog()" />
        </q-card-section>

        <q-separator />

        <!-- List Siswa yang S/I/A -->
        <q-list v-if="absentStudents.length > 0" separator>
          <q-item v-for="student in absentStudents" :key="student.id">
            <q-item-section>
              <q-item-label class="text-weight-bold">{{ student.studentName }}</q-item-label>
              <q-item-label caption v-if="student.keterangan">
                <q-icon name="notes" size="xs" /> {{ student.keterangan }}
              </q-item-label>
            </q-item-section>
            <q-item-section side>
              <div class="row items-center q-gutter-sm">
                <q-badge :color="getStatusColor(student.status)" :label="student.status" />
                <q-btn flat dense round icon="edit" color="grey-7" size="sm" @click="openAbsenceDialog(student)" />
                <q-btn flat dense round icon="check_circle" color="positive" size="sm"
                  @click="handleClearAbsence(student.studentId)">
                  <q-tooltip>Tandai Hadir</q-tooltip>
                </q-btn>
              </div>
            </q-item-section>
          </q-item>
        </q-list>

        <!-- Empty State untuk Ketidakhadiran -->
        <q-card-section v-else class="text-center q-pa-lg text-grey-6">
          <q-icon name="check_circle_outline" size="40px" color="positive" />
          <div class="text-body2 q-mt-sm">Semua siswa tercatat hadir. Luar biasa!</div>
        </q-card-section>
      </q-card>

      <!-- TOMBOL SIMPAN -->
      <q-btn unelevated color="primary" label="Simpan Absensi" icon="save"
        class="full-width q-py-sm text-weight-bold shadow-3" :loading="loadingSubmit" @click="handleSubmit" />
    </template>

    <!-- DIALOG INPUT KETIDAKHADIRAN -->
    <q-dialog v-model="showAbsenceDialog" persistent>
      <q-card style="min-width: 350px; max-width: 500px;">
        <q-card-section class="row items-center justify-between">
          <div class="text-h6">{{ isEditing ? 'Edit' : 'Catat' }} Ketidakhadiran</div>
          <q-btn flat dense round icon="close" v-close-popup />
        </q-card-section>

        <q-card-section class="q-gutter-md">
          <BaseSelectSearch v-model="formAbsence.studentId" :options="studentList" option-label="studentName"
            option-value="studentId" label="Pilih atau Ketik Nama Siswa" :disable="isEditing"
            :rules="[val => !!val || 'Wajib pilih siswa']" />

          <q-btn-toggle v-model="formAbsence.status" no-caps unelevated toggle-color="primary" color="white"
            text-color="primary" class="" :options="[
              { label: 'Sakit', value: 'sakit' },
              { label: 'Izin', value: 'izin' },
              { label: 'Alpha', value: 'alpha' }
            ]" />

          <q-input v-model="formAbsence.keterangan" outlined dense
            label="Keterangan (Misal: Sakit demam, Izin pernikahan)" type="textarea" rows="2" />
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md">
          <q-btn flat label="Batal" color="grey" v-close-popup />
          <q-btn unelevated label="Simpan Status" color="primary" @click="handleSaveAbsence" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useQuasar } from 'quasar'
import { useAttendance } from '@/composables/useAttendance'
import { useTeaching } from '@/composables/useTeaching'
import BaseSelectSearch from '@/components/ui/BaseSelectSearch.vue'

const $q = useQuasar()
const route = useRoute()
const router = useRouter()

const attendance = useAttendance()
const teaching = useTeaching()
// const context = useContext()

// State filter
const today = new Date().toISOString().slice(0, 10)
const selectedDate = ref(today)
const selectedClass = ref(null)
const cameFromSchedule = ref(false)
const loadingSubmit = ref(false)

// Getter yang disesuaikan
const students = ref(attendance.students)
const summary = ref(attendance.summary)
const studentList = computed(() => students.value)
const absentStudents = computed(() => students.value.filter(s => s.status !== 'hadir'))
const attendanceSummary = computed(() => summary.value || { hadir: 0, izin: 0, sakit: 0, alpha: 0 })

// Opsi kelas dari sesi mengajar hari ini
const classOptions = computed(() => {
  const sessions = teaching.sessions.value
  console.log('AttencancePage', sessions)
  // Ambil unique kelas dari jadwal
  const unique = new Map()
  sessions.forEach(s => {
    if (!unique.has(s.className)) {
      unique.set(s.className, { id: s.className, name: s.className }) // class_name sebagai id
    }
  })
  return Array.from(unique.values())
})

// Periode dari operational context
// const activePeriod = computed(() => {
//   const op = context.operational.value
//   return `${op.academicYearName} - ${op.semesterName}`
// })

// State dialog
const showAbsenceDialog = ref(false)
const isEditing = ref(false)
const formAbsence = ref({ studentId: null, status: 'S', keterangan: '' })

const summaryLabels = {
  hadir: { text: 'Hadir', color: 'positive' },
  sakit: { text: 'Sakit', color: 'warning' },
  izin: { text: 'Izin', color: 'info' },
  alpha: { text: 'Alpha', color: 'negative' }
}

// Load jadwal untuk mendapatkan kelas
async function loadSessions() {
  if (!selectedDate.value) return
  await teaching.loadSessions(selectedDate.value)
}

// Load attendance berdasarkan sesi yang sesuai dengan kelas dan tanggal
async function loadAttendance() {
  if (!selectedClass.value || !selectedDate.value) return
  // Cari sesi yang cocok
  const session = teaching.sessions.value.find(
    s => s.className === selectedClass.value && s.date === selectedDate.value
  )
  if (session) {
    await attendance.loadOrCreateAttendance(session)
  }
}

// Event handler
async function onDateChange() {
  await loadSessions()
  // Reset pilihan kelas
  selectedClass.value = null
}

async function onClassChange() {
  await loadAttendance()
}

// Buka dialog untuk mencatat ketidakhadiran
function openAbsenceDialog(student = null) {
  if (student) {
    isEditing.value = true
    formAbsence.value = {
      studentId: student.studentId,
      status: student.status,
      keterangan: student.note || ''
    }
  } else {
    isEditing.value = false
    formAbsence.value = { studentId: null, status: 'S', keterangan: '' }
  }
  showAbsenceDialog.value = true
}
const getStatusColor = (status) => {
  const colors = { sakit: 'warning', izin: 'info', alpha: 'negative' }
  return colors[status] || 'grey'
}
// Simpan status ketidakhadiran (panggil updateStatus)
function handleSaveAbsence() {
  if (!formAbsence.value.studentId) {
    $q.notify({ type: 'warning', message: 'Pilih siswa terlebih dahulu' })
    return
  }
  attendance.updateStudentStatus(formAbsence.value.studentId, formAbsence.value.status)
  if (formAbsence.value.keterangan) {
    attendance.updateStudentNote(formAbsence.value.studentId, formAbsence.value.keterangan)
  }
  showAbsenceDialog.value = false
  $q.notify({ type: 'positive', message: 'Status ketidakhadiran berhasil dicatat' })
}

// Tandai hadir (clear absence)
function handleClearAbsence(studentId) {
  attendance.updateStudentStatus(studentId, 'hadir')
  $q.notify({ type: 'info', message: 'Siswa ditandai sebagai Hadir' })
}

// Submit optimized payload
async function handleSubmit() {
  if (!selectedClass.value) return

  // Cari session ID
  const session = teaching.sessions.value.find(
    s => s.className === selectedClass.value && s.date === selectedDate.value
  )
  if (!session) {
    $q.notify({ type: 'warning', message: 'Tidak ada sesi mengajar untuk kelas ini' })
    return
  }

  // Pastikan attendance sudah ter-load untuk sesi ini
  if (!attendance.currentRecord || attendance.currentRecord.sessionId !== session.id) {
    await attendance.loadOrCreateAttendance(session)
  }

  // Kumpulkan data absences (hanya yang tidak hadir)
  const absences = students.value
    .filter(s => s.status !== 'hadir')
    .map(s => ({
      student_id: s.studentId,
      status: s.status,
      keterangan: s.note
    }))

  // eslint-disable-next-line no-unused-vars
  const payload = {
    class_id: selectedClass.value,
    date: selectedDate.value,
    total_students: students.value.length,
    absences
  }

  loadingSubmit.value = true
  try {
    // Gunakan action submitAttendance dari store, tapi kita perlu mengirim payload khusus
    // Untuk saat ini, kita akan submit melalui store dengan asumsi data sudah tersimpan.
    await attendance.submitAttendance()  // ini akan submit record yang sudah ada
    $q.notify({ type: 'positive', message: 'Absensi berhasil disimpan!', position: 'top' })
  } catch (e) {
    console.error(e)
    $q.notify({ type: 'negative', message: e.message || 'Gagal menyimpan absensi', position: 'top' })
  } finally {
    loadingSubmit.value = false
  }
}

// Lifecycle
onMounted(async () => {
  // Cek query params jika dari jadwal
  if (route.query.class_id) {
    cameFromSchedule.value = true
    selectedClass.value = route.query.class_id
  }
  await loadSessions()
  if (selectedClass.value) await loadAttendance()
})

// Watch tanggal untuk refresh
watch(selectedDate, async () => {
  await loadSessions()
  selectedClass.value = null
})
</script>

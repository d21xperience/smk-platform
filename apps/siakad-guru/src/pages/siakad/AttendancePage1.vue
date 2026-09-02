<template>
  <q-page class="q-pa-md">
    <!-- HEADER & FILTER (Sama seperti sebelumnya) -->
    <div class="row items-center justify-between q-mb-md">
      <div class="col">
        <div class="row items-center q-gutter-sm">
          <q-btn
            v-if="cameFromSchedule"
            flat
            dense
            round
            icon="arrow_back"
            color="grey-7"
            @click="router.back()"
          />
          <div>
            <div class="text-h5 text-weight-bold text-grey-9">Absensi Siswa</div>
            <div class="text-caption text-grey-7">
              Periode:
              <span class="text-weight-bold text-primary">{{
                academicStore.activePeriodCode
              }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <q-card flat bordered class="q-mb-md">
      <q-card-section class="row q-col-gutter-md">
        <div class="col-12 col-md-6">
          <q-input
            v-model="selectedDate"
            outlined
            dense
            label="Tanggal"
            mask="####-##-##"
            :disable="loading"
          >
            <template v-slot:append>
              <q-icon name="event" class="cursor-pointer">
                <q-popup-proxy cover transition-show="scale" transition-hide="scale">
                  <q-date
                    v-model="selectedDate"
                    mask="YYYY-MM-DD"
                    @update:model-value="onDateChange"
                  />
                </q-popup-proxy>
              </q-icon>
            </template>
          </q-input>
        </div>
        <div class="col-12 col-md-6">
          <q-select
            v-model="selectedClass"
            :options="classOptions"
            option-label="name"
            option-value="id"
            outlined
            dense
            label="Pilih Kelas"
            emit-value
            map-options
            :loading="loading"
            @update:model-value="onClassChange"
          />
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
          <q-btn
            unelevated
            no-caps
            color="primary"
            icon="add"
            label="Catat Tidak Hadir"
            size="md"
            @click="openAbsenceDialog()"
          />
        </q-card-section>

        <q-separator />

        <!-- List Siswa yang S/I/A -->
        <q-list v-if="absentStudents.length > 0" separator>
          <q-item v-for="student in absentStudents" :key="student.id">
            <q-item-section>
              <q-item-label class="text-weight-bold">{{ student.name }}</q-item-label>
              <q-item-label caption v-if="student.keterangan">
                <q-icon name="notes" size="xs" /> {{ student.keterangan }}
              </q-item-label>
            </q-item-section>
            <q-item-section side>
              <div class="row items-center q-gutter-sm">
                <q-badge :color="getStatusColor(student.status)" :label="student.status" />
                <q-btn
                  flat
                  dense
                  round
                  icon="edit"
                  color="grey-7"
                  size="sm"
                  @click="openAbsenceDialog(student)"
                />
                <q-btn
                  flat
                  dense
                  round
                  icon="check_circle"
                  color="positive"
                  size="sm"
                  @click="handleClearAbsence(student.id)"
                >
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
      <q-btn
        unelevated
        color="primary"
        label="Simpan Absensi"
        icon="save"
        class="full-width q-py-sm text-weight-bold shadow-3"
        :loading="loadingSubmit"
        @click="handleSubmit"
      />
    </template>

    <!-- DIALOG INPUT KETIDAKHADIRAN -->
    <q-dialog v-model="showAbsenceDialog" persistent>
      <q-card style="min-width: 350px; max-width: 500px">
        <q-card-section class="row items-center justify-between">
          <div class="text-h6">{{ isEditing ? 'Edit' : 'Catat' }} Ketidakhadiran</div>
          <q-btn flat dense round icon="close" v-close-popup />
        </q-card-section>

        <q-card-section class="q-gutter-md">
          <!-- <q-select v-model="formAbsence.studentId" :options="studentList" option-label="name" option-value="id"
            emit-value map-options outlined dense label="Pilih atau Ketik Nama Siswa" :disable="isEditing" use-input
            fill-input hide-selected input-debounce="0" @filter="filterFn"
            :rules="[val => !!val || 'Wajib pilih siswa']" /> -->
          <BaseSelectSearch
            v-model="formAbsence.studentId"
            :options="studentList"
            option-label="name"
            option-value="id"
            label="Pilih atau Ketik Nama Siswa"
            :disable="isEditing"
            :rules="[(val) => !!val || 'Wajib pilih siswa']"
          />

          <q-btn-toggle
            v-model="formAbsence.status"
            no-caps
            unelevated
            toggle-color="primary"
            color="white"
            text-color="primary"
            class=""
            :options="[
              { label: 'Sakit (S)', value: 'S' },
              { label: 'Izin (I)', value: 'I' },
              { label: 'Alpha (A)', value: 'A' },
            ]"
          />

          <q-input
            v-model="formAbsence.keterangan"
            outlined
            dense
            label="Keterangan (Misal: Sakit demam, Izin pernikahan)"
            type="textarea"
            rows="2"
          />
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
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useQuasar } from 'quasar'
import { useAttendance } from '@/composables/useAttendance'
import { useAcademicStore } from '@/stores/academicStore'
import BaseSelectSearch from '@/components/ui/BaseSelectSearch.vue'

const $q = useQuasar()
const route = useRoute()
const router = useRouter()
const academicStore = useAcademicStore()

const {
  loading,
  classOptions,
  studentList,
  absentStudents,
  attendanceSummary,
  fetchSchedule,
  fetchStudents,
  setAbsence,
  clearAbsence,
  submitAttendance,
} = useAttendance()

// --- STATE FILTER ---
const today = new Date().toISOString().split('T')[0]
const selectedDate = ref(today)
const selectedClass = ref(null)
const cameFromSchedule = ref(false)
const loadingSubmit = ref(false)

const summaryLabels = {
  hadir: { text: 'Hadir', color: 'positive' },
  sakit: { text: 'Sakit', color: 'warning' },
  izin: { text: 'Izin', color: 'info' },
  alpha: { text: 'Alpha', color: 'negative' },
}

// --- STATE DIALOG ABSENCE ---
const showAbsenceDialog = ref(false)
const isEditing = ref(false)
const formAbsence = ref({
  studentId: null,
  status: 'S',
  keterangan: '',
})

// --- METHODS FILTER ---
const onDateChange = () => {
  studentList.value = []
  if (selectedClass.value) fetchStudents(selectedClass.value, selectedDate.value)
}

const onClassChange = (classId) => {
  if (classId) fetchStudents(classId, selectedDate.value)
  else studentList.value = []
}

// --- METHODS DIALOG ---
const openAbsenceDialog = (student = null) => {
  if (student) {
    // Mode Edit
    isEditing.value = true
    formAbsence.value = {
      studentId: student.id,
      status: student.status,
      keterangan: student.keterangan,
    }
  } else {
    // Mode Tambah Baru
    isEditing.value = false
    formAbsence.value = { studentId: null, status: 'S', keterangan: '' }
  }
  showAbsenceDialog.value = true
}

const handleSaveAbsence = () => {
  if (!formAbsence.value.studentId) {
    $q.notify({ type: 'warning', message: 'Pilih siswa terlebih dahulu' })
    return
  }

  setAbsence(formAbsence.value.studentId, formAbsence.value.status, formAbsence.value.keterangan)
  showAbsenceDialog.value = false
  $q.notify({ type: 'positive', message: 'Status ketidakhadiran berhasil dicatat' })
}

const handleClearAbsence = (studentId) => {
  clearAbsence(studentId)
  $q.notify({ type: 'info', message: 'Siswa ditandai sebagai Hadir' })
}

const getStatusColor = (status) => {
  const colors = { S: 'warning', I: 'info', A: 'negative' }
  return colors[status] || 'grey'
}

// --- METHODS SUBMIT (OPTIMIZED PAYLOAD) ---
const handleSubmit = async () => {
  if (!selectedClass.value) return

  // 🚀 OPTIMIZED PAYLOAD FOR GOLANG BACKEND
  // Hanya mengirim data yang TIDAK HADIR. Backend akan mengasumsikan sisanya adalah 'H' (Hadir).
  const payload = {
    class_id: selectedClass.value,
    date: selectedDate.value,
    total_students: studentList.value.length,
    absences: absentStudents.value.map((s) => ({
      student_id: s.id,
      status: s.status,
      keterangan: s.keterangan,
    })),
  }

  loadingSubmit.value = true
  try {
    await submitAttendance(payload)
    $q.notify({ type: 'positive', message: 'Absensi berhasil disimpan!', position: 'top' })
  } catch (e) {
    console.log(e)
    $q.notify({ type: 'negative', message: 'Gagal menyimpan absensi', position: 'top' })
  } finally {
    loadingSubmit.value = false
  }
}

// --- LIFECYCLE ---
onMounted(async () => {
  if (route.query.class_id) {
    cameFromSchedule.value = true
    selectedClass.value = route.query.class_id
    if (route.query.class_name) {
      const exists = classOptions.value.find((c) => c.id == route.query.class_id)
      if (!exists)
        classOptions.value.push({ id: route.query.class_id, name: route.query.class_name })
    }
  }

  await fetchSchedule(selectedDate.value)
  if (selectedClass.value) await fetchStudents(selectedClass.value, selectedDate.value)
})
</script>

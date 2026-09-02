<!--
FILE: src/pages/attendance/AttendanceSessionPage.vue
STATUS: MODIFY
STATUS IMPLEMENTASI: COMPLETE
-->

<template>
  <q-page padding>
    <div class="row items-center q-mb-md">
      <q-btn flat dense icon="arrow_back" label="Kembali" :to="{ name: 'attendance' }" class="q-mr-sm" />

      <div class="text-h5">Absensi Siswa</div>
    </div>

    <ReadOnlyContextBanner />

    <q-banner v-if="error" class="bg-negative text-white q-mb-md">
      {{ error }}
    </q-banner>

    <q-card v-if="isLoading" class="q-mb-md">
      <q-card-section class="text-center">
        <q-spinner color="primary" size="2em" />

        <div class="q-mt-sm">Loading...</div>
      </q-card-section>
    </q-card>

    <template v-else-if="session">
      <q-card class="q-mb-md">
        <q-card-section>
          <div class="text-h6">
            {{ session.subjectName }} - {{ session.className }}
          </div>

          <div class="text-caption text-grey">
            {{ session.date }} | Status:

            <q-badge :color="isSubmitted ? 'blue' : 'orange'">
              {{ session.status }}
            </q-badge>
          </div>
        </q-card-section>
      </q-card>

      <q-card v-if="summary" class="q-mb-md">
        <q-card-section>
          <div class="row q-gutter-md text-center">
            <div class="col">
              <div class="text-h6">{{ summary.total }}</div>
              <div class="text-caption">Total</div>
            </div>

            <div class="col">
              <div class="text-h6 text-positive">{{ summary.present }}</div>
              <div class="text-caption">Hadir</div>
            </div>

            <div class="col">
              <div class="text-h6 text-warning">{{ summary.sick }}</div>
              <div class="text-caption">Sakit</div>
            </div>

            <div class="col">
              <div class="text-h6 text-info">{{ summary.permission }}</div>
              <div class="text-caption">Izin</div>
            </div>

            <div class="col">
              <div class="text-h6 text-negative">{{ summary.absent }}</div>
              <div class="text-caption">Absen</div>
            </div>
          </div>
        </q-card-section>
      </q-card>

      <q-card class="q-mb-md">
        <q-card-section>
          <div class="text-h6 q-mb-md">Daftar Siswa</div>

          <q-list separator>
            <q-item v-for="student in students" :key="student.studentId">
              <q-item-section>
                <q-item-label>{{ student.studentName }}</q-item-label>
              </q-item-section>

              <q-item-section side class="q-gutter-sm row items-center">
                <q-select :model-value="student.status" :options="statusOptions" dense outlined style="width: 130px"
                  emit-value map-options :disable="isSubmitted || isReadOnlyContext"
                  @update:model-value="handleStatusChange(student.studentId, $event)" />

                <q-input :model-value="student.note" dense outlined placeholder="Catatan" style="width: 200px"
                  :disable="isSubmitted || isReadOnlyContext"
                  @update:model-value="handleNoteChange(student.studentId, $event)" @blur="saveDraft" />
              </q-item-section>
            </q-item>
          </q-list>
        </q-card-section>
      </q-card>

      <q-card v-if="isReadOnlyContext">
        <q-card-section class="text-center text-grey">
          Absensi historis hanya dapat dilihat. Simpan draft dan submit tidak tersedia.
        </q-card-section>
      </q-card>

      <q-card v-else-if="isDraft">
        <q-card-section class="row q-gutter-sm">
          <q-btn color="primary" label="Simpan Draft" :loading="isSaving" @click="saveDraft" />

          <q-btn color="positive" label="Submit Absensi" :loading="isLoading" @click="handleSubmit" />
        </q-card-section>
      </q-card>

      <q-card v-else-if="isSubmitted">
        <q-card-section class="text-center text-grey">
          Absensi telah disubmit dan tidak dapat diubah.
        </q-card-section>
      </q-card>
    </template>

    <q-card v-else-if="isReadOnlyContext">
      <q-card-section class="text-center text-grey">
        Data absensi untuk periode historis ini tidak tersedia.
      </q-card-section>
    </q-card>

    <q-card v-else>
      <q-card-section class="text-center text-grey">
        Sesi absensi tidak ditemukan.
      </q-card-section>
    </q-card>
  </q-page>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useQuasar } from 'quasar'
import { useAttendance } from '../../composables/useAttendance.js'
import { useContext } from '../../composables/useContext.js'
import ReadOnlyContextBanner from '../../components/common/ReadOnlyContextBanner.vue'

const route = useRoute()
const router = useRouter()
const $q = useQuasar()

const { isReadOnlyContext } = useContext()

const {
  session,
  students,
  isLoading,
  isSaving,
  error,
  isDraft,
  isSubmitted,
  summary,
  loadOrCreateSession,
  updateRecordStatus,
  updateRecordNote,
  saveDraft,
  submitAttendance,
} = useAttendance()

const teachingSessionId = computed(() => route.params.teachingSessionId)

const statusOptions = [
  { label: 'Hadir', value: 'present' },
  { label: 'Sakit', value: 'sick' },
  { label: 'Izin', value: 'permission' },
  { label: 'Absen', value: 'absent' },
]

onMounted(async () => {
  try {
    await loadOrCreateSession(teachingSessionId.value)
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: err.message,
    })

    router.push({ name: 'attendance' })
  }
})

const handleStatusChange = async (studentId, status) => {
  if (isReadOnlyContext.value) return

  await updateRecordStatus(studentId, status)
}

const handleNoteChange = (studentId, note) => {
  if (isReadOnlyContext.value) return

  updateRecordNote(studentId, note)
}

const handleSubmit = async () => {
  if (isReadOnlyContext.value) {
    $q.notify({
      type: 'warning',
      message: 'Data historis bersifat read-only',
    })

    return
  }

  try {
    await submitAttendance()

    $q.notify({
      type: 'positive',
      message: 'Absensi berhasil disubmit',
    })

    router.push({ name: 'attendance' })
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: err.message,
    })
  }
}
</script>

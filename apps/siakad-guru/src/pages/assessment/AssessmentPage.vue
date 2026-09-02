<!--
FILE: src/pages/assessment/AssessmentPage.vue
STATUS: MODIFY
STATUS IMPLEMENTASI: COMPLETE
-->

<template>
  <q-page padding>
    <ReadOnlyContextBanner />

    <div class="text-h5 q-mb-md">Assessment</div>

    <q-banner v-if="error" class="bg-negative text-white q-mb-md">
      {{ error }}
    </q-banner>

    <q-card class="q-mb-md">
      <q-card-section>
        <div class="text-h6 q-mb-md">Buat / Muat Penilaian</div>

        <q-form class="q-gutter-md" @submit="handleLoadOrCreate">
          <q-select v-model="selectedClass" :options="classOptions" label="Kelas" outlined dense emit-value map-options
            :disable="isReadOnlyContext" />

          <q-select v-model="selectedSubject" :options="subjectOptions" label="Mata Pelajaran" outlined dense emit-value
            map-options :disable="isReadOnlyContext" />

          <q-input v-model="assessmentDate" label="Tanggal" type="date" outlined dense :disable="isReadOnlyContext" />

          <q-btn type="submit" label="Load / Create" color="primary" :loading="isLoading"
            :disable="isReadOnlyContext" />
        </q-form>
      </q-card-section>
    </q-card>

    <q-card v-if="assessment">
      <q-card-section>
        <div class="text-h6 q-mb-md">Penilaian Ditemukan</div>

        <q-list dense>
          <q-item>
            <q-item-section>
              <q-item-label>Kelas</q-item-label>
            </q-item-section>

            <q-item-section side>
              <q-item-label caption>
                {{ assessment.className }}
              </q-item-label>
            </q-item-section>
          </q-item>

          <q-item>
            <q-item-section>
              <q-item-label>Mata Pelajaran</q-item-label>
            </q-item-section>

            <q-item-section side>
              <q-item-label caption>
                {{ assessment.subjectName }}
              </q-item-label>
            </q-item-section>
          </q-item>

          <q-item>
            <q-item-section>
              <q-item-label>Status</q-item-label>
            </q-item-section>

            <q-item-section side>
              <q-badge :color="isFinalized ? 'blue' : 'orange'">
                {{ assessment.status }}
              </q-badge>
            </q-item-section>
          </q-item>
        </q-list>

        <q-btn color="primary" :label="isReadOnlyContext ? 'Lihat Penilaian' : 'Kelola Penilaian'" class="q-mt-md"
          @click="handleNavigate" />
      </q-card-section>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref } from 'vue'
import { useQuasar } from 'quasar'
import { useAssessment } from '../../composables/useAssessment.js'
import { useContext } from '../../composables/useContext.js'
import ReadOnlyContextBanner from '../../components/common/ReadOnlyContextBanner.vue'

const $q = useQuasar()

const { isReadOnlyContext } = useContext()

const {
  assessment,
  isLoading,
  error,
  isFinalized,
  loadAssessment,
  createAssessment,
  navigateToSession,
} = useAssessment()

const selectedClass = ref('')
const selectedSubject = ref('')
const assessmentDate = ref(new Date().toISOString().split('T')[0])

const classOptions = [
  { label: 'X-A', value: 'CLS-X-A' },
  { label: 'XI-B', value: 'CLS-XI-B' },
  { label: 'XII-C', value: 'CLS-XII-C' },
]

const subjectOptions = [
  { label: 'Matematika', value: 'SUB-MATH' },
  { label: 'Fisika', value: 'SUB-PHYS' },
  { label: 'Kimia', value: 'SUB-CHEM' },
]

const defaultComponents = [
  { id: 'COMP-001', name: 'Tugas', type: 'assignment', weight: 30 },
  { id: 'COMP-002', name: 'UTS', type: 'midterm', weight: 30 },
  { id: 'COMP-003', name: 'UAS', type: 'final', weight: 40 },
]

const handleLoadOrCreate = async () => {
  if (isReadOnlyContext.value) {
    $q.notify({
      type: 'warning',
      message: 'Data historis bersifat read-only',
    })

    return
  }

  if (!selectedClass.value || !selectedSubject.value) {
    $q.notify({
      type: 'warning',
      message: 'Pilih kelas dan mata pelajaran terlebih dahulu',
    })

    return
  }

  try {
    await loadAssessment({
      classId: selectedClass.value,
      subjectId: selectedSubject.value,
    })

    if (!assessment.value) {
      const classLabel = classOptions.find((item) => item.value === selectedClass.value)?.label || ''
      const subjectLabel = subjectOptions.find((item) => item.value === selectedSubject.value)?.label || ''

      await createAssessment({
        classId: selectedClass.value,
        className: classLabel,
        subjectId: selectedSubject.value,
        subjectName: subjectLabel,
        date: assessmentDate.value,
        components: defaultComponents,
      })

      $q.notify({
        type: 'positive',
        message: 'Penilaian baru berhasil dibuat',
      })
    }
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: err.message,
    })
  }
}

const handleNavigate = () => {
  if (!assessment.value) return

  navigateToSession(assessment.value.id)
}
</script>

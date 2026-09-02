<!--
FILE: src/pages/progress/ProgressPage.vue
STATUS: MODIFY
STATUS IMPLEMENTASI: COMPLETE
-->

<template>
  <q-page padding>
    <ReadOnlyContextBanner />

    <div class="text-h5 q-mb-md">Student Progress</div>

    <q-banner v-if="error" class="bg-negative text-white q-mb-md">
      {{ error }}
    </q-banner>

    <q-card class="q-mb-md">
      <q-card-section>
        <q-select v-model="selectedClass" :options="classOptions" label="Pilih Kelas" outlined dense emit-value
          map-options @update:model-value="handleClassChange" />
      </q-card-section>
    </q-card>

    <q-card v-if="isLoading" class="q-mb-md">
      <q-card-section class="text-center">
        <q-spinner color="primary" size="2em" />

        <div class="q-mt-sm">Loading...</div>
      </q-card-section>
    </q-card>

    <template v-else>
      <q-card v-if="students.length === 0">
        <q-card-section class="text-center text-grey">
          Pilih kelas untuk melihat daftar siswa.
        </q-card-section>
      </q-card>

      <q-list v-else separator>
        <q-item v-for="student in students" :key="student.studentId" clickable v-ripple
          @click="handleNavigate(student)">
          <q-item-section avatar>
            <q-icon name="person" color="primary" size="2em" />
          </q-item-section>

          <q-item-section>
            <q-item-label>{{ student.studentName }}</q-item-label>
            <q-item-label caption>{{ student.className }}</q-item-label>
          </q-item-section>

          <q-item-section side>
            <q-icon name="chevron_right" />
          </q-item-section>
        </q-item>
      </q-list>
    </template>
  </q-page>
</template>

<script setup>
import { ref } from 'vue'
import { useProgress } from '../../composables/useProgress.js'
import ReadOnlyContextBanner from '../../components/common/ReadOnlyContextBanner.vue'

const {
  students,
  isLoading,
  error,
  loadStudents,
  navigateToStudent,
} = useProgress()

const selectedClass = ref('')

const classOptions = [
  { label: 'X-A', value: 'CLS-X-A' },
  { label: 'XI-B', value: 'CLS-XI-B' },
  { label: 'XII-C', value: 'CLS-XII-C' },
]

const handleClassChange = async (classId) => {
  await loadStudents({ classId })
}

const handleNavigate = (student) => {
  navigateToStudent(student)
}
</script>

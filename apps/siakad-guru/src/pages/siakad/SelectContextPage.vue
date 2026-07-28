<template>
  <q-page class="q-pa-md flex flex-center">
    <q-card style="min-width: 400px; max-width: 90vw;">
      <q-card-section class="text-h6">Pilih Konteks Pembelajaran</q-card-section>
      <q-card-section>
        <!-- {{ ctx.academicYears }} -->
        <q-select v-model="selectedYear" :options="ctx.academicYears.value" label="Tahun Ajaran" option-value="id"
          option-label="name" @update:model-value="onYearChange" />
        <q-select v-model="selectedSemester" :options="ctx.semesters.value" label="Semester" option-value="id"
          option-label="name" class="q-mt-md" :disable="!selectedYear" @update:model-value="onSemesterChange" />
        <q-select v-model="selectedPeriod" :options="ctx.academicPeriods.value" label="Periode (opsional)"
          option-value="id" option-label="name" class="q-mt-md" :disable="!selectedSemester" />
        <q-btn label="Simpan & Lanjutkan" color="primary" class="full-width q-mt-lg"
          :disable="!selectedYear || !selectedSemester" @click="saveContext" />
      </q-card-section>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useContext } from '@/composables/useContext'
import { useRouter, useRoute } from 'vue-router'

const ctx = useContext()
const router = useRouter()
const route = useRoute()

const selectedYear = ref(null)
const selectedSemester = ref(null)
const selectedPeriod = ref(null)

onMounted(async () => {
  await ctx.loadAcademicYears()
})

function onYearChange(yearId) {
  selectedSemester.value = null
  selectedPeriod.value = null
  ctx.loadSemesters(yearId)
}
function onSemesterChange(semesterId) {
  selectedPeriod.value = null
  ctx.loadAcademicPeriods(semesterId)
}

function saveContext() {
  const year = ctx.academicYears.value.find(y => y.id === selectedYear.value.id)
  const semester = ctx.semesters.value.find(s => s.id === selectedSemester.value.id)
  const period = ctx.academicPeriods.value.find(p => p.id === selectedPeriod.value.id) || null
  ctx.setOperationalContext(year, semester, period)
  const redirect = route.query.redirect || '/siakad/jadwal-pelajaran'
  router.push(redirect)
}
</script>

<template>
  <div class="login-wrapper flex flex-center">
    <!-- <q-page class="flex flex-center bg-primary"> -->
    <q-card style="width: 500px; max-width: 90vw">
      <q-card-section class="text-center q-pb-none">
        <div class="text-h5">Pilih Tahun Ajaran</div>
        <div class="text-subtitle2 text-grey">Tentukan Sekolah, Tahun Pelajaran, dan Semester</div>
      </q-card-section>
      <q-card-section>
        <q-banner v-if="error" class="bg-negative text-white q-mb-md">
          <template #avatar>
            <q-icon name="error" />
          </template>
          {{ error }}
        </q-banner>
        <q-form @submit="handleSelectContext" class="q-gutter-md">
          <q-select v-model="selectedSchool" :options="schoolOptions" label="Sekolah" outlined dense emit-value
            map-options :disable="isLoading" />
          <q-select v-model="selectedAcademicYear" :options="academicYearOptions" label="Tahun Pelajaran" outlined dense
            emit-value map-options :disable="isLoading" />
          <q-select v-model="selectedSemester" :options="semesterOptions" label="Semester" outlined dense emit-value
            map-options :disable="isLoading || !selectedAcademicYear" />
          <q-btn type="submit" label="Lanjutkan" color="primary" :loading="isLoading" :disable="!isFormValid" unelevated
            class="full-width" />
        </q-form>
      </q-card-section>
      <q-card-section class="text-center text-caption text-grey">
        <p class="q-mb-none">Anda login sebagai: {{ currentUser?.name }}</p>
        <p>Role: {{ currentUser?.role }}</p>
      </q-card-section>
    </q-card>
    <!-- </q-page> -->
  </div>
</template>
<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useContext } from '../../composables/useContext.js'
import { useAuth } from '../../composables/useAuth.js'
import { useQuasar } from 'quasar'

const router = useRouter()
const $q = useQuasar()

const {
  availableContexts,
  isLoading,
  error,
  loadAvailableContexts,
  selectContext
} = useContext()

const { currentUser } = useAuth()

const selectedSchool = ref('')
const selectedAcademicYear = ref('')
const selectedSemester = ref('')
const isSelecting = ref(false)

const schoolOptions = computed(() => {
  if (!availableContexts.value?.school) return []
  return [{
    label: availableContexts.value.school.name,
    value: availableContexts.value.school.id
  }]
})

const academicYearOptions = computed(() => {
  if (!availableContexts.value?.academicYears) return []
  return availableContexts.value.academicYears.map(ay => ({
    label: ay.name,
    value: ay.id
  }))
})

const semesterOptions = computed(() => {
  if (!availableContexts.value?.academicYears) return []
  const selectedAY = availableContexts.value.academicYears.find(
    ay => ay.id === selectedAcademicYear.value
  )
  if (!selectedAY) return []
  return selectedAY.semesters.map(sem => ({
    label: sem.name,
    value: sem.id
  }))
})

const isFormValid = computed(() => {
  return !!(selectedSchool.value && selectedAcademicYear.value && selectedSemester.value)
})

// Load contexts saat component mounted
onMounted(async () => {
  try {
    await loadAvailableContexts()
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: 'Gagal memuat data konteks: ' + err.message
    })
  }
})

// Auto-select school pertama (jika hanya ada satu)
watch(availableContexts, (newVal) => {
  if (newVal?.school) {
    selectedSchool.value = newVal.school.id
  }
}, { immediate: true })

// Reset semester saat academic year berubah
watch(selectedAcademicYear, () => {
  selectedSemester.value = ''
})

const handleSelectContext = async () => {
  if (!isFormValid.value) {
    $q.notify({
      type: 'warning',
      message: 'Pilih semua konteks terlebih dahulu'
    })
    return
  }

  isSelecting.value = true
  try {
    await selectContext({
      schoolId: selectedSchool.value,
      academicYearId: selectedAcademicYear.value,
      semesterId: selectedSemester.value
    })

    router.push({ name: 'dashboard' })
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: err.message
    })
  } finally {
    isSelecting.value = false
  }
}
</script>
<style scoped>
.login-wrapper {
  min-height: 100vh;
}
</style>

<!--
FILE: src/pages/progress/ProgressStudentPage.vue
STATUS: MODIFY
STATUS IMPLEMENTASI: COMPLETE
-->

<template>
  <q-page padding>
    <div class="row items-center q-mb-md">
      <q-btn flat dense icon="arrow_back" label="Kembali" :to="{ name: 'progress' }" class="q-mr-sm" />

      <div class="text-h5">Progress Siswa</div>
    </div>

    <ReadOnlyContextBanner />

    <q-banner v-if="error" class="bg-negative text-white q-mb-md">
      {{ error }}
    </q-banner>

    <q-card class="q-mb-md">
      <q-card-section>
        <div class="text-h6">{{ studentName }}</div>
        <div class="text-caption text-grey">Kelas: {{ className }}</div>
      </q-card-section>
    </q-card>

    <q-card v-if="isLoading" class="q-mb-md">
      <q-card-section class="text-center">
        <q-spinner color="primary" size="2em" />

        <div class="q-mt-sm">Loading...</div>
      </q-card-section>
    </q-card>

    <template v-else>
      <div class="row items-center justify-between q-mb-sm">
        <div class="text-h6">Riwayat Progress</div>

        <q-btn v-if="!isReadOnlyContext" color="primary" label="Tambah Progress" size="sm"
          @click="showAddDialog = true" />

        <q-badge v-else color="grey">
          Read-only
        </q-badge>
      </div>

      <q-card v-if="progressRecords.length === 0" class="q-mb-md">
        <q-card-section class="text-center text-grey">
          Belum ada progress untuk siswa ini.
        </q-card-section>
      </q-card>

      <q-list v-else separator class="q-mb-md">
        <q-item v-for="record in progressRecords" :key="record.id">
          <q-item-section avatar>
            <q-icon :name="getTypeIcon(record.type)" :color="getTypeColor(record.type)" size="2em" />
          </q-item-section>

          <q-item-section>
            <q-item-label>{{ record.title }}</q-item-label>
            <q-item-label caption>{{ record.description }}</q-item-label>
            <q-item-label caption class="text-grey">
              {{ record.date }} | {{ getTypeLabel(record.type) }}
            </q-item-label>
          </q-item-section>
        </q-item>
      </q-list>
    </template>

    <q-dialog v-model="showAddDialog">
      <q-card style="min-width: 500px">
        <q-card-section>
          <div class="text-h6">Tambah Progress</div>
        </q-card-section>

        <q-card-section class="q-gutter-md">
          <q-select v-model="formType" :options="typeOptions" label="Jenis Progress" outlined dense emit-value
            map-options />

          <q-input v-model="formTitle" label="Judul" outlined dense />

          <q-input v-model="formDescription" label="Deskripsi" type="textarea" outlined rows="3" />

          <q-input v-model="formDate" label="Tanggal" type="date" outlined dense />
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="Batal" color="primary" v-close-popup />

          <q-btn flat label="Simpan" color="positive" :loading="isLoading" @click="handleAddProgress" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { useQuasar } from 'quasar'
import { useProgress } from '../../composables/useProgress.js'
import { useContext } from '../../composables/useContext.js'
import { getProgressTypeLabel } from '../../domain/progress/models/ProgressType.js'
import ReadOnlyContextBanner from '../../components/common/ReadOnlyContextBanner.vue'

const route = useRoute()
const $q = useQuasar()

const { isReadOnlyContext } = useContext()

const {
  progressRecords,
  isLoading,
  error,
  loadProgress,
  createProgress,
} = useProgress()

const studentId = computed(() => route.params.studentId)
const studentName = computed(() => route.query.studentName || '')
const classId = computed(() => route.query.classId || '')
const className = computed(() => route.query.className || '')

const showAddDialog = ref(false)
const formType = ref('')
const formTitle = ref('')
const formDescription = ref('')
const formDate = ref(new Date().toISOString().split('T')[0])

const typeOptions = [
  { label: 'Catatan Guru', value: 'note' },
  { label: 'Prestasi', value: 'achievement' },
  { label: 'Pelanggaran', value: 'violation' },
  { label: 'Bimbingan Konseling', value: 'counseling' },
]

onMounted(async () => {
  await loadProgress({ studentId: studentId.value })
})

const handleAddProgress = async () => {
  if (isReadOnlyContext.value) {
    $q.notify({
      type: 'warning',
      message: 'Data historis bersifat read-only',
    })

    return
  }

  if (!formType.value || !formTitle.value || !formDescription.value) {
    $q.notify({
      type: 'warning',
      message: 'Lengkapi semua field terlebih dahulu',
    })

    return
  }

  try {
    await createProgress({
      studentId: studentId.value,
      studentName: studentName.value,
      classId: classId.value,
      className: className.value,
      type: formType.value,
      title: formTitle.value,
      description: formDescription.value,
      date: formDate.value,
    })

    showAddDialog.value = false
    formType.value = ''
    formTitle.value = ''
    formDescription.value = ''

    $q.notify({
      type: 'positive',
      message: 'Progress berhasil ditambahkan',
    })
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: err.message,
    })
  }
}

const getTypeIcon = (type) => {
  const iconMap = {
    note: 'edit_note',
    achievement: 'emoji_events',
    violation: 'warning',
    counseling: 'psychology',
  }

  return iconMap[type] || 'info'
}

const getTypeColor = (type) => {
  const colorMap = {
    note: 'blue',
    achievement: 'positive',
    violation: 'negative',
    counseling: 'purple',
  }

  return colorMap[type] || 'grey'
}

const getTypeLabel = (type) => {
  return getProgressTypeLabel(type)
}
</script>

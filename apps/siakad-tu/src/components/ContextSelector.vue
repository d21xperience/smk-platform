<template>
  <div>
    <q-btn flat dense no-caps class="context-selector-btn" @click="showPopup = true">
      <div class="row items-center no-wrap q-gutter-xs">
        <q-icon name="event" size="20px" />
        <div class="text-left">
          <div class="text-caption text-weight-medium" style="line-height: 1.2">
            {{ displayLabel || 'Pilih Context' }}
          </div>
        </div>
        <q-icon name="arrow_drop_down" size="18px" />
      </div>
    </q-btn>

    <q-dialog v-model="showPopup" persistent>
      <q-card style="min-width: 350px">
        <q-card-section class="row items-center q-pb-none">
          <div class="text-h6">Pilih Tahun Ajaran</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pt-md">
          <div class="q-gutter-md">
            <q-select v-model="localYear" :options="yearOptions" label="Tahun Ajaran" outlined emit-value map-options
              :disable="!hasYearChangeAuthority"
              :hint="!hasYearChangeAuthority ? 'Anda tidak memiliki akses mengubah tahun ajaran' : ''">
              <template v-slot:prepend>
                <q-icon name="calendar_today" />
              </template>
            </q-select>

            <q-select v-model="localSemester" :options="semesterOptions" label="Semester" outlined emit-value
              map-options>
              <template v-slot:prepend>
                <q-icon name="school" />
              </template>
            </q-select>
          </div>
        </q-card-section>

        <q-card-actions align="right" class="q-px-md q-pb-md">
          <q-btn flat label="Batal" color="grey" v-close-popup />
          <q-btn label="Simpan" color="primary" @click="applySelection" :loading="isLoading"
            :disable="!localYear || !localSemester" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useQuasar } from 'quasar'
import { useOperationalContext } from '@/composables/context/useOperationalContext'

const $q = useQuasar()

const {
  displayLabel,
  selectedYear,
  selectedSemester,
  yearOptions,
  semesterOptions,
  hasYearChangeAuthority,
  isLoading,
  error,
  onYearChange,
  onSemesterChange,
} = useOperationalContext()

const showPopup = ref(false)

// Temporary state for the form
const localYear = ref(selectedYear.value)
const localSemester = ref(selectedSemester.value)

// Sync temp state when dialog opens
watch(showPopup, (isOpen) => {
  if (isOpen) {
    localYear.value = selectedYear.value
    localSemester.value = selectedSemester.value
  }
})

async function applySelection() {
  if (localYear.value !== selectedYear.value && hasYearChangeAuthority.value) {
    await onYearChange(localYear.value)
  }
  if (localSemester.value !== selectedSemester.value) {
    await onSemesterChange(localSemester.value)
  }

  if (!error.value) {
    $q.notify({ type: 'positive', message: 'Tahun Ajaran & Semester berhasil diperbarui.' })
    showPopup.value = false
  } else {
    $q.notify({ type: 'negative', message: error.value })
  }
}
</script>

<style scoped>
.context-selector-btn {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  transition: background 0.2s;
}

.context-selector-btn:hover {
  background: rgba(255, 255, 255, 0.2);
}
</style>

<template>
  <div>
    <!-- Compact Button -->
    <q-btn
      flat
      dense
      no-caps
      class="context-selector-btn"
      @click="showPopup = true"
    >
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

    <!-- Popup Dialog -->
    <q-dialog v-model="showPopup" persistent>
      <q-card style="min-width: 350px">
        <q-card-section class="row items-center q-pb-none">
          <div class="text-h6">Pilih Tahun Ajaran</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pt-md">
          <div class="q-gutter-md">
            <!-- Tahun Ajaran -->
            <q-select
              v-model="tempSelectedYear"
              :options="yearOptions"
              label="Tahun Ajaran"
              outlined
              emit-value
              map-options
            >
              <template v-slot:prepend>
                <q-icon name="calendar_today" />
              </template>
            </q-select>

            <!-- Semester -->
            <q-select
              v-model="tempSelectedSemester"
              :options="semesterOptions"
              label="Semester"
              outlined
              emit-value
              map-options
            >
              <template v-slot:prepend>
                <q-icon name="school" />
              </template>
            </q-select>
          </div>
        </q-card-section>

        <q-card-actions align="right" class="q-px-md q-pb-md">
          <q-btn flat label="Batal" color="grey" v-close-popup />
          <q-btn
            label="Simpan"
            color="primary"
            @click="applySelection"
            :disable="!tempSelectedYear || !tempSelectedSemester"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useOperationalContext } from '@/composables/context/useOperationalContext'

const {
  selectedYear,
  selectedSemester,
  yearOptions,
  semesterOptions,
  onYearChange,
  onSemesterChange,
  displayLabel,
} = useOperationalContext()

// === STATE ===
const showPopup = ref(false)

// Temporary state untuk popup (belum di-apply)
const tempSelectedYear = ref(selectedYear.value)
const tempSelectedSemester = ref(selectedSemester.value)

// === WATCHERS ===
// Sinkronisasi temp state dengan actual state saat popup dibuka
watch(showPopup, (isOpen) => {
  if (isOpen) {
    tempSelectedYear.value = selectedYear.value
    tempSelectedSemester.value = selectedSemester.value
  }
})

// === ACTIONS ===
function applySelection() {
  if (tempSelectedYear.value !== selectedYear.value) {
    onYearChange(tempSelectedYear.value)
  }
  if (tempSelectedSemester.value !== selectedSemester.value) {
    onSemesterChange(tempSelectedSemester.value)
  }
  showPopup.value = false
}
</script>

<style scoped>
.context-selector-btn {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  /* padding: 4px 12px; */
  transition: background 0.2s;
}

.context-selector-btn:hover {
  background: rgba(255, 255, 255, 0.2);
}
</style>

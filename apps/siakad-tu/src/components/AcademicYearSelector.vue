<template>
  <div>
    <!-- TOMBOL UTAMA HEADER (Menampilkan Tahun & Semester Aktif) -->
    <q-btn flat no-caps color="primary" icon="event"
      :label="`TA: ${academicStore.activeYearLabel} - ${academicStore.activeSemesterLabel}`"
      class="text-weight-bold font-tab-match text-white" @click="openDialog">

      <q-tooltip>Klik untuk mengubah Tahun Ajaran & Semester</q-tooltip>
    </q-btn>

    <!-- JENDELA DIALOG DUA PILIHAN TERPADU -->
    <q-dialog v-model="state.isDialogOpen" transition-show="scale" transition-hide="scale">
      <q-card style="width: 500px; max-width: 90vw;" class="q-pa-sm">
        <q-card-section class="row items-center q-pb-none">
          <div>
            <div class="text-subtitle1 text-weight-bold text-primary">Konfigurasi Periode Administrasi</div>
            <div class="text-caption text-grey-6">Target Value Backend: <q-badge color="grey-3" text-color="dark">{{
              calculatedPeriodCode }}</q-badge></div>
          </div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="row q-col-gutter-md q-pt-md">
          <!-- KOLOM KIRI: PILIHAN TAHUN AJARAN -->
          <div class="col-12 col-sm-7">
            <div class="text-caption text-weight-bold text-grey-7 q-mb-xs">1. Pilih Tahun Ajaran</div>
            <!-- <q-list bordered separator class="rounded-borders">
              <q-item v-for="year in props.availableYears" :key="year.value" tag="label" v-ripple
                :class="state.selectedStartYear === year.value ? 'bg-blue-1 text-weight-bold' : ''">
                <q-item-section avatar>
                  <q-radio v-model="state.selectedStartYear" :val="year.value" color="primary" />
                </q-item-section>
                <q-item-section>
                  <q-item-label>{{ year.label }}</q-item-label>
                </q-item-section>
              </q-item>
            </q-list> -->
            <q-select v-model="state.selectedStartYear" dense outlined :options="props.availableYears"
              :option-label="props.availableYears.label" :option-value="props.availableYears.value" />

          </div>

          <!-- KOLOM KANAN: PILIHAN SEMESTER -->
          <div class="col-12 col-sm-5">
            <div class="text-caption text-weight-bold text-grey-7 q-mb-xs">2. Pilih Semester</div>
            <q-list bordered separator class="rounded-borders">
              <q-item tag="label" v-ripple :class="state.selectedSemester === '1' ? 'bg-teal-1' : ''">
                <q-item-section avatar>
                  <q-radio v-model="state.selectedSemester" val="1" color="teal" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-medium">Ganjil (1)</q-item-label>
                </q-item-section>
              </q-item>

              <q-item tag="label" v-ripple :class="state.selectedSemester === '2' ? 'bg-teal-1' : ''">
                <q-item-section avatar>
                  <q-radio v-model="state.selectedSemester" val="2" color="teal" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-medium">Genap (2)</q-item-label>
                </q-item-section>
              </q-item>
            </q-list>
          </div>
        </q-card-section>

        <!-- TOMBOL AKSI EKSEKUSI -->
        <q-card-actions align="right" class="q-px-md q-pb-md q-pt-none">
          <q-btn flat label="Batal" color="grey-7" v-close-popup />
          <q-btn label="Simpan Sesi TU" color="primary" icon="check" @click="applyPeriodChange" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </div>
</template>

<script setup>
import { reactive, computed } from 'vue'
import { useAcademicStore } from '@/stores/academicStore'
import { useQuasar } from 'quasar'

const academicStore = useAcademicStore()
const $q = useQuasar()

const props = defineProps({
  availableYears: {
    type: Array,
    default: () => [
      { label: '2024/2025', value: 2024 },
      { label: '2025/2026', value: 2025 },
      { label: '2026/2027', value: 2026 },
      { label: '2027/2028', value: 2027 }
    ]
  }
})

const state = reactive({
  isDialogOpen: false,
  selectedStartYear: {
    value: 2026,
    label: '20226/2027'
  },
  selectedSemester: '1'
})

// Fungsi pembuat value gabungan otomatis (e.g., 2026 + "2" = "20262")
const calculatedPeriodCode = computed(() => {
  return `${state.selectedStartYear.value}${state.selectedSemester}`
})

// Sinkronisasi data form saat dialog dibuka oleh user
function openDialog() {
  const currentCode = academicStore.activePeriodCode
  state.selectedStartYear = parseInt(currentCode.substring(0, 4))
  state.selectedSemester = currentCode.slice(-1)
  state.isDialogOpen = true
}

function applyPeriodChange() {
  const finalCode = calculatedPeriodCode.value

  // Update data ke Pinia Store
  academicStore.setActivePeriod(finalCode)
  state.isDialogOpen = false

  $q.notify({
    type: 'positive',
    icon: 'swap_horiz',
    message: `Sesi beralih ke TA ${academicStore.activeYearLabel} Semester ${academicStore.activeSemesterLabel} (Kode: ${finalCode})`,
    position: 'top'
  })
}
</script>

<style scoped>
.font-tab-match {
  font-weight: 600;
}
</style>

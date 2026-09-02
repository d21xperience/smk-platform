<!--
FILE: src/components/common/HistoricalContextSelector.vue
STATUS: NEW
STATUS IMPLEMENTASI: COMPLETE
-->

<template>
  <q-dialog :model-value="modelValue" persistent @update:model-value="handleVisibilityChange">
    <q-card style="width: 560px; max-width: 92vw">
      <q-card-section class="row items-center no-wrap">
        <div>
          <div class="text-h6">Riwayat Context</div>
          <div class="text-caption text-grey">
            Pilih tahun pelajaran dan semester untuk melihat data historis.
          </div>
        </div>

        <q-space />

        <q-btn flat dense round icon="close" aria-label="Tutup" @click="closeDialog" />
      </q-card-section>

      <q-separator />

      <q-card-section>
        <q-banner v-if="historyError" class="bg-negative text-white q-mb-md">
          {{ historyError }}
        </q-banner>

        <div v-if="currentContext" class="q-mb-md">
          <div class="text-caption text-grey-7 q-mb-xs">
            Context aktif saat ini
          </div>

          <q-badge color="positive">
            {{ activeContextLabel }}
          </q-badge>
        </div>

        <q-select v-model="selectedEntry" :options="entryOptions" option-label="label"
          label="Tahun Pelajaran & Semester Historis" outlined dense clearable :loading="isLoadingHistory"
          :disable="isLoadingHistory" />

        <q-banner v-if="!isLoadingHistory && entryOptions.length === 0" class="bg-grey-3 q-mt-md">
          Belum ada data historis yang tersedia.
        </q-banner>

        <div class="q-mt-md text-caption text-grey-7">
          Data historis bersifat read-only. Create, update, dan delete tidak tersedia pada mode ini.
        </div>
      </q-card-section>

      <q-card-actions align="right" class="q-px-md q-pb-md">
        <q-btn flat label="Batal" :disable="isLoadingHistory" @click="closeDialog" />

        <q-btn color="primary" label="Lihat Riwayat" :loading="isLoadingHistory" :disable="!selectedEntry"
          @click="handleSelectHistory" />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { useQuasar } from 'quasar'
import { useContext } from '../../composables/useContext.js'
import { formatContextLabel } from '../../utils/contextFormatter.js'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:modelValue', 'selected'])

const $q = useQuasar()

const {
  currentContext,
  availableHistory,
  isLoadingHistory,
  historyError,
  loadAvailableHistory,
  enterHistoricalContext,
} = useContext()

const selectedEntry = ref(null)

const entryOptions = computed(() => {
  return availableHistory.value?.entries || []
})

const activeContextLabel = computed(() => {
  return formatContextLabel(currentContext.value)
})

const handleVisibilityChange = (value) => {
  emit('update:modelValue', value)
}

const closeDialog = () => {
  emit('update:modelValue', false)
}

watch(
  () => props.modelValue,
  async (isOpen) => {
    if (!isOpen) return

    selectedEntry.value = null

    try {
      await loadAvailableHistory()
    } catch (err) {
      $q.notify({
        type: 'negative',
        message: `Gagal memuat riwayat context: ${err.message}`,
      })
    }
  }
)

const handleSelectHistory = async () => {
  if (!selectedEntry.value) {
    $q.notify({
      type: 'warning',
      message: 'Pilih tahun pelajaran dan semester historis terlebih dahulu',
    })
    return
  }

  try {
    const historicalContext = await enterHistoricalContext({
      academicYearId: selectedEntry.value.academicYearId,
      semesterId: selectedEntry.value.semesterId,
    })

    emit('selected', historicalContext)
    closeDialog()
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: `Gagal masuk ke context historis: ${err.message}`,
    })
  }
}
</script>

<!--
FILE: src/pages/assessment/AssessmentSessionPage.vue
STATUS: MODIFY
STATUS IMPLEMENTASI: COMPLETE
-->

<template>
  <q-page padding>
    <div class="row items-center q-mb-md">
      <q-btn flat dense icon="arrow_back" label="Kembali" :to="{ name: 'assessment' }" class="q-mr-sm" />

      <div class="text-h5">Kelola Penilaian</div>
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

    <template v-else-if="assessment">
      <q-card class="q-mb-md">
        <q-card-section>
          <div class="text-h6">
            {{ assessment.subjectName }} - {{ assessment.className }}
          </div>

          <div class="text-caption text-grey">
            {{ assessment.date }} | Status:

            <q-badge :color="isFinalized ? 'blue' : 'orange'">
              {{ assessment.status }}
            </q-badge>
          </div>
        </q-card-section>
      </q-card>

      <q-card class="q-mb-md">
        <q-card-section>
          <div class="text-h6 q-mb-md">Komponen Penilaian</div>

          <q-list separator>
            <q-item v-for="component in assessment.components" :key="component.id">
              <q-item-section>
                <q-item-label>{{ component.name }}</q-item-label>
                <q-item-label caption>Bobot: {{ component.weight }}%</q-item-label>
              </q-item-section>

              <q-item-section side>
                <q-input :model-value="assessment.scores[component.id]" type="number" dense outlined
                  style="width: 100px" :min="0" :max="100" :disable="isFinalized || isReadOnlyContext"
                  @update:model-value="handleScoreChange(component.id, $event)" />
              </q-item-section>
            </q-item>
          </q-list>
        </q-card-section>
      </q-card>

      <q-card v-if="previewResult" class="q-mb-md">
        <q-card-section>
          <div class="text-h6 q-mb-md">Hasil Perhitungan</div>

          <div class="row q-gutter-md text-center">
            <div class="col">
              <div class="text-h5">{{ previewResult.finalScore }}</div>
              <div class="text-caption">Nilai Akhir</div>
            </div>

            <div class="col">
              <div class="text-h5">{{ previewResult.grade }}</div>
              <div class="text-caption">Grade</div>
            </div>

            <div class="col">
              <div class="text-h6">{{ previewResult.predicate }}</div>
              <div class="text-caption">Predikat</div>
            </div>
          </div>
        </q-card-section>
      </q-card>

      <q-card v-if="result" class="q-mb-md">
        <q-card-section>
          <div class="text-h6 q-mb-md">Hasil Final</div>

          <div class="row q-gutter-md text-center">
            <div class="col">
              <div class="text-h5">{{ result.finalScore }}</div>
              <div class="text-caption">Nilai Akhir</div>
            </div>

            <div class="col">
              <div class="text-h5">{{ result.grade }}</div>
              <div class="text-caption">Grade</div>
            </div>

            <div class="col">
              <div class="text-h6">{{ result.predicate }}</div>
              <div class="text-caption">Predikat</div>
            </div>
          </div>
        </q-card-section>
      </q-card>

      <q-card v-if="isReadOnlyContext">
        <q-card-section class="text-center text-grey">
          Penilaian historis hanya dapat dilihat. Perubahan nilai dan finalisasi tidak tersedia.
        </q-card-section>
      </q-card>

      <q-card v-else-if="isDraft">
        <q-card-section class="row q-gutter-sm">
          <q-btn color="positive" label="Finalisasi Penilaian" :loading="isLoading" @click="handleFinalize" />
        </q-card-section>
      </q-card>

      <q-card v-else-if="isFinalized">
        <q-card-section class="text-center text-grey">
          Penilaian telah difinalisasi dan tidak dapat diubah.
        </q-card-section>
      </q-card>
    </template>

    <q-card v-else>
      <q-card-section class="text-center text-grey">
        Penilaian tidak ditemukan.
      </q-card-section>
    </q-card>
  </q-page>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useQuasar } from 'quasar'
import { useAssessment } from '../../composables/useAssessment.js'
import { useContext } from '../../composables/useContext.js'
import ReadOnlyContextBanner from '../../components/common/ReadOnlyContextBanner.vue'

const route = useRoute()
const router = useRouter()
const $q = useQuasar()

const { isReadOnlyContext } = useContext()

const {
  assessment,
  result,
  isLoading,
  error,
  isDraft,
  isFinalized,
  previewResult,
  loadAssessmentById,
  updateScore,
  finalizeAssessment,
} = useAssessment()

const assessmentId = computed(() => route.params.assessmentId)

onMounted(async () => {
  try {
    await loadAssessmentById(assessmentId.value)
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: err.message,
    })

    router.push({ name: 'assessment' })
  }
})

const handleScoreChange = async (componentId, score) => {
  if (isReadOnlyContext.value) return

  const numericScore = Number(score)

  if (Number.isNaN(numericScore)) return

  try {
    await updateScore({
      componentId,
      score: numericScore,
    })
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: err.message,
    })
  }
}

const handleFinalize = async () => {
  if (isReadOnlyContext.value) {
    $q.notify({
      type: 'warning',
      message: 'Data historis bersifat read-only',
    })

    return
  }

  try {
    await finalizeAssessment()

    $q.notify({
      type: 'positive',
      message: 'Penilaian berhasil difinalisasi',
    })
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: err.message,
    })
  }
}
</script>

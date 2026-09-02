<!--
FILE: src/pages/journal/JournalSessionPage.vue
STATUS: MODIFY
STATUS IMPLEMENTASI: COMPLETE
-->

<template>
  <q-page padding>
    <div class="row items-center q-mb-md">
      <q-btn flat dense icon="arrow_back" label="Kembali" :to="{ name: 'journal' }" class="q-mr-sm" />

      <div class="text-h5">Jurnal Mengajar</div>
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

    <template v-else-if="journal">
      <q-card class="q-mb-md">
        <q-card-section>
          <div class="text-h6">
            {{ journal.subjectName }} - {{ journal.className }}
          </div>

          <div class="text-caption text-grey">
            {{ journal.date }} | Status:

            <q-badge :color="isSubmitted ? 'blue' : 'orange'">
              {{ journal.status }}
            </q-badge>
          </div>
        </q-card-section>
      </q-card>

      <q-card class="q-mb-md">
        <q-card-section>
          <div class="text-h6 q-mb-md">Isi Jurnal</div>

          <q-input v-model="material" label="Materi Pembelajaran" type="textarea" outlined rows="3" class="q-mb-md"
            :disable="isSubmitted || isReadOnlyContext" />

          <q-input v-model="activities" label="Kegiatan Pembelajaran" type="textarea" outlined rows="4" class="q-mb-md"
            :disable="isSubmitted || isReadOnlyContext" />

          <q-input v-model="reflection" label="Refleksi" type="textarea" outlined rows="3"
            :disable="isSubmitted || isReadOnlyContext" />
        </q-card-section>
      </q-card>

      <q-card v-if="isReadOnlyContext">
        <q-card-section class="text-center text-grey">
          Jurnal historis hanya dapat dilihat. Simpan draft dan submit tidak tersedia.
        </q-card-section>
      </q-card>

      <q-card v-else-if="isDraft">
        <q-card-section class="row q-gutter-sm">
          <q-btn color="primary" label="Simpan Draft" :loading="isSaving" @click="handleSave" />

          <q-btn color="positive" label="Submit Jurnal" :loading="isLoading" @click="handleSubmit" />
        </q-card-section>
      </q-card>

      <q-card v-else-if="isSubmitted">
        <q-card-section class="text-center text-grey">
          Jurnal telah disubmit dan tidak dapat diubah.
        </q-card-section>
      </q-card>
    </template>

    <q-card v-else-if="isReadOnlyContext">
      <q-card-section class="text-center text-grey">
        Data jurnal untuk periode historis ini tidak tersedia.
      </q-card-section>
    </q-card>

    <q-card v-else>
      <q-card-section class="text-center text-grey">
        Jurnal tidak ditemukan.
      </q-card-section>
    </q-card>
  </q-page>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useQuasar } from 'quasar'
import { useJournal } from '../../composables/useJournal.js'
import { useContext } from '../../composables/useContext.js'
import ReadOnlyContextBanner from '../../components/common/ReadOnlyContextBanner.vue'

const route = useRoute()
const router = useRouter()
const $q = useQuasar()

const { isReadOnlyContext } = useContext()

const {
  journal,
  isLoading,
  isSaving,
  error,
  isDraft,
  isSubmitted,
  loadOrCreateJournal,
  updateJournal,
  submitJournal,
} = useJournal()

const teachingSessionId = computed(() => route.params.teachingSessionId)

const material = ref('')
const activities = ref('')
const reflection = ref('')

onMounted(async () => {
  try {
    await loadOrCreateJournal(teachingSessionId.value)

    if (journal.value) {
      material.value = journal.value.material
      activities.value = journal.value.activities
      reflection.value = journal.value.reflection
    }
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: err.message,
    })

    router.push({ name: 'journal' })
  }
})

const handleSave = async () => {
  if (isReadOnlyContext.value) {
    $q.notify({
      type: 'warning',
      message: 'Data historis bersifat read-only',
    })

    return
  }

  try {
    await updateJournal({
      material: material.value,
      activities: activities.value,
      reflection: reflection.value,
    })

    $q.notify({
      type: 'positive',
      message: 'Jurnal berhasil disimpan',
    })
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: err.message,
    })
  }
}

const handleSubmit = async () => {
  if (isReadOnlyContext.value) {
    $q.notify({
      type: 'warning',
      message: 'Data historis bersifat read-only',
    })

    return
  }

  try {
    await submitJournal()

    $q.notify({
      type: 'positive',
      message: 'Jurnal berhasil disubmit',
    })

    router.push({ name: 'journal' })
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: err.message,
    })
  }
}
</script>

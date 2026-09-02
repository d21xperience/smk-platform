<!--
FILE: src/pages/teaching/TeachingSessionPage.vue
STATUS: MODIFY
STATUS IMPLEMENTASI: COMPLETE
-->

<template>
  <q-page padding>
    <div class="row items-center q-mb-md">
      <q-btn flat dense icon="arrow_back" label="Kembali" :to="{ name: 'teaching' }" class="q-mr-sm" />

      <div class="text-h5">Detail Sesi Mengajar</div>
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

    <template v-else-if="session">
      <q-card class="q-mb-md">
        <q-card-section>
          <div class="text-h6">
            {{ session.subjectName }} - {{ session.className }}
          </div>

          <q-separator class="q-my-md" />

          <q-list dense>
            <q-item>
              <q-item-section>
                <q-item-label>Tanggal</q-item-label>
              </q-item-section>

              <q-item-section side>
                <q-item-label caption>
                  {{ session.date }}
                </q-item-label>
              </q-item-section>
            </q-item>

            <q-item>
              <q-item-section>
                <q-item-label>Waktu</q-item-label>
              </q-item-section>

              <q-item-section side>
                <q-item-label caption>
                  {{ session.startTime }} - {{ session.endTime }}
                </q-item-label>
              </q-item-section>
            </q-item>

            <q-item>
              <q-item-section>
                <q-item-label>Status</q-item-label>
              </q-item-section>

              <q-item-section side>
                <q-badge :color="getStatusColor(session.status)">
                  {{ session.status }}
                </q-badge>
              </q-item-section>
            </q-item>

            <q-item v-if="session.teacherPresence">
              <q-item-section>
                <q-item-label>Kehadiran Guru</q-item-label>
              </q-item-section>

              <q-item-section side>
                <q-item-label caption>
                  {{ session.teacherPresence }}
                </q-item-label>
              </q-item-section>
            </q-item>

            <q-item v-if="session.notes">
              <q-item-section>
                <q-item-label>Catatan</q-item-label>
              </q-item-section>

              <q-item-section side>
                <q-item-label caption>
                  {{ session.notes }}
                </q-item-label>
              </q-item-section>
            </q-item>
          </q-list>
        </q-card-section>
      </q-card>

      <template v-if="!isReadOnlyContext">
        <q-card v-if="session.status === 'scheduled'">
          <q-card-section>
            <div class="text-h6 q-mb-md">Mulai Sesi</div>

            <q-select v-model="teacherPresence" :options="presenceOptions" label="Kehadiran Guru" outlined dense
              emit-value map-options class="q-mb-md" />

            <div class="row q-gutter-sm">
              <q-btn color="positive" label="Mulai Sesi" :disable="!teacherPresence || isLoading" :loading="isLoading"
                @click="handleStartSession" />

              <q-btn color="negative" label="Batalkan" outline @click="showCancelDialog = true" />
            </div>
          </q-card-section>
        </q-card>

        <q-card v-else-if="session.status === 'active'">
          <q-card-section>
            <div class="text-h6 q-mb-md">Akhiri Sesi</div>

            <q-input v-model="endNotes" label="Catatan Mengajar" type="textarea" outlined rows="3" class="q-mb-md" />

            <div class="row q-gutter-sm">
              <q-btn color="primary" label="Akhiri Sesi" :loading="isLoading" @click="handleEndSession" />

              <q-btn color="negative" label="Batalkan" outline @click="showCancelDialog = true" />
            </div>
          </q-card-section>
        </q-card>

        <q-card v-else-if="session.status === 'completed'">
          <q-card-section class="text-center text-grey">
            Sesi ini telah selesai.
          </q-card-section>
        </q-card>

        <q-card v-else-if="session.status === 'cancelled'">
          <q-card-section class="text-center text-grey">
            Sesi ini telah dibatalkan.
          </q-card-section>
        </q-card>
      </template>

      <q-card v-else>
        <q-card-section class="text-center text-grey">
          Data historis hanya dapat dilihat. Aksi mulai, akhiri, dan batalkan sesi tidak tersedia.
        </q-card-section>
      </q-card>
    </template>

    <q-dialog v-model="showCancelDialog">
      <q-card style="min-width: 400px">
        <q-card-section>
          <div class="text-h6">Batalkan Sesi</div>
        </q-card-section>

        <q-card-section>
          <q-input v-model="cancelReason" label="Alasan Pembatalan" type="textarea" outlined rows="3" />
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="Batal" color="primary" v-close-popup />

          <q-btn flat label="Konfirmasi" color="negative" :loading="isLoading" @click="handleCancelSession" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useQuasar } from 'quasar'
import { useTeaching } from '../../composables/useTeaching.js'
import { useContext } from '../../composables/useContext.js'
import ReadOnlyContextBanner from '../../components/common/ReadOnlyContextBanner.vue'

const route = useRoute()
const router = useRouter()
const $q = useQuasar()

const { isReadOnlyContext } = useContext()

const {
  sessions,
  isLoading,
  error,
  loadSessions,
  startSession,
  endSession,
  cancelSession,
} = useTeaching()

const sessionId = computed(() => route.params.sessionId)

const session = computed(() => {
  return sessions.value.find((s) => s.id === sessionId.value)
})

const teacherPresence = ref('')
const endNotes = ref('')
const cancelReason = ref('')
const showCancelDialog = ref(false)

const presenceOptions = [
  { label: 'Hadir', value: 'present' },
  { label: 'Terlambat', value: 'late' },
  { label: 'Tidak Hadir', value: 'absent' },
  { label: 'Delegasi', value: 'delegated' },
]

onMounted(async () => {
  if (sessions.value.length === 0) {
    const today = new Date().toISOString().split('T')[0]

    await loadSessions(today)
  }
})

const handleStartSession = async () => {
  if (isReadOnlyContext.value) {
    $q.notify({
      type: 'warning',
      message: 'Data historis bersifat read-only',
    })

    return
  }

  try {
    await startSession(sessionId.value, teacherPresence.value)

    $q.notify({
      type: 'positive',
      message: 'Sesi mengajar dimulai',
    })
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: err.message,
    })
  }
}

const handleEndSession = async () => {
  if (isReadOnlyContext.value) {
    $q.notify({
      type: 'warning',
      message: 'Data historis bersifat read-only',
    })

    return
  }

  try {
    await endSession(sessionId.value, endNotes.value)

    $q.notify({
      type: 'positive',
      message: 'Sesi mengajar selesai',
    })

    router.push({ name: 'teaching' })
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: err.message,
    })
  }
}

const handleCancelSession = async () => {
  if (isReadOnlyContext.value) {
    $q.notify({
      type: 'warning',
      message: 'Data historis bersifat read-only',
    })

    return
  }

  try {
    await cancelSession(sessionId.value, cancelReason.value)

    showCancelDialog.value = false

    $q.notify({
      type: 'info',
      message: 'Sesi mengajar dibatalkan',
    })

    router.push({ name: 'teaching' })
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: err.message,
    })
  }
}

const getStatusColor = (status) => {
  const colorMap = {
    scheduled: 'grey',
    active: 'green',
    completed: 'blue',
    cancelled: 'red',
  }

  return colorMap[status] || 'grey'
}
</script>

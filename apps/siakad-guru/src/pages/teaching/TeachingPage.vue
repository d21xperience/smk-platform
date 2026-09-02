<!--
FILE: src/pages/teaching/TeachingPage.vue
STATUS: MODIFY
STATUS IMPLEMENTASI: COMPLETE
-->

<template>
  <q-page padding>
    <ReadOnlyContextBanner />

    <div class="row items-center justify-between q-mb-md">
      <div class="text-h5">Teaching</div>

      <q-input v-model="selectedDate" type="date" outlined dense style="width: 200px"
        @update:model-value="handleDateChange" />
    </div>

    <q-banner v-if="error" class="bg-negative text-white q-mb-md">
      {{ error }}
    </q-banner>

    <q-card v-if="isLoading" class="q-mb-md">
      <q-card-section class="text-center">
        <q-spinner color="primary" size="2em" />

        <div class="q-mt-sm">Loading...</div>
      </q-card-section>
    </q-card>

    <template v-else>
      <div class="text-h6 q-mb-sm">Sesi Hari Ini</div>

      <q-card v-if="sessions.length === 0" class="q-mb-md">
        <q-card-section class="text-center text-grey">
          Tidak ada sesi mengajar untuk tanggal ini.
        </q-card-section>
      </q-card>

      <q-list v-else separator class="q-mb-md">
        <q-item v-for="session in sessions" :key="session.id" clickable v-ripple @click="navigateToSession(session.id)">
          <q-item-section avatar>
            <q-icon :name="getSessionIcon(session.status)" :color="getStatusColor(session.status)" size="2em" />
          </q-item-section>

          <q-item-section>
            <q-item-label>
              {{ session.subjectName }} - {{ session.className }}
            </q-item-label>

            <q-item-label caption>
              {{ session.startTime }} - {{ session.endTime }}
            </q-item-label>
          </q-item-section>

          <q-item-section side>
            <q-badge :color="getStatusColor(session.status)">
              {{ session.status }}
            </q-badge>
          </q-item-section>
        </q-item>
      </q-list>

      <div class="text-h6 q-mb-sm">Jadwal Tersedia</div>

      <q-card v-if="getAvailableSchedules.length === 0">
        <q-card-section class="text-center text-grey">
          Tidak ada jadwal yang tersedia untuk tanggal ini.
        </q-card-section>
      </q-card>

      <q-list v-else separator>
        <q-item v-for="schedule in getAvailableSchedules" :key="schedule.id">
          <q-item-section avatar>
            <q-icon name="schedule" color="grey" size="2em" />
          </q-item-section>

          <q-item-section>
            <q-item-label>
              {{ schedule.subjectName }} - {{ schedule.className }}
            </q-item-label>

            <q-item-label caption>
              {{ schedule.startTime }} - {{ schedule.endTime }}
            </q-item-label>
          </q-item-section>

          <q-item-section side>
            <q-btn v-if="!isReadOnlyContext" color="primary" label="Buat Sesi" size="sm"
              @click="handleCreateSession(schedule.id)" />

            <q-badge v-else color="grey">
              Read-only
            </q-badge>
          </q-item-section>
        </q-item>
      </q-list>
    </template>
  </q-page>
</template>

<script setup>
import { onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useTeaching } from '../../composables/useTeaching.js'
import { useContext } from '../../composables/useContext.js'
import ReadOnlyContextBanner from '../../components/common/ReadOnlyContextBanner.vue'

const $q = useQuasar()

const { isReadOnlyContext } = useContext()

const {
  selectedDate,
  sessions,
  isLoading,
  error,
  getAvailableSchedules,
  loadSchedules,
  loadSessions,
  createSession,
  navigateToSession,
} = useTeaching()

onMounted(async () => {
  await Promise.all([
    loadSchedules(),
    loadSessions(),
  ])
})

const handleDateChange = async () => {
  await loadSessions()
}

const handleCreateSession = async (scheduleId) => {
  if (isReadOnlyContext.value) {
    $q.notify({
      type: 'warning',
      message: 'Data historis bersifat read-only',
    })

    return
  }

  try {
    const session = await createSession(scheduleId)

    $q.notify({
      type: 'positive',
      message: 'Sesi mengajar berhasil dibuat',
    })

    navigateToSession(session.id)
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

const getSessionIcon = (status) => {
  const iconMap = {
    scheduled: 'schedule',
    active: 'play_circle',
    completed: 'check_circle',
    cancelled: 'cancel',
  }

  return iconMap[status] || 'help'
}
</script>

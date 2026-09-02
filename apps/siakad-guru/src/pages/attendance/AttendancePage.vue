<!--
FILE: src/pages/attendance/AttendancePage.vue
STATUS: MODIFY
STATUS IMPLEMENTASI: COMPLETE
-->

<template>
  <q-page padding>
    <ReadOnlyContextBanner />

    <div class="row items-center justify-between q-mb-md">
      <div class="text-h5">Attendance</div>

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
      <div class="text-h6 q-mb-sm">Sesi Mengajar - {{ selectedDate }}</div>

      <q-card v-if="sessions.length === 0">
        <q-card-section class="text-center text-grey">
          Tidak ada sesi mengajar untuk tanggal ini.
        </q-card-section>
      </q-card>

      <q-list v-else separator>
        <q-item v-for="session in sessions" :key="session.id" clickable v-ripple @click="handleNavigate(session.id)">
          <q-item-section avatar>
            <q-icon name="how_to_reg" color="primary" size="2em" />
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
            <q-btn color="primary" :label="isReadOnlyContext ? 'Lihat' : 'Absensi'" size="sm"
              @click.stop="handleNavigate(session.id)" />
          </q-item-section>
        </q-item>
      </q-list>
    </template>
  </q-page>
</template>

<script setup>
import { onMounted } from 'vue'
import { useTeaching } from '../../composables/useTeaching.js'
import { useAttendance } from '../../composables/useAttendance.js'
import { useContext } from '../../composables/useContext.js'
import ReadOnlyContextBanner from '../../components/common/ReadOnlyContextBanner.vue'

const {
  sessions,
  selectedDate,
  isLoading,
  loadSessions,
} = useTeaching()

const {
  navigateToSession,
  error,
} = useAttendance()

const { isReadOnlyContext } = useContext()

onMounted(async () => {
  await loadSessions()
})

const handleDateChange = async () => {
  await loadSessions()
}

const handleNavigate = (teachingSessionId) => {
  navigateToSession(teachingSessionId)
}
</script>

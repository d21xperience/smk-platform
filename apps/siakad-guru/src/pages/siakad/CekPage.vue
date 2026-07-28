<template>
  <q-page class="q-pa-md">
    <div class="text-h6 q-mb-md">Dashboard Guru</div>
    <div v-if="!ctx.operational.value.academicYearId" class="text-grey">
      Silakan pilih konteks pembelajaran terlebih dahulu.
      <q-btn label="Pilih Konteks" to="/siakad/contex" flat color="primary" />
    </div>
    <div v-else>
      <div class="text-subtitle1 text-grey q-mb-sm">
        {{ ctx.operational.value.academicYearName }} | {{ ctx.operational.value.semesterName }}
      </div>
      <div class="text-subtitle2 q-mb-md">Jadwal Hari Ini: {{ today }}</div>
      <q-list bordered separator>
        <q-item v-for="session in sessions" :key="session.id" clickable @click="openSession(session)">
          <q-item-section>
            <q-item-label>{{ session.subject }} - {{ session.className }}</q-item-label>
            <q-item-label caption>{{ session.startTime }} - {{ session.endTime }}</q-item-label>
          </q-item-section>
          <q-item-section side>
            <q-badge :color="statusColor(session.status)">{{ statusLabel(session.status) }}</q-badge>
          </q-item-section>
        </q-item>
      </q-list>
      <q-btn v-if="!sessions.length" label="Muat Jadwal" color="primary" @click="loadTodaySessions" class="q-mt-md" />
    </div>
  </q-page>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useTeaching } from '@/composables/useTeaching'
import { useContext } from '@/composables/useContext'
import { useRouter } from 'vue-router'
// import { SessionStatus } from '@/engine/TeachingEngine'

const teaching = useTeaching()
const ctx = useContext()
const router = useRouter()
const today = new Date().toISOString().slice(0, 10)

const sessions = computed(() => teaching.sessions.value)

async function loadTodaySessions() {
  // await teaching.loadSessions(today)
  await teaching.loadSessions('2026-07-18')
}

onMounted(async () => {
  if (ctx.isOperationalContextReady?.value) {
    await loadTodaySessions()
  }
})

function openSession(session) {
  // Navigasi ke halaman sesi dengan tabs
  router.push(`/teaching/${session.id}`)
}

// function statusColor(status) {
//   switch (status) {
//     case SessionStatus.SCHEDULED: return 'grey'
//     case SessionStatus.STARTED: return 'blue'
//     case SessionStatus.IN_PROGRESS: return 'green'
//     case SessionStatus.COMPLETED: return 'orange'
//     case SessionStatus.LOCKED: return 'red'
//     default: return 'grey'
//   }
// }

function statusLabel(status) {
  switch (status) {
    case 'scheduled': return 'Terjadwal'
    case 'started': return 'Dimulai'
    case 'in_progress': return 'Berlangsung'
    case 'completed': return 'Selesai'
    case 'locked': return 'Terkunci'
    default: return status
  }
}
</script>

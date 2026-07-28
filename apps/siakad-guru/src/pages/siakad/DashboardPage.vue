<template>
  <q-page class="q-pa-md dashboard-page">
    <!-- LOADING STATE -->
    <div v-if="loading" class="row justify-center items-center" style="min-height: 60vh">
      <q-spinner-dots color="primary" size="50px" />
      <div class="text-grey-7 q-ml-sm">Memuat data dashboard...</div>
    </div>

    <!-- KONTEN DASHBOARD (Tampil setelah data dimuat) -->
    <template v-else>
      <!-- WELCOME BANNER -->
      <div class="welcome-banner q-pa-lg q-mb-md">
        <div class="row items-center justify-between">
          <div>
            <div class="text-h6 text-weight-bold text-white">
              Selamat {{ greetingTime }}, {{ userName }} 👋
            </div>
            <div class="text-body2 text-grey-3 q-mt-xs">
              {{ todayLabel }} &middot; {{ roleLabel }}
            </div>
          </div>
          <q-icon name="dashboard" size="48px" class="text-gold gt-xs" />
        </div>
      </div>

      <!-- STAT CARDS -->
      <div class="row q-col-gutter-md q-mb-md">
        <div class="col-6 col-md-3" v-for="(stat, index) in stats" :key="index">
          <q-card flat bordered class="stat-card">
            <q-card-section class="row items-center no-wrap">
              <div class="stat-icon q-mr-md" :class="`bg-${stat.color}-1`">
                <q-icon :name="stat.icon" :color="stat.color" size="26px" />
              </div>
              <div>
                <div class="text-h5 text-weight-bold">{{ stat.value }}</div>
                <div class="text-caption text-grey-7">{{ stat.label }}</div>
              </div>
            </q-card-section>
          </q-card>
        </div>
      </div>

      <div class="row q-col-gutter-md">
        <!-- WEEKLY ATTENDANCE CHART -->
        <div class="col-12 col-md-7">
          <q-card flat bordered class="full-height">
            <q-card-section>
              <div class="row items-center justify-between q-mb-md">
                <div class="text-subtitle1 text-weight-bold">Rekap Kehadiran Mingguan</div>
                <q-badge color="blue-1" text-color="primary" class="q-px-sm">Minggu Ini</q-badge>
              </div>

              <div class="bar-chart" v-if="weeklyAttendance.length > 0">
                <div class="bar-chart__col" v-for="day in weeklyAttendance" :key="day.day">
                  <div class="bar-chart__bars">
                    <div class="bar-chart__bar bar-chart__bar--hadir" :style="{ height: day.hadir + '%' }">
                      <q-tooltip>Hadir: {{ day.hadir }}%</q-tooltip>
                    </div>
                    <div class="bar-chart__bar bar-chart__bar--alpha" :style="{ height: day.alpha + '%' }">
                      <q-tooltip>Alpha: {{ day.alpha }}%</q-tooltip>
                    </div>
                  </div>
                  <div class="text-caption text-grey-7 q-mt-xs">{{ day.day }}</div>
                </div>
              </div>
              <div v-else class="text-center text-grey-6 q-pa-lg">Tidak ada data kehadiran minggu ini.</div>

              <div class="row items-center q-gutter-md q-mt-md">
                <div class="row items-center">
                  <div class="legend-dot bg-primary q-mr-xs" />
                  <span class="text-caption text-grey-8">Hadir</span>
                </div>
                <div class="row items-center">
                  <div class="legend-dot bg-red q-mr-xs" />
                  <span class="text-caption text-grey-8">Alpha</span>
                </div>
              </div>
            </q-card-section>
          </q-card>
        </div>

        <!-- JADWAL HARI INI -->
        <div class="col-12 col-md-5">
          <q-card flat bordered class="full-height">
            <q-card-section>
              <div class="text-subtitle1 text-weight-bold q-mb-md">Jadwal Mengajar Hari Ini</div>

              <q-list separator v-if="todaySchedule.length">
                <q-item v-for="item in todaySchedule" :key="item.id">
                  <q-item-section avatar top>
                    <q-chip square color="grey-2" text-color="grey-9" class="text-weight-bold">
                      {{ item.startTime }} - {{ item.endTime }}
                    </q-chip>
                  </q-item-section>
                  <q-item-section>
                    <q-item-label class="text-weight-medium">{{ item.subject }}</q-item-label>
                    <q-item-label caption>{{ item.className }}</q-item-label>
                  </q-item-section>
                  <q-item-section side>
                    <q-badge :color="item.status === 'completed' ? 'positive' : 'grey-5'" outline>
                      {{ item.status === 'completed' ? 'Selesai' : statusLabel(item.status) }}
                    </q-badge>
                  </q-item-section>
                </q-item>
              </q-list>

              <div v-else class="text-center text-grey-6 q-pa-lg">
                <q-icon name="event_available" size="40px" class="q-mb-sm" />
                <div>Tidak ada jadwal mengajar hari ini</div>
              </div>
            </q-card-section>
          </q-card>
        </div>
      </div>

      <!-- ALERTS: SISWA PERLU PERHATIAN -->
      <div class="q-mt-md">
        <div class="text-subtitle1 text-weight-bold q-mb-sm">Perlu Perhatian</div>
        <div class="row q-col-gutter-md">
          <div class="col-12 col-md-6" v-for="alert in criticalAlerts" :key="alert.id">
            <q-card flat bordered class="alert-card">
              <q-card-section class="row items-center no-wrap">
                <q-icon :name="alert.icon" color="red" size="28px" class="q-mr-md" />
                <div>
                  <div class="text-weight-medium">{{ alert.title }}</div>
                  <div class="text-caption text-grey-7">{{ alert.description }}</div>
                </div>
              </q-card-section>
            </q-card>
          </div>

          <div class="col-12" v-if="criticalAlerts.length === 0">
            <q-card flat bordered class="q-pa-md text-center text-grey-6">
              <q-icon name="check_circle" color="positive" size="28px" class="q-mr-sm" />
              Tidak ada catatan penting saat ini
            </q-card>
          </div>
        </div>
      </div>

      <!-- QUICK ACCESS -->
      <div class="q-mt-md">
        <div class="text-subtitle1 text-weight-bold q-mb-sm">Akses Cepat</div>
        <div class="row q-col-gutter-md">
          <div class="col-6 col-md-3" v-for="link in quickLinks" :key="link.label">
            <q-btn no-caps unelevated class="quick-link-btn full-width" align="left" :to="link.to">
              <q-icon :name="link.icon" color="primary" size="24px" class="q-mr-sm" />
              <span class="text-weight-medium">{{ link.label }}</span>
            </q-btn>
          </div>
        </div>
      </div>
    </template>
  </q-page>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useAuth } from '@/composables/useAuth'
import { useTeaching } from '@/composables/useTeaching'
// import { useDashboard } from '@/composables/useDashboard'      // dari statistics cache
import { useStudentProgressStore } from '@/stores/studentProgress.store'
import { useUtils } from '@/composables/helper/useUtils'        // untuk greetingTime, todayLabel

const auth = useAuth()
const teaching = useTeaching()
// const dashboardCache = useDashboard()
const studentProgressStore = useStudentProgressStore()
const { greetingTime, todayLabel } = useUtils()

const loading = computed(() => teaching.loading.value.sessions)

const userName = computed(() => auth.user.value?.name || 'Guru')
const roleLabel = computed(() => {
  const roleMap = {
    guru: 'Guru Mata Pelajaran',
    wali_kelas: 'Wali Kelas',
    admin: 'Administrator SIAKAD'
  }
  return roleMap[auth.user.value?.role] || 'Guru'
})

// STATS: dihitung dari sessions hari ini dan global stats
const stats = computed(() => {
  const sessions = teaching.sessions.value
  const completed = sessions.filter(s => s.status === 'completed' || s.status === 'locked').length
  const ongoing = sessions.filter(s => s.status === 'started' || s.status === 'in_progress').length
  const scheduled = sessions.filter(s => s.status === 'scheduled').length

  return [
    { label: 'Sesi Hari Ini', value: sessions.length, icon: 'event', color: 'primary' },
    { label: 'Sedang Berlangsung', value: ongoing, icon: 'play_circle', color: 'green' },
    { label: 'Selesai', value: completed, icon: 'check_circle', color: 'orange' },
    { label: 'Terjadwal', value: scheduled, icon: 'schedule', color: 'blue' }
  ]
})

// WEEKLY ATTENDANCE: sementara dummy (belum ada proyeksi)
const weeklyAttendance = computed(() => {
  // Di sprint berikutnya akan diganti dengan data dari proyeksi
  return []
})

// TODAY SCHEDULE: dari sessions yang sudah dimuat
const todaySchedule = computed(() => {
  return teaching.sessions.value.map(session => ({
    id: session.id,
    startTime: session.startTime,
    endTime: session.endTime,
    subject: session.subject,
    className: session.className,
    status: session.status
  }))
})

// CRITICAL ALERTS: dari student progress
const criticalAlerts = computed(() => {
  const alerts = []
  const allProgress = studentProgressStore.getAllProgress
  Object.values(allProgress).forEach(student => {
    const attendance = student.attendanceSummary
    if (attendance && attendance.alpha > 3) {
      alerts.push({
        id: student.studentId,
        title: `${student.studentName} - Alpha Tinggi`,
        description: `Tercatat ${attendance.alpha} kali alpha. Segera tindak lanjuti.`,
        icon: 'warning'
      })
    }
  })
  return alerts.slice(0, 4) // maks 4 alert
})

const quickLinks = [
  { label: 'Input Absensi', icon: 'fact_check', to: { name: 'input-absensi' } },
  { label: 'Input Nilai', icon: 'assessment', to: '/siakad/nilai' },
  { label: 'Jadwal Mengajar', icon: 'schedule', to: '/siakad/jadwal-pelajaran' },
  { label: 'Koreksi Data', icon: 'edit_note', to: '/siakad/koreksi' }
]

// Helper untuk label status
function statusLabel(status) {
  const map = { scheduled: 'Terjadwal', started: 'Dimulai', in_progress: 'Berlangsung', completed: 'Selesai', locked: 'Terkunci' }
  return map[status] || status
}

onMounted(async () => {
  // Muat sesi hari ini
  const today = new Date().toISOString().slice(0, 10)
  await teaching.loadSessions(today)
  // Student progress sudah terisi otomatis oleh proyeksi, jadi tidak perlu load manual
})
</script>

<style scoped>
/* style sama seperti sebelumnya, tidak berubah */
.dashboard-page {
  max-width: 1400px;
  margin: 0 auto;
}

.welcome-banner {
  border-radius: 8px;
  background: linear-gradient(135deg, #0d2b4e 0%, #1a4a7a 100%);
  position: relative;
  overflow: hidden;
}

.text-gold {
  color: #d4af37;
}

.stat-card,
.alert-card {
  border-radius: 8px;
  height: 100%;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.full-height {
  height: 100%;
  border-radius: 8px;
}

.alert-card {
  border-left: 4px solid var(--q-red, #c10015);
}

.bar-chart {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  height: 160px;
  gap: 8px;
}

.bar-chart__col {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  height: 100%;
}

.bar-chart__bars {
  flex: 1;
  width: 100%;
  display: flex;
  align-items: flex-end;
  justify-content: center;
  gap: 4px;
}

.bar-chart__bar {
  width: 14px;
  border-radius: 4px 4px 0 0;
  transition: height 0.3s ease;
}

.bar-chart__bar--hadir {
  background-color: var(--q-primary, #1a4a7a);
}

.bar-chart__bar--alpha {
  background-color: #c10015;
}

.legend-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
}

.quick-link-btn {
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  padding: 12px 16px;
}
</style>

<!--
FILE: src/layouts/MainLayout.vue
STATUS: MODIFY
STATUS IMPLEMENTASI: COMPLETE
-->

<template>
  <q-layout view="hHh lpR fFf">
    <q-header elevated class="bg-primary text-white">
      <q-toolbar>
        <q-btn
          flat
          dense
          round
          icon="menu"
          aria-label="Menu"
          @click="toggleLeftDrawer"
        />

        <q-toolbar-title>
          SIAKAD Guru
        </q-toolbar-title>

        <div
          v-if="displayContext"
          class="q-mr-md row items-center q-gutter-xs"
        >
          <q-badge :color="isHistoryModeActive ? 'warning' : 'secondary'">
            {{ contextLabel }}
          </q-badge>

          <q-badge
            v-if="isReadOnlyContext"
            color="negative"
          >
            READ-ONLY
          </q-badge>
        </div>

        <q-btn
          v-if="isHistoryModeActive"
          flat
          dense
          icon="undo"
          label="Kembali ke Aktif"
          class="q-mr-sm"
          @click="handleExitHistory"
        />

        <q-btn
          flat
          dense
          icon="history"
          label="Riwayat"
          class="q-mr-sm"
          @click="handleOpenHistorySelector"
        />

        <q-btn
          v-if="isAuthenticated"
          flat
          dense
          icon="logout"
          label="Logout"
          @click="handleLogout"
        />
      </q-toolbar>
    </q-header>

    <q-drawer
      v-model="leftDrawerOpen"
      show-if-above
      bordered
      class="bg-grey-2"
    >
      <q-list>
        <q-item-label header>Menu Utama</q-item-label>

        <q-item
          clickable
          v-ripple
          :to="{ name: 'dashboard' }"
          exact
        >
          <q-item-section avatar>
            <q-icon name="dashboard" />
          </q-item-section>
          <q-item-section>
            <q-item-label>Dashboard</q-item-label>
          </q-item-section>
        </q-item>

        <q-item
          clickable
          v-ripple
          :to="{ name: 'teaching' }"
        >
          <q-item-section avatar>
            <q-icon name="school" />
          </q-item-section>
          <q-item-section>
            <q-item-label>Teaching</q-item-label>
          </q-item-section>
        </q-item>

        <q-item
          clickable
          v-ripple
          :to="{ name: 'attendance' }"
        >
          <q-item-section avatar>
            <q-icon name="how_to_reg" />
          </q-item-section>
          <q-item-section>
            <q-item-label>Attendance</q-item-label>
          </q-item-section>
        </q-item>

        <q-item
          clickable
          v-ripple
          :to="{ name: 'journal' }"
        >
          <q-item-section avatar>
            <q-icon name="book" />
          </q-item-section>
          <q-item-section>
            <q-item-label>Journal</q-item-label>
          </q-item-section>
        </q-item>

        <q-item
          clickable
          v-ripple
          :to="{ name: 'assessment' }"
        >
          <q-item-section avatar>
            <q-icon name="grading" />
          </q-item-section>
          <q-item-section>
            <q-item-label>Assessment</q-item-label>
          </q-item-section>
        </q-item>

        <q-item
          clickable
          v-ripple
          :to="{ name: 'progress' }"
        >
          <q-item-section avatar>
            <q-icon name="trending_up" />
          </q-item-section>
          <q-item-section>
            <q-item-label>Student Progress</q-item-label>
          </q-item-section>
        </q-item>

        <q-item
          clickable
          v-ripple
          :to="{ name: 'inventory' }"
        >
          <q-item-section avatar>
            <q-icon name="inventory_2" />
          </q-item-section>
          <q-item-section>
            <q-item-label>Inventaris Kelas</q-item-label>
          </q-item-section>
        </q-item>

        <q-item
          clickable
          v-ripple
          :to="{ name: 'academic-calender' }"
        >
          <q-item-section avatar>
            <q-icon name="event" />
          </q-item-section>
          <q-item-section>
            <q-item-label>Kalender Akademik</q-item-label>
          </q-item-section>
        </q-item>

        <template v-if="canViewHomeroomSection">
          <q-separator />

          <q-item-label header>Wali Kelas</q-item-label>

          <q-item
            v-if="canViewHomeroomBilling"
            clickable
            v-ripple
            :to="{ name: 'homeroom-billing' }"
          >
            <q-item-section avatar>
              <q-icon name="payments" />
            </q-item-section>
            <q-item-section>
              <q-item-label>Homeroom Billing</q-item-label>
            </q-item-section>
          </q-item>

          <q-item
            v-if="canViewHomeroomStudents"
            clickable
            v-ripple
            :to="{ name: 'homeroom-students' }"
          >
            <q-item-section avatar>
              <q-icon name="groups" />
            </q-item-section>
            <q-item-section>
              <q-item-label>Daftar Siswa</q-item-label>
            </q-item-section>
          </q-item>

          <q-item
            v-if="canViewHomeroomProgress"
            clickable
            v-ripple
            :to="{ name: 'homeroom-progress' }"
          >
            <q-item-section avatar>
              <q-icon name="analytics" />
            </q-item-section>
            <q-item-section>
              <q-item-label>Rekap Nilai</q-item-label>
            </q-item-section>
          </q-item>

          <q-item
            v-if="canViewHomeroomAttendance"
            clickable
            v-ripple
            :to="{ name: 'homeroom-attendance' }"
          >
            <q-item-section avatar>
              <q-icon name="event_available" />
            </q-item-section>
            <q-item-section>
              <q-item-label>Rekap Kehadiran</q-item-label>
            </q-item-section>
          </q-item>

          <q-item
            v-if="canViewHomeroomJournal"
            clickable
            v-ripple
            :to="{ name: 'homeroom-journal' }"
          >
            <q-item-section avatar>
              <q-icon name="sticky_note_2" />
            </q-item-section>
            <q-item-section>
              <q-item-label>Catatan Siswa</q-item-label>
            </q-item-section>
          </q-item>

          <q-item
            v-if="canViewHomeroomCommunication"
            clickable
            v-ripple
            :to="{ name: 'homeroom-communication' }"
          >
            <q-item-section avatar>
              <q-icon name="contact_phone" />
            </q-item-section>
            <q-item-section>
              <q-item-label>Komunikasi Ortu</q-item-label>
            </q-item-section>
          </q-item>

          <q-item
            v-if="canViewHomeroomReport"
            clickable
            v-ripple
            :to="{ name: 'homeroom-print-report' }"
          >
            <q-item-section avatar>
              <q-icon name="print" />
            </q-item-section>
            <q-item-section>
              <q-item-label>Cetak Laporan</q-item-label>
            </q-item-section>
          </q-item>
        </template>

        <q-separator />

        <q-item-label header>Dokumen</q-item-label>

        <q-item
          clickable
          v-ripple
          :to="{ name: 'teacher-documents' }"
        >
          <q-item-section avatar>
            <q-icon name="folder" />
          </q-item-section>
          <q-item-section>
            <q-item-label>Teacher Documents</q-item-label>
          </q-item-section>
        </q-item>
      </q-list>
    </q-drawer>

    <q-page-container>
      <router-view />
    </q-page-container>

    <HistoricalContextSelector
      v-model="isHistorySelectorOpen"
      @selected="handleHistorySelected"
    />
  </q-layout>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useQuasar } from 'quasar'
import { useAuth } from '../composables/useAuth.js'
import { useContext } from '../composables/useContext.js'
import HistoricalContextSelector from '../components/common/HistoricalContextSelector.vue'
import { formatContextLabel } from '../utils/contextFormatter.js'

const $q = useQuasar()

const { isAuthenticated, logout, currentUser } = useAuth()

const {
  displayContext,
  isHistoryModeActive,
  isReadOnlyContext,
  exitHistoricalContext,
} = useContext()

const leftDrawerOpen = ref(true)
const isHistorySelectorOpen = ref(false)

const toggleLeftDrawer = () => {
  leftDrawerOpen.value = !leftDrawerOpen.value
}

const contextLabel = computed(() => {
  return formatContextLabel(displayContext.value)
})

const hasPermission = (permission) => {
  if (!currentUser.value) return false

  const permissions = currentUser.value.permissions || []
  const role = currentUser.value.role
  const isAdmin = role === 'admin' || role === 'superadmin'

  return isAdmin || permissions.includes(permission)
}

const canViewHomeroomBilling = computed(() => {
  return hasPermission('homeroom.billing.read')
})

const canViewHomeroomStudents = computed(() => {
  return hasPermission('homeroom.student.read')
})

const canViewHomeroomProgress = computed(() => {
  return hasPermission('homeroom.progress.read')
})

const canViewHomeroomAttendance = computed(() => {
  return hasPermission('homeroom.attendance.read')
})

const canViewHomeroomJournal = computed(() => {
  return hasPermission('homeroom.journal.read')
})

const canViewHomeroomCommunication = computed(() => {
  return hasPermission('homeroom.communication.read')
})

const canViewHomeroomReport = computed(() => {
  return hasPermission('homeroom.report.print')
})

const canViewHomeroomSection = computed(() => {
  return canViewHomeroomBilling.value ||
    canViewHomeroomStudents.value ||
    canViewHomeroomProgress.value ||
    canViewHomeroomAttendance.value ||
    canViewHomeroomJournal.value ||
    canViewHomeroomCommunication.value ||
    canViewHomeroomReport.value
})

const handleOpenHistorySelector = () => {
  isHistorySelectorOpen.value = true
}

const handleHistorySelected = (historicalContext) => {
  $q.notify({
    type: 'info',
    message: `Context historis aktif: ${formatContextLabel(historicalContext)}. Data bersifat read-only.`,
  })
}

const handleExitHistory = () => {
  exitHistoricalContext()

  $q.notify({
    type: 'positive',
    message: 'Kembali ke context aktif.',
  })
}

const handleLogout = () => {
  $q.dialog({
    title: 'Konfirmasi',
    message: 'Apakah Anda yakin ingin keluar?',
    cancel: true,
    persistent: true,
  }).onOk(async () => {
    await logout()
  })
}
</script>

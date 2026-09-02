<template>
  <q-layout view="hHh LpR fFf">
    <!-- Topbar -->
    <q-header elevated>
      <q-toolbar>
        <q-btn flat @click="drawerLeft = !drawerLeft" round dense icon="menu" />
        <q-toolbar-title>Operator Sekolah</q-toolbar-title>

        <!-- Indikator Koneksi Dapodik -->
        <q-chip
          :color="dapodikStatus.color"
          text-color="white"
          dense
          icon-right="circle"
          class="q-mr-sm"
        >
          {{ dapodikStatus.label }}
        </q-chip>

        <q-btn flat round dense icon="dark_mode" @click="toggleDarkMode" />
      </q-toolbar>
    </q-header>

    <!-- Sidebar Drawer -->
    <q-drawer
      v-model="drawerLeft"
      show-if-above
      :mini="miniState"
      @mouseover="miniState = false"
      @mouseout="miniState = true"
      :width="250"
      :breakpoint="700"
      bordered
    >
      <q-scroll-area class="fit">
        <q-list padding>
          <!-- Dashboard -->
          <q-item
            clickable
            v-ripple
            :to="{ name: 'ops-dashboard' }"
            exact
            active-class="text-primary"
          >
            <q-item-section avatar><q-icon name="dashboard" /></q-item-section>
            <q-item-section>Dashboard</q-item-section>
          </q-item>

          <!-- Antrean Tugas (dengan badge) -->
          <q-item
            clickable
            v-ripple
            :to="{ name: 'ops-antrean-tugas' }"
            active-class="text-primary"
          >
            <q-item-section avatar><q-icon name="assignment" /></q-item-section>
            <q-item-section>Antrean Tugas</q-item-section>
            <q-item-section side>
              <q-badge color="red" floating>{{ pendingCount }}</q-badge>
            </q-item-section>
          </q-item>

          <!-- Sinkronisasi -->
          <q-item clickable v-ripple :to="{ name: 'ops-sinkronisasi' }" active-class="text-primary">
            <q-item-section avatar><q-icon name="sync" /></q-item-section>
            <q-item-section>Sinkronisasi Dapodik</q-item-section>
          </q-item>

          <!-- Riwayat Aktivitas -->
          <q-item clickable v-ripple :to="{ name: 'ops-riwayat' }" active-class="text-primary">
            <q-item-section avatar><q-icon name="history" /></q-item-section>
            <q-item-section>Riwayat Aktivitas</q-item-section>
          </q-item>
        </q-list>
      </q-scroll-area>
    </q-drawer>

    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script setup>
import { useQuasar } from 'quasar'
import { ref } from 'vue'

const $q = useQuasar()

// Drawer state
const drawerLeft = ref(true)
const miniState = ref(false)

// Dummy pending count untuk badge
const pendingCount = ref(5)

// Status koneksi Dapodik (contoh, nanti bisa dari store)
const dapodikStatus = ref({ color: 'green', label: 'Terhubung' })

// Dark mode toggle
const isDark = ref(false)
function toggleDarkMode() {
  isDark.value = !isDark.value
  $q.dark.set(isDark.value)
}
</script>

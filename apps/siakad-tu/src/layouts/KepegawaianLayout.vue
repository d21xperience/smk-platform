<template>
  <q-layout view="lHh Lpr lFf">
    <q-header elevated class="bg-primary text-white">
      <q-toolbar>
        <q-btn dense flat round icon="menu" @click="toggleLeftDrawer" />
        <q-toolbar-title>SIAKAD SMK Pasundan</q-toolbar-title>
        <div>{{ user?.name || 'Operator' }}</div>
        <q-btn flat round icon="logout" @click="logout" class="q-ml-sm" />
      </q-toolbar>
    </q-header>

    <q-drawer v-model="leftDrawerOpen" show-if-above bordered class="bg-grey-2">
      <q-list>
        <q-item-label header class="text-primary">Modul Kepegawaian</q-item-label>
        <q-item clickable v-ripple to="/kepegawaian/guru">
          <q-item-section avatar><q-icon name="people" /></q-item-section>
          <q-item-section>Data Induk GTK</q-item-section>
        </q-item>
        <q-item clickable v-ripple to="/kepegawaian/beban-mengajar">
          <q-item-section avatar><q-icon name="schedule" /></q-item-section>
          <q-item-section>Beban Mengajar</q-item-section>
        </q-item>
        <q-item clickable v-ripple to="/kepegawaian/cuti">
          <q-item-section avatar><q-icon name="event_busy" /></q-item-section>
          <q-item-section>Pengajuan Cuti/SPT</q-item-section>
        </q-item>
        <q-item clickable v-ripple to="/kepegawaian/arsip">
          <q-item-section avatar><q-icon name="folder" /></q-item-section>
          <q-item-section>Arsip Dokumen</q-item-section>
        </q-item>
        <q-item clickable v-ripple to="/kepegawaian/kinerja">
          <q-item-section avatar><q-icon name="insights" /></q-item-section>
          <q-item-section>Dashboard Kinerja</q-item-section>
        </q-item>
      </q-list>
    </q-drawer>

    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from 'stores/authStore'

const router = useRouter()
const authStore = useAuthStore()
const leftDrawerOpen = ref(true)

function toggleLeftDrawer() {
  leftDrawerOpen.value = !leftDrawerOpen.value
}

function logout() {
  authStore.logout()
  router.push('/login')
}

const user = authStore.user
</script>

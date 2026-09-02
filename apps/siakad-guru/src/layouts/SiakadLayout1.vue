<template>
  <q-layout view="hHh lpR fFf">
    <!-- HEADER -->
    <q-header elevated class="bg-primary text-white">
      <q-toolbar>
        <q-btn
          flat
          dense
          round
          icon="menu"
          aria-label="Menu"
          @click="toggleLeftDrawer"
          class="lt-md"
        />

        <q-toolbar-title class="row items-center no-wrap">
          <q-avatar size="32px" class="q-mr-sm">
            <img src="~@/assets/logo-smk.png" alt="Logo" />
          </q-avatar>
          <span class="text-weight-bold">SIAKAD - SMK Pasundan Jatinangor</span>
        </q-toolbar-title>
        <div>
          <!-- <q-toolbar> -->
          <!-- <q-btn flat label="Batal" v-close-popup class="text-weight-bold" /> -->
          <!-- <q-toolbar-title class="text-center text-weight-bold">Ruang Kerja Penulis Artikel</q-toolbar-title> -->
          <!-- </q-toolbar> -->
        </div>
        <!-- User Info & Actions (Desktop) -->
        <div class="gt-sm row items-center q-gutter-sm">
          <q-btn flat round dense icon="notifications" />
          <q-chip outline color="white" text-color="white" icon="person">
            {{ userName }}
          </q-chip>
          <q-btn flat round dense icon="logout" @click="handleLogout">
            <q-tooltip>Keluar</q-tooltip>
          </q-btn>
        </div>
      </q-toolbar>
    </q-header>

    <!-- SIDEBAR DRAWER -->
    <q-drawer v-model="leftDrawerOpen" show-if-above bordered class="bg-grey-1" :width="260">
      <q-scroll-area class="fit">
        <q-list padding>
          <q-item-label header class="text-grey-8 text-weight-bold">MENU UTAMA</q-item-label>

          <q-item
            clickable
            v-ripple
            :to="{ name: 'siakad-dashboard' }"
            exact
            active-class="text-primary bg-blue-1"
          >
            <q-item-section avatar><q-icon name="dashboard" /></q-item-section>
            <q-item-section>Dashboard</q-item-section>
          </q-item>

          <q-item
            clickable
            v-ripple
            :to="{ name: 'contex' }"
            exact
            active-class="text-primary bg-blue-1"
          >
            <q-item-section avatar><q-icon name="fact_check" /></q-item-section>
            <q-item-section>Pilih kontek</q-item-section>
          </q-item>

          <q-item
            clickable
            v-ripple
            :to="{ name: 'input-absensi' }"
            exact
            active-class="text-primary bg-blue-1"
          >
            <q-item-section avatar><q-icon name="fact_check" /></q-item-section>
            <q-item-section>Absensi Siswa</q-item-section>
          </q-item>

          <!-- <q-separator class="q-my-md" /> -->

          <q-item
            clickable
            v-ripple
            :to="{ name: 'jadwal-pelajaran' }"
            active-class="text-primary bg-blue-1"
          >
            <q-item-section avatar><q-icon name="schedule" /></q-item-section>
            <q-item-section>Jadwal Mengajar</q-item-section>
          </q-item>
          <q-item
            clickable
            v-ripple
            :to="{ name: 'kelender-akademik' }"
            active-class="text-primary bg-blue-1"
          >
            <q-item-section avatar><q-icon name="event" /></q-item-section>
            <q-item-section>Kalender Akademik</q-item-section>
          </q-item>
          <q-item-label header class="text-grey-8 text-weight-bold">LAINNYA</q-item-label>

          <q-item
            clickable
            v-ripple
            :to="{ name: 'catatan-siswa' }"
            active-class="text-primary bg-blue-1"
          >
            <q-item-section avatar><q-icon name="history" /></q-item-section>
            <q-item-section>Catatan Siswa</q-item-section>
          </q-item>

          <q-item clickable v-ripple to="/siakad/nilai" active-class="text-primary bg-blue-1">
            <q-item-section avatar><q-icon name="assessment" /></q-item-section>
            <q-item-section>Input Nilai</q-item-section>
          </q-item>
          <q-item clickable v-ripple :to="{ name: 'profil' }" active-class="text-primary bg-blue-1">
            <q-item-section avatar><q-icon name="person" /></q-item-section>
            <q-item-section>Profil</q-item-section>
          </q-item>
          <q-separator class="q-my-md" />
          <!-- Menu untuk keluar -->
          <q-item clickable v-ripple @click="handleLogout" active-class="text-primary bg-blue-1">
            <q-item-section avatar><q-icon name="logout" /></q-item-section>
            <q-item-section>Keluar</q-item-section>
          </q-item>
        </q-list>
      </q-scroll-area>
    </q-drawer>

    <!-- MAIN CONTENT -->
    <q-page-container class="bg-grey-2">
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script setup>
import { computed, ref } from 'vue'
// import { useRouter } from 'vue-router'
import { useAuth } from '@/composables/useAuth'
import { useQuasar } from 'quasar'
import { useRouter } from 'vue-router'

const router = useRouter()
const $q = useQuasar()
// const router = useRouter()
const auth = useAuth()

const leftDrawerOpen = ref(false)

const userName = computed(() => auth.user?.name || 'Guru')

function toggleLeftDrawer() {
  leftDrawerOpen.value = !leftDrawerOpen.value
}

function handleLogout() {
  $q.dialog({
    title: 'Konfirmasi',
    message: 'Apakah Anda yakin ingin keluar dari SIAKAD?',
    cancel: true,
    persistent: true,
  }).onOk(() => {
    auth.logout()
    router.push('/auth/siakad')
  })
}
</script>

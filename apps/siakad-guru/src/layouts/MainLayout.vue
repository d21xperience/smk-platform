<template>
  <q-layout view="hHh LpR fFf">
    <!-- Header -->
    <q-header elevated class="bg-primary text-white">
      <q-toolbar>
        <!-- TOMBOL HAMBURGER: Hanya muncul di layar HP/Tablet kecil (lt-sm) -->
        <q-btn flat dense round icon="menu" aria-label="Menu" class="lt-sm q-mr-sm" @click="toggleLeftDrawer" />

        <!-- Logo dan Judul Sekolah -->
        <q-toolbar-title class="row items-center no-wrap">
          <q-avatar size="32px">
            <img src="~@/assets/logo-smk.png" alt="Logo" onerror="this.style.display='none'" />
          </q-avatar>
          <!-- text-subtitle1 pada HP agar teks tidak terlalu besar dan terpotong -->
          <span class="q-ml-sm text-weight-bold text-subtitle1 text-sm-h6 ellipsis">
            SMK Pasundan Jatinangor
          </span>
        </q-toolbar-title>
        <div class="gt-xs row items-center no-wrap text-weight-medium">

          <!-- 1. GRUP TAB UTAMA -->
          <q-tabs v-model="currentTab" align="left" shrink stretch class="text-white">
            <q-route-tab to="/" label="Beranda" exact />
            <q-route-tab to="/mading" label="Mading Digital" />
            <!-- <q-route-tab to="/faq" label="PSB" /> -->
            <!-- <q-route-tab to="/simulasi" label="Simulasi" /> -->
          </q-tabs>

          <!-- 2. SUB-MENU DROPDOWN 1: JURUSAN -->
          <q-btn-dropdown flat stretch no-caps label="PROFIL" content-class="bg-white text-grey-9 shadow-5"
            class="text-white font-tab-match">
            <q-list min-width="180px">
              <q-item clickable v-ripple to="/jurusan/tbsm">
                <!-- <q-item-section avatar><q-icon name="motorcycle" color="red-7" /></q-item-section> -->
                <q-item-section>Sekolah</q-item-section>
              </q-item>
              <q-item clickable v-ripple to="/jurusan/tkj">
                <!-- <q-item-section avatar><q-icon name="computer" color="blue-7" /></q-item-section> -->
                <q-item-section>Visi & Misi</q-item-section>
              </q-item>
              <q-item clickable v-ripple to="/jurusan/tkr">
                <!-- <q-item-section avatar><q-icon name="directions_car" color="orange-8" /></q-item-section> -->
                <q-item-section>Teknik Kendaraan Ringan</q-item-section>
              </q-item>
              <q-item clickable v-ripple to="/jurusan/akl">
                <!-- <q-item-section avatar><q-icon name="payments" color="green-7" /></q-item-section> -->
                <q-item-section>Akuntansi & Keuangan</q-item-section>
              </q-item>
              <q-item clickable v-ripple to="/jurusan/otkp">
                <!-- <q-item-section avatar><q-icon name="work" color="purple-7" /></q-item-section> -->
                <q-item-section>Manajemen Perkantoran</q-item-section>
              </q-item>
            </q-list>
          </q-btn-dropdown>

          <!-- 3. SUB-MENU DROPDOWN 2: SIAKAD (PORTAL INTERNAL) -->
          <q-btn-dropdown flat stretch no-caps label="PORTAL" content-class="bg-white text-grey-9 shadow-5"
            class="text-white font-tab-match">
            <q-list min-width="180px">
              <!-- Jika SIAKAD menggunakan link eksternal web lain, pakai href dan target="_blank" -->
              <q-item clickable v-ripple :to="{ name: 'siakadLogin' }" target="_blank">
                <!-- <q-item-section avatar><q-icon name="school" color="primary" /></q-item-section> -->
                <q-item-section>SIAKAD</q-item-section>
              </q-item>

              <q-item clickable v-ripple href="https://e-rapor.smkpasjat.my.id" target="_blank">
                <!-- <q-item-section avatar><q-icon name="co_present" color="secondary" /></q-item-section> -->
                <q-item-section>e-Rapor</q-item-section>
              </q-item>

              <q-item clickable v-ripple href="https://sekolah.sch.id" target="_blank">
                <!-- <q-item-section avatar><q-icon name="family_restroom" color="accent" /></q-item-section> -->
                <q-item-section>CBT</q-item-section>
              </q-item>
            </q-list>
          </q-btn-dropdown>

        </div>

        <q-space class="gt-xs" />
        <!-- Tombol Dark Mode -->
        <q-btn flat round dense :icon="$q.dark.isActive ? 'light_mode' : 'dark_mode'" @click="toggleDarkMode" />
        <q-btn flat round dense :icon="loggedin ? 'logout' : 'login'" @click="handleLogin">
          <q-tooltip>{{ loggedin ? 'Logout' : 'Login' }}</q-tooltip>
        </q-btn>
      </q-toolbar>
    </q-header>

    <!-- SIDEBAR DRAWER: Menu navigasi khusus untuk pengguna HP -->
    <q-drawer v-model="leftDrawerOpen" side="left" bordered behavior="mobile" class="bg-grey-1">
      <q-scroll-area class="fit">
        <q-list padding class="text-grey-8">
          <q-item-label header class="text-weight-bold text-primary">
            MENU NAVIGASI
          </q-item-label>

          <q-item clickable v-ripple to="/" exact active-class="text-primary text-weight-bold">
            <q-item-section avatar><q-icon name="home" /></q-item-section>
            <q-item-section>Beranda</q-item-section>
          </q-item>

          <q-item clickable v-ripple to="/mading" active-class="text-primary text-weight-bold">
            <q-item-section avatar><q-icon name="collections" /></q-item-section>
            <q-item-section>Mading Digital</q-item-section>
          </q-item>

          <q-item clickable v-ripple to="/faq" active-class="text-primary text-weight-bold">
            <q-item-section avatar><q-icon name="quiz" /></q-item-section>
            <q-item-section>FAQ</q-item-section>
          </q-item>

          <q-separator q-my-md />

          <!-- Tombol PPDB Darurat di dalam Menu HP jika status aktif -->
          <q-item v-if="ppdbStatus === 'ACTIVE'" clickable v-ripple @click="bukaFormPPDB"
            class="bg-accent text-white rounded-borders q-mx-md">
            <q-item-section avatar><q-icon name="school" color="white" /></q-item-section>
            <q-item-section class="text-weight-bold">Daftar PPDB Online</q-item-section>
          </q-item>
        </q-list>
      </q-scroll-area>
    </q-drawer>

    <!-- Konten Halaman Aktif -->
    <q-page-container>
      <!-- Banner Pengumuman Dinamis Sebelum PPDB Dibuka -->
      <q-banner v-if="ppdbStatus === 'UPCOMING'"
        class="bg-amber text-black text-center text-weight-bold text-caption text-sm-body2">
        <q-icon name="info" size="xs" /> PPDB TA 2026/2027 dibuka dalam {{ daysLeft }} hari lagi!
      </q-banner>

      <router-view />
    </q-page-container>

    <!-- Sticky Button PPDB DINAMIS (Responsif dengan ukuran mini pada HP) -->
    <q-page-sticky v-if="ppdbStatus === 'ACTIVE'" position="bottom-right" :offset="[16, 16]">
      <q-btn fab :mini="$q.screen.xs" icon="school" color="accent" :label="$q.screen.xs ? '' : 'Daftar PPDB'"
        @click="bukaFormPPDB" class="shadow-5" :class="{ 'q-px-md': !$q.screen.xs }">
        <q-badge color="red" floating animate>Buka</q-badge>
        <q-tooltip v-if="$q.screen.xs" anchor="center left" self="center right">Daftar PPDB</q-tooltip>
      </q-btn>
    </q-page-sticky>

  </q-layout>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useQuasar } from 'quasar'
import { useRouter } from 'vue-router'
const router = useRouter()
const $q = useQuasar()
const currentTab = ref('beranda')
const leftDrawerOpen = ref(false)
const loggedin = ref(false)
// Pengaturan Jadwal PPDB
const ppdbConfig = ref({
  startDate: new Date('2026-06-01T00:00:00'),
  endDate: new Date('2026-07-15T23:59:59')
})

const ppdbStatus = computed(() => {
  const now = new Date()
  if (now < ppdbConfig.value.startDate) return 'UPCOMING'
  if (now >= ppdbConfig.value.startDate && now <= ppdbConfig.value.endDate) return 'ACTIVE'
  return 'CLOSED'
})

const daysLeft = computed(() => {
  const diff = ppdbConfig.value.startDate - new Date()
  return Math.ceil(diff / (1000 * 60 * 60 * 24))
})

function toggleLeftDrawer() {
  leftDrawerOpen.value = !leftDrawerOpen.value
}

function toggleDarkMode() {
  $q.dark.toggle()
}

function bukaFormPPDB() {
  window.open('https://sekolah.sch.id', '_blank')
}

const handleLogin = () => {
  router.push('/login')
}




</script>

<style scoped>
/* Menjamin teks panjang pada judul sekolah terpotong rapi dengan titik (...) di layar super kecil */
.ellipsis {
  text-overflow: ellipsis;
  overflow: hidden;
  white-space: nowrap;
}

.bg-gradient-cta {
  background: linear-gradient(135deg, #1976D2 0%, #0D47A1 100%);
}

.z-top {
  position: relative;
  z-index: 2;
}
</style>

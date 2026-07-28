<template>
  <q-layout view="lHh Lpr lFf">
    <!-- Header Atas Portal Siswa -->
    <q-header elevated class="bg-teal-8 text-white">
      <q-toolbar>
        <q-btn flat dense round icon="menu" aria-label="Menu" @click="toggleLeftDrawer" />

        <q-avatar size="32px" class="q-ml-sm bg-white q-pa-xs">
          <q-icon name="school" color="teal-8" size="20px" />
        </q-avatar>

        <q-toolbar-title class="text-subtitle1 text-weight-bold">
          Portal Mandiri Siswa SMK
        </q-toolbar-title>

        <q-badge color="white" text-color="teal-9" label="Calon Siswa Baru" class="q-pa-sm text-weight-bold" />
      </q-toolbar>
    </q-header>

    <!-- Sidebar Navigasi Khusus Fase 1 PPDB -->
    <q-drawer v-model="leftDrawerOpen" show-if-above bordered :width="250" class="bg-grey-1">
      <q-scroll-area class="fit">
        <!-- Banner Mini Profil Pendaftar -->
        <div class="q-pa-md bg-teal-1 text-teal-10 flex flex-center text-center">
          <div>
            <q-avatar size="50px" color="teal-8" text-color="white" icon="person" class="q-mb-xs" />
            <div class="text-weight-bold text-subtitle2">Portal Calon Siswa</div>
            <div class="text-grey-7" style="font-size: 11px;">TA 2026/2027</div>
          </div>
        </div>

        <q-separator />

        <!-- Daftar Navigasi Menu Phase 1 -->
        <q-list padding class="menu-list">
          <q-item-label header class="text-weight-bold text-uppercase text-grey-6" style="font-size: 10px;">
            Alur Pendaftaran PPDB
          </q-item-label>

          <!-- 1. Dashboard -->
          <q-item clickable v-ripple to="/siswa-portal/dashboard-ppdb" exact
            active-class="bg-teal-1 text-teal text-weight-bold" class="text-grey-8">
            <q-item-section avatar><q-icon name="dashboard" /></q-item-section>
            <q-item-section>Dashboard PPDB</q-item-section>
          </q-item>

          <!-- 2. Isi Biodata -->
          <q-item clickable v-ripple to="/siswa-portal/biodata" exact
            active-class="bg-teal-1 text-teal text-weight-bold" class="text-grey-8">
            <q-item-section avatar><q-icon name="assignment" /></q-item-section>
            <q-item-section>Isi Biodata Lengkap</q-item-section>
          </q-item>

          <!-- 3. Upload Dokumen -->
          <q-item clickable v-ripple to="/siswa-portal/dokumen" exact
            active-class="bg-teal-1 text-teal text-weight-bold" class="text-grey-8">
            <q-item-section avatar><q-icon name="cloud_upload" /></q-item-section>
            <q-item-section>Upload Berkas Digital</q-item-section>
          </q-item>

          <!-- 4. Riwayat Keuangan -->
          <q-item clickable v-ripple to="/siswa-portal/pembayaran" exact
            active-class="bg-teal-1 text-teal text-weight-bold" class="text-grey-8">
            <q-item-section avatar><q-icon name="payments" /></q-item-section>
            <q-item-section>Biaya Pendaftaran</q-item-section>
          </q-item>

          <q-separator class="q-my-md" />

          <!-- Tombol Logout Keluar Sistem -->
          <q-item clickable v-ripple to="/logout-session" class="text-red-9 bg-red-1">
            <q-item-section avatar><q-icon name="power_settings_new" color="red-9" /></q-item-section>
            <q-item-section class="text-weight-medium">Keluar Aplikasi</q-item-section>
          </q-item>
        </q-list>
      </q-scroll-area>
    </q-drawer>

    <!-- Konten Halaman Pengisian Mandiri Siswa -->
    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script setup>
import { ref } from 'vue'

const leftDrawerOpen = ref(true)
const toggleLeftDrawer = () => {
  leftDrawerOpen.value = !leftDrawerOpen.value
}
</script>

<style scoped>
.menu-list .q-item {
  border-radius: 0 24px 24px 0;
  margin-right: 12px;
  margin-bottom: 4px;
}
</style>

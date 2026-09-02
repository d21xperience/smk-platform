<template>
  <q-layout view="clR LpR lFf">
    <!-- Header Admin -->
    <q-header elevated class="bg-dark text-white">
      <q-toolbar>
        <q-btn flat dense round icon="menu" @click="toggleLeftDrawer" />
        <q-toolbar-title>Panel Admin Sekolah</q-toolbar-title>

        <!-- Info User yang Login -->
        <div class="q-gutter-sm row items-center">
          <q-avatar color="primary" text-color="white" size="32px">
            {{ user.name.charAt(0) }}
          </q-avatar>
          <div class="text-subtitle2">{{ user.name }} ({{ user.role.toUpperCase() }})</div>
          <q-btn flat round dense icon="logout" @click="handleLogout" />
        </div>
      </q-toolbar>
    </q-header>

    <!-- Sidebar Menu Navigasi Berdasarkan Role -->
    <q-drawer v-model="leftDrawerOpen" show-if-above bordered class="bg-grey-1">
      <q-list>
        <q-item-label header class="text-weight-bold">Menu Navigasi</q-item-label>

        <!-- MENU UNTUK SEMUA ROLE (Dashboard Utama) -->
        <q-item clickable v-ripple to="/admin/dashboard">
          <q-item-section avatar><q-icon name="dashboard" /></q-item-section>
          <q-item-section>Ringkasan</q-item-section>
        </q-item>

        <!-- MENU KHUSUS ADMIN & TATA USAHA -->
        <q-expansion-item v-if="['admin', 'tu'].includes(user.role)" icon="assignment" label="Manajemen PPDB"
          header-class="text-weight-medium">
          <q-list class="q-pl-md">
            <!-- Hanya Admin yang bisa atur tanggal buka-tutup -->
            <q-item v-if="user.role === 'admin'" clickable v-ripple to="/admin/ppdb-config">
              <q-item-section><q-icon name="settings_suggest" /></q-item-section>
              <q-item-section>Pengaturan Jadwal PPDB</q-item-section>
            </q-item>
            <q-item clickable v-ripple to="/admin/ppdb-data">
              <q-item-section><q-icon name="people" /></q-item-section>
              <q-item-section>Data Pendaftar & Berkas</q-item-section>
            </q-item>
          </q-list>
        </q-expansion-item>

        <!-- MENU UNTUK OSIS & KAPROG (Upload Poster/Mading) -->
        <q-item v-if="['admin', 'osis', 'kaprog'].includes(user.role)" clickable v-ripple to="/admin/posters">
          <q-item-section avatar><q-icon name="collections" /></q-item-section>
          <q-item-section>Kelola Poster & Mading</q-item-section>
        </q-item>

        <!-- MENU UNTUK GURU & KAPROG (Artikel & Berita) -->
        <q-item v-if="['admin', 'guru', 'kaprog'].includes(user.role)" clickable v-ripple to="/admin/articles">
          <q-item-section avatar><q-icon name="article" /></q-item-section>
          <q-item-section>Tulis Artikel / Berita</q-item-section>
        </q-item>
      </q-list>
    </q-drawer>

    <!-- Konten Halaman Kerja Admin -->
    <q-page-container class="bg-grey-2">
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script setup>
import { ref } from 'vue'

const leftDrawerOpen = ref(false)

// Simulasi Data User yang sedang login (Anda bisa mengubah role ini untuk uji coba)
// Pilihan: 'admin', 'tu', 'kaprog', 'osis', 'guru'
const user = ref({
  name: 'Pak Budi',
  role: 'kaprog'
})

const toggleLeftDrawer = () => {
  leftDrawerOpen.value = !leftDrawerOpen.value
}

const handleLogout = () => {
  console.log('User logged out')
}
</script>

<template>
  <q-layout view="hHh Lpr lFf">
    <!-- Header Atas Khusus PPDB -->
    <q-header elevated class="bg-indigo-9 text-white">
      <q-toolbar>
        <q-btn flat dense round icon="menu" aria-label="Menu" @click="toggleLeftDrawer" />

        <q-avatar size="32px" class="q-ml-sm bg-white q-pa-xs">
          <q-icon name="app_registration" color="indigo-9" size="20px" />
        </q-avatar>

        <q-toolbar-title class="text-subtitle1 text-weight-bold">
          Sistem PPDB & Pendaftaran Angkatan Baru
        </q-toolbar-title>

        <!-- Informasi Tahun Ajaran Aktif -->
        <q-badge color="orange-9" label="Tahun Ajaran 2026/2027" class="q-pa-sm text-weight-bold" />
      </q-toolbar>
    </q-header>

    <!-- Sidebar Navigasi Khusus Panitia PPDB -->
    <q-drawer v-model="leftDrawerOpen" show-if-above bordered :width="260" class="bg-grey-1">
      <q-scroll-area class="fit">
        <!-- Banner Mini Profil Petugas PPDB -->
        <div class="q-pa-md bg-indigo-1 text-indigo-10 flex flex-center text-center">
          <div>
            <q-avatar size="48px" color="indigo" text-color="white" icon="badge" class="q-mb-xs" />
            <div class="text-weight-bold text-caption">PANITIA PPDB</div>
            <div class="text-grey-7" style="font-size: 11px;">SMK Swasta Indonesia</div>
          </div>
        </div>

        <q-separator />

        <!-- Daftar Menu Terbatas Hanya untuk PPDB -->
        <q-list padding class="menu-list">
          <q-item-label header class="text-weight-bold text-uppercase text-grey-6" style="font-size: 11px;">
            Modul Pendaftaran
          </q-item-label>

          <!-- Menu Karantina Siswa Baru & Kelulusan NIS -->
          <q-item clickable v-ripple :to="{ name: 'tambah-siswa' }" exact
            active-class="bg-indigo-1 text-indigo text-weight-bold" class="text-grey-8">
            <q-item-section avatar>
              <q-icon name="person_add_alt_1" />
            </q-item-section>
            <q-item-section>
              Karantina & NIS Massal
            </q-item-section>
          </q-item>

          <!-- Contoh Menu Tambahan Khusus PPDB di Masa Depan -->
          <!-- <q-item clickable v-ripple to="/ppdb/statistik-asal-sekolah"
            active-class="bg-indigo-1 text-indigo text-weight-bold" class="text-grey-5" disable>
            <q-item-section avatar>
              <q-icon name="school" />
            </q-item-section>
            <q-item-section>
              Analisis Asal Sekolah
              <q-item-label caption class="text-orange-9">Segera Hadir</q-item-label>
            </q-item-section>
          </q-item> -->
          <!-- SEBELUMNYA DISABLE, SEKARANG AKTIFKAN SEPERTI INI -->
          <q-item clickable v-ripple :to="{ name: 'analisis-asal-sekolah' }"
            active-class="bg-indigo-1 text-indigo text-weight-bold" class="text-grey-8">
            <q-item-section avatar>
              <q-icon name="analytics" /> <!-- Ubah ikon menjadi analytics -->
            </q-item-section>
            <q-item-section>
              Analisis Asal Sekolah
            </q-item-section>
          </q-item>

          <q-separator class="q-my-md" />
          <q-item clickable v-ripple :to="{ name: 'generate-kelas' }"
            active-class="bg-indigo-1 text-indigo text-weight-bold" class="text-grey-8">
            <q-item-section avatar>
              <q-icon name="group_work" />
            </q-item-section>
            <q-item-section>
              Otomatisasi Kelas X
            </q-item-section>
          </q-item>
          <q-separator class="q-my-md" />

          <!-- Tombol Kembali ke Dashboard Utama Kesiswaan Sekolah -->
          <q-item clickable v-ripple to="/dashboard" class="text-red-9 bg-red-1"
            style="border-radius: 0 24px 24px 0; margin-right: 12px;">
            <q-item-section avatar>
              <q-icon name="arrow_back" color="red-9" />
            </q-item-section>
            <q-item-section class="text-weight-medium">
              Kembali ke SIAKAD Utama
            </q-item-section>
          </q-item>
        </q-list>
      </q-scroll-area>
    </q-drawer>

    <!-- Tempat Merender Halaman PpdbSiswa.vue -->
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

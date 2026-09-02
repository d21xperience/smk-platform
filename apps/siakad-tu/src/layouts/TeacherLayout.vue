<!-- <template>
  <q-layout view="lHh Lpr lFf">
    <q-header elevated>
      <q-toolbar>
        <q-btn flat dense round icon="menu" aria-label="Menu" @click="toggleLeftDrawer" />

        <q-toolbar-title>
          Quasar App
        </q-toolbar-title>

        <div>Quasar v{{ $q.version }}</div>
      </q-toolbar>
    </q-header>

    <q-drawer v-model="leftDrawerOpen" show-if-above bordered>
      <q-list>
        <q-item-label header>
          Essential Links
        </q-item-label>

        <EssentialLink v-for="link in linksList" :key="link.title" v-bind="link" />
      </q-list>
    </q-drawer>

    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script setup>
import { ref } from 'vue'
import EssentialLink from 'components/EssentialLink.vue'

const linksList = [
  {
    title: 'Docs',
    caption: 'quasar.dev',
    icon: 'school',
    link: 'https://quasar.dev'
  },
  {
    title: 'Github',
    caption: 'github.com/quasarframework',
    icon: 'code',
    link: 'https://github.com/quasarframework'
  },
  {
    title: 'Discord Chat Channel',
    caption: 'chat.quasar.dev',
    icon: 'chat',
    link: 'https://chat.quasar.dev'
  },
  {
    title: 'Forum',
    caption: 'forum.quasar.dev',
    icon: 'record_voice_over',
    link: 'https://forum.quasar.dev'
  },
  {
    title: 'Twitter',
    caption: '@quasarframework',
    icon: 'rss_feed',
    link: 'https://twitter.quasar.dev'
  },
  {
    title: 'Facebook',
    caption: '@QuasarFramework',
    icon: 'public',
    link: 'https://facebook.quasar.dev'
  },
  {
    title: 'Quasar Awesome',
    caption: 'Community Quasar projects',
    icon: 'favorite',
    link: 'https://awesome.quasar.dev'
  }
]

const leftDrawerOpen = ref(false)

function toggleLeftDrawer() {
  leftDrawerOpen.value = !leftDrawerOpen.value
}
</script>
 -->
<template>
  <q-layout view="lHh Lpr lFf">
    <q-header elevated>
      <q-toolbar>
        <q-btn flat dense round icon="menu" aria-label="Menu" @click="toggleLeftDrawer" />
        <q-toolbar-title>CBT Ujian Online</q-toolbar-title>
        <div>Guru: {{ teacherName }}</div>
      </q-toolbar>
    </q-header>

    <q-drawer v-model="leftDrawerOpen" show-if-above bordered>
      <q-list>
        <!-- Dashboard -->
        <q-item clickable v-ripple to="/dashboard" exact>
          <q-item-section avatar>
            <q-icon name="dashboard" />
          </q-item-section>
          <q-item-section>Dashboard</q-item-section>
        </q-item>

        <!-- Kelola Soal (dengan submenu) -->
        <q-expansion-item
          icon="edit_note"
          label="Kelola Soal"
          caption="Buat, edit, hapus soal"
          group="menuGroup"
        >
          <q-item clickable v-ripple to="/soal/buat">
            <q-item-section avatar><q-icon name="add_circle" /></q-item-section>
            <q-item-section>Buat Soal</q-item-section>
          </q-item>
          <q-item clickable v-ripple to="/soal/bank">
            <q-item-section avatar><q-icon name="folder" /></q-item-section>
            <q-item-section>Bank Soal</q-item-section>
          </q-item>
          <q-item clickable v-ripple to="/soal/kategori">
            <q-item-section avatar><q-icon name="category" /></q-item-section>
            <q-item-section>Kategori & Subjek</q-item-section>
          </q-item>
        </q-expansion-item>

        <!-- Monitoring Ujian (submenu dari fitur yang sudah dibuat) -->
        <q-expansion-item
          icon="monitoring"
          label="Monitoring Ujian"
          caption="Analisis hasil dan performa"
          group="menuGroup"
        >
          <q-item clickable v-ripple to="/monitoring/statistik-soal">
            <q-item-section avatar><q-icon name="bar_chart" /></q-item-section>
            <q-item-section>Statistik Soal</q-item-section>
            <q-item-section side><q-badge color="blue">TK, DB, %Benar</q-badge></q-item-section>
          </q-item>
          <q-item clickable v-ripple to="/monitoring/distribusi-nilai">
            <q-item-section avatar><q-icon name="insert_chart" /></q-item-section>
            <q-item-section>Distribusi Nilai & Kelulusan</q-item-section>
          </q-item>
          <q-item clickable v-ripple to="/monitoring/heatmap">
            <q-item-section avatar><q-icon name="grid_on" /></q-item-section>
            <q-item-section>Heatmap Jawaban Siswa</q-item-section>
          </q-item>
          <q-item clickable v-ripple to="/monitoring/export">
            <q-item-section avatar><q-icon name="file_download" /></q-item-section>
            <q-item-section>Export CSV/Excel</q-item-section>
          </q-item>
        </q-expansion-item>

        <!-- Kelola Ujian -->
        <q-expansion-item
          icon="quiz"
          label="Ujian"
          caption="Jadwal & pelaksanaan"
          group="menuGroup"
        >
          <q-item clickable v-ripple to="/ujian/buat">
            <q-item-section avatar><q-icon name="assignment_add" /></q-item-section>
            <q-item-section>Buat Ujian Baru</q-item-section>
          </q-item>
          <q-item clickable v-ripple to="/ujian/daftar">
            <q-item-section avatar><q-icon name="list_alt" /></q-item-section>
            <q-item-section>Daftar Ujian</q-item-section>
          </q-item>
          <q-item clickable v-ripple to="/ujian/jadwal">
            <q-item-section avatar><q-icon name="event" /></q-item-section>
            <q-item-section>Jadwal Ujian</q-item-section>
          </q-item>
        </q-expansion-item>

        <!-- Data Siswa -->
        <q-item clickable v-ripple to="/siswa">
          <q-item-section avatar>
            <q-icon name="people" />
          </q-item-section>
          <q-item-section>Data Siswa</q-item-section>
        </q-item>

        <!-- Laporan -->
        <q-item clickable v-ripple to="/laporan">
          <q-item-section avatar>
            <q-icon name="description" />
          </q-item-section>
          <q-item-section>Laporan Keseluruhan</q-item-section>
        </q-item>

        <!-- Pengaturan -->
        <q-item clickable v-ripple to="/pengaturan">
          <q-item-section avatar>
            <q-icon name="settings" />
          </q-item-section>
          <q-item-section>Pengaturan</q-item-section>
        </q-item>
      </q-list>

      <template #mini>
        <q-list>
          <q-item clickable v-ripple to="/dashboard" class="justify-center">
            <q-item-section avatar><q-icon name="dashboard" /></q-item-section>
          </q-item>
          <!-- versi mini untuk item lain bisa ditambahkan sesuai kebutuhan -->
        </q-list>
      </template>
    </q-drawer>

    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script setup>
import { ref } from 'vue'

const leftDrawerOpen = ref(false)
const teacherName = ref('Bapak/Ibu Guru') // nanti bisa diambil dari store/auth

function toggleLeftDrawer() {
  leftDrawerOpen.value = !leftDrawerOpen.value
}
</script>

<style scoped>
/* opsional: styling tambahan */
.q-expansion-item__content .q-item {
  padding-left: 32px;
}
</style>

<!-- <template>
  <q-layout view="hHh lpR fFf">
    <q-header elevated class="bg-dark text-white">
      <q-toolbar>
        <q-toolbar-title>Monitoring: {{ store.user?.name }}</q-toolbar-title>
        <q-btn flat icon="logout" label="Keluar" @click="handleLogout" />
      </q-toolbar>
    </q-header>
    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useExamStore } from 'stores/exam-store'
const router = useRouter()
const store = useExamStore()
const handleLogout = () => { store.logout(); router.push('/login') }
</script> -->

<template>
  <q-layout view="lHh Lpr lFf">
    <q-header elevated>
      <q-toolbar>
        <q-btn flat dense round icon="menu" aria-label="Menu" @click="toggleLeftDrawer" />
        <q-toolbar-title>SIAKAD - Admin Panel</q-toolbar-title>
        <div>Admin: {{ adminName }}</div>
        <q-btn flat round icon="logout" @click="logout" class="q-ml-md" />
      </q-toolbar>
    </q-header>

    <q-drawer v-model="leftDrawerOpen" show-if-above bordered>

      <EssentialLink v-for="link in linksList" :key="link.title" v-bind="link" />

      <!-- Mini drawer version -->
      <template #mini>
        <q-list>
          <q-item clickable v-ripple to="/admin/dashboard" class="justify-center">
            <q-item-section avatar><q-icon name="dashboard" /></q-item-section>
          </q-item>
          <!-- tambahkan item mini lainnya sesuai kebutuhan -->
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
// import { useRouter } from 'vue-router'
import { useQuasar } from 'quasar'
import EssentialLink from 'components/EssentialLink.vue'

const $q = useQuasar()
// const router = useRouter()

const leftDrawerOpen = ref(false)
const adminName = ref('Admin Utama') // nanti dari store/auth

function toggleLeftDrawer() {
  leftDrawerOpen.value = !leftDrawerOpen.value
}

function logout() {
  $q.dialog({
    title: 'Konfirmasi',
    message: 'Apakah Anda yakin ingin keluar?',
    cancel: true,
    persistent: true
  }).onOk(() => {
    // hapus token, redirect ke login
    // router.push('/login')
  })
}
const linksList = [
  {
    title: 'Dashboard',
    icon: 'dashboard',
    link: { name: 'surat-menyurat-dashboard' },
  },
  {
    title: 'Surat masuk',
    icon: 'inbox',
    link: { name: 'surat-masuk' },
  },
  {
    title: 'Surat keluar',
    icon: 'outbox',
    link: { name: 'surat-keluar' },
  },

  {
    title: 'Disposisi',
    icon: 'assignment_turned_in',
    link: { name: 'disposisi' },
  },
  {
    title: 'Klasifikasi',
    icon: 'category',
    link: { name: 'klasifikasi' },
  },

  // {
  //   title: 'Manajemen Pengguna',
  //   icon: 'people',
  //   caption: 'Kelola akun & role',
  //   group: 'adminGroup',
  //   children: [
  //     { title: 'Akun Guru', icon: 'school', link: '/admin/pengguna/guru' },
  //     // { title: 'Akun Siswa', icon: 'face', link: '/admin/pengguna/siswa' },
  //     { title: 'Admin Lain', icon: 'admin_panel_settings', link: '/admin/pengguna/admin' },
  //     { title: 'Manajemen Role & Hak Akses', icon: 'lock', link: { name: 'admin-pengaturan' } },
  //   ]
  // },
  // {
  //   title: 'Kepegawaian',
  //   icon: 'event_note',
  //   caption: 'Kelola kepegawaian',
  //   group: 'adminGroup',
  //   children: [
  //     { title: 'Data Guru', icon: 'school', link: { name: 'admin-data-guru' } },
  //     { title: 'Beban Mengajar', icon: 'schedule', link: { name: 'admin-beban-mengajar' } },
  //     { title: 'Cuti', icon: 'event_busy', link: { name: 'admin-cuti-mengajar' } },
  //     { title: 'Arsip', icon: 'folder', link: { name: 'admin-arsip-dokumen' } },
  //   ]
  // },
  // {
  //   title: 'Kurikulum',
  //   icon: 'assignment',
  //   caption: 'Kelola Kurikulum',
  //   group: 'adminGroup',
  //   children: [
  //     { title: 'Daftar Ujian', icon: 'list', link: '/admin/ujian' },
  //     { title: 'Bank Soal', icon: 'quiz', link: '/admin/bank-soal' },
  //     { title: 'Hasil Ujian', icon: 'bar_chart', link: '/admin/hasil-ujian' },
  //   ]
  // },
  // {
  //   title: 'Kesiswaan',
  //   icon: 'people',
  //   caption: 'Buat & kelola ujian',
  //   group: 'adminGroup',
  //   children: [
  //     { title: 'Daftar Ujian', icon: 'list', link: '/admin/ujian' },
  //     { title: 'Bank Soal', icon: 'quiz', link: '/admin/bank-soal' },
  //     { title: 'Hasil Ujian', icon: 'bar_chart', link: '/admin/hasil-ujian' },
  //   ]
  // },
  {
    title: 'Pengaturan',
    icon: 'settings',
    link: '/admin/settings',
  },
]
</script>

<style scoped>
.q-expansion-item__content .q-item {
  padding-left: 32px;
}

/* untuk submenu dalam submenu (soal di dalam konten ujian) */
.q-expansion-item__content .q-expansion-item .q-expansion-item__content .q-item {
  padding-left: 48px;
}
</style>

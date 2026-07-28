<template>
  <AppLayout app-title="SIAKAD - SMK Pasundan Jatinangor" :logo="logoUrl" :user-name="userName"
    :menu-groups="menuGroups" @logout="handleLogout">
    <!-- Slot default: konten utama halaman -->
    <router-view />
  </AppLayout>
</template>

<script setup>
import { useAuth } from '@/composables/useAuth'
import { AppLayout } from '@smk-platform/ui-components'
import { useQuasar } from 'quasar'
import { computed } from 'vue'
import { useRouter } from 'vue-router'

// ============================================================
// 1. Konfigurasi Aplikasi
// ============================================================

const router = useRouter()
const $q = useQuasar()
const auth = useAuth()

// Logo (gunakan path absolut dari public/assets atau alias @)
const logoUrl = new URL('@/assets/logo-smk.png', import.meta.url).href

// Nama User (dari Auth context)
const userName = computed(() => auth.user?.name || 'Guru')

// ============================================================
// 2. Definisi Menu (sesuai kebutuhan SIAKAD Guru)
// ============================================================

const menuGroups = [
  {
    label: 'MENU UTAMA',
    items: [
      { label: 'Dashboard', icon: 'dashboard', to: { name: 'siakad-dashboard' } },
      { label: 'Pilih Konteks', icon: 'fact_check', to: { name: 'contex' } },
      { label: 'Absensi Siswa', icon: 'fact_check', to: { name: 'input-absensi' } },
    ],
  },
  {
    label: 'AKADEMIK',
    items: [
      { label: 'Jadwal Mengajar', icon: 'schedule', to: { name: 'jadwal-pelajaran' } },
      { label: 'Kalender Akademik', icon: 'event', to: { name: 'kelender-akademik' } },
    ],
  },
  {
    label: 'LAINNYA',
    items: [
      { label: 'Catatan Siswa', icon: 'history', to: { name: 'catatan-siswa' } },
      { label: 'Input Nilai', icon: 'assessment', to: '/siakad/nilai' },
      { label: 'Profil', icon: 'person', to: { name: 'profil' } },
    ],
    separatorAfter: true,
  },
]

// ============================================================
// 3. Event Handler Logout
// ============================================================

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

<!--
  Tidak perlu style di sini karena semua style dari komponen shared,
  atau bisa tambahkan style spesifik jika perlu
-->
<style scoped>
/* Tambahkan style lokal jika diperlukan */
</style>

<template>
  <q-layout view="lHh Lpr lFf">
    <!-- 1. HEADER UTAMA DASHBOARD ADMIN -->
    <q-header elevated class="bg-dark text-white">
      <q-toolbar>
        <!-- Tombol Hamburger Menu (Buka/Tutup Sidebar) -->
        <q-btn flat dense round icon="menu" aria-label="Menu" @click="toggleLeftDrawer" />

        <q-toolbar-title class="text-subtitle1 text-sm-h6 text-weight-bold">
          Dashboard <span class="text-amber">CMS Pasundan Jatinangor</span>
        </q-toolbar-title>

        <q-space />

        <!-- INFO PROFIL USER YANG SEDANG LOGIN -->
        <div class="row items-center no-wrap q-gutter-sm">
          <q-chip
            icon="account_circle"
            color="grey-8"
            text-color="white"
            class="q-px-md text-weight-medium"
          >
            {{ currentUserName }}
            <q-badge color="amber-9" class="q-ml-sm text-weight-bold uppercase-text">
              {{ currentUserRole }}
            </q-badge>
          </q-chip>

          <!-- Tombol Keluar Sesi (Logout) -->
          <q-btn flat round dense icon="logout" color="negative" @click="confirmLogout">
            <q-tooltip>Keluar Aplikasi</q-tooltip>
          </q-btn>
        </div>
      </q-toolbar>
    </q-header>

    <!-- 2. SIDEBAR NAVIGASI DINAMIS BERDASARKAN ROLE (RBAC) -->
    <q-drawer v-model="leftDrawerOpen" show-if-above bordered :width="260" class="bg-grey-1">
      <q-scroll-area class="fit">
        <div class="q-pa-md text-center bg-grey-2 border-bottom">
          <div class="text-weight-bold text-primary">SMK Pasundan Jatinangor</div>
          <div class="text-caption text-grey-6">Panel Kontrol Konten Web</div>
        </div>

        <q-list padding class="menu-list">
          <q-item-label header class="text-weight-bold text-grey-7">MENU UTAMA</q-item-label>

          <!-- MENU UNTUK SEMUA ROLE (Dashboard Ringkasan) -->
          <q-item
            clickable
            v-ripple
            to="/admin/dashboard"
            exact
            active-class="text-primary text-weight-bold bg-blue-1"
          >
            <q-item-section avatar><q-icon name="dashboard" /></q-item-section>
            <q-item-section>Ringkasan Sistem</q-item-section>
          </q-item>

          <!-- ========================================== -->
          <!-- KELOMPOK MENU: HANYA UNTUK SUPER ADMIN & TU -->
          <!-- ========================================== -->
          <q-expansion-item
            v-if="['admin', 'tu'].includes(currentUserRole)"
            icon="assignment"
            label="Manajemen PPDB"
            header-class="text-weight-medium text-grey-8"
            default-opened
          >
            <q-list class="q-pl-md">
              <!-- Menu Setting Waktu Hanya untuk Super Admin -->
              <q-item
                v-if="currentUserRole === 'admin'"
                clickable
                v-ripple
                to="/admin/ppdb-config"
                active-class="text-primary text-weight-bold"
              >
                <q-item-section avatar><q-icon name="date_range" size="xs" /></q-item-section>
                <q-item-section class="text-caption text-sm-body2"
                  >Jadwal & Kuota PPDB</q-item-section
                >
              </q-item>

              <!-- Menu Validasi Data Berkas untuk Admin & TU -->
              <q-item
                clickable
                v-ripple
                to="/admin/ppdb-data"
                active-class="text-primary text-weight-bold"
              >
                <q-item-section avatar><q-icon name="groups" size="xs" /></q-item-section>
                <q-item-section class="text-caption text-sm-body2">Data Pendaftar</q-item-section>
              </q-item>
            </q-list>
          </q-expansion-item>

          <!-- ========================================== -->
          <!-- KELOMPOK MENU: UNTUK OSIS & KAPROG (POSTER)-->
          <!-- ========================================== -->
          <q-item
            v-if="['admin', 'osis', 'kaprog'].includes(currentUserRole)"
            clickable
            v-ripple
            to="/admin/posters"
            active-class="text-primary text-weight-bold bg-blue-1"
          >
            <q-item-section avatar><q-icon name="add_photo_alternate" /></q-item-section>
            <q-item-section>Kelola Poster Mading</q-item-section>
          </q-item>

          <!-- ========================================== -->
          <!-- KELOMPOK MENU: UNTUK GURU & KAPROG (ARTIKEL)-->
          <!-- ========================================== -->
          <q-item
            v-if="['admin', 'guru', 'kaprog'].includes(currentUserRole)"
            clickable
            v-ripple
            to="/admin/articles"
            active-class="text-primary text-weight-bold bg-blue-1"
          >
            <q-item-section avatar><q-icon name="newspaper" /></q-item-section>
            <q-item-section>Tulis Artikel Berita</q-item-section>
          </q-item>

          <q-separator q-my-md />

          <!-- MENU UTK KEMBALI LIHAT HALAMAN DEPAN WEB PUBLIK -->
          <q-item clickable v-ripple to="/" target="_blank">
            <q-item-section avatar><q-icon name="open_in_new" color="grey-6" /></q-item-section>
            <q-item-section class="text-grey-6 text-weight-medium"
              >Lihat Website Utama</q-item-section
            >
          </q-item>
        </q-list>
      </q-scroll-area>
    </q-drawer>

    <!-- 3. AREA RENDERING KONTEN HALAMAN INDIVIDUAL ADMIN -->
    <q-page-container class="bg-grey-3">
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useQuasar } from 'quasar'

const router = useRouter()
const $q = useQuasar()

const leftDrawerOpen = ref(false)
const currentUserRole = ref('guest')
const currentUserName = ref('Staf Sekolah')

onMounted(() => {
  // Ambil data session role dan nama yang di-set sewaktu login
  // currentUserRole.value = localStorage.getItem('user_role') || 'guest'
  // currentUserName.value = localStorage.getItem('user_name') || 'User Internal'
  // // Pengamanan Darurat: Jika tidak sengaja bypass url tanpa login, tendang balik ke /login
  // if (currentUserRole.value === 'guest') {
  //   router.push('/login')
  // }
})

function toggleLeftDrawer() {
  leftDrawerOpen.value = !leftDrawerOpen.value
}

// Dialog Konfirmasi Keluar dari Sistem Admin
function confirmLogout() {
  $q.dialog({
    title: 'Konfirmasi Keluar',
    message: 'Apakah Anda yakin ingin keluar dari sesi administrasi ini?',
    cancel: { label: 'Batal', flat: true },
    ok: { label: 'Ya, Keluar', color: 'negative', unelevated: true },
    persistent: true,
  }).onOk(() => {
    // Bersihkan sesi di browser
    localStorage.clear()
    $q.notify({
      type: 'info',
      message: 'Sesi ditutup. Anda berhasil keluar.',
      position: 'top',
    })
    // Kembalikan ke halaman login utama
    router.push('/login')
  })
}
</script>

<style scoped>
.border-bottom {
  border-bottom: 1px solid #e0e0e0;
}

.uppercase-text {
  text-transform: uppercase;
}

.menu-list .q-item {
  border-radius: 0 24px 24px 0;
  margin-right: 12px;
}
</style>

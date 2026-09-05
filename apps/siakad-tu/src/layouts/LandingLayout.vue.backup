<template>
  <q-layout view="hHh lpR fFf">
    <q-header elevated>
      <q-toolbar>
        <!-- TOMBOL HAMBURGER: Hanya muncul di layar HP/Tablet kecil (lt-sm) -->
        <q-btn flat dense round icon="menu" aria-label="Menu" class="lt-sm q-mr-sm" @click="toggleLeftDrawer" />
        <q-toolbar-title>SIAKAD - SMK</q-toolbar-title>
        <div class="gt-xs row items-center no-wrap text-weight-medium">
          <q-tabs align="center">
            <q-route-tab :to="{ name: 'landing-home' }" icon="home" />
            <q-route-tab :to="{ name: 'landing-kurikulum' }" label="Kurikulum" />
            <q-route-tab :to="{ name: 'landing-kesiswaan' }" label="Kesiswaan" />
            <q-route-tab :to="{ name: 'landing-sarpras' }" label="Sarpras" />
            <q-route-tab :to="{ name: 'landing-hubin' }" label="hubin" />
            <q-route-tab :to="{ name: 'landing-psb' }" label="psb" />
            <DropDownMenu label="PORTAL" :items="itemPortal" align="center" />
            <q-btn flat icon="search" label="cari..." @click="isSearchOpen = true" class="q-ml-md">
              <q-tooltip>Cari siswa</q-tooltip>
            </q-btn>
          </q-tabs>
        </div>
        <q-space class="gt-xs" />
        <div class="text-caption text-grey-4 q-mt-xs">
          <div v-if="isAuthenticated">
            <AcademicYearSelector />
          </div>
          <div v-else>
            T.A. 2026/2027
          </div>
        </div>
        <q-btn flat round :icon="isAuthenticated ? 'logout' : 'login'" @click="handleAuthenticated" class="q-ml-md">
          <q-tooltip>{{ isAuthenticated ? 'Keluar' : 'Login' }}</q-tooltip>
        </q-btn>
      </q-toolbar>
      <q-toolbar v-show="true" class="gt-xs row items-center no-wrap text-weight-medium bg-grey-12">
        <DropDownMenu label="Ur. Kesiswaan" :items="itemKesiswaan" class="text-black" />
        <DropDownMenu label="Ur. Kepegawaian" :items="itemKepegawaian" class="text-black" />
        <DropDownMenu label="Ur. Sarpras" :items="itemSarpras" class="text-black" />
        <DropDownMenu label="Ur. Hubin" :items="itemHubin" class="text-black" />
        <DropDownMenu label="Ur. PSB" :items="itemPsb" class="text-black" />
        <DropDownMenu label="Persuratan" :items="itemPersuratan" class="text-black" />

      </q-toolbar>
    </q-header>
    <q-page-container>
      <router-view />
    </q-page-container>
    <!-- Footer -->
    <q-footer class="text-white q-py-sm">
      <div class="text-center text-caption">
        <div class="q-mt-xs">SIAKAD SMK &copy; 2026. (Made with <span class="text-red-14">❤</span> by. Deden M.J.)</div>
      </div>
    </q-footer>


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
          <!-- <q-item v-if="ppdbStatus === 'ACTIVE'" clickable v-ripple @click="bukaFormPPDB"
            class="bg-accent text-white rounded-borders q-mx-md">
            <q-item-section avatar><q-icon name="school" color="white" /></q-item-section>
            <q-item-section class="text-weight-bold">Daftar PPDB Online</q-item-section>
          </q-item> -->
        </q-list>
      </q-scroll-area>
    </q-drawer>
    <!-- KONSUMSI KOMPONEN TERPISAH DI SINI -->
    <SearchStudentDialog v-model="isSearchOpen" />
    <!-- ✨ KOMPONEN TERPISAH: Dialog Login ✨ -->
    <LoginDialog v-model="authStore.showLoginDialog" />
    <!-- Logout Confirmation Dialog -->
    <q-dialog v-model="showLogoutDialog" :aria-label="'Konfirmasi keluar'">
      <q-card>
        <q-card-section class="row items-center">
          <q-avatar icon="logout" color="primary" text-color="white" aria-hidden="true" />
          <span class="q-ml-sm text-h6">Konfirmasi Keluar</span>
        </q-card-section>

        <q-card-section>
          Apakah Anda yakin ingin keluar dari panel?
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="Batal" color="grey" v-close-popup aria-label="Batal" />
          <q-btn flat label="Keluar" color="primary" @click="logout" aria-label="Keluar" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-layout>
</template>

<script setup>
import AcademicYearSelector from '@/components/AcademicYearSelector.vue'
import DropDownMenu from '@/components/DropDownMenu.vue'
import LoginDialog from '@/components/LoginDialog.vue'
import SearchStudentDialog from '@/components/SearchStudentDialog.vue'
import { useAuthStore } from '@/stores/authStore'
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
const router = useRouter()
const leftDrawerOpen = ref(false)
// const isAuthenticated = ref(false)
const showLogoutDialog = ref(false)
const isSearchOpen = ref(false)
const authStore = useAuthStore()
const itemPersuratan = [
  { title: 'Dashboard', to: '/persuratan/dashboard' },
  { title: 'Surat Keluar', to: '/persuratan/surat-keluar' },
  { title: 'Surat Masuk', to: '/persuratan/surat-masuk' },
  { title: 'Disposisi', to: '/persuratan/disposisi' },
  { title: 'Klasifikasi', to: '/persuratan/klasifikasi' }
]
const itemKepegawaian = [
  { title: 'Data Induk GTK', to: '/kepegawaian/data-guru' },
  { title: 'Beban Mengajar', to: '/kepegawaian/beban-mengajar' },
  { title: 'Arsip', to: '/kepegawaian/arsip' },
]
const itemHubin = [
  { title: 'Mitra', to: '/hubin/mitra' },
  { title: 'Penempatan', to: '/hubin/penempatan' },
  { title: 'Logbook', to: '/hubin/logbook' },
  { title: 'Penilaian', to: '/hubin/penilaian' },
  { title: 'Sertifikat', to: '/hubin/sertifikat' },
]
const itemPsb = [
  { title: 'Tambah Siswa', to: '/psb/siswa/baru' },
  { title: 'Rekap Asal Sekolah', to: '/psb/analisis-asal-sekolah' },
  { title: 'Kelas', to: '/psb/generate-kelas' },
]
const itemPortal = [
  { title: 'Our Website', url: 'https://kemdikbud.go.id', target: '_blank' },
  { title: 'CBT', url: 'https://kemdikbud.go.id', target: '_blank' },
  { title: 'e-Rapor', url: 'https:e-rapor.smkpasjat.my.id', target: '_blank' },
  { title: 'Ruang GTK', url: 'https://guru.kemendikdasmen.go.id/', target: '_blank' },
  { title: 'Info GTK', url: 'https://info.gtk.kemendikdasmen.go.id/', target: '_blank' },
  { title: 'SIMPKB', url: 'https://paspor-gtk.simpkb.id/casgpo/login?service=https%3A%2F%2Fppg-backend.simpkb.id%2Fauth%2Flogin', target: '_blank' },
]

const itemKesiswaan = [
  { title: 'Dashboard', to: '/kesiswaan/dashboard' },
  { title: 'Direktori Siswa', to: '/kesiswaan/data-siswa' },
  { title: 'Buku Induk Siswa', to: '/kesiswaan/buku-induk-siswa' },
  { title: 'Mutasi Siswa', to: '/kesiswaan/manajemen-mutasi-siswa' },
  { title: 'Ploting Berkas', to: '/kesiswaan/ploting-berkas-siswa' },
  { title: 'Sertifikat UKK', to: '/kesiswaan/sertifikat-ukk-siswa' },
  { title: 'Asuransi Kerja', to: '/kesiswaan/asuransi-kerja-siswa' },
  { title: 'Surat Aktif', to: '/kesiswaan/surat-aktif-siswa' },
  { title: 'Legalisir', to: '/kesiswaan/legalisir-alumni' },
  { title: 'Tracer Study', to: '/kesiswaan/tracer-study-alumni' },
  { title: 'Loker', to: '/kesiswaan/loker' },
  { title: 'Pelaporan', to: '/kesiswaan/rekap-pelaporan' },
]
const itemSarpras = [
  { title: 'Dashboard', to: '/sarpras/dashboard' },
  {
    title: 'Inventaris', icon: 'wallet', children: [
      // { title: 'Tambah', to: '/sarpras/inventaris/tambah' },
      { title: 'Logistik', to: '/sarpras/inventaris/logistik' },
      { title: 'Ruangan', to: '/sarpras/inventaris/ruangan' },
      { title: 'Cetak QR', to: '/sarpras/inventaris/cetak-qr' },
    ]
  },
  {
    title: 'Servis', icon: 'room_service', children: [
      { title: 'Tiket', to: '/sarpras/servis/tiket' },
      { title: 'Jadwal', to: '/sarpras/servis/jadwal' },
      { title: 'Vendor', to: '/sarpras/servis/vendor' }
    ]
  },
  {
    title: 'Sirkulasi', icon: 'sync', children: [
      { title: 'Pinjam', to: '/sarpras/sirkulasi/pinjam' },
      { title: 'Jadwal', to: '/sarpras/sirkulasi/minta-atk' },
    ]
  },
  {
    title: 'Mutasi', icon: 'delete_sweep', children: [
      { title: 'Log', to: '/sarpras/mutasi/log' },
      { title: 'Hapus', to: '/sarpras/mutasi/hapus' },
    ]
  },
]
const isAuthenticated = computed(() => authStore.isAuthenticated)

const handleAuthenticated = () => {
  if (isAuthenticated.value) {
    showLogoutDialog.value = true
  } else {
    authStore.triggerLogin()
  }
}



const toggleLeftDrawer = () => {
  leftDrawerOpen.value = !leftDrawerOpen.value
}
const logout = () => {
  authStore.logout()
  showLogoutDialog.value = false
  router.push({ name: 'landing-home' })
}



</script>

<style scoped></style>

<template>
  <q-layout view="hHh lpR fFf">
    <!-- Header -->
    <q-header elevated class="bg-white text-dark">
      <q-toolbar>
        <!-- Mobile menu button -->
        <q-btn flat dense round icon="menu" aria-label="Menu" class="lt-sm" @click="toggleLeftDrawer" />

        <!-- Logo -->
        <q-toolbar-title class="flex no-wrap items-center">
          <div>
            <q-icon name="school" size="2em" class="q-mr-sm" />
            <div class="column">
              <div class="text-h6 text-weight-bold">SDP</div>
              <div class="text-caption gt-xs">{{ schoolLabel }}</div>
            </div>
          </div>
        </q-toolbar-title>

        <q-space />

        <!-- Desktop Navigation -->
        <div class="gt-xs row items-center no-wrap text-weight-medium">
          <q-tabs align="center" active-color="primary" indicator-color="primary">
            <q-route-tab v-for="item in menuItems" :key="item.routeName" :to="{ name: item.routeName }"
              :icon="item.icon" :label="item.name" />
            <DropDownMenu label="PORTAL" :items="itemPortal" align="center" />
            <q-btn flat icon="search" label="cari..." @click="openSearch" class="q-ml-md">
              <q-tooltip>Cari siswa</q-tooltip>
            </q-btn>
          </q-tabs>
        </div>
        <q-space />
        <!-- Context Selector (Desktop) -->
        <div class="gt-sm q-mr-md">
          <ContextSelector />
        </div>
        <!-- User menu (jika sudah login) -->
        <q-space class="lt-sm" />

        <div v-if="isAuthenticated" class="row items-center q-gutter-sm">
          <q-btn flat round>
            <q-avatar color="primary" text-color="white" size="2em">
              {{ currentUser?.name?.charAt(0) || 'U' }}
            </q-avatar>
            <q-menu>
              <q-list style="min-width: 150px">
                <q-item clickable v-close-popup @click="handleLogout">
                  <q-item-section avatar>
                    <q-icon name="logout" />
                  </q-item-section>
                  <q-item-section>Logout</q-item-section>
                </q-item>
              </q-list>
            </q-menu>
          </q-btn>
        </div>
        <div v-else class="row items-center q-gutter-sm">
          <q-btn flat icon="login" @click="handleLoginDialog" />
        </div>
      </q-toolbar>

      <TuMenuToolbar :visible="isAuthenticated" @menu-click="onMenuSelected" />

    </q-header>

    <!-- Mobile Drawer -->
    <q-drawer v-model="leftDrawerOpen" bordered class="bg-grey-1">
      <q-list>
        <q-item-label header class="text-weight-bold text-primary">
          Menu Navigasi
        </q-item-label>

        <q-item v-for="item in menuItems" :key="item.routeName" clickable :to="{ name: item.routeName }" v-close-popup>
          <q-item-section avatar>
            <q-icon :name="item.icon" />
          </q-item-section>
          <q-item-section>
            {{ item.name }}
          </q-item-section>
        </q-item>

        <q-separator />

        <q-item-label header>Portal</q-item-label>

        <q-item v-for="(portal, index) in itemPortal" :key="index" clickable @click="portal.action" v-close-popup>
          <q-item-section avatar>
            <q-icon :name="portal.icon" />
          </q-item-section>
          <q-item-section>
            {{ portal.label }}
          </q-item-section>
        </q-item>
      </q-list>
    </q-drawer>

    <!-- Main Content -->
    <q-page-container>
      <router-view />
    </q-page-container>

    <!-- Footer -->
    <q-footer class="bg-grey-2 text-grey-8 q-py-sm">
      <div class="row items-center q-px-md text-center">
        <!-- 1. Kolom Kiri: Dibuat kosong sebagai penyeimbang -->
        <div class="col-4 gt-xs"></div>
        <div class="col-xs-12 col-sm-4 text-weight-bold">
          <div class="q-mt-xs">SDP SMK &copy; 2026. (Made with <span class="text-red-14">❤</span> by. Deden Moh.J.)
          </div>
        </div>
        <q-space />
        <div class="col-xs-12 col-sm-4 text-right text-weight-medium">
          <q-icon name="calendar_today" size="xs" class="q-mr-xs" color="primary" />
          {{ currentDate }}, {{ currentTime }}
        </div>
      </div>
    </q-footer>

    <!-- Search Dialog -->
    <q-dialog v-model="isSearchOpen" persistent>
      <q-card style="width: 500px; max-width: 90vw">
        <q-card-section class="row items-center q-pb-none">
          <div class="text-h6">Cari Siswa</div>
          <q-space />
          <q-btn icon="close" flat round dense @click="closeSearch" />
        </q-card-section>

        <q-card-section>
          <q-input autofocus outlined label="Masukkan NISN, NIS, atau Nama Siswa" @keyup.enter="handleSearch">
            <template v-slot:append>
              <q-icon name="search" />
            </template>
          </q-input>
          <div class="text-caption text-grey-7 q-mt-sm">
            * Fitur pencarian hanya tersedia untuk pengguna yang sudah login
          </div>
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="Batal" color="grey" @click="closeSearch" />
          <q-btn label="Cari" color="primary" @click="handleSearch" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Login Required Dialog -->
    <q-dialog v-model="showLoginDialog" persistent>
      <q-card style="width: 400px; max-width: 90vw">
        <q-card-section class="row items-center q-pb-none">
          <q-avatar icon="lock" color="warning" text-color="white" />
          <div class="text-h6 q-ml-sm">Akses Terbatas</div>
        </q-card-section>

        <q-card-section class="q-pt-md">
          <div class="text-body1">
            Halaman detail hanya dapat diakses oleh pegawai Tata Usaha yang sudah login.
          </div>
          <div class="text-caption text-grey-7 q-mt-sm">
            Silakan login terlebih dahulu untuk melanjutkan.
          </div>
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="Batal" color="grey" @click="cancelLogin" />
          <q-btn label="Login" color="primary" icon="login" @click="proceedToLogin" />
        </q-card-actions>
      </q-card>
    </q-dialog>


    <!-- ✅ REUSABLE LOGIN DIALOG COMPONENT -->
    <LoginDialog v-model="showLoginDialog1"
      message="Halaman ini hanya dapat diakses oleh pegawai Tata Usaha yang sudah login." :is-loading="isLoggingIn"
      :error-message="loginError" @login="handleLoginSubmit" @cancel="handleLoginCancel" />


  </q-layout>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useLandingLayout } from '@/composables/landing/useLandingLayout'
import { useOperationalContext } from '@/composables/context/useOperationalContext'
import { useCurrentDateTime } from '@/composables/ui/useCurrentDateTime'
import { useAuthStore } from '@/stores/authStore'

// Components
import DropDownMenu from '@/components/DropDownMenu.vue'
import ContextSelector from '@/components/ContextSelector.vue'
import LoginDialog from '@/components/LoginDialog.vue'
import TuMenuToolbar from '@/components/dashboard/TuMenuToolbar.vue'

const router = useRouter()
const authStore = useAuthStore()
const { currentDate, currentTime } = useCurrentDateTime()
const { schoolLabel } = useOperationalContext()

const {
  isSearchOpen,
  showLoginRequiredDialog,
  currentUser,
  menuItems,
  itemPortal,
  openSearch,
  closeSearch,
  cancelLogin,
} = useLandingLayout()

// === STATE ===
const leftDrawerOpen = ref(false)
const pendingNavigation = ref(null)
const isLoggingIn = ref(false)
const loginError = ref('')
const showLoginDialog = computed({
  get: () => authStore.showLoginDialog,
  set: (val) => { if (!val) authStore.closeLoginDialog() }
})

const showLoginDialog1 = ref(false)

const handleLoginDialog = () => {
  showLoginDialog1.value = true
}
// === COMPUTED AUTH ===
// Ini akan otomatis berubah (reaktif) jika user login/logout
const isAuthenticated = computed(() => authStore.isAuthenticated)

// === ACTIONS ===
function toggleLeftDrawer() {
  leftDrawerOpen.value = !leftDrawerOpen.value
}

// Fungsi untuk memicu pop-up login dari halaman/komponen lain
function requestAccess(routeName = null, params = {}) {
  if (isAuthenticated.value) {
    if (routeName) router.push({ name: routeName, params })
  } else {
    if (routeName) pendingNavigation.value = { name: routeName, params }
    showLoginDialog.value = true
  }
}

function handleSearch() {
  if (!isAuthenticated.value) {
    closeSearch()
    requestAccess() // Panggil pop-up login
  } else {
    console.log('Search executed')
  }
}

function proceedToLogin() {
  showLoginRequiredDialog.value = false
  showLoginDialog.value = true
}

async function handleLoginSubmit(credentials) {
  isLoggingIn.value = true
  loginError.value = ''

  try {
    const success = await authStore.login(credentials)
    if (success) {
      authStore.closeLoginDialog() // Tutup popup via store
      // Redirect ke halaman yang sebelumnya tertunda, atau ke home
      if (authStore.redirectUrl) {
        router.push(authStore.redirectUrl)
        authStore.redirectUrl = null // Reset
      } else if (pendingNavigation.value) {
        router.push(pendingNavigation.value)
        pendingNavigation.value = null
      } else {
        router.push({ name: 'landing-home' })
      }
    } else {
      loginError.value = authStore.error?.message || 'Username atau password salah.'
    }
  } catch (err) {
    loginError.value = `Terjadi kesalahan jaringan. Silakan coba lagi. ${err}`
  } finally {
    isLoggingIn.value = false
  }
}

function handleLoginCancel() {
  showLoginDialog.value = false
  pendingNavigation.value = null
  loginError.value = ''
}

async function handleLogout() {
  await authStore.logout()
  router.push({ name: 'landing-home' })
}

// Menangkap event emit ketika menu diklik
const onMenuSelected = (data) => {
  if (!isAuthenticated.value) {
    requestAccess() // Jika belum login, munculkan pop-up
  } else {
    // Lakukan aksi routing di sini jika diperlukan
    console.log(`Menu ${data.menu} diklik dengan item:`, data.item)
  }
}

</script>

<style scoped>
.q-tabs {
  font-size: 0.9rem;
}
</style>

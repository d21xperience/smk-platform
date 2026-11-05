<template>
  <q-layout view="hHh lpR fFf">
    <q-header elevated class="bg-white text-dark">
      <!-- ========================================== -->
      <!-- BARIS 1: HEADER UTAMA                      -->
      <!-- ========================================== -->
      <q-toolbar class="q-py-sm">
        <!-- 1. Tombol Hamburger (Hanya muncul di Mobile/Tablet) -->
        <q-btn flat dense round icon="menu" aria-label="Menu" @click="toggleLeftDrawer"
          class="lt-md q-mr-sm text-grey-8" />

        <!-- 2. Brand / Logo Aplikasi -->
        <q-toolbar-title class="flex no-wrap items-center gt-xs">
          <div class="row items-center no-wrap">
            <q-icon name="school" size="1.8em" class="q-mr-sm text-primary" />
            <div class="column">
              <div class="text-h6 text-weight-bold text-primary" style="line-height: 1.1;">SDP</div>
              <div class="text-caption text-grey-7" style="line-height: 1.1;">{{ schoolLabel }}</div>
            </div>
          </div>
        </q-toolbar-title>

        <!-- 3. Menu Utama Horizontal (Hanya muncul di Desktop/Layar Lebar) -->
        <div class="gt-sm row items-center no-wrap text-weight-medium q-ml-lg">
          <q-tabs align="left" active-color="primary" indicator-color="primary" class="text-grey-8">
            <q-route-tab v-for="item in menuItems" :key="item.routeName" :to="{ name: item.routeName }"
              :icon="item.icon" :label="item.name" no-caps />
            <DropDownMenu label="PORTAL" :items="itemPortal" align="center" />
          </q-tabs>
        </div>

        <q-space />

        <!-- 4. Area Kanan (Search, Context, User) -->
        <div class="row items-center no-wrap q-gutter-sm">
          <!-- Tombol Search -->
          <q-btn flat round dense icon="search" @click="openSearch" class="text-grey-8">
            <q-tooltip>Cari siswa</q-tooltip>
          </q-btn>

          <!-- Context Selector (Tahun Ajaran - Semester) -->
          <div class="gt-xs">
            <ContextSelector />
          </div>

          <!-- User Menu / Login -->
          <div v-if="isAuthenticated" class="row items-center q-ml-sm">
            <q-btn flat round>
              <q-avatar color="primary" text-color="white" size="2em" class="text-weight-bold">
                {{ currentUser?.name?.charAt(0).toUpperCase() || 'U' }}
              </q-avatar>
              <q-menu>
                <q-list style="min-width: 150px">
                  <q-item clickable v-close-popup @click="handleLogout">
                    <q-item-section avatar>
                      <q-icon name="logout" color="negative" />
                    </q-item-section>
                    <q-item-section class="text-negative text-weight-medium">Logout</q-item-section>
                  </q-item>
                </q-list>
              </q-menu>
            </q-btn>
          </div>
          <div v-else class="row items-center q-ml-sm">
            <q-btn flat round icon="login" @click="handleLoginDialog" class="text-primary" />
          </div>
        </div>
      </q-toolbar>

      <!-- ========================================== -->
      <!-- BARIS 2: TOOLBAR KEDUA (Ur. Kesiswaan, dll)-->
      <!-- ========================================== -->
      <TuMenuToolbar :visible="isAuthenticated" @menu-click="onMenuSelected" class="gt-xs" />
    </q-header>

    <!-- ========================================== -->
    <!-- SIDENAV (KIRI) - Hanya untuk Mobile        -->
    <!-- ========================================== -->
    <q-drawer v-model="leftDrawerOpen" bordered class="bg-white" :width="280" overlay>
      <q-scroll-area class="fit">
        <!-- Header Sidenav (Mini Profil) -->
        <div class="q-pa-md row items-center bg-grey-1">
          <q-avatar color="primary" text-color="white" icon="school" size="40px" class="q-mr-md" />
          <div>
            <div class="text-weight-bold text-primary">SDP SMK</div>
            <div class="text-caption text-grey-7">{{ schoolLabel }}</div>
          </div>
        </div>
        <q-separator />

        <!-- Menu Navigasi Utama (Karena di Mobile Top Menu disembunyikan, menu ini yang dipakai) -->
        <q-list padding class="menu-list">
          <q-item-label header class="text-weight-bold text-caption text-grey-6 text-uppercase tracking-wide">
            Menu Utama
          </q-item-label>

          <q-item v-for="item in menuItems" :key="item.routeName" clickable v-ripple :to="{ name: item.routeName }"
            active-class="bg-blue-1 text-primary text-weight-bold" class="q-mb-xs rounded-borders"
            @click="leftDrawerOpen = false">
            <q-item-section avatar>
              <q-icon :name="item.icon" :color="item.icon ? 'primary' : 'grey-7'" />
            </q-item-section>
            <q-item-section>
              {{ item.name }}
            </q-item-section>
          </q-item>

          <q-separator class="q-my-md" />

          <!-- Menu Portal -->
          <q-item-label header class="text-weight-bold text-caption text-grey-6 text-uppercase tracking-wide">
            Portal Akses
          </q-item-label>

          <!-- <q-item v-for="(portal, index) in itemPortal" :key="index" clickable v-ripple
            @click="portal.action; leftDrawerOpen = false" class="q-mb-xs rounded-borders">
            <q-item-section avatar>
              <q-icon :name="portal.icon" color="secondary" />
            </q-item-section>
            <q-item-section>
              {{ portal.label }}
            </q-item-section>
          </q-item> -->
        </q-list>
      </q-scroll-area>
    </q-drawer>

    <!-- Konten Utama Router -->
    <q-page-container>
      <router-view />
    </q-page-container>

    <!-- Footer -->
    <q-footer class="bg-grey-2 text-grey-8 q-py-sm">
      <div class="row items-center q-px-md text-center">
        <div class="col-4 gt-xs"></div>
        <div class="col-xs-12 col-sm-4 text-weight-bold">
          <div class="q-mt-xs">SDP Made with <span class="text-red-14">❤</span> by. Deden Moh.J. - &copy; 2026</div>
        </div>
        <q-space />
        <div class="col-xs-12 col-sm-4 text-right text-weight-medium">
          <q-icon name="calendar_today" size="xs" class="q-mr-xs" color="primary" />
          {{ currentDate }}, {{ currentTime }}
        </div>
      </div>
    </q-footer>

    <!-- Dialog Search -->
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

    <!-- Dialog Login Required -->
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

    <!-- Komponen Login Dialog -->
    <LoginDialog v-model="showLoginDialog1"
      message="Halaman ini hanya dapat diakses oleh pegawai Tata Usaha yang sudah login." :is-loading="isLoggingIn"
      :error-message="loginError" @login="handleLoginSubmit" @cancel="handleLoginCancel" />
  </q-layout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useLandingLayout } from '@/composables/landing/useLandingLayout'
import { useOperationalContext } from '@/composables/context/useOperationalContext'
import { useCurrentDateTime } from '@/composables/ui/useCurrentDateTime'
import { useAuthStore } from '@/stores/authStore'
import { useContextStore } from '@/stores/contextStore'

import DropDownMenu from '@/components/DropDownMenu.vue'
import ContextSelector from '@/components/ContextSelector.vue'
import LoginDialog from '@/components/LoginDialog.vue'
import TuMenuToolbar from '@/components/dashboard/TuMenuToolbar.vue'

const router = useRouter()
const authStore = useAuthStore()
const contextStore = useContextStore()

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

const isAuthenticated = computed(() => authStore.isAuthenticated)

onMounted(() => {
  contextStore.initialize()
})

function toggleLeftDrawer() {
  leftDrawerOpen.value = !leftDrawerOpen.value
}

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
    requestAccess()
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
      await contextStore.initialize()

      authStore.closeLoginDialog()
      if (authStore.redirectUrl) {
        router.push(authStore.redirectUrl)
        authStore.redirectUrl = null
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

const onMenuSelected = (data) => {
  if (!isAuthenticated.value) {
    requestAccess()
  } else {
    console.log(`Menu ${data.menu} diklik dengan item:`, data.item)
  }
}
</script>

<style scoped>
.menu-list .q-item {
  border-radius: 8px;
  margin: 4px 12px;
}

.q-item.q-router-link--active {
  color: var(--q-primary);
  background-color: rgba(25, 118, 210, 0.08);
  font-weight: 600;
}

.tracking-wide {
  letter-spacing: 0.05em;
}

/* Memperbaiki tinggi toolbar agar tidak terlalu tinggi di mobile */
@media (max-width: 599px) {
  .q-toolbar {
    min-height: 56px;
  }
}
</style>

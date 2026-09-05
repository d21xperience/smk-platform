<template>
  <q-layout view="hHh lpR fFf">
    <!-- Header -->
    <q-header elevated class="bg-white text-dark">
      <q-toolbar>
        <!-- Mobile menu button -->
        <q-btn flat dense round icon="menu" aria-label="Menu" class="lt-sm" @click="toggleLeftDrawer" />

        <!-- Logo -->
        <q-toolbar-title class="text-weight-bold text-primary">
          <q-icon name="school" size="2em" class="q-mr-sm" />
          SDP SMK
        </q-toolbar-title>

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
        <div v-else class="row items-center q-gutter-sm lt-sm">
          <q-btn flat label="Login" @click="proceedToLogin" />
        </div>
      </q-toolbar>
    </q-header>

    <!-- Mobile Drawer -->
    <!-- <q-drawer v-model="leftDrawerOpen" show-if-above bordered class="bg-grey-1"> -->
    <!-- <q-drawer v-model="leftDrawerOpen" bordered class="bg-grey-1">
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
    </q-drawer> -->

    <!-- Main Content -->
    <q-page-container>
      <router-view />
    </q-page-container>

    <!-- Footer -->
    <q-footer class="bg-grey-14 text-white q-py-sm">
      <div class="text-center text-caption">
        <div class="q-mt-xs">SDP SMK &copy; 2026. (Made with <span class="text-red-14">❤</span> by. Deden M.J.)</div>
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
    <q-dialog v-model="showLoginRequiredDialog" persistent>
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
  </q-layout>
</template>

<script setup>
import { useLandingLayout } from '@/composables/landing/useLandingLayout'
import DropDownMenu from '@/components/DropDownMenu.vue'

const {
  isSearchOpen,
  showLoginRequiredDialog,
  isAuthenticated,
  currentUser,
  menuItems,
  itemPortal,
  toggleLeftDrawer,
  openSearch,
  closeSearch,
  proceedToLogin,
  cancelLogin,
} = useLandingLayout()

function handleSearch() {
  // Cek apakah user sudah login
  if (!isAuthenticated.value) {
    closeSearch()
    showLoginRequiredDialog.value = true
  } else {
    // Implementasi search logic
    console.log('Search executed')
  }
}

function handleLogout() {
  // Implementasi logout logic
  console.log('Logout')
}
</script>

<style scoped>
.q-tabs {
  font-size: 0.9rem;
}
</style>

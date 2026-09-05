<template>
  <q-layout view="lHh Lpr lFf">
    <!-- Header Aplikasi -->
    <q-header elevated class="bg-primary text-white">
      <q-toolbar>
        <q-btn flat dense round icon="menu" aria-label="Menu" @click="toggleLeftDrawer" />
        <q-toolbar-title>
          SARPRAS Digital Hub
          <span class="text-subtitle2 text-weight-light q-ml-sm">v1.0 (Swasta)</span>
        </q-toolbar-title>
        <div>Yayasan Pendidikan</div>
      </q-toolbar>
    </q-header>

    <!-- Sidebar Komponen (q-drawer) -->
    <q-drawer v-model="leftDrawerOpen" show-if-above bordered class="bg-grey-1" :width="280">
      <q-scroll-area class="fit">
        <!-- Profil Singkat Pengguna di Atas Sidebar -->
        <div class="q-pa-md bg-grey-3 text-grey-9 row no-wrap items-center">
          <q-avatar size="48px" class="q-mr-md bg-secondary text-white">TU</q-avatar>
          <div>
            <div class="text-weight-bold">Staf Tata Usaha</div>
            <div class="text-caption text-grey-7">Sarana & Prasarana</div>
          </div>
        </div>

        <!-- Daftar Menu Navigasi -->
        <q-list padding class="text-grey-8">

          <!-- Loop Menu Utama -->
          <template v-for="(menu, index) in menuList" :key="index">

            <!-- Jika Menu Memiliki Sub-menu (Dropdown) -->
            <q-expansion-item v-if="menu.children && menu.children.length > 0" :icon="menu.icon" :label="menu.title"
              header-class="text-weight-medium" group="menu-group">
              <q-list class="q-pl-md">
                <q-item v-for="(subMenu, subIndex) in menu.children" :key="subIndex" clickable v-ripple :to="subMenu.to"
                  active-class="text-primary bg-blue-1 text-weight-bold">
                  <q-item-section avatar v-if="subMenu.icon">
                    <q-icon :name="subMenu.icon" size="xs" />
                  </q-item-section>
                  <q-item-section>
                    {{ subMenu.title }}
                  </q-item-section>
                </q-item>
              </q-list>
            </q-expansion-item>

            <!-- Jika Menu Biasa (Tanpa Sub-menu) -->
            <q-item v-else clickable v-ripple :to="menu.to" exact active-class="text-primary bg-blue-1 text-weight-bold"
              class="text-weight-medium">
              <q-item-section avatar>
                <q-icon :name="menu.icon" />
              </q-item-section>
              <q-item-section>
                {{ menu.title }}
              </q-item-section>
            </q-item>

          </template>

        </q-list>
      </q-scroll-area>
    </q-drawer>

    <!-- Area Konten Halaman Utama -->
    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script setup>
import { ref } from 'vue'

// State untuk membuka/menutup sidebar
const leftDrawerOpen = ref(false)

function toggleLeftDrawer() {
  leftDrawerOpen.value = !leftDrawerOpen.value
}

// Struktur Data Menu Navigasi (Sesuai Kebutuhan Sarpras)
const menuList = ref([
  {
    title: 'Dashboard',
    icon: 'dashboard',
    to: '/sarpras/dashboard'
  },
  {
    title: 'Inventaris & Registrasi',
    icon: 'inventory_2',
    children: [
      { title: 'Registrasi Aset Baru', icon: 'add_box', to: '/sarpras/inventaris/tambah' },
      { title: 'Katalog Barang Habis Pakai', icon: 'Layers', to: '/sarpras/inventaris/logistik' },
      { title: 'Kartu Inventaris Ruang (KIR)', icon: 'room', to: '/sarpras/inventaris/ruangan' },
      { title: 'Cetak QR Code', icon: 'qr_code_2', to: '/sarpras/inventaris/cetak-qr' }
    ]
  },
  {
    title: 'Pemeliharaan & Servis',
    icon: 'build',
    children: [
      { title: 'Laporan Kerusakan', icon: 'report_problem', to: '/sarpras/servis/tiket' },
      { title: 'Jadwal Perawatan Berkala', icon: 'event', to: '/sarpras/servis/jadwal' },
      { title: 'Data Vendor/Teknisi', icon: 'engineering', to: '/sarpras/servis/vendor' }
    ]
  },
  {
    title: 'Logistik & Sirkulasi',
    icon: 'swap_horiz',
    children: [
      { title: 'Peminjaman Barang', icon: 'outbound', to: '/sarpras/sirkulasi/pinjam' },
      { title: 'Permintaan ATK Guru', icon: 'assignment', to: '/sarpras/sirkulasi/minta-atk' }
    ]
  },
  {
    title: 'Mutasi & Penghapusan',
    icon: 'delete_sweep',
    children: [
      { title: 'Log Mutasi Lokasi', icon: 'trending_flat', to: '/sarpras/mutasi/log' },
      { title: 'Usulan Hapus Aset', icon: 'gavel', to: '/sarpras/mutasi/hapus' }
    ]
  },
  {
    title: 'Laporan',
    icon: 'analytics',
    to: '/sarpras/laporan'
  }
])
</script>

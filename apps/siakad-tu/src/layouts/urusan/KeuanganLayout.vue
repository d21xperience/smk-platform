<template>
  <q-layout view="hHh Lpr lFf">
    <!-- Topbar -->
    <q-header elevated>
      <q-toolbar>
        <q-btn flat @click="drawerLeft = !drawerLeft" round dense icon="menu" />
        <q-toolbar-title>Keuangan</q-toolbar-title>

        <!-- Indikator Koneksi Dapodik -->
        <!-- <q-chip :color="dapodikStatus.color" text-color="white" dense icon-right="circle" class="q-mr-sm">
          {{ dapodikStatus.label }}
        </q-chip> -->

        <q-btn flat round dense icon="dark_mode" @click="toggleDarkMode" />
      </q-toolbar>
    </q-header>

    <q-drawer v-model="drawer" show-if-above :mini="miniState" @mouseover="miniState = false"
      @mouseout="miniState = true" :width="250" :breakpoint="700" bordered>
      <q-scroll-area class="fit">
        <q-list padding>
          <q-item clickable v-ripple to="/keuangan/dashboard">
            <q-item-section avatar><q-icon name="dashboard" /></q-item-section>
            <q-item-section>Dashboard</q-item-section>
          </q-item>

          <q-expansion-item label="Master Data" icon="folder" default-opened>
            <q-item clickable v-ripple to="/master/siswa">
              <q-item-section avatar><q-icon name="people" /></q-item-section>
              <q-item-section>Siswa</q-item-section>
            </q-item>
            <q-item clickable v-ripple to="/master/kelas">
              <q-item-section avatar><q-icon name="class" /></q-item-section>
              <q-item-section>Kelas</q-item-section>
            </q-item>
            <q-item clickable v-ripple to="/master/tahun-ajaran">
              <q-item-section avatar><q-icon name="event" /></q-item-section>
              <q-item-section>Tahun Ajaran</q-item-section>
            </q-item>
          </q-expansion-item>

          <q-expansion-item label="Keuangan" icon="account_balance">
            <q-item clickable to="/keuangan/tagihan">
              <q-item-section avatar><q-icon name="receipt" /></q-item-section>
              <q-item-section>Tagihan</q-item-section>
            </q-item>
            <q-item clickable to="/keuangan/pembayaran">
              <q-item-section avatar><q-icon name="payments" /></q-item-section>
              <q-item-section>Pembayaran</q-item-section>
            </q-item>
            <q-item clickable to="/keuangan/pengeluaran">
              <q-item-section avatar><q-icon name="money_off" /></q-item-section>
              <q-item-section>Pengeluaran</q-item-section>
            </q-item>
          </q-expansion-item>
          <!-- <q-expansion-item label="Transaksi" icon="receipt">
          <q-item clickable to="/transaksi/spp">
            <q-item-section avatar><q-icon name="payments" /></q-item-section>
            <q-item-section>Pembayaran SPP</q-item-section>
          </q-item>
          <q-item clickable to="/transaksi/lainnya">
            <q-item-section avatar><q-icon name="add_card" /></q-item-section>
            <q-item-section>Pembayaran Lainnya</q-item-section>
          </q-item>
          <q-item clickable to="/transaksi/pengeluaran">
            <q-item-section avatar><q-icon name="money_off" /></q-item-section>
            <q-item-section>Pengeluaran</q-item-section>
          </q-item>
        </q-expansion-item> -->

          <q-expansion-item label="Laporan" icon="analytics">
            <q-item clickable to="/laporan/pemasukan">
              <q-item-section avatar><q-icon name="trending_up" /></q-item-section>
              <q-item-section>Pemasukan</q-item-section>
            </q-item>
            <q-item clickable to="/laporan/pengeluaran">
              <q-item-section avatar><q-icon name="trending_down" /></q-item-section>
              <q-item-section>Pengeluaran</q-item-section>
            </q-item>
            <q-item clickable to="/laporan/bukukas">
              <q-item-section avatar><q-icon name="book" /></q-item-section>
              <q-item-section>Buku Kas</q-item-section>
            </q-item>
          </q-expansion-item>

          <q-expansion-item label="Pengaturan" icon="settings">
            <q-item clickable to="/pengaturan/user">
              <q-item-section avatar><q-icon name="person" /></q-item-section>
              <q-item-section>User</q-item-section>
            </q-item>
            <q-item clickable to="/pengaturan/role">
              <q-item-section avatar><q-icon name="badge" /></q-item-section>
              <q-item-section>Role</q-item-section>
            </q-item>
          </q-expansion-item>
        </q-list>
      </q-scroll-area>
    </q-drawer>

    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script setup>
import { ref } from 'vue';
import { useQuasar } from 'quasar'

const $q = useQuasar()
const drawer = ref(true);
const miniState = ref(false)

// Dark mode toggle
const isDark = ref(false)
function toggleDarkMode() {
  isDark.value = !isDark.value
  $q.dark.set(isDark.value)
}
</script>

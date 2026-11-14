<template>
  <q-layout view="hHh Lpr lFf">
    <q-header elevated class="bg-teal-9 text-white">
      <q-toolbar>
        <q-btn flat dense round icon="menu" @click="drawer = !drawer" />
        <q-toolbar-title>
          <div class="row items-center no-wrap">
            <q-icon name="account_balance_wallet" size="24px" class="q-mr-sm" />
            <div>
              <div class="text-subtitle1 text-weight-bold">Keuangan SMK</div>
              <div class="text-caption" style="opacity: 0.85">Bendahara &amp; Kepala Sekolah</div>
            </div>
          </div>
        </q-toolbar-title>

        <q-chip v-if="saldoKas !== null" :color="saldoKas < 0 ? 'negative' : 'positive'" text-color="white" dense
          class="q-mr-sm">
          <q-icon name="payments" left size="14px" />
          {{ formatRp(saldoKas) }}
        </q-chip>

        <q-btn flat round dense icon="dark_mode" @click="$q.dark.toggle()">
          <q-tooltip>Toggle dark mode</q-tooltip>
        </q-btn>

        <q-btn flat round dense icon="logout" @click="logout">
          <q-tooltip>Keluar</q-tooltip>
        </q-btn>
      </q-toolbar>
    </q-header>

    <q-drawer v-model="drawer" show-if-above :width="270" :breakpoint="700" bordered class="bg-grey-1">
      <q-scroll-area class="fit">
        <q-list padding>
          <template v-for="(item, idx) in menuItems" :key="idx">
            <q-expansion-item v-if="item.children" :icon="item.icon" :label="item.label"
              :default-opened="isGroupOpen(item)" header-class="text-weight-medium">
              <q-item v-for="(sub, i) in item.children" :key="i" clickable v-ripple :to="sub.to" exact
                active-class="bg-teal-1 text-teal-10 text-weight-bold">
                <q-item-section avatar>
                  <q-icon :name="sub.icon" size="20px" />
                </q-item-section>
                <q-item-section>
                  <q-item-label>{{ sub.label }}</q-item-label>
                </q-item-section>
              </q-item>
            </q-expansion-item>

            <q-item v-else clickable v-ripple :to="item.to" exact
              active-class="bg-teal-1 text-teal-10 text-weight-bold">
              <q-item-section avatar>
                <q-icon :name="item.icon" />
              </q-item-section>
              <q-item-section>
                <q-item-label>{{ item.label }}</q-item-label>
              </q-item-section>
            </q-item>
          </template>
        </q-list>
      </q-scroll-area>
    </q-drawer>

    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useQuasar } from 'quasar'
import { useAuthStore } from '@/stores/authStore'
import { api } from '@/boot/axios'

const $q = useQuasar()
const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const drawer = ref(true)
const saldoKas = ref(null)

// Defensive: ambil dari authStore kalau ada, fallback seed
const currentUser = computed(() => authStore.user || authStore.currentUser || {})
const schoolID = computed(() => currentUser.value.school_id || '11111111-1111-1111-1111-111111111111')
const periodID = computed(() => currentUser.value.academic_period_id || '20261')

const menuItems = [
  { label: 'Dashboard', icon: 'dashboard', to: '/keuangan/dashboard' },
  {
    label: 'Penerimaan & Tagihan', icon: 'receipt_long', children: [
      { label: 'Input Penerimaan', icon: 'add_circle', to: '/keuangan/penerimaan/input' },
      { label: 'Rekap Harian', icon: 'event', to: '/keuangan/penerimaan/rekap-harian' },
      { label: 'Rekap Per Kelas', icon: 'class', to: '/keuangan/penerimaan/rekap-kelas' },
      { label: 'Daftar Tagihan', icon: 'list_alt', to: '/keuangan/tagihan' },
    ]
  },
  {
    label: 'Buku Kas & Pengeluaran', icon: 'book', children: [
      { label: 'Buku Kas Umum', icon: 'menu_book', to: '/keuangan/bku' },
      { label: 'Ajukan Pengeluaran', icon: 'add_shopping_cart', to: '/keuangan/pengeluaran/ajukan' },
      { label: 'Verifikasi Bendahara', icon: 'fact_check', to: '/keuangan/pengeluaran/verifikasi' },
      { label: 'Persetujuan Kepsek', icon: 'approval', to: '/keuangan/pengeluaran/persetujuan' },
      { label: 'Daftar Pengeluaran', icon: 'payments', to: '/keuangan/pengeluaran' },
    ]
  },
  {
    label: 'RAB & Realisasi', icon: 'assessment', children: [
      { label: 'RAB Tahunan', icon: 'edit_document', to: '/keuangan/rab' },
      { label: 'Realisasi Anggaran', icon: 'trending_up', to: '/keuangan/rab/realisasi' },
    ]
  },
  { label: 'Rekonsiliasi Bank', icon: 'sync_alt', to: '/keuangan/rekonsiliasi' },
  { label: 'Buku Piutang', icon: 'receipt', to: '/keuangan/piutang' },
  { label: 'Penutupan Buku Kas', icon: 'lock', to: '/keuangan/penutupan' },
  {
    label: 'Laporan Keuangan', icon: 'analytics', children: [
      { label: 'Neraca', icon: 'account_balance', to: '/keuangan/laporan/neraca' },
      { label: 'LRA', icon: 'summarize', to: '/keuangan/laporan/lra' },
      { label: 'Ringkasan Bulan', icon: 'calendar_month', to: '/keuangan/laporan/ringkasan-bulan' },
    ]
  },
  {
    label: 'Master Data', icon: 'folder', children: [
      { label: 'Jenis Iuran', icon: 'category', to: '/keuangan/master/jenis-iuran' },
      { label: 'Kode Rekening (COA)', icon: 'account_tree', to: '/keuangan/master/kode-rekening' },
      { label: 'Rekening Bank', icon: 'account_balance', to: '/keuangan/master/rekening-bank' },
      { label: 'Tarif Iuran', icon: 'price_change', to: '/keuangan/master/tarif-iuran' },
    ]
  },
  { label: 'PIN TTD Digital', icon: 'fingerprint', to: '/keuangan/pengaturan/pin-ttd' },
]

function isGroupOpen(group) {
  return group.children?.some(c => route.path.startsWith(c.to)) || false
}

function formatRp(n) {
  if (n === null || n === undefined) return '-'
  return 'Rp ' + Number(n).toLocaleString('id-ID')
}

async function logout() {
  try {
    if (typeof authStore.logout === 'function') {
      await authStore.logout()
    }
  } catch { console.log('...') }
  router.push('/')
}

async function loadSaldoKas() {
  try {
    const res = await api.get('/keuangan/laporan/dashboard', {
      params: {
        school_id: schoolID.value,
        academic_period_id: periodID.value
      }
    })
    saldoKas.value = res.data?.data?.kas?.saldo_terkini ?? 0
  } catch {
    // silent — saldo chip tidak wajib
    saldoKas.value = null
  }
}

onMounted(() => {
  loadSaldoKas()
})
</script>

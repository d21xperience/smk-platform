<template>
  <q-layout view="lHh Lpr lFf">

    <!-- ─── HEADER ─────────────────────────────────────────────── -->
    <q-header class="siakad-header">
      <q-toolbar class="q-px-lg">

        <!-- Logo + Brand -->
        <div class="row items-center q-gutter-sm">
          <q-avatar size="34px" color="blue-6" text-color="white" icon="school" />
          <div>
            <div class="text-weight-medium text-white" style="font-size:14px;line-height:1.1">
              SIAKAD SMK Pasundan
            </div>
            <div class="text-blue-3" style="font-size:11px">Sistem Informasi Akademik</div>
          </div>
        </div>

        <q-space />

        <!-- Nav pills (desktop only) -->
        <div class="row q-gutter-xs gt-sm">
          <q-btn v-for="nav in navItems" :key="nav.label" flat dense :label="nav.label" :icon="nav.icon"
            :class="nav.active ? 'nav-pill-active' : 'nav-pill'" no-caps size="sm" />
        </div>

        <q-space />

        <!-- Login button -->
        <q-btn unelevated color="blue-6" label="Login" icon="login" no-caps size="sm" class="q-ml-sm"
          style="border-radius:8px" @click="$router.push('/login')" />

      </q-toolbar>
    </q-header>

    <!-- ─── PAGE ───────────────────────────────────────────────── -->
    <q-page-container>
      <q-page class="siakad-page q-pa-lg">

        <!-- Breadcrumb -->
        <q-breadcrumbs class="q-mb-xs text-grey-5" style="font-size:12px">
          <q-breadcrumbs-el icon="home" />
          <q-breadcrumbs-el label="Dashboard" />
        </q-breadcrumbs>

        <div class="text-h6 text-weight-medium q-mb-lg" style="color:var(--q-dark)">
          Ringkasan Akademik
        </div>

        <!-- ── STAT CARDS ─────────────────────────────────────── -->
        <div class="row q-col-gutter-md q-mb-lg">

          <div v-for="stat in statCards" :key="stat.label" class="col-xs-12 col-sm-6 col-md-3">
            <q-card flat bordered class="stat-card">
              <!-- accent bar -->
              <div class="stat-accent" :style="{ background: stat.color }" />

              <q-card-section class="q-pa-md">
                <q-avatar :color="stat.lightColor" :text-color="stat.iconColor" :icon="stat.icon" size="36px"
                  style="border-radius:8px" class="q-mb-sm" />
                <div class="stat-label">{{ stat.label }}</div>
                <div class="stat-value">{{ stat.value }}</div>
                <div class="stat-sub">
                  <q-icon v-if="stat.trend" name="trending_up" size="12px" color="positive" />
                  <span>{{ stat.sub }}</span>
                </div>
              </q-card-section>
            </q-card>
          </div>

        </div>

        <!-- ── ROW 2: KELAS TABLE + SURAT ────────────────────── -->
        <div class="row q-col-gutter-md q-mb-lg">

          <!-- Kelas table -->
          <div class="col-xs-12 col-md-7">
            <q-card flat bordered class="siakad-card full-height">
              <q-card-section class="q-pb-none">
                <div class="row items-center justify-between">
                  <div class="section-title">
                    <q-icon name="grid_view" size="16px" class="q-mr-xs" />
                    Kelas &amp; Jumlah Siswa
                  </div>
                  <q-btn flat dense no-caps size="xs" label="Lihat semua" color="primary" icon-right="chevron_right" />
                </div>
              </q-card-section>

              <q-table :rows="kelasRows" :columns="kelasColumns" flat dense hide-bottom class="kelas-table q-mt-xs"
                row-key="kelas">
                <template #body-cell-jurusan="props">
                  <q-td :props="props">
                    <q-badge :color="jurusanColor(props.value)" :label="props.value"
                      style="border-radius:20px;font-size:11px;padding:3px 10px" />
                  </q-td>
                </template>
                <template #body-cell-siswa="props">
                  <q-td :props="props" class="text-right">
                    <span class="text-weight-medium">{{ props.value }}</span>
                  </q-td>
                </template>
              </q-table>
            </q-card>
          </div>

          <!-- Surat -->
          <div class="col-xs-12 col-md-5">
            <q-card flat bordered class="siakad-card full-height">
              <q-card-section>
                <div class="section-title q-mb-md">
                  <q-icon name="description" size="16px" class="q-mr-xs" />
                  Surat Keluar Terakhir
                </div>

                <div class="surat-box text-center q-py-sm">
                  <div class="text-caption text-grey-5 q-mb-xs">Nomor surat</div>
                  <div class="surat-number">{{ nomorSuratTerakhir }}</div>
                  <div class="text-caption text-grey-5 q-mt-xs">Surat Keterangan Aktif Siswa</div>
                </div>

                <q-separator class="q-my-sm" />

                <div class="row items-center justify-between" style="font-size:12px">
                  <div class="row items-center q-gutter-xs text-grey-6">
                    <q-icon name="calendar_today" size="12px" />
                    <span>{{ tglSuratTerakhir }}</span>
                  </div>
                  <div class="row items-center q-gutter-xs text-grey-6">
                    <q-icon name="person_outline" size="12px" />
                    <span>TU Admin</span>
                  </div>
                  <q-btn flat dense no-caps size="xs" label="Cetak" icon="print" color="primary" />
                </div>
              </q-card-section>
            </q-card>
          </div>

        </div>

        <!-- ── SEARCH SISWA ───────────────────────────────────── -->
        <q-card flat bordered class="siakad-card q-mb-lg">
          <q-card-section>
            <div class="section-title q-mb-md">
              <q-icon name="manage_search" size="16px" class="q-mr-xs" />
              Cari Siswa
            </div>

            <q-input v-model="searchSiswa" placeholder="Ketik nama atau NIS siswa…" outlined dense debounce="300"
              style="max-width:480px">
              <template #prepend>
                <q-icon name="search" />
              </template>
              <template #append>
                <q-icon v-if="searchSiswa" name="close" class="cursor-pointer" @click="searchSiswa = ''" />
              </template>
            </q-input>

            <!-- Results -->
            <q-table v-if="filteredSiswa.length" :rows="filteredSiswa" :columns="siswaColumns" flat dense
              class="q-mt-md" row-key="nis" hide-bottom />

            <!-- Empty state -->
            <div v-else class="empty-state column items-center q-py-xl text-grey-5">
              <q-icon name="person_search" size="40px" class="q-mb-sm" style="opacity:.35" />
              <div style="font-size:13px">Ketik minimal 3 karakter untuk mulai pencarian</div>
            </div>

          </q-card-section>
        </q-card>

        <!-- Footer -->
        <div class="text-center text-grey-5 q-py-md" style="font-size:11px">
          &copy; {{ new Date().getFullYear() }} SMK Pasundan Jatinangor. All rights reserved.
        </div>

      </q-page>
    </q-page-container>

  </q-layout>
</template>

<script setup>
import { ref, computed } from 'vue'

// ── Nav ──────────────────────────────────────────────────────────
const navItems = [
  { label: 'Dashboard', icon: 'dashboard', active: true },
  { label: 'Siswa', icon: 'people', active: false },
  { label: 'Guru', icon: 'school', active: false },
  { label: 'Kelas', icon: 'grid_view', active: false },
  { label: 'Surat', icon: 'mail', active: false },
]

// ── Stat cards ───────────────────────────────────────────────────
const statCards = [
  {
    label: 'Siswa Aktif',
    value: '842',
    sub: '+12 dari bulan lalu',
    trend: true,
    icon: 'people',
    color: '#3b82f6',
    lightColor: 'blue-1',
    iconColor: 'blue-8',
  },
  {
    label: 'Guru Aktif',
    value: '48',
    sub: 'Semua bidang studi',
    trend: false,
    icon: 'school',
    color: '#10b981',
    lightColor: 'green-1',
    iconColor: 'green-8',
  },
  {
    label: 'Pegawai TU',
    value: '14',
    sub: 'Tenaga kependidikan',
    trend: false,
    icon: 'badge',
    color: '#8b5cf6',
    lightColor: 'purple-1',
    iconColor: 'purple-8',
  },
  {
    label: 'Rata-rata Kehadiran',
    value: '91%',
    sub: 'Semester berjalan',
    trend: false,
    icon: 'event_available',
    color: '#f59e0b',
    lightColor: 'amber-1',
    iconColor: 'amber-8',
  },
]

// ── Kelas table ──────────────────────────────────────────────────
const kelasColumns = [
  { name: 'kelas', label: 'Kelas', field: 'kelas', align: 'left' },
  { name: 'jurusan', label: 'Jurusan', field: 'jurusan', align: 'left' },
  { name: 'waliKelas', label: 'Wali Kelas', field: 'waliKelas', align: 'left' },
  { name: 'siswa', label: 'Siswa', field: 'siswa', align: 'right' },
]

const kelasRows = [
  { kelas: 'X RPL 1', jurusan: 'RPL', waliKelas: 'Pak Andi', siswa: 34 },
  { kelas: 'X RPL 2', jurusan: 'RPL', waliKelas: 'Bu Sari', siswa: 32 },
  { kelas: 'XI TKJ 1', jurusan: 'TKJ', waliKelas: 'Pak Dodi', siswa: 36 },
  { kelas: 'XI TKJ 2', jurusan: 'TKJ', waliKelas: 'Bu Lina', siswa: 33 },
  { kelas: 'XII MM 1', jurusan: 'MM', waliKelas: 'Bu Rina', siswa: 30 },
]

function jurusanColor(jurusan) {
  const map = { RPL: 'blue-6', TKJ: 'teal-6', MM: 'amber-7', AK: 'purple-6' }
  return map[jurusan] ?? 'grey-6'
}

// ── Surat ─────────────────────────────────────────────────────────
const nomorSuratTerakhir = ref('035/SMK-PSD/VI/2025')
const tglSuratTerakhir = ref('28 Mei 2025')

// ── Search siswa ─────────────────────────────────────────────────
const searchSiswa = ref('')

const siswaData = [
  { nis: '2023001', nama: 'Ahmad Fauzi', kelas: 'XI TKJ 1', status: 'Aktif' },
  { nis: '2023002', nama: 'Budi Santoso', kelas: 'X RPL 2', status: 'Aktif' },
  { nis: '2023003', nama: 'Cahya Ramadhani', kelas: 'XII MM 1', status: 'Aktif' },
  { nis: '2023004', nama: 'Dini Permatasari', kelas: 'X RPL 1', status: 'Aktif' },
  { nis: '2023005', nama: 'Eko Prasetyo', kelas: 'XI TKJ 2', status: 'Aktif' },
]

const siswaColumns = [
  { name: 'nis', label: 'NIS', field: 'nis', align: 'left' },
  { name: 'nama', label: 'Nama', field: 'nama', align: 'left' },
  { name: 'kelas', label: 'Kelas', field: 'kelas', align: 'left' },
  { name: 'status', label: 'Status', field: 'status', align: 'left' },
]

const filteredSiswa = computed(() => {
  const q = searchSiswa.value.trim()
  if (q.length < 3) return []
  const lower = q.toLowerCase()
  return siswaData.filter(s =>
    s.nama.toLowerCase().includes(lower) || s.nis.includes(q)
  )
})
</script>

<style scoped>
/* ── Layout ─────────────────────────────────────────────────── */
.siakad-header {
  background: #1e293b !important;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.siakad-page {
  background: #f1f5f9;
  min-height: 100vh;
}

/* ── Nav pills ──────────────────────────────────────────────── */
.nav-pill {
  color: #94a3b8 !important;
  border-radius: 8px !important;
  padding: 4px 12px !important;
}

.nav-pill:hover {
  background: rgba(255, 255, 255, 0.07) !important;
  color: #f1f5f9 !important;
}

.nav-pill-active {
  background: rgba(255, 255, 255, 0.1) !important;
  color: #f1f5f9 !important;
  border-radius: 8px !important;
  padding: 4px 12px !important;
}

/* ── Cards ──────────────────────────────────────────────────── */
.siakad-card {
  background: #ffffff;
  border: 1px solid #e2e8f0 !important;
  border-radius: 12px !important;
}

/* ── Stat card ──────────────────────────────────────────────── */
.stat-card {
  background: #ffffff;
  border: 1px solid #e2e8f0 !important;
  border-radius: 12px !important;
  position: relative;
  overflow: hidden;
  transition: box-shadow 0.2s;
}

.stat-card:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08) !important;
}

.stat-accent {
  position: absolute;
  top: 0;
  left: 0;
  width: 3px;
  height: 100%;
  border-radius: 4px 0 0 4px;
}

.stat-label {
  font-size: 11px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  color: #94a3b8;
  margin-bottom: 4px;
  font-weight: 500;
}

.stat-value {
  font-size: 26px;
  font-weight: 500;
  color: #1e293b;
  line-height: 1;
  margin-bottom: 6px;
}

.stat-sub {
  font-size: 11px;
  color: #b0bec5;
  display: flex;
  align-items: center;
  gap: 3px;
}

/* ── Section title ──────────────────────────────────────────── */
.section-title {
  font-size: 13px;
  font-weight: 500;
  color: #334155;
  display: flex;
  align-items: center;
}

/* ── Kelas table ────────────────────────────────────────────── */
.kelas-table :deep(thead tr th) {
  font-size: 11px !important;
  text-transform: uppercase;
  letter-spacing: 0.4px;
  color: #94a3b8 !important;
  font-weight: 500 !important;
  padding: 8px 12px !important;
  border-bottom: 1px solid #f1f5f9 !important;
}

.kelas-table :deep(tbody tr td) {
  font-size: 13px !important;
  padding: 10px 12px !important;
  border-bottom: 1px solid #f8fafc !important;
}

.kelas-table :deep(tbody tr:hover td) {
  background: #f8fafc !important;
}

/* ── Surat ──────────────────────────────────────────────────── */
.surat-number {
  font-family: 'Courier New', monospace;
  font-size: 16px;
  font-weight: 600;
  color: #3b82f6;
  letter-spacing: 0.3px;
}

/* ── Empty state ────────────────────────────────────────────── */
.empty-state {
  border: 1px dashed #e2e8f0;
  border-radius: 10px;
  margin-top: 12px;
}

/* ── Dark mode support ──────────────────────────────────────── */
.body--dark .siakad-page {
  background: #0f172a;
}

.body--dark .siakad-card,
.body--dark .stat-card {
  background: #1e293b !important;
  border-color: #334155 !important;
}

.body--dark .stat-value {
  color: #f1f5f9;
}

.body--dark .stat-label {
  color: #64748b;
}

.body--dark .section-title {
  color: #cbd5e1;
}

.body--dark .kelas-table :deep(thead tr th) {
  color: #64748b !important;
  border-bottom-color: #1e293b !important;
}

.body--dark .kelas-table :deep(tbody tr td) {
  border-bottom-color: #1e293b !important;
}

.body--dark .kelas-table :deep(tbody tr:hover td) {
  background: #263248 !important;
}

.body--dark .empty-state {
  border-color: #334155;
}
</style>

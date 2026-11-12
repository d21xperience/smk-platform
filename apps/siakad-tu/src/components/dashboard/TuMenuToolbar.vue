<template>
  <q-toolbar v-if="visible" class="gt-xs row items-center no-wrap text-weight-medium dashboard-header">
    <DropDownMenu v-for="(menu, index) in menus" :key="index" :label="menu.label" :items="menu.items"
      class="text-grey-4" @update:item="handleItemClick(menu.label, $event)" />
  </q-toolbar>
</template>

<script setup>
import { computed } from 'vue'
import DropDownMenu from '@/components/DropDownMenu.vue'
import { useContextStore } from '@/stores/contextStore'
const contextStore = useContextStore()

const isPsbActive = computed(() => {
  const current = contextStore.availablePeriods.find(p => String(p.semester_id) === String(contextStore.selectedSemesterId))
  return current?.is_psb_active || false
})


defineProps({ visible: { type: Boolean, default: true } })
const emit = defineEmits(['menu-click'])

const handleItemClick = (menuLabel, itemData) => {
  emit('menu-click', { menu: menuLabel, item: itemData })
}

const itemPengaturan = [
  { title: 'Profil Sekolah', to: '/pengaturan/sekolah' },
  { title: 'Tahun Ajaran & Semester', to: '/pengaturan/tahun-ajaran' },
  { title: 'Manajemen Pengguna', to: '/pengaturan/pengguna' }
]
const itemPersuratan = [
  { title: 'Dashboard', to: '/persuratan/dashboard' },
  { title: 'Surat Keluar', to: '/persuratan/surat-keluar' },
  { title: 'Surat Masuk', to: '/persuratan/surat-masuk' },
  { title: 'Buku Agenda', to: '/persuratan/buku-agenda' },          // <-- BARU
  { title: 'Master Jenis Surat', to: '/persuratan/master-jenis-surat' }, // <-- BARU
  { title: 'Disposisi', to: '/persuratan/disposisi' },
  { title: 'Klasifikasi', to: '/persuratan/klasifikasi' }
]
const itemKepegawaian = [
  { title: 'Dashboard', to: { name: 'dashboard-kepegawaian' } },
  { title: 'Direktori GTK', to: { name: 'direktori-gtk' } },
  { title: 'Buku Induk Pegawai', to: { name: 'buku-induk' } },
  { title: 'Statistik Keberadaan', to: { name: 'statistik-keberadaan' } },
  { title: 'SK Pembagian Tugas', to: { name: 'sk-pembagian-tugas' } },
  { title: 'Penilaian SKP', to: { name: 'skp-periode' } },
  { title: 'Buku Kontrol Kenaikan', to: { name: 'buku-kontrol' } },
  { title: 'Usulan Kenaikan', to: { name: 'usulan-kenaikan' } },
  { title: 'Rekap Kehadiran GTK', to: { name: 'kehadiran-pegawai' } },
  { title: 'Beban Mengajar', to: { name: 'beban-mengajar' } },
  { title: 'Cuti Kerja', to: { name: 'cuti-kerja' } },
  { title: 'Arsip', to: { name: 'arsip-kepegawaian' } },
  { title: 'Laporan Berkala', to: { name: 'laporan-berkala' } },
]

const itemKurikulum = [
  { title: 'Dashboard Kurikulum', to: '/preview/kurikulum' },
  { title: 'Mata Pelajaran', to: '/kurikulum/mapel' },
  { title: 'Penugasan Mengajar', to: '/kurikulum/pembelajaran' },
  { title: 'Kelas', to: '/kurikulum/kelas' },
]

const itemHubin = [
  { title: 'Mitra', to: '/hubin/mitra' },
  { title: 'Penempatan', to: '/hubin/penempatan' },
  { title: 'Logbook', to: '/hubin/logbook' },
  { title: 'Penilaian', to: '/hubin/penilaian' },
  { title: 'Sertifikat', to: '/hubin/sertifikat' },
  { title: 'Ploting Berkas', to: { name: 'ploting-berkas-siswa' } },
  { title: 'Sertifikat UKK', to: { name: 'sertifikat-ukk-siswa' } },
  { title: 'Asuransi Kerja', to: { name: 'asuransi-kerja-siswa' } },
  { title: 'Tracer Study', to: { name: 'tracer-study-alumni' } },
  { title: 'Loker', to: { name: 'loker' } },
]
const itemPsb = [
  { title: 'Tambah Siswa', to: '/psb/siswa/baru' },
  { title: 'Rekap Asal Sekolah', to: '/psb/analisis-asal-sekolah' },
  { title: 'Generate Kelas', to: '/psb/generate-kelas' },
]

const itemKesiswaan = [
  { title: 'Dashboard', to: '/kesiswaan/dashboard' },
  { title: 'Direktori Siswa', to: '/kesiswaan/siswa' },
  { title: 'Buku Induk Siswa', to: '/kesiswaan/buku-induk-siswa' },
  { title: 'Rekap Kehadiran', to: '/kesiswaan/rekap-kehadiran' },
  { title: 'Mutasi Siswa', to: '/kesiswaan/mutasi-siswa' },
  { title: 'Surat Aktif', to: '/kesiswaan/surat-aktif-siswa' },
  { title: 'PIP', to: 'beasiswa-pip-siswa' },
  { title: 'Legalisir', to: '/kesiswaan/legalisir-alumni' },
  { title: 'Pelaporan', to: '/kesiswaan/rekap-pelaporan' },
]
const itemSarpras = [
  { title: 'Dashboard', to: '/sarpras/dashboard' },
  {
    title: 'Inventaris', icon: 'wallet', children: [
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

// const menus = computed(() => [
//   { label: 'Ur. Kesiswaan', items: itemKesiswaan },
//   { label: 'Ur. Kurikulum', items: itemKurikulum },
//   { label: 'Ur. Kepegawaian', items: itemKepegawaian },
//   { label: 'Ur. Sarpras', items: itemSarpras },
//   { label: 'Ur. Hubin', items: itemHubin },
//   { label: 'Persuratan', items: itemPersuratan },
//   { label: 'PSB', items: itemPsb },
//   { label: 'Pengaturan Sistem', items: itemPengaturan }
// ])
const menus = computed(() => {
  const list = [
    { label: 'Ur. Kesiswaan', items: itemKesiswaan },
    { label: 'Ur. Kurikulum', items: itemKurikulum },
    { label: 'Ur. Kepegawaian', items: itemKepegawaian },
    { label: 'Ur. Sarpras', items: itemSarpras },
    { label: 'Ur. Hubin', items: itemHubin },
    { label: 'Persuratan', items: itemPersuratan },
  ]

  // Hanya tampilkan PSB jika aktif
  if (isPsbActive.value) {
    list.push({ label: 'PSB', items: itemPsb })
  }

  list.push({ label: 'Pengaturan Sistem', items: itemPengaturan })
  return list
})
</script>

<style scoped>
.dashboard-header {
  background: linear-gradient(135deg, #0a192f 0%, #123a63 100%);
}
</style>

<template>
  <q-toolbar v-show="visible" class="gt-xs row items-center no-wrap text-weight-medium dashboard-header">
    <!-- Melakukan looping untuk setiap item menu yang dikirim lewat props -->
    <DropDownMenu v-for="(menu, index) in menus" :key="index" :label="menu.label" :items="menu.items"
      class="text-grey-4" @update:item="handleItemClick(menu.label, $event)" />
  </q-toolbar>
</template>

<script setup>
import { defineProps, defineEmits, computed } from 'vue'
import DropDownMenu from '@/components/DropDownMenu.vue' // Sesuaikan jalur folder Anda

// 1. Definisikan Props
defineProps({
  visible: {
    type: Boolean,
    default: true
  },
})

// 2. Definisikan Emit
const emit = defineEmits(['menu-click'])

// 3. Fungsi untuk menangkap klik dan mengirimnya ke komponen induk
const handleItemClick = (menuLabel, itemData) => {
  emit('menu-click', {
    menu: menuLabel,
    item: itemData
  })
}

const itemPersuratan = [
  { title: 'Dashboard', to: '/persuratan/dashboard' },
  { title: 'Surat Keluar', to: '/persuratan/surat-keluar' },
  { title: 'Surat Masuk', to: '/persuratan/surat-masuk' },
  { title: 'Disposisi', to: '/persuratan/disposisi' },
  { title: 'Klasifikasi', to: '/persuratan/klasifikasi' }
]
const itemKepegawaian = [
  { title: 'Dashboard', to: { name: 'dashboard-kepegawaian' } },
  { title: 'Direktori GTK', to: { name: 'direktori-gtk' } },
  { title: 'Kinerja GTK', to: { name: '' } },
  { title: 'Beban Mengajar', to: { name: 'beban-mengajar' } },
  { title: 'Cuti Kerja', to: { name: 'cuti-kerja' } },
  { title: 'Arsip', to: { name: 'arsip-kepegawaian' } },
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
  { title: 'Kelas', to: '/psb/generate-kelas' },
]

const itemKesiswaan = [
  { title: 'Dashboard', to: '/kesiswaan/dashboard' },
  { title: 'Direktori Siswa', to: '/kesiswaan/siswa' },
  { title: 'Buku Induk Siswa', to: '/kesiswaan/buku-induk-siswa' },
  { title: 'Mutasi Siswa', to: { name: 'student-mutation' } },
  { title: 'Surat Aktif', to: '/kesiswaan/surat-aktif-siswa' },
  { title: 'PIP', to: { name: 'beasiswa-pip-siswa' } },
  { title: 'Legalisir', to: '/kesiswaan/legalisir-alumni' },
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

// Menyatukan data menjadi satu array terstruktur untuk dikirim ke props
const menus = computed(() => [
  { label: 'Ur. Kesiswaan', items: itemKesiswaan },
  { label: 'Ur. Kepegawaian', items: itemKepegawaian },
  { label: 'Ur. Sarpras', items: itemSarpras },
  { label: 'Ur. Hubin', items: itemHubin },
  { label: 'Persuratan', items: itemPersuratan },
  { label: 'PSB', items: itemPsb },
  // { label: 'portal', items: itemPortal }
])
</script>

<style scoped>
.dashboard-header {
  background: linear-gradient(135deg, #0a192f 0%, #123a63 100%);
}
</style>

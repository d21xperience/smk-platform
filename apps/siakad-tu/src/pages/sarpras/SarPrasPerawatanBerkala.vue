<template>
  <q-page class="q-pa-md">
    <!-- PANEL ACTIONS & STATISTIK JADWAL -->
    <div class="row q-col-gutter-sm q-mb-md items-center">
      <div class="col-12 col-sm-4">
        <q-card flat bordered class="bg-indigo-1 text-indigo-9">
          <q-card-section class="row items-center no-wrap q-py-sm">
            <div>
              <div class="text-caption text-grey-7">Total Agenda Perawatan</div>
              <div class="text-h6 text-weight-bold">{{ daftarJadwal.length }} Rutinitas</div>
            </div>
            <q-space />
            <q-icon name="date_range" size="sm" />
          </q-card-section>
        </q-card>
      </div>

      <div class="col-12 col-sm-4">
        <q-card flat bordered class="bg-deep-orange-1 text-deep-orange-9">
          <q-card-section class="row items-center no-wrap q-py-sm">
            <div>
              <div class="text-caption text-grey-7">Jatuh Tempo Bulan Ini</div>
              <div class="text-h6 text-weight-bold">{{ totalJatuhTempo }} Aset</div>
            </div>
            <q-space />
            <q-icon name="notification_important" size="sm" />
          </q-card-section>
        </q-card>
      </div>

      <div class="col-12 col-sm-4 text-right">
        <q-btn color="primary" icon="add_alarm" label="Tambah Jadwal Servis" class="full-width-xs"
          @click="bukaDialogTambah" />
      </div>
    </div>

    <!-- TABEL MANAJEMEN AGENDA PERAWATAN BERKALA -->
    <q-table title="Jadwal Perawatan & Servis Rutin" :rows="daftarJadwal" :columns="columns" row-key="id"
      :filter="filterText" flat bordered>
      <!-- Slot Cari Data -->
      <template v-slot:top-right>
        <q-input borderless dense debounce="300" v-model="filterText" placeholder="Cari nama aset/ruang...">
          <template v-slot:append>
            <q-icon name="search" />
          </template>
        </q-input>
      </template>

      <!-- Kustomisasi Tampilan Kolom Frekuensi Rutinitas -->
      <template v-slot:body-cell-frekuensi="props">
        <q-td :props="props">
          <q-badge :color="getFrekuensiColor(props.value)" class="text-weight-bold">
            {{ props.value }}
          </q-badge>
        </q-td>
      </template>

      <!-- Kustomisasi Tampilan Tanggal Servis Berikutnya -->
      <template v-slot:body-cell-nextServis="props">
        <q-td :props="props" class="text-center text-weight-medium">
          <div class="row items-center justify-center no-wrap">
            <q-icon name="event" size="xs" class="q-mr-xs text-grey-6" />
            <span>{{ props.value }}</span>
            <!-- Deteksi jika sudah mendekati atau lewat tanggal target -->
            <q-badge v-if="apakahJatuhTempo(props.value)" color="deep-orange" text-color="white" floating rounded />
          </div>
        </q-td>
      </template>

      <!-- Kolom Tombol Aksi Cepat Selesai Servis Rutin -->
      <template v-slot:body-cell-aksi="props">
        <q-td :props="props" class="text-center">
          <q-btn dense size="sm" color="positive" icon="verified" label="Selesai Servis"
            @click="konfirmasiSelesaiServis(props.row)">
            <q-tooltip>Catat bahwa servis rutin periode ini sudah beres</q-tooltip>
          </q-btn>
        </q-td>
      </template>
    </q-table>
  </q-page>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const filterText = ref('')

// 1. Mock Data Agenda Perawatan Aset Sekolah
const daftarJadwal = ref([
  { id: 'PM-01', namaAset: 'AC Split Sharp 1 PK (Semua Ruang Kelas)', lokasi: 'Gedung A & B', frekuensi: '3 Bulanan', lastServis: '2026/03/10', nextServis: '2026/06/10' },
  { id: 'PM-02', namaAset: 'Deep Cleaning & Defregment PC Lab', lokasi: 'Lab Komputer', frekuensi: 'Bulanan', lastServis: '2026/05/01', nextServis: '2026/06/01' },
  { id: 'PM-03', namaAset: 'Servis Rutin & Ganti Oli Mobil Sekolah', lokasi: 'Garasi Samping', frekuensi: '6 Bulanan', lastServis: '2026/01/15', nextServis: '2026/07/15' },
  { id: 'PM-04', namaAset: 'Kalibrasi Mikroskop Lab Biologi', lokasi: 'Lab IPA', frekuensi: 'Tahunan', lastServis: '2025/11/20', nextServis: '2026/11/20' }
])

// 2. Definisi Struktur Kolom Tabel Quasar
const columns = [
  { name: 'id', label: 'ID Jadwal', align: 'left', field: 'id', sortable: true },
  { name: 'namaAset', label: 'Nama Program Perawatan / Nama Aset', align: 'left', field: 'namaAset', sortable: true },
  { name: 'lokasi', label: 'Area Cakupan', align: 'left', field: 'lokasi' },
  { name: 'frekuensi', label: 'Siklus Berkala', align: 'center', field: 'frekuensi' },
  { name: 'lastServis', label: 'Servis Terakhir', align: 'center', field: 'lastServis' },
  { name: 'nextServis', label: 'Jadwal Berikutnya', align: 'center', field: 'nextServis', sortable: true },
  { name: 'aksi', label: 'Status Eksekusi', align: 'center' }
]

// 3. Helper Penentuan Warna Siklus Berkala
function getFrekuensiColor(frekuensi) {
  if (frekuensi === 'Bulanan') return 'teal'
  if (frekuensi === '3 Bulanan') return 'indigo'
  if (frekuensi === '6 Bulanan') return 'blue-grey'
  return 'purple'
}

// 4. Logika Perhitungan Jatuh Tempo (Asumsi Tanggal Hari Ini: 06 Juni 2026)
const TANGGAL_SEKARANG = new Date('2026/06/06')

function apakahJatuhTempo(tanggalTargetTeks) {
  const targetDate = new Date(tanggalTargetTeks)
  // Dianggap kritis jika tanggal target sudah terlewat atau tersisa kurang dari 5 hari lagi
  const selisihHari = (targetDate - TANGGAL_SEKARANG) / (1000 * 60 * 60 * 24)
  return selisihHari <= 5
}

const totalJatuhTempo = computed(() => {
  return daftarJadwal.value.filter(item => apakahJatuhTempo(item.nextServis)).length
})

// 5. Aksi Penanganan Klik Selesai Servis
function konfirmasiSelesaiServis(row) {
  $q.dialog({
    title: 'Konfirmasi Selesai Perawatan',
    message: `Apakah Anda ingin memperbarui tanggal pengerjaan rutin untuk "${row.namaAset}"? Tanggal pengerjaan terakhir akan diubah menjadi hari ini.`,
    cancel: true,
    persistent: true
  }).onOk(() => {
    // Simulasi pembaruan tanggal otomatis
    row.lastServis = '2026/06/06'

    // Geser tanggal berikutnya berdasarkan tipe frekuensi
    if (row.frekuensi === 'Bulanan') row.nextServis = '2026/07/06'
    else if (row.frekuensi === '3 Bulanan') row.nextServis = '2026/09/06'
    else if (row.frekuensi === '6 Bulanan') row.nextServis = '2026/12/06'

    $q.notify({
      color: 'success',
      icon: 'check_circle',
      message: 'Jurnal jadwal pemeliharaan berhasil dimutakhirkan!',
      backgroundColor: 'green-6'
    })
  })
}

function bukaDialogTambah() {
  $q.notify({
    message: 'Membuka formulir pembuatan siklus perawatan berkala baru.',
    color: 'primary'
  })
}
</script>

<style scoped>
@media (max-width: 599px) {
  .full-width-xs {
    width: 100%;
  }
}
</style>

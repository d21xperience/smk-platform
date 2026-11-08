<template>
  <q-page class="q-pa-md bg-grey-2">
    <!-- Header Halaman -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Manajemen Asuransi Kerja PKL</h1>
        <p class="text-caption text-grey-7 q-mb-none">Administrasi perlindungan jaminan kecelakaan dan keselamatan kerja
          siswa magang SMK</p>
      </div>
      <div>
        <q-btn color="green-8" icon="how_to_reg" label="Aktivasi Massal" @click="bulkActivation" />
      </div>
    </div>

    <!-- Ringkasan Statistik Perlindungan Siswa -->
    <div class="row q-col-gutter-sm q-mb-md">
      <div class="col-12 col-sm-4">
        <q-card class="bg-teal-1 text-teal-10 flat shadow-1">
          <q-card-section class="q-py-sm flex justify-between items-center">
            <div>
              <div class="text-caption text-weight-medium">Terlindungi Aktif</div>
              <div class="text-h5 text-weight-bold">330 / 342 Siswa</div>
            </div>
            <q-icon name="verified_user" size="md" />
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-sm-4">
        <q-card class="bg-red-1 text-red-9 flat shadow-1">
          <q-card-section class="q-py-sm flex justify-between items-center">
            <div>
              <div class="text-caption text-weight-medium">Belum Aktif / Pending</div>
              <div class="text-h5 text-weight-bold">12 Siswa</div>
            </div>
            <q-icon name="gpp_maybe" size="md" />
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-sm-4">
        <q-card class="bg-purple-1 text-purple-9 flat shadow-1">
          <q-card-section class="q-py-sm flex justify-between items-center">
            <div>
              <div class="text-caption text-weight-medium">Total Klaim Aktif</div>
              <div class="text-h5 text-weight-bold">1 Kasus</div>
            </div>
            <q-icon name="medical_services" size="md" />
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Panel Filter & Tabel Utama -->
    <q-card class="shadow-1">
      <q-card-section class="q-py-sm">
        <div class="row q-col-gutter-sm items-center">
          <div class="col-12 col-md-5">
            <q-input v-model="search" dense outlined placeholder="Cari Nama Siswa, Perusahaan, atau No Polis..."
              clearable>
              <template v-slot:prepend><q-icon name="search" /></template>
            </q-input>
          </div>
          <div class="col-12 col-md-4">
            <q-select v-model="filterStatus"
              :options="['Semua Status', 'Aktif Dilindungi', 'Masa Aktif Habis', 'Belum Terdaftar']" dense outlined
              label="Status Asuransi" />
          </div>
          <div class="col-12 col-md-3 flex justify-end text-caption text-grey-7">
            Provider Utama: <strong class="text-primary q-ml-xs">BPJS Ketenagakerjaan</strong>
          </div>
        </div>
      </q-card-section>

      <!-- Tabel Jaminan Keselamatan Kerja -->
      <q-card-section class="q-pa-none">
        <q-table flat :rows="filteredAsuransi" :columns="columns" row-key="id" selection="multiple"
          v-model:selected="selectedSiswa" :pagination="{ rowsPerPage: 10 }">
          <!-- Custom Tampilan Nomor Polis & Provider -->
          <template v-slot:body-cell-polis="props">
            <q-td :props="props">
              <div class="text-weight-bold text-grey-9">{{ props.row.noPolis || '—' }}</div>
              <div class="text-caption text-grey-6">{{ props.row.provider }}</div>
            </q-td>
          </template>

          <!-- Custom Badge Status Perlindungan -->
          <template v-slot:body-cell-status="props">
            <q-td :props="props" class="text-center">
              <q-badge
                :color="props.row.status === 'Aktif Dilindungi' ? 'teal' : props.row.status === 'Masa Aktif Habis' ? 'amber-9' : 'red'"
                :label="props.row.status" />
            </q-td>
          </template>

          <!-- Custom Kolom Aksi Perlindungan & Penanganan Klaim Insiden -->
          <template v-slot:body-cell-aksi="props">
            <q-td :props="props" class="text-center q-gutter-xs">
              <q-btn size="sm" color="primary" icon="badge" round flat @click="openEditPolisDialog(props.row)">
                <q-tooltip>Input / Ubah Nomor Polis</q-tooltip>
              </q-btn>
              <q-btn size="sm" color="red-8" icon="gavel" round flat :disable="props.row.status === 'Belum Terdaftar'"
                @click="triggerClaim(props.row)">
                <q-tooltip>Ajukan Klaim JKK (Kecelakaan Kerja)</q-tooltip>
              </q-btn>
            </q-td>
          </template>
        </q-table>
      </q-card-section>
    </q-card>

    <!-- DIALOG POPUP: Ubah/Input Nomor Polis Asuransi -->
    <q-dialog v-model="dialogPolisOpen" persistent>
      <q-card style="min-width: 380px;">
        <q-card-section class="bg-primary text-white">
          <div class="text-h6">Update Nomor Polis Siswa</div>
          <div class="text-caption text-blue-2">{{ activeRow.nama }} ({{ activeRow.kelas }})</div>
        </q-card-section>

        <q-card-section class="q-gutter-md q-pt-md">
          <q-input v-model="formPolis.noPolis" label="Nomor Kartu / Polis Asuransi" outlined dense
            placeholder="Masukkan nomor registrasi polis..." />
          <q-input v-model="formPolis.masaAktif" label="Masa Berlaku Hingga" outlined dense
            placeholder="Contoh: 31 Desember 2026" />
          <q-select v-model="formPolis.status" :options="['Aktif Dilindungi', 'Masa Aktif Habis', 'Belum Terdaftar']"
            label="Status Kepesertaan" outlined dense />
        </q-card-section>

        <q-card-actions align="right" class="q-pb-md q-pr-md">
          <q-btn flat label="Batal" v-close-popup color="grey" />
          <q-btn color="primary" label="Simpan Polis" @click="savePolisInfo" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()

// State Komponen
const search = ref('')
const filterStatus = ref('Semua Status')
const selectedSiswa = ref([])
const dialogPolisOpen = ref(false)

const activeRow = ref({})
const formPolis = ref({ noPolis: '', masaAktif: '', status: '' })

// Definisi Judul Kolom Tabel q-table
const columns = [
  { name: 'nama', align: 'left', label: 'Nama Siswa', field: 'nama', sortable: true },
  { name: 'perusahaan', align: 'left', label: 'Lokasi PKL (DUDI)', field: 'perusahaan', sortable: true },
  { name: 'polis', align: 'left', label: 'Informasi Jaminan Keamanan', field: 'noPolis' },
  { name: 'masaAktif', align: 'left', label: 'Masa Berlaku', field: 'masaAktif' },
  { name: 'status', align: 'center', label: 'Status JKK', field: 'status', sortable: true },
  { name: 'aksi', align: 'center', label: 'Aksi & Klaim' }
]

// Data Mock / Dummy Kepesertaan Jaminan Kecelakaan Kerja (JKK) Magang
const dataAsuransi = ref([
  { id: 1, nama: 'Ahmad Fauzi', kelas: 'XII RPL 1', perusahaan: 'PT. Telekomunikasi Indonesia Tbk', provider: 'BPJS Ketenagakerjaan', noPolis: '000123456789', masaAktif: '30 Sep 2026', status: 'Aktif Dilindungi' },
  { id: 2, nama: 'Budi Santoso', kelas: 'XII TKRO 3', perusahaan: 'PT. Astra Honda Motor', provider: 'BPJS Ketenagakerjaan', noPolis: '000123456790', masaAktif: '30 Sep 2026', status: 'Aktif Dilindungi' },
  { id: 3, nama: 'Chandra Wijaya', kelas: 'XI RPL 2', perusahaan: 'PT. Solusi Teknologi Nusantara', provider: 'BPJS Ketenagakerjaan', noPolis: '000123456112', masaAktif: '14 Juli 2026', status: 'Masa Aktif Habis' },
  { id: 4, nama: 'Dinda Lestari', kelas: 'XI AKL 1', perusahaan: 'Bank Mandiri (Persero) Tbk', provider: 'Belum Terdaftar', noPolis: '', masaAktif: '—', status: 'Belum Terdaftar' }
])

// Logika Compute Penyaringan Filter Frontend
const filteredAsuransi = computed(() => {
  return dataAsuransi.value.filter(item => {
    const matchSearch = search.value
      ? item.nama.toLowerCase().includes(search.value.toLowerCase()) ||
      item.perusahaan.toLowerCase().includes(search.value.toLowerCase()) ||
      item.noPolis.includes(search.value)
      : true

    const matchStatus = filterStatus.value !== 'Semua Status'
      ? item.status === filterStatus.value
      : true

    return matchSearch && matchStatus
  })
})

// Fungsi Interaksi Tombol Administrasi TU
const openEditPolisDialog = (row) => {
  activeRow.value = row
  formPolis.value = {
    noPolis: row.noPolis,
    masaAktif: row.masaAktif === '—' ? '30 September 2026' : row.masaAktif,
    status: row.status
  }
  dialogPolisOpen.value = true
}

const savePolisInfo = () => {
  const target = dataAsuransi.value.find(item => item.id === activeRow.value.id)
  if (target) {
    target.noPolis = formPolis.value.noPolis
    target.masaAktif = formPolis.value.masaAktif
    target.status = formPolis.value.status
    if (target.noPolis && target.provider === 'Belum Terdaftar') {
      target.provider = 'BPJS Ketenagakerjaan'
    }
  }
  $q.notify({ color: 'primary', message: `Data kartu jaminan untuk ${activeRow.value.nama} sukses disimpan!`, icon: 'save' })
  dialogPolisOpen.value = false
}

const triggerClaim = (row) => {
  $q.dialog({
    title: 'Konfirmasi Pengajuan Klaim JKK',
    message: `Apakah Anda ingin mencetak dokumen berkas F3-BPJS untuk memproses klaim pengobatan kecelakaan kerja atas nama siswa ${row.nama}?`,
    cancel: true,
    persistent: true
  }).onOk(() => {
    $q.notify({ color: 'purple-8', message: `Berkas klaim kecelakaan berhasil digenerate. Silakan kirimkan ke RS instansi terafiliasi.`, icon: 'medical_services' })
  })
}

const bulkActivation = () => {
  if (selectedSiswa.value.length === 0) {
    $q.notify({ color: 'negative', message: 'Silakan pilih/centang nama siswa yang ingin diaktivasi preminya!', icon: 'warning' })
    return
  }
  selectedSiswa.value.forEach(siswa => {
    const target = dataAsuransi.value.find(item => item.id === siswa.id)
    if (target && target.status === 'Belum Terdaftar') {
      target.status = 'Aktif Dilindungi'
      target.noPolis = 'BPJS-MAGANG-' + Math.floor(100000 + Math.random() * 900000)
      target.provider = 'BPJS Ketenagakerjaan'
      target.masaAktif = '30 September 2026'
    }
  })
  $q.notify({ color: 'green-8', message: `Berhasil mendaftarkan jaminan aktif masal untuk ${selectedSiswa.value.length} siswa terpilih ke sistem BPJS Ketenagakerjaan.`, icon: 'done_all' })
  selectedSiswa.value = []
}
</script>

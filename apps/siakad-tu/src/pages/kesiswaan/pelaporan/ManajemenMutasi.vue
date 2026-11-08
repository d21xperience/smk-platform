<template>
  <q-page class="q-pa-md bg-grey-2">
    <!-- Header Halaman -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Manajemen Mutasi Siswa</h1>
        <p class="text-caption text-grey-7 q-mb-none">Pengelolaan administrasi keluar-masuk peserta didik SMK</p>
      </div>
      <div>
        <q-btn :color="activeTab === 'masuk' ? 'green-8' : 'red-9'" :icon="activeTab === 'masuk' ? 'login' : 'logout'"
          :label="activeTab === 'masuk' ? 'Catat Mutasi Masuk' : 'Catat Mutasi Keluar'" @click="openMutasiDialog" />
      </div>
    </div>

    <!-- Tab Filter Alur Mutasi -->
    <q-card class="shadow-1 q-mb-md">
      <q-tabs v-model="activeTab" dense class="text-grey" active-color="primary" indicator-color="primary" align="left"
        narrow-indicator>
        <q-tab name="masuk" icon="call_received" label="Mutasi Masuk (Pindahan)" />
        <q-tab name="keluar" icon="call_made" label="Mutasi Keluar (Pindah Sekolah)" />
      </q-tabs>
    </q-card>

    <!-- Tabel Riwayat Mutasi -->
    <q-card class="shadow-1">
      <q-card-section class="q-pa-none">
        <q-table flat :rows="activeTab === 'masuk' ? riwayatMasuk : riwayatKeluar"
          :columns="activeTab === 'masuk' ? kolomMasuk : kolomKeluar" row-key="id" :pagination="{ rowsPerPage: 10 }">
          <!-- Kustomisasi Kolom Berkas / Status -->
          <template v-slot:body-cell-berkas="props">
            <q-td :props="props" class="text-center">
              <q-chip :color="props.row.berkasLengkap ? 'green-1' : 'amber-1'"
                :text-color="props.row.berkasLengkap ? 'green-9' : 'amber-9'" size="sm" dense
                :icon="props.row.berkasLengkap ? 'check_circle' : 'pending'">
                {{ props.row.berkasLengkap ? 'Lengkap' : 'Belum Lengkap' }}
              </q-chip>
            </q-td>
          </template>

          <!-- Kustomisasi Kolom Aksi Cetak Dokumen -->
          <template v-slot:body-cell-aksi="props">
            <q-td :props="props" class="text-center">
              <q-btn size="sm" color="primary" icon="print" round flat @click="cetakSuratMutasi(props.row)">
                <q-tooltip>Cetak Surat Keterangan / Arsip</q-tooltip>
              </q-btn>
            </q-td>
          </template>
        </q-table>
      </q-card-section>
    </q-card>

    <!-- DIALOG FORM INPUT MUTASI BARU -->
    <q-dialog v-model="dialogOpen" persistent>
      <q-card style="min-width: 500px; max-width: 600px;">
        <q-card-section class="row items-center q-pb-none">
          <div class="text-h6 text-weight-bold">
            {{ activeTab === 'masuk' ? 'Formulir Siswa Masuk Pindahan' : 'Formulir Siswa Pindah Keluar' }}
          </div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-separator class="q-mt-sm" />

        <q-card-section class="q-scroll q-py-md" style="max-height: 70vh;">
          <q-form @submit.prevent="simpanMutasi" class="row q-col-gutter-sm">

            <!-- Jika MUTASI MASUK (Input Data Baru) -->
            <template v-if="activeTab === 'masuk'">
              <div class="col-12"><q-input v-model="form.nama" label="Nama Lengkap Siswa" outlined dense required />
              </div>
              <div class="col-6"><q-input v-model="form.nisn" label="NISN" outlined dense required type="number" />
              </div>
              <div class="col-6">
                <q-select v-model="form.jurusan" :options="opsiJurusan" label="Pilihan Jurusan SMK" outlined dense
                  required />
              </div>
              <div class="col-12"><q-input v-model="form.sekolahAsal" label="Sekolah Asal (SMK/SMA Asal)" outlined dense
                  required /></div>
              <div class="col-12"><q-input v-model="form.noSuratRekomendasi" label="No. Surat Rekomendasi Dinas"
                  outlined dense /></div>
              <div class="col-12">
                <q-checkbox v-model="form.berkasLengkap"
                  label="Seluruh berkas fisik/rapor mutasi sudah diverifikasi lengkap" color="green" />
              </div>
            </template>

            <!-- Jika MUTASI KELUAR (Pilih dari data yang sudah ada) -->
            <template v-else>
              <div class="col-12">
                <q-input v-model="form.nama" label="Nama Siswa Aktif" placeholder="Ketik nama atau NISN siswa..."
                  outlined dense required />
                <p class="text-caption text-grey-6 q-mt-xs q-mb-none">*Sistem otomatis menonaktifkan status di Dapodik
                  setelah disimpan.</p>
              </div>
              <div class="col-12"><q-input v-model="form.sekolahTujuan" label="Sekolah Tujuan (Tempat Pindah)" outlined
                  dense required /></div>
              <div class="col-12"><q-input v-model="form.alasan" label="Alasan Pindah / Mutasi Keluar" outlined dense
                  type="textarea" rows="2" required /></div>
            </template>

            <!-- Tombol Submit -->
            <div class="col-12 flex justify-end q-mt-md">
              <q-btn label="Batal" color="grey" flat v-close-popup class="q-mr-sm" />
              <q-btn :label="activeTab === 'masuk' ? 'Simpan & Daftarkan' : 'Proses Keluar'"
                :color="activeTab === 'masuk' ? 'green-8' : 'red-9'" type="submit" />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()

// State Navigasi Tab Aktif ('masuk' atau 'keluar')
const activeTab = ref('masuk')

// State Mengontrol Dialog Form
const dialogOpen = ref(false)

// Struktur Reactive Form Object
const form = ref({
  nama: '', nisn: '', jurusan: '', sekolahAsal: '', sekolahTujuan: '',
  noSuratRekomendasi: '', alasan: '', berkasLengkap: false
})

const opsiJurusan = ['RPL', 'TKRO', 'AKL']

// 1. Defenisi Kolom & Data Mock MUTASI MASUK
const kolomMasuk = [
  { name: 'tgl', align: 'left', label: 'Tanggal Masuk', field: 'tgl', sortable: true },
  { name: 'nama', align: 'left', label: 'Nama Siswa', field: 'nama' },
  { name: 'nisn', align: 'left', label: 'NISN', field: 'nisn' },
  { name: 'jurusan', align: 'center', label: 'Jurusan', field: 'jurusan' },
  { name: 'sekolahAsal', align: 'left', label: 'Sekolah Asal', field: 'sekolahAsal' },
  { name: 'berkas', align: 'center', label: 'Status Berkas' },
  { name: 'aksi', align: 'center', label: 'Cetak' }
]

const riwayatMasuk = ref([
  { id: 101, tgl: '10/02/2026', nama: 'Rian Hidayat', nisn: '0078811223', jurusan: 'RPL', sekolahAsal: 'SMKN 2 Bandung', berkasLengkap: true },
  { id: 102, tgl: '22/01/2026', nama: 'Siti Sarah', nisn: '0075544332', jurusan: 'AKL', sekolahAsal: 'SMA Pasundan 1', berkasLengkap: false }
])

// 2. Defenisi Kolom & Data Mock MUTASI KELUAR
const kolomKeluar = [
  { name: 'tgl', align: 'left', label: 'Tanggal Keluar', field: 'tgl', sortable: true },
  { name: 'nama', align: 'left', label: 'Nama Siswa', field: 'nama' },
  { name: 'kelasAsal', align: 'center', label: 'Kelas Asal', field: 'kelasAsal' },
  { name: 'sekolahTujuan', align: 'left', label: 'Sekolah Tujuan', field: 'sekolahTujuan' },
  { name: 'alasan', align: 'left', label: 'Alasan Mutasi', field: 'alasan' },
  { name: 'aksi', align: 'center', label: 'Cetak' }
]

const riwayatKeluar = ref([
  { id: 201, tgl: '05/02/2026', nama: 'Fajar Ramadhan', kelasAsal: 'XII TKRO 1', sekolahTujuan: 'SMKN 1 Sumedang', alasan: 'Ikut Domisili Orang Tua' }
])

// Membuka Dialog dan Mereset Isi Form
const openMutasiDialog = () => {
  form.value = {
    nama: '', nisn: '', jurusan: '', sekolahAsal: '', sekolahTujuan: '',
    noSuratRekomendasi: '', alasan: '', berkasLengkap: false
  }
  dialogOpen.value = true
}

// Menangani Submit Form Simpan Data Mutasi
const simpanMutasi = () => {
  const tglHariIni = new Date().toLocaleDateString('id-ID')

  if (activeTab.value === 'masuk') {
    riwayatMasuk.value.unshift({
      id: Date.now(),
      tgl: tglHariIni,
      nama: form.value.nama,
      nisn: form.value.nisn,
      jurusan: form.value.jurusan,
      sekolahAsal: form.value.sekolahAsal,
      berkasLengkap: form.value.berkasLengkap
    })
    $q.notify({ color: 'green-8', message: 'Siswa masuk berhasil terdata di buku mutasi!', icon: 'check' })
  } else {
    riwayatKeluar.value.unshift({
      id: Date.now(),
      tgl: tglHariIni,
      nama: form.value.nama,
      kelasAsal: 'Deteksi Otomatis System',
      sekolahTujuan: form.value.sekolahTujuan,
      alasan: form.value.alasan
    })
    $q.notify({ color: 'red-9', message: 'Siswa berhasil dimutasikan keluar sistem.', icon: 'logout' })
  }

  dialogOpen.value = false
}

// Fungsi Cetak Berkas Surat Keterangan Pindah
const cetakSuratMutasi = (row) => {
  $q.notify({
    color: 'primary',
    message: `Mencetak dokumen mutasi untuk: ${row.nama}`,
    icon: 'print'
  })
}
</script>

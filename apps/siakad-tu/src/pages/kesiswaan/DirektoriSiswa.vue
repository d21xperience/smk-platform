<template>
  <q-page class="q-pa-md bg-grey-2">
    <!-- Header Halaman -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Direktori Siswa</h1>
        <p class="text-caption text-grey-7 q-mb-none">Manajemen data komputerisasi Buku Klaper dan data aktif siswa SMK
        </p>
      </div>
      <div>
        <q-btn color="primary" icon="add" label="Tambah Siswa" @click="openAddDialog" />
        <q-btn color="green-8" icon="file_download" label="Ekspor Excel" class="q-ml-sm" flat @click="exportToExcel" />
      </div>
    </div>

    <!-- Panel Filter & Pencarian -->
    <q-card class="shadow-1 q-mb-md">
      <q-card-section class="q-py-sm">
        <div class="row q-col-gutter-sm items-center">
          <!-- Input Pencarian -->
          <div class="col-12 col-md-4">
            <q-input v-model="filterText" dense outlined placeholder="Cari Nama, NIS, atau NISN..." clearable>
              <template v-slot:prepend>
                <q-icon name="search" />
              </template>
            </q-input>
          </div>

          <!-- Filter Jurusan -->
          <div class="col-12 col-sm-6 col-md-3">
            <q-select v-model="selectedJurusan" :options="opsiJurusan" dense outlined
              label="Kompetensi Keahlian (Jurusan)" emit-value map-options />
          </div>

          <!-- Filter Tingkat -->
          <div class="col-12 col-sm-6 col-md-2">
            <q-select v-model="selectedTingkat" :options="opsiTingkat" dense outlined label="Tingkat Kelas" emit-value
              map-options />
          </div>

          <!-- Tombol Reset Filter -->
          <div class="col-12 col-md-3 flex justify-end">
            <q-btn flat color="grey-7" label="Reset Filter" icon="refresh" @click="resetFilters"
              :disable="!isFilterActive" />
          </div>
        </div>
      </q-card-section>
    </q-card>

    <!-- Tabel Utama Data Siswa -->
    <q-card class="shadow-1">
      <q-card-section class="q-pa-none">
        <q-table flat :rows="filteredSiswa" :columns="columns" row-key="nisn" :loading="loading" selection="multiple"
          v-model:selected="selectedRows" :pagination="initialPagination">
          <!-- Kustomisasi Tampilan Kolom Status -->
          <template v-slot:body-cell-status="props">
            <q-td :props="props">
              <q-badge :color="props.row.status === 'Aktif' ? 'green' : 'orange'" :label="props.row.status" />
            </q-td>
          </template>

          <!-- Kustomisasi Tampilan Kolom Aksi -->
          <template v-slot:body-cell-aksi="props">
            <q-td :props="props" class="q-gutter-xs text-center">
              <q-btn size="sm" color="blue" icon="visibility" round flat @click="viewDetail(props.row)">
                <q-tooltip>Buku Induk Lengkap</q-tooltip>
              </q-btn>
              <q-btn size="sm" color="amber-9" icon="edit" round flat @click="editSiswa(props.row)">
                <q-tooltip>Ubah Data</q-tooltip>
              </q-btn>
              <q-btn size="sm" color="red" icon="person_remove" round flat @click="confirmDelete(props.row)">
                <q-tooltip>Proses Keluar / Mutasi</q-tooltip>
              </q-btn>
            </q-td>
          </template>
        </q-table>
      </q-card-section>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()

// State Loading dan Seleksi Baris Tabel
const loading = ref(false)
const selectedRows = ref([])

// State untuk Filter & Pencarian
const filterText = ref('')
const selectedJurusan = ref('Semua')
const selectedTingkat = ref('Semua')

// Opsi Pilihan Filter (Khas SMK)
const opsiJurusan = [
  { label: 'Semua Jurusan', value: 'Semua' },
  { label: 'Rekayasa Perangkat Lunak (RPL)', value: 'RPL' },
  { label: 'Teknik Kendaraan Ringan (TKRO)', value: 'TKRO' },
  { label: 'Akuntansi & Keuangan Lembaga (AKL)', value: 'AKL' }
]

const opsiTingkat = [
  { label: 'Semua Tingkat', value: 'Semua' },
  { label: 'Kelas X', value: 'X' },
  { label: 'Kelas XI', value: 'XI' },
  { label: 'Kelas XII', value: 'XII' }
]

// Konfigurasi Kolom Tabel Quasar
const columns = [
  { name: 'nis', align: 'left', label: 'NIS', field: 'nis', sortable: true },
  { name: 'nisn', align: 'left', label: 'NISN', field: 'nisn', sortable: true },
  { name: 'nama', align: 'left', label: 'Nama Lengkap', field: 'nama', sortable: true },
  { name: 'jk', align: 'center', label: 'JK', field: 'jk' },
  { name: 'jurusan', align: 'left', label: 'Jurusan', field: 'jurusan', sortable: true },
  { name: 'kelas', align: 'left', label: 'Kelas', field: 'kelas', sortable: true },
  { name: 'status', align: 'center', label: 'Status', field: 'status' },
  { name: 'aksi', align: 'center', label: 'Aksi' }
]

// Data Mock / Dummy Siswa Aktif SMK
const dataSiswa = ref([
  { nis: '22230101', nisn: '0061234561', nama: 'Ahmad Fauzi', jk: 'L', jurusan: 'RPL', kelas: 'XII RPL 1', status: 'Aktif' },
  { nis: '22230102', nisn: '0061234562', nama: 'Budi Santoso', jk: 'L', jurusan: 'TKRO', kelas: 'XII TKRO 3', status: 'Aktif' },
  { nis: '23240201', nisn: '0071234563', nama: 'Chandra Wijaya', jk: 'L', jurusan: 'RPL', kelas: 'XI RPL 2', status: 'Aktif' },
  { nis: '23240202', nisn: '0071234564', nama: 'Dinda Lestari', jk: 'P', jurusan: 'AKL', kelas: 'XI AKL 1', status: 'Aktif' },
  { nis: '24250301', nisn: '0081234565', nama: 'Eka Putri', jk: 'P', jurusan: 'AKL', kelas: 'X AKL 2', status: 'Aktif' },
  { nis: '22230105', nisn: '0061234569', nama: 'Fajar Ramadhan', jk: 'L', jurusan: 'TKRO', kelas: 'XII TKRO 1', status: 'Mutasi' }
])

// Pagination Awal q-table
const initialPagination = {
  sortBy: 'nama',
  descending: false,
  page: 1,
  rowsPerPage: 10
}

// Deteksi Apakah Filter Sedang Aktif
const isFilterActive = computed(() => {
  return filterText.value !== '' || selectedJurusan.value !== 'Semua' || selectedTingkat.value !== 'Semua'
})

// Fungsi Logika Pengolahan Filter Data di Frontend (Composition API)
const filteredSiswa = computed(() => {
  return dataSiswa.value.filter(siswa => {
    // 1. Filter Teks Global (Nama / NIS / NISN)
    const matchText = filterText.value
      ? siswa.nama.toLowerCase().includes(filterText.value.toLowerCase()) ||
      siswa.nis.includes(filterText.value) ||
      siswa.nisn.includes(filterText.value)
      : true

    // 2. Filter Jurusan
    const matchJurusan = selectedJurusan.value !== 'Semua'
      ? siswa.jurusan === selectedJurusan.value
      : true

    // 3. Filter Tingkat Kelas (Mengambil kata pertama dari string kelas, misal "XII RPL 1" diambil "XII")
    const matchTingkat = selectedTingkat.value !== 'Semua'
      ? siswa.kelas.split(' ')[0] === selectedTingkat.value
      : true

    return matchText && matchJurusan && matchTingkat
  })
})

// Fungsi Interaksi Aksi Tombol
const resetFilters = () => {
  filterText.value = ''
  selectedJurusan.value = 'Semua'
  selectedTingkat.value = 'Semua'
}

const openAddDialog = () => {
  $q.notify({ color: 'primary', message: 'Fitur Tambah Siswa Baru (Form PPDB)', icon: 'add' })
}

const viewDetail = (siswa) => {
  $q.notify({ color: 'info', message: `Membuka detail Buku Induk: ${siswa.nama}`, icon: 'visibility' })
}

const editSiswa = (siswa) => {
  $q.notify({ color: 'warning', message: `Ubah data siswa: ${siswa.nama}`, icon: 'edit', textColor: 'dark' })
}

const confirmDelete = (siswa) => {
  $q.dialog({
    title: 'Konfirmasi Mutasi Keluar',
    message: `Apakah Anda yakin ingin memproses status mutasi/keluar untuk siswa ${siswa.nama}?`,
    cancel: true,
    persistent: true
  }).onOk(() => {
    $q.notify({ color: 'negative', message: `${siswa.nama} berhasil diproses keluar.`, icon: 'person_remove' })
  })
}

const exportToExcel = () => {
  $q.notify({ color: 'green-8', message: 'Mengekspor data ke format .xlsx...', icon: 'file_download' })
}
</script>

<template>
  <q-page class="q-pa-md bg-grey-2">
    <!-- Header Halaman -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Surat Keterangan Aktif Belajar</h1>
        <p class="text-caption text-grey-7 q-mb-none">Pelayanan penerbitan, verifikasi data, dan pencetakan surat
          keterangan siswa aktif SMK</p>
      </div>
    </div>

    <div class="row q-col-gutter-md">
      <!-- 1. PANEL KIRI: Form Pembuatan / Pengajuan Surat Baru -->
      <div class="col-12 col-md-5">
        <q-card class="shadow-1">
          <q-card-section class="bg-primary text-white">
            <div class="text-subtitle1 text-weight-bold">Buat Surat Keterangan Baru</div>
            <div class="text-caption text-blue-2">Input cepat pencatatan agenda surat keluar kesiswaan</div>
          </q-card-section>

          <q-card-section class="q-pa-md">
            <q-form @submit.prevent="generateSurat" class="q-gutter-sm">
              <!-- Autocomplete pencarian siswa dari database lokal -->
              <q-select v-model="formSurat.siswaSelected" :options="opsiSiswaAktif" label="Cari & Pilih Siswa Aktif"
                outlined dense emit-value map-options use-input @filter="filterSiswa"
                :rules="[val => !!val || 'Siswa wajib dipilih']" />

              <!-- Menampilkan informasi kelas otomatis setelah siswa dipilih -->
              <div v-if="formSurat.siswaSelected"
                class="bg-blue-1 text-primary q-pa-sm rounded-borders text-caption q-mb-sm">
                <q-icon name="info" /> Terpilih: <strong>{{ formSurat.siswaSelected.nama }}</strong> ({{
                  formSurat.siswaSelected.kelas }}) • NISN: {{ formSurat.siswaSelected.nisn }}
              </div>

              <q-input v-model="formSurat.noSurat" label="Nomor Surat Resmi" outlined dense
                placeholder="Contoh: 421.5/102-SMK/{{ currentYear }}"
                :rules="[val => !!val || 'Nomor surat wajib diisi']" />

              <q-input v-model="formSurat.namaOrtu" label="Nama Orang Tua / Wali" outlined dense
                placeholder="Nama ayah atau ibu kandung..." :rules="[val => !!val || 'Nama orang tua wajib diisi']" />

              <q-select v-model="formSurat.keperluan"
                :options="['Tunjangan Gaji Orang Tua (PNS/BUMN)', 'Persyaratan BPJS Kesehatan', 'Pengajuan Beasiswa Luar Sekolah', 'Pembuatan Paspor Magang / Kunjungan Kerja']"
                label="Tujuan / Keperluan Surat" outlined dense
                :rules="[val => !!val || 'Keperluan surat wajib ditentukan']" />

              <div class="row justify-end q-mt-md">
                <q-btn label="Reset Form" type="reset" color="grey" flat class="q-mr-sm" @click="resetForm" />
                <q-btn label="Tambahkan ke Antrean" type="submit" color="primary" icon="playlist_add" />
              </div>
            </q-form>
          </q-card-section>
        </q-card>
      </div>

      <!-- 2. PANEL KANAN: Daftar Antrean Penerbitan Surat -->
      <div class="col-12 col-md-7">
        <q-card class="shadow-1 fit">
          <q-card-section class="q-py-sm flex justify-between items-center">
            <div>
              <div class="text-subtitle1 text-weight-bold text-grey-9">Antrean Penerbitan Surat</div>
              <div class="text-caption text-grey-6">Gunakan tombol aksi kanan untuk memicu modul printer pengolah
                dokumen</div>
            </div>
            <q-input v-model="searchQueue" dense outlined placeholder="Cari nama di antrean..."
              style="max-width: 200px;">
              <template v-slot:prepend><q-icon name="search" size="xs" /></template>
            </q-input>
          </q-card-section>

          <q-card-section class="q-pa-none">
            <q-table flat :rows="filteredQueue" :columns="columns" row-key="id" :pagination="{ rowsPerPage: 7 }">
              <!-- Kustomisasi Badge Status Proses -->
              <template v-slot:body-cell-status="props">
                <q-td :props="props" class="text-center">
                  <q-chip size="xs" :color="props.row.status === 'Selesai' ? 'green-1' : 'orange-1'"
                    :text-color="props.row.status === 'Selesai' ? 'green-9' : 'orange-9'">
                    {{ props.row.status }}
                  </q-chip>
                </q-td>
              </template>

              <!-- Kustomisasi Kolom Aksi Cetak & Update Status -->
              <template v-slot:body-cell-aksi="props">
                <q-td :props="props" class="text-center q-gutter-xs">
                  <q-btn size="sm" color="indigo" icon="print" round flat @click="printSurat(props.row)">
                    <q-tooltip>Cetak Lembar Surat Aktif</q-tooltip>
                  </q-btn>
                  <q-btn size="sm" color="green-8" icon="check_circle" round flat
                    :disable="props.row.status === 'Selesai'" @click="completeSurat(props.row)">
                    <q-tooltip>Tandai Selesai Diambil</q-tooltip>
                  </q-btn>
                </q-td>
              </template>
            </q-table>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()
// const currentYear = ref(new Date().getFullYear())

// State Penampung Form
const formSurat = ref({
  siswaSelected: null,
  noSurat: `421.5/084-SMK/2026`,
  namaOrtu: '',
  keperluan: 'Persyaratan BPJS Kesehatan'
})

const searchQueue = ref('')

// Database Dummy Frontend untuk Autocomplete Dropdown Siswa
const masterSiswaAktif = [
  { label: 'Ahmad Fauzi (XII RPL 1)', value: { nama: 'Ahmad Fauzi', kelas: 'XII RPL 1', nisn: '0061234561' } },
  { label: 'Budi Santoso (XII TKRO 3)', value: { nama: 'Budi Santoso', kelas: 'XII TKRO 3', nisn: '0061234562' } },
  { label: 'Chandra Wijaya (XI RPL 2)', value: { nama: 'Chandra Wijaya', kelas: 'XI RPL 2', nisn: '0071234563' } },
  { label: 'Dinda Lestari (XI AKL 1)', value: { nama: 'Dinda Lestari', kelas: 'XI AKL 1', nisn: '0071234564' } },
  { label: 'Eka Putri (X AKL 2)', value: { nama: 'Eka Putri', kelas: 'X AKL 2', nisn: '0081234565' } }
]

const opsiSiswaAktif = ref(masterSiswaAktif)

// Data List Antrean Cetak Surat
const dataAntreanSurat = ref([
  { id: 1, noSurat: '421.5/081-SMK/2026', nama: 'Muhammad Reyhan', kelas: 'XII RPL 2', keperluan: 'Tunjangan Gaji Orang Tua', status: 'Selesai' },
  { id: 2, noSurat: '421.5/082-SMK/2026', nama: 'Siti Aminah', kelas: 'XI AKL 1', keperluan: 'Persyaratan BPJS Kesehatan', status: 'Menunggu Cetak' },
  { id: 3, noSurat: '421.5/083-SMK/2026', nama: 'Riska Amelia', kelas: 'X ULP 1', keperluan: 'Pembuatan Paspor Magang', status: 'Menunggu Cetak' }
])

// Judul Struktur Kolom Tabel Antrean
const columns = [
  { name: 'noSurat', align: 'left', label: 'No. Surat Keluar', field: 'noSurat', sortable: true },
  { name: 'nama', align: 'left', label: 'Nama Siswa', field: 'nama', sortable: true },
  { name: 'kelas', align: 'left', label: 'Kelas', field: 'kelas' },
  { name: 'keperluan', align: 'left', label: 'Keperluan Dokumen', field: 'keperluan' },
  { name: 'status', align: 'center', label: 'Status Proses', field: 'status' },
  { name: 'aksi', align: 'center', label: 'Aksi' }
]

// Fungsi Filter Pintar untuk Fitur Pencarian Dropdown Autocomplete
const filterSiswa = (val, update) => {
  if (val === '') {
    update(() => { opsiSiswaAktif.value = masterSiswaAktif })
    return
  }
  update(() => {
    const needle = val.toLowerCase()
    opsiSiswaAktif.value = masterSiswaAktif.filter(v => v.label.toLowerCase().includes(needle))
  })
}

// Logika Compute Filter Pencarian Tabel Antrean Sisi Klien
const filteredQueue = computed(() => {
  if (!searchQueue.value) return dataAntreanSurat.value
  return dataAntreanSurat.value.filter(item =>
    item.nama.toLowerCase().includes(searchQueue.value.toLowerCase()) ||
    item.noSurat.includes(searchQueue.value)
  )
})

// Fungsi Eksekusi Administrasi Persuratan
const generateSurat = () => {
  const payload = formSurat.value

  // Masukkan pengajuan surat baru ke baris antrean kanan
  dataAntreanSurat.value.unshift({
    id: Date.now(),
    noSurat: payload.noSurat,
    nama: payload.siswaSelected.nama,
    kelas: payload.siswaSelected.kelas,
    keperluan: payload.keperluan,
    status: 'Menunggu Cetak'
  })

  $q.notify({
    color: 'green-8',
    message: `Surat Keterangan untuk ${payload.siswaSelected.nama} berhasil didaftarkan ke antrean!`,
    icon: 'check_circle'
  })

  resetForm()
}

const resetForm = () => {
  formSurat.value = {
    siswaSelected: null,
    noSurat: `421.5/08${dataAntreanSurat.value.length + 1}-SMK/2026`,
    namaOrtu: '',
    keperluan: 'Persyaratan BPJS Kesehatan'
  }
}

const printSurat = (row) => {
  $q.notify({
    color: 'indigo-8',
    message: `Menghubungkan ke printer kantor... Mencetak dokumen resmi nomor: ${row.noSurat}`,
    icon: 'print'
  })
}

const completeSurat = (row) => {
  row.status = 'Selesai'
  $q.notify({
    color: 'green-9',
    message: `Status berkas surat atas nama ${row.nama} dinyatakan SELESAI diambil.`,
    icon: 'assignment_turned_in'
  })
}
</script>

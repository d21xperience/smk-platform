<template>
  <q-page class="q-pa-md bg-grey-2">
    <!-- Header Halaman -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Manajemen Beasiswa & PIP</h1>
        <p class="text-caption text-grey-7 q-mb-none">Verifikasi berkas persyaratan, pendataan rekening, dan pemantauan
          pencairan dana Program Indonesia Pintar</p>
      </div>
      <div>
        <q-btn color="primary" icon="add" label="Usulkan Siswa" @click="openUsulanDialog" />
        <q-btn color="green-8" icon="download" label="Ekspor Format Dapodik" class="q-ml-sm" flat
          @click="exportDapodik" />
      </div>
    </div>

    <!-- Panel Informasi Ringkasan Kuota PIP -->
    <div class="row q-col-gutter-sm q-mb-md">
      <div class="col-12 col-sm-3">
        <q-card class="bg-white text-dark flat shadow-1">
          <q-card-section class="q-py-sm flex justify-between items-center">
            <div>
              <div class="text-caption text-grey-7 text-uppercase">Total Penerima</div>
              <div class="text-h5 text-weight-bold text-primary">124 Siswa</div>
            </div>
            <q-avatar color="blue-1" text-color="primary" icon="payments" size="40px" />
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-sm-3">
        <q-card class="bg-white text-dark flat shadow-1">
          <q-card-section class="q-py-sm flex justify-between items-center">
            <div>
              <div class="text-caption text-grey-7 text-uppercase">Berkas Terverifikasi</div>
              <div class="text-h5 text-weight-bold text-green-7">112 Siswa</div>
            </div>
            <q-avatar color="green-1" text-color="green-7" icon="verified" size="40px" />
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-sm-3">
        <q-card class="bg-white text-dark flat shadow-1">
          <q-card-section class="q-py-sm flex justify-between items-center">
            <div>
              <div class="text-caption text-grey-7 text-uppercase">Menunggu Review</div>
              <div class="text-h5 text-weight-bold text-orange-8">12 Siswa</div>
            </div>
            <q-avatar color="orange-1" text-color="orange-8" icon="rate_review" size="40px" />
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-sm-3">
        <q-card class="bg-white text-dark flat shadow-1">
          <q-card-section class="q-py-sm flex justify-between items-center">
            <div>
              <div class="text-caption text-grey-7 text-uppercase">Dana Cair (SK Pemberian)</div>
              <div class="text-h5 text-weight-bold text-purple-7">84 Siswa</div>
            </div>
            <q-avatar color="purple-1" text-color="purple-7" icon="price_check" size="40px" />
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Panel Filter & Tabel -->
    <q-card class="shadow-1">
      <q-card-section class="q-py-sm">
        <div class="row q-col-gutter-sm items-center">
          <div class="col-12 col-md-5">
            <q-input v-model="search" dense outlined placeholder="Cari nama siswa, NISN, atau virtual account..."
              clearable>
              <template v-slot:prepend><q-icon name="search" /></template>
            </q-input>
          </div>
          <div class="col-12 col-md-4">
            <q-select v-model="filterTahap"
              :options="['Semua Tahap', 'SK Nominasi (Belum Aktivasi)', 'SK Pemberian (Sudah Cair)', 'Ditolak / Tidak Layak']"
              dense outlined label="Status Pencairan Dana" />
          </div>
          <div class="col-12 col-md-3 flex justify-end text-caption text-grey-6">
            Bank Penyalur SMK: <strong class="text-red-9 q-ml-xs">BNI (Bank Negara Indonesia)</strong>
          </div>
        </div>
      </q-card-section>

      <!-- Tabel Data Beasiswa -->
      <q-card-section class="q-pa-none">
        <q-table flat :rows="filteredBeasiswa" :columns="columns" row-key="id" :pagination="{ rowsPerPage: 10 }">
          <!-- Kolom Detail Kategori Jaminan Sosial -->
          <template v-slot:body-cell-kategori="props">
            <q-td :props="props">
              <div class="text-weight-medium">{{ props.row.kategori }}</div>
              <div class="text-caption text-grey-6">ID: {{ props.row.noKartuJaminan || 'Usulan Sekolah' }}</div>
            </q-td>
          </template>

          <!-- Kolom Status Berkas Fisik Persyaratan -->
          <template v-slot:body-cell-berkas="props">
            <q-td :props="props" class="text-center">
              <q-btn dense flat size="sm" :color="props.row.berkasValid ? 'green' : 'orange-9'"
                :icon="props.row.berkasValid ? 'check_box' : 'indeterminate_check_box'"
                :label="props.row.berkasValid ? 'Lengkap & Sah' : 'Belum Diverifikasi'"
                @click="toggleBerkas(props.row)" />
            </q-td>
          </template>

          <!-- Kolom Status Tahap Pencairan Beasiswa -->
          <template v-slot:body-cell-status="props">
            <q-td :props="props" class="text-center">
              <q-badge
                :color="props.row.status === 'SK Pemberian (Cair)' ? 'purple' : props.row.status === 'SK Nominasi' ? 'orange' : 'red'"
                :label="props.row.status" />
            </q-td>
          </template>

          <!-- Kolom Aksi Administratif TU -->
          <template v-slot:body-cell-aksi="props">
            <q-td :props="props" class="text-center q-gutter-xs">
              <q-btn size="sm" color="primary" icon="edit" round flat @click="editBeasiswa(props.row)">
                <q-tooltip>Update Rekening & No. SimPINTAR</q-tooltip>
              </q-btn>
              <q-btn size="sm" color="indigo-7" icon="print" round flat @click="printSuratPengantar(props.row)">
                <q-tooltip>Cetak Surat Pengantar Aktivasi Bank</q-tooltip>
              </q-btn>
            </q-td>
          </template>
        </q-table>
      </q-card-section>
    </q-card>

    <!-- DIALOG POPUP: Usulkan Penerima Beasiswa Baru -->
    <q-dialog v-model="dialogUsulanOpen" persistent>
      <q-card style="min-width: 420px;">
        <q-card-section class="bg-primary text-white">
          <div class="text-h6">Form Usulan Beasiswa PIP</div>
          <div class="text-caption text-blue-2">Menambahkan entitas siswa layak bantu ke draf pelaporan kesiswaan</div>
        </q-card-section>

        <q-card-section class="q-gutter-md q-pt-md">
          <q-select v-model="formUsulan.siswa" :options="opsiSiswaLayak" label="Pilih Nama Siswa" outlined dense />
          <q-select v-model="formUsulan.kategori"
            :options="['Pemegang KIP', 'Keluarga PKH / KKS', 'Pertimbangan Khusus (SKTM)']"
            label="Dasar Kategori Kelayakan" outlined dense />
          <q-input v-model="formUsulan.noKartu" label="Nomor Kartu Jaminan Sosial (KIP/KKS)" outlined dense
            placeholder="Kosongkan jika menggunakan SKTM..." />
          <q-input v-model="formUsulan.alasan" type="textarea" rows="2" label="Alasan Pertimbangan Usulan" outlined
            dense placeholder="Contoh: Orang tua terkena PHK / Yatim Piatu..." />
        </q-card-section>

        <q-card-actions align="right" class="q-pb-md q-pr-md">
          <q-btn flat label="Batal" v-close-popup color="grey" />
          <q-btn color="primary" label="Simpan Usulan" @click="saveUsulanBaru" />
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
const filterTahap = ref('Semua Tahap')
const dialogUsulanOpen = ref(false)

const formUsulan = ref({ siswa: null, kategori: 'Pemegang KIP', noKartu: '', alasan: '' })

// Definisi Struktur Kolom Tabel q-table
const columns = [
  { name: 'nama', align: 'left', label: 'Nama Siswa', field: 'nama', sortable: true },
  { name: 'kelas', align: 'left', label: 'Kelas', field: 'kelas', sortable: true },
  { name: 'kategori', align: 'left', label: 'Kategori Kelayakan', field: 'kategori' },
  { name: 'berkas', align: 'center', label: 'Status Dokumen Fisik', field: 'berkasValid' },
  { name: 'status', align: 'center', label: 'Status Dana PIP', field: 'status', sortable: true },
  { name: 'aksi', align: 'center', label: 'Aksi' }
]

// Data Mock / Dummy Pelacakan Administrasi Beasiswa PIP SMK
const dataBeasiswa = ref([
  { id: 1, nama: 'Ahmad Fauzi', kelas: 'XII RPL 1', kategori: 'Pemegang KIP', noKartuJaminan: 'KIP-2026-8871A', berkasValid: true, status: 'SK Pemberian (Cair)' },
  { id: 2, nama: 'Budi Santoso', kelas: 'XII TKRO 3', kategori: 'Keluarga PKH / KKS', noKartuJaminan: 'KKS-3204-1192B', berkasValid: true, status: 'SK Nominasi' },
  { id: 3, nama: 'Anisa Rahmawati', kelas: 'XII AKL 2', kategori: 'Pertimbangan Khusus (SKTM)', noKartuJaminan: '', berkasValid: false, status: 'SK Nominasi' },
  { id: 4, nama: 'Dinda Lestari', kelas: 'XI AKL 1', kategori: 'Pemegang KIP', noKartuJaminan: 'KIP-2026-9921C', berkasValid: true, status: 'SK Pemberian (Cair)' }
])

const opsiSiswaLayak = [
  { label: 'Chandra Wijaya (XI RPL 2)', value: 'Chandra Wijaya' },
  { label: 'Eka Putri (X AKL 2)', value: 'Eka Putri' }
]

// Logika Compute Filter Pencarian Frontend
const filteredBeasiswa = computed(() => {
  return dataBeasiswa.value.filter(item => {
    const matchSearch = search.value
      ? item.nama.toLowerCase().includes(search.value.toLowerCase()) ||
      item.noKartuJaminan.includes(search.value)
      : true

    const matchTahap = filterTahap.value !== 'Semua Tahap'
      ? item.status.includes(filterTahap.value.split(' ')[1]) // Mencocokkan kata kunci seperti "Pemberian" atau "Nominasi"
      : true
    return matchSearch && matchTahap
  })
})
// Fungsi Penanganan Event Aksi Kesiswaan
const toggleBerkas = (row) => {
  row.berkasValid = !row.berkasValid
  $q.notify({ color: row.berkasValid ? 'green-8' : 'orange-9', message: `Berkas ${row.nama} dinyatakan: ${row.berkasValid ? 'VALID & LENGKAP' : 'BELUM VALID'}`, icon: 'assignment_turned_in' })
}
const openUsulanDialog = () => {
  formUsulan.value = { siswa: null, kategori: 'Pemegang KIP', noKartu: '', alasan: '' }
  dialogUsulanOpen.value = true
}
const saveUsulanBaru = () => {
  if (!formUsulan.value.siswa) {
    $q.notify({ color: 'negative', message: 'Nama siswa wajib dipilih!', icon: 'warning' })
    return
  }
  dataBeasiswa.value.push({ id: Date.now(), nama: formUsulan.value.siswa, kelas: 'Tingkat Aktif', kategori: formUsulan.value.kategori, noKartuJaminan: formUsulan.value.noKartu || 'Sedang Diproses Dinas', berkasValid: false, status: 'SK Nominasi' })
  $q.notify({ color: 'primary', message: `Sukses mendaftarkan draf usulan beasiswa PIP untuk siswa: ${formUsulan.value.siswa}`, icon: 'done' })
  dialogUsulanOpen.value = false
}
const editBeasiswa = (row) => {
  $q.dialog({ title: 'Update Buku Tabungan & SimPINTAR', message: ` Masukkan nomor rekening BNI aktif untuk siswa ${row.nama}:`, prompt: { model: '', type: 'text', placeholder: 'Contoh: 117283621...' }, cancel: true, persistent: true }).onOk(data => {
    $q.notify({ color: 'primary', message: `Nomor rekening ${data} sukses disinkronkan ke profil penerima dana.`, icon: 'save' })
  })
}
const printSuratPengantar = (row) => { $q.notify({ color: 'indigo-8', message: `Mencetak Surat Pengantar Aktivasi Rekening / Pencairan Kolektif ke Bank BNI untuk ${row.nama}.`, icon: 'print' }) }
const exportDapodik = () => { $q.notify({ color: 'green-8', message: 'Mengunduh berkas kompilasi data usulan berformat .csv untuk diunggah ke server Dapodik.', icon: 'file_download' }) }
</script>

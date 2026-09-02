<template>
  <q-page class="q-pa-md bg-grey-2">
    <!-- Header Halaman -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Manajemen Lowongan Kerja BKK</h1>
        <p class="text-caption text-grey-7 q-mb-none">Pengelolaan informasi karir, kemitraan industri, dan fasilitasi
          penempatan kerja alumni SMK</p>
      </div>
      <div>
        <q-btn color="primary" icon="add_box" label="Tambah Info Loker" @click="openLokerDialog" />
      </div>
    </div>

    <!-- Ringkasan Statistik Loker BKK -->
    <div class="row q-col-gutter-sm q-mb-md">
      <div class="col-12 col-sm-3">
        <q-card class="bg-white text-dark flat shadow-1">
          <q-card-section class="q-py-sm flex justify-between items-center">
            <div>
              <div class="text-caption text-grey-7 text-uppercase">Loker Aktif</div>
              <div class="text-h5 text-weight-bold text-primary">12 Lowongan</div>
            </div>
            <q-avatar color="blue-1" text-color="primary" icon="campaign" size="40px" />
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-sm-3">
        <q-card class="bg-white text-dark flat shadow-1">
          <q-card-section class="q-py-sm flex justify-between items-center">
            <div>
              <div class="text-caption text-grey-7 text-uppercase">Pelamar Terdata</div>
              <div class="text-h5 text-weight-bold text-green-7">84 Alumni</div>
            </div>
            <q-avatar color="green-1" text-color="green-7" icon="people" size="40px" />
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-sm-3">
        <q-card class="bg-white text-dark flat shadow-1">
          <q-card-section class="q-py-sm flex justify-between items-center">
            <div>
              <div class="text-caption text-grey-7 text-uppercase">Mitra Baru (Bulan Ini)</div>
              <div class="text-h5 text-weight-bold text-purple-7">4 Perusahaan</div>
            </div>
            <q-avatar color="purple-1" text-color="purple-7" icon="handshake" size="40px" />
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-sm-3">
        <q-card class="bg-white text-dark flat shadow-1">
          <q-card-section class="q-py-sm flex justify-between items-center">
            <div>
              <div class="text-caption text-grey-7 text-uppercase">Mendekati Deadline</div>
              <div class="text-h5 text-weight-bold text-red-8">2 Loker</div>
            </div>
            <q-avatar color="red-1" text-color="red-8" icon="alarm" size="40px" />
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Panel Filter Pencarian & Tabel Lowongan -->
    <q-card class="shadow-1">
      <q-card-section class="q-py-sm">
        <div class="row q-col-gutter-sm items-center">
          <div class="col-12 col-md-5">
            <q-input v-model="search" dense outlined placeholder="Cari posisi kerja, nama perusahaan, atau lokasi..."
              clearable>
              <template v-slot:prepend><q-icon name="search" /></template>
            </q-input>
          </div>
          <div class="col-12 col-md-4">
            <q-select v-model="filterJurusan" :options="['Semua Jurusan', 'RPL', 'TKRO', 'AKL']" dense outlined
              label="Kualifikasi Kompetensi" />
          </div>
          <div class="col-12 col-md-3 flex justify-end text-caption text-grey-6">
            Status Sistem: <q-badge color="green" label="Online Terbuka" class="q-ml-sm" />
          </div>
        </div>
      </q-card-section>

      <!-- Tabel Data Informasi Lowongan Kerja -->
      <q-card-section class="q-pa-none">
        <q-table flat :rows="filteredLoker" :columns="columns" row-key="id" :pagination="{ rowsPerPage: 10 }">
          <!-- Custom Tampilan Kolom Perusahaan & Posisi -->
          <template v-slot:body-cell-posisi="props">
            <q-td :props="props">
              <div class="text-weight-bold text-primary">{{ props.row.posisi }}</div>
              <div class="text-caption text-grey-9">{{ props.row.perusahaan }} ({{ props.row.lokasi }})</div>
            </q-td>
          </template>

          <!-- Custom Tampilan Badge Kualifikasi Jurusan SMK -->
          <template v-slot:body-cell-kualifikasi="props">
            <q-td :props="props" class="text-center">
              <q-badge color="teal" :label="'Lulusan ' + props.row.kualifikasi" />
            </q-td>
          </template>

          <!-- Custom Tampilan Status Publikasi Informasi -->
          <template v-slot:body-cell-status="props">
            <q-td :props="props" class="text-center">
              <q-chip clickable @click="toggleStatusPublikasi(props.row)" size="sm"
                :color="props.row.status === 'Aktif Terbit' ? 'green-1' : 'red-1'"
                :text-color="props.row.status === 'Aktif Terbit' ? 'green-9' : 'red-9'"
                :icon="props.row.status === 'Aktif Terbit' ? 'visibility' : 'visibility_off'">
                {{ props.row.status }}
              </q-chip>
            </q-td>
          </template>

          <!-- Custom Kolom Aksi Staf BKK -->
          <template v-slot:body-cell-aksi="props">
            <q-td :props="props" class="text-center q-gutter-xs">
              <q-btn size="sm" color="indigo" icon="group" round flat @click="lihatPelamar(props.row)">
                <q-tooltip>Lihat Daftar Pelamar Alumni</q-tooltip>
              </q-btn>
              <q-btn size="sm" color="red" icon="delete" round flat @click="hapusLoker(props.row)">
                <q-tooltip>Turunkan / Hapus Loker</q-tooltip>
              </q-btn>
            </q-td>
          </template>
        </q-table>
      </q-card-section>
    </q-card>

    <!-- DIALOG POPUP: Form Pemasukan Lowongan Kerja Baru -->
    <q-dialog v-model="dialogLokerOpen" persistent>
      <q-card style="min-width: 420px;">
        <q-card-section class="bg-primary text-white">
          <div class="text-h6">Tambah Informasi Lowongan Kerja</div>
          <div class="text-caption text-blue-2">Publikasikan peluang karir industri ke portal alumni BKK</div>
        </q-card-section>

        <q-card-section class="q-gutter-md q-pt-md">
          <q-input v-model="formLoker.perusahaan" label="Nama Perusahaan / Industri Mitra" outlined dense
            placeholder="Contoh: PT. Astra Diatsu Motor..." />
          <q-input v-model="formLoker.posisi" label="Posisi / Jabatan Pekerjaan" outlined dense
            placeholder="Contoh: Maintenance Staff / Jr. Programmer..." />

          <div class="row q-col-gutter-sm">
            <div class="col-6">
              <q-select v-model="formLoker.kualifikasi" :options="['RPL', 'TKRO', 'AKL']" label="Kualifikasi Jurusan"
                outlined dense />
            </div>
            <div class="col-6">
              <q-input v-model="formLoker.lokasi" label="Kota Penempatan" outlined dense
                placeholder="Contoh: Bekasi / Karawang" />
            </div>
          </div>

          <q-input v-model="formLoker.deadline" label="Batas Akhir Pendaftaran (Deadline)" outlined dense
            placeholder="Contoh: 31 Juli 2026" />
        </q-card-section>

        <q-card-actions align="right" class="q-pb-md q-pr-md">
          <q-btn flat label="Batal" v-close-popup color="grey" />
          <q-btn color="primary" label="Terbitkan Loker" @click="saveLokerBaru" />
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
const filterJurusan = ref('Semua Jurusan')
const dialogLokerOpen = ref(false)

const formLoker = ref({ perusahaan: '', posisi: '', kualifikasi: 'RPL', lokasi: '', deadline: '' })

// Struktur Kolom Tabel q-table Log Lowongan Kerja
const columns = [
  { name: 'posisi', align: 'left', label: 'Detail Lowongan & Perusahaan', field: 'posisi', sortable: true },
  { name: 'kualifikasi', align: 'center', label: 'Kualifikasi Jurusan', field: 'kualifikasi' },
  { name: 'deadline', align: 'left', label: 'Batas Pendaftaran', field: 'deadline', sortable: true },
  { name: 'status', align: 'center', label: 'Status Portal', field: 'status' },
  { name: 'aksi', align: 'center', label: 'Aksi' }
]

// Data Mock / Dummy Portal Karir Loker BKK SMK
const dataLoker = ref([
  { id: 1, perusahaan: 'PT. Toyota Motor Manufacturing', posisi: 'Operator Perakitan Otomotif', kualifikasi: 'TKRO', lokasi: 'Karawang', deadline: '30 Juni 2026', status: 'Aktif Terbit' },
  { id: 2, perusahaan: 'PT. Tokopedia Indonesia', posisi: 'Junior QA Automation', kualifikasi: 'RPL', lokasi: 'Jakarta Selatan', deadline: '15 Juli 2026', status: 'Aktif Terbit' },
  { id: 3, perusahaan: 'KAP Tanudiredja (PwC Indonesia)', posisi: 'Junior Audit Associate', kualifikasi: 'AKL', lokasi: 'Jakarta Pusat', deadline: '05 Juli 2026', status: 'Draf Peninjauan' },
  { id: 4, perusahaan: 'PT. Epson Indonesia', posisi: 'Staff QC Elektronik', kualifikasi: 'TKRO', lokasi: 'Bekasi', deadline: '22 Juni 2026', status: 'Aktif Terbit' }
])

// Logika Compute Filter Pencarian Lowongan Kerja di Sisi Klien
const filteredLoker = computed(() => {
  return dataLoker.value.filter(item => {
    const matchSearch = search.value
      ? item.perusahaan.toLowerCase().includes(search.value.toLowerCase()) ||
      item.posisi.toLowerCase().includes(search.value.toLowerCase()) ||
      item.lokasi.toLowerCase().includes(search.value.toLowerCase())
      : true

    const matchJurusan = filterJurusan.value !== 'Semua Jurusan'
      ? item.kualifikasi === filterJurusan.value
      : true

    return matchSearch && matchJurusan
  })
})


// Fungsi Penanganan Event Aksi Kesiswaan & BKK
const toggleStatusPublikasi = (row) => {
  row.status = row.status === 'Aktif Terbit' ? 'Draf Peninjauan' : 'Aktif Terbit'
  $q.notify({ color: row.status === 'Aktif Terbit' ? 'green-8' : 'orange-9', message: `Informasi loker ${row.perusahaan} diubah menjadi: ${row.status.toUpperCase()}`, icon: 'sync' })
}
const openLokerDialog = () => {
  formLoker.value = { perusahaan: '', posisi: '', kualifikasi: 'RPL', lokasi: '', deadline: '31 Juli 2026' }
  dialogLokerOpen.value = true
}
const saveLokerBaru = () => {
  if (!formLoker.value.perusahaan || !formLoker.value.posisi) {
    $q.notify({ color: 'negative', message: 'Nama perusahaan dan Posisi jabatan wajib diisi!', icon: 'warning' })
    return
  }
  dataLoker.value.unshift({ id: Date.now(), perusahaan: formLoker.value.perusahaan, posisi: formLoker.value.posisi, kualifikasi: formLoker.value.kualifikasi, lokasi: formLoker.value.lokasi || 'Nasional', deadline: formLoker.value.deadline, status: 'Aktif Terbit' })
  $q.notify({ color: 'primary', message: `Peluang karir di ${formLoker.value.perusahaan} resmi dipublikasikan ke kanal alumni.`, icon: 'campaign' })
  dialogLokerOpen.value = false
}
const lihatPelamar = (row) => { $q.notify({ color: 'indigo-8', message: `Membuka draf berkas lamaran masuk dari alumni untuk lowongan: ${row.posisi}`, icon: 'folder_shared' }) }
const hapusLoker = (row) => {
  $q.dialog({ title: 'Konfirmasi Penghapusan', message: `Apakah Anda yakin ingin menghapus atau menurunkan iklan lowongan ${row.posisi} dari ${row.perusahaan} ?`, cancel: true, persistent: true }).onOk(() => {
    dataLoker.value = dataLoker.value.filter(item => item.id !== row.id)
    $q.notify({ color: 'negative', message: 'Informasi lowongan berhasil dihapus dari papan pengumuman digital BKK.', icon: 'delete' })
  })
}

</script>

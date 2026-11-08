<template>
  <q-page class="q-pa-md bg-grey-2">
    <!-- Header Halaman -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Plotting & Berkas PKL</h1>
        <p class="text-caption text-grey-7 q-mb-none">Manajemen hubungan industri, penempatan magang, dan administrasi
          dokumen PKL siswa SMK</p>
      </div>
      <div>
        <q-btn color="primary" icon="assignment_ind" label="Plotting Siswa" @click="openPlottingDialog" />
      </div>
    </div>

    <!-- Statistik Ringkas PKL -->
    <div class="row q-col-gutter-sm q-mb-md">
      <div class="col-12 col-sm-4">
        <q-card class="bg-blue-1 text-primary flat shadow-1">
          <q-card-section class="q-py-sm flex justify-between items-center">
            <div>
              <div class="text-caption text-weight-medium">Total Siswa Terplotting</div>
              <div class="text-h5 text-weight-bold">342 / 350 Siswa</div>
            </div>
            <q-icon name="check_circle" size="md" />
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-sm-4">
        <q-card class="bg-orange-1 text-orange-9 flat shadow-1">
          <q-card-section class="q-py-sm flex justify-between items-center">
            <div>
              <div class="text-caption text-weight-medium">Belum Dapat Tempat</div>
              <div class="text-h5 text-weight-bold">8 Siswa</div>
            </div>
            <q-icon name="pending" size="md" />
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-sm-4">
        <q-card class="bg-green-1 text-green-9 flat shadow-1">
          <q-card-section class="q-py-sm flex justify-between items-center">
            <div>
              <div class="text-caption text-weight-medium">Mitra Industri Aktif (DUDI)</div>
              <div class="text-h5 text-weight-bold">48 Perusahaan</div>
            </div>
            <q-icon name="corporate_fare" size="md" />
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Panel Filter & Tabel Utama -->
    <q-card class="shadow-1">
      <q-card-section class="q-py-sm">
        <div class="row q-col-gutter-sm items-center">
          <div class="col-12 col-md-5">
            <q-input v-model="search" dense outlined placeholder="Cari nama siswa, kelas, atau nama perusahaan..."
              clearable>
              <template v-slot:prepend><q-icon name="search" /></template>
            </q-input>
          </div>
          <div class="col-12 col-md-4">
            <q-select v-model="filterJurusan" :options="['Semua Jurusan', 'RPL', 'TKRO', 'AKL']" dense outlined
              label="Filter Kompetensi Keahlian" />
          </div>
          <div class="col-12 col-md-3 flex justify-end">
            <q-btn-dropdown color="green-8" icon="print" label="Cetak Berkas Masal">
              <q-list>
                <q-item clickable v-close-popup @click="bulkPrint('Surat Pengantar')">
                  <q-item-section><q-item-label>Surat Pengantar PKL</q-item-label></q-item-section>
                </q-item>
                <q-item clickable v-close-popup @click="bulkPrint('Surat Izin')">
                  <q-item-section><q-item-label>Surat Izin Orang Tua</q-item-label></q-item-section>
                </q-item>
                <q-item clickable v-close-popup @click="bulkPrint('Form Nilai')">
                  <q-item-section><q-item-label>Lembar Penilaian Industri</q-item-label></q-item-section>
                </q-item>
              </q-list>
            </q-btn-dropdown>
          </div>
        </div>
      </q-card-section>

      <!-- Tabel Plotting -->
      <q-card-section class="q-pa-none">
        <q-table flat :rows="filteredPlotting" :columns="columns" row-key="id" selection="multiple"
          v-model:selected="selectedSiswa" :pagination="{ rowsPerPage: 10 }">
          <!-- Custom badge untuk nama perusahaan industri -->
          <template v-slot:body-cell-perusahaan="props">
            <q-td :props="props">
              <div class="text-weight-bold text-grey-9">{{ props.row.perusahaan }}</div>
              <div class="text-caption text-grey-6">Kota: {{ props.row.kota }}</div>
            </q-td>
          </template>

          <!-- Custom Kolom Aksi Cetak Dokumen Mandiri -->
          <template v-slot:body-cell-aksi="props">
            <q-td :props="props" class="text-center q-gutter-xs">
              <q-btn size="sm" color="indigo-7" icon="description" round flat
                @click="printSingle(props.row, 'Surat Pengantar')">
                <q-tooltip>Cetak Surat Pengantar</q-tooltip>
              </q-btn>
              <q-btn size="sm" color="teal-7" icon="family_restroom" round flat
                @click="printSingle(props.row, 'Surat Izin')">
                <q-tooltip>Cetak Surat Izin Orang Tua</q-tooltip>
              </q-btn>
              <q-btn size="sm" color="amber-9" icon="edit" round flat @click="editPlotting(props.row)">
                <q-tooltip>Ubah Penempatan / Pembimbing</q-tooltip>
              </q-btn>
            </q-td>
          </template>
        </q-table>
      </q-card-section>
    </q-card>

    <!-- DIALOG POPUP: Form Plotting Baru / Ubah Penempatan -->
    <q-dialog v-model="dialogOpen" persistent>
      <q-card style="min-width: 450px;">
        <q-card-section class="bg-primary text-white">
          <div class="text-h6">{{ isEdit ? 'Ubah Penempatan PKL' : 'Plotting Penempatan PKL Baru' }}</div>
        </q-card-section>

        <q-card-section class="q-gutter-md q-pt-md">
          <!-- Input hanya tampil jika membuat data baru -->
          <q-select v-if="!isEdit" v-model="form.siswa" :options="opsiSiswaBelumPkl" label="Pilih Siswa (Belum PKL)"
            outlined dense emit-value map-options />
          <q-input v-else v-model="form.namaSiswa" label="Nama Siswa" outlined dense readonly />

          <q-select v-model="form.perusahaan" :options="opsiPerusahaanMitra" label="Perusahaan Mitra (DUDI)" outlined
            dense emit-value map-options />
          <q-input v-model="form.pembimbing" label="Guru Pembimbing Internal" outlined dense
            placeholder="Nama Guru Pembimbing..." />
        </q-card-section>

        <q-card-actions align="right" class="text-primary q-pb-md q-pr-md">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn color="primary" label="Simpan Penempatan" @click="savePlotting" />
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
const selectedSiswa = ref([])
const dialogOpen = ref(false)
const isEdit = ref(false)

// Struktur Form Dialog
const form = ref({ id: null, siswa: null, namaSiswa: '', perusahaan: '', pembimbing: '' })

// Definisi Struktur Kolom Tabel q-table
const columns = [
  { name: 'nama', align: 'left', label: 'Nama Siswa', field: 'nama', sortable: true },
  { name: 'kelas', align: 'left', label: 'Kelas', field: 'kelas', sortable: true },
  { name: 'perusahaan', align: 'left', label: 'Perusahaan Mitra (DUDI)', field: 'perusahaan', sortable: true },
  { name: 'pembimbing', align: 'left', label: 'Guru Pembimbing', field: 'pembimbing' },
  { name: 'aksi', align: 'center', label: 'Cetak & Aksi' }
]

// Data Mock / Dummy Hubungan Penempatan PKL SMK
const dataPlotting = ref([
  { id: 1, nama: 'Ahmad Fauzi', kelas: 'XII RPL 1', jurusan: 'RPL', perusahaan: 'PT. Telekomunikasi Indonesia Tbk', kota: 'Bandung', pembimbing: 'Budi Rahardjo, S.T.' },
  { id: 2, nama: 'Budi Santoso', kelas: 'XII TKRO 3', jurusan: 'TKRO', perusahaan: 'PT. Astra Honda Motor', kota: 'Jakarta', pembimbing: 'Eko Sulistyo, M.Pd.' },
  { id: 3, nama: 'Chandra Wijaya', kelas: 'XI RPL 2', jurusan: 'RPL', perusahaan: 'PT. Solusi Teknologi Nusantara', kota: 'Cimahi', pembimbing: 'Budi Rahardjo, S.T.' },
  { id: 4, nama: 'Dinda Lestari', kelas: 'XI AKL 1', jurusan: 'AKL', perusahaan: 'Bank Mandiri (Persero) Tbk', kota: 'Bandung', pembimbing: 'Hj. Nenden, S.E.' }
])

// Opsi Data Dropdown untuk Dialog
const opsiSiswaBelumPkl = [
  { label: 'Eka Putri (X AKL 2)', value: { nama: 'Eka Putri', kelas: 'X AKL 2', jurusan: 'AKL' } },
  { label: 'Fajar Ramadhan (XII TKRO 1)', value: { nama: 'Fajar Ramadhan', kelas: 'XII TKRO 1', jurusan: 'TKRO' } }
]

const opsiPerusahaanMitra = [
  { label: 'PT. Pindad (Persero)', value: 'PT. Pindad (Persero)' },
  { label: 'PT. Kereta Api Indonesia', value: 'PT. Kereta Api Indonesia' },
  { label: 'Otoritas Jasa Keuangan (OJK)', value: 'Otoritas Jasa Keuangan (OJK)' }
]

// Logika Compute Filter Pencarian Di Sisi Klien (Frontend)
const filteredPlotting = computed(() => {
  return dataPlotting.value.filter(item => {
    const matchSearch = search.value
      ? item.nama.toLowerCase().includes(search.value.toLowerCase()) ||
      item.perusahaan.toLowerCase().includes(search.value.toLowerCase())
      : true

    const matchJurusan = filterJurusan.value !== 'Semua Jurusan'
      ? item.jurusan === filterJurusan.value
      : true

    return matchSearch && matchJurusan
  })
})

// Fungsi Penanganan Event Aksi
const openPlottingDialog = () => {
  isEdit.value = false
  form.value = { id: null, siswa: null, namaSiswa: '', perusahaan: '', pembimbing: '' }
  dialogOpen.value = true
}

const editPlotting = (row) => {
  isEdit.value = true
  form.value = {
    id: row.id,
    siswa: null,
    namaSiswa: row.nama,
    perusahaan: row.perusahaan,
    pembimbing: row.pembimbing
  }
  dialogOpen.value = true
}

const savePlotting = () => {
  if (isEdit.value) {
    // Logika simpan perubahan plotting
    const target = dataPlotting.value.find(item => item.id === form.value.id)
    if (target) {
      target.perusahaan = form.value.perusahaantarget.pembimbing = form.value.pembimbing
    } $q.notify({ color: 'green-8', message: 'Pembaruan lokasi PKL & Pembimbing berhasil disimpan!', icon: 'check' })
  }
  else {
    // Logika simpan plotting siswa baru
    if (!form.value.siswa) return
    dataPlotting.value.push({ id: Date.now(), nama: form.value.siswa.nama, kelas: form.value.siswa.kelas, jurusan: form.value.siswa.jurusan, perusahaan: form.value.perusahaan, kota: 'Bandung', pembimbing: form.value.pembimbing || 'Belum Ditentukan' })
    $q.notify({ color: 'primary', message: `Berhasil menempatkan ${form.value.siswa.nama} di lokasi PKL baru!`, icon: 'assignment_turned_in' })
  }
  dialogOpen.value = false
}
const printSingle = (row, jenisSurat) => {
  $q.notify({ color: 'indigo-8', message: `Sedang memproses generator cetak [${jenisSurat}]untuk siswa: ${row.nama}`, icon: 'print' })
}
const bulkPrint = (jenisSurat) => {
  if (selectedSiswa.value.length === 0) {
    $q.notify({ color: 'negative', message: 'Silakan pilih/centang nama siswa pada tabel terlebih dahulu!', icon: 'warning' })
    return
  }
  $q.notify({ color: 'green-9', message: `Memproses cetak massal berkas [${jenisSurat}]untuk ${selectedSiswa.value.length} siswa terpilih.`, icon: 'print' })
}
</script>

<template>
  <q-page class="q-pa-md bg-grey-2">
    <!-- Header Halaman -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Sertifikasi & Uji Kompetensi Keahlian (UKK)</h1>
        <p class="text-caption text-grey-7 q-mb-none">Pencatatan nomor registrasi, verifikasi berkas portfolio, dan
          status kompetensi keahlian siswa SMK</p>
      </div>
      <div>
        <q-btn color="primary" icon="assignment_turned_in" label="Input Reg. Massal" @click="openMassalDialog" />
      </div>
    </div>

    <!-- Banner Informasi Metrik UKK -->
    <div class="row q-col-gutter-sm q-mb-md">
      <div class="col-12 col-sm-3">
        <q-card class="bg-white text-dark flat shadow-1">
          <q-card-section class="q-py-sm flex justify-between items-center">
            <div>
              <div class="text-caption text-grey-7 text-uppercase">Total Peserta</div>
              <div class="text-h5 text-weight-bold text-primary">285 Siswa</div>
            </div>
            <q-avatar color="blue-1" text-color="primary" icon="person" size="40px" />
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-sm-3">
        <q-card class="bg-white text-dark flat shadow-1">
          <q-card-section class="q-py-sm flex justify-between items-center">
            <div>
              <div class="text-caption text-grey-7 text-uppercase">Berkas Lengkap</div>
              <div class="text-h5 text-weight-bold text-green-7">260 Siswa</div>
            </div>
            <q-avatar color="green-1" text-color="green-7" icon="fact_check" size="40px" />
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-sm-3">
        <q-card class="bg-white text-dark flat shadow-1">
          <q-card-section class="q-py-sm flex justify-between items-center">
            <div>
              <div class="text-caption text-grey-7 text-uppercase">Belum Verifikasi</div>
              <div class="text-h5 text-weight-bold text-orange-8">25 Siswa</div>
            </div>
            <q-avatar color="orange-1" text-color="orange-8" icon="pending_actions" size="40px" />
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-sm-3">
        <q-card class="bg-white text-dark flat shadow-1">
          <q-card-section class="q-py-sm flex justify-between items-center">
            <div>
              <div class="text-caption text-grey-7 text-uppercase">Sudah Kompeten</div>
              <div class="text-h5 text-weight-bold text-purple-7">145 Siswa</div>
            </div>
            <q-avatar color="purple-1" text-color="purple-7" icon="workspace_premium" size="40px" />
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Panel Filter Pencarian & Tabel -->
    <q-card class="shadow-1">
      <q-card-section class="q-py-sm">
        <div class="row q-col-gutter-sm items-center">
          <div class="col-12 col-md-5">
            <q-input v-model="search" dense outlined placeholder="Cari Nama Siswa, NISN, atau No. Peserta..." clearable>
              <template v-slot:prepend><q-icon name="search" /></template>
            </q-input>
          </div>
          <div class="col-12 col-md-4">
            <q-select v-model="filterSkema"
              :options="['Semua Skema', 'LSP-P1 (Badan Nasional Sertifikasi)', 'UKK Mandiri Dinas', 'Mitra Industri / Axioo']"
              dense outlined label="Skema Ujian" />
          </div>
          <div class="col-12 col-md-3 flex justify-end">
            <q-btn flat color="secondary" icon="download" label="Unduh Format Berkas" @click="downloadTemplate" />
          </div>
        </div>
      </q-card-section>

      <!-- Tabel Uji Kompetensi Keahlian -->
      <q-card-section class="q-pa-none">
        <q-table flat :rows="filteredUkk" :columns="columns" row-key="id" :pagination="{ rowsPerPage: 10 }">
          <!-- Kolom Nomor Registrasi / Keterangan Ujian -->
          <template v-slot:body-cell-noPeserta="props">
            <q-td :props="props">
              <div class="text-weight-bold text-primary">{{ props.row.noPeserta || 'Belum Digenerate' }}</div>
              <div class="text-caption text-grey-6">Skema: {{ props.row.skema }}</div>
            </q-td>
          </template>

          <!-- Kolom Status Verifikasi Administrasi Berkas -->
          <template v-slot:body-cell-berkas="props">
            <q-td :props="props" class="text-center">
              <q-chip clickable @click="toggleBerkas(props.row)" size="sm"
                :color="props.row.berkasLengkap ? 'green-1' : 'red-1'"
                :text-color="props.row.berkasLengkap ? 'green-9' : 'red-9'"
                :icon="props.row.berkasLengkap ? 'check_circle' : 'cancel'">
                {{ props.row.berkasLengkap ? 'Lengkap (Terverifikasi)' : 'Belum Lengkap' }}
              </q-chip>
            </q-td>
          </template>

          <!-- Kolom Hasil Penilaian Akhir -->
          <template v-slot:body-cell-hasil="props">
            <q-td :props="props" class="text-center">
              <q-badge
                :color="props.row.hasilUjian === 'Kompeten' ? 'purple' : props.row.hasilUjian === 'Belum Ujian' ? 'grey-6' : 'red-8'"
                :label="props.row.hasilUjian" />
            </q-td>
          </template>

          <!-- Kolom Aksi -->
          <template v-slot:body-cell-aksi="props">
            <q-td :props="props" class="text-center q-gutter-xs">
              <q-btn size="sm" color="primary" icon="edit_note" round flat @click="editSiswaUkk(props.row)">
                <q-tooltip>Input Nilai & No. Sertifikat</q-tooltip>
              </q-btn>
              <q-btn size="sm" color="green-8" icon="print" round flat :disable="props.row.hasilUjian !== 'Kompeten'"
                @click="printSertifikat(props.row)">
                <q-tooltip>Cetak Lembar Transkrip UKK</q-tooltip>
              </q-btn>
            </q-td>
          </template>
        </q-table>
      </q-card-section>
    </q-card>

    <!-- DIALOG POPUP: Input Registrasi Massal / Kolektif -->
    <q-dialog v-model="dialogMassalOpen" persistent>
      <q-card style="min-width: 400px;">
        <q-card-section class="bg-primary text-white">
          <div class="text-h6">Input Nomor Ujian Massal</div>
          <div class="text-caption text-blue-2">Menyisipkan kode registrasi untuk siswa jurusan terpilih</div>
        </q-card-section>

        <q-card-section class="q-gutter-md q-pt-md">
          <q-select v-model="formMassal.jurusan" :options="['RPL', 'TKRO', 'AKL']" label="Pilih Kompetensi Keahlian"
            outlined dense />
          <q-input v-model="formMassal.prefix" label="Prefix / Awalan Nomor" outlined dense
            placeholder="Contoh: UKK-LSP-2026-" />
          <p class="text-caption text-amber-9 bg-amber-1 q-pa-sm rounded-borders">
            <q-icon name="info" /> Nomor ujian otomatis dibuat berurutan berdasarkan nomor urut Buku Klaper/Induk siswa.
          </p>
        </q-card-section>

        <q-card-actions align="right" class="q-pb-md q-pr-md">
          <q-btn flat label="Batal" v-close-popup color="grey" />
          <q-btn color="primary" label="Generate Nomor" @click="processGenerateMassal" />
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
const filterSkema = ref('Semua Skema')
const dialogMassalOpen = ref(false)

const formMassal = ref({ jurusan: '', prefix: '' })

// Struktur Kolom Tabel q-table
const columns = [
  { name: 'nama', align: 'left', label: 'Nama Siswa', field: 'nama', sortable: true },
  { name: 'kelas', align: 'left', label: 'Kelas', field: 'kelas', sortable: true },
  { name: 'noPeserta', align: 'left', label: 'No. Registrasi Peserta', field: 'noPeserta' },
  { name: 'berkas', align: 'center', label: 'Verifikasi Administrasi', field: 'berkasLengkap' },
  { name: 'hasil', align: 'center', label: 'Hasil Kelulusan', field: 'hasilUjian' },
  { name: 'aksi', align: 'center', label: 'Aksi' }
]

// Data Mock / Dummy Sertifikasi Kelulusan Ujian SMK
const dataUkk = ref([
  { id: 1, nama: 'Ahmad Fauzi', kelas: 'XII RPL 1', skema: 'LSP-P1 (Badan Nasional Sertifikasi)', noPeserta: 'REG-LSP-2026001', berkasLengkap: true, hasilUjian: 'Kompeten' },
  { id: 2, nama: 'Budi Santoso', kelas: 'XII TKRO 3', skema: 'UKK Mandiri Dinas', noPeserta: 'REG-DND-2026084', berkasLengkap: true, hasilUjian: 'Kompeten' },
  { id: 3, nama: 'Anisa Rahmawati', kelas: 'XII AKL 2', skema: 'LSP-P1 (Badan Nasional Sertifikasi)', noPeserta: 'REG-LSP-2026042', berkasLengkap: false, hasilUjian: 'Belum Ujian' },
  { id: 4, nama: 'Fajar Ramadhan', kelas: 'XII TKRO 1', skema: 'UKK Mandiri Dinas', noPeserta: '', berkasLengkap: false, hasilUjian: 'Belum Ujian' }
])

// Logika Compute Filter Pencarian Frontend
const filteredUkk = computed(() => {
  return dataUkk.value.filter(item => {
    const matchSearch = search.value
      ? item.nama.toLowerCase().includes(search.value.toLowerCase()) ||
      item.noPeserta.includes(search.value)
      : true

    const matchSkema = filterSkema.value !== 'Semua Skema'
      ? item.skema.includes(filterSkema.value.split(' ')[0]) // Mencocokkan kata depan seperti "LSP-P1" atau "UKK"
      : true

    return matchSearch && matchSkema
  })
})

// Fungsi Penanganan Event Aksi Kesiswaan
const toggleBerkas = (row) => {
  row.berkasLengkap = !row.berkasLengkap
  $q.notify({
    color: row.berkasLengkap ? 'green-8' : 'orange-9',
    message: `Status berkas ${row.nama} diubah menjadi: ${row.berkasLengkap ? 'LENGKAP' : 'BELUM LENGKAP'}`,
    icon: 'sync'
  })
}

const openMassalDialog = () => {
  formMassal.value = { jurusan: '', prefix: 'UKK-SMK-2026-' }
  dialogMassalOpen.value = true
}
const processGenerateMassal = () => {
  if (!formMassal.value.jurusan) {
    $q.notify({ color: 'negative', message: 'Silakan pilih jurusan terlebih dahulu!', icon: 'warning' })
    return
  }
  // Cari siswa yang nomor pesertanya masih kosong dan sesuai jurusan pilihan
  let count = 0
  dataUkk.value.forEach((siswa, index) => {
    if (siswa.kelas.includes(formMassal.value.jurusan) && !siswa.noPeserta) {
      count++
      siswa.noPeserta = `${formMassal.value.prefix}00${index + 1}`
    }
  })
  $q.notify({
    color: 'green-8', message: `Berhasil men - generate nomor registrasi ujian otomatis untuk ${count} siswa
    ${formMassal.value.jurusan}.`, icon: 'done_all'
  })
  dialogMassalOpen.value = false
}
const editSiswaUkk = (row) => {
  $q.dialog({ title: 'Input Hasil Kelulusan UKK', message: `Masukkan status kelulusan kompetensi untuk ${row.nama}:`, options: { type: 'radio', model: row.hasilUjian, items: [{ label: 'Belum Ujian', value: 'Belum Ujian' }, { label: 'Kompeten (Lulus)', value: 'Kompeten' }, { label: 'Belum Kompeten (Mengulang)', value: 'Belum Kompeten' }] }, cancel: true, persistent: true }).onOk(data => {
    row.hasilUjian = data
    $q.notify({ color: 'primary', message: `Status kelulusan ${row.nama} diperbarui menjadi: ${data}`, icon: 'save' })
  })
}
const printSertifikat = (row) => { $q.notify({ color: 'purple-8', message: `Mengunduh Transkrip Nilai Kompetensi & sertifikat pendamping untuk ${row.nama}`, icon: 'download' }) }
const downloadTemplate = () => { $q.notify({ color: 'grey-8', message: 'Mengunduh format lampiran F-02 portofolio pendaftaran LSP-P1.', icon: 'file_download' }) }
</script>

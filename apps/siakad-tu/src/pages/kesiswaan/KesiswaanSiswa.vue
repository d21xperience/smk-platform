<template>
  <q-page class="q-pa-md">
    <div class="row justify-between q-mb-md">
      <div class="text-h5">Data Induk Siswa</div>
      <div class="q-gutter-sm">
        <q-btn color="green-7" icon="table_view" label="Export Excel" @click="openExportExcel" />
        <q-btn color="deep-orange-7" icon="print" label="Cetak PDF" @click="openCetakPdf" />
        <q-btn color="primary" icon="add" label="Tambah Siswa" @click="openForm" />
      </div>
    </div>

    <q-table :rows="siswaStore.list" :columns="columns" row-key="id" flat bordered dense>
      <template v-slot:body-cell-aksi="props">
        <q-td :props="props">
          <q-btn flat round dense icon="edit" color="info" @click="openForm(props.row)" class="q-mr-sm" />
          <q-btn flat round dense icon="history" color="secondary" @click="showMutasi(props.row)" class="q-mr-sm">
            <q-tooltip>Riwayat Mutasi</q-tooltip>
          </q-btn>
          <q-btn flat round dense icon="delete" color="negative" @click="confirmDelete(props.row.id)" />
        </q-td>
      </template>
    </q-table>

    <!-- ==================== DIALOG FORM SISWA ==================== -->
    <q-dialog v-model="formDialog" persistent>
      <q-card style="min-width: 700px">
        <q-card-section>
          <div class="text-h6">{{ editMode ? 'Edit Siswa' : 'Tambah Siswa' }}</div>
        </q-card-section>
        <q-card-section>
          <q-form @submit="submitForm" class="q-gutter-md">
            <div class="row q-col-gutter-md">
              <div class="col-6">
                <q-input v-model="form.nis" label="NIS" outlined :rules="[v => !!v]" />
              </div>
              <div class="col-6">
                <q-input v-model="form.nama" label="Nama Lengkap" outlined :rules="[v => !!v]" />
              </div>
              <div class="col-6">
                <q-input v-model="form.tempat_lahir" label="Tempat Lahir" outlined />
              </div>
              <div class="col-6">
                <q-input v-model="form.tanggal_lahir" label="Tanggal Lahir" type="date" outlined />
              </div>
              <div class="col-6">
                <q-select v-model="form.jenis_kelamin" :options="['Laki-laki', 'Perempuan']" label="Jenis Kelamin"
                  outlined />
              </div>
              <div class="col-6">
                <q-input v-model="form.no_hp" label="No. HP Siswa" outlined />
              </div>
              <div class="col-12">
                <q-input v-model="form.alamat" label="Alamat" type="textarea" rows="2" outlined />
              </div>
              <div class="col-12">
                <q-select v-model="form.kelas_id" :options="kelasOptions" label="Kelas" option-value="id"
                  option-label="nama" outlined />
              </div>
              <div class="col-12">
                <div class="text-subtitle1">Data Orang Tua</div>
              </div>
              <div class="col-6">
                <q-input v-model="form.orang_tua.ayah" label="Nama Ayah" outlined />
              </div>
              <div class="col-6">
                <q-input v-model="form.orang_tua.ibu" label="Nama Ibu" outlined />
              </div>
              <div class="col-6">
                <q-input v-model="form.orang_tua.pekerjaan_ayah" label="Pekerjaan Ayah" outlined />
              </div>
              <div class="col-6">
                <q-input v-model="form.orang_tua.pekerjaan_ibu" label="Pekerjaan Ibu" outlined />
              </div>
              <div class="col-12">
                <q-input v-model="form.orang_tua.no_hp_ortu" label="No. HP Orang Tua" outlined />
              </div>
              <div class="col-12">
                <q-select v-model="form.status" :options="['aktif', 'pindah', 'keluar']" label="Status" outlined />
              </div>
            </div>
            <div class="q-gutter-sm">
              <q-btn label="Simpan" type="submit" color="primary" />
              <q-btn label="Batal" flat color="negative" v-close-popup @click="resetForm" />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- ==================== DIALOG RIWAYAT MUTASI ==================== -->
    <q-dialog v-model="mutasiDialog">
      <q-card style="min-width: 600px">
        <q-card-section>
          <div class="text-h6">Riwayat Mutasi Siswa</div>
        </q-card-section>
        <q-card-section>
          <q-table :rows="mutasiList" :columns="mutasiColumns" dense flat />
          <q-btn label="Tambah Mutasi" icon="add" flat color="primary" @click="openMutasiForm" class="q-mt-md" />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Tutup" color="primary" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- ==================== DIALOG FORM MUTASI ==================== -->
    <q-dialog v-model="mutasiFormDialog" persistent>
      <q-card style="min-width: 500px">
        <q-card-section>
          <div class="text-h6">Tambah Riwayat Mutasi</div>
        </q-card-section>
        <q-card-section>
          <q-form @submit="submitMutasi" class="q-gutter-md">
            <q-input v-model="mutasiForm.tanggal" label="Tanggal Mutasi" type="date" outlined />
            <q-input v-model="mutasiForm.asal_sekolah" label="Asal Sekolah / Tujuan" outlined />
            <q-input v-model="mutasiForm.alasan" label="Alasan" type="textarea" rows="2" outlined />
            <q-file v-model="mutasiForm.dokumen" label="Dokumen (PDF/Image)" accept=".pdf,.jpg,.png" outlined />
            <div class="q-gutter-sm">
              <q-btn label="Simpan" type="submit" color="primary" />
              <q-btn label="Batal" flat color="negative" v-close-popup />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- ==================== DIALOG EXPORT EXCEL ==================== -->
    <q-dialog v-model="exportExcelDialog" persistent>
      <q-card style="min-width: 500px">
        <q-card-section class="bg-green-7 text-white">
          <div class="text-h6"><q-icon name="table_view" class="q-mr-sm" />Export Data ke Excel</div>
        </q-card-section>
        <q-card-section>
          <div class="text-subtitle2 q-mb-sm">Pilih Kolom yang Ingin Diekspor:</div>
          <div class="row q-col-gutter-sm">
            <div class="col-6" v-for="col in excelColumnOptions" :key="col.field">
              <q-checkbox v-model="selectedExcelColumns" :val="col.field" :label="col.label" dense />
            </div>
          </div>
          <q-separator class="q-my-md" />
          <div class="row q-col-gutter-sm">
            <div class="col-12">
              <q-input v-model="excelFileName" label="Nama File Excel" outlined dense hint="Tanpa ekstensi .xlsx" />
            </div>
          </div>
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" color="negative" v-close-popup />
          <q-btn label="Export" icon="download" color="green-7" :disable="selectedExcelColumns.length === 0"
            @click="doExportExcel" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- ==================== DIALOG CETAK PDF ==================== -->
    <q-dialog v-model="cetakPdfDialog" persistent>
      <q-card style="min-width: 600px">
        <q-card-section class="bg-deep-orange-7 text-white">
          <div class="text-h6"><q-icon name="print" class="q-mr-sm" />Pengaturan Cetak PDF</div>
        </q-card-section>
        <q-card-section>
          <div class="text-subtitle2 q-mb-sm text-grey-7">Header Instansi</div>
          <div class="row q-col-gutter-md q-mb-md">
            <div class="col-12">
              <q-input v-model="printConfig.nama_instansi" label="Nama Sekolah / Instansi" outlined dense />
            </div>
            <div class="col-12">
              <q-input v-model="printConfig.alamat_instansi" label="Alamat Instansi" outlined dense type="textarea"
                rows="2" />
            </div>
            <div class="col-6">
              <q-input v-model="printConfig.telp_instansi" label="No. Telp / Fax" outlined dense />
            </div>
            <div class="col-6">
              <q-input v-model="printConfig.website_instansi" label="Website / Email" outlined dense />
            </div>
          </div>

          <q-separator class="q-mb-md" />
          <div class="text-subtitle2 q-mb-sm text-grey-7">Kustomisasi Dokumen</div>
          <div class="row q-col-gutter-md q-mb-md">
            <div class="col-12">
              <q-input v-model="printConfig.judul_dokumen" label="Judul Dokumen" outlined dense
                placeholder="Contoh: Daftar Hadir Siswa" />
            </div>
            <div class="col-6">
              <q-input v-model="printConfig.nama_kelas" label="Nama Kelas" outlined dense
                placeholder="Contoh: XII RPL 1" />
            </div>
            <div class="col-6">
              <q-input v-model="printConfig.nama_wali_kelas" label="Nama Wali Kelas" outlined dense />
            </div>
            <div class="col-6">
              <q-input v-model="printConfig.tahun_ajaran" label="Tahun Ajaran" outlined dense
                placeholder="Contoh: 2024/2025" />
            </div>
            <div class="col-6">
              <q-input v-model="printConfig.semester" label="Semester" outlined dense
                placeholder="Contoh: Ganjil / Genap" />
            </div>
            <div class="col-12">
              <q-input v-model="printConfig.catatan" label="Catatan Tambahan (opsional)" outlined dense type="textarea"
                rows="2" />
            </div>
          </div>

          <q-separator class="q-mb-md" />
          <div class="text-subtitle2 q-mb-sm text-grey-7">Pilih Kolom yang Dicetak:</div>
          <div class="row q-col-gutter-sm">
            <div class="col-6" v-for="col in printColumnOptions" :key="col.field">
              <q-checkbox v-model="selectedPrintColumns" :val="col.field" :label="col.label" dense />
            </div>
          </div>
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" color="negative" v-close-popup />
          <q-btn label="Cetak" icon="print" color="deep-orange-7" :disable="selectedPrintColumns.length === 0"
            @click="doCetakPdf" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useSiswaStore } from 'stores/kesiswaan/siswa'
import { useKelasStore } from '@/stores/kesiswaan/kelasStore'
import * as XLSX from 'xlsx'

const $q = useQuasar()
const siswaStore = useSiswaStore()
const kelasStore = useKelasStore()

// ==================== STATE EXISTING ====================
const formDialog = ref(false)
const editMode = ref(false)
const editId = ref(null)
const form = ref({
  nis: '', nama: '', tempat_lahir: '', tanggal_lahir: '', jenis_kelamin: '', alamat: '', no_hp: '', kelas_id: null,
  orang_tua: { ayah: '', ibu: '', pekerjaan_ayah: '', pekerjaan_ibu: '', no_hp_ortu: '' },
  status: 'aktif'
})
const mutasiDialog = ref(false)
const currentSiswaId = ref(null)
const mutasiList = ref([])
const mutasiFormDialog = ref(false)
const mutasiForm = ref({ tanggal: '', asal_sekolah: '', alasan: '', dokumen: null })

const columns = [
  { name: 'nis', label: 'NIS', field: 'nis' },
  { name: 'nama', label: 'Nama', field: 'nama' },
  { name: 'kelas_nama', label: 'Kelas', field: 'kelas_nama' },
  { name: 'status', label: 'Status', field: 'status' },
  { name: 'aksi', label: 'Aksi', field: 'aksi' }
]
const mutasiColumns = [
  { name: 'tanggal', label: 'Tanggal', field: 'tanggal' },
  { name: 'asal_sekolah', label: 'Asal/Tujuan', field: 'asal_sekolah' },
  { name: 'alasan', label: 'Alasan', field: 'alasan' }
]
const kelasOptions = computed(() => kelasStore.list)

// ==================== STATE EXPORT EXCEL ====================
const exportExcelDialog = ref(false)
const excelFileName = ref('Data_Siswa')

// Semua kolom yang tersedia untuk export
const excelColumnOptions = [
  { field: 'nis', label: 'NIS' },
  { field: 'nama', label: 'Nama Lengkap' },
  { field: 'kelas_nama', label: 'Kelas' },
  { field: 'tempat_lahir', label: 'Tempat Lahir' },
  { field: 'tanggal_lahir', label: 'Tanggal Lahir' },
  { field: 'jenis_kelamin', label: 'Jenis Kelamin' },
  { field: 'alamat', label: 'Alamat' },
  { field: 'no_hp', label: 'No. HP Siswa' },
  { field: 'status', label: 'Status' },
  { field: 'orang_tua.ayah', label: 'Nama Ayah' },
  { field: 'orang_tua.ibu', label: 'Nama Ibu' },
  { field: 'orang_tua.pekerjaan_ayah', label: 'Pekerjaan Ayah' },
  { field: 'orang_tua.pekerjaan_ibu', label: 'Pekerjaan Ibu' },
  { field: 'orang_tua.no_hp_ortu', label: 'No. HP Orang Tua' },
]

// Kolom default yang sudah tercentang
const selectedExcelColumns = ref(['nis', 'nama', 'kelas_nama', 'jenis_kelamin', 'status'])

// ==================== STATE CETAK PDF ====================
const cetakPdfDialog = ref(false)

const printConfig = ref({
  nama_instansi: 'SMK Pasundan Jatinangor',
  alamat_instansi: 'Jl. Raya Jatinangor No. 1, Kab. Sumedang, Jawa Barat',
  telp_instansi: '(022) 12345678',
  website_instansi: 'www.smkpasundan.sch.id',
  judul_dokumen: 'Daftar Siswa',
  nama_kelas: '',
  nama_wali_kelas: '',
  tahun_ajaran: '',
  semester: '',
  catatan: ''
})

const printColumnOptions = [
  { field: 'no', label: 'No.' },
  { field: 'nis', label: 'NIS' },
  { field: 'nama', label: 'Nama Lengkap' },
  { field: 'kelas_nama', label: 'Kelas' },
  { field: 'tempat_lahir', label: 'Tempat Lahir' },
  { field: 'tanggal_lahir', label: 'Tanggal Lahir' },
  { field: 'jenis_kelamin', label: 'Jenis Kelamin' },
  { field: 'alamat', label: 'Alamat' },
  { field: 'no_hp', label: 'No. HP' },
  { field: 'status', label: 'Status' },
  { field: 'orang_tua.no_hp_ortu', label: 'No. HP Ortu' },
]

const selectedPrintColumns = ref(['no', 'nis', 'nama', 'kelas_nama', 'jenis_kelamin', 'status'])

// ==================== FUNGSI EXISTING ====================
function openForm(row = null) {
  if (row) {
    editMode.value = true
    editId.value = row.id
    form.value = JSON.parse(JSON.stringify(row))
  } else {
    editMode.value = false
    resetForm()
  }
  formDialog.value = true
}
function resetForm() {
  form.value = {
    nis: '', nama: '', tempat_lahir: '', tanggal_lahir: '', jenis_kelamin: '', alamat: '', no_hp: '', kelas_id: null,
    orang_tua: { ayah: '', ibu: '', pekerjaan_ayah: '', pekerjaan_ibu: '', no_hp_ortu: '' },
    status: 'aktif'
  }
}
function submitForm() {
  const kelas = kelasStore.list.find(k => k.id === form.value.kelas_id)
  const data = { ...form.value, kelas_nama: kelas?.nama }
  if (editMode.value) {
    siswaStore.update(editId.value, data)
    $q.notify({ type: 'positive', message: 'Data siswa diupdate' })
  } else {
    siswaStore.tambah(data)
    $q.notify({ type: 'positive', message: 'Siswa ditambahkan' })
  }
  formDialog.value = false
}
function confirmDelete(id) {
  $q.dialog({ title: 'Konfirmasi', message: 'Hapus siswa ini?', cancel: true }).onOk(() => {
    siswaStore.hapus(id)
    $q.notify({ type: 'positive', message: 'Siswa dihapus' })
  })
}
function showMutasi(siswa) {
  currentSiswaId.value = siswa.id
  mutasiList.value = siswa.riwayat_mutasi || []
  mutasiDialog.value = true
}
function openMutasiForm() {
  mutasiForm.value = { tanggal: '', asal_sekolah: '', alasan: '', dokumen: null }
  mutasiFormDialog.value = true
}
function submitMutasi() {
  let dokumenBase64 = ''
  if (mutasiForm.value.dokumen) {
    const reader = new FileReader()
    reader.onload = (e) => { dokumenBase64 = e.target.result }
    reader.readAsDataURL(mutasiForm.value.dokumen)
  }
  const mutasi = {
    tanggal: mutasiForm.value.tanggal,
    asal_sekolah: mutasiForm.value.asal_sekolah,
    alasan: mutasiForm.value.alasan,
    dokumen: dokumenBase64 || mutasiForm.value.dokumen?.name
  }
  siswaStore.tambahMutasi(currentSiswaId, mutasi)
  $q.notify({ type: 'positive', message: 'Riwayat mutasi ditambahkan' })
  mutasiFormDialog.value = false
  const siswa = siswaStore.list.find(s => s.id === currentSiswaId.value)
  mutasiList.value = siswa.riwayat_mutasi
}

// ==================== FUNGSI EXPORT EXCEL ====================
function openExportExcel() {
  exportExcelDialog.value = true
}

/**
 * Ambil nilai nested dari objek menggunakan dot-notation
 * Contoh: getNestedValue(row, 'orang_tua.ayah') → row.orang_tua.ayah
 */
function getNestedValue(obj, path) {
  return path.split('.').reduce((acc, key) => acc?.[key] ?? '', obj)
}

function doExportExcel() {
  const activeColumns = excelColumnOptions.filter(c => selectedExcelColumns.value.includes(c.field))

  // Buat array of objects sesuai kolom terpilih
  const exportData = siswaStore.list.map((row, index) => {
    const entry = { 'No.': index + 1 }
    activeColumns.forEach(col => {
      entry[col.label] = getNestedValue(row, col.field)
    })
    return entry
  })

  const ws = XLSX.utils.json_to_sheet(exportData)
  const wb = XLSX.utils.book_new()
  XLSX.utils.book_append_sheet(wb, ws, 'Data Siswa')

  // Auto-width kolom berdasarkan konten
  const colWidths = Object.keys(exportData[0] || {}).map(key => ({
    wch: Math.max(key.length, ...exportData.map(r => String(r[key] ?? '').length)) + 2
  }))
  ws['!cols'] = colWidths

  XLSX.writeFile(wb, `${excelFileName.value || 'Data_Siswa'}.xlsx`)
  exportExcelDialog.value = false
  $q.notify({ type: 'positive', message: 'File Excel berhasil diunduh', icon: 'check_circle' })
}

// ==================== FUNGSI CETAK PDF ====================
function openCetakPdf() {
  cetakPdfDialog.value = true
}

function doCetakPdf() {
  const cfg = printConfig.value
  const activeCols = printColumnOptions.filter(c => selectedPrintColumns.value.includes(c.field))
  const data = siswaStore.list

  // Buat tanggal cetak
  const tglCetak = new Date().toLocaleDateString('id-ID', {
    day: '2-digit', month: 'long', year: 'numeric'
  })

  // Bangun baris tabel
  const headerRow = activeCols.map(c => `<th>${c.label}</th>`).join('')
  const bodyRows = data.map((row, i) =>
    `<tr>${activeCols.map(c => {
      if (c.field === 'no') return `<td style="text-align:center">${i + 1}</td>`
      return `<td>${getNestedValue(row, c.field) || '-'}</td>`
    }).join('')}</tr>`
  ).join('')

  // Bangun info kustomisasi
  const infoKustom = [
    cfg.nama_kelas ? `<tr><td>Kelas</td><td>: ${cfg.nama_kelas}</td></tr>` : '',
    cfg.nama_wali_kelas ? `<tr><td>Wali Kelas</td><td>: ${cfg.nama_wali_kelas}</td></tr>` : '',
    cfg.tahun_ajaran ? `<tr><td>Tahun Ajaran</td><td>: ${cfg.tahun_ajaran}</td></tr>` : '',
    cfg.semester ? `<tr><td>Semester</td><td>: ${cfg.semester}</td></tr>` : '',
  ].filter(Boolean).join('')

  const printHTML = `
    <!DOCTYPE html>
    <html lang="id">
    <head>
      <meta charset="UTF-8">
      <title>${cfg.judul_dokumen}</title>
      <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: Arial, sans-serif; font-size: 11px; color: #000; padding: 20px; }

        /* ===== HEADER INSTANSI ===== */
        .header-instansi {
          display: flex;
          align-items: center;
          border-bottom: 3px solid #000;
          padding-bottom: 10px;
          margin-bottom: 8px;
        }
        .header-instansi .logo {
          width: 70px;
          height: 70px;
          margin-right: 16px;
          object-fit: contain;
        }
        .header-instansi .logo-placeholder {
          width: 70px;
          height: 70px;
          margin-right: 16px;
          border: 2px solid #555;
          display: flex;
          align-items: center;
          justify-content: center;
          font-size: 9px;
          color: #555;
          text-align: center;
          flex-shrink: 0;
        }
        .header-info { flex: 1; }
        .header-info .nama-instansi {
          font-size: 16px;
          font-weight: bold;
          text-transform: uppercase;
          letter-spacing: 0.5px;
        }
        .header-info .detail-instansi {
          font-size: 10px;
          line-height: 1.6;
          color: #333;
          margin-top: 2px;
        }

        /* ===== JUDUL DOKUMEN ===== */
        .judul-dokumen {
          text-align: center;
          margin: 14px 0 10px;
        }
        .judul-dokumen h2 {
          font-size: 14px;
          font-weight: bold;
          text-transform: uppercase;
          letter-spacing: 1px;
        }
        .judul-dokumen .garis-judul {
          width: 60px;
          height: 2px;
          background: #000;
          margin: 4px auto 0;
        }

        /* ===== INFO KUSTOM ===== */
        .info-kustom {
          margin-bottom: 12px;
          font-size: 11px;
        }
        .info-kustom table { border-collapse: collapse; }
        .info-kustom td { padding: 1px 8px 1px 0; }

        /* ===== TABEL DATA ===== */
        table.tabel-data {
          width: 100%;
          border-collapse: collapse;
          font-size: 10px;
          margin-top: 8px;
        }
        table.tabel-data th {
          background-color: #2c2c2c;
          color: #fff;
          padding: 6px 8px;
          text-align: left;
          border: 1px solid #555;
        }
        table.tabel-data td {
          padding: 5px 8px;
          border: 1px solid #ccc;
          vertical-align: top;
        }
        table.tabel-data tr:nth-child(even) td {
          background-color: #f7f7f7;
        }

        /* ===== FOOTER ===== */
        .footer-cetak {
          margin-top: 16px;
          font-size: 10px;
          color: #555;
          display: flex;
          justify-content: space-between;
          align-items: flex-end;
        }
        .ttd-box { text-align: center; }
        .ttd-box .ttd-nama { margin-top: 50px; font-weight: bold; border-top: 1px solid #000; padding-top: 4px; }
        .ttd-box .ttd-jabatan { font-size: 9px; color: #555; }
        .catatan-box { font-style: italic; max-width: 50%; }

        @media print {
          body { padding: 0; }
          @page { margin: 1.5cm; size: A4 landscape; }
        }
      </style>
    </head>
    <body>
      <!-- HEADER INSTANSI -->
      <div class="header-instansi">
        <div class="logo-placeholder">LOGO<br>INSTANSI</div>
        <div class="header-info">
          <div class="nama-instansi">${cfg.nama_instansi}</div>
          <div class="detail-instansi">
            ${cfg.alamat_instansi}<br>
            Telp/Fax: ${cfg.telp_instansi} &nbsp;|&nbsp; ${cfg.website_instansi}
          </div>
        </div>
      </div>

      <!-- JUDUL DOKUMEN -->
      <div class="judul-dokumen">
        <h2>${cfg.judul_dokumen}</h2>
        <div class="garis-judul"></div>
      </div>

      <!-- INFO KUSTOM (kelas, wali kelas, dll.) -->
      ${infoKustom ? `<div class="info-kustom"><table>${infoKustom}</table></div>` : ''}

      <!-- TABEL DATA -->
      <table class="tabel-data">
        <thead><tr>${headerRow}</tr></thead>
        <tbody>${bodyRows}</tbody>
      </table>

      <!-- FOOTER -->
      <div class="footer-cetak">
        <div class="catatan-box">
          ${cfg.catatan ? `<strong>Catatan:</strong><br>${cfg.catatan}` : ''}
        </div>
        <div>
          <div>Dicetak pada: ${tglCetak}</div>
          <br>
          <div class="ttd-box">
            ${cfg.nama_wali_kelas
      ? `<div class="ttd-jabatan">Wali Kelas ${cfg.nama_kelas}</div>
                 <div class="ttd-nama">${cfg.nama_wali_kelas}</div>`
      : `<div class="ttd-jabatan">Mengetahui,</div>
                 <div class="ttd-nama">(____________________)</div>`
    }
          </div>
        </div>
      </div>

    </body>
    </html>
  `

  // Buka window print baru
  const printWindow = window.open('', '_blank', 'width=1024,height=768')
  printWindow.document.write(printHTML)
  printWindow.document.close()
  printWindow.focus()
  setTimeout(() => {
    printWindow.print()
    // Tutup window setelah print dialog ditutup (opsional)
    // printWindow.close()
  }, 500)

  cetakPdfDialog.value = false
}

onMounted(() => {
  siswaStore.loadData()
  kelasStore.loadData()
})
</script>

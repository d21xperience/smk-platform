<template>
  <q-btn color="deep-orange-7" icon="print" label="Cetak PDF" @click="openCetakPdf" />
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
</template>

<script setup>
import { ref } from 'vue'
const props = defineProps({
  siswa: { type: Array, required: true }
})
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

function openCetakPdf() {
  cetakPdfDialog.value = true
}

/**
 * Ambil nilai nested dari objek menggunakan dot-notation
 * Contoh: getNestedValue(row, 'orang_tua.ayah') → row.orang_tua.ayah
 */
function getNestedValue(obj, path) {
  return path.split('.').reduce((acc, key) => acc?.[key] ?? '', obj)
}

function doCetakPdf() {
  const cfg = printConfig.value
  const activeCols = printColumnOptions.filter(c => selectedPrintColumns.value.includes(c.field))
  const data = props.siswa.list

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
</script>

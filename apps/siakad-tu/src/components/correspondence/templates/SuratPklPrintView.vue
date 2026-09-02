<template>
  <div class="print-area a4-paper">
    <!-- KOP SURAT (Mock Data, nanti diambil dari School Profile) -->
    <div class="kop-surat">
      <div class="logo">
        <!-- <img src="/logo-sekolah.png" alt="Logo" /> -->
        <div class="fake-logo">LOGO</div>
      </div>
      <div class="text-kop">
        <h1>PEMERINTAH PROVINSI / YAYASAN</h1>
        <h2>SEKOLAH MENENGAH KEJURUAN (SMK) - 01</h2>
        <p>Jl. Pendidikan No. 123, Kota Contoh, Kode Pos 12345</p>
        <p>Telp: (021) 1234567 | Email: info@smk01.sch.id | Website: www.smk01.sch.id</p>
      </div>
    </div>
    <hr class="garis-kop" />

    <!-- ISI SURAT -->
    <div class="isi-surat">
      <div class="header-surat">
        <table class="no-border">
          <tbody>
            <tr>
              <td width="150">Nomor</td>
              <td width="10">:</td>
              <td>{{ letter.nomorSurat }}</td>
            </tr>
            <tr>
              <td>Lampiran</td>
              <td>:</td>
              <td>1 Berkas</td>
            </tr>
            <tr>
              <td>Perihal</td>
              <td>:</td>
              <td><strong>Pengantar Praktek Kerja Lapangan (PKL)</strong></td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="tujuan">
        <p>Kepada Yth.<br />
          Pimpinan <strong>{{ letter.templateData?.dudi }}</strong><br />
          di<br />
          Tempat</p>
      </div>

      <div class="pembuka">
        <p>Dengan hormat,</p>
        <p>Sebagai bagian dari kurikulum Sekolah Menengah Kejuruan (SMK), kami bermaksud menugaskan peserta didik kami
          untuk melaksanakan Praktek Kerja Lapangan (PKL) di instansi/perusahaan yang Bapak/Ibu pimpin.</p>
        <p>Adapun nama-nama peserta didik yang akan melaksanakan PKL adalah sebagai berikut:</p>
      </div>

      <!-- TABEL SISWA -->
      <table class="tabel-siswa" border="1" cellpadding="5" cellspacing="0">
        <thead>
          <tr>
            <th width="40" align="center">No</th>
            <th align="center">Nama Lengkap Siswa</th>
            <th width="100" align="center">Kelas</th>
            <th width="150" align="center">No. Telepon / HP</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(siswa, index) in letter.templateData?.siswaList" :key="index">
            <td align="center">{{ index + 1 }}</td>
            <td>{{ siswa.nama }}</td>
            <td align="center">{{ siswa.kelas }}</td>
            <td>{{ siswa.noHp }}</td>
          </tr>
        </tbody>
      </table>

      <div class="penutup">
        <p>Demikian surat pengantar ini kami sampaikan. Atas perhatian dan kerjasama yang baik, kami ucapkan terima
          kasih.</p>
      </div>

      <!-- TANDA TANGAN -->
      <div class="tanda-tangan">
        <div class="tempat-tanggal">
          Kota Contoh, {{ formatTanggal(letter.tanggal) }}
        </div>
        <div class="jabatan">Kepala Sekolah,</div>

        <!-- QR Code / Tanda Tangan Digital bisa ditaruh di sini nanti -->
        <div class="nama-ttd">
          <strong><u>Dr. Nama Kepala Sekolah, M.Pd.</u></strong><br />
          NIP. 196501011990011001
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  letter: {
    type: Object,
    required: true
  }
})

function formatTanggal(dateStr) {
  if (!dateStr) return ''
  const bulan = ['Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni', 'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember']
  const d = new Date(dateStr)
  return `${d.getDate()} ${bulan[d.getMonth()]} ${d.getFullYear()}`
}
</script>

<style scoped>
/* Layout A4 Presisi */
.a4-paper {
  width: 210mm;
  min-height: 297mm;
  padding: 15mm 20mm 20mm 25mm;
  /* Margin kiri lebih besar untuk binding */
  margin: 0 auto;
  background: white;
  box-shadow: 0 0 15px rgba(0, 0, 0, 0.1);
  font-family: 'Times New Roman', Times, serif;
  font-size: 12pt;
  line-height: 1.5;
  color: black;
}

/* Kop Surat */
.kop-surat {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
}

.logo {
  margin-right: 15px;
}

.fake-logo {
  width: 80px;
  height: 80px;
  background: #ddd;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
}

.text-kop {
  text-align: center;
  flex: 1;
}

.text-kop h1 {
  margin: 0;
  font-size: 14pt;
  font-weight: bold;
  text-transform: uppercase;
}

.text-kop h2 {
  margin: 0;
  font-size: 16pt;
  font-weight: bold;
}

.text-kop p {
  margin: 0;
  font-size: 11pt;
}

.garis-kop {
  border: 0;
  border-top: 3px double black;
  margin-bottom: 20px;
}

/* Isi Surat */
.no-border {
  border-collapse: collapse;
}

.no-border td {
  padding: 2px 0;
  vertical-align: top;
}

.header-surat {
  margin-bottom: 15px;
}

.tujuan {
  margin-bottom: 15px;
}

.pembuka {
  margin-bottom: 15px;
  text-align: justify;
}

/* Tabel */
.tabel-siswa {
  width: 100%;
  border-collapse: collapse;
  margin-bottom: 20px;
}

.tabel-siswa th {
  background-color: #f4f4f4;
  font-weight: bold;
  text-align: center;
}

.tabel-siswa td,
.tabel-siswa th {
  border: 1px solid black;
  padding: 6px;
}

/* Tanda Tangan */
.penutup {
  margin-bottom: 30px;
  text-align: justify;
}

.tanda-tangan {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  width: 45%;
  margin-left: auto;
}

.tempat-tanggal {
  margin-bottom: 5px;
}

.jabatan {
  margin-bottom: 60px;
}

/* Ruang untuk TTD basah / QR */
.nama-ttd {
  text-align: center;
}


/* ================= GLOBAL PRINT STYLES ================= */
@media print {

  /* Sembunyikan semua elemen kecuali area cetak */
  body * {
    visibility: hidden;
  }

  .print-area,
  .print-area * {
    visibility: visible;
  }

  .print-area {
    position: absolute;
    left: 0;
    top: 0;
    width: 100%;
    margin: 0;
    padding: 0;
    box-shadow: none;
  }

  /* Reset margin browser */
  @page {
    size: A4;
    margin: 0;
  }

  /* Pastikan background warna (jika ada) tetap tercetak */
  * {
    -webkit-print-color-adjust: exact !important;
    print-color-adjust: exact !important;
  }
}
</style>

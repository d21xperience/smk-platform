<template>
  <q-page class="q-pa-md">
    <!-- PANEL HEADER DENGAN FITUR EKSPOR DOKUMEN -->
    <q-card flat bordered class="q-mb-md bg-teal-1 text-teal-9">
      <q-card-section class="row items-center justify-between q-col-gutter-sm">
        <div class="row items-center">
          <q-avatar icon="analytics" color="teal" text-color="white" class="q-mr-sm" size="md" />
          <div>
            <div class="text-subtitle1 text-weight-bold">Konsolidasi Laporan Keuangan Aset</div>
            <div class="text-caption text-grey-7">Tahun Buku Operasional Sensus: 2026</div>
          </div>
        </div>

        <!-- Tombol Ekspor Masal Dokumen untuk Yayasan -->
        <div class="row q-gutter-sm">
          <q-btn color="indigo" icon="description" label="Ekspor Excel (Buku Induk)" @click="eksporData('Excel')" />
          <q-btn color="negative" icon="picture_as_pdf" label="Cetak PDF Laporan" @click="eksporData('PDF')" />
        </div>
      </q-card-section>
    </q-card>

    <!-- PILIHAN SUB-LAPORAN (TABS) -->
    <q-tabs v-model="subTab" dense class="text-grey-7 bg-white q-mb-md rounded-borders shadow-1" active-color="teal"
      indicator-color="teal" align="justify">
      <q-tab name="buku_induk" icon="menu_book" label="Buku Induk Rekapitulasi Inventaris" />
      <q-tab name="depresiasi" icon="trending_down" label="Tabel Penyusutan & Nilai Buku Aset" />
    </q-tabs>

    <!-- PANEL KONTEN TAB 1: BUKU INDUK INVENTARIS -->
    <q-tab-panels v-model="subTab" animated class="bg-transparent">
      <q-tab-panel name="buku_induk" class="q-pa-none">
        <q-table title="Buku Induk Register Aset Sekolah Swasta" :rows="masterLaporanAset" :columns="kolomBukuInduk"
          row-key="id" flat bordered>
          <template v-slot:body-cell-harga="props">
            <q-td :props="props" class="text-right">
              Rp {{ props.value.toLocaleString('id-ID') }}
            </q-td>
          </template>
        </q-table>
      </q-tab-panel>

      <!-- PANEL KONTEN TAB 2: TABEL PENYUSUTAN NILAI (DEPRESIASI) -->
      <q-tab-panel name="depresiasi" class="q-pa-none">
        <q-table title="Kalkulasi Akumulasi Depresiasi (Metode Garis Lurus)" :rows="masterLaporanAset"
          :columns="kolomDepresiasi" row-key="id" flat bordered>
          <!-- Perhitungan Dinamis Frontend untuk Beban Penyusutan -->
          <template v-slot:body-cell-penyusutanTahunan="props">
            <q-td :props="props" class="text-right text-orange-9 text-weight-medium">
              Rp {{ hitungPenyusutanTahunan(props.row).toLocaleString('id-ID') }}
            </q-td>
          </template>

          <!-- Perhitungan Dinamis Frontend untuk Sisa Nilai Buku Aktiva -->
          <template v-slot:body-cell-nilaiBuku="props">
            <q-td :props="props" class="text-right text-green-9 text-weight-bold bg-green-1">
              Rp {{ hitungNilaiBukuSaatIni(props.row).toLocaleString('id-ID') }}
            </q-td>
          </template>
        </q-table>
      </q-tab-panel>
    </q-tab-panels>
  </q-page>
</template>

<script setup>
import { ref } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const subTab = ref('buku_induk')

// Asumsi Parameter Waktu Berjalan Saat Ini untuk Akuntansi (Tahun Sensus: 2026)
const TAHUN_SEKARANG = 2026

// 1. Master Data Finansial Aset Tetap Sekolah Swasta
const masterLaporanAset = ref([
  { id: 'YYS/SAR-COMP/2023/001', name: 'PC Komputer i5 Lab Server', tglBeli: 2023, harga: 12000000, umurEkonomis: 4, nilaiSisaAset: 2000000, asalDana: 'Dana Yayasan' },
  { id: 'YYS/SAR-PROJ/2024/015', name: 'Proyektor BenQ Wall Mount', tglBeli: 2024, harga: 7500000, umurEkonomis: 5, nilaiSisaAset: 1000000, asalDana: 'Dana BOS' },
  { id: 'YYS/SAR-ACUN/2022/008', name: 'AC Split Panasonic 1.5 PK', tglBeli: 2022, harga: 5500000, umurEkonomis: 5, nilaiSisaAset: 500000, asalDana: 'Dana Yayasan' },
  { id: 'YYS/SAR-MEBL/2025/012', name: 'Meja Rapat Jati Kursi Set', tglBeli: 2025, harga: 25000000, umurEkonomis: 10, nilaiSisaAset: 5000000, asalDana: 'Dana Yayasan' }
])

// 2. Desain Struktur Kolom Tab 1: Buku Induk
const kolomBukuInduk = [
  { name: 'id', label: 'Kode Register Barang', align: 'left', field: 'id', sortable: true },
  { name: 'name', label: 'Nama Spesifikasi Barang', align: 'left', field: 'name', sortable: true },
  { name: 'tglBeli', label: 'Tahun Perolehan', align: 'center', field: 'tglBeli', sortable: true },
  { name: 'harga', label: 'Harga Beli Awal (Mula-Mula)', align: 'right', field: 'harga', sortable: true },
  { name: 'asalDana', label: 'Sumber Pembiayaan', align: 'center', field: 'asalDana' }
]

// 3. Desain Struktur Kolom Tab 2: Penyusutan Finansial
const kolomDepresiasi = [
  { name: 'id', label: 'Kode Barang', align: 'left', field: 'id' },
  { name: 'name', label: 'Nama Aset', align: 'left', field: 'name' },
  { name: 'harga', label: 'Harga Perolehan', align: 'right', field: 'harga', format: val => `Rp ${val.toLocaleString('id-ID')}` },
  { name: 'umurEkonomis', label: 'Masa Umur (Thn)', align: 'center', field: 'umurEkonomis' },
  { name: 'penyusutanTahunan', label: 'Penyusutan / Tahun', align: 'right' },
  { name: 'nilaiBuku', label: 'Nilai Buku (Sisa Aktiva 2026)', align: 'right' }
]

// 4. ALGORITMA AKUNTANSI SWASTA: METODE GARIS LURUS
// Formula: (Harga Perolehan - Nilai Sisa Akhir) / Umur Ekonomis
function hitungPenyusutanTahunan(row) {
  return (row.harga - row.nilaiSisaAset) / row.umurEkonomis
}

function hitungNilaiBukuSaatIni(row) {
  const penyusutanPerTahun = hitungPenyusutanTahunan(row)
  const jumlahTahunBerjalan = TAHUN_SEKARANG - row.tglBeli

  // Total Akumulasi Depresiasi selama aset dimiliki
  const totalAkumulasiPenyusutan = penyusutanPerTahun * jumlahTahunBerjalan

  // Nilai Buku = Harga Awal - Total Akumulasi Penyusutan
  const sisaNilaiBuku = row.harga - totalAkumulasiPenyusutan

  // Amankan agar nilai buku tidak jatuh di bawah nilai sisa minimal barang rongsokan
  return sisaNilaiBuku < row.nilaiSisaAset ? row.nilaiSisaAset : sisaNilaiBuku
}

// 5. Fungsi Aksi Ekspor Tombol Dokumen Luar
function eksporData(tipeFormat) {
  $q.notify({
    color: 'teal',
    icon: 'cloud_download',
    message: `Menyusun berkas neraca saldo, sukses mengeksport seluruh tabel ke dalam format ${tipeFormat} Yayasan.`
  })
}
</script>

<style scoped>
.text-weight-bold {
  font-family: inherit;
}
</style>

<template>
  <q-page class="q-pa-md class-print-area">
    <!-- AREA NON-PRINT: Sembunyikan kontrol ini saat proses cetak browser aktif -->
    <div class="no-print">
      <!-- ALERT PANDUAN -->
      <q-badge color="info" class="q-pa-sm q-mb-md full-width text-subtitle2">
        <q-icon name="info" class="q-mr-xs" />
        Pilih aset dengan mencentang kotak di kiri tabel, lalu klik tombol "Cetak Stiker Terpilih" untuk mencetak stiker
        label.
      </q-badge>

      <!-- TABEL SELEKSI MASSAL -->
      <q-table title="Generator Stiker QR Code" :rows="daftarAset" :columns="columns" row-key="id" selection="multiple"
        v-model:selected="asetTerpilih" :filter="filterText" flat bordered class="q-mb-lg">
        <!-- Slot Atas Kanan untuk Pencarian & Aksi Massal -->
        <template v-slot:top-right>
          <div class="row q-gutter-sm items-center">
            <q-input borderless dense debounce="300" v-model="filterText" placeholder="Cari Kode/Nama Aset...">
              <template v-slot:append>
                <q-icon name="search" />
              </template>
            </q-input>

            <q-btn color="primary" icon="qr_code_scanner" :label="`Cetak Stiker Terpilih (${asetTerpilih.length})`"
              :disabled="asetTerpilih.length === 0" @click="eksekusiCetak" />
          </div>
        </template>
      </q-table>
    </div>

    <!-- AREA PRATINJAU CETAKAN (PRINT PREVIEW GRID) -->
    <!-- Bagian ini akan otomatis ditata rapi oleh CSS saat masuk ke mode cetak kertas -->
    <div v-if="asetTerpilih.length > 0" class="print-preview-section">
      <div class="text-h6 text-grey-8 q-mb-md no-print row items-center">
        <q-icon name="visibility" class="q-mr-sm" /> Pratinjau Tata Letak Stiker Label (Kertas)
      </div>

      <div class="stiker-grid">
        <div v-for="item in asetTerpilih" :key="item.id" class="stiker-box q-pa-sm row items-center justify-between">
          <!-- Sisi Kiri Stiker: Informasi Teks Sekolah & Barang -->
          <div class="stiker-info col">
            <div class="text-sekolah text-weight-bolder text-uppercase">SMA SW KARYA YAYASAN</div>
            <div class="text-nama-barang text-weight-bold q-mt-xs">{{ item.name }}</div>
            <div class="text-kode-aset text-monospace q-mt-xs">{{ item.id }}</div>
            <div class="text-lokasi text-caption text-grey-7 q-mt-xs">Lokasi: {{ item.lokasi }}</div>
          </div>

          <!-- Sisi Kanan Stiker: Simulasi Gambar QR Code -->
          <div class="stiker-qr col-auto q-pl-sm flex flex-center">
            <!-- Di masa depan, ganti div ini dengan library QR Code generator seperti 'qrcode.vue' atau API gratis -->
            <div class="mock-qr-border flex flex-center">
              <q-icon name="qr_code" size="64px" color="black" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const filterText = ref('')
const asetTerpilih = ref([])

// 1. Data Mock Master Aset yang Siap Dibuatkan Label QR
const daftarAset = ref([
  { id: 'YYS/SAR-COMP/2026/001', name: 'PC Komputer Laboratorium A', lokasi: 'Lab Komputer' },
  { id: 'YYS/SAR-COMP/2026/002', name: 'PC Komputer Laboratorium B', lokasi: 'Lab Komputer' },
  { id: 'YYS/SAR-PROJ/2025/012', name: 'Proyektor BenQ Wall Mount', lokasi: 'Ruang Kelas 11-A' },
  { id: 'YYS/SAR-ACUN/2024/005', name: 'AC Split Panasonic 1.5 PK', lokasi: 'Ruang Guru' },
  { id: 'YYS/SAR-MEBL/2026/089', name: 'Meja Rapat Kayu Jati Yayasan', lokasi: 'Ruang Rapat Utama' }
])

// 2. Kolom Utama Tabel Generator
const columns = [
  { name: 'id', label: 'Kode Unik Aset', align: 'left', field: 'id', sortable: true },
  { name: 'name', label: 'Nama Inventaris', align: 'left', field: 'name', sortable: true },
  { name: 'lokasi', label: 'Penempatan Ruang', align: 'left', field: 'lokasi' }
]

// 3. Fungsi Pemicu Cetak Kertas Perangkat Keras
function eksekusiCetak() {
  $q.notify({
    message: 'Membuka dialog pencetakan printer sistem...',
    color: 'positive',
    icon: 'print',
    timeout: 1500
  })

  // Menunggu notifikasi selesai, lalu memicu fungsi cetak bawaan browser
  setTimeout(() => {
    window.print()
  }, 500)
}
</script>

<style scoped>
.text-monospace {
  font-family: 'Courier New', Courier, monospace;
}

/* KUSTOMISASI TATA LETAK STIKER LABEL BARCODE */
.stiker-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 15px;
}

.stiker-box {
  border: 2px dashed #333;
  border-radius: 6px;
  background-color: #fff;
  min-height: 120px;
  width: 100%;
  max-width: 360px;
  page-break-inside: avoid;
  /* Mencegah stiker terpotong antar halaman kertas */
}

.text-sekolah {
  font-size: 10px;
  border-bottom: 1px solid #ccc;
  padding-bottom: 2px;
  letter-spacing: 0.5px;
}

.text-nama-barang {
  font-size: 14px;
  line-height: 1.2;
}

.text-kode-aset {
  font-size: 11px;
  color: #1976D2;
}

.mock-qr-border {
  border: 1px solid #444;
  padding: 4px;
  background: #fff;
}

/* CSS PRINT MEDIA MANAGER: Mengontrol apa yang terlihat di kertas printer */
@media print {

  /* Sembunyikan elemen navigasi, sidebar, header Quasar dan tombol cari */
  .no-print,
  .q-header,
  .q-drawer,
  .q-btn,
  footer {
    display: none !important;
  }

  /* Atur ulang margin halaman cetak agar pas dengan ukuran kertas stiker */
  .class-print-area {
    padding: 0 !important;
  }

  .stiker-grid {
    grid-template-columns: repeat(2, 1fr);
    /* Paksa 2 baris stiker per baris horizontal di kertas A4 */
    gap: 10px;
  }

  .stiker-box {
    border: 1px solid #000;
    /* Ubah ke garis tegas saat dicetak */
  }
}
</style>

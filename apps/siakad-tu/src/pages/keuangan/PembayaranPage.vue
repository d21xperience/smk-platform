<template>
  <q-page class="q-pa-md">
    <div class="text-h5 q-mb-md">Pembayaran SPP</div>

    <!-- ========== 1. Pencarian Siswa ========== -->
    <q-card class="q-mb-md">
      <q-card-section>
        <div class="text-subtitle1">Cari Siswa</div>
        <q-input v-model="keyword" label="Nama atau NIS" outlined dense @keyup.enter="cariSiswa">
          <template v-slot:append>
            <q-icon name="search" class="cursor-pointer" @click="cariSiswa" />
          </template>
        </q-input>
      </q-card-section>
    </q-card>

    <!-- Hasil pencarian -->
    <div v-if="hasilPencarian.length > 0" class="q-mb-md">
      <div class="text-subtitle1">Hasil Pencarian:</div>
      <q-list bordered separator>
        <q-item v-for="siswa in hasilPencarian" :key="siswa.id" clickable @click="pilihSiswa(siswa)"
          :active="siswaTerpilih?.id === siswa.id" active-class="bg-primary text-white">
          <q-item-section>
            <q-item-label>{{ siswa.nama }}</q-item-label>
            <q-item-label caption>NIS: {{ siswa.nis }} | Kelas: {{ siswa.kelas }}</q-item-label>
          </q-item-section>
        </q-item>
      </q-list>
    </div>

    <!-- Jika sudah memilih siswa -->
    <div v-if="siswaTerpilih">
      <q-separator class="q-my-md" />

      <!-- Informasi siswa terpilih -->
      <q-banner class="bg-grey-2 q-mb-md">
        <template v-slot:avatar>
          <q-icon name="school" color="primary" />
        </template>
        <div><strong>Siswa:</strong> {{ siswaTerpilih.nama }} (NIS: {{ siswaTerpilih.nis }}) - {{ siswaTerpilih.kelas }}
        </div>
      </q-banner>

      <!-- ========== 2. Input Pembayaran ========== -->
      <q-card class="q-mb-md">
        <q-card-section>
          <div class="text-subtitle1">Input Pembayaran Baru</div>
          <div class="row q-col-gutter-md q-mt-sm">
            <div class="col-12 col-md-4">
              <q-select v-model="formBulan" :options="daftarBulan" label="Bulan" outlined dense />
            </div>
            <div class="col-12 col-md-4">
              <q-input v-model.number="formTahun" type="number" label="Tahun" outlined dense :min="2020" :max="2030" />
            </div>
            <div class="col-12 col-md-4">
              <q-input v-model.number="formJumlah" type="number" label="Jumlah (Rp)" outlined dense prefix="Rp " />
            </div>
          </div>
          <div class="q-mt-md">
            <q-btn label="Simpan Pembayaran" color="primary" @click="simpanPembayaran" :loading="loading" />
          </div>
        </q-card-section>
      </q-card>

      <!-- ========== 3. Riwayat Pembayaran ========== -->
      <q-card>
        <q-card-section>
          <div class="text-subtitle1">Riwayat Pembayaran SPP</div>
          <q-table :rows="riwayat" :columns="riwayatColumns" row-key="id" :loading="loadingRiwayat" flat bordered dense>
            <template v-slot:body-cell-actions="props">
              <q-td auto-width>
                <q-btn flat round dense icon="print" @click="cetakBukti(props.row)" color="info" />
              </q-td>
            </template>
          </q-table>
        </q-card-section>
      </q-card>
    </div>

    <!-- ========== 4. Dialog Cetak Bukti (Print Preview) ========== -->
    <q-dialog v-model="dialogCetak" persistent>
      <q-card style="min-width: 400px">
        <q-card-section class="row items-center">
          <div class="text-h6">Bukti Pembayaran SPP</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>
        <q-card-section id="buktiPembayaran">
          <div style="padding: 16px; font-family: monospace">
            <div style="text-align: center; margin-bottom: 20px">
              <strong>SEKOLAH MENENGAH ATAS</strong><br />
              <small>Jl. Pendidikan No. 123, Kota</small>
            </div>
            <hr />
            <p><strong>Bukti Pembayaran SPP</strong></p>
            <p><strong>Nama Siswa:</strong> {{ siswaTerpilih?.nama }}</p>
            <p><strong>NIS:</strong> {{ siswaTerpilih?.nis }}</p>
            <p><strong>Kelas:</strong> {{ siswaTerpilih?.kelas }}</p>
            <p><strong>Bulan/Tahun:</strong> {{ buktiData?.bulan }} {{ buktiData?.tahun }}</p>
            <p><strong>Jumlah:</strong> Rp {{ buktiData?.jumlah?.toLocaleString() }}</p>
            <p><strong>Tanggal Bayar:</strong> {{ buktiData?.tanggalBayar }}</p>
            <p><strong>Petugas:</strong> {{ buktiData?.petugas || 'Admin' }}</p>
            <hr />
            <div style="text-align: center; margin-top: 20px">
              <small>Terima kasih atas pembayaran Anda.</small>
            </div>
          </div>
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Tutup" v-close-popup />
          <q-btn flat label="Cetak" color="primary" @click="printBukti" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref } from 'vue';
import { useQuasar } from 'quasar';
import { siswaService } from '../../services/siswaService';
import { usePembayaranSPP } from '../../composables/keuangan/usePembayaranSPP';

const $q = useQuasar();

// State untuk pencarian
const keyword = ref('');
const hasilPencarian = ref([]);
const siswaTerpilih = ref(null);

// State untuk form pembayaran
const formBulan = ref('Januari');
const formTahun = ref(2025);
const formJumlah = ref(200000);
const daftarBulan = ['Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni', 'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'];

// Composable riwayat
const { riwayat, loading, fetchRiwayatBySiswa, simpanPembayaran } = usePembayaranSPP();
const loadingRiwayat = ref(false);

// Table columns riwayat
const riwayatColumns = [
  { name: 'bulan', label: 'Bulan', field: 'bulan', align: 'left' },
  { name: 'tahun', label: 'Tahun', field: 'tahun', align: 'center' },
  { name: 'jumlah', label: 'Jumlah', field: 'jumlah', align: 'right', format: val => `Rp ${val.toLocaleString()}` },
  { name: 'tanggalBayar', label: 'Tgl Bayar', field: 'tanggalBayar', align: 'center' },
  { name: 'actions', label: 'Cetak', align: 'center', field: 'actions' }
];

// Dialog cetak
const dialogCetak = ref(false);
const buktiData = ref(null);

// Fungsi pencarian
async function cariSiswa() {
  if (!keyword.value.trim()) {
    hasilPencarian.value = [];
    return;
  }
  try {
    const result = await siswaService.search(keyword.value);
    hasilPencarian.value = result;
  } catch (err) {
    $q.notify({ type: 'negative', message: `Gagal mencari siswa ${err}` });
  }
}

// Pilih siswa
async function pilihSiswa(siswa) {
  siswaTerpilih.value = siswa;
  // Load riwayat pembayaran siswa tersebut
  loadingRiwayat.value = true;
  try {
    await fetchRiwayatBySiswa(siswa.id);
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Gagal memuat riwayat' + err });
  } finally {
    loadingRiwayat.value = false;
  }
}

// Simpan pembayaran baru
// async function simpanPembayaran() {
//   if (!siswaTerpilih.value) {
//     $q.notify({ type: 'warning', message: 'Pilih siswa terlebih dahulu' });
//     return;
//   }
//   if (!formBulan.value || !formTahun.value || !formJumlah.value || formJumlah.value <= 0) {
//     $q.notify({ type: 'warning', message: 'Isi semua data dengan benar' });
//     return;
//   }

//   try {
//     const dataBaru = {
//       siswaId: siswaTerpilih.value.id,
//       bulan: formBulan.value,
//       tahun: formTahun.value,
//       jumlah: formJumlah.value,
//       petugas: 'Admin' // nanti bisa dari user login
//     };
//     await simpanPembayaran(dataBaru);
//     $q.notify({ type: 'positive', message: 'Pembayaran berhasil disimpan' });
//     // Refresh riwayat
//     await fetchRiwayatBySiswa(siswaTerpilih.value.id);
//     // Reset form
//     formBulan.value = 'Januari';
//     formTahun.value = 2025;
//     formJumlah.value = 200000;
//   } catch (err) {
//     $q.notify({ type: 'negative', message: 'Gagal menyimpan pembayaran' });
//   }
// }

// Cetak bukti
function cetakBukti(row) {
  buktiData.value = {
    ...row,
    namaSiswa: siswaTerpilih.value.nama,
    nis: siswaTerpilih.value.nis,
    kelas: siswaTerpilih.value.kelas
  };
  dialogCetak.value = true;
}

function printBukti() {
  const printContent = document.getElementById('buktiPembayaran').innerHTML;
  // const originalTitle = document.title;
  const printWindow = window.open('', '_blank');
  printWindow.document.write(`
    <html>
      <head><title>Cetak Bukti Pembayaran</title></head>
      <body>${printContent}</body>
    </html>
  `);
  printWindow.document.close();
  printWindow.print();
  // Tidak menutup otomatis, biarkan user menutup sendiri
}
</script>

<style scoped>
/* optional */
</style>

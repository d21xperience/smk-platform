<template>
  <q-page class="q-pa-md">
    <div class="text-h5 q-mb-md">Penerimaan Lainnya</div>

    <!-- Tab: Input Penerimaan & Manajemen Kategori -->
    <q-tabs v-model="tab" dense class="q-mb-md">
      <q-tab name="input" label="Input Penerimaan" />
      <q-tab name="kategori" label="Kelola Kategori" />
    </q-tabs>

    <!-- TAB INPUT PENERIMAAN -->
    <q-tab-panels v-model="tab" animated>
      <q-tab-panel name="input" class="q-pa-none">
        <!-- Pencarian Siswa (sama seperti SPP) -->
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

        <div v-if="siswaTerpilih">
          <q-separator class="q-my-md" />
          <q-banner class="bg-grey-2 q-mb-md">
            <template v-slot:avatar>
              <q-icon name="school" color="primary" />
            </template>
            <div><strong>Siswa:</strong> {{ siswaTerpilih.nama }} (NIS: {{ siswaTerpilih.nis }}) - {{
              siswaTerpilih.kelas }}
            </div>
          </q-banner>

          <!-- Form Input Penerimaan -->
          <q-card class="q-mb-md">
            <q-card-section>
              <div class="text-subtitle1">Tambah Penerimaan Lainnya</div>
              <div class="row q-col-gutter-md q-mt-sm">
                <div class="col-12 col-md-4">
                  <q-select v-model="formKategoriId" :options="kategoriOptions" label="Kategori" option-value="id"
                    option-label="nama" emit-value map-options outlined dense />
                </div>
                <div class="col-12 col-md-4">
                  <q-input v-model.number="formJumlah" type="number" label="Jumlah (Rp)" outlined dense prefix="Rp " />
                </div>
                <div class="col-12 col-md-4">
                  <q-input v-model="formKeterangan" label="Keterangan (opsional)" outlined dense />
                </div>
              </div>
              <div class="q-mt-md">
                <q-btn label="Simpan Penerimaan" color="primary" @click="simpanPembayaran" :loading="loading" />
              </div>
            </q-card-section>
          </q-card>

          <!-- Riwayat Penerimaan -->
          <q-card>
            <q-card-section>
              <div class="text-subtitle1">Riwayat Penerimaan Lainnya</div>
              <q-table :rows="riwayat" :columns="riwayatColumns" row-key="id" :loading="loadingRiwayat" flat bordered
                dense>
                <template v-slot:body-cell-actions="props">
                  <q-td auto-width>
                    <q-btn flat round dense icon="print" @click="cetakBukti(props.row)" color="info" />
                    <q-btn flat round dense icon="delete" @click="hapusPembayaranHandler(props.row.id)"
                      color="negative" />
                  </q-td>
                </template>
              </q-table>
            </q-card-section>
          </q-card>
        </div>
      </q-tab-panel>

      <!-- TAB KELOLA KATEGORI -->
      <q-tab-panel name="kategori" class="q-pa-none">
        <q-card>
          <q-card-section>
            <div class="row justify-between items-center">
              <div class="text-subtitle1">Daftar Kategori Penerimaan</div>
              <q-btn label="Tambah Kategori" color="primary" @click="openKategoriDialog" />
            </div>
            <q-table :rows="kategoriList" :columns="kategoriColumns" row-key="id" :loading="kategoriLoading" flat
              bordered dense class="q-mt-md">
              <template v-slot:body-cell-actions="props">
                <q-td auto-width>
                  <q-btn flat round dense icon="edit" @click="editKategori(props.row)" />
                  <q-btn flat round dense icon="delete" @click="hapusKategoriHandler(props.row.id)" color="negative" />
                </q-td>
              </template>
            </q-table>
          </q-card-section>
        </q-card>
      </q-tab-panel>
    </q-tab-panels>

    <!-- Dialog Kategori -->
    <q-dialog v-model="kategoriDialog">
      <q-card style="min-width: 400px">
        <q-card-section>
          <div class="text-h6">{{ isEditKategori ? 'Edit' : 'Tambah' }} Kategori</div>
        </q-card-section>
        <q-card-section>
          <q-input v-model="kategoriForm.nama" label="Nama Kategori" outlined dense />
          <q-input v-model="kategoriForm.deskripsi" label="Deskripsi" outlined dense type="textarea" />
          <q-input v-model.number="kategoriForm.nominalDefault" type="number" label="Nominal Default (Rp)" outlined
            dense prefix="Rp " />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn flat label="Simpan" color="primary" @click="simpanKategori" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Dialog Cetak Bukti -->
    <q-dialog v-model="dialogCetak" persistent>
      <q-card style="min-width: 400px">
        <q-card-section class="row items-center">
          <div class="text-h6">Bukti Penerimaan Lainnya</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>
        <q-card-section id="buktiPenerimaanLain">
          <div style="padding: 16px; font-family: monospace">
            <div style="text-align: center; margin-bottom: 20px">
              <strong>SEKOLAH MENENGAH ATAS</strong><br />
              <small>Jl. Pendidikan No. 123, Kota</small>
            </div>
            <hr />
            <p><strong>Bukti Penerimaan Lainnya</strong></p>
            <p><strong>Nama Siswa:</strong> {{ siswaTerpilih?.nama }}</p>
            <p><strong>NIS:</strong> {{ siswaTerpilih?.nis }}</p>
            <p><strong>Kategori:</strong> {{ buktiData?.kategoriNama }}</p>
            <p><strong>Jumlah:</strong> Rp {{ buktiData?.jumlah?.toLocaleString() }}</p>
            <p><strong>Tanggal Bayar:</strong> {{ buktiData?.tanggalBayar }}</p>
            <p><strong>Keterangan:</strong> {{ buktiData?.keterangan || '-' }}</p>
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
import { ref, onMounted } from 'vue';
import { useQuasar } from 'quasar';
import { siswaService } from '../../services/siswaService';
import { useKategoriPenerimaan } from '../../composables/keuangan/useKategoriPenerimaan';
import { usePembayaranLain } from '../../composables/keuangan/usePembayaranLain';

const $q = useQuasar();
const tab = ref('input');

// Pencarian siswa
const keyword = ref('');
const hasilPencarian = ref([]);
const siswaTerpilih = ref(null);

// Kategori
const { kategoriList, loading: kategoriLoading, fetchAll: fetchKategori, tambahKategori, updateKategori, hapusKategori } = useKategoriPenerimaan();
const kategoriDialog = ref(false);
const isEditKategori = ref(false);
const kategoriForm = ref({ nama: '', deskripsi: '', nominalDefault: 0 });
const selectedKategoriId = ref(null);

// Pembayaran Lain
const { riwayat, loading, fetchRiwayatBySiswa, simpanPembayaran, hapusPembayaran } = usePembayaranLain();
const loadingRiwayat = ref(false);
const formKategoriId = ref(null);
const formJumlah = ref(0);
const formKeterangan = ref('');

// Columns riwayat
const riwayatColumns = [
  { name: 'kategori', label: 'Kategori', field: row => getKategoriName(row.kategoriId), align: 'left' },
  { name: 'jumlah', label: 'Jumlah', field: 'jumlah', align: 'right', format: val => `Rp ${val.toLocaleString()}` },
  { name: 'tanggalBayar', label: 'Tgl Bayar', field: 'tanggalBayar', align: 'center' },
  { name: 'keterangan', label: 'Keterangan', field: 'keterangan', align: 'left' },
  { name: 'actions', label: 'Aksi', align: 'center' }
];

// Columns kategori
const kategoriColumns = [
  { name: 'nama', label: 'Nama Kategori', field: 'nama' },
  { name: 'deskripsi', label: 'Deskripsi', field: 'deskripsi' },
  { name: 'nominalDefault', label: 'Nominal Default', field: 'nominalDefault', format: val => `Rp ${val?.toLocaleString()}` },
  { name: 'actions', label: 'Aksi', align: 'center' }
];

// Helper
const kategoriOptions = ref([]);
function getKategoriName(id) {
  const kat = kategoriList.value.find(k => k.id === id);
  return kat ? kat.nama : '';
}

async function cariSiswa() {
  if (!keyword.value.trim()) {
    hasilPencarian.value = [];
    return;
  }
  try {
    hasilPencarian.value = await siswaService.search(keyword.value);
  } catch (err) {
    $q.notify({ type: 'negative', message: `Gagal mencari siswa ${err}` });
  }
}

async function pilihSiswa(siswa) {
  siswaTerpilih.value = siswa;
  loadingRiwayat.value = true;
  try {
    await fetchRiwayatBySiswa(siswa.id);
  } catch (err) {
    $q.notify({ type: 'negative', message: `Gagal memuat riwayat ${err}` });
  } finally {
    loadingRiwayat.value = false;
  }
}

// async function simpanPembayaranHandler() {
//   if (!siswaTerpilih.value) {
//     $q.notify({ type: 'warning', message: 'Pilih siswa terlebih dahulu' });
//     return;
//   }
//   if (!formKategoriId.value || !formJumlah.value || formJumlah.value <= 0) {
//     $q.notify({ type: 'warning', message: 'Pilih kategori dan isi jumlah' });
//     return;
//   }
//   try {
//     await simpanPembayaran({
//       siswaId: siswaTerpilih.value.id,
//       kategoriId: formKategoriId.value,
//       jumlah: formJumlah.value,
//       keterangan: formKeterangan.value,
//       petugas: 'Admin'
//     });
//     $q.notify({ type: 'positive', message: 'Penerimaan berhasil disimpan' });
//     await fetchRiwayatBySiswa(siswaTerpilih.value.id);
//     formKategoriId.value = null;
//     formJumlah.value = 0;
//     formKeterangan.value = '';
//   } catch (err) {
//     $q.notify({ type: 'negative', message: `Gagal menyimpan ${err}` });
//   }
// }

async function hapusPembayaranHandler(id) {
  $q.dialog({
    title: 'Konfirmasi',
    message: 'Hapus penerimaan ini?',
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await hapusPembayaran(id);
      $q.notify({ type: 'positive', message: 'Data dihapus' });
    } catch (err) {
      $q.notify({ type: 'negative', message: `Gagal hapus ${err}` });
    }
  });
}

// Manajemen Kategori
function openKategoriDialog() {
  isEditKategori.value = false;
  kategoriForm.value = { nama: '', deskripsi: '', nominalDefault: 0 };
  kategoriDialog.value = true;
}

function editKategori(kategori) {
  isEditKategori.value = true;
  selectedKategoriId.value = kategori.id;
  kategoriForm.value = { ...kategori };
  kategoriDialog.value = true;
}

async function simpanKategori() {
  if (!kategoriForm.value.nama) {
    $q.notify({ type: 'warning', message: 'Nama kategori wajib diisi' });
    return;
  }
  try {
    if (isEditKategori.value) {
      await updateKategori(selectedKategoriId.value, kategoriForm.value);
      $q.notify({ type: 'positive', message: 'Kategori diperbarui' });
    } else {
      await tambahKategori(kategoriForm.value);
      $q.notify({ type: 'positive', message: 'Kategori ditambahkan' });
    }
    kategoriDialog.value = false;
    await fetchKategori(); // refresh
    updateKategoriOptions();
  } catch (err) {
    $q.notify({ type: 'negative', message: err.message });
  }
}

async function hapusKategoriHandler(id) {
  $q.dialog({
    title: 'Konfirmasi',
    message: 'Hapus kategori? Data penerimaan dengan kategori ini tidak akan terhapus namun akan kehilangan referensi.',
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await hapusKategori(id);
      $q.notify({ type: 'positive', message: 'Kategori dihapus' });
      await fetchKategori();
      updateKategoriOptions();
    } catch (err) {
      $q.notify({ type: 'negative', message: err.message });
    }
  });
}

function updateKategoriOptions() {
  kategoriOptions.value = kategoriList.value.map(k => ({ id: k.id, nama: k.nama, nominalDefault: k.nominalDefault }));
}

// Cetak Bukti
const dialogCetak = ref(false);
const buktiData = ref(null);
function cetakBukti(row) {
  const kategori = getKategoriName(row.kategoriId);
  buktiData.value = {
    ...row,
    kategoriNama: kategori,
    namaSiswa: siswaTerpilih.value.nama,
    nis: siswaTerpilih.value.nis,
    kelas: siswaTerpilih.value.kelas
  };
  dialogCetak.value = true;
}
function printBukti() {
  const printContent = document.getElementById('buktiPenerimaanLain').innerHTML;
  const printWindow = window.open('', '_blank');
  printWindow.document.write(`<html><head><title>Cetak Bukti Penerimaan</title></head><body>${printContent}</body></html>`);
  printWindow.document.close();
  printWindow.print();
}

onMounted(async () => {
  await fetchKategori();
  updateKategoriOptions();
});
</script>

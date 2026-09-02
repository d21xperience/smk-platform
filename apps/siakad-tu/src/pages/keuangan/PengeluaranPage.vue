<template>
  <q-page class="q-pa-md">
    <div class="text-h5 q-mb-md">Pengeluaran Keuangan</div>

    <!-- Tab: Data Pengeluaran & Kelola Kategori -->
    <q-tabs v-model="tab" dense>
      <q-tab name="pengeluaran" label="Data Pengeluaran" />
      <q-tab name="kategori" label="Kelola Kategori" />
    </q-tabs>

    <q-separator />

    <q-tab-panels v-model="tab" animated>
      <!-- Panel Pengeluaran -->
      <q-tab-panel name="pengeluaran">
        <!-- Form Input -->
        <q-card class="q-mb-md">
          <q-card-section>
            <div class="text-subtitle1">Tambah / Edit Pengeluaran</div>
            <div class="row q-col-gutter-md q-mt-sm">
              <div class="col-12 col-md-3">
                <q-input v-model="form.tanggal" label="Tanggal" type="date" outlined dense />
              </div>
              <div class="col-12 col-md-3">
                <q-select v-model="form.kategoriId" :options="kategoriOptions" label="Kategori" option-value="id"
                  option-label="nama" emit-value map-options outlined dense />
              </div>
              <div class="col-12 col-md-3">
                <q-input v-model.number="form.nominal" label="Nominal (Rp)" type="number" outlined dense prefix="Rp " />
              </div>
              <div class="col-12 col-md-3">
                <q-uploader label="Upload Bukti" url="#" :auto-upload="false" @added="onFileAdded"
                  style="max-width: 100%" />
              </div>
            </div>
            <div class="row q-mt-sm">
              <div class="col-12">
                <q-input v-model="form.keterangan" label="Keterangan" outlined dense autogrow />
              </div>
            </div>
            <div class="q-mt-md">
              <q-btn label="Simpan" color="primary" @click="simpanPengeluaran" :loading="loadingSimpan" />
              <q-btn v-if="form.id" label="Batal" flat @click="resetForm" class="q-ml-sm" />
            </div>
          </q-card-section>
        </q-card>

        <!-- Filter dan Tabel -->
        <q-card>
          <q-card-section>
            <div class="row q-col-gutter-md">
              <div class="col-12 col-md-4">
                <q-input v-model="filterTanggalMulai" label="Tanggal Mulai" type="date" outlined dense />
              </div>
              <div class="col-12 col-md-4">
                <q-input v-model="filterTanggalAkhir" label="Tanggal Akhir" type="date" outlined dense />
              </div>
              <div class="col-12 col-md-4">
                <q-select v-model="filterKategoriId" :options="kategoriOptions" label="Filter Kategori" clearable
                  outlined dense />
              </div>
            </div>
            <div class="q-mt-md">
              <q-btn label="Cetak Laporan" icon="print" color="info" @click="cetakLaporan" />
            </div>
          </q-card-section>

          <q-table :rows="filteredPengeluaran" :columns="columns" row-key="id" :loading="loading" flat bordered>
            <template v-slot:body-cell-bukti="props">
              <q-td auto-width>
                <div v-if="props.row.bukti">
                  <q-icon name="attach_file" />
                  <span class="text-caption">{{ props.row.bukti }}</span>
                </div>
                <div v-else>-</div>
              </q-td>
            </template>
            <template v-slot:body-cell-actions="props">
              <q-td auto-width>
                <q-btn flat round dense icon="edit" @click="editPengeluaran(props.row)" />
                <q-btn flat round dense icon="delete" color="negative" @click="hapusPengeluaran(props.row.id)" />
              </q-td>
            </template>
          </q-table>
        </q-card>
      </q-tab-panel>

      <!-- Panel Kelola Kategori -->
      <q-tab-panel name="kategori">
        <div class="q-mb-md">
          <q-btn label="Tambah Kategori" color="primary" @click="dialogKategori = true" />
        </div>
        <q-table :rows="kategoriStore.items" :columns="kategoriColumns" row-key="id" flat bordered>
          <template v-slot:body-cell-actions="props">
            <q-td auto-width>
              <q-btn flat round dense icon="edit" @click="editKategori(props.row)" />
              <q-btn flat round dense icon="delete" color="negative" @click="hapusKategori(props.row.id)" />
            </q-td>
          </template>
        </q-table>
      </q-tab-panel>
    </q-tab-panels>

    <!-- Dialog Tambah/Edit Kategori -->
    <q-dialog v-model="dialogKategori">
      <q-card style="min-width: 400px">
        <q-card-section>
          <div class="text-h6">{{ isEditKategori ? 'Edit' : 'Tambah' }} Kategori</div>
        </q-card-section>
        <q-card-section>
          <q-input v-model="formKategori.nama" label="Nama Kategori" outlined dense />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn flat label="Simpan" color="primary" @click="simpanKategori" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useQuasar } from 'quasar';
import { useKategoriPengeluaranStore } from 'stores/keuangan/kategoriPengeluaranStore';
import { usePengeluaranStore } from 'stores//keuangan/pengeluaranStore';

const $q = useQuasar();
const kategoriStore = useKategoriPengeluaranStore();
const pengeluaranStore = usePengeluaranStore();

const tab = ref('pengeluaran');
const loading = computed(() => pengeluaranStore.loading);
const loadingSimpan = ref(false);

// Form pengeluaran
const form = ref({ id: null, tanggal: '', kategoriId: null, nominal: null, keterangan: '', bukti: null });
const resetForm = () => { form.value = { id: null, tanggal: '', kategoriId: null, nominal: null, keterangan: '', bukti: null }; };

// Upload bukti - simpan nama file (simulasi)
function onFileAdded(files) {
  if (files.length) form.value.bukti = files[0].name;
}

// Filter
const filterTanggalMulai = ref('');
const filterTanggalAkhir = ref('');
const filterKategoriId = ref(null);

const kategoriOptions = computed(() => kategoriStore.items);

const filteredPengeluaran = computed(() => {
  let data = pengeluaranStore.items;
  if (filterTanggalMulai.value) data = data.filter(p => p.tanggal >= filterTanggalMulai.value);
  if (filterTanggalAkhir.value) data = data.filter(p => p.tanggal <= filterTanggalAkhir.value);
  if (filterKategoriId.value) data = data.filter(p => p.kategoriId === filterKategoriId.value);
  return data;
});

const columns = [
  { name: 'tanggal', label: 'Tanggal', field: 'tanggal', align: 'left', sortable: true },
  { name: 'kategori', label: 'Kategori', field: row => kategoriStore.items.find(k => k.id === row.kategoriId)?.nama || '-' },
  { name: 'nominal', label: 'Nominal', field: 'nominal', align: 'right', sortable: true, format: val => `Rp ${val?.toLocaleString()}` },
  { name: 'keterangan', label: 'Keterangan', field: 'keterangan' },
  { name: 'bukti', label: 'Bukti', field: 'bukti', align: 'center' },
  { name: 'actions', label: 'Aksi', align: 'center' }
];

// CRUD Pengeluaran
async function simpanPengeluaran() {
  if (!form.value.tanggal || !form.value.kategoriId || !form.value.nominal) {
    $q.notify({ type: 'warning', message: 'Lengkapi data wajib' });
    return;
  }
  loadingSimpan.value = true;
  try {
    await pengeluaranStore.simpan({ ...form.value });
    $q.notify({ type: 'positive', message: 'Data tersimpan' });
    resetForm();
  } catch (err) {
    $q.notify({ type: 'negative', message: `Gagal simpan ${err}` });
  } finally {
    loadingSimpan.value = false;
  }
}

function editPengeluaran(row) {
  form.value = { ...row };
}

async function hapusPengeluaran(id) {
  $q.dialog({ title: 'Konfirmasi', message: 'Hapus pengeluaran?', cancel: true }).onOk(async () => {
    await pengeluaranStore.hapus(id);
    $q.notify({ type: 'positive', message: 'Terhapus' });
  });
}

// Cetak laporan (print table)
function cetakLaporan() {
  const printContent = document.querySelector('.q-table').cloneNode(true);
  const win = window.open('', '_blank');
  win.document.write(`
    <html><head><title>Laporan Pengeluaran</title>
    <style>table { border-collapse: collapse; width:100%; } th, td { border:1px solid #ddd; padding:8px; }</style>
    </head><body><h3>Laporan Pengeluaran</h3>${printContent.outerHTML}</body></html>
  `);
  win.document.close();
  win.print();
}

// Kelola Kategori
const dialogKategori = ref(false);
const isEditKategori = ref(false);
const formKategori = ref({ id: null, nama: '' });
const kategoriColumns = [
  { name: 'nama', label: 'Nama Kategori', field: 'nama' },
  { name: 'actions', label: 'Aksi', align: 'center' }
];

function editKategori(kat) {
  isEditKategori.value = true;
  formKategori.value = { ...kat };
  dialogKategori.value = true;
}
async function simpanKategori() {
  if (!formKategori.value.nama) return;
  if (isEditKategori.value) {
    await kategoriStore.update(formKategori.value.id, formKategori.value);
  } else {
    await kategoriStore.tambah({ nama: formKategori.value.nama });
  }
  dialogKategori.value = false;
  resetFormKategori();
  await kategoriStore.fetchAll(); // refresh
}
function resetFormKategori() { formKategori.value = { id: null, nama: '' }; isEditKategori.value = false; }
// eslint-disable-next-line no-unused-vars
async function hapusKategori(id) { /* konfirmasi + panggil service delete */ }

onMounted(async () => {
  await kategoriStore.fetchAll();
  await pengeluaranStore.fetchAll();
});
</script>

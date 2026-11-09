<template>
  <q-page class="q-pa-md">
    <!-- PANEL SELEKTOR RUANGAN -->
    <q-card flat bordered class="q-mb-md bg-blue-grey-1">
      <q-card-section class="row items-center q-col-gutter-md">
        <div class="col-12 col-sm-6 col-md-4">
          <q-select outlined dense bg-color="white" v-model="ruanganTerpilih" :options="kirStore.rooms" option-label="nama"
            label="Pilih Lokasi Ruangan / Kelas" @update:model-value="onRuanganChange" :loading="kirStore.isLoading">
            <template v-slot:prepend><q-icon name="room" color="primary" /></template>
          </q-select>
        </div>
        <q-space />
        <div class="col-12 col-sm-auto">
          <q-btn color="secondary" icon="print" label="Cetak Lembar KIR" :disabled="!ruanganTerpilih" @click="cetakKIR" :loading="isDownloading" />
        </div>
      </q-card-section>
    </q-card>

    <!-- KONTEN DATA KIR -->
    <div v-if="ruanganTerpilih">
      <!-- KARTU REKAPITULASI -->
      <div class="row q-col-gutter-md q-mb-md">
        <div class="col-12 col-md-3">
          <q-card class="bg-primary text-white">
            <q-card-section>
              <div class="text-caption">Total Aset</div>
              <div class="text-h5 text-weight-bold">{{ kirStore.items.length }} Unit</div>
            </q-card-section>
          </q-card>
        </div>
        <div class="col-12 col-md-3">
          <q-card class="bg-green-7 text-white">
            <q-card-section>
              <div class="text-caption">Kondisi Baik</div>
              <div class="text-h5 text-weight-bold">{{ countBaik }} Unit</div>
            </q-card-section>
          </q-card>
        </div>
        <div class="col-12 col-md-3">
          <q-card class="bg-orange-7 text-white">
            <q-card-section>
              <div class="text-caption">Rusak Ringan</div>
              <div class="text-h5 text-weight-bold">{{ countRR }} Unit</div>
            </q-card-section>
          </q-card>
        </div>
        <div class="col-12 col-md-3">
          <q-card class="bg-red-7 text-white">
            <q-card-section>
              <div class="text-caption">Rusak Berat</div>
              <div class="text-h5 text-weight-bold">{{ countRB }} Unit</div>
            </q-card-section>
          </q-card>
        </div>
      </div>

      <q-card flat bordered class="q-mb-md border-primary shadow-1">
        <q-card-section class="row items-center justify-between q-py-md">
          <div class="row items-center">
            <q-avatar icon="person" color="primary" text-color="white" size="md" class="q-mr-sm" />
            <div>
              <div class="text-caption text-grey-7">Penanggung Jawab Ruangan (PJR)</div>
              <div class="text-subtitle1 text-weight-bold text-grey-9">{{ ruanganTerpilih.penanggungJawab }}</div>
            </div>
          </div>
          <q-separator vertical class="gt-xs q-mx-md" />
          <div>
            <div class="text-caption text-grey-7 text-right-gt-xs">Total Nilai Buku</div>
            <div class="text-subtitle1 text-weight-bold text-primary text-right-gt-xs">Rp {{ totalNilai.toLocaleString('id-ID') }}</div>
          </div>
        </q-card-section>
      </q-card>

      <q-table :title="`Daftar Inventaris: ${ruanganTerpilih.nama}`" :rows="kirStore.items" :columns="columns"
        row-key="id" flat bordered :loading="kirStore.isLoading" :rows-per-page-options="[0]">
        
        <template v-slot:top-right>
          <q-btn color="primary" icon="add" label="Tambah Barang" @click="bukaDialogTambah" no-caps />
        </template>

        <template v-slot:body-cell-condition="props">
          <q-td :props="props">
            <q-chip :color="props.value === 'Bagus' ? 'green-2' : props.value === 'Rusak Ringan' ? 'orange-2' : 'red-2'"
              :text-color="props.value === 'Bagus' ? 'green-9' : props.value === 'Rusak Ringan' ? 'orange-9' : 'red-9'"
              dense class="text-weight-medium">{{ props.value }}</q-chip>
          </q-td>
        </template>

        <template v-slot:body-cell-id="props">
          <q-td :props="props"><span class="text-monospace text-weight-bold text-primary">{{ props.value }}</span></q-td>
        </template>

        <template v-slot:body-cell-actions="props">
          <q-td :props="props" class="q-gutter-x-sm">
            <q-btn flat round dense color="primary" icon="edit" @click="bukaDialogEdit(props.row)">
              <q-tooltip>Edit Data</q-tooltip>
            </q-btn>
            <q-btn flat round dense color="secondary" icon="folder_open" @click="bukaDokumen(props.row)">
              <q-tooltip>Dokumen & Detail</q-tooltip>
            </q-btn>
            <q-btn flat round dense color="negative" icon="delete" @click="konfirmasiHapus(props.row)">
              <q-tooltip>Hapus</q-tooltip>
            </q-btn>
          </q-td>
        </template>
      </q-table>
    </div>

    <!-- DIALOG DETAIL & DOKUMEN ASET -->
    <q-dialog v-model="showDokumenDialog" maximized transition-show="slide-up" transition-hide="slide-down">
      <q-card class="column no-wrap">
        <q-card-section class="row items-center q-pb-none bg-primary text-white">
          <div class="text-h6">Detail Aset & Dokumen</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>
        <q-card-section class="col scroll q-pt-md">
          <AssetDocumentManager v-if="selectedAsset" :assetId="selectedAsset.id" :assetCode="selectedAsset.id" />
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- DIALOG FORM CRUD -->
    <q-dialog v-model="showDialog" persistent>
      <q-card style="min-width: 450px">
        <q-card-section class="text-h6">{{ isEdit ? 'Edit Barang Inventaris' : 'Tambah Barang Inventaris' }}</q-card-section>
        <q-card-section class="q-gutter-y-md">
          <q-input v-if="isEdit" v-model="formData.id" label="Kode Unik Aset" outlined dense disable />
          
          <q-select v-if="!isEdit" v-model="formData.catalog_id" :options="kirStore.catalogs" option-value="id" option-label="name"
            label="Pilih Katalog Barang" outlined dense emit-value map-options>
            <template v-slot:option="scope">
              <q-item v-bind="scope.itemProps">
                <q-item-section>
                  <q-item-label>{{ scope.opt.name }}</q-item-label>
                  <q-item-label caption>Merk: {{ scope.opt.brand || '-' }}</q-item-label>
                </q-item-section>
              </q-item>
            </template>
          </q-select>

          <q-input v-if="!isEdit" v-model="formData.purchase_date" type="date" label="Tanggal Pembelian (Tahun untuk Kode)" outlined dense />
          
          <q-select v-model="formData.condition" :options="['Bagus', 'Rusak Ringan', 'Rusak Berat']" label="Kondisi Fisik" outlined dense />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn color="primary" :label="isEdit ? 'Simpan Perubahan' : 'Simpan'" @click="submitForm" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useKirStore } from '@/stores/sarpras/kir'
import AssetDocumentManager from '@/components/sarpras/AssetDocumentManager.vue'

const $q = useQuasar()
const kirStore = useKirStore()
const ruanganTerpilih = ref(null)
const isDownloading = ref(false)

const showDialog = ref(false)
const showDokumenDialog = ref(false)
const isEdit = ref(false)
const selectedAsset = ref(null)
const formData = ref({ id: '', catalog_id: '', purchase_date: new Date().toISOString().split('T')[0], condition: 'Bagus' })

const columns = [
  { name: 'id', label: 'Kode Unik Aset (Kodifikasi)', align: 'left', field: 'id', sortable: true },
  { name: 'name', label: 'Nama Barang', align: 'left', field: 'name', sortable: true },
  { name: 'brand', label: 'Merk/Brand', align: 'left', field: 'brand' },
  { name: 'category_name', label: 'Kategori', align: 'center', field: 'category_name' },
  { name: 'nilai_buku', label: 'Nilai Buku', align: 'right', field: 'nilai_buku', format: val => `Rp ${val.toLocaleString('id-ID')}` },
  { name: 'condition', label: 'Kondisi Fisik', align: 'center', field: 'condition', sortable: true },
  { name: 'actions', label: 'Aksi', align: 'center', field: 'id' }
]

// Computed untuk Rekapitulasi
const countBaik = computed(() => kirStore.items.filter(i => i.condition === 'Bagus').length)
const countRR = computed(() => kirStore.items.filter(i => i.condition === 'Rusak Ringan').length)
const countRB = computed(() => kirStore.items.filter(i => i.condition === 'Rusak Berat').length)
const totalNilai = computed(() => kirStore.items.reduce((sum, item) => sum + (item.nilai_buku || 0), 0))

onMounted(() => { 
  kirStore.fetchRooms()
  kirStore.fetchMasterData()
})

function onRuanganChange(val) { if (val) kirStore.fetchItemsByRoom(val.id) }

function bukaDokumen(item) {
  selectedAsset.value = item
  showDokumenDialog.value = true
}

function bukaDialogTambah() {
  isEdit.value = false
  formData.value = { id: '', catalog_id: '', purchase_date: new Date().toISOString().split('T')[0], condition: 'Bagus' }
  showDialog.value = true
}

function bukaDialogEdit(item) {
  isEdit.value = true
  formData.value = { ...item }
  showDialog.value = true
}

async function submitForm() {
  try {
    if (isEdit.value) {
      await kirStore.updateItem(formData.value.id, formData.value)
      $q.notify({ type: 'positive', message: 'Barang berhasil diupdate' })
    } else {
      formData.value.room_id = ruanganTerpilih.value.id
      const res = await kirStore.createItem(formData.value)
      $q.notify({ type: 'positive', message: res.data.message || 'Barang berhasil ditambahkan' })
    }
    showDialog.value = false
    kirStore.fetchItemsByRoom(ruanganTerpilih.value.id)
  } catch (error) {
    $q.notify({ type: 'negative', message: error.response?.data?.error || 'Terjadi kesalahan' })
  }
}

function konfirmasiHapus(item) {
  $q.dialog({ title: 'Konfirmasi Hapus', message: `Hapus aset "${item.name}"?`, cancel: true, persistent: true })
    .onOk(async () => {
      await kirStore.deleteItem(item.id)
      $q.notify({ type: 'positive', message: 'Barang berhasil dihapus' })
      kirStore.fetchItemsByRoom(ruanganTerpilih.value.id)
    })
}

async function cetakKIR() {
  isDownloading.value = true
  try {
    await kirStore.downloadKIRReport(ruanganTerpilih.value.id)
    $q.notify({ type: 'positive', message: 'Laporan KIR berhasil diunduh' })
  } catch (error) {
    $q.notify({ type: 'negative', message: 'Gagal mengunduh laporan KIR' })
  } finally {
    isDownloading.value = false
  }
}
</script>

<style scoped>
.border-primary { border-left: 5px solid var(--q-primary); }
.text-monospace { font-family: 'Courier New', Courier, monospace; }
@media gt-xs { .text-right-gt-xs { text-align: right; } }
</style>

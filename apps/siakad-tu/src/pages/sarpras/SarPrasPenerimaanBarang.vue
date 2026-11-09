<template>
  <q-page class="q-pa-md">
    <!-- HEADER & TOMBOL TAMBAH -->
    <div class="row q-mb-md items-center justify-between">
      <div class="text-h6 text-weight-bold">Buku Penerimaan Barang (BAST)</div>
      <q-btn color="primary" icon="add" label="Buat Penerimaan Baru" @click="bukaForm()" no-caps />
    </div>

    <!-- TABEL PENERIMAAN -->
    <q-table :rows="store.receipts" :columns="columns" row-key="id" flat bordered :loading="store.isLoading"
      :filter="filterText">
      <template v-slot:top-right>
        <q-input borderless dense debounce="300" v-model="filterText" placeholder="Cari BAST/Supplier...">
          <template v-slot:append><q-icon name="search" /></template>
        </q-input>
      </template>

      <!-- Status Badge -->
      <template v-slot:body-cell-status="props">
        <q-td :props="props" class="text-center">
          <q-chip :color="getStatusColor(props.value)" text-color="white" dense class="text-weight-bold">
            {{ props.value }}
          </q-chip>
        </q-td>
      </template>

      <!-- Kolom Aksi -->
      <template v-slot:body-cell-aksi="props">
        <q-td :props="props" class="text-center">
          <div v-if="props.row.status === 'DRAFT'" class="row justify-center q-gutter-xs">
            <q-btn dense round size="sm" color="positive" icon="check" @click="verifikasiData(props.row)">
              <q-tooltip>Verifikasi (Integrasi Stok/Aset)</q-tooltip>
            </q-btn>
            <q-btn dense round size="sm" color="negative" icon="close" @click="bukaTolak(props.row)">
              <q-tooltip>Tolak Penerimaan</q-tooltip>
            </q-btn>
          </div>
          <div class="row justify-center q-gutter-xs">
            <q-btn dense round size="sm" color="primary" icon="picture_as_pdf" @click="store.downloadPDF(props.row.id)">
              <q-tooltip>Cetak BAST PDF</q-tooltip>
            </q-btn>
          </div>
        </q-td>
      </template>
    </q-table>

    <!-- DIALOG FORM PENERIMAAN -->
    <q-dialog v-model="showForm" persistent maximized transition-show="slide-up" transition-hide="slide-down">
      <q-card class="column no-wrap">
        <q-card-section class="row items-center q-pb-none bg-primary text-white">
          <div class="text-h6">Form Penerimaan Barang</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="col scroll q-gutter-y-md">
          <!-- INFO HEADER -->
          <div class="row q-col-gutter-md">
            <div class="col-12 col-md-4">
              <q-input outlined dense v-model="form.receipt_date" type="date" label="Tanggal Terima *" />
            </div>
            <div class="col-12 col-md-4">
              <q-select outlined dense v-model="form.source_type"
                :options="['PROCUREMENT', 'HIBAH', 'DONASI', 'MUTASI', 'OTHER']" label="Sumber *" />
            </div>
            <div class="col-12 col-md-4">
              <q-input outlined dense v-model="form.supplier_name" label="Supplier / Donor *" />
            </div>
            <div class="col-12 col-md-6">
              <q-input outlined dense v-model="form.invoice_number" label="No. Faktur / Surat Jalan" />
            </div>
            <div class="col-12 col-md-6">
              <q-input outlined dense v-model="form.funding_source" label="Sumber Dana (BOS, BOP, Hibah)" />
            </div>
          </div>

          <q-separator />

          <!-- DETAIL BARANG (MULTI-ITEM) -->
          <div class="row items-center justify-between">
            <div class="text-subtitle1 text-weight-bold">Detail Barang</div>
            <q-btn color="secondary" icon="add" label="Tambah Baris" @click="tambahBaris" size="sm" no-caps />
          </div>

          <div v-for="(item, index) in form.items" :key="index"
            class="row q-col-gutter-sm items-center bg-grey-1 q-pa-sm rounded-borders">
            <div class="col-12 col-md-3">
              <q-select outlined dense v-model="item.catalog_id" :options="kirStore.catalogs" option-value="id"
                option-label="name" label="Barang *" emit-value map-options
                @update:model-value="onPilihBarang(index, $event)" />
            </div>
            <div class="col-6 col-md-2">
              <q-input outlined dense type="number" v-model.number="item.quantity" label="Qty *" />
            </div>
            <div class="col-6 col-md-1">
              <q-input outlined dense v-model="item.unit" label="Satuan" />
            </div>
            <div class="col-6 col-md-2">
              <q-input outlined dense type="number" v-model.number="item.unit_price" label="Harga Satuan" />
            </div>
            <div class="col-6 col-md-2">
              <q-select outlined dense v-model="item.condition" :options="['BAIK', 'RUSAK_RINGAN', 'RUSAK_BERAT']"
                label="Kondisi" />
            </div>
            <div class="col-12 col-md-2 text-right">
              <q-toggle v-model="item.is_inventory" label="Aset Tetap?" color="primary" dense />
              <q-btn flat round dense color="negative" icon="delete" @click="hapusBaris(index)" class="q-ml-sm" />
            </div>
          </div>

          <q-separator />

          <!-- UPLOAD DOKUMEN -->
          <div class="row items-center q-gutter-md">
            <q-file v-model="fileDokumen" label="Upload Scan BAST / Faktur" outlined dense
              @update:model-value="uploadFile" :loading="isUploading">
              <template v-slot:prepend><q-icon name="cloud_upload" /></template>
            </q-file>
            <div v-if="form.document_urls.length > 0" class="text-positive text-caption">
              <q-icon name="check_circle" /> {{ form.document_urls.length }} Dokumen terupload
            </div>
          </div>
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md bg-grey-2">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn color="primary" label="Simpan Penerimaan (DRAFT)" @click="simpanPenerimaan" :loading="isSubmitting" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- DIALOG TOLAK -->
    <q-dialog v-model="showTolak">
      <q-card style="width: 400px;">
        <q-card-section class="text-h6">Tolak Penerimaan</q-card-section>
        <q-card-section>
          <q-input outlined dense type="textarea" v-model="rejectReason" label="Alasan Penolakan (Wajib) *" />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn color="negative" label="Tolak" @click="submitTolak" :loading="isSubmitting" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useQuasar } from 'quasar'
import { useGoodsReceiptStore } from '@/stores/sarpras/goodsReceipt'
import { useKirStore } from '@/stores/sarpras/kir'
import { useContextStore } from '@/stores/contextStore'

const $q = useQuasar()
const store = useGoodsReceiptStore()
const kirStore = useKirStore()
const contextStore = useContextStore()

const filterText = ref('')
const showForm = ref(false)
const showTolak = ref(false)
const isSubmitting = ref(false)
const isUploading = ref(false)
const fileDokumen = ref(null)
const rejectReason = ref('')
const selectedId = ref('')

const form = ref({
  receipt_date: new Date().toISOString().split('T')[0],
  source_type: 'PROCUREMENT',
  supplier_name: '',
  invoice_number: '',
  funding_source: '',
  document_urls: [],
  items: []
})

const columns = [
  { name: 'bast_number', label: 'No. BAST', field: 'bast_number', align: 'left', sortable: true },
  { name: 'receipt_date', label: 'Tanggal', field: 'receipt_date', align: 'center', format: val => new Date(val).toLocaleDateString('id-ID') },
  { name: 'supplier_name', label: 'Supplier/Donor', field: 'supplier_name', align: 'left' },
  { name: 'source_type', label: 'Sumber', field: 'source_type', align: 'center' },
  { name: 'status', label: 'Status', field: 'status', align: 'center' },
  { name: 'aksi', label: 'Aksi', align: 'center' }
]

const getStatusColor = (s) => ({ 'DRAFT': 'amber', 'VERIFIED': 'green', 'REJECTED': 'red', 'DISTRIBUTED': 'blue' }[s] || 'grey')

onMounted(() => {
  store.fetchReceipts()
  kirStore.fetchMasterData() // Untuk ambil list katalog barang
})

watch(() => contextStore.selectedSemesterId, (newVal) => {
  if (newVal) store.fetchReceipts()
})

function bukaForm() {
  form.value = {
    receipt_date: new Date().toISOString().split('T')[0],
    source_type: 'PROCUREMENT',
    supplier_name: '',
    invoice_number: '',
    funding_source: '',
    document_urls: [],
    items: []
  }
  fileDokumen.value = null
  tambahBaris()
  showForm.value = true
}

function tambahBaris() {
  form.value.items.push({ catalog_id: '', item_name: '', specification: '', quantity: 1, unit: '', unit_price: 0, condition: 'BAIK', is_inventory: true })
}

function hapusBaris(index) {
  form.value.items.splice(index, 1)
}

function onPilihBarang(index, catalogId) {
  const cat = kirStore.catalogs.find(c => c.id === catalogId)
  if (cat) {
    form.value.items[index].item_name = cat.name
    form.value.items[index].unit = 'Unit' // Default
    // Kita asumsikan default INVENTORY, user bisa ubah toggle-nya
  }
}

async function uploadFile(file) {
  if (!file) return
  isUploading.value = true
  try {
    const formData = new FormData()
    formData.append('document', file)
    const url = await store.uploadDocument(formData)
    form.value.document_urls.push(url)
    $q.notify({ type: 'positive', message: 'Dokumen berhasil diupload' })
    fileDokumen.value = null
  } catch (e) {
    $q.notify({ type: 'negative', message: `Gagal upload dokumen ${e}` })
  } finally {
    isUploading.value = false
  }
}

async function simpanPenerimaan() {
  if (!form.value.supplier_name || form.value.items.length === 0) {
    $q.notify({ type: 'negative', message: 'Supplier dan Detail Barang wajib diisi' })
    return
  }
  isSubmitting.value = true
  try {
    const res = await store.createReceipt(form.value)
    $q.notify({ type: 'positive', message: res.message })
    showForm.value = false
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || 'Gagal menyimpan' })
  } finally {
    isSubmitting.value = false
  }
}

async function verifikasiData(row) {
  $q.dialog({
    title: 'Konfirmasi Verifikasi',
    message: `Verifikasi BAST ${row.bast_number}? Tindakan ini akan menambah stok / aset secara otomatis.`,
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await store.verifyReceipt(row.id)
      $q.notify({ type: 'positive', message: 'Penerimaan berhasil diverifikasi' })
    } catch (e) {
      $q.notify({ type: 'negative', message: e.response?.data?.error || 'Gagal verifikasi' })
    }
  })
}

function bukaTolak(row) {
  selectedId.value = row.id
  rejectReason.value = ''
  showTolak.value = true
}

async function submitTolak() {
  if (!rejectReason.value) {
    $q.notify({ type: 'negative', message: 'Alasan penolakan wajib diisi' })
    return
  }
  isSubmitting.value = true
  try {
    await store.rejectReceipt(selectedId.value, rejectReason.value)
    $q.notify({ type: 'info', message: 'Penerimaan ditolak' })
    showTolak.value = false
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || 'Gagal menolak' })
  } finally {
    isSubmitting.value = false
  }
}
</script>

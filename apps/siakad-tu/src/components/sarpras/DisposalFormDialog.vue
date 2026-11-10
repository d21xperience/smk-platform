<template>
  <q-dialog v-model="show" persistent maximized transition-show="slide-up" transition-hide="slide-down">
    <q-card class="column no-wrap">
      <q-card-section class="row items-center q-pb-none bg-primary text-white">
        <div class="text-h6">Form Usulan Penghapusan Aset</div>
        <q-space />
        <q-btn icon="close" flat round dense v-close-popup />
      </q-card-section>

      <q-card-section class="col scroll q-gutter-y-md">
        <!-- HEADER USULAN -->
        <div class="row q-col-gutter-md">
          <div class="col-12 col-md-6">
            <q-select outlined dense v-model="form.reason_category" 
              :options="[
                { label: 'Rusak Berat', value: 'RUSAK_BERAT' },
                { label: 'Usia > Masa Manfaat', value: 'USIA_MASA_MANFAAT' },
                { label: 'Tidak Ekonomis', value: 'TIDAK_EKONOMIS' },
                { label: 'Hilang', value: 'HILANG' },
                { label: 'Usang / Teknologi Tertinggal', value: 'USANG' },
                { label: 'Bencana Alam', value: 'BENCANA' }
              ]" 
              label="Kategori Alasan *" emit-value map-options />
          </div>
          <div class="col-12 col-md-6">
            <q-select outlined dense v-model="form.proposed_method" 
              :options="[
                { label: 'Dimusnahkan', value: 'DIMUSNAHKAN' },
                { label: 'Dijual', value: 'DIJUAL' },
                { label: 'Disumbangkan', value: 'DISUMBANGKAN' },
                { label: 'Ditukar', value: 'DITUKAR' },
                { label: 'Dimanfaatkan (Sparepart)', value: 'DIMANFAATKAN' }
              ]" 
              label="Usulan Metode Penghapusan *" emit-value map-options />
          </div>
        </div>

        <q-separator />

        <!-- DAFTAR ASET YANG DIUSULKAN -->
        <div class="row items-center justify-between">
          <div class="text-subtitle1 text-weight-bold">Daftar Aset yang Diusulkan</div>
          <q-btn color="secondary" icon="add" label="Tambah Aset" @click="bukaPilihAset" size="sm" no-caps />
        </div>

        <q-table :rows="form.items" :columns="columns" row-key="inventory_item_id" flat bordered :rows-per-page-options="[0]">
          <template v-slot:body-cell-actions="props">
            <q-td :props="props">
              <q-btn flat round dense color="negative" icon="delete" @click="hapusAset(props.rowIndex)" />
            </q-td>
          </template>
          <template v-slot:body-cell-specific_reason="props">
            <q-td :props="props">
              <q-input v-model="props.row.specific_reason" dense outlined placeholder="Alasan spesifik..." />
            </q-td>
          </template>
        </q-table>
      </q-card-section>

      <q-card-actions align="right" class="q-pa-md bg-grey-2">
        <q-btn flat label="Batal" v-close-popup />
        <q-btn color="primary" label="Kirim Usulan" @click="submit" :loading="isSubmitting" />
      </q-card-actions>
    </q-card>
  </q-dialog>

  <!-- Dialog Pilih Aset -->
  <q-dialog v-model="showPilihAset">
    <q-card style="width: 700px; max-width: 90vw;">
      <q-card-section class="text-h6">Pilih Aset untuk Diusulkan</q-card-section>
      <q-card-section>
        <q-table :rows="store.activeAssets" :columns="columnsPilih" row-key="id" selection="multiple" v-model:selected="selectedAssets" flat bordered :filter="filterText" :rows-per-page-options="[5, 10]">
          <template v-slot:top-right>
            <q-input borderless dense debounce="300" v-model="filterText" placeholder="Cari Kode/Nama...">
              <template v-slot:append><q-icon name="search" /></template>
            </q-input>
          </template>
        </q-table>
      </q-card-section>
      <q-card-actions align="right">
        <q-btn flat label="Batal" v-close-popup />
        <q-btn color="primary" label="Tambahkan" @click="konfirmasiPilihAset" />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useDisposalStore } from '@/stores/sarpras/disposal'

const props = defineProps({ modelValue: Boolean })
const emit = defineEmits(['update:modelValue', 'saved'])

const $q = useQuasar()
const store = useDisposalStore()

// PERBAIKAN: Gunakan computed untuk reaktivitas v-model
const show = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const showPilihAset = ref(false)
const isSubmitting = ref(false)
const filterText = ref('')
const selectedAssets = ref([])

const form = ref({
  reason_category: 'RUSAK_BERAT',
  proposed_method: 'DIMUSNAHKAN',
  items: []
})

const columns = [
  { name: 'asset_code', label: 'Kode Aset', field: 'asset_code', align: 'left' },
  { name: 'asset_name', label: 'Nama Barang', field: 'asset_name', align: 'left' },
  { name: 'specific_reason', label: 'Alasan Spesifik *', field: 'specific_reason', align: 'left', style: 'width: 300px' },
  { name: 'actions', label: 'Aksi', align: 'center' }
]

const columnsPilih = [
  { name: 'id', label: 'Kode Aset', field: 'id', align: 'left' },
  { name: 'name', label: 'Nama Barang', field: 'name', align: 'left' },
  { name: 'category_name', label: 'Kategori', field: 'category_name', align: 'left' }
]

onMounted(() => {
  store.fetchActiveAssets()
})

function bukaPilihAset() {
  selectedAssets.value = []
  showPilihAset.value = true
}

function konfirmasiPilihAset() {
  const existingIds = form.value.items.map(i => i.inventory_item_id)
  selectedAssets.value.forEach(asset => {
    if (!existingIds.includes(asset.id)) {
      form.value.items.push({
        inventory_item_id: asset.id,
        asset_code: asset.id,
        asset_name: asset.name,
        acquisition_price: asset.nilai_buku || 0,
        specific_reason: '',
        photo_urls: '[]'
      })
    }
  })
  showPilihAset.value = false
}

function hapusAset(index) {
  form.value.items.splice(index, 1)
}

async function submit() {
  if (form.value.items.length === 0) {
    $q.notify({ type: 'negative', message: 'Minimal 1 aset harus diusulkan' })
    return
  }
  const missingReason = form.value.items.some(i => !i.specific_reason)
  if (missingReason) {
    $q.notify({ type: 'negative', message: 'Alasan spesifik wajib diisi untuk semua aset' })
    return
  }

  isSubmitting.value = true
  try {
    await store.createRequest(form.value)
    $q.notify({ type: 'positive', message: 'Usulan berhasil dibuat' })
    show.value = false
    form.value.items = []
    emit('saved')
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || 'Gagal membuat usulan' })
  } finally {
    isSubmitting.value = false
  }
}
</script>

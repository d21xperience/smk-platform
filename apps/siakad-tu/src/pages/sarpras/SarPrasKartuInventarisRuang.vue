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
          <q-btn color="secondary" icon="print" label="Cetak Lembar KIR Pintu" :disabled="!ruanganTerpilih" @click="cetakKIR" />
        </div>
      </q-card-section>
    </q-card>

    <!-- KONTEN DATA KIR -->
    <div v-if="ruanganTerpilih">
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
            <div class="text-caption text-grey-7 text-right-gt-xs">Total Inventaris Ruang</div>
            <div class="text-subtitle1 text-weight-bold text-primary text-right-gt-xs">{{ kirStore.items.length }} Unit Barang</div>
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
            <q-btn flat round dense color="primary" icon="edit" @click="bukaDialogEdit(props.row)" />
            <q-btn flat round dense color="negative" icon="delete" @click="konfirmasiHapus(props.row)" />
          </q-td>
        </template>
      </q-table>
    </div>

    <!-- DIALOG FORM CRUD -->
    <q-dialog v-model="showDialog" persistent>
      <q-card style="min-width: 450px">
        <q-card-section class="text-h6">{{ isEdit ? 'Edit Barang Inventaris' : 'Tambah Barang Inventaris' }}</q-card-section>
        <q-card-section class="q-gutter-y-md">
          <!-- Jika Edit, ID tidak bisa diubah -->
          <q-input v-if="isEdit" v-model="formData.id" label="Kode Unik Aset" outlined dense disable />
          
          <!-- Jika Tambah, Pilih Katalog Barang -->
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
import { ref, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useKirStore } from '@/stores/sarpras/kir'

const $q = useQuasar()
const kirStore = useKirStore()
const ruanganTerpilih = ref(null)

const showDialog = ref(false)
const isEdit = ref(false)
const formData = ref({ id: '', catalog_id: '', purchase_date: new Date().toISOString().split('T')[0], condition: 'Bagus' })

const columns = [
  { name: 'id', label: 'Kode Unik Aset (Kodifikasi)', align: 'left', field: 'id', sortable: true },
  { name: 'name', label: 'Nama Barang', align: 'left', field: 'name', sortable: true },
  { name: 'brand', label: 'Merk/Brand', align: 'left', field: 'brand' },
  { name: 'category_name', label: 'Kategori', align: 'center', field: 'category_name' },
  { name: 'condition', label: 'Kondisi Fisik', align: 'center', field: 'condition', sortable: true },
  { name: 'actions', label: 'Aksi', align: 'center', field: 'id' }
]

onMounted(() => { 
  kirStore.fetchRooms()
  kirStore.fetchMasterData()
})

function onRuanganChange(val) { if (val) kirStore.fetchItemsByRoom(val.id) }

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

function cetakKIR() {
  $q.notify({ message: `Menyiapkan cetakan lembar KIR untuk ${ruanganTerpilih.value.nama}...`, color: 'secondary', icon: 'local_printshop' })
}
</script>

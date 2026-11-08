<template>
  <q-page class="q-pa-md">
    <q-card flat bordered>
      <q-tabs v-model="tab" dense class="text-grey" active-color="primary" indicator-color="primary" align="left"
        narrow-indicator>
        <q-tab name="kategori" label="Kategori Barang" />
        <q-tab name="katalog" label="Katalog Barang" />
      </q-tabs>

      <q-separator />

      <q-tab-panels v-model="tab" animated>
        <!-- TAB KATEGORI -->
        <q-tab-panel name="kategori">
          <div class="row q-mb-md justify-between items-center">
            <div class="text-h6">Daftar Kategori</div>
            <q-btn color="primary" icon="add" label="Tambah Kategori" @click="bukaDialogKategori()" no-caps />
          </div>
          <q-table :rows="kirStore.categories" :columns="columnsKategori" row-key="id" flat bordered
            :loading="kirStore.isLoading">
            <template v-slot:body-cell-actions="props">
              <q-td :props="props">
                <q-btn flat round dense color="primary" icon="edit" @click="bukaDialogKategori(props.row)" />
                <q-btn flat round dense color="negative" icon="delete" @click="hapusKategori(props.row)" />
              </q-td>
            </template>
          </q-table>
        </q-tab-panel>

        <!-- TAB KATALOG -->
        <q-tab-panel name="katalog">
          <div class="row q-mb-md justify-between items-center">
            <div class="text-h6">Daftar Katalog (Blueprint Barang)</div>
            <q-btn color="primary" icon="add" label="Tambah Katalog" @click="bukaDialogKatalog()" no-caps />
          </div>
          <q-table :rows="kirStore.catalogs" :columns="columnsKatalog" row-key="id" flat bordered
            :loading="kirStore.isLoading">
            <template v-slot:body-cell-image_url="props">
              <q-td :props="props">
                <q-avatar v-if="props.value" size="40px" class="shadow-1">
                  <img :src="props.value" />
                </q-avatar>
                <q-icon v-else name="image_not_supported" size="24px" color="grey-4" />
              </q-td>
            </template>
            <template v-slot:body-cell-category_id="props">
              <q-td :props="props">{{ getCategoryName(props.value) }}</q-td>
            </template>
            <template v-slot:body-cell-actions="props">
              <q-td :props="props">
                <q-btn flat round dense color="primary" icon="edit" @click="bukaDialogKatalog(props.row)" />
                <q-btn flat round dense color="negative" icon="delete" @click="hapusKatalog(props.row)" />
              </q-td>
            </template>
          </q-table>
        </q-tab-panel>
      </q-tab-panels>
    </q-card>

    <!-- Dialog Kategori -->
    <q-dialog v-model="showDialogKategori" persistent>
      <q-card style="min-width: 350px">
        <q-card-section class="text-h6">{{ isEditKategori ? 'Edit Kategori' : 'Tambah Kategori' }}</q-card-section>
        <q-card-section class="q-gutter-y-md">
          <q-input v-model="formKategori.code" label="Kode Kategori (Cth: ELN)" outlined dense />
          <q-input v-model="formKategori.name" label="Nama Kategori" outlined dense />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn color="primary" label="Simpan" @click="submitKategori" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Dialog Katalog -->
    <q-dialog v-model="showDialogKatalog" persistent>
      <q-card style="min-width: 400px">
        <q-card-section class="text-h6">{{ isEditKatalog ? 'Edit Katalog' : 'Tambah Katalog' }}</q-card-section>
        <q-card-section class="q-gutter-y-md">
          <q-select v-model="formKatalog.category_id" :options="kirStore.categories" option-value="id"
            option-label="name" label="Pilih Kategori" outlined dense emit-value map-options />
          <q-input v-model="formKatalog.name" label="Nama Barang (Cth: Proyektor Epson)" outlined dense />
          <q-input v-model="formKatalog.brand" label="Merk (Opsional)" outlined dense />

          <!-- FITUR UPLOAD FOTO -->
          <div class="row items-center q-gutter-md">
            <q-avatar v-if="formKatalog.image_url" size="60px" class="shadow-1">
              <img :src="formKatalog.image_url" />
            </q-avatar>
            <q-file v-model="imageFile" label="Upload Foto Barang" outlined dense accept="image/*"
              @update:model-value="uploadImage" :loading="isUploading">
              <template v-slot:prepend><q-icon name="cloud_upload" /></template>
            </q-file>
          </div>
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn color="primary" label="Simpan" @click="submitKatalog" :loading="isSubmitting" />
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
const tab = ref('kategori')

// State Kategori
const showDialogKategori = ref(false)
const isEditKategori = ref(false)
const formKategori = ref({ id: '', code: '', name: '' })

// State Katalog
const showDialogKatalog = ref(false)
const isEditKatalog = ref(false)
const formKatalog = ref({ id: '', category_id: '', name: '', brand: '', image_url: '' })
const imageFile = ref(null)
const isUploading = ref(false)
const isSubmitting = ref(false)

const columnsKategori = [
  { name: 'code', label: 'Kode Kategori', field: 'code', align: 'left', sortable: true },
  { name: 'name', label: 'Nama Kategori', field: 'name', align: 'left', sortable: true },
  { name: 'actions', label: 'Aksi', field: 'id', align: 'center' }
]

const columnsKatalog = [
  { name: 'image_url', label: 'Foto', field: 'image_url', align: 'center', style: 'width: 60px' },
  { name: 'name', label: 'Nama Barang', field: 'name', align: 'left', sortable: true },
  { name: 'brand', label: 'Merk', field: 'brand', align: 'left' },
  { name: 'category_id', label: 'Kategori', field: 'category_id', align: 'left' },
  { name: 'actions', label: 'Aksi', field: 'id', align: 'center' }
]

onMounted(() => kirStore.fetchMasterData())

const getCategoryName = (id) => {
  const cat = kirStore.categories.find(c => c.id === id)
  return cat ? cat.name : '-'
}

// --- HANDLER KATEGORI ---
function bukaDialogKategori(item = null) {
  if (item) {
    isEditKategori.value = true
    formKategori.value = { ...item }
  } else {
    isEditKategori.value = false
    formKategori.value = { id: '', code: '', name: '' }
  }
  showDialogKategori.value = true
}

async function submitKategori() {
  try {
    if (isEditKategori.value) {
      await kirStore.updateCategory(formKategori.value.id, formKategori.value)
      $q.notify({ type: 'positive', message: 'Kategori diupdate' })
    } else {
      await kirStore.createCategory(formKategori.value)
      $q.notify({ type: 'positive', message: 'Kategori ditambahkan' })
    }
    showDialogKategori.value = false
    kirStore.fetchMasterData()
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || 'Error' })
  }
}

function hapusKategori(item) {
  $q.dialog({ title: 'Konfirmasi', message: `Hapus kategori ${item.name}?`, cancel: true, persistent: true })
    .onOk(async () => {
      await kirStore.deleteCategory(item.id)
      $q.notify({ type: 'positive', message: 'Kategori dihapus' })
      kirStore.fetchMasterData()
    })
}

// --- HANDLER KATALOG ---
function bukaDialogKatalog(item = null) {
  imageFile.value = null
  if (item) {
    isEditKatalog.value = true
    formKatalog.value = { ...item }
  } else {
    isEditKatalog.value = false
    formKatalog.value = { id: '', category_id: '', name: '', brand: '', image_url: '' }
  }
  showDialogKatalog.value = true
}

async function uploadImage(file) {
  if (!file) return
  isUploading.value = true
  try {
    const formData = new FormData()
    formData.append('image', file)
    const url = await kirStore.uploadCatalogImage(formData)
    formKatalog.value.image_url = url
    $q.notify({ type: 'positive', message: 'Foto berhasil diupload' })
  } catch (e) {
    $q.notify({ type: 'negative', message: `Gagal upload foto ${e}` })
    imageFile.value = null
  } finally {
    isUploading.value = false
  }
}

async function submitKatalog() {
  isSubmitting.value = true
  try {
    if (isEditKatalog.value) {
      await kirStore.updateCatalog(formKatalog.value.id, formKatalog.value)
      $q.notify({ type: 'positive', message: 'Katalog diupdate' })
    } else {
      await kirStore.createCatalog(formKatalog.value)
      $q.notify({ type: 'positive', message: 'Katalog ditambahkan' })
    }
    showDialogKatalog.value = false
    kirStore.fetchMasterData()
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || 'Error' })
  } finally {
    isSubmitting.value = false
  }
}

function hapusKatalog(item) {
  $q.dialog({ title: 'Konfirmasi', message: `Hapus katalog ${item.name}?`, cancel: true, persistent: true })
    .onOk(async () => {
      await kirStore.deleteCatalog(item.id)
      $q.notify({ type: 'positive', message: 'Katalog dihapus' })
      kirStore.fetchMasterData()
    })
}
</script>

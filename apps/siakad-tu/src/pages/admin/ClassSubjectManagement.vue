<template>
  <q-page class="q-pa-md">
    <div class="text-h5 q-mb-md">Manajemen Kelas & Mata Pelajaran</div>

    <q-tabs v-model="tab" dense class="q-mb-md" align="left">
      <q-tab name="kelas" label="Kelas" />
      <q-tab name="mapel" label="Mata Pelajaran" />
    </q-tabs>

    <q-tab-panels v-model="tab" animated>
      <!-- Tab Kelas -->
      <q-tab-panel name="kelas" class="q-pa-none">
        <div class="q-mb-md">
          <q-btn color="primary" icon="add" label="Tambah Kelas" @click="openFormKelas(null)" />
        </div>
        <q-table :rows="kelasList" :columns="kelasColumns" row-key="id" flat bordered dense>
          <template v-slot:body-cell-actions="props">
            <q-td :props="props">
              <q-btn flat round dense icon="edit" color="info" @click="openFormKelas(props.row)" class="q-mr-sm" />
              <q-btn flat round dense icon="delete" color="negative" @click="confirmDeleteKelas(props.row)" />
            </q-td>
          </template>
        </q-table>
      </q-tab-panel>

      <!-- Tab Mata Pelajaran -->
      <q-tab-panel name="mapel" class="q-pa-none">
        <div class="q-mb-md">
          <q-btn color="primary" icon="add" label="Tambah Mata Pelajaran" @click="openFormMapel(null)" />
        </div>
        <q-table :rows="mapelList" :columns="mapelColumns" row-key="id" flat bordered dense>
          <template v-slot:body-cell-actions="props">
            <q-td :props="props">
              <q-btn flat round dense icon="edit" color="info" @click="openFormMapel(props.row)" class="q-mr-sm" />
              <q-btn flat round dense icon="delete" color="negative" @click="confirmDeleteMapel(props.row)" />
            </q-td>
          </template>
        </q-table>
      </q-tab-panel>
    </q-tab-panels>

    <!-- Dialog Form Kelas -->
    <q-dialog v-model="kelasDialog" persistent>
      <q-card style="min-width: 400px;">
        <q-card-section>
          <div class="text-h6">{{ kelasFormTitle }}</div>
        </q-card-section>
        <q-card-section>
          <q-form @submit="submitKelas" class="q-gutter-md">
            <q-input v-model="kelasForm.nama" label="Nama Kelas (contoh: 10 IPA 1)" lazy-rules
              :rules="[val => !!val || 'Nama kelas harus diisi']" outlined autofocus />
            <div class="q-gutter-sm">
              <q-btn label="Simpan" type="submit" color="primary" />
              <q-btn label="Batal" flat color="negative" v-close-popup />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Dialog Form Mata Pelajaran -->
    <q-dialog v-model="mapelDialog" persistent>
      <q-card style="min-width: 400px;">
        <q-card-section>
          <div class="text-h6">{{ mapelFormTitle }}</div>
        </q-card-section>
        <q-card-section>
          <q-form @submit="submitMapel" class="q-gutter-md">
            <q-input v-model="mapelForm.nama" label="Nama Mata Pelajaran" lazy-rules
              :rules="[val => !!val || 'Nama mata pelajaran harus diisi']" outlined autofocus />
            <q-input v-model="mapelForm.kode" label="Kode (opsional)" outlined />
            <div class="q-gutter-sm">
              <q-btn label="Simpan" type="submit" color="primary" />
              <q-btn label="Batal" flat color="negative" v-close-popup />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Dialog Konfirmasi Hapus -->
    <q-dialog v-model="deleteDialog" persistent>
      <q-card>
        <q-card-section>
          <div class="text-h6">Konfirmasi Hapus</div>
        </q-card-section>
        <q-card-section>
          Apakah Anda yakin ingin menghapus <strong>{{ deleteItem?.nama || deleteItem?.nama_mapel }}</strong>?
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" color="negative" v-close-popup />
          <q-btn flat label="Hapus" color="primary" @click="executeDelete" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()

// Data dummy kelas
const kelasList = ref([
  { id: 1, nama: '10 IPA 1' },
  { id: 2, nama: '10 IPA 2' },
  { id: 3, nama: '10 IPS 1' },
  { id: 4, nama: '11 IPA 1' },
  { id: 5, nama: '11 IPA 2' },
  { id: 6, nama: '11 IPS 1' },
  { id: 7, nama: '12 IPA 1' },
  { id: 8, nama: '12 IPA 2' },
  { id: 9, nama: '12 IPS 1' }
])

// Data dummy mata pelajaran
const mapelList = ref([
  { id: 1, nama: 'Matematika', kode: 'MATH01' },
  { id: 2, nama: 'Fisika', kode: 'PHY01' },
  { id: 3, nama: 'Kimia', kode: 'CHEM01' },
  { id: 4, nama: 'Biologi', kode: 'BIO01' },
  { id: 5, nama: 'Bahasa Indonesia', kode: 'IND01' },
  { id: 6, nama: 'Bahasa Inggris', kode: 'ENG01' },
  { id: 7, nama: 'Sejarah', kode: 'HIS01' },
  { id: 8, nama: 'Geografi', kode: 'GEO01' },
  { id: 9, nama: 'Ekonomi', kode: 'ECO01' }
])

const tab = ref('kelas')

// Kelas form state
const kelasDialog = ref(false)
const isEditKelas = ref(false)
const editKelasId = ref(null)
const kelasForm = ref({ nama: '' })
const kelasFormTitle = ref('Tambah Kelas')

// Mapel form state
const mapelDialog = ref(false)
const isEditMapel = ref(false)
const editMapelId = ref(null)
const mapelForm = ref({ nama: '', kode: '' })
const mapelFormTitle = ref('Tambah Mata Pelajaran')

// Delete state
const deleteDialog = ref(false)
const deleteType = ref('') // 'kelas' atau 'mapel'
const deleteItem = ref(null)

const kelasColumns = [
  { name: 'nama', label: 'Nama Kelas', field: 'nama', align: 'left' },
  { name: 'actions', label: 'Aksi', field: 'actions', align: 'center' }
]

const mapelColumns = [
  { name: 'nama', label: 'Mata Pelajaran', field: 'nama', align: 'left' },
  { name: 'kode', label: 'Kode', field: 'kode', align: 'left' },
  { name: 'actions', label: 'Aksi', field: 'actions', align: 'center' }
]

// Kelas methods
function openFormKelas(kelas) {
  if (kelas) {
    isEditKelas.value = true
    editKelasId.value = kelas.id
    kelasForm.value = { nama: kelas.nama }
    kelasFormTitle.value = 'Edit Kelas'
  } else {
    isEditKelas.value = false
    editKelasId.value = null
    kelasForm.value = { nama: '' }
    kelasFormTitle.value = 'Tambah Kelas'
  }
  kelasDialog.value = true
}

function submitKelas() {
  if (!kelasForm.value.nama) return
  if (isEditKelas.value) {
    const index = kelasList.value.findIndex(k => k.id === editKelasId.value)
    if (index !== -1) {
      kelasList.value[index].nama = kelasForm.value.nama
      $q.notify({ type: 'positive', message: 'Kelas berhasil diupdate' })
    }
  } else {
    const newId = Date.now()
    kelasList.value.push({ id: newId, nama: kelasForm.value.nama })
    $q.notify({ type: 'positive', message: 'Kelas berhasil ditambahkan' })
  }
  kelasDialog.value = false
}

function confirmDeleteKelas(kelas) {
  deleteType.value = 'kelas'
  deleteItem.value = kelas
  deleteDialog.value = true
}

// Mapel methods
function openFormMapel(mapel) {
  if (mapel) {
    isEditMapel.value = true
    editMapelId.value = mapel.id
    mapelForm.value = { nama: mapel.nama, kode: mapel.kode || '' }
    mapelFormTitle.value = 'Edit Mata Pelajaran'
  } else {
    isEditMapel.value = false
    editMapelId.value = null
    mapelForm.value = { nama: '', kode: '' }
    mapelFormTitle.value = 'Tambah Mata Pelajaran'
  }
  mapelDialog.value = true
}

function submitMapel() {
  if (!mapelForm.value.nama) return
  if (isEditMapel.value) {
    const index = mapelList.value.findIndex(m => m.id === editMapelId.value)
    if (index !== -1) {
      mapelList.value[index].nama = mapelForm.value.nama
      mapelList.value[index].kode = mapelForm.value.kode
      $q.notify({ type: 'positive', message: 'Mata pelajaran berhasil diupdate' })
    }
  } else {
    const newId = Date.now()
    mapelList.value.push({ id: newId, nama: mapelForm.value.nama, kode: mapelForm.value.kode })
    $q.notify({ type: 'positive', message: 'Mata pelajaran berhasil ditambahkan' })
  }
  mapelDialog.value = false
}

function confirmDeleteMapel(mapel) {
  deleteType.value = 'mapel'
  deleteItem.value = mapel
  deleteDialog.value = true
}

// Delete execution
function executeDelete() {
  if (deleteType.value === 'kelas') {
    const index = kelasList.value.findIndex(k => k.id === deleteItem.value.id)
    if (index !== -1) kelasList.value.splice(index, 1)
    $q.notify({ type: 'positive', message: 'Kelas berhasil dihapus' })
  } else if (deleteType.value === 'mapel') {
    const index = mapelList.value.findIndex(m => m.id === deleteItem.value.id)
    if (index !== -1) mapelList.value.splice(index, 1)
    $q.notify({ type: 'positive', message: 'Mata pelajaran berhasil dihapus' })
  }
  deleteDialog.value = false
  deleteItem.value = null
}
</script>

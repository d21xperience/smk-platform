<template>
  <q-page class="q-pa-md">
    <div class="row justify-between q-mb-md">
      <div class="text-h5">Database Mitra Industri</div>
      <q-btn color="primary" icon="add" label="Tambah Mitra" @click="openDialog" />
    </div>

    <q-table :rows="mitraStore.list" :columns="columns" row-key="id" flat bordered dense>
      <template v-slot:body-cell-aksi="props">
        <q-td :props="props">
          <q-btn flat round dense icon="edit" color="info" @click="openDialog(props.row)" class="q-mr-sm" />
          <q-btn flat round dense icon="delete" color="negative" @click="confirmDelete(props.row.id)" />
        </q-td>
      </template>
    </q-table>

    <q-dialog v-model="dialog" persistent>
      <q-card style="min-width: 500px">
        <q-card-section>
          <div class="text-h6">{{ editMode ? 'Edit Mitra' : 'Tambah Mitra' }}</div>
        </q-card-section>
        <q-card-section>
          <q-form @submit="submitForm" class="q-gutter-md">
            <q-input v-model="form.nama" label="Nama Perusahaan" outlined :rules="[v => !!v]" />
            <q-input v-model="form.bidang" label="Bidang Usaha" outlined />
            <q-input v-model="form.alamat" label="Alamat" outlined />
            <q-input v-model.number="form.kuota" label="Kuota Siswa" type="number" outlined />
            <q-input v-model="form.kontak" label="Kontak (Telp/Email)" outlined />
            <q-input v-model.number="form.rating" label="Rating (0-5)" type="number" step="0.1" min="0" max="5"
              outlined />
            <div class="q-gutter-sm">
              <q-btn label="Simpan" type="submit" color="primary" />
              <q-btn label="Batal" flat color="negative" v-close-popup @click="resetForm" />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useMitraStore } from 'stores/pkl/mitra'

const $q = useQuasar()
const mitraStore = useMitraStore()

const dialog = ref(false)
const editMode = ref(false)
const editId = ref(null)
const form = ref({ nama: '', bidang: '', alamat: '', kuota: 0, kontak: '', rating: 0 })

const columns = [
  { name: 'nama', label: 'Nama Mitra', field: 'nama', align: 'left' },
  { name: 'bidang', label: 'Bidang', field: 'bidang', align: 'left' },
  { name: 'alamat', label: 'Alamat', field: 'alamat', align: 'left' },
  { name: 'kuota', label: 'Kuota', field: 'kuota', align: 'center' },
  { name: 'kontak', label: 'Kontak', field: 'kontak', align: 'left' },
  { name: 'rating', label: 'Rating', field: 'rating', align: 'center' },
  { name: 'aksi', label: 'Aksi', field: 'aksi', align: 'center' }
]

function openDialog(row = null) {
  if (row) {
    editMode.value = true
    editId.value = row.id
    form.value = { ...row }
  } else {
    editMode.value = false
    resetForm()
  }
  dialog.value = true
}

function resetForm() {
  form.value = { nama: '', bidang: '', alamat: '', kuota: 0, kontak: '', rating: 0 }
}

function submitForm() {
  if (editMode.value) {
    mitraStore.update(editId.value, form.value)
    $q.notify({ type: 'positive', message: 'Mitra diupdate' })
  } else {
    mitraStore.tambah(form.value)
    $q.notify({ type: 'positive', message: 'Mitra ditambahkan' })
  }
  dialog.value = false
  resetForm()
}

function confirmDelete(id) {
  $q.dialog({ title: 'Konfirmasi', message: 'Hapus mitra ini?', cancel: true }).onOk(() => {
    mitraStore.hapus(id)
    $q.notify({ type: 'positive', message: 'Mitra dihapus' })
  })
}

onMounted(() => mitraStore.loadData())
</script>

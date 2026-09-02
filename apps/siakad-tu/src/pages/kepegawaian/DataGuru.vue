<template>
  <q-page class="q-pa-md">
    <div class="row justify-between q-mb-md">
      <div class="text-h5">Data Induk GTK (Guru/Tenaga Kependidikan)</div>
      <q-btn color="primary" icon="add" label="Tambah Guru" @click="openDialog()" />
    </div>

    <q-table :rows="guruStore.list" :columns="columns" row-key="id" flat bordered dense>
      <template v-slot:body-cell-aksi="props">
        <q-td :props="props">
          <q-btn flat round dense icon="edit" color="info" @click="openDialog(props.row)" class="q-mr-sm" />
          <q-btn flat round dense icon="delete" color="negative" @click="confirmDelete(props.row.id)" />
        </q-td>
      </template>
    </q-table>

    <!-- Dialog form tambah/edit -->
    <q-dialog v-model="dialog" persistent>
      <q-card style="min-width: 500px">
        <q-card-section>
          <div class="text-h6">{{ editMode ? 'Edit Guru' : 'Tambah Guru' }}</div>
        </q-card-section>
        <q-card-section>
          <q-form @submit="submitForm" class="q-gutter-md">
            <q-input v-model="form.nama" label="Nama Lengkap" outlined :rules="[v => !!v || 'Wajib diisi']" />
            <q-input v-model="form.nuptk" label="NUPTK" outlined mask="################"
              :rules="[v => v && v.length === 16 || '16 digit']" />
            <q-input v-model="form.jabatan" label="Jabatan" outlined />
            <q-input v-model.number="form.masa_tugas" label="Masa Tugas (tahun)" type="number" outlined />
            <q-select v-model="form.status" :options="['aktif', 'nonaktif']" label="Status" outlined />
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
import { useGuruStore } from 'stores/kepegawaian/guru'

const $q = useQuasar()
const guruStore = useGuruStore()

const dialog = ref(false)
const editMode = ref(false)
const editId = ref(null)
const form = ref({
  nama: '',
  nuptk: '',
  jabatan: '',
  masa_tugas: 0,
  status: 'aktif'
})

const columns = [
  { name: 'nama', label: 'Nama', field: 'nama', align: 'left' },
  { name: 'nuptk', label: 'NUPTK', field: 'nuptk', align: 'left' },
  { name: 'jabatan', label: 'Jabatan', field: 'jabatan', align: 'left' },
  { name: 'masa_tugas', label: 'Masa Tugas (thn)', field: 'masa_tugas', align: 'center' },
  { name: 'status', label: 'Status', field: 'status', align: 'center' },
  { name: 'aksi', label: 'Aksi', field: 'aksi', align: 'center' }
]

function openDialog(row = null) {
  if (row) {
    editMode.value = true
    editId.value = row.id
    form.value = { ...row }
  } else {
    editMode.value = false
    editId.value = null
    resetForm()
  }
  dialog.value = true
}

function resetForm() {
  form.value = {
    nama: '',
    nuptk: '',
    jabatan: '',
    masa_tugas: 0,
    status: 'aktif'
  }
}

function submitForm() {
  if (editMode.value) {
    guruStore.update(editId.value, form.value)
    $q.notify({ type: 'positive', message: 'Data guru diupdate' })
  } else {
    guruStore.tambah(form.value)
    $q.notify({ type: 'positive', message: 'Guru baru ditambahkan' })
  }
  dialog.value = false
  resetForm()
}

function confirmDelete(id) {
  $q.dialog({
    title: 'Konfirmasi',
    message: 'Hapus data guru ini?',
    cancel: true,
    persistent: true
  }).onOk(() => {
    guruStore.hapus(id)
    $q.notify({ type: 'positive', message: 'Data guru dihapus' })
  })
}

onMounted(() => {
  guruStore.loadData()
})
</script>

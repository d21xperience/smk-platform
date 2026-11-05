<template>
  <q-page class="q-pa-md bg-grey-2">
    <div class="row justify-between items-center q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Data GTK</h1>
        <p class="text-caption text-grey-7 q-mb-none">Guru & Tenaga Kependidikan (DAPODIK)</p>
      </div>
      <q-btn color="primary" icon="add" label="Tambah Guru" @click="openDialog()" />
    </div>

    <q-card class="shadow-1">
      <q-table :rows="guruStore.list" :columns="columns" row-key="id" :loading="guruStore.isLoading" flat bordered
        dense>
        <template v-slot:body-cell-gender="props">
          <q-td :props="props" class="text-center">
            <q-badge :color="props.value === 'L' ? 'blue-6' : 'pink-6'">
              {{ props.value === 'L' ? 'Laki-laki' : 'Perempuan' }}
            </q-badge>
          </q-td>
        </template>

        <template v-slot:body-cell-aksi="props">
          <q-td :props="props" class="text-center">
            <q-btn flat round dense icon="edit" color="info" @click="openDialog(props.row)" class="q-mr-sm" />
            <q-btn flat round dense icon="delete" color="negative" @click="confirmDelete(props.row)" />
          </q-td>
        </template>
      </q-table>
    </q-card>

    <!-- Dialog Form Tambah/Edit -->
    <q-dialog v-model="dialogOpen" persistent>
      <q-card style="min-width: 600px">
        <q-card-section class="row items-center q-pb-none bg-primary text-white">
          <div class="text-h6">{{ isEdit ? 'Edit Data Guru' : 'Tambah Data Guru' }}</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section>
          <q-form @submit="onSubmit" class="q-gutter-md">
            <div class="row q-col-gutter-sm">
              <div class="col-4">
                <q-input v-model="form.gelarDepan" label="Gelar Depan" hint="Cth: Dr." filled />
              </div>
              <div class="col-8">
                <q-input v-model="form.name" label="Nama Lengkap (Tanpa Gelar)" filled
                  :rules="[val => !!val || 'Wajib diisi']" />
              </div>
            </div>
            <div class="row q-col-gutter-sm">
              <div class="col-8">
                <q-input v-model="form.gelarBelakang" label="Gelar Belakang" hint="Cth: S.Kom., M.Pd." filled />
              </div>
              <div class="col-4">
                <q-select v-model="form.gender"
                  :options="[{ label: 'Laki-laki', value: 'L' }, { label: 'Perempuan', value: 'P' }]" emit-value map-options
                  label="Jenis Kelamin" filled :rules="[val => !!val || 'Wajib diisi']" />
              </div>
            </div>

            <div class="row q-col-gutter-sm">
              <div class="col-6">
                <q-input v-model="form.nip" label="NIP" filled mask="##################" />
              </div>
              <div class="col-6">
                <q-input v-model="form.nuptk" label="NUPTK" filled mask="################" />
              </div>
            </div>

            <div class="row q-col-gutter-sm">
              <div class="col-6">
                <q-input v-model="form.birthPlace" label="Tempat Lahir" filled />
              </div>
              <div class="col-6">
                <q-input v-model="form.birthDate" label="Tanggal Lahir" filled mask="####-##-##" hint="YYYY-MM-DD" />
              </div>
            </div>

            <q-input v-model="form.address" label="Alamat" type="textarea" filled />

            <div class="row justify-end q-mt-lg">
              <q-btn label="Batal" color="grey" flat v-close-popup class="q-mr-sm" />
              <q-btn label="Simpan Data" color="primary" type="submit" :loading="guruStore.isLoading" />
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
import { useGuruStore } from '@/stores/kepegawaian/guru'

const $q = useQuasar()
const guruStore = useGuruStore()

const dialogOpen = ref(false)
const isEdit = ref(false)
const form = ref({
  id: '',
  name: '',
  gelarDepan: '',
  gelarBelakang: '',
  nip: '',
  nuptk: '',
  gender: 'L',
  birthPlace: '',
  birthDate: '',
  address: ''
})

// Sesuaikan kolom tabel dengan format camelCase dari Backend
const columns = [
  { name: 'name', label: 'Nama Lengkap & Gelar', field: 'name', align: 'left', sortable: true },
  { name: 'nip', label: 'NIP', field: 'nip', align: 'left' },
  { name: 'nuptk', label: 'NUPTK', field: 'nuptk', align: 'left' },
  { name: 'gender', label: 'JK', field: 'gender', align: 'center' },
  { name: 'birthPlace', label: 'Tempat Lahir', field: 'birthPlace', align: 'left' },
  { name: 'aksi', label: 'Aksi', align: 'center' }
]

onMounted(() => {
  guruStore.fetchTeachers()
})

const openDialog = (row = null) => {
  if (row) {
    isEdit.value = true
    // Konversi format tanggal dari backend (ISO) ke string YYYY-MM-DD
    const tgl = row.birthDate ? new Date(row.birthDate).toISOString().split('T')[0] : ''

    // FIX: Petakan data dari row (camelCase) ke form
    form.value = {
      id: row.id,
      name: row.name,
      gelarDepan: row.gelarDepan,
      gelarBelakang: row.gelarBelakang,
      nip: row.nip,
      nuptk: row.nuptk,
      gender: row.gender,
      birthPlace: row.birthPlace,
      birthDate: tgl,
      address: row.address
    }
  } else {
    isEdit.value = false
    form.value = {
      id: '', name: '', gelarDepan: '', gelarBelakang: '',
      nip: '', nuptk: '', gender: 'L', birthPlace: '', birthDate: '', address: ''
    }
  }
  dialogOpen.value = true
}

const onSubmit = async () => {
  try {
    // Pastikan payload sesuai dengan yang diharapkan backend
    const payload = {
      name: form.value.name,
      gelarDepan: form.value.gelarDepan,
      gelarBelakang: form.value.gelarBelakang,
      nip: form.value.nip,
      nuptk: form.value.nuptk,
      gender: form.value.gender,
      birthPlace: form.value.birthPlace,
      birthDate: form.value.birthDate,
      address: form.value.address
    }

    if (isEdit.value) {
      await guruStore.updateTeacher(form.value.id, payload)
      $q.notify({ color: 'positive', message: 'Data Guru diperbarui' })
    } else {
      await guruStore.createTeacher(payload)
      $q.notify({ color: 'positive', message: 'Guru baru ditambahkan' })
    }
    dialogOpen.value = false
  } catch (err) {
    $q.notify({ color: 'negative', message: err.response?.data?.error || 'Terjadi kesalahan' })
  }
}

const confirmDelete = (row) => {
  $q.dialog({
    title: 'Konfirmasi Hapus',
    message: `Hapus data guru ${row.name}?`,
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await guruStore.deleteTeacher(row.id)
      $q.notify({ color: 'positive', message: 'Data guru dihapus' })
    } catch (err) {
      $q.notify({ color: 'negative', message: `Gagal menghapus data guru ${err}` })
    }
  })
}
</script>

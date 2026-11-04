<template>
  <q-page class="q-pa-md bg-grey-2">
    <!-- Header Halaman -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Manajemen Data Kelas</h1>
        <p class="text-caption text-grey-7 q-mb-none">Kelola data rombongan belajar sesuai semester aktif</p>
      </div>
      <div class="q-gutter-x-sm">
        <q-btn color="primary" icon="add" label="Tambah Kelas" @click="openDialog()" />
      </div>
    </div>

    <!-- Tabel Kelas -->
    <q-card class="shadow-1">
      <q-table :rows="kelasStore.list" :columns="columns" row-key="id" :loading="kelasStore.isLoading" flat bordered>
        <!-- Search Bar -->
        <template v-slot:top-right>
          <q-input v-model="filterKeyword" placeholder="Cari kelas..." dense outlined clearable
            @update:model-value="searchData">
            <template v-slot:append>
              <q-icon name="search" />
            </template>
          </q-input>
        </template>

        <!-- Custom Render: Level/Grade -->
        <template v-slot:body-cell-level="props">
          <q-td :props="props">
            <q-badge color="blue-6">Kelas {{ props.row.level || props.row.grade }}</q-badge>
          </q-td>
        </template>

        <!-- Custom Render: Aksi (Edit & Hapus) -->
        <template v-slot:body-cell-actions="props">
          <q-td :props="props" class="q-gutter-x-sm">
            <q-btn flat round dense color="orange" icon="edit" @click="openDialog(props.row)">
              <q-tooltip>Edit Kelas</q-tooltip>
            </q-btn>
            <q-btn flat round dense color="negative" icon="delete" @click="confirmDelete(props.row)">
              <q-tooltip>Hapus Kelas</q-tooltip>
            </q-btn>
          </q-td>
        </template>

        <!-- Empty State -->
        <template v-slot:no-data>
          <div class="full-width row flex-center text-grey-6 q-gutter-y-sm q-pa-lg">
            <q-icon size="48px" name="domain_disabled" />
            <span class="text-subtitle1">Belum ada data kelas di semester ini.</span>
          </div>
        </template>
      </q-table>
    </q-card>

    <!-- Dialog Form Tambah/Edit -->
    <q-dialog v-model="dialogOpen" persistent>
      <q-card style="min-width: 450px">
        <q-card-section class="row items-center q-pb-none bg-primary text-white">
          <div class="text-h6">{{ isEdit ? 'Edit Data Kelas' : 'Tambah Kelas Baru' }}</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section>
          <q-form @submit="onSubmit" class="q-gutter-md">
            <q-input v-model="form.name" label="Nama Kelas" hint="Contoh: X TKJ 1" filled
              :rules="[val => !!val || 'Nama kelas wajib diisi']" />
            <div class="row q-col-gutter-sm">
              <div class="col-6">
                <q-input v-model.number="form.grade" label="Tingkat" type="number" hint="10, 11, 12" filled
                  :rules="[val => !!val || 'Wajib diisi']" />
              </div>
              <div class="col-6">
                <q-input v-model.number="form.capacity" label="Kapasitas" type="number" filled
                  :rules="[val => val > 0 || 'Minimal 1']" />
              </div>
            </div>
            <q-input v-model="form.major" label="Jurusan (Opsional)" hint="Contoh: Teknik Komputer dan Jaringan"
              filled />

            <div class="row justify-end q-mt-lg">
              <q-btn label="Batal" color="grey" flat v-close-popup class="q-mr-sm" />
              <q-btn label="Simpan Data" color="primary" type="submit" :loading="kelasStore.isLoading" />
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
import { useKelasStore } from '@/stores/kesiswaan/kelasStore'

const $q = useQuasar()
const kelasStore = useKelasStore()

const filterKeyword = ref('')

// Definisi Kolom Tabel
const columns = [
  { name: 'name', label: 'Nama Kelas', field: 'name', align: 'left', sortable: true },
  { name: 'level', label: 'Tingkat', field: 'grade', align: 'center' },
  { name: 'major', label: 'Jurusan', field: 'major', align: 'left' },
  { name: 'capacity', label: 'Kapasitas', field: 'capacity', align: 'center' },
  { name: 'actions', label: 'Aksi', align: 'center' }
]

// State Dialog & Form
const dialogOpen = ref(false)
const isEdit = ref(false)
const form = ref({
  id: '',
  name: '',
  grade: 10,
  major: '',
  capacity: 36
})

onMounted(() => {
  loadData()
})

// Fungsi Load Data
const loadData = async (keyword = '') => {
  try {
    await kelasStore.fetchKelas({ keyword, limit: 100 })
  } catch (err) {
    $q.notify({ color: 'negative', message: `Gagal memuat data kelas dari server ${err}` })
  }
}

// Pencarian (Bisa disambungkan ke backend via parameter keyword)
const searchData = () => {
  loadData(filterKeyword.value)
}

// Buka Dialog (Tambah / Edit)
const openDialog = (row = null) => {
  if (row) {
    isEdit.value = true
    form.value = { ...row }
  } else {
    isEdit.value = false
    form.value = { id: '', name: '', grade: 10, major: '', capacity: 36 }
  }
  dialogOpen.value = true
}

// Submit Form
const onSubmit = async () => {
  try {
    const payload = {
      name: form.value.name,
      grade: form.value.grade,
      level: form.value.level || '',
      major: form.value.major,
      capacity: form.value.capacity,
      homeroomId: ""
    }

    if (isEdit.value) {
      await kelasStore.updateKelas(form.value.id, payload)
      $q.notify({ color: 'positive', message: 'Data kelas berhasil diperbarui', icon: 'check_circle' })
    } else {
      await kelasStore.createKelas(payload)
      $q.notify({ color: 'positive', message: 'Kelas baru berhasil ditambahkan', icon: 'check_circle' })
    }
    dialogOpen.value = false
  } catch (err) {
    const errorMsg = err.response?.data?.error || 'Terjadi kesalahan sistem'
    $q.notify({ color: 'negative', message: errorMsg })
  }
}

// Konfirmasi Hapus
const confirmDelete = (row) => {
  $q.dialog({
    title: 'Konfirmasi Hapus',
    message: `Apakah Anda yakin ingin menghapus kelas ${row.name}? Tindakan ini tidak dapat dibatalkan.`,
    cancel: true,
    persistent: true,
    color: 'negative'
  }).onOk(async () => {
    try {
      await kelasStore.deleteKelas(row.id)
      $q.notify({ color: 'positive', message: 'Kelas berhasil dihapus', icon: 'delete' })
    } catch (err) {
      $q.notify({ color: 'negative', message: `Gagal menghapus kelas ${err}` })
    }
  })
}
</script>

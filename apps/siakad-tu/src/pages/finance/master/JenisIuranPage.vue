<template>
  <q-page class="q-pa-md">
    <q-linear-progress v-if="store.jenisIuran.loading" indeterminate color="teal" class="q-mb-sm" />

    <div class="row items-center justify-between q-mb-md">
      <div class="text-h5">
        <q-icon name="category" class="q-mr-sm" />Jenis Iuran
      </div>
      <q-btn color="teal" icon="add" label="Tambah" @click="openCreate" />
    </div>

    <q-banner v-if="store.jenisIuran.error" class="bg-negative text-white q-mb-md" rounded>
      {{ store.jenisIuran.error }}
    </q-banner>

    <q-card>
      <q-table :rows="store.jenisIuran.items" :columns="columns" row-key="id" :pagination="pagination"
        :loading="store.jenisIuran.loading" flat bordered @request="onRequest">
        <template v-slot:body-cell-kategori="props">
          <q-td :props="props">
            <q-chip dense :color="chipColor(props.row.kategori)" text-color="white">
              {{ props.row.kategori }}
            </q-chip>
          </q-td>
        </template>
        <template v-slot:body-cell-is_aktif="props">
          <q-td :props="props">
            <q-icon :name="props.row.is_aktif ? 'check_circle' : 'cancel'"
              :color="props.row.is_aktif ? 'positive' : 'grey'" />
          </q-td>
        </template>
        <template v-slot:body-cell-actions="props">
          <q-td :props="props">
            <q-btn flat dense round icon="edit" @click="openEdit(props.row)" />
            <q-btn flat dense round icon="delete" color="negative" @click="confirmDelete(props.row)" />
          </q-td>
        </template>
      </q-table>
    </q-card>

    <q-dialog v-model="dialog" persistent>
      <q-card style="width: 500px; max-width: 90vw">
        <q-card-section>
          <div class="text-h6">{{ editing.id ? 'Edit' : 'Tambah' }} Jenis Iuran</div>
        </q-card-section>
        <q-card-section class="q-gutter-sm">
          <q-input v-model="form.kode" label="Kode *" outlined dense :rules="[v => !!v || 'Wajib diisi']" />
          <q-input v-model="form.nama" label="Nama *" outlined dense :rules="[v => !!v || 'Wajib diisi']" />
          <q-select v-model="form.kategori" :options="kategoriOptions" label="Kategori *" outlined dense emit-value
            map-options :rules="[v => !!v || 'Wajib diisi']" />
          <div class="row q-gutter-sm">
            <q-toggle v-model="form.is_wajib" label="Wajib" />
            <q-toggle v-model="form.is_aktif" label="Aktif" />
          </div>
          <q-input v-model="form.keterangan" label="Keterangan" outlined dense type="textarea" rows="2" />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn color="teal" :label="editing.id ? 'Simpan' : 'Tambah'" :loading="store.jenisIuran.loading"
            @click="submit" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useMasterStore } from '@/stores/finance/masterStore'

const $q = useQuasar()
const store = useMasterStore()

const dialog = ref(false)
const editing = ref({})
const form = reactive({ kode: '', nama: '', kategori: '', is_wajib: true, is_aktif: true, keterangan: '' })
const pagination = ref({ page: 1, rowsPerPage: 50, rowsNumber: 0 })

const kategoriOptions = [
  { label: 'SPP', value: 'SPP' },
  { label: 'Uang Gedung', value: 'GEDUNG' },
  { label: 'Ekstrakurikuler', value: 'EKSTRA' },
  { label: 'Seragam', value: 'SERAGAM' },
  { label: 'Kegiatan', value: 'KEGIATAN' },
  { label: 'Lainnya', value: 'LAINNYA' },
]

const columns = [
  { name: 'kode', label: 'Kode', field: 'kode', align: 'left', sortable: true },
  { name: 'nama', label: 'Nama', field: 'nama', align: 'left', sortable: true },
  { name: 'kategori', label: 'Kategori', field: 'kategori', align: 'center' },
  {
    name: 'is_wajib', label: 'Wajib', field: 'is_wajib', align: 'center',
    format: v => v ? 'Ya' : 'Tidak'
  },
  { name: 'is_aktif', label: 'Aktif', field: 'is_aktif', align: 'center' },
  { name: 'actions', label: 'Aksi', align: 'center' },
]

function chipColor(k) {
  return { SPP: 'blue-6', GEDUNG: 'green-6', EKSTRA: 'orange-6', SERAGAM: 'purple-6', KEGIATAN: 'teal-6' }[k] || 'grey-6'
}

function resetForm() {
  form.kode = ''
  form.nama = ''
  form.kategori = ''
  form.is_wajib = true
  form.is_aktif = true
  form.keterangan = ''
}

function openCreate() {
  editing.value = {}
  resetForm()
  dialog.value = true
}

function openEdit(row) {
  editing.value = row
  form.kode = row.kode
  form.nama = row.nama
  form.kategori = row.kategori
  form.is_wajib = row.is_wajib
  form.is_aktif = row.is_aktif
  form.keterangan = row.keterangan || ''
  dialog.value = true
}

async function submit() {
  if (!form.kode || !form.nama || !form.kategori) {
    $q.notify({ type: 'warning', message: 'Lengkapi field wajib' })
    return
  }
  try {
    if (editing.value.id) {
      await store.jenisIuran.update(editing.value.id, { ...form })
      $q.notify({ type: 'positive', message: 'Berhasil diperbarui' })
    } else {
      await store.jenisIuran.create({ ...form })
      $q.notify({ type: 'positive', message: 'Berhasil ditambahkan' })
    }
    dialog.value = false
    await loadData()
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || e.message })
  }
}

async function confirmDelete(row) {
  $q.dialog({
    title: 'Konfirmasi Hapus',
    message: `Hapus jenis iuran ${row.nama}?`,
    cancel: true,
    persistent: true,
  }).onOk(async () => {
    try {
      await store.jenisIuran.remove(row.id)
      $q.notify({ type: 'positive', message: 'Berhasil dihapus' })
      await loadData()
    } catch (e) {
      $q.notify({ type: 'negative', message: e.response?.data?.error || e.message })
    }
  })
}

async function loadData() {
  try {
    const res = await store.jenisIuran.list({
      limit: pagination.value.rowsPerPage,
      offset: (pagination.value.page - 1) * pagination.value.rowsPerPage,
    })
    pagination.value.rowsNumber = res.pagination?.total || 0
  } catch {
    console.log('JenisIuranPage...')
  }
}

function onRequest(props) {
  pagination.value.page = props.pagination.page
  pagination.value.rowsPerPage = props.pagination.rowsPerPage
  loadData()
}

onMounted(loadData)
</script>

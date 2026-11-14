<template>
  <q-page class="q-pa-md">
    <q-linear-progress v-if="store.kodeRekening.loading" indeterminate color="teal" class="q-mb-sm" />

    <div class="row items-center justify-between q-mb-md">
      <div class="text-h5"><q-icon name="account_tree" class="q-mr-sm" />Kode Rekening (COA)</div>
      <div class="row q-gutter-sm">
        <q-select v-model="filterTipe" :options="tipeFilter" label="Filter Tipe" outlined dense clearable
          style="width: 180px" @update:model-value="loadData" />
        <q-btn color="teal" icon="add" label="Tambah" @click="openCreate" />
      </div>
    </div>

    <q-card>
      <q-table :rows="store.kodeRekening.items" :columns="columns" row-key="id" flat bordered :pagination="pagination"
        :loading="store.kodeRekening.loading" @request="onRequest" :filter="search" :filter-method="serverFilter">
        <template v-slot:top-right>
          <q-input v-model="search" placeholder="Cari kode/nama..." outlined dense @update:model-value="loadData">
            <template v-slot:append><q-icon name="search" /></template>
          </q-input>
        </template>
        <template v-slot:body-cell-tipe="props">
          <q-td :props="props">
            <q-chip dense :color="chipColor(props.row.tipe)" text-color="white">
              {{ props.row.tipe }}
            </q-chip>
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
          <div class="text-h6">{{ editing.id ? 'Edit' : 'Tambah' }} Kode Rekening</div>
        </q-card-section>
        <q-card-section class="q-gutter-sm">
          <q-input v-model="form.kode" label="Kode *" outlined dense hint="Contoh: 1.1.1" />
          <q-input v-model="form.nama" label="Nama *" outlined dense />
          <q-select v-model="form.tipe" :options="tipeOptions" label="Tipe *" outlined dense emit-value map-options />
          <q-input v-model.number="form.level" label="Level" outlined dense type="number" />
          <q-toggle v-model="form.is_aktif" label="Aktif" />
          <q-input v-model="form.keterangan" label="Keterangan" outlined dense />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn color="teal" :label="editing.id ? 'Simpan' : 'Tambah'" :loading="store.kodeRekening.loading"
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
const search = ref('')
const filterTipe = ref(null)
const form = reactive({ kode: '', nama: '', tipe: '', parent_id: '', level: 1, is_aktif: true, keterangan: '' })
const pagination = ref({ page: 1, rowsPerPage: 100, rowsNumber: 0 })

const tipeOptions = [
  { label: 'Aset', value: 'ASET' },
  { label: 'Kewajiban', value: 'KEWAJIBAN' },
  { label: 'Ekuitas', value: 'EKUITAS' },
  { label: 'Pendapatan', value: 'PENDAPATAN' },
  { label: 'Belanja', value: 'BELANJA' },
]
const tipeFilter = [
  { label: 'Semua', value: null },
  ...tipeOptions,
]

const columns = [
  { name: 'kode', label: 'Kode', field: 'kode', align: 'left', sortable: true },
  { name: 'nama', label: 'Nama', field: 'nama', align: 'left', sortable: true },
  { name: 'tipe', label: 'Tipe', field: 'tipe', align: 'center' },
  { name: 'level', label: 'Level', field: 'level', align: 'center' },
  {
    name: 'is_aktif', label: 'Aktif', field: 'is_aktif', align: 'center',
    format: v => v ? '✓' : '✗'
  },
  { name: 'actions', label: 'Aksi', align: 'center' },
]

function chipColor(t) {
  return { ASET: 'blue-6', KEWAJIBAN: 'red-6', EKUITAS: 'purple-6', PENDAPATAN: 'green-6', BELANJA: 'orange-6' }[t] || 'grey-6'
}
function serverFilter() { return true } // server-side already filtered

function resetForm() {
  form.kode = ''
  form.nama = ''
  form.tipe = ''
  form.parent_id = ''
  form.level = 1
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
  Object.assign(form, row)
  form.keterangan = row.keterangan || ''
  dialog.value = true
}

async function submit() {
  if (!form.kode || !form.nama || !form.tipe) {
    $q.notify({ type: 'warning', message: 'Lengkapi field wajib' })
    return
  }
  try {
    if (editing.value.id) {
      await store.kodeRekening.update(editing.value.id, { ...form })
    } else {
      await store.kodeRekening.create({ ...form })
    }
    $q.notify({ type: 'positive', message: 'Berhasil disimpan' })
    dialog.value = false
    await loadData()
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || e.message })
  }
}

async function confirmDelete(row) {
  $q.dialog({ title: 'Konfirmasi Hapus', message: `Hapus COA ${row.kode} ${row.nama}?`, cancel: true })
    .onOk(async () => {
      try {
        await store.kodeRekening.remove(row.id)
        $q.notify({ type: 'positive', message: 'Berhasil dihapus' })
        await loadData()
      } catch (e) {
        $q.notify({ type: 'negative', message: e.response?.data?.error || e.message })
      }
    })
}

async function loadData() {
  try {
    const params = { limit: pagination.value.rowsPerPage }
    if (filterTipe.value) params.tipe = filterTipe.value
    if (search.value) params.search = search.value
    const res = await store.kodeRekening.list(params)
    pagination.value.rowsNumber = res.pagination?.total || 0
  } catch {
    console.log('KodeRekeningPage...')
  }
}

function onRequest(props) {
  pagination.value.page = props.pagination.page
  loadData()
}

onMounted(loadData)
</script>

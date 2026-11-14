<template>
  <q-page class="q-pa-md">
    <q-linear-progress v-if="store.rekeningBank.loading" indeterminate color="teal" class="q-mb-sm" />

    <div class="row items-center justify-between q-mb-md">
      <div class="text-h5"><q-icon name="account_balance" class="q-mr-sm" />Rekening Bank</div>
      <q-btn color="teal" icon="add" label="Tambah" @click="openCreate" />
    </div>

    <q-card>
      <q-table :rows="store.rekeningBank.items" :columns="columns" row-key="id" flat bordered :pagination="pagination"
        :loading="store.rekeningBank.loading" @request="onRequest">
        <template v-slot:body-cell-saldo_terkini="props">
          <q-td :props="props">
            <span :class="props.row.saldo_terkini < 0 ? 'text-negative' : ''">
              {{ formatRp(props.row.saldo_terkini) }}
            </span>
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
          <div class="text-h6">{{ editing.id ? 'Edit' : 'Tambah' }} Rekening Bank</div>
        </q-card-section>
        <q-card-section class="q-gutter-sm">
          <q-input v-model="form.nama_bank" label="Nama Bank *" outlined dense />
          <q-input v-model="form.nomor_rekening" label="Nomor Rekening *" outlined dense />
          <q-input v-model="form.atas_nama" label="Atas Nama *" outlined dense />
          <q-input v-model="form.cabang" label="Cabang" outlined dense />
          <q-input v-model.number="form.saldo_awal" label="Saldo Awal" outlined dense type="number" />
          <q-select v-model="form.format_file_mutasi" :options="formatOptions" label="Format File Mutasi" outlined dense
            emit-value map-options clearable />
          <q-toggle v-model="form.is_aktif" label="Aktif" />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn color="teal" :label="editing.id ? 'Simpan' : 'Tambah'" :loading="store.rekeningBank.loading"
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
const form = reactive({
  nama_bank: '', nomor_rekening: '', atas_nama: '', cabang: '',
  saldo_awal: 0, saldo_terkini: 0, format_file_mutasi: '', is_aktif: true,
})
const pagination = ref({ page: 1, rowsPerPage: 50, rowsNumber: 0 })

const formatOptions = [
  { label: 'CSV BCA', value: 'CSV_BCA' },
  { label: 'CSV Mandiri', value: 'CSV_MANDIRI' },
  { label: 'CSV BRI', value: 'CSV_BRI' },
  { label: 'Excel Generic', value: 'XLSX' },
]

const columns = [
  { name: 'nama_bank', label: 'Bank', field: 'nama_bank', align: 'left', sortable: true },
  { name: 'nomor_rekening', label: 'No. Rekening', field: 'nomor_rekening', align: 'left' },
  { name: 'atas_nama', label: 'Atas Nama', field: 'atas_nama', align: 'left' },
  { name: 'saldo_terkini', label: 'Saldo', field: 'saldo_terkini', align: 'right', sortable: true },
  { name: 'is_aktif', label: 'Aktif', field: 'is_aktif', align: 'center' },
  { name: 'actions', label: 'Aksi', align: 'center' },
]

function formatRp(n) {
  return n === null || n === undefined ? '-' : 'Rp ' + Number(n).toLocaleString('id-ID')
}

function resetForm() {
  form.nama_bank = ''
  form.nomor_rekening = ''
  form.atas_nama = ''
  form.cabang = ''
  form.saldo_awal = 0
  form.saldo_terkini = 0
  form.format_file_mutasi = ''
  form.is_aktif = true
}

function openCreate() {
  editing.value = {}
  resetForm()
  dialog.value = true
}

function openEdit(row) {
  editing.value = row
  Object.assign(form, row)
  form.saldo_awal = row.saldo_awal || 0
  form.saldo_terkini = row.saldo_terkini || 0
  dialog.value = true
}

async function submit() {
  if (!form.nama_bank || !form.nomor_rekening || !form.atas_nama) {
    $q.notify({ type: 'warning', message: 'Lengkapi field wajib' })
    return
  }
  try {
    if (editing.value.id) {
      await store.rekeningBank.update(editing.value.id, { ...form })
    } else {
      await store.rekeningBank.create({ ...form })
    }
    $q.notify({ type: 'positive', message: 'Berhasil disimpan' })
    dialog.value = false
    await loadData()
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || e.message })
  }
}

async function confirmDelete(row) {
  $q.dialog({ title: 'Konfirmasi Hapus', message: `Hapus rekening ${row.nama_bank}?`, cancel: true })
    .onOk(async () => {
      try {
        await store.rekeningBank.remove(row.id)
        $q.notify({ type: 'positive', message: 'Berhasil dihapus' })
        await loadData()
      } catch (e) {
        $q.notify({ type: 'negative', message: e.response?.data?.error || e.message })
      }
    })
}

async function loadData() {
  try {
    const res = await store.rekeningBank.list({ limit: pagination.value.rowsPerPage })
    pagination.value.rowsNumber = res.pagination?.total || 0
  } catch {
    console.log('RekeningBankPage....')
  }
}

function onRequest(props) {
  pagination.value.page = props.pagination.page
  loadData()
}

onMounted(loadData)
</script>

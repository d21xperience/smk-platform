<template>
  <q-page class="q-pa-md">
    <q-linear-progress v-if="store.tarifIuran.loading" indeterminate color="teal" class="q-mb-sm" />

    <div class="row items-center justify-between q-mb-md">
      <div class="text-h5"><q-icon name="price_change" class="q-mr-sm" />Tarif Iuran</div>
      <q-btn color="teal" icon="add" label="Tambah" @click="openCreate" />
    </div>

    <q-card>
      <q-table :rows="store.tarifIuran.items" :columns="columns" row-key="id" flat bordered :pagination="pagination"
        :loading="store.tarifIuran.loading" @request="onRequest">
        <template v-slot:body-cell-nominal="props">
          <q-td :props="props">
            <span class="text-weight-medium">{{ formatRp(props.row.nominal) }}</span>
          </q-td>
        </template>
        <template v-slot:body-cell-kelas="props">
          <q-td :props="props">
            {{ props.row.kelas || '(semua kelas)' }}
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
          <div class="text-h6">{{ editing.id ? 'Edit' : 'Tambah' }} Tarif Iuran</div>
        </q-card-section>
        <q-card-section class="q-gutter-sm">
          <q-select v-model="form.jenis_iuran_id" :options="jenisOptions" label="Jenis Iuran *" outlined dense
            emit-value map-options :loading="store.jenisIuran.loading" :rules="[v => !!v || 'Wajib diisi']" />
          <q-input v-model="form.kelas" label="Kelas (kosongkan untuk semua kelas)" outlined dense />
          <q-input v-model.number="form.nominal" label="Nominal *" outlined dense type="number" prefix="Rp"
            :rules="[v => v > 0 || 'Harus > 0']" />
          <q-input v-model.number="form.tahun_berlaku" label="Tahun Berlaku *" outlined dense type="number"
            :rules="[v => v >= 2000 && v <= 2100 || 'Tahun tidak valid']" />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn color="teal" :label="editing.id ? 'Simpan' : 'Tambah'" :loading="store.tarifIuran.loading"
            @click="submit" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useMasterStore } from '@/stores/finance/masterStore'

const $q = useQuasar()
const store = useMasterStore()

const dialog = ref(false)
const editing = ref({})
const form = reactive({ jenis_iuran_id: '', kelas: '', nominal: 0, tahun_berlaku: new Date().getFullYear() })
const pagination = ref({ page: 1, rowsPerPage: 50, rowsNumber: 0 })

const jenisOptions = computed(() =>
  store.jenisIuran.items.map(j => ({ label: `${j.kode} — ${j.nama}`, value: j.id }))
)

const columns = [
  { name: 'jenis_iuran_nama', label: 'Jenis Iuran', field: 'jenis_iuran_nama', align: 'left', sortable: true },
  { name: 'kelas', label: 'Kelas', field: 'kelas', align: 'left' },
  { name: 'nominal', label: 'Nominal', field: 'nominal', align: 'right', sortable: true },
  { name: 'tahun_berlaku', label: 'Tahun', field: 'tahun_berlaku', align: 'center', sortable: true },
  { name: 'actions', label: 'Aksi', align: 'center' },
]

function formatRp(n) {
  return n === null || n === undefined ? '-' : 'Rp ' + Number(n).toLocaleString('id-ID')
}

function resetForm() {
  form.jenis_iuran_id = ''
  form.kelas = ''
  form.nominal = 0
  form.tahun_berlaku = new Date().getFullYear()
}

function openCreate() {
  editing.value = {}
  resetForm()
  dialog.value = true
}
function openEdit(row) {
  editing.value = row
  Object.assign(form, row)
  form.kelas = row.kelas || ''
  dialog.value = true
}

async function submit() {
  if (!form.jenis_iuran_id || !form.nominal) {
    $q.notify({ type: 'warning', message: 'Lengkapi field wajib' })
    return
  }
  try {
    if (editing.value.id) {
      await store.tarifIuran.update(editing.value.id, { ...form })
    } else {
      await store.tarifIuran.create({ ...form })
    }
    $q.notify({ type: 'positive', message: 'Berhasil disimpan' })
    dialog.value = false
    await loadData()
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || e.message })
  }
}

async function confirmDelete(row) {
  $q.dialog({ title: 'Konfirmasi Hapus', message: `Hapus tarif ${row.jenis_iuran_nama}?`, cancel: true })
    .onOk(async () => {
      try {
        await store.tarifIuran.remove(row.id)
        $q.notify({ type: 'positive', message: 'Berhasil dihapus' })
        await loadData()
      } catch (e) {
        $q.notify({ type: 'negative', message: e.response?.data?.error || e.message })
      }
    })
}

async function loadData() {
  try {
    // Load jenis iuran untuk dropdown
    if (store.jenisIuran.items.length === 0) {
      await store.jenisIuran.list({ limit: 100 })
    }
    const res = await store.tarifIuran.list({ limit: pagination.value.rowsPerPage })
    pagination.value.rowsNumber = res.pagination?.total || 0
  } catch {
    console.log('TarifIuranPage...')
  }
}

function onRequest(props) {
  pagination.value.page = props.pagination.page
  loadData()
}

onMounted(loadData)
</script>

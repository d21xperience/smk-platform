<template>
  <q-page class="q-pa-md">
    <q-linear-progress v-if="store.tagihanLoading" indeterminate color="teal" class="q-mb-sm" />

    <div class="row items-center justify-between q-mb-md">
      <div class="text-h5"><q-icon name="list_alt" class="q-mr-sm" />Daftar Tagihan</div>
    </div>

    <q-card class="q-mb-md">
      <q-card-section>
        <div class="row q-col-gutter-sm">
          <div class="col-12 col-sm-4 col-md-3">
            <q-input v-model="filter.nis" label="NIS" outlined dense clearable @keyup.enter="loadData" />
          </div>
          <div class="col-12 col-sm-4 col-md-3">
            <q-select v-model="filter.status" :options="statusOptions" label="Status" outlined dense emit-value
              map-options clearable />
          </div>
          <div class="col-12 col-sm-4 col-md-3">
            <q-select v-model="filter.jenis_iuran_id" :options="jenisOptions" label="Jenis Iuran" outlined dense
              emit-value map-options clearable />
          </div>
          <div class="col-12 col-sm-4 col-md-2">
            <q-input v-model.number="filter.bulan" label="Bulan" outlined dense type="number" min="1" max="12" />
          </div>
          <div class="col-12 col-sm-4 col-md-1">
            <q-btn color="teal" icon="search" @click="loadData" class="full-width" />
          </div>
        </div>
      </q-card-section>
    </q-card>

    <q-card>
      <q-table :rows="store.tagihan" :columns="columns" row-key="id" flat bordered :pagination="pagination"
        :loading="store.tagihanLoading" @request="onRequest">
        <template v-slot:body-cell-sisa="props">
          <q-td :props="props">
            <span class="text-weight-medium">{{ formatRp(props.row.nominal - props.row.total_dibayar) }}</span>
          </q-td>
        </template>
        <template v-slot:body-cell-status="props">
          <q-td :props="props">
            <q-chip dense :color="statusColor(props.row.status)" text-color="white">
              {{ props.row.status }}
            </q-chip>
          </q-td>
        </template>
      </q-table>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useTransaksiStore } from '@/stores/finance/transaksiStore'
import { useMasterStore } from '@/stores/finance/masterStore'

const store = useTransaksiStore()
const masterStore = useMasterStore()

const filter = reactive({ nis: '', status: '', jenis_iuran_id: '', bulan: null })
const pagination = ref({ page: 1, rowsPerPage: 50, rowsNumber: 0 })

const statusOptions = [
  { label: 'Belum Lunas', value: 'BELUM_LUNAS' },
  { label: 'Lunas', value: 'LUNAS' },
  { label: 'Sebagian', value: 'SEBAGIAN' },
  { label: 'Dibatalkan', value: 'DIBATALKAN' },
]

const jenisOptions = computed(() =>
  masterStore.jenisIuran.items.map(j => ({ label: j.nama, value: j.id }))
)

const columns = [
  { name: 'nis', label: 'NIS', field: 'nis', align: 'left', sortable: true },
  { name: 'kelas_nama', label: 'Kelas', field: 'kelas_nama', align: 'left' },
  { name: 'jenis_iuran_nama', label: 'Jenis Iuran', field: 'jenis_iuran_nama', align: 'left' },
  {
    name: 'periode', label: 'Periode', align: 'center',
    field: r => `${r.periode_bulan}/${r.periode_tahun}`
  },
  {
    name: 'nominal', label: 'Nominal', field: 'nominal', align: 'right',
    format: v => formatRp(v)
  },
  { name: 'sisa', label: 'Sisa', align: 'right' },
  { name: 'status', label: 'Status', field: 'status', align: 'center' },
  { name: 'tanggal_jatuh_tempo', label: 'Jatuh Tempo', field: 'tanggal_jatuh_tempo', align: 'center' },
]

function formatRp(n) {
  return n === null || n === undefined ? '-' : 'Rp ' + Number(n).toLocaleString('id-ID')
}
function statusColor(s) {
  return { BELUM_LUNAS: 'negative', LUNAS: 'positive', SEBAGIAN: 'warning', DIBATALKAN: 'grey' }[s] || 'grey'
}

async function loadData() {
  try {
    const params = {
      limit: pagination.value.rowsPerPage,
      offset: (pagination.value.page - 1) * pagination.value.rowsPerPage,
    }
    if (filter.nis) params.nis = filter.nis
    if (filter.status) params.status = filter.status
    if (filter.jenis_iuran_id) params.jenis_iuran_id = filter.jenis_iuran_id
    if (filter.bulan) params.bulan = filter.bulan
    const res = await store.listTagihan(params)
    pagination.value.rowsNumber = res.pagination?.total || 0
  } catch {
    console.log('TagihanPage....')
  }
}

function onRequest(props) {
  pagination.value.page = props.pagination.page
  pagination.value.rowsPerPage = props.pagination.rowsPerPage
  loadData()
}

onMounted(async () => {
  if (masterStore.jenisIuran.items.length === 0) {
    await masterStore.jenisIuran.list({ limit: 100 })
  }
  await loadData()
})
</script>

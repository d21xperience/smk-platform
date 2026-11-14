<template>
  <q-page class="q-pa-md">
    <q-linear-progress v-if="store.bkuLoading" indeterminate color="teal" class="q-mb-sm" />

    <div class="row items-center justify-between q-mb-md">
      <div class="text-h5"><q-icon name="menu_book" class="q-mr-sm" />Buku Kas Umum (BKU)</div>
      <q-chip v-if="store.saldoKas !== 0" :color="store.saldoKas < 0 ? 'negative' : 'positive'" text-color="white">
        <q-icon name="payments" left size="16px" />Saldo: {{ formatRp(store.saldoKas) }}
      </q-chip>
    </div>

    <q-card class="q-mb-md">
      <q-card-section>
        <div class="row q-col-gutter-sm items-end">
          <div class="col-12 col-sm-3 col-md-2">
            <q-input v-model="filter.tanggal_dari" label="Dari" outlined dense type="date" />
          </div>
          <div class="col-12 col-sm-3 col-md-2">
            <q-input v-model="filter.tanggal_sampai" label="Sampai" outlined dense type="date" />
          </div>
          <div class="col-12 col-sm-3 col-md-2">
            <q-btn color="teal" icon="search" label="Filter" @click="loadData" :loading="store.bkuLoading"
              class="full-width" />
          </div>
          <div class="col-12 col-sm-3 col-md-2">
            <q-btn flat icon="refresh" label="Reset" @click="resetFilter" class="full-width" />
          </div>
        </div>
      </q-card-section>
    </q-card>

    <q-card>
      <q-table :rows="store.bku" :columns="columns" row-key="id" flat bordered dense :pagination="pagination"
        :loading="store.bkuLoading" @request="onRequest">
        <template v-slot:body-cell-debet="props">
          <q-td :props="props" class="text-right">
            <span v-if="props.row.debet > 0" class="text-positive text-weight-medium">
              +{{ formatRp(props.row.debet) }}
            </span>
          </q-td>
        </template>
        <template v-slot:body-cell-kredit="props">
          <q-td :props="props" class="text-right">
            <span v-if="props.row.kredit > 0" class="text-negative text-weight-medium">
              -{{ formatRp(props.row.kredit) }}
            </span>
          </q-td>
        </template>
        <template v-slot:body-cell-saldo_berjalan="props">
          <q-td :props="props" class="text-right">
            <span :class="props.row.saldo_berjalan < 0 ? 'text-negative' : ''">
              {{ formatRp(props.row.saldo_berjalan) }}
            </span>
          </q-td>
        </template>
        <template v-slot:body-cell-is_locked="props">
          <q-td :props="props" class="text-center">
            <q-icon v-if="props.row.is_locked" name="lock" color="grey-7" size="18px">
              <q-tooltip>Periode sudah ditutup</q-tooltip>
            </q-icon>
            <q-icon v-else name="lock_open" color="positive" size="18px" />
          </q-td>
        </template>
        <template v-slot:body-cell-referensi_tabel="props">
          <q-td :props="props">
            <q-chip v-if="props.row.referensi_tabel" dense :color="refColor(props.row.referensi_tabel)"
              text-color="white">
              {{ props.row.referensi_tabel }}
            </q-chip>
          </q-td>
        </template>
      </q-table>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useTransaksiStore } from '@/stores/finance/transaksiStore'

const store = useTransaksiStore()
const filter = reactive({ tanggal_dari: '', tanggal_sampai: '' })
const pagination = ref({ page: 1, rowsPerPage: 100, rowsNumber: 0 })

const columns = [
  { name: 'tanggal', label: 'Tanggal', field: 'tanggal', align: 'left', sortable: true },
  { name: 'no_bukti', label: 'No. Bukti', field: 'no_bukti', align: 'left' },
  { name: 'uraian', label: 'Uraian', field: 'uraian', align: 'left', style: 'max-width: 300px; white-space: normal' },
  { name: 'kode_rekening_kode', label: 'Kode Akun', field: 'kode_rekening_kode', align: 'center' },
  { name: 'debet', label: 'Debet', align: 'right', sortable: true },
  { name: 'kredit', label: 'Kredit', align: 'right', sortable: true },
  { name: 'saldo_berjalan', label: 'Saldo', align: 'right', sortable: true },
  { name: 'referensi_tabel', label: 'Ref', align: 'center' },
  { name: 'is_locked', label: 'Lock', align: 'center' },
]

function formatRp(n) {
  return n === null || n === undefined ? '-' : 'Rp ' + Number(n).toLocaleString('id-ID')
}
function refColor(ref) {
  return { penerimaan_iuran: 'positive', pengeluaran: 'negative' }[ref] || 'grey'
}

async function loadData() {
  try {
    const params = {
      limit: pagination.value.rowsPerPage,
      offset: (pagination.value.page - 1) * pagination.value.rowsPerPage,
    }
    if (filter.tanggal_dari) params.tanggal_dari = filter.tanggal_dari
    if (filter.tanggal_sampai) params.tanggal_sampai = filter.tanggal_sampai
    const res = await store.listBKU(params)
    pagination.value.rowsNumber = res.pagination?.total || 0
    // Load saldo terkini
    await store.getSaldo({ tanggal: filter.tanggal_sampai || undefined })
  } catch {
    console.log('BkuPage...')
  }
}

function resetFilter() {
  filter.tanggal_dari = ''
  filter.tanggal_sampai = ''
  loadData()
}

function onRequest(props) {
  pagination.value.page = props.pagination.page
  pagination.value.rowsPerPage = props.pagination.rowsPerPage
  loadData()
}

onMounted(loadData)
</script>

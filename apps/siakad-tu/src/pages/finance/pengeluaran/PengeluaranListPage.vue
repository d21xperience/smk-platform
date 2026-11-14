<template>
  <q-page class="q-pa-md">
    <q-linear-progress v-if="store.pengeluaranLoading" indeterminate color="teal" class="q-mb-sm" />

    <div class="row items-center justify-between q-mb-md">
      <div class="text-h5"><q-icon name="payments" class="q-mr-sm" />Daftar Pengeluaran (BKK)</div>
      <q-btn flat icon="refresh" label="Refresh" @click="loadData" />
    </div>

    <q-card class="q-mb-md">
      <q-card-section>
        <div class="row q-col-gutter-sm items-end">
          <div class="col-12 col-sm-3 col-md-2">
            <q-select v-model="filter.status" :options="statusOptions" label="Status" outlined dense emit-value
              map-options clearable @update:model-value="loadData" />
          </div>
          <div class="col-12 col-sm-3 col-md-2">
            <q-input v-model="filter.tanggal_dari" label="Dari" outlined dense type="date" />
          </div>
          <div class="col-12 col-sm-3 col-md-2">
            <q-input v-model="filter.tanggal_sampai" label="Sampai" outlined dense type="date" />
          </div>
          <div class="col-12 col-sm-3 col-md-2">
            <q-btn color="teal" icon="search" label="Filter" @click="loadData" class="full-width" />
          </div>
        </div>
      </q-card-section>
    </q-card>

    <q-card v-if="store.pengeluaran.length === 0 && !store.pengeluaranLoading">
      <q-card-section class="text-center text-grey-6 q-pa-xl">
        <q-icon name="payments" size="48px" class="q-mb-sm" />
        <div>Tidak ada pengeluaran.</div>
      </q-card-section>
    </q-card>

    <q-card v-else>
      <q-table :rows="store.pengeluaran" :columns="columns" row-key="id" flat bordered :pagination="{ rowsPerPage: 50 }"
        :loading="store.pengeluaranLoading">
        <template v-slot:body-cell-jumlah="props">
          <q-td :props="props" class="text-right text-negative text-weight-medium">
            -{{ formatRp(props.row.jumlah) }}
          </q-td>
        </template>
        <template v-slot:body-cell-status="props">
          <q-td :props="props" class="text-center">
            <q-chip dense :color="statusColor(props.row.status)" text-color="white">
              {{ props.row.status }}
            </q-chip>
          </q-td>
        </template>
        <template v-slot:body-cell-no_bkk="props">
          <q-td :props="props">
            <span class="text-weight-medium">{{ props.row.no_bkk }}</span>
          </q-td>
        </template>
        <template v-slot:body-cell-actions="props">
          <q-td :props="props" class="text-center">
            <q-btn flat dense round icon="visibility" color="info" @click="openDetail(props.row)" />
          </q-td>
        </template>
      </q-table>
    </q-card>

    <q-dialog v-model="detailDialog">
      <q-card style="width: 550px; max-width: 90vw">
        <q-card-section>
          <div class="text-h6">Detail Pengeluaran</div>
        </q-card-section>
        <q-card-section v-if="currentRow">
          <q-list dense>
            <q-item><q-item-section>
                <div class="text-caption">No. BKK</div>
                <div class="text-subtitle2 text-weight-bold">{{ currentRow.no_bkk }}</div>
              </q-item-section></q-item>
            <q-item><q-item-section>
                <div class="text-caption">Tanggal</div>
                <div>{{ currentRow.tanggal }}</div>
              </q-item-section></q-item>
            <q-item><q-item-section>
                <div class="text-caption">Keperluan</div>
                <div>{{ currentRow.keperluan }}</div>
              </q-item-section></q-item>
            <q-item><q-item-section>
                <div class="text-caption">Jumlah</div>
                <div class="text-h6 text-negative">{{ formatRp(currentRow.jumlah) }}</div>
              </q-item-section></q-item>
            <q-item><q-item-section>
                <div class="text-caption">Status</div>
                <q-chip dense :color="statusColor(currentRow.status)" text-color="white">
                  {{ currentRow.status }}</q-chip>
              </q-item-section></q-item>
          </q-list>
        </q-card-section>
        <q-card-actions align="right"><q-btn flat label="Tutup" v-close-popup /></q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { usePengeluaranStore } from '@/stores/finance/pengeluaranStore'

const store = usePengeluaranStore()

const filter = reactive({ status: '', tanggal_dari: '', tanggal_sampai: '' })
const detailDialog = ref(false)
const currentRow = ref(null)

const statusOptions = [
  { label: 'Siap Cair', value: 'SIAP_CAIR' },
  { label: 'Cair', value: 'CAIR' },
  { label: 'Dibatalkan', value: 'DIBATALKAN' },
]

const columns = [
  { name: 'no_bkk', label: 'No. BKK', align: 'left', sortable: true },
  { name: 'tanggal', label: 'Tanggal', field: 'tanggal', align: 'left', sortable: true },
  {
    name: 'keperluan', label: 'Keperluan', field: 'keperluan', align: 'left',
    style: 'max-width: 350px; white-space: normal'
  },
  { name: 'jumlah', label: 'Jumlah', align: 'right', sortable: true },
  { name: 'status', label: 'Status', field: 'status', align: 'center' },
  { name: 'actions', label: 'Detail', align: 'center' },
]

function formatRp(n) {
  return n === null || n === undefined ? '-' : 'Rp ' + Number(n).toLocaleString('id-ID')
}
function statusColor(s) {
  return { SIAP_CAIR: 'orange', CAIR: 'positive', DIBATALKAN: 'grey' }[s] || 'grey'
}

function openDetail(row) { currentRow.value = row; detailDialog.value = true }

async function loadData() {
  try {
    const params = { limit: 100 }
    if (filter.status) params.status = filter.status
    if (filter.tanggal_dari) params.tanggal_dari = filter.tanggal_dari
    if (filter.tanggal_sampai) params.tanggal_sampai = filter.tanggal_sampai
    await store.listPengeluaran(params)
  } catch {
    console.log('PengeluaranListPage...')
  }
}

onMounted(loadData)
</script>

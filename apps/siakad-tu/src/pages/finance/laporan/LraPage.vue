<template>
  <q-page class="q-pa-md">
    <q-linear-progress v-if="store.lraLoading" indeterminate color="teal" class="q-mb-sm" />

    <div class="row items-center justify-between q-mb-md">
      <div class="text-h5"><q-icon name="summarize" class="q-mr-sm" />LRA (Laporan Realisasi Anggaran)</div>
    </div>

    <q-card class="q-mb-md">
      <q-card-section>
        <div class="row q-col-gutter-sm items-end">
          <div class="col-6 col-md-2">
            <q-input v-model.number="tahun" label="Tahun Anggaran" outlined dense type="number" />
          </div>
          <div class="col-6 col-md-2">
            <q-select v-model="sampaiBulan" :options="bulanOptions" label="Sampai Bulan" outlined dense
              emit-value map-options />
          </div>
          <div class="col-12 col-md-2">
            <q-btn color="teal" icon="refresh" label="Muat" @click="loadData"
              :loading="store.lraLoading" class="full-width" />
          </div>
        </div>
      </q-card-section>
    </q-card>

    <div v-if="lra">
      <q-card class="q-mb-md">
        <q-card-section>
          <div class="row q-col-gutter-md items-center">
            <div class="col-6 col-md-3">
              <div class="text-caption text-grey-7">RAB</div>
              <div class="text-subtitle1">{{ lra.rab_id.substring(0, 8) }}...</div>
            </div>
            <div class="col-6 col-md-3">
              <div class="text-caption text-grey-7">Status</div>
              <q-chip dense :color="lra.rab_status === 'DISAHKAN' ? 'positive' : 'grey'" text-color="white">
                {{ lra.rab_status }}
              </q-chip>
            </div>
            <div class="col-6 col-md-3">
              <div class="text-caption text-grey-7">Tanggal Pengesahan</div>
              <div>{{ lra.tanggal_pengesahan || '-' }}</div>
            </div>
            <div class="col-6 col-md-3">
              <div class="text-caption text-grey-7">Sampai Bulan</div>
              <div>{{ lra.sampai_bulan }}</div>
            </div>
          </div>
        </q-card-section>
      </q-card>

      <q-card v-if="lra.realisasi">
        <q-card-section>
          <div class="row q-col-gutter-md q-mb-md">
            <div class="col-12 col-md-6">
              <q-card class="bg-green-1" flat>
                <q-card-section>
                  <div class="text-caption">Pendapatan</div>
                  <div class="text-caption">
                    Pagu: {{ formatRp(lra.realisasi.total_pagu_pendapatan) }} ·
                    Realisasi: {{ formatRp(lra.realisasi.total_realisasi_pendapatan) }}
                  </div>
                  <q-linear-progress :value="ratio(lra.realisasi.total_realisasi_pendapatan, lra.realisasi.total_pagu_pendapatan)"
                    color="green-8" class="q-mt-sm" size="8px" />
                </q-card-section>
              </q-card>
            </div>
            <div class="col-12 col-md-6">
              <q-card class="bg-orange-1" flat>
                <q-card-section>
                  <div class="text-caption">Belanja</div>
                  <div class="text-caption">
                    Pagu: {{ formatRp(lra.realisasi.total_pagu_belanja) }} ·
                    Realisasi: {{ formatRp(lra.realisasi.total_realisasi_belanja) }}
                  </div>
                  <q-linear-progress :value="ratio(lra.realisasi.total_realisasi_belanja, lra.realisasi.total_pagu_belanja)"
                    color="orange-8" class="q-mt-sm" size="8px" />
                </q-card-section>
              </q-card>
            </div>
          </div>

          <q-table :rows="lra.realisasi.items" :columns="columns" row-key="rab_detail_id" flat bordered dense
            :pagination="{ rowsPerPage: 50 }" :filter="filter" :filter-method="serverFilter">
            <template v-slot:top-right>
              <q-input v-model="filter" placeholder="Cari kode/uraian..." outlined dense>
                <template v-slot:append><q-icon name="search" /></template>
              </q-input>
            </template>
            <template v-slot:body-cell-pagu="props">
              <q-td :props="props" class="text-right">{{ formatRp(props.row.pagu) }}</q-td>
            </template>
            <template v-slot:body-cell-realisasi="props">
              <q-td :props="props" class="text-right">{{ formatRp(props.row.realisasi) }}</q-td>
            </template>
            <template v-slot:body-cell-sisa="props">
              <q-td :props="props" class="text-right">{{ formatRp(props.row.sisa) }}</q-td>
            </template>
            <template v-slot:body-cell-persen_realisasi="props">
              <q-td :props="props" class="text-right">{{ props.row.persen_realisasi.toFixed(1) }}%</q-td>
            </template>
            <template v-slot:body-cell-status_warning="props">
              <q-td :props="props" class="text-center">
                <q-chip v-if="props.row.status_warning !== 'NORMAL'" dense
                  :color="props.row.status_warning === 'OVER_BUDGET_WARNING' ? 'negative' : 'warning'"
                  text-color="white">
                  {{ props.row.status_warning === 'OVER_BUDGET_WARNING' ? 'Over' : 'Under' }}
                </q-chip>
                <q-chip v-else dense color="positive" text-color="white">OK</q-chip>
              </q-td>
            </template>
          </q-table>
        </q-card-section>
      </q-card>
    </div>

    <q-card v-else-if="!store.lraLoading && !errorMsg">
      <q-card-section class="text-center text-grey-6 q-pa-xl">
        <q-icon name="summarize" size="48px" class="q-mb-sm" />
        <div>Pilih tahun dan klik Muat untuk melihat LRA.</div>
      </q-card-section>
    </q-card>

    <q-card v-if="errorMsg" class="bg-negative text-white">
      <q-card-section>{{ errorMsg }}</q-card-section>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useLaporanStore } from '@/stores/finance/laporanStore'

const store = useLaporanStore()

const tahun = ref(new Date().getFullYear())
const sampaiBulan = ref(new Date().getMonth() + 1)
const filter = ref('')
const errorMsg = ref('')

const bulanOptions = [
  'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
  'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
].map((label, i) => ({ label, value: i + 1 }))

const lra = computed(() => store.lra)

const columns = [
  { name: 'kode_rekening_kode', label: 'Kode', field: 'kode_rekening_kode', align: 'left', sortable: true },
  { name: 'uraian', label: 'Uraian', field: 'uraian', align: 'left' },
  { name: 'tipe', label: 'Tipe', field: 'tipe', align: 'center' },
  { name: 'pagu', label: 'Pagu', align: 'right', sortable: true },
  { name: 'realisasi', label: 'Realisasi', align: 'right', sortable: true },
  { name: 'sisa', label: 'Sisa', align: 'right' },
  { name: 'persen_realisasi', label: '%', align: 'right', sortable: true },
  { name: 'status_warning', label: 'Status', align: 'center' },
]

function formatRp(n) {
  return n === null || n === undefined ? '-' : 'Rp ' + Number(n).toLocaleString('id-ID')
}
function ratio(real, pagu) {
  return pagu > 0 ? Math.min(real / pagu, 1) : 0
}
function serverFilter() { return true }

async function loadData() {
  errorMsg.value = ''
  try {
    await store.getLRA(tahun.value, sampaiBulan.value)
  } catch (e) {
    errorMsg.value = e.response?.data?.error || e.message
  }
}

onMounted(loadData)
</script>

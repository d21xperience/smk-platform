<template>
  <q-page class="q-pa-md">
    <q-linear-progress v-if="store.piutangLoading || store.perSiswaLoading || store.agingLoading" indeterminate
      color="teal" class="q-mb-sm" />

    <div class="row items-center justify-between q-mb-md">
      <div class="text-h5"><q-icon name="receipt" class="q-mr-sm" />Buku Piutang Siswa</div>
      <q-btn flat icon="refresh" label="Refresh" @click="loadAll" />
    </div>

    <q-tabs v-model="tab" dense class="text-grey" active-color="teal" indicator-color="teal" align="left">
      <q-tab name="flat" icon="list" label="Per Tagihan" />
      <q-tab name="siswa" icon="person" label="Per Siswa" />
      <q-tab name="aging" icon="assessment" label="Aging Report" />
    </q-tabs>

    <q-separator />

    <!-- TAB: FLAT LIST -->
    <q-tab-panels v-model="tab" animated>
      <q-tab-panel name="flat" class="q-pa-none q-pt-md">
        <q-card class="q-mb-md">
          <q-card-section>
            <div class="row q-col-gutter-sm items-end">
              <div class="col-12 col-sm-3 col-md-2">
                <q-input v-model="filterFlat.nis" label="NIS" outlined dense clearable />
              </div>
              <div class="col-12 col-sm-3 col-md-2">
                <q-select v-model="filterFlat.aging_group" :options="agingOptions" label="Aging" outlined dense
                  emit-value map-options clearable />
              </div>
              <div class="col-12 col-sm-3 col-md-2">
                <q-btn color="teal" icon="search" label="Filter" @click="loadFlat" :loading="store.piutangLoading"
                  class="full-width" />
              </div>
            </div>
          </q-card-section>
        </q-card>

        <q-card>
          <q-table :rows="store.piutang" :columns="flatColumns" row-key="id" flat bordered dense
            :loading="store.piutangLoading" :pagination="{ rowsPerPage: 50 }">
            <template v-slot:body-cell-sisa="props">
              <q-td :props="props" class="text-right text-weight-medium">{{ formatRp(props.row.sisa) }}</q-td>
            </template>
            <template v-slot:body-cell-aging_group="props">
              <q-td :props="props" class="text-center">
                <q-chip dense :color="agingColor(props.row.aging_group)" text-color="white">
                  {{ props.row.aging_group }}
                </q-chip>
              </q-td>
            </template>
            <template v-slot:body-cell-hari_terlambat="props">
              <q-td :props="props" class="text-center">{{ props.row.hari_terlambat }} hari</q-td>
            </template>
          </q-table>
        </q-card>
      </q-tab-panel>

      <!-- TAB: PER SISWA -->
      <q-tab-panel name="siswa" class="q-pa-none q-pt-md">
        <q-card>
          <q-table :rows="store.perSiswa" :columns="siswaColumns" row-key="nis" flat bordered
            :loading="store.perSiswaLoading" :pagination="{ rowsPerPage: 50 }">
            <template v-slot:body-cell-total_sisa="props">
              <q-td :props="props" class="text-right text-weight-bold">{{ formatRp(props.row.total_sisa) }}</q-td>
            </template>
            <template v-slot:body-cell-oldest_aging_group="props">
              <q-td :props="props" class="text-center">
                <q-chip dense :color="agingColor(props.row.oldest_aging_group)" text-color="white">
                  {{ props.row.oldest_aging_group }}
                </q-chip>
              </q-td>
            </template>
            <template v-slot:body-cell-actions="props">
              <q-td :props="props" class="text-center">
                <q-btn flat dense round icon="visibility" @click="openDetail(props.row.nis)" />
              </q-td>
            </template>
          </q-table>
        </q-card>
      </q-tab-panel>

      <!-- TAB: AGING REPORT -->
      <q-tab-panel name="aging" class="q-pa-none q-pt-md">
        <div v-if="store.agingReport">
          <div class="row q-col-gutter-md q-mb-md">
            <div v-for="b in store.agingReport.buckets" :key="b.group" class="col-6 col-md-2">
              <q-card :class="agingBg(b.group)" flat>
                <q-card-section class="q-pa-sm">
                  <div class="text-caption text-grey-7">{{ b.group }}</div>
                  <div class="text-h6">{{ b.item_count }}</div>
                  <div class="text-caption">{{ formatRp(b.total_sisa) }}</div>
                </q-card-section>
              </q-card>
            </div>
          </div>
          <q-card>
            <q-card-section>
              <div class="row items-center justify-between">
                <div class="text-subtitle1 text-weight-medium">
                  Total Piutang: <span class="text-negative">{{ formatRp(store.agingReport.total_sisa) }}</span>
                </div>
                <div class="text-caption text-grey-7">Per {{ store.agingReport.as_of_date }}</div>
              </div>
            </q-card-section>
          </q-card>
        </div>
      </q-tab-panel>
    </q-tab-panels>

    <!-- Dialog Detail Siswa -->
    <q-dialog v-model="detailDialog">
      <q-card style="width: 700px; max-width: 90vw">
        <q-card-section>
          <div class="text-h6">Piutang NIS: {{ detailNIS }}</div>
        </q-card-section>
        <q-card-section v-if="detailSiswa">
          <div class="text-subtitle1 q-mb-sm">
            {{ detailSiswa.nama_siswa }} — {{ detailSiswa.kelas_nama }}
            · Total: <strong>{{ formatRp(detailSiswa.total_sisa) }}</strong>
          </div>
          <q-list separator dense>
            <q-item v-for="item in detailSiswa.items" :key="item.id">
              <q-item-section>
                <q-item-label>{{ item.jenis_iuran_nama }} ({{ item.periode_bulan }}/{{ item.periode_tahun
                  }})</q-item-label>
                <q-item-label caption>Jatuh tempo: {{ item.tanggal_jatuh_tempo }} · {{ item.hari_terlambat }}
                  hari</q-item-label>
              </q-item-section>
              <q-item-section side>
                <div class="text-weight-medium">{{ formatRp(item.sisa) }}</div>
                <q-chip dense :color="agingColor(item.aging_group)" text-color="white">{{ item.aging_group }}</q-chip>
              </q-item-section>
            </q-item>
          </q-list>
        </q-card-section>
        <q-card-actions align="right"><q-btn flat label="Tutup" v-close-popup /></q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { usePiutangStore } from '@/stores/finance/piutangStore'

const store = usePiutangStore()

const tab = ref('flat')
const detailDialog = ref(false)
const detailNIS = ref('')
const detailSiswa = ref(null)
const filterFlat = reactive({ nis: '', aging_group: '' })

const agingOptions = [
  { label: '<30', value: '<30' }, { label: '30-60', value: '30-60' },
  { label: '60-90', value: '60-90' }, { label: '>90', value: '>90' },
]

const flatColumns = [
  { name: 'nis', label: 'NIS', field: 'nis', align: 'left', sortable: true },
  { name: 'nama_siswa', label: 'Nama', field: 'nama_siswa', align: 'left' },
  { name: 'kelas_nama', label: 'Kelas', field: 'kelas_nama', align: 'left' },
  { name: 'jenis_iuran_nama', label: 'Jenis Iuran', field: 'jenis_iuran_nama', align: 'left' },
  { name: 'periode', label: 'Periode', align: 'center', field: r => `${r.periode_bulan}/${r.periode_tahun}` },
  { name: 'sisa', label: 'Sisa', align: 'right', sortable: true },
  { name: 'hari_terlambat', label: 'Terlambat', align: 'center' },
  { name: 'aging_group', label: 'Aging', align: 'center' },
]

const siswaColumns = [
  { name: 'nis', label: 'NIS', field: 'nis', align: 'left', sortable: true },
  { name: 'nama_siswa', label: 'Nama', field: 'nama_siswa', align: 'left' },
  { name: 'kelas_nama', label: 'Kelas', field: 'kelas_nama', align: 'left' },
  { name: 'jumlah_tagihan', label: 'Tagihan', field: 'jumlah_tagihan', align: 'center' },
  { name: 'total_sisa', label: 'Total Sisa', align: 'right', sortable: true },
  { name: 'oldest_aging_group', label: 'Aging Terparah', align: 'center' },
  { name: 'actions', label: 'Detail', align: 'center' },
]

function formatRp(n) {
  return n === null || n === undefined ? '-' : 'Rp ' + Number(n).toLocaleString('id-ID')
}
function agingColor(g) {
  return { '>90': 'negative', '60-90': 'deep-orange', '30-60': 'orange', '<30': 'warning', BELUM_JATUH_TEMPO: 'green' }[g] || 'grey'
}
function agingBg(g) {
  return { '>90': 'bg-red-1', '60-90': 'bg-orange-1', '30-60': 'bg-yellow-1', '<30': 'bg-green-1', BELUM_JATUH_TEMPO: 'bg-teal-1' }[g] || 'bg-grey-1'
}

async function loadFlat() {
  const params = {}
  if (filterFlat.nis) params.nis = filterFlat.nis
  if (filterFlat.aging_group) params.aging_group = filterFlat.aging_group
  await store.listPiutang(params)
}

async function loadAll() {
  try {
    await Promise.all([
      loadFlat(),
      store.listPerSiswa(),
      store.getAgingReport(),
    ])
  } catch {
    console.log('PiutangPage.vue')
  }
}

async function openDetail(nis) {
  try {
    const res = await store.getPiutangByNIS(nis)
    detailSiswa.value = res.data
    detailNIS.value = nis
    detailDialog.value = true
  } catch {
    console.log('PiutangPage.vue')
  }
}

onMounted(loadAll)
</script>

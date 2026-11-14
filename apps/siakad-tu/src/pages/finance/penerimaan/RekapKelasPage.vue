<template>
  <q-page class="q-pa-md">
    <q-linear-progress v-if="loading" indeterminate color="teal" class="q-mb-sm" />

    <div class="row items-center justify-between q-mb-md">
      <div class="text-h5"><q-icon name="class" class="q-mr-sm" />Rekap Penerimaan Per Kelas</div>
    </div>

    <q-card class="q-mb-md">
      <q-card-section>
        <div class="row q-col-gutter-sm items-end">
          <div class="col-6 col-sm-3 col-md-2">
            <q-select v-model="bulan" :options="bulanOptions" label="Bulan" outlined dense emit-value map-options />
          </div>
          <div class="col-6 col-sm-3 col-md-2">
            <q-input v-model.number="tahun" label="Tahun" outlined dense type="number" />
          </div>
          <div class="col-12 col-sm-3 col-md-2">
            <q-btn color="teal" icon="refresh" label="Muat" @click="loadData" :loading="loading" class="full-width" />
          </div>
        </div>
      </q-card-section>
    </q-card>

    <q-card v-if="rekap.length > 0">
      <q-card-section>
        <q-table :rows="rekap" :columns="columns" row-key="kelas" flat bordered :pagination="{ rowsPerPage: 50 }">
          <template v-slot:body-cell-total="props">
            <q-td :props="props">
              <span class="text-weight-bold">{{ formatRp(props.row.total) }}</span>
            </q-td>
          </template>
          <template v-slot:body-cell-actions="props">
            <q-td :props="props">
              <q-btn flat dense icon="expand" @click="showDetail(props.row)" />
            </q-td>
          </template>
        </q-table>
      </q-card-section>
    </q-card>

    <q-card v-else-if="!loading">
      <q-card-section class="text-center text-grey-6 q-pa-xl">
        <q-icon name="class" size="48px" class="q-mb-sm" />
        <div>Tidak ada penerimaan pada {{ bulanNama }} {{ tahun }}</div>
      </q-card-section>
    </q-card>

    <q-dialog v-model="detailDialog">
      <q-card style="width: 700px; max-width: 90vw">
        <q-card-section>
          <div class="text-h6">Detail Kelas: {{ detailKelas }}</div>
        </q-card-section>
        <q-card-section>
          <q-list dense separator>
            <q-item v-for="p in detailItems" :key="p.id">
              <q-item-section>
                <q-item-label>{{ p.nis }} — {{ p.jenis_iuran_nama }}</q-item-label>
                <q-item-label caption>{{ p.no_bkm }} · {{ p.tanggal }}</q-item-label>
              </q-item-section>
              <q-item-section side>{{ formatRp(p.jumlah) }}</q-item-section>
            </q-item>
          </q-list>
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Tutup" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useTransaksiStore } from '@/stores/finance/transaksiStore'

const store = useTransaksiStore()

const bulan = ref(new Date().getMonth() + 1)
const tahun = ref(new Date().getFullYear())
const loading = ref(false)
const rekap = ref([])

const detailDialog = ref(false)
const detailKelas = ref('')
const detailItems = ref([])

const bulanOptions = [
  'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
  'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
].map((label, i) => ({ label, value: i + 1 }))

const bulanNama = ref('')
const columns = [
  { name: 'kelas', label: 'Kelas', field: 'kelas', align: 'left', sortable: true },
  { name: 'count', label: 'Jumlah Transaksi', field: 'count', align: 'center', sortable: true },
  { name: 'unique_nis', label: 'Siswa Bayar', field: 'unique_nis', align: 'center' },
  { name: 'total', label: 'Total', align: 'right', sortable: true },
  { name: 'actions', label: 'Detail', align: 'center' },
]

function formatRp(n) {
  return n === null || n === undefined ? '-' : 'Rp ' + Number(n).toLocaleString('id-ID')
}

function showDetail(row) {
  detailKelas.value = row.kelas
  detailItems.value = row.items
  detailDialog.value = true
}

async function loadData() {
  loading.value = true
  bulanNama.value = bulanOptions.find(b => b.value === bulan.value)?.label || ''
  try {
    // Hitung tanggal awal & akhir bulan
    const y = tahun.value
    const m = String(bulan.value).padStart(2, '0')
    const firstDay = `${y}-${m}-01`
    const lastDay = new Date(y, bulan.value, 0).toISOString().split('T')[0]

    const res = await store.listPenerimaan({
      tanggal_dari: firstDay,
      tanggal_sampai: lastDay,
      limit: 5000,
    })
    const list = res.data || []

    // Group by kelas
    const map = {}
    list.forEach(p => {
      const k = p.kelas_nama || '(tanpa kelas)'
      if (!map[k]) map[k] = { kelas: k, count: 0, total: 0, items: [], nisSet: new Set() }
      map[k].count++
      map[k].total += p.jumlah
      map[k].items.push(p)
      map[k].nisSet.add(p.nis)
    })

    rekap.value = Object.values(map).map(k => ({
      kelas: k.kelas,
      count: k.count,
      unique_nis: k.nisSet.size,
      total: k.total,
      items: k.items,
    })).sort((a, b) => a.kelas.localeCompare(b.kelas))
  } catch {
    rekap.value = []
  } finally {
    loading.value = false
  }
}

onMounted(loadData)
</script>

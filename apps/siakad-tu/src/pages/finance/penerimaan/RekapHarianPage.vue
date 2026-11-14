<template>
  <q-page class="q-pa-md">
    <q-linear-progress v-if="loading" indeterminate color="teal" class="q-mb-sm" />

    <div class="row items-center justify-between q-mb-md">
      <div class="text-h5"><q-icon name="event" class="q-mr-sm" />Rekap Penerimaan Harian</div>
    </div>

    <q-card class="q-mb-md">
      <q-card-section>
        <div class="row q-col-gutter-sm items-end">
          <div class="col-12 col-sm-4 col-md-3">
            <q-input v-model="tanggal" label="Tanggal *" outlined dense type="date" />
          </div>
          <div class="col-12 col-sm-4 col-md-3">
            <q-btn color="teal" icon="refresh" label="Muat" @click="loadData" :loading="loading" class="full-width" />
          </div>
        </div>
      </q-card-section>
    </q-card>

    <div v-if="penerimaan.length > 0">
      <div class="row q-col-gutter-md q-mb-md">
        <div class="col-6 col-md-3">
          <q-card class="bg-teal-1">
            <q-card-section>
              <div class="text-caption text-grey-7">Total Penerimaan</div>
              <div class="text-h6 text-teal-10">{{ formatRp(totalHarian) }}</div>
            </q-card-section>
          </q-card>
        </div>
        <div class="col-6 col-md-3">
          <q-card class="bg-blue-1">
            <q-card-section>
              <div class="text-caption text-grey-7">Jumlah Transaksi</div>
              <div class="text-h6 text-blue-10">{{ penerimaan.length }}</div>
            </q-card-section>
          </q-card>
        </div>
        <div class="col-6 col-md-3">
          <q-card class="bg-orange-1">
            <q-card-section>
              <div class="text-caption text-grey-7">Siswa Bayar</div>
              <div class="text-h6 text-orange-10">{{ uniqueNIS }}</div>
            </q-card-section>
          </q-card>
        </div>
        <div class="col-6 col-md-3">
          <q-card class="bg-purple-1">
            <q-card-section>
              <div class="text-caption text-grey-7">Kelas Terlibat</div>
              <div class="text-h6 text-purple-10">{{ rekapPerKelas.length }}</div>
            </q-card-section>
          </q-card>
        </div>
      </div>

      <q-card class="q-mb-md">
        <q-card-section>
          <div class="text-subtitle1 text-weight-medium q-mb-sm">Rekap Per Kelas</div>
          <q-list separator>
            <q-expansion-item v-for="k in rekapPerKelas" :key="k.kelas" :label="`${k.kelas} — ${k.count} transaksi`"
              :caption="formatRp(k.total)" default-opened>
              <q-card flat>
                <q-card-section>
                  <q-list dense>
                    <q-item v-for="p in k.items" :key="p.id">
                      <q-item-section>
                        <q-item-label>{{ p.nis }} — {{ p.jenis_iuran_nama }}</q-item-label>
                        <q-item-label caption>{{ p.no_bkm }} · {{ p.metode_bayar }}</q-item-label>
                      </q-item-section>
                      <q-item-section side>
                        <div class="text-weight-medium">{{ formatRp(p.jumlah) }}</div>
                      </q-item-section>
                    </q-item>
                  </q-list>
                </q-card-section>
              </q-card>
            </q-expansion-item>
          </q-list>
        </q-card-section>
      </q-card>
    </div>

    <q-card v-else-if="!loading">
      <q-card-section class="text-center text-grey-6 q-pa-xl">
        <q-icon name="event_busy" size="48px" class="q-mb-sm" />
        <div>Tidak ada penerimaan pada tanggal {{ tanggal }}</div>
      </q-card-section>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useTransaksiStore } from '@/stores/finance/transaksiStore'

const store = useTransaksiStore()

const tanggal = ref(new Date().toISOString().split('T')[0])
const penerimaan = ref([])
const loading = ref(false)

const totalHarian = computed(() => penerimaan.value.reduce((s, p) => s + p.jumlah, 0))
const uniqueNIS = computed(() => new Set(penerimaan.value.map(p => p.nis)).size)

const rekapPerKelas = computed(() => {
  const map = {}
  penerimaan.value.forEach(p => {
    const k = p.kelas_nama || '(tanpa kelas)'
    if (!map[k]) map[k] = { kelas: k, count: 0, total: 0, items: [] }
    map[k].count++
    map[k].total += p.jumlah
    map[k].items.push(p)
  })
  return Object.values(map).sort((a, b) => a.kelas.localeCompare(b.kelas))
})

function formatRp(n) {
  return n === null || n === undefined ? '-' : 'Rp ' + Number(n).toLocaleString('id-ID')
}

async function loadData() {
  loading.value = true
  try {
    const res = await store.listPenerimaan({
      tanggal_dari: tanggal.value,
      tanggal_sampai: tanggal.value,
      limit: 1000,
    })
    penerimaan.value = res.data || []
  } catch {
    penerimaan.value = []
  } finally {
    loading.value = false
  }
}

onMounted(loadData)
</script>

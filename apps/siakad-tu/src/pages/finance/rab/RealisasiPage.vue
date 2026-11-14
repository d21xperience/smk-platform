<template>
  <q-page class="q-pa-md">
    <q-linear-progress v-if="store.realisasiLoading" indeterminate color="teal" class="q-mb-sm" />

    <div class="row items-center justify-between q-mb-md">
      <div class="text-h5"><q-icon name="trending_up" class="q-mr-sm" />Realisasi Anggaran</div>
    </div>

    <q-card class="q-mb-md">
      <q-card-section>
        <div class="row q-col-gutter-sm items-end">
          <div class="col-6 col-sm-3 col-md-2">
            <q-input v-model.number="tahun" label="Tahun Anggaran" outlined dense type="number"
              min="2000" max="2100" />
          </div>
          <div class="col-6 col-sm-3 col-md-2">
            <q-select v-model="sampaiBulan" :options="bulanOptions" label="Sampai Bulan" outlined dense
              emit-value map-options />
          </div>
          <div class="col-12 col-sm-3 col-md-2">
            <q-btn color="teal" icon="refresh" label="Muat" @click="loadData" :loading="store.realisasiLoading"
              class="full-width" />
          </div>
        </div>
      </q-card-section>
    </q-card>

    <div v-if="report">
      <!-- Summary -->
      <div class="row q-col-gutter-md q-mb-md">
        <div class="col-12 col-md-6">
          <q-card class="bg-green-1">
            <q-card-section>
              <div class="text-caption text-grey-7">Pendapatan (Pagu → Realisasi)</div>
              <div class="text-h6 text-green-10">{{ formatRp(report.total_pagu_pendapatan) }}</div>
              <div class="text-subtitle2">{{ formatRp(report.total_realisasi_pendapatan) }}
                ({{ persen(report.total_realisasi_pendapatan, report.total_pagu_pendapatan) }}%)</div>
              <q-linear-progress :value="progressRatio(report.total_realisasi_pendapatan, report.total_pagu_pendapatan)"
                color="green-8" track-color="green-2" class="q-mt-sm" size="8px" />
            </q-card-section>
          </q-card>
        </div>
        <div class="col-12 col-md-6">
          <q-card class="bg-orange-1">
            <q-card-section>
              <div class="text-caption text-grey-7">Belanja (Pagu → Realisasi)</div>
              <div class="text-h6 text-orange-10">{{ formatRp(report.total_pagu_belanja) }}</div>
              <div class="text-subtitle2">{{ formatRp(report.total_realisasi_belanja) }}
                ({{ persen(report.total_realisasi_belanja, report.total_pagu_belanja) }}%)</div>
              <q-linear-progress :value="progressRatio(report.total_realisasi_belanja, report.total_pagu_belanja)"
                :color="belanjaBarColor" track-color="orange-2" class="q-mt-sm" size="8px" />
            </q-card-section>
          </q-card>
        </div>
      </div>

      <!-- Per Pos -->
      <q-card>
        <q-card-section>
          <div class="text-subtitle1 text-weight-medium q-mb-sm">Realisasi per Pos Anggaran</div>
          <q-list separator>
            <q-item v-for="item in report.items" :key="item.rab_detail_id" class="q-py-md">
              <q-item-section>
                <q-item-label>
                  <span class="text-weight-medium">{{ item.kode_rekening_kode }}</span>
                  — {{ item.uraian }}
                  <q-chip dense :color="item.tipe === 'PENDAPATAN' ? 'green-6' : 'orange-6'" text-color="white"
                    class="q-ml-sm">
                    {{ item.tipe }}
                  </q-chip>
                  <q-badge v-if="item.status_warning !== 'NORMAL'"
                    :color="warningColor(item.status_warning)" class="q-ml-sm">
                    {{ warningLabel(item.status_warning) }}
                  </q-badge>
                </q-item-label>
                <q-item-label caption>
                  Pagu: {{ formatRp(item.pagu) }} · Realisasi: {{ formatRp(item.realisasi) }} ·
                  Sisa: {{ formatRp(item.sisa) }}
                </q-item-label>
                <q-linear-progress :value="progressRatio(item.realisasi, item.pagu)"
                  :color="rowBarColor(item)" track-color="grey-3" class="q-mt-sm" size="10px" />
              </q-item-section>
              <q-item-section side style="min-width: 100px">
                <div class="text-h6" :class="rowBarColor(item) === 'negative' ? 'text-negative' : ''">
                  {{ item.persen_realisasi.toFixed(1) }}%
                </div>
              </q-item-section>
            </q-item>
          </q-list>
        </q-card-section>
      </q-card>
    </div>

    <q-card v-else-if="!store.realisasiLoading">
      <q-card-section class="text-center text-grey-6 q-pa-xl">
        <q-icon name="assessment" size="48px" class="q-mb-sm" />
        <div>Belum ada data realisasi. Pilih tahun dan klik Muat.</div>
      </q-card-section>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useRabStore } from '@/stores/finance/rabStore'

const $q = useQuasar()
const store = useRabStore()

const tahun = ref(new Date().getFullYear())
const sampaiBulan = ref(new Date().getMonth() + 1)

const bulanOptions = [
  'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
  'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
].map((label, i) => ({ label, value: i + 1 }))

const report = computed(() => store.realisasi)

const belanjaBarColor = computed(() => {
  if (!report.value) return 'orange-8'
  const pct = persen(report.value.total_realisasi_belanja, report.value.total_pagu_belanja)
  if (pct > 95) return 'negative'
  if (pct > 80) return 'warning'
  return 'orange-8'
})

function formatRp(n) {
  return n === null || n === undefined ? '-' : 'Rp ' + Number(n).toLocaleString('id-ID')
}
function persen(real, pagu) {
  return pagu > 0 ? ((real / pagu) * 100).toFixed(1) : '0.0'
}
function progressRatio(real, pagu) {
  return pagu > 0 ? Math.min(real / pagu, 1) : 0
}
function rowBarColor(item) {
  if (item.persen_realisasi > 95) return 'negative'
  if (item.persen_realisasi > 80) return 'warning'
  if (item.tipe === 'PENDAPATAN') return 'green-8'
  return 'orange-8'
}
function warningColor(w) {
  return { OVER_BUDGET_WARNING: 'negative', UNDER_BUDGET_WARNING: 'warning' }[w] || 'grey'
}
function warningLabel(w) {
  return { OVER_BUDGET_WARNING: 'Over Budget', UNDER_BUDGET_WARNING: 'Under Budget' }[w] || w
}

async function loadData() {
  try {
    await store.getRealisasi(tahun.value, sampaiBulan.value)
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || e.message, timeout: 5000 })
  }
}

onMounted(loadData)
</script>

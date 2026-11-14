<template>
  <q-page class="q-pa-md">
    <q-linear-progress v-if="store.neracaLoading" indeterminate color="teal" class="q-mb-sm" />

    <div class="row items-center justify-between q-mb-md">
      <div class="text-h5"><q-icon name="account_balance" class="q-mr-sm" />Neraca Keuangan</div>
    </div>

    <q-card class="q-mb-md">
      <q-card-section>
        <div class="row q-col-gutter-sm items-end">
          <div class="col-6 col-md-2">
            <q-select v-model="mode" :options="modeOptions" label="Mode" outlined dense emit-value map-options
              @update:model-value="loadData" />
          </div>
          <div v-if="mode === 'LIVE'" class="col-6 col-md-2">
            <q-input v-model="tanggal" label="Per Tanggal (kosong = hari ini)" outlined dense type="date" />
          </div>
          <div v-if="mode === 'SNAPSHOT'" class="col-12 col-md-4">
            <q-select v-model="penutupanID" :options="penutupanOptions" label="Pilih Penutupan" outlined dense
              emit-value map-options />
          </div>
          <div class="col-12 col-md-2">
            <q-btn color="teal" icon="refresh" label="Muat" @click="loadData" :loading="store.neracaLoading"
              class="full-width" />
          </div>
        </div>
      </q-card-section>
    </q-card>

    <div v-if="report">
      <!-- Balance Banner -->
      <q-banner :class="report.is_balance ? 'bg-positive' : 'bg-negative'" class="text-white q-mb-md" rounded>
        <template v-slot:avatar>
          <q-icon :name="report.is_balance ? 'check_circle' : 'error'" />
        </template>
        <div class="text-weight-bold">
          {{ report.is_balance ? 'BALANCE' : 'TIDAK BALANCE' }}
          — Aset: {{ formatRp(report.total_aset) }} =
          Kewajiban: {{ formatRp(report.total_kewajiban) }} +
          Ekuitas: {{ formatRp(report.total_ekuitas) }}
        </div>
        <div v-if="!report.is_balance" class="text-caption">
          Selisih: {{ formatRp(report.selisih) }}
        </div>
      </q-banner>

      <div class="row q-col-gutter-md">
        <!-- ASET -->
        <div class="col-12 col-md-6">
          <q-card>
            <q-card-section class="bg-blue-1">
              <div class="text-subtitle1 text-weight-bold">ASET</div>
            </q-card-section>
            <q-list separator>
              <q-item v-for="item in report.aset" :key="item.kode_rekening_id">
                <q-item-section>
                  <q-item-label>{{ item.kode }} — {{ item.nama }}</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <span class="text-weight-medium" :class="item.saldo < 0 ? 'text-negative' : ''">
                    {{ formatRp(item.saldo) }}
                  </span>
                </q-item-section>
              </q-item>
              <q-item v-if="report.aset.length === 0">
                <q-item-section class="text-grey-6">Tidak ada data.</q-item-section>
              </q-item>
            </q-list>
            <q-separator />
            <q-card-section>
              <div class="row justify-between text-h6">
                <div>Total Aset</div>
                <div class="text-blue-10">{{ formatRp(report.total_aset) }}</div>
              </div>
            </q-card-section>
          </q-card>
        </div>

        <!-- KEWAJIBAN + EKUITAS -->
        <div class="col-12 col-md-6">
          <q-card class="q-mb-md">
            <q-card-section class="bg-red-1">
              <div class="text-subtitle1 text-weight-bold">KEWAJIBAN</div>
            </q-card-section>
            <q-list separator>
              <q-item v-for="item in report.kewajiban" :key="item.kode_rekening_id">
                <q-item-section>
                  <q-item-label>{{ item.kode }} — {{ item.nama }}</q-item-label>
                </q-item-section>
                <q-item-section side>{{ formatRp(item.saldo) }}</q-item-section>
              </q-item>
              <q-item v-if="report.kewajiban.length === 0">
                <q-item-section class="text-grey-6">Tidak ada kewajiban.</q-item-section>
              </q-item>
            </q-list>
            <q-separator />
            <q-card-section>
              <div class="row justify-between text-subtitle1">
                <div>Total Kewajiban</div>
                <div class="text-red-10 text-weight-bold">{{ formatRp(report.total_kewajiban) }}</div>
              </div>
            </q-card-section>
          </q-card>

          <q-card>
            <q-card-section class="bg-purple-1">
              <div class="text-subtitle1 text-weight-bold">EKUITAS</div>
            </q-card-section>
            <q-list separator>
              <q-item v-for="item in report.ekuitas" :key="item.kode_rekening_id">
                <q-item-section>
                  <q-item-label>{{ item.kode }} — {{ item.nama }}</q-item-label>
                </q-item-section>
                <q-item-section side>{{ formatRp(item.saldo) }}</q-item-section>
              </q-item>
            </q-list>
            <q-separator />
            <q-card-section>
              <div class="row justify-between text-subtitle1">
                <div>Total Ekuitas</div>
                <div class="text-purple-10 text-weight-bold">{{ formatRp(report.total_ekuitas) }}</div>
              </div>
            </q-card-section>
          </q-card>
        </div>

        <!-- PENDAPATAN & BELANJA -->
        <div class="col-12 col-md-6">
          <q-card>
            <q-card-section class="bg-green-1">
              <div class="text-subtitle1 text-weight-bold">PENDAPATAN</div>
            </q-card-section>
            <q-list separator dense>
              <q-item v-for="item in report.pendapatan" :key="item.kode_rekening_id">
                <q-item-section><q-item-label>{{ item.kode }} — {{ item.nama }}</q-item-label></q-item-section>
                <q-item-section side>{{ formatRp(item.saldo) }}</q-item-section>
              </q-item>
              <q-item v-if="report.pendapatan.length === 0">
                <q-item-section class="text-grey-6">—</q-item-section>
              </q-item>
            </q-list>
            <q-separator />
            <q-card-section>
              <div class="row justify-between text-subtitle1">
                <div>Total Pendapatan</div>
                <div class="text-green-10 text-weight-bold">{{ formatRp(report.total_pendapatan) }}</div>
              </div>
            </q-card-section>
          </q-card>
        </div>

        <div class="col-12 col-md-6">
          <q-card>
            <q-card-section class="bg-orange-1">
              <div class="text-subtitle1 text-weight-bold">BELANJA</div>
            </q-card-section>
            <q-list separator dense>
              <q-item v-for="item in report.belanja" :key="item.kode_rekening_id">
                <q-item-section><q-item-label>{{ item.kode }} — {{ item.nama }}</q-item-label></q-item-section>
                <q-item-section side>{{ formatRp(item.saldo) }}</q-item-section>
              </q-item>
              <q-item v-if="report.belanja.length === 0">
                <q-item-section class="text-grey-6">—</q-item-section>
              </q-item>
            </q-list>
            <q-separator />
            <q-card-section>
              <div class="row justify-between text-subtitle1">
                <div>Total Belanja</div>
                <div class="text-orange-10 text-weight-bold">{{ formatRp(report.total_belanja) }}</div>
              </div>
              <div class="row justify-between text-subtitle1 q-mt-xs">
                <div>Surplus/Defisit</div>
                <div :class="report.surplus_defisit >= 0 ? 'text-positive' : 'text-negative'">
                  {{ formatRp(report.surplus_defisit) }}
                </div>
              </div>
            </q-card-section>
          </q-card>
        </div>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useLaporanStore } from '@/stores/finance/laporanStore'
import { usePenutupanStore } from '@/stores/finance/penutupanStore'

const store = useLaporanStore()
const penutupanStore = usePenutupanStore()

const mode = ref('LIVE')
const tanggal = ref('')
const penutupanID = ref('')
const report = computed(() => store.neraca)

const modeOptions = [
  { label: 'Live', value: 'LIVE' },
  { label: 'Snapshot', value: 'SNAPSHOT' },
]

const penutupanOptions = computed(() =>
  penutupanStore.penutupanList.map(p => ({
    label: `${p.periode_bulan}/${p.periode_tahun} — ${p.status}`,
    value: p.id,
  }))
)

function formatRp(n) {
  return n === null || n === undefined ? '-' : 'Rp ' + Number(n).toLocaleString('id-ID')
}

async function loadData() {
  try {
    if (mode.value === 'SNAPSHOT') {
      if (!penutupanID.value) return
      await store.getNeracaSnapshot(penutupanID.value)
    } else {
      const params = {}
      if (tanggal.value) params.tanggal = tanggal.value
      await store.getNeraca(params)
    }
  } catch {
    console.log('NeracaPage.vue')
  }
}

onMounted(async () => {
  if (penutupanStore.penutupanList.length === 0) {
    await penutupanStore.listPenutupan()
  }
  await loadData()
})
</script>

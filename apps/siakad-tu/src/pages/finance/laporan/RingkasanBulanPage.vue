<template>
  <q-page class="q-pa-md">
    <q-linear-progress v-if="store.ringkasanLoading" indeterminate color="teal" class="q-mb-sm" />

    <div class="row items-center justify-between q-mb-md">
      <div class="text-h5"><q-icon name="calendar_month" class="q-mr-sm" />Ringkasan Bulanan</div>
    </div>

    <q-card class="q-mb-md">
      <q-card-section>
        <div class="row q-col-gutter-sm items-end">
          <div class="col-6 col-md-2">
            <q-select v-model="bulan" :options="bulanOptions" label="Bulan" outlined dense emit-value map-options />
          </div>
          <div class="col-6 col-md-2">
            <q-input v-model.number="tahun" label="Tahun" outlined dense type="number" />
          </div>
          <div class="col-12 col-md-2">
            <q-btn color="teal" icon="refresh" label="Muat" @click="loadData" :loading="store.ringkasanLoading"
              class="full-width" />
          </div>
        </div>
      </q-card-section>
    </q-card>

    <div v-if="ringkasan">
      <div class="row q-col-gutter-md q-mb-md">
        <div class="col-6 col-md-3">
          <q-card class="bg-blue-1"><q-card-section>
              <div class="text-caption">Saldo Awal</div>
              <div class="text-h6 text-blue-10">{{ formatRp(ringkasan.saldo_awal) }}</div>
            </q-card-section></q-card>
        </div>
        <div class="col-6 col-md-3">
          <q-card class="bg-green-1"><q-card-section>
              <div class="text-caption">Total Debet</div>
              <div class="text-h6 text-green-10">+{{ formatRp(ringkasan.total_debet) }}</div>
            </q-card-section></q-card>
        </div>
        <div class="col-6 col-md-3">
          <q-card class="bg-red-1"><q-card-section>
              <div class="text-caption">Total Kredit</div>
              <div class="text-h6 text-negative">-{{ formatRp(ringkasan.total_kredit) }}</div>
            </q-card-section></q-card>
        </div>
        <div class="col-6 col-md-3">
          <q-card class="bg-teal-1"><q-card-section>
              <div class="text-caption">Saldo Akhir</div>
              <div class="text-h6" :class="ringkasan.saldo_akhir < 0 ? 'text-negative' : 'text-teal-10'">
                {{ formatRp(ringkasan.saldo_akhir) }}
              </div>
            </q-card-section></q-card>
        </div>
      </div>

      <q-card class="q-mb-md">
        <q-card-section>
          <div class="text-subtitle2 text-grey-7">
            <q-icon name="receipt" /> {{ ringkasan.jumlah_transaksi }} transaksi pada {{ bulanName }} {{
              ringkasan.periode_tahun }}
          </div>
        </q-card-section>
      </q-card>

      <div class="row q-col-gutter-md">
        <div class="col-12 col-md-6">
          <q-card>
            <q-card-section class="bg-green-1">
              <div class="text-subtitle1 text-weight-medium">Top 3 Debet (Penerimaan)</div>
            </q-card-section>
            <q-list separator>
              <q-item v-for="(item, i) in ringkasan.top_debet" :key="i">
                <q-item-section avatar>
                  <q-avatar color="green-7" text-color="white">{{ i + 1 }}</q-avatar>
                </q-item-section>
                <q-item-section>
                  <q-item-label>{{ item.kode_rekening }} — {{ item.nama }}</q-item-label>
                  <q-item-label caption>{{ item.jumlah_transaksi }} transaksi</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <span class="text-green-10 text-weight-bold">{{ formatRp(item.total) }}</span>
                </q-item-section>
              </q-item>
              <q-item v-if="ringkasan.top_debet.length === 0">
                <q-item-section class="text-grey-6">Tidak ada penerimaan.</q-item-section>
              </q-item>
            </q-list>
          </q-card>
        </div>
        <div class="col-12 col-md-6">
          <q-card>
            <q-card-section class="bg-orange-1">
              <div class="text-subtitle1 text-weight-medium">Top 3 Kredit (Pengeluaran)</div>
            </q-card-section>
            <q-list separator>
              <q-item v-for="(item, i) in ringkasan.top_kredit" :key="i">
                <q-item-section avatar>
                  <q-avatar color="orange-7" text-color="white">{{ i + 1 }}</q-avatar>
                </q-item-section>
                <q-item-section>
                  <q-item-label>{{ item.kode_rekening }} — {{ item.nama }}</q-item-label>
                  <q-item-label caption>{{ item.jumlah_transaksi }} transaksi</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <span class="text-orange-10 text-weight-bold">{{ formatRp(item.total) }}</span>
                </q-item-section>
              </q-item>
              <q-item v-if="ringkasan.top_kredit.length === 0">
                <q-item-section class="text-grey-6">Tidak ada pengeluaran.</q-item-section>
              </q-item>
            </q-list>
          </q-card>
        </div>
      </div>
    </div>

    <q-card v-else-if="!store.ringkasanLoading">
      <q-card-section class="text-center text-grey-6 q-pa-xl">
        <q-icon name="calendar_month" size="48px" class="q-mb-sm" />
        <div>Pilih bulan & tahun, klik Muat.</div>
      </q-card-section>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useLaporanStore } from '@/stores/finance/laporanStore'

const store = useLaporanStore()

const bulan = ref(new Date().getMonth() + 1)
const tahun = ref(new Date().getFullYear())

const bulanOptions = [
  'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
  'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
].map((label, i) => ({ label, value: i + 1 }))

const ringkasan = computed(() => store.ringkasanBulan)
const bulanName = computed(() => bulanOptions.find(b => b.value === bulan.value)?.label || '')

function formatRp(n) {
  return n === null || n === undefined ? '-' : 'Rp ' + Number(n).toLocaleString('id-ID')
}

async function loadData() {
  try {
    await store.getRingkasanBulan(bulan.value, tahun.value)
  } catch {
    console.log('RingkasanBulanPage.vue')
  }
}

onMounted(loadData)
</script>

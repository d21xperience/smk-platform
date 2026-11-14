<template>
  <q-page class="q-pa-md">
    <q-linear-progress v-if="loading" indeterminate color="teal" class="q-mb-sm" />

    <div class="row items-center justify-between q-mb-md">
      <div class="text-h5">
        <q-icon name="dashboard" class="q-mr-sm" />
        Dashboard Keuangan
      </div>
      <div class="text-caption text-grey-7">{{ asOfDate }}</div>
    </div>

    <q-banner v-if="error" class="bg-negative text-white q-mb-md" rounded>
      <template v-slot:avatar><q-icon name="error" /></template>
      {{ error }}
    </q-banner>

    <div class="row q-col-gutter-md q-mb-md">
      <div class="col-12 col-sm-6 col-md-3">
        <q-card :class="saldoKas < 0 ? 'bg-red-1' : 'bg-teal-1'">
          <q-card-section>
            <div class="text-caption text-grey-7">Saldo Kas</div>
            <div class="text-h5 text-weight-bold" :class="saldoKas < 0 ? 'text-negative' : 'text-teal-10'">
              {{ formatRp(saldoKas) }}
            </div>
            <div class="text-caption text-grey-6 q-mt-xs">
              Debet: {{ formatRp(kasDebet) }} · Kredit: {{ formatRp(kasKredit) }}
            </div>
          </q-card-section>
        </q-card>
      </div>

      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="bg-orange-1">
          <q-card-section>
            <div class="text-caption text-grey-7">Total Piutang</div>
            <div class="text-h5 text-weight-bold text-orange-10">{{ formatRp(totalPiutang) }}</div>
            <div class="text-caption text-grey-6 q-mt-xs">
              {{ jumlahPiutangItem }} tagihan · {{ jumlahPiutangSiswa }} siswa
            </div>
          </q-card-section>
        </q-card>
      </div>

      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="bg-blue-1">
          <q-card-section>
            <div class="text-caption text-grey-7">Realisasi RAB</div>
            <div class="text-h5 text-weight-bold text-blue-10">{{ persenRealisasi.toFixed(2) }}%</div>
            <q-linear-progress :value="Math.min(persenRealisasi, 100) / 100" color="blue-8" track-color="blue-2" class="q-mt-xs" />
          </q-card-section>
        </q-card>
      </div>

      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="bg-purple-1">
          <q-card-section>
            <div class="text-caption text-grey-7">Pengajuan Pending</div>
            <div class="text-h5 text-weight-bold text-purple-10">{{ totalPengajuanPending }}</div>
            <div class="text-caption text-grey-6 q-mt-xs">{{ formatRp(nominalPending) }}</div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <q-card class="q-mb-md">
      <q-card-section>
        <div class="text-subtitle1 text-weight-medium q-mb-sm">
          <q-icon name="warning" class="q-mr-sm" />Alerts
        </div>
        <div v-if="alerts.length === 0" class="text-grey-6">Tidak ada alert.</div>
        <q-list v-else dense separator>
          <q-item v-for="(a, i) in alerts" :key="i">
            <q-item-section avatar>
              <q-icon :name="alertIcon(a.severity)" :color="alertColor(a.severity)" />
            </q-item-section>
            <q-item-section>
              <q-item-label>{{ a.message }}</q-item-label>
              <q-item-label caption>{{ a.category }}</q-item-label>
            </q-item-section>
          </q-item>
        </q-list>
      </q-card-section>
    </q-card>

    <q-card>
      <q-card-section>
        <div class="text-subtitle1 text-weight-medium q-mb-sm">
          <q-icon name="calendar_month" class="q-mr-sm" />Aging Piutang
        </div>
        <div class="row q-col-gutter-sm">
          <div v-for="b in agingBuckets" :key="b.group" class="col-6 col-sm-4 col-md-2">
            <q-card :class="bucketColor(b.group)" flat bordered>
              <q-card-section class="text-center q-pa-sm">
                <div class="text-caption text-grey-7">{{ b.group }}</div>
                <div class="text-h6">{{ b.item_count }}</div>
                <div class="text-caption">{{ formatRp(b.total_sisa) }}</div>
              </q-card-section>
            </q-card>
          </div>
        </div>
      </q-card-section>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { api } from '@/boot/axios'
import { useAuthStore } from '@/stores/authStore'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const authStore = useAuthStore()

const loading = ref(false)
const error = ref('')
const asOfDate = ref('')
const saldoKas = ref(0)
const kasDebet = ref(0)
const kasKredit = ref(0)
const totalPiutang = ref(0)
const jumlahPiutangItem = ref(0)
const jumlahPiutangSiswa = ref(0)
const persenRealisasi = ref(0)
const totalPengajuanPending = ref(0)
const nominalPending = ref(0)
const alerts = ref([])
const agingBuckets = ref([])

const currentUser = computed(() => authStore.user || authStore.currentUser || {})
const schoolID = computed(() => currentUser.value.school_id || '11111111-1111-1111-1111-111111111111')
const periodID = computed(() => currentUser.value.academic_period_id || '20261')

function formatRp(n) {
  if (n === null || n === undefined) return '-'
  return 'Rp ' + Number(n).toLocaleString('id-ID')
}
function alertIcon(s) {
  return s === 'DANGER' ? 'error' : s === 'WARNING' ? 'warning' : 'info'
}
function alertColor(s) {
  return s === 'DANGER' ? 'negative' : s === 'WARNING' ? 'orange' : 'info'
}
function bucketColor(g) {
  if (g === '>90') return 'bg-red-1'
  if (g === '60-90') return 'bg-orange-1'
  if (g === '30-60') return 'bg-yellow-1'
  return 'bg-green-1'
}

async function loadDashboard() {
  loading.value = true
  error.value = ''
  try {
    const res = await api.get('/keuangan/laporan/dashboard', {
      params: { school_id: schoolID.value, academic_period_id: periodID.value }
    })
    const d = res.data.data
    asOfDate.value = d.as_of_date || ''
    saldoKas.value = d.kas?.saldo_terkini ?? 0
    kasDebet.value = d.kas?.total_debet_bulan_ini ?? 0
    kasKredit.value = d.kas?.total_kredit_bulan_ini ?? 0
    totalPiutang.value = d.piutang?.total_sisa ?? 0
    jumlahPiutangItem.value = d.piutang?.jumlah_item ?? 0
    jumlahPiutangSiswa.value = d.piutang?.jumlah_siswa ?? 0
    persenRealisasi.value = d.rab?.persen_realisasi ?? 0
    totalPengajuanPending.value = (d.pengajuan?.menunggu_verifikasi ?? 0) + (d.pengajuan?.menunggu_persetujuan ?? 0)
    nominalPending.value = d.pengajuan?.total_nominal_pending ?? 0
    alerts.value = d.alerts ?? []
    agingBuckets.value = d.piutang?.aging_summary ?? []
  } catch (e) {
    error.value = e.response?.data?.error || e.message
    $q.notify({ type: 'negative', message: 'Gagal memuat dashboard: ' + error.value })
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadDashboard()
})
</script>

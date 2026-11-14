<template>
  <q-page class="q-pa-md">
    <q-linear-progress v-if="store.penutupanLoading || store.previewLoading" indeterminate color="teal"
      class="q-mb-sm" />

    <!-- LIST MODE -->
    <div v-if="mode === 'list'">
      <div class="row items-center justify-between q-mb-md">
        <div class="text-h5"><q-icon name="lock" class="q-mr-sm" />Penutupan Buku Kas</div>
        <q-btn color="teal" icon="add" label="Tutup Buku Baru" @click="openWizard" />
      </div>

      <q-card v-if="store.penutupanList.length === 0 && !store.penutupanLoading">
        <q-card-section class="text-center text-grey-6 q-pa-xl">
          <q-icon name="lock" size="48px" class="q-mb-sm" />
          <div>Belum ada penutupan.</div>
        </q-card-section>
      </q-card>

      <q-card v-else>
        <q-table :rows="store.penutupanList" :columns="listColumns" row-key="id" flat bordered
          :pagination="{ rowsPerPage: 20 }" :loading="store.penutupanLoading">
          <template v-slot:body-cell-status="props">
            <q-td :props="props" class="text-center">
              <q-chip dense :color="statusColor(props.row.status)" text-color="white">
                {{ statusLabel(props.row.status) }}
              </q-chip>
            </q-td>
          </template>
          <template v-slot:body-cell-saldo_akhir_buku="props">
            <q-td :props="props" class="text-right">{{ formatRp(props.row.saldo_akhir_buku) }}</q-td>
          </template>
          <template v-slot:body-cell-selisih_kas="props">
            <q-td :props="props" class="text-right">
              <span :class="props.row.selisih_kas !== 0 ? 'text-negative text-weight-bold' : ''">
                {{ formatRp(props.row.selisih_kas) }}
              </span>
            </q-td>
          </template>
          <template v-slot:body-cell-actions="props">
            <q-td :props="props" class="text-center">
              <q-btn flat dense round icon="visibility" @click="openDetail(props.row.id)" />
              <q-btn v-if="props.row.status === 'CLOSED'" flat dense round icon="fingerprint" color="purple"
                @click="openSign(props.row)" />
            </q-td>
          </template>
        </q-table>
      </q-card>
    </div>

    <!-- WIZARD MODE -->
    <div v-else-if="mode === 'wizard'">
      <div class="row items-center justify-between q-mb-md">
        <div class="text-h5"><q-icon name="lock" class="q-mr-sm" />Tutup Buku Baru</div>
        <q-btn flat icon="arrow_back" label="Kembali" @click="mode = 'list'" />
      </div>

      <q-stepper v-model="step" vertical color="teal" animated>
        <!-- Step 1: Preview Snapshot -->
        <q-step :name="1" title="Preview Snapshot" icon="preview" :done="step > 1">
          <div class="row q-col-gutter-sm q-mb-md">
            <div class="col-6 col-md-2">
              <q-select v-model="wizard.bulan" :options="bulanOptions" label="Bulan" outlined dense emit-value
                map-options />
            </div>
            <div class="col-6 col-md-2">
              <q-input v-model.number="wizard.tahun" label="Tahun" outlined dense type="number" />
            </div>
            <div class="col-12 col-md-3">
              <q-btn color="teal" icon="refresh" label="Load Preview" @click="loadPreview"
                :loading="store.previewLoading" class="full-width" />
            </div>
          </div>

          <div v-if="store.previewSnapshot" class="row q-col-gutter-md">
            <div class="col-12 col-md-6">
              <q-card class="bg-blue-1"><q-card-section>
                  <div class="text-caption">Snapshot BKU</div>
                  <div class="text-caption">Saldo Awal: {{ formatRp(store.previewSnapshot.snapshot_bku.SaldoAwal) }}
                  </div>
                  <div class="text-caption">Total Debet: {{ formatRp(store.previewSnapshot.snapshot_bku.TotalDebet) }}
                  </div>
                  <div class="text-caption">Total Kredit: {{ formatRp(store.previewSnapshot.snapshot_bku.TotalKredit) }}
                  </div>
                  <div class="text-subtitle1 text-weight-bold q-mt-sm">Saldo Akhir: {{
                    formatRp(store.previewSnapshot.snapshot_bku.SaldoAkhirBuku) }}</div>
                </q-card-section></q-card>
            </div>
            <div class="col-12 col-md-6">
              <q-card class="bg-purple-1"><q-card-section>
                  <div class="text-caption">Posisi Keuangan</div>
                  <div class="text-caption">Aset: {{ formatRp(store.previewSnapshot.snapshot_posisi.TotalAset) }}</div>
                  <div class="text-caption">Kewajiban: {{ formatRp(store.previewSnapshot.snapshot_posisi.TotalKewajiban)
                    }}
                  </div>
                  <div class="text-subtitle1 text-weight-bold q-mt-sm">Ekuitas: {{
                    formatRp(store.previewSnapshot.snapshot_posisi.TotalEkuitas) }}</div>
                </q-card-section></q-card>
            </div>
          </div>

          <q-btn v-if="store.previewSnapshot" color="teal" label="Lanjut" icon-right="arrow_forward" @click="step = 2"
            class="q-mt-md" />
        </q-step>

        <!-- Step 2: Cash Count -->
        <q-step :name="2" title="Cash Count" icon="payments" :done="step > 2">
          <div class="text-body2 q-mb-md">
            Hitung uang fisik di kas/brankas sekolah, lalu masukkan nominalnya.
            Saldo buku saat ini: <strong>{{ formatRp(store.previewSnapshot?.snapshot_bku.SaldoAkhirBuku) }}</strong>
          </div>
          <q-input v-model.number="wizard.saldo_fisik_kas" label="Saldo Fisik Kas *" outlined dense type="number"
            prefix="Rp" class="q-mb-md" style="max-width: 300px" />
          <div v-if="selisih !== 0" class="q-banner bg-warning text-white q-mb-md" rounded>
            <q-icon name="warning" class="q-mr-sm" />
            Selisih: {{ formatRp(selisih) }}. Akan dicatat sebagai selisih kas.
          </div>
          <q-input v-model="wizard.catatan" label="Catatan" outlined dense type="textarea" rows="2" class="q-mb-md" />
          <div class="q-gutter-sm">
            <q-btn flat label="Kembali" @click="step = 1" />
            <q-btn color="teal" label="Tutup Buku" icon="lock" @click="submitTutup" :loading="store.penutupanLoading"
              :disable="wizard.saldo_fisik_kas === null" />
          </div>
        </q-step>
      </q-stepper>
    </div>

    <!-- DETAIL MODE -->
    <div v-else-if="mode === 'detail' && detail">
      <div class="row items-center justify-between q-mb-md">
        <div class="text-h5">
          <q-icon name="receipt" class="q-mr-sm" />
          Penutupan {{ detail.periode_bulan }}/{{ detail.periode_tahun }}
          <q-chip dense :color="statusColor(detail.status)" text-color="white" class="q-ml-sm">
            {{ statusLabel(detail.status) }}
          </q-chip>
        </div>
        <q-btn flat icon="arrow_back" label="Kembali" @click="mode = 'list'" />
      </div>

      <q-card class="q-mb-md">
        <q-card-section>
          <div class="row q-col-gutter-md">
            <div class="col-6 col-md-3">
              <div class="text-caption text-grey-7">Saldo Awal</div>
              <div class="text-h6">{{ formatRp(detail.saldo_awal) }}</div>
            </div>
            <div class="col-6 col-md-3">
              <div class="text-caption text-grey-7">Debet</div>
              <div class="text-h6 text-positive">+{{ formatRp(detail.total_debet) }}</div>
            </div>
            <div class="col-6 col-md-3">
              <div class="text-caption text-grey-7">Kredit</div>
              <div class="text-h6 text-negative">-{{ formatRp(detail.total_kredit) }}</div>
            </div>
            <div class="col-6 col-md-3">
              <div class="text-caption text-grey-7">Saldo Akhir Buku</div>
              <div class="text-h6">{{ formatRp(detail.saldo_akhir_buku) }}</div>
            </div>
          </div>
        </q-card-section>
      </q-card>

      <q-card class="q-mb-md">
        <q-card-section>
          <div class="row q-col-gutter-md">
            <div class="col-6 col-md-3">
              <div class="text-caption text-grey-7">Saldo Fisik</div>
              <div class="text-h6">{{ formatRp(detail.saldo_fisik_kas) }}</div>
            </div>
            <div class="col-6 col-md-3">
              <div class="text-caption text-grey-7">Selisih</div>
              <div class="text-h6" :class="detail.selisih_kas !== 0 ? 'text-negative' : 'text-positive'">
                {{ formatRp(detail.selisih_kas) }}
              </div>
            </div>
            <div class="col-6 col-md-3">
              <div class="text-caption text-grey-7">Total Aset</div>
              <div class="text-h6">{{ formatRp(detail.total_aset) }}</div>
            </div>
            <div class="col-6 col-md-3">
              <div class="text-caption text-grey-7">Ekuitas</div>
              <div class="text-h6">{{ formatRp(detail.total_ekuitas) }}</div>
            </div>
          </div>
        </q-card-section>
      </q-card>

      <q-card v-if="detail.ttd_kepsek_hash">
        <q-card-section>
          <div class="text-subtitle2 q-mb-sm"><q-icon name="verified" color="positive" /> Dokumen Ditandatangani</div>
          <div class="text-caption">Signed by: <strong>{{ detail.ttd_kepsek_nama || detail.ttd_kepsek_by }}</strong>
          </div>
          <div class="text-caption">Signed at: {{ detail.ttd_kepsek_at }}</div>
          <div class="text-caption" style="word-break: break-all">
            Hash: <code style="font-size: 10px">{{ detail.ttd_kepsek_hash }}</code>
          </div>
          <q-btn flat dense icon="qr_code" label="Verifikasi Publik" color="primary"
            :href="`/#/keuangan/verify/${detail.ttd_kepsek_hash}`" target="_blank" class="q-mt-sm" />
        </q-card-section>
      </q-card>

      <q-card v-else class="bg-orange-1">
        <q-card-section>
          <q-icon name="warning" color="orange" />
          Belum ditandatangani Kepala Sekolah.
          <q-btn color="purple" icon="fingerprint" label="Tanda Tangan Sekarang" @click="openSign(detail)"
            class="q-ml-md" />
        </q-card-section>
      </q-card>
    </div>

    <!-- SIGN DIALOG -->
    <q-dialog v-model="signDialog" persistent>
      <q-card style="width: 400px; max-width: 90vw">
        <q-card-section>
          <div class="text-h6"><q-icon name="fingerprint" class="q-mr-sm" />Tanda Tangan Digital</div>
          <div class="text-caption text-grey-7 q-mt-xs">
            {{ currentSign?.periode_bulan }}/{{ currentSign?.periode_tahun }}
          </div>
        </q-card-section>
        <q-card-section>
          <q-input v-model="pinInput" label="PIN TTD (6 digit)" outlined dense type="password" maxlength="6"
            inputmode="numeric" autofocus>
            <template v-slot:prepend><q-icon name="lock" /></template>
          </q-input>
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn color="purple" icon="verified" label="Tanda Tangan" :loading="signing" @click="submitSign" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { usePenutupanStore } from '@/stores/finance/penutupanStore'
import { useTTDStore } from '@/stores/finance/ttdStore'

const $q = useQuasar()
const store = usePenutupanStore()
const ttdStore = useTTDStore()

const mode = ref('list')
const step = ref(1)
const detail = ref(null)
const signDialog = ref(false)
const currentSign = ref(null)
const pinInput = ref('')
const signing = ref(false)

const now = new Date()
const wizard = reactive({
  bulan: now.getMonth() + 1,
  tahun: now.getFullYear(),
  saldo_fisik_kas: null,
  catatan: '',
})

const bulanOptions = [
  'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
  'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
].map((label, i) => ({ label, value: i + 1 }))

const listColumns = [
  {
    name: 'periode', label: 'Periode', align: 'center',
    field: r => `${r.periode_bulan}/${r.periode_tahun}`
  },
  { name: 'saldo_akhir_buku', label: 'Saldo Akhir', align: 'right' },
  {
    name: 'saldo_fisik_kas', label: 'Fisik', field: 'saldo_fisik_kas', align: 'right',
    format: v => formatRp(v)
  },
  { name: 'selisih_kas', label: 'Selisih', align: 'right' },
  { name: 'status', label: 'Status', field: 'status', align: 'center' },
  { name: 'actions', label: 'Aksi', align: 'center' },
]

const selisih = computed(() => {
  if (!store.previewSnapshot || wizard.saldo_fisik_kas === null) return 0
  return wizard.saldo_fisik_kas - (store.previewSnapshot.snapshot_bku.SaldoAkhirBuku || 0)
})

function formatRp(n) {
  return n === null || n === undefined ? '-' : 'Rp ' + Number(n).toLocaleString('id-ID')
}
function statusColor(s) {
  return {
    CLOSED: 'blue', MENUNGGU_TTD_KEPALA: 'orange',
    CLOSED_FULLY_SIGNED: 'positive', DIBUKA_KEMBALI: 'grey'
  }[s] || 'grey'
}
function statusLabel(s) {
  return {
    CLOSED: 'Tertutup', MENUNGGU_TTD_KEPALA: 'Menunggu TTD',
    CLOSED_FULLY_SIGNED: 'Signed', DIBUKA_KEMBALI: 'Dibuka Kembali'
  }[s] || s
}

function openWizard() {
  wizard.bulan = now.getMonth() + 1
  wizard.tahun = now.getFullYear()
  wizard.saldo_fisik_kas = null
  wizard.catatan = ''
  step.value = 1
  mode.value = 'wizard'
}

async function loadPreview() {
  try {
    await store.loadPreview(wizard.bulan, wizard.tahun)
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || e.message })
  }
}

async function submitTutup() {
  try {
    await store.tutupBuku({
      periode_bulan: wizard.bulan,
      periode_tahun: wizard.tahun,
      saldo_fisik_kas: wizard.saldo_fisik_kas,
      catatan: wizard.catatan || undefined,
    })
    $q.notify({ type: 'positive', message: 'Buku kas berhasil ditutup' })
    mode.value = 'list'
    await loadList()
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || e.message, timeout: 5000 })
  }
}

async function openDetail(id) {
  try {
    const res = await store.getPenutupan(id)
    detail.value = res.data
    mode.value = 'detail'
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || e.message })
  }
}

function openSign(row) {
  currentSign.value = row
  pinInput.value = ''
  signDialog.value = true
}

async function submitSign() {
  if (!/^\d{6}$/.test(pinInput.value)) {
    $q.notify({ type: 'warning', message: 'PIN harus 6 digit angka' })
    return
  }
  signing.value = true
  try {
    await store.signPenutupan(currentSign.value.id, pinInput.value)
    $q.notify({ type: 'positive', message: 'Penutupan berhasil ditandatangani' })
    signDialog.value = false
    if (mode.value === 'detail') await openDetail(currentSign.value.id)
    await loadList()
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || e.message })
  } finally {
    signing.value = false
  }
}

async function loadList() {
  try { await store.listPenutupan() } catch {
    console.log('PenutupanPage..')
  }
}

onMounted(async () => {
  await ttdStore.getPinStatus()
  await loadList()
})
</script>

<template>
  <q-page class="q-pa-md">
    <q-linear-progress v-if="store.pengajuanLoading" indeterminate color="teal" class="q-mb-sm" />

    <div class="row items-center justify-between q-mb-md">
      <div class="text-h5"><q-icon name="approval" class="q-mr-sm" />Persetujuan Kepala Sekolah</div>
      <q-btn flat icon="refresh" label="Refresh" @click="loadData" />
    </div>

    <q-banner v-if="pinStatus && !pinStatus.is_setup" class="bg-warning text-white q-mb-md" rounded>
      <template v-slot:avatar><q-icon name="warning" /></template>
      PIN TTD belum disetup. <q-btn flat dense label="Setup PIN" to="/keuangan/pengaturan/pin-ttd" class="text-white" />
    </q-banner>

    <q-card v-if="pendingList.length === 0 && !store.pengajuanLoading">
      <q-card-section class="text-center text-grey-6 q-pa-xl">
        <q-icon name="check_circle" size="48px" class="q-mb-sm" color="positive" />
        <div>Tidak ada pengajuan menunggu persetujuan.</div>
      </q-card-section>
    </q-card>

    <q-card v-else>
      <q-table :rows="pendingList" :columns="columns" row-key="id" flat bordered :pagination="{ rowsPerPage: 50 }"
        :loading="store.pengajuanLoading">
        <template v-slot:body-cell-jumlah="props">
          <q-td :props="props" class="text-right text-weight-medium">{{ formatRp(props.row.jumlah) }}</q-td>
        </template>
        <template v-slot:body-cell-actions="props">
          <q-td :props="props" class="text-center">
            <q-btn flat dense round icon="visibility" color="info" @click="openDetail(props.row)" />
            <q-btn color="positive" icon="approval" label="Setujui" size="sm" dense :disable="!pinStatus?.is_setup"
              @click="openSignDialog(props.row, 'SETUJUI')" />
            <q-btn flat color="negative" icon="close" label="Tolak" size="sm" dense class="q-ml-xs"
              @click="openSignDialog(props.row, 'TOLAK')" />
          </q-td>
        </template>
      </q-table>
    </q-card>

    <!-- Dialog Detail -->
    <q-dialog v-model="detailDialog">
      <q-card style="width: 600px; max-width: 90vw">
        <q-card-section>
          <div class="text-h6">Detail Pengajuan</div>
        </q-card-section>
        <q-card-section v-if="currentRow">
          <q-list dense>
            <q-item><q-item-section>
                <div class="text-caption">Nomor</div>
                <div class="text-subtitle2">{{ currentRow.nomor_pengajuan }}</div>
              </q-item-section></q-item>
            <q-item><q-item-section>
                <div class="text-caption">Pemohon</div>
                <div>{{ currentRow.pemohon_nama }}</div>
              </q-item-section></q-item>
            <q-item><q-item-section>
                <div class="text-caption">Keperluan</div>
                <div>{{ currentRow.keperluan }}</div>
              </q-item-section></q-item>
            <q-item><q-item-section>
                <div class="text-caption">Jumlah</div>
                <div class="text-h6 text-primary">{{ formatRp(currentRow.jumlah) }}</div>
              </q-item-section></q-item>
            <q-item v-if="currentRow.catatan_bendahara">
              <q-item-section>
                <div class="text-caption">Catatan Bendahara</div>
                <div>{{ currentRow.catatan_bendahara }}</div>
              </q-item-section>
            </q-item>
          </q-list>
        </q-card-section>
        <q-card-actions align="right"><q-btn flat label="Tutup" v-close-popup /></q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Dialog Sign (PIN) -->
    <q-dialog v-model="signDialog" persistent>
      <q-card style="width: 420px; max-width: 90vw">
        <q-card-section>
          <div class="text-h6">
            <q-icon :name="signAksi === 'SETUJUI' ? 'approval' : 'cancel'"
              :color="signAksi === 'SETUJUI' ? 'positive' : 'negative'" class="q-mr-sm" />
            {{ signAksi === 'SETUJUI' ? 'Setujui & Tanda Tangan' : 'Tolak Pengajuan' }}
          </div>
          <div class="text-caption text-grey-7 q-mt-xs">
            {{ currentRow?.nomor_pengajuan }} — {{ formatRp(currentRow?.jumlah) }}
          </div>
        </q-card-section>
        <q-card-section>
          <div class="text-caption text-grey-7 q-mb-sm">
            <q-icon name="info" /> Aksi ini membutuhkan <strong>PIN TTD</strong> (6 digit).
            <span v-if="signAksi === 'SETUJUI'"> Sistem akan otomatis generate BKK.</span>
          </div>

          <q-input v-model="pinInput" label="PIN TTD (6 digit)" outlined dense type="password" maxlength="6"
            inputmode="numeric" :rules="[v => /^\d{6}$/.test(v) || 'PIN harus 6 digit angka']" autofocus>
            <template v-slot:prepend><q-icon name="lock" /></template>
          </q-input>

          <q-input v-model="catatan" label="Catatan (wajib jika tolak)" outlined dense type="textarea" rows="2"
            class="q-mt-sm" />

          <q-banner v-if="verifyResult && !verifyResult.valid" class="bg-negative text-white q-mt-sm" rounded>
            {{ verifyResult.message }} (sisa: {{ verifyResult.remaining_attempt }})
          </q-banner>
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" v-close-popup @click="resetSign" />
          <q-btn :color="signAksi === 'SETUJUI' ? 'positive' : 'negative'"
            :label="signAksi === 'SETUJUI' ? 'Setujui' : 'Tolak'" :loading="signing" @click="submitSign" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { usePengeluaranStore } from '@/stores/finance/pengeluaranStore'
import { useTTDStore } from '@/stores/finance/ttdStore'

const $q = useQuasar()
const store = usePengeluaranStore()
const ttdStore = useTTDStore()

const detailDialog = ref(false)
const signDialog = ref(false)
const currentRow = ref(null)
const signAksi = ref('SETUJUI')
const pinInput = ref('')
const catatan = ref('')
const signing = ref(false)
const verifyResult = ref(null)

const pendingList = computed(() =>
  store.pengajuan.filter(p => p.status === 'MENUNGGU_PERSETUJUAN_KEPALA')
)
const pinStatus = computed(() => ttdStore.pinStatus)

const columns = [
  { name: 'nomor_pengajuan', label: 'No. Pengajuan', field: 'nomor_pengajuan', align: 'left', sortable: true },
  { name: 'pemohon_nama', label: 'Pemohon', field: 'pemohon_nama', align: 'left' },
  {
    name: 'keperluan', label: 'Keperluan', field: 'keperluan', align: 'left',
    style: 'max-width: 250px; white-space: normal'
  },
  { name: 'jumlah', label: 'Jumlah', align: 'right', sortable: true },
  { name: 'tanggal', label: 'Tanggal', field: 'tanggal', align: 'center' },
  { name: 'actions', label: 'Aksi', align: 'center' },
]

function formatRp(n) {
  return n === null || n === undefined ? '-' : 'Rp ' + Number(n).toLocaleString('id-ID')
}

function openDetail(row) { currentRow.value = row; detailDialog.value = true }
function openSignDialog(row, aksi) {
  currentRow.value = row
  signAksi.value = aksi
  pinInput.value = ''
  catatan.value = ''
  verifyResult.value = null
  signDialog.value = true
}

function resetSign() {
  pinInput.value = ''
  catatan.value = ''
  verifyResult.value = null
}

async function submitSign() {
  if (!/^\d{6}$/.test(pinInput.value)) {
    $q.notify({ type: 'warning', message: 'PIN harus 6 digit angka' })
    return
  }
  if (signAksi.value === 'TOLAK' && !catatan.value.trim()) {
    $q.notify({ type: 'warning', message: 'Catatan wajib untuk penolakan' })
    return
  }
  signing.value = true
  verifyResult.value = null
  try {
    // 1. Verify PIN dulu
    const vr = await ttdStore.verifyPin(pinInput.value)
    verifyResult.value = vr
    if (!vr.valid) {
      $q.notify({ type: 'negative', message: vr.message || 'PIN tidak valid' })
      return
    }
    // 2. Kalau valid → call approval
    await store.approvalKepsek(currentRow.value.id, signAksi.value, catatan.value)
    $q.notify({ type: 'positive', message: signAksi.value === 'SETUJUI' ? 'Pengajuan disetujui, BKK ter-generate' : 'Pengajuan ditolak' })
    signDialog.value = false
    await loadData()
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || e.message })
  } finally {
    signing.value = false
  }
}

async function loadData() {
  try {
    await store.listPengajuan({ limit: 100 })
    await ttdStore.getPinStatus()
  } catch {
    console.log('PersetuJuanKepsek...')
  }
}

onMounted(loadData)
</script>

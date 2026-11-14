<template>
  <q-page class="q-pa-md">
    <q-linear-progress v-if="store.pengajuanLoading" indeterminate color="teal" class="q-mb-sm" />

    <div class="row items-center justify-between q-mb-md">
      <div class="text-h5"><q-icon name="fact_check" class="q-mr-sm" />Verifikasi Pengajuan</div>
      <q-btn flat icon="refresh" label="Refresh" @click="loadData" />
    </div>

    <q-card v-if="pendingList.length === 0 && !store.pengajuanLoading">
      <q-card-section class="text-center text-grey-6 q-pa-xl">
        <q-icon name="check_circle" size="48px" class="q-mb-sm" color="positive" />
        <div>Tidak ada pengajuan menunggu verifikasi.</div>
      </q-card-section>
    </q-card>

    <q-card v-else>
      <q-table :rows="pendingList" :columns="columns" row-key="id" flat bordered :pagination="{ rowsPerPage: 50 }"
        :loading="store.pengajuanLoading">
        <template v-slot:body-cell-jumlah="props">
          <q-td :props="props" class="text-right text-weight-medium">{{ formatRp(props.row.jumlah) }}</q-td>
        </template>
        <template v-slot:body-cell-prioritas="props">
          <q-td :props="props" class="text-center">
            <q-chip dense :color="prioritasColor(props.row.prioritas)" text-color="white">
              {{ props.row.prioritas }}
            </q-chip>
          </q-td>
        </template>
        <template v-slot:body-cell-actions="props">
          <q-td :props="props" class="text-center">
            <q-btn flat dense round icon="visibility" color="info" @click="openDetail(props.row)" />
            <q-btn flat dense round icon="check" color="positive" @click="openAksi(props.row, 'SETUJUI')" />
            <q-btn flat dense round icon="edit" color="warning" @click="openAksi(props.row, 'REVISI')" />
            <q-btn flat dense round icon="close" color="negative" @click="openAksi(props.row, 'TOLAK')" />
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
            <q-item><q-item-section>
                <div class="text-caption">Prioritas</div>
                <q-chip dense :color="prioritasColor(currentRow.prioritas)" text-color="white">
                  {{ currentRow.prioritas }}</q-chip>
              </q-item-section></q-item>
            <q-item v-if="currentRow.catatan_pemohon">
              <q-item-section>
                <div class="text-caption">Catatan Pemohon</div>
                <div>{{ currentRow.catatan_pemohon }}</div>
              </q-item-section>
            </q-item>
          </q-list>
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Tutup" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Dialog Aksi -->
    <q-dialog v-model="aksiDialog" persistent>
      <q-card style="width: 500px; max-width: 90vw">
        <q-card-section>
          <div class="text-h6">
            <q-icon :name="aksiIcon" :color="aksiColor" class="q-mr-sm" />
            {{ aksiLabel }}
          </div>
        </q-card-section>
        <q-card-section v-if="currentRow">
          <div class="text-caption text-grey-7 q-mb-sm">
            {{ currentRow.nomor_pengajuan }} — {{ formatRp(currentRow.jumlah) }}
          </div>
          <q-input v-model="aksiCatatan" label="Catatan *" outlined dense type="textarea" rows="3" :hint="aksiHint" />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn :color="aksiColor" :label="aksiLabel" :loading="store.pengajuanLoading" @click="submitAksi" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { usePengeluaranStore } from '@/stores/finance/pengeluaranStore'

const $q = useQuasar()
const store = usePengeluaranStore()

const detailDialog = ref(false)
const aksiDialog = ref(false)
const currentRow = ref(null)
const currentAksi = ref('SETUJUI')
const aksiCatatan = ref('')

const pendingList = computed(() =>
  store.pengajuan.filter(p => p.status === 'MENUNGGU_VERIFIKASI_BENDAHARA')
)

const aksiLabel = computed(() => ({
  SETUJUI: 'Setujui & Teruskan',
  TOLAK: 'Tolak Pengajuan',
  REVISI: 'Kembalikan untuk Revisi',
}[currentAksi.value] || 'Aksi'))

const aksiIcon = computed(() => ({
  SETUJUI: 'check_circle', TOLAK: 'cancel', REVISI: 'edit',
}[currentAksi.value] || 'help'))

const aksiColor = computed(() => ({
  SETUJUI: 'positive', TOLAK: 'negative', REVISI: 'warning',
}[currentAksi.value] || 'primary'))

const aksiHint = computed(() => ({
  SETUJUI: 'Catatan opsional',
  TOLAK: 'Wajib: jelaskan alasan penolakan',
  REVISI: 'Wajib: jelaskan apa yang perlu diperbaiki',
}[currentAksi.value] || ''))

const columns = [
  { name: 'nomor_pengajuan', label: 'No. Pengajuan', field: 'nomor_pengajuan', align: 'left', sortable: true },
  { name: 'pemohon_nama', label: 'Pemohon', field: 'pemohon_nama', align: 'left' },
  {
    name: 'keperluan', label: 'Keperluan', field: 'keperluan', align: 'left',
    style: 'max-width: 250px; white-space: normal'
  },
  { name: 'jumlah', label: 'Jumlah', align: 'right', sortable: true },
  { name: 'prioritas', label: 'Prioritas', field: 'prioritas', align: 'center' },
  { name: 'tanggal', label: 'Tanggal', field: 'tanggal', align: 'center' },
  { name: 'actions', label: 'Aksi', align: 'center' },
]

function formatRp(n) {
  return n === null || n === undefined ? '-' : 'Rp ' + Number(n).toLocaleString('id-ID')
}
function prioritasColor(p) {
  return { RENDAH: 'grey', SEDANG: 'blue', TINGGI: 'orange', DARURAT: 'negative' }[p] || 'grey'
}

function openDetail(row) {
  currentRow.value = row
  detailDialog.value = true
}

function openAksi(row, aksi) {
  currentRow.value = row
  currentAksi.value = aksi
  aksiCatatan.value = ''
  aksiDialog.value = true
}

async function submitAksi() {
  if (currentAksi.value !== 'SETUJUI' && !aksiCatatan.value.trim()) {
    $q.notify({ type: 'warning', message: 'Catatan wajib untuk aksi ini' })
    return
  }
  try {
    await store.verifikasiBendahara(currentRow.value.id, currentAksi.value, aksiCatatan.value)
    $q.notify({ type: 'positive', message: `Berhasil ${aksiLabel.value.toLowerCase()}` })
    aksiDialog.value = false
    await loadData()
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || e.message })
  }
}

async function loadData() {
  try {
    await store.listPengajuan({ limit: 100 })
  } catch {
    console.log('VerifikasiBendahara...')
  }
}

onMounted(loadData)
</script>

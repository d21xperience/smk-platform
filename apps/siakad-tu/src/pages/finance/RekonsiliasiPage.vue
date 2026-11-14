<template>
  <q-page class="q-pa-md">
    <q-linear-progress v-if="store.rekonLoading" indeterminate color="teal" class="q-mb-sm" />

    <!-- LIST MODE -->
    <div v-if="mode === 'list'">
      <div class="row items-center justify-between q-mb-md">
        <div class="text-h5"><q-icon name="sync_alt" class="q-mr-sm" />Rekonsiliasi Bank</div>
        <q-btn color="teal" icon="add" label="Rekonsiliasi Baru" @click="openCreate" />
      </div>

      <q-card v-if="store.rekonList.length === 0 && !store.rekonLoading">
        <q-card-section class="text-center text-grey-6 q-pa-xl">
          <q-icon name="sync_alt" size="48px" class="q-mb-sm" />
          <div>Belum ada rekonsiliasi.</div>
        </q-card-section>
      </q-card>

      <q-card v-else>
        <q-table :rows="store.rekonList" :columns="columns" row-key="id" flat bordered :pagination="{ rowsPerPage: 20 }"
          :loading="store.rekonLoading">
          <template v-slot:body-cell-status="props">
            <q-td :props="props" class="text-center">
              <q-chip dense :color="props.row.status === 'BALANCE' ? 'positive' : 'negative'" text-color="white">
                {{ props.row.status }}
              </q-chip>
            </q-td>
          </template>
          <template v-slot:body-cell-saldo_buku="props">
            <q-td :props="props" class="text-right">{{ formatRp(props.row.saldo_buku) }}</q-td>
          </template>
          <template v-slot:body-cell-saldo_bank="props">
            <q-td :props="props" class="text-right">{{ formatRp(props.row.saldo_bank) }}</q-td>
          </template>
          <template v-slot:body-cell-actions="props">
            <q-td :props="props" class="text-center">
              <q-btn flat dense round icon="visibility" @click="openDetail(props.row.id)" />
              <q-btn flat dense round icon="delete" color="negative" @click="confirmDelete(props.row)" />
            </q-td>
          </template>
        </q-table>
      </q-card>
    </div>

    <!-- FORM MODE -->
    <div v-else-if="mode === 'form'">
      <div class="row items-center justify-between q-mb-md">
        <div class="text-h5"><q-icon name="add_circle" class="q-mr-sm" />Rekonsiliasi Baru</div>
        <q-btn flat icon="arrow_back" label="Kembali" @click="mode = 'list'" />
      </div>

      <q-card class="q-mb-md">
        <q-card-section>
          <div class="row q-col-gutter-sm">
            <div class="col-12 col-md-4">
              <q-select v-model="form.rekening_bank_id" :options="rekeningOptions" label="Rekening Bank *" outlined
                dense emit-value map-options />
            </div>
            <div class="col-6 col-md-2">
              <q-select v-model="form.periode_bulan" :options="bulanOptions" label="Bulan *" outlined dense emit-value
                map-options />
            </div>
            <div class="col-6 col-md-2">
              <q-input v-model.number="form.periode_tahun" label="Tahun *" outlined dense type="number" />
            </div>
            <div class="col-12 col-md-4">
              <q-input v-model="form.catatan" label="Catatan" outlined dense />
            </div>
          </div>
        </q-card-section>
      </q-card>

      <q-card>
        <q-card-section>
          <div class="row items-center justify-between q-mb-sm">
            <div class="text-subtitle1 text-weight-medium">Mutasi Bank (Input Manual)</div>
            <q-btn color="green-7" icon="add" label="Tambah Mutasi" dense @click="addMutasi" />
          </div>

          <q-markup-table separator="cell" flat bordered dense>
            <thead>
              <tr class="bg-grey-3">
                <th style="width: 130px">Tanggal</th>
                <th>Keterangan</th>
                <th style="width: 120px">Debet</th>
                <th style="width: 120px">Kredit</th>
                <th style="width: 130px">Referensi</th>
                <th style="width: 50px"></th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(m, idx) in form.mutasi_items" :key="idx">
                <td><q-input v-model="m.tanggal" outlined dense type="date" /></td>
                <td><q-input v-model="m.keterangan" outlined dense /></td>
                <td><q-input v-model.number="m.debet" outlined dense type="number" min="0" /></td>
                <td><q-input v-model.number="m.kredit" outlined dense type="number" min="0" /></td>
                <td><q-input v-model="m.referensi" outlined dense placeholder="No. BKM/BKK" /></td>
                <td class="text-center">
                  <q-btn flat dense round icon="delete" color="negative" @click="form.mutasi_items.splice(idx, 1)" />
                </td>
              </tr>
              <tr v-if="form.mutasi_items.length === 0">
                <td colspan="6" class="text-center text-grey-6 q-pa-md">Klik "+ Tambah Mutasi" untuk mulai.</td>
              </tr>
            </tbody>
          </q-markup-table>

          <div class="row justify-end q-mt-md q-gutter-sm">
            <q-btn flat label="Batal" @click="mode = 'list'" />
            <q-btn color="teal" icon="sync_alt" label="Jalankan Rekonsiliasi" :loading="store.rekonLoading"
              :disable="!form.rekening_bank_id || form.mutasi_items.length === 0" @click="submit" />
          </div>
        </q-card-section>
      </q-card>
    </div>

    <!-- DETAIL MODE -->
    <div v-else-if="mode === 'detail' && detail">
      <div class="row items-center justify-between q-mb-md">
        <div class="text-h5">
          <q-icon name="checklist" class="q-mr-sm" />
          Rekonsiliasi {{ detail.periode_bulan }}/{{ detail.periode_tahun }}
          <q-chip dense :color="detail.status === 'BALANCE' ? 'positive' : 'negative'" text-color="white"
            class="q-ml-sm">{{
              detail.status }}</q-chip>
        </div>
        <q-btn flat icon="arrow_back" label="Kembali" @click="mode = 'list'" />
      </div>

      <div class="row q-col-gutter-md q-mb-md">
        <div class="col-6 col-md-3">
          <q-card class="bg-blue-1"><q-card-section>
              <div class="text-caption">Saldo Buku</div>
              <div class="text-h6 text-blue-10">{{ formatRp(detail.saldo_buku) }}</div>
            </q-card-section></q-card>
        </div>
        <div class="col-6 col-md-3">
          <q-card class="bg-purple-1"><q-card-section>
              <div class="text-caption">Saldo Bank</div>
              <div class="text-h6 text-purple-10">{{ formatRp(detail.saldo_bank) }}</div>
            </q-card-section></q-card>
        </div>
        <div class="col-6 col-md-3">
          <q-card :class="detail.status === 'BALANCE' ? 'bg-green-1' : 'bg-red-1'"><q-card-section>
              <div class="text-caption">Saldo Rekonsiliasi</div>
              <div class="text-h6">Bank: {{ formatRp(detail.saldo_rekonsiliasi_bank) }}</div>
              <div class="text-caption">Buku: {{ formatRp(detail.saldo_rekonsiliasi_buku) }}</div>
            </q-card-section></q-card>
        </div>
        <div class="col-6 col-md-3">
          <q-card><q-card-section>
              <div class="text-caption">Match Summary</div>
              <div class="text-h6">{{ detail.jumlah_matched }} ✓</div>
              <div class="text-caption">
                Unmatched: {{ detail.jumlah_unmatched_bank }} bank · {{ detail.jumlah_unmatched_buku }} buku
              </div>
            </q-card-section></q-card>
        </div>
      </div>

      <q-card v-if="detail.details && detail.details.length" class="q-mb-md">
        <q-card-section>
          <div class="text-subtitle1 text-weight-medium q-mb-sm">
            <q-icon name="check_circle" color="positive" /> Matched Pairs ({{ detail.details.length }})
          </div>
          <q-markup-table separator="cell" flat bordered dense>
            <thead>
              <tr class="bg-green-1">
                <th>Tanggal</th>
                <th>Keterangan Bank</th>
                <th class="text-right">Debet</th>
                <th>No. Bukti Buku</th>
                <th class="text-right">Debet Buku</th>
                <th>Δ Hari</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="d in detail.details" :key="d.id">
                <td>{{ d.mutasi_tanggal }}</td>
                <td>{{ d.mutasi_keterangan }}</td>
                <td class="text-right">{{ formatRp(d.mutasi_debet) }}</td>
                <td><strong>{{ d.no_bukti_bku }}</strong></td>
                <td class="text-right">{{ formatRp(d.mutasi_debet) }}</td>
                <td class="text-center">{{ d.selisih_hari }}</td>
              </tr>
            </tbody>
          </q-markup-table>
        </q-card-section>
      </q-card>

      <q-card v-if="detail.unmatched && detail.unmatched.length">
        <q-card-section>
          <div class="text-subtitle1 text-weight-medium q-mb-sm">
            <q-icon name="warning" color="orange" /> Unmatched Items ({{ detail.unmatched.length }})
          </div>
          <q-list separator>
            <q-item v-for="u in detail.unmatched" :key="u.id">
              <q-item-section avatar>
                <q-chip dense :color="u.tipe === 'BANK' ? 'purple' : 'orange'" text-color="white">
                  {{ u.tipe }}
                </q-chip>
              </q-item-section>
              <q-item-section>
                <q-item-label class="text-caption" style="font-family: monospace">{{ u.data_json }}</q-item-label>
              </q-item-section>
            </q-item>
          </q-list>
        </q-card-section>
      </q-card>
    </div>
  </q-page>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useRekonsiliasiStore } from '@/stores/finance/rekonsiliasiStore'
import { useMasterStore } from '@/stores/finance/masterStore'

const $q = useQuasar()
const store = useRekonsiliasiStore()
const masterStore = useMasterStore()

const mode = ref('list')
const detail = ref(null)

const now = new Date()
const form = reactive({
  rekening_bank_id: '',
  periode_bulan: now.getMonth() + 1,
  periode_tahun: now.getFullYear(),
  catatan: '',
  mutasi_items: [],
})

const bulanOptions = [
  'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
  'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
].map((label, i) => ({ label, value: i + 1 }))

const rekeningOptions = computed(() =>
  masterStore.rekeningBank.items.map(r => ({ label: `${r.nama_bank} — ${r.nomor_rekening}`, value: r.id }))
)

const columns = [
  { name: 'rekening_bank_nama', label: 'Rekening', field: 'rekening_bank_nama', align: 'left' },
  {
    name: 'periode', label: 'Periode', align: 'center',
    field: r => `${r.periode_bulan}/${r.periode_tahun}`
  },
  { name: 'saldo_buku', label: 'Saldo Buku', align: 'right' },
  { name: 'saldo_bank', label: 'Saldo Bank', align: 'right' },
  { name: 'jumlah_matched', label: 'Matched', field: 'jumlah_matched', align: 'center' },
  { name: 'status', label: 'Status', field: 'status', align: 'center' },
  { name: 'actions', label: 'Aksi', align: 'center' },
]

function formatRp(n) {
  return n === null || n === undefined ? '-' : 'Rp ' + Number(n).toLocaleString('id-ID')
}

function openCreate() {
  form.rekening_bank_id = rekeningOptions.value[0]?.value || ''
  form.periode_bulan = now.getMonth() + 1
  form.periode_tahun = now.getFullYear()
  form.catatan = ''
  form.mutasi_items = []
  mode.value = 'form'
}

function addMutasi() {
  form.mutasi_items.push({
    tanggal: new Date().toISOString().split('T')[0],
    keterangan: '', debet: 0, kredit: 0, referensi: '',
  })
}

async function submit() {
  try {
    const payload = {
      rekening_bank_id: form.rekening_bank_id,
      periode_bulan: form.periode_bulan,
      periode_tahun: form.periode_tahun,
      catatan: form.catatan || undefined,
      mutasi_items: form.mutasi_items.map(m => ({
        tanggal: m.tanggal, keterangan: m.keterangan,
        debet: Number(m.debet) || 0, kredit: Number(m.kredit) || 0,
        referensi: m.referensi || undefined,
      })),
    }
    const res = await store.jalankanRekonsiliasi(payload)
    $q.notify({ type: 'positive', message: `Rekonsiliasi selesai: ${res.data.status}` })
    await openDetail(res.data.id)
    await loadList()
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || e.message, timeout: 5000 })
  }
}

async function openDetail(id) {
  try {
    const res = await store.getRekonsiliasi(id)
    detail.value = res.data
    mode.value = 'detail'
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || e.message })
  }
}

function confirmDelete(row) {
  $q.dialog({
    title: 'Hapus Rekonsiliasi',
    message: `Hapus rekonsiliasi ${row.periode_bulan}/${row.periode_tahun}?`,
    cancel: true, persistent: true,
  }).onOk(async () => {
    try {
      await store.deleteRekonsiliasi(row.id)
      $q.notify({ type: 'positive', message: 'Dihapus' })
      await loadList()
    } catch (e) {
      $q.notify({ type: 'negative', message: e.response?.data?.error || e.message })
    }
  })
}

async function loadList() {
  try { await store.listRekonsiliasi() } catch {
    console.log('RekonsiliasiPage...')
  }
}

onMounted(async () => {
  if (masterStore.rekeningBank.items.length === 0) {
    await masterStore.rekeningBank.list({ limit: 50 })
  }
  await loadList()
})
</script>

<template>
  <q-page class="q-pa-md">
    <q-linear-progress v-if="store.rabLoading" indeterminate color="teal" class="q-mb-sm" />

    <!-- ============================================================ -->
    <!-- MODE: LIST                                                    -->
    <!-- ============================================================ -->
    <div v-if="mode === 'list'">
      <div class="row items-center justify-between q-mb-md">
        <div class="text-h5"><q-icon name="edit_document" class="q-mr-sm" />RAB Tahunan</div>
        <q-btn color="teal" icon="add" label="Buat RAB Baru" @click="openCreate" />
      </div>

      <q-card>
        <q-table :rows="store.rabList" :columns="listColumns" row-key="id" flat bordered
          :pagination="{ rowsPerPage: 20 }" :loading="store.rabLoading">
          <template v-slot:body-cell-status="props">
            <q-td :props="props" class="text-center">
              <q-chip dense :color="statusColor(props.row.status)" text-color="white">
                {{ props.row.status }}
              </q-chip>
            </q-td>
          </template>
          <template v-slot:body-cell-total_pendapatan="props">
            <q-td :props="props" class="text-right">{{ formatRp(props.row.total_pendapatan) }}</q-td>
          </template>
          <template v-slot:body-cell-total_belanja="props">
            <q-td :props="props" class="text-right">{{ formatRp(props.row.total_belanja) }}</q-td>
          </template>
          <template v-slot:body-cell-surplus_defisit="props">
            <q-td :props="props" class="text-right">
              <span :class="props.row.surplus_defisit < 0 ? 'text-negative' : 'text-positive'">
                {{ formatRp(props.row.surplus_defisit) }}
              </span>
            </q-td>
          </template>
          <template v-slot:body-cell-actions="props">
            <q-td :props="props" class="text-center">
              <q-btn flat dense round icon="visibility" @click="openEdit(props.row)" />
              <q-btn v-if="props.row.status === 'DRAFT' || props.row.status === 'DIREVISI'" flat dense round icon="edit"
                color="primary" @click="openEdit(props.row)" />
              <q-btn v-if="props.row.status === 'DRAFT' || props.row.status === 'DIAJUKAN'" flat dense round
                icon="verified" color="positive" @click="confirmSahkan(props.row)" />
              <q-btn v-if="props.row.status === 'DRAFT' || props.row.status === 'DIREVISI'" flat dense round
                icon="delete" color="negative" @click="confirmDelete(props.row)" />
            </q-td>
          </template>
        </q-table>
      </q-card>
    </div>

    <!-- ============================================================ -->
    <!-- MODE: FORM                                                    -->
    <!-- ============================================================ -->
    <div v-else>
      <div class="row items-center justify-between q-mb-md">
        <div class="text-h5">
          <q-icon :name="editingId ? 'edit' : 'add_circle'" class="q-mr-sm" />
          {{ editingId ? 'Edit' : 'Buat' }} RAB
          <span v-if="form.tahun_anggaran" class="text-h6 text-grey-7 q-ml-sm">
            Tahun {{ form.tahun_anggaran }}
          </span>
        </div>
        <q-btn flat icon="arrow_back" label="Kembali ke List" @click="backToList" />
      </div>

      <q-card class="q-mb-md">
        <q-card-section>
          <div class="row q-col-gutter-sm">
            <div class="col-12 col-sm-3 col-md-2">
              <q-input v-model.number="form.tahun_anggaran" label="Tahun Anggaran *" outlined dense type="number"
                min="2000" max="2100" :disable="!!editingId" />
            </div>
            <div class="col-12 col-md-6">
              <q-input v-model="form.catatan" label="Catatan" outlined dense />
            </div>
          </div>
        </q-card-section>
      </q-card>

      <!-- Items Table -->
      <q-card>
        <q-card-section>
          <div class="row items-center justify-between q-mb-sm">
            <div class="text-subtitle1 text-weight-medium">Item RAB</div>
            <div class="q-gutter-sm">
              <q-btn color="green-7" icon="add" label="+ Pendapatan" dense @click="addItem('PENDAPATAN')" />
              <q-btn color="orange-7" icon="add" label="+ Belanja" dense @click="addItem('BELANJA')" />
            </div>
          </div>

          <q-markup-table separator="cell" flat bordered dense>
            <thead>
              <tr class="bg-grey-3">
                <th style="width: 100px">Tipe</th>
                <th style="width: 180px">Kode Rekening</th>
                <th>Uraian</th>
                <th style="width: 90px">Volume</th>
                <th style="width: 80px">Satuan</th>
                <th style="width: 130px">Harga Satuan</th>
                <th style="width: 130px">Subtotal</th>
                <th style="width: 110px">Sumber</th>
                <th style="width: 50px"></th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(item, idx) in form.items" :key="idx">
                <td class="text-center">
                  <q-chip dense :color="item.tipe === 'PENDAPATAN' ? 'green-6' : 'orange-6'" text-color="white">
                    {{ item.tipe === 'PENDAPATAN' ? 'Pendapatan' : 'Belanja' }}
                  </q-chip>
                </td>
                <td>
                  <q-select v-model="item.kode_rekening_id" :options="kodeRekeningOptions(item.tipe)" outlined dense
                    emit-value map-options use-input input-debounce="0" @filter="onFilterKode"
                    :loading="masterStore.kodeRekening.loading" />
                </td>
                <td>
                  <q-input v-model="item.uraian" outlined dense />
                </td>
                <td>
                  <q-input v-model.number="item.volume" outlined dense type="number" min="0.01" step="0.01" />
                </td>
                <td>
                  <q-input v-model="item.satuan" outlined dense placeholder="unit" />
                </td>
                <td>
                  <q-input v-model.number="item.harga_satuan" outlined dense type="number" min="1" />
                </td>
                <td class="text-right">
                  <span class="text-weight-medium">{{ formatRp(subtotal(item)) }}</span>
                </td>
                <td>
                  <q-select v-model="item.sumber_dana" :options="sumberOptions" outlined dense emit-value map-options
                    clearable />
                </td>
                <td class="text-center">
                  <q-btn flat dense round icon="delete" color="negative" @click="removeItem(idx)" />
                </td>
              </tr>
              <tr v-if="form.items.length === 0">
                <td colspan="9" class="text-center text-grey-6 q-pa-md">
                  Belum ada item. Klik "+ Pendapatan" atau "+ Belanja" untuk menambahkan.
                </td>
              </tr>
            </tbody>
          </q-markup-table>
        </q-card-section>

        <!-- Footer Totals -->
        <q-separator />
        <q-card-section>
          <div class="row q-col-gutter-md">
            <div class="col-12 col-md-4">
              <q-card flat class="bg-green-1">
                <q-card-section class="q-pa-sm">
                  <div class="text-caption text-grey-7">Total Pendapatan</div>
                  <div class="text-h6 text-green-10">{{ formatRp(totalPendapatan) }}</div>
                </q-card-section>
              </q-card>
            </div>
            <div class="col-12 col-md-4">
              <q-card flat class="bg-orange-1">
                <q-card-section class="q-pa-sm">
                  <div class="text-caption text-grey-7">Total Belanja</div>
                  <div class="text-h6 text-orange-10">{{ formatRp(totalBelanja) }}</div>
                </q-card-section>
              </q-card>
            </div>
            <div class="col-12 col-md-4">
              <q-card flat :class="surplus >= 0 ? 'bg-blue-1' : 'bg-red-1'">
                <q-card-section class="q-pa-sm">
                  <div class="text-caption text-grey-7">Surplus / Defisit</div>
                  <div class="text-h6" :class="surplus >= 0 ? 'text-blue-10' : 'text-negative'">
                    {{ formatRp(surplus) }}
                  </div>
                </q-card-section>
              </q-card>
            </div>
          </div>

          <q-banner v-if="totalBelanja > totalPendapatan && form.items.length > 0"
            class="bg-negative text-white q-mt-md" rounded>
            <template v-slot:avatar><q-icon name="error" /></template>
            Total belanja melebihi pendapatan — tidak dapat disimpan.
          </q-banner>

          <div class="row justify-end q-mt-md q-gutter-sm">
            <q-btn flat label="Batal" @click="backToList" />
            <q-btn color="teal" icon="save" :label="editingId ? 'Simpan Perubahan' : 'Simpan sebagai DRAFT'"
              :loading="store.rabLoading" :disable="totalBelanja > totalPendapatan || form.items.length === 0"
              @click="submit" />
          </div>
        </q-card-section>
      </q-card>
    </div>
  </q-page>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useRabStore } from '@/stores/finance/rabStore'
import { useMasterStore } from '@/stores/finance/masterStore'

const $q = useQuasar()
const store = useRabStore()
const masterStore = useMasterStore()

const mode = ref('list')
const editingId = ref(null)

const form = reactive({
  tahun_anggaran: new Date().getFullYear(),
  catatan: '',
  items: [],
})

const listColumns = [
  { name: 'tahun_anggaran', label: 'Tahun', field: 'tahun_anggaran', align: 'center', sortable: true },
  { name: 'total_pendapatan', label: 'Pendapatan', align: 'right', sortable: true },
  { name: 'total_belanja', label: 'Belanja', align: 'right', sortable: true },
  { name: 'surplus_defisit', label: 'Surplus', align: 'right', sortable: true },
  { name: 'status', label: 'Status', field: 'status', align: 'center' },
  { name: 'tanggal_pengesahan', label: 'Disahkan', field: 'tanggal_pengesahan', align: 'center' },
  { name: 'actions', label: 'Aksi', align: 'center' },
]

const sumberOptions = [
  { label: 'BOS', value: 'BOS' },
  { label: 'BOP', value: 'BOP' },
  { label: 'Hibah', value: 'HIBAH' },
  { label: 'Mandiri', value: 'MANDIRI' },
]

const kodeFilter = ref('')

function onFilterKode(val, update) {
  update(() => { kodeFilter.value = val })
}

function kodeRekeningOptions(tipe) {
  return masterStore.kodeRekening.items
    .filter(k => k.tipe === tipe && k.is_aktif)
    .filter(k => !kodeFilter.value || k.kode.toLowerCase().includes(kodeFilter.value.toLowerCase()) ||
      k.nama.toLowerCase().includes(kodeFilter.value.toLowerCase()))
    .map(k => ({ label: `${k.kode} — ${k.nama}`, value: k.id }))
}

function formatRp(n) {
  return n === null || n === undefined ? '-' : 'Rp ' + Number(n).toLocaleString('id-ID')
}
function statusColor(s) {
  return { DRAFT: 'grey', DIAJUKAN: 'blue', DISAHKAN: 'positive', DIREVISI: 'warning' }[s] || 'grey'
}
function subtotal(item) {
  return (Number(item.volume) || 0) * (Number(item.harga_satuan) || 0)
}

const totalPendapatan = computed(() =>
  form.items.filter(i => i.tipe === 'PENDAPATAN').reduce((s, i) => s + subtotal(i), 0)
)
const totalBelanja = computed(() =>
  form.items.filter(i => i.tipe === 'BELANJA').reduce((s, i) => s + subtotal(i), 0)
)
const surplus = computed(() => totalPendapatan.value - totalBelanja.value)

function addItem(tipe) {
  form.items.push({
    tipe,
    kode_rekening_id: '',
    uraian: '',
    volume: 1,
    satuan: 'unit',
    harga_satuan: 0,
    sumber_dana: '',
  })
}
function removeItem(idx) {
  form.items.splice(idx, 1)
}

function openCreate() {
  editingId.value = null
  form.tahun_anggaran = new Date().getFullYear()
  form.catatan = ''
  form.items = []
  mode.value = 'form'
}

async function openEdit(row) {
  try {
    const res = await store.getRab(row.id)
    const detail = res.data
    editingId.value = detail.id
    form.tahun_anggaran = detail.tahun_anggaran
    form.catatan = detail.catatan || ''
    form.items = (detail.details || []).map(d => ({
      tipe: d.tipe,
      kode_rekening_id: d.kode_rekening_id,
      uraian: d.uraian,
      volume: d.volume,
      satuan: d.satuan,
      harga_satuan: d.harga_satuan,
      sumber_dana: d.sumber_dana || '',
    }))
    mode.value = 'form'
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || e.message })
  }
}

function backToList() {
  mode.value = 'list'
  editingId.value = null
}

async function submit() {
  if (form.items.length === 0) {
    $q.notify({ type: 'warning', message: 'Item RAB tidak boleh kosong' })
    return
  }
  if (totalBelanja.value > totalPendapatan.value) {
    $q.notify({ type: 'warning', message: 'Total belanja melebihi pendapatan' })
    return
  }
  try {
    const payload = {
      academic_period_id: store.periodID,
      tahun_anggaran: form.tahun_anggaran,
      catatan: form.catatan || undefined,
      items: form.items.map(i => ({
        kode_rekening_id: i.kode_rekening_id,
        uraian: i.uraian,
        volume: i.volume,
        satuan: i.satuan || 'unit',
        harga_satuan: i.harga_satuan,
        sumber_dana: i.sumber_dana || undefined,
      })),
    }
    if (editingId.value) {
      await store.updateRab(editingId.value, payload)
      $q.notify({ type: 'positive', message: 'RAB berhasil diperbarui' })
    } else {
      await store.createRab(payload)
      $q.notify({ type: 'positive', message: 'RAB berhasil disimpan sebagai DRAFT' })
    }
    backToList()
    await loadList()
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || e.message, timeout: 5000 })
  }
}

function confirmSahkan(row) {
  $q.dialog({
    title: 'Sahkan RAB',
    message: `Sahkan RAB tahun ${row.tahun_anggaran}? Setelah disahkan tidak dapat diubah tanpa revisi.`,
    cancel: true, persistent: true,
  }).onOk(async () => {
    try {
      await store.sahkanRab(row.id, 'Disahkan melalui UI')
      $q.notify({ type: 'positive', message: 'RAB berhasil disahkan' })
      await loadList()
    } catch (e) {
      $q.notify({ type: 'negative', message: e.response?.data?.error || e.message })
    }
  })
}

function confirmDelete(row) {
  $q.dialog({
    title: 'Hapus RAB',
    message: `Hapus RAB tahun ${row.tahun_anggaran}?`,
    cancel: true, persistent: true,
  }).onOk(async () => {
    try {
      await store.deleteRab(row.id)
      $q.notify({ type: 'positive', message: 'RAB berhasil dihapus' })
      await loadList()
    } catch (e) {
      $q.notify({ type: 'negative', message: e.response?.data?.error || e.message })
    }
  })
}

async function loadList() {
  try { await store.listRab() } catch {
    console.log('RabFormPage...')
  }
}

onMounted(async () => {
  if (masterStore.kodeRekening.items.length === 0) {
    await masterStore.kodeRekening.list({ limit: 100 })
  }
  await loadList()
})
</script>

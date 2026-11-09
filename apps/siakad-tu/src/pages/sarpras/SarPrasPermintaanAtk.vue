<template>
  <q-page class="q-pa-md">
    <!-- PANEL RINGKASAN -->
    <div class="row q-col-gutter-sm q-mb-md items-center justify-between">
      <div class="col-12 col-sm-4">
        <q-card flat bordered class="bg-amber-1 text-amber-9">
          <q-card-section class="row items-center no-wrap q-py-sm">
            <div>
              <div class="text-caption text-grey-7">Menunggu Persetujuan</div>
              <div class="text-h6 text-weight-bold">{{ totalPending }} Pengajuan</div>
            </div>
            <q-space />
            <q-icon name="pending_actions" size="sm" />
          </q-card-section>
        </q-card>
      </div>

      <div class="col-12 col-sm-auto">
        <q-btn color="primary" icon="add" label="Ajukan Kebutuhan Barang" @click="bukaDialog()" no-caps />
      </div>
    </div>

    <!-- TABEL PENGAJUAN -->
    <q-table title="Daftar Pengajuan Kebutuhan Barang" :rows="store.requests" :columns="columns" row-key="id" flat bordered :loading="store.isLoading" :filter="filterText">
      <template v-slot:top-right>
        <q-input borderless dense debounce="300" v-model="filterText" placeholder="Cari nomor/item/pengaju...">
          <template v-slot:append><q-icon name="search" /></template>
        </q-input>
      </template>

      <!-- Badge Prioritas -->
      <template v-slot:body-cell-priority="props">
        <q-td :props="props" class="text-center">
          <q-badge :color="getPriorityColor(props.value)">{{ props.value }}</q-badge>
        </q-td>
      </template>

      <!-- Badge Status -->
      <template v-slot:body-cell-status="props">
        <q-td :props="props" class="text-center">
          <q-chip :color="getStatusColor(props.value)" text-color="white" dense class="text-weight-bold">
            {{ props.value }}
          </q-chip>
        </q-td>
      </template>

      <!-- Kolom Aksi -->
      <template v-slot:body-cell-aksi="props">
        <q-td :props="props" class="text-center">
          <div v-if="props.row.status === 'PENDING'" class="row justify-center q-gutter-xs">
            <q-btn dense round size="sm" color="positive" icon="check" @click="prosesApproval(props.row, 'APPROVED')">
              <q-tooltip>Setujui</q-tooltip>
            </q-btn>
            <q-btn dense round size="sm" color="negative" icon="close" @click="bukaTolakDialog(props.row)">
              <q-tooltip>Tolak</q-tooltip>
            </q-btn>
            <q-btn dense round size="sm" color="primary" icon="edit" @click="bukaDialog(props.row)">
              <q-tooltip>Edit</q-tooltip>
            </q-btn>
            <q-btn dense round size="sm" color="grey" icon="delete" @click="konfirmasiHapus(props.row)">
              <q-tooltip>Hapus</q-tooltip>
            </q-btn>
          </div>
          <div v-else-if="props.row.status === 'APPROVED'">
            <q-btn dense size="sm" color="blue" label="Tandai Diadakan" @click="prosesApproval(props.row, 'PROCURED')" no-caps />
          </div>
          <span v-else class="text-caption text-grey-5">Selesai</span>
        </q-td>
      </template>
    </q-table>

    <!-- DIALOG FORM PENGAJUAN -->
    <q-dialog v-model="dialogForm" persistent>
      <q-card style="width: 500px; max-width: 90vw;">
        <q-card-section class="bg-primary text-white">
          <div class="text-h6"><q-icon name="post_add" /> {{ isEdit ? 'Edit' : 'Formulir' }} Pengajuan Kebutuhan</div>
        </q-card-section>

        <q-card-section class="q-gutter-md q-pt-md">
          <q-select 
            outlined dense 
            v-model="form.requester_name" 
            :options="store.ptkList" 
            option-label="name" 
            option-value="name"
            label="Nama Guru / Staff Pemohon *" 
            emit-value map-options
            :disable="isEdit" 
          />
          <q-select 
            outlined dense 
            v-model="form.item_name" 
            :options="store.catalogList" 
            option-label="name" 
            option-value="name"
            label="Nama Barang Kebutuhan *" 
            emit-value map-options
          />
          <q-input outlined dense v-model="form.specification" label="Spesifikasi (Opsional)" />
          
          <div class="row q-col-gutter-sm">
            <div class="col-6">
              <q-input outlined dense type="number" v-model.number="form.quantity" label="Jumlah *" />
            </div>
            <div class="col-6">
              <q-input outlined dense v-model="form.unit" label="Satuan (Cth: Pcs, Rim) *" />
            </div>
          </div>
          <div class="row q-col-gutter-sm">
            <div class="col-6">
              <q-input outlined dense type="number" v-model.number="form.estimated_price" label="Estimasi Harga (Rp)" />
            </div>
            <div class="col-6">
              <q-select outlined dense v-model="form.priority" :options="['LOW', 'NORMAL', 'HIGH', 'URGENT']" label="Prioritas *" />
            </div>
          </div>
          <q-input outlined dense type="textarea" v-model="form.justification" label="Alasan Kebutuhan *" />
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn color="primary" :label="isEdit ? 'Simpan Perubahan' : 'Kirim Pengajuan'" @click="submitForm" :loading="isSubmitting" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- DIALOG TOLAK -->
    <q-dialog v-model="dialogTolak">
      <q-card style="width: 400px;">
        <q-card-section class="text-h6">Tolak Pengajuan</q-card-section>
        <q-card-section>
          <q-input outlined dense type="textarea" v-model="rejectionReason" label="Alasan Penolakan (Wajib) *" />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn color="negative" label="Tolak Pengajuan" @click="submitPenolakan" :loading="isSubmitting" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useQuasar } from 'quasar'
import { useProcurementStore } from '@/stores/sarpras/procurement'
import { useContextStore } from '@/stores/contextStore' // <-- IMPORT CONTEXT STORE

const $q = useQuasar()
const store = useProcurementStore()
const contextStore = useContextStore() // <-- INISIALISASI

const filterText = ref('')
const dialogForm = ref(false)
const dialogTolak = ref(false)
const isEdit = ref(false)
const isSubmitting = ref(false)
const rejectionReason = ref('')
const selectedRequest = ref(null)

const form = ref({ id: '', requester_name: '', item_name: '', specification: '', quantity: 1, unit: '', estimated_price: 0, priority: 'NORMAL', justification: '' })

const columns = [
  { name: 'request_number', label: 'No. Pengajuan', align: 'left', field: 'request_number', sortable: true },
  { name: 'created_at', label: 'Tanggal', align: 'center', field: 'created_at', format: val => new Date(val).toLocaleDateString('id-ID') },
  { name: 'requester_name', label: 'Pemohon', align: 'left', field: 'requester_name', sortable: true },
  { name: 'item_name', label: 'Barang', align: 'left', field: 'item_name' },
  { name: 'quantity', label: 'Vol', align: 'center', field: 'quantity', format: (val, row) => `${val} ${row.unit}` },
  { name: 'priority', label: 'Prioritas', align: 'center', field: 'priority' },
  { name: 'status', label: 'Status', align: 'center', field: 'status' },
  { name: 'aksi', label: 'Tindakan', align: 'center' }
]

const totalPending = computed(() => store.requests.filter(r => r.status === 'PENDING').length)

const getPriorityColor = (p) => ({ 'LOW': 'grey', 'NORMAL': 'blue', 'HIGH': 'orange', 'URGENT': 'red' }[p] || 'grey')
const getStatusColor = (s) => ({ 'PENDING': 'amber', 'APPROVED': 'green', 'REJECTED': 'red', 'PROCURED': 'blue' }[s] || 'grey')

onMounted(() => {
  store.fetchRequests()
  store.fetchDropdownData()
})

// WATCHER UNTUK REAKTIVITAS PERUBAHAN SEMESTER
watch(() => contextStore.selectedSemesterId, (newVal) => {
  if (newVal) {
    store.fetchRequests()
  }
})

function bukaDialog(row = null) {
  if (row) {
    isEdit.value = true
    form.value = { ...row }
  } else {
    isEdit.value = false
    form.value = { id: '', requester_name: '', item_name: '', specification: '', quantity: 1, unit: '', estimated_price: 0, priority: 'NORMAL', justification: '' }
  }
  dialogForm.value = true
}

async function submitForm() {
  if (!form.value.requester_name || !form.value.item_name || !form.value.unit || !form.value.justification) {
    $q.notify({ type: 'negative', message: 'Mohon lengkapi kolom bertanda *' })
    return
  }
  isSubmitting.value = true
  try {
    if (isEdit.value) {
      await store.updateRequest(form.value.id, form.value)
      $q.notify({ type: 'positive', message: 'Pengajuan berhasil diupdate' })
    } else {
      const res = await store.createRequest(form.value)
      $q.notify({ type: 'positive', message: res.message })
    }
    dialogForm.value = false
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || 'Gagal menyimpan pengajuan' })
  } finally {
    isSubmitting.value = false
  }
}

function konfirmasiHapus(row) {
  $q.dialog({
    title: 'Konfirmasi Hapus',
    message: `Hapus pengajuan ${row.request_number}?`,
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await store.deleteRequest(row.id)
      $q.notify({ type: 'positive', message: 'Pengajuan berhasil dihapus' })
    } catch (e) {
      $q.notify({ type: 'negative', message: e.response?.data?.error || 'Gagal menghapus' })
    }
  })
}

async function prosesApproval(row, newStatus) {
  try {
    await store.updateStatus(row.id, newStatus)
    $q.notify({ type: 'positive', message: `Status diubah menjadi ${newStatus}` })
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || 'Gagal update status' })
  }
}

function bukaTolakDialog(row) {
  selectedRequest.value = row
  rejectionReason.value = ''
  dialogTolak.value = true
}

async function submitPenolakan() {
  if (!rejectionReason.value) {
    $q.notify({ type: 'negative', message: 'Alasan penolakan wajib diisi' })
    return
  }
  isSubmitting.value = true
  try {
    await store.updateStatus(selectedRequest.value.id, 'REJECTED', rejectionReason.value)
    $q.notify({ type: 'info', message: 'Pengajuan ditolak' })
    dialogTolak.value = false
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || 'Gagal menolak' })
  } finally {
    isSubmitting.value = false
  }
}
</script>

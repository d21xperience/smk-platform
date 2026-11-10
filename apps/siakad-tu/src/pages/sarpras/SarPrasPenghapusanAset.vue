<template>
  <q-page class="q-pa-md">
    <div class="row q-mb-md items-center justify-between">
      <div class="text-h6 text-weight-bold">Usulan Penghapusan Aset</div>
      <q-btn color="primary" icon="add" label="Buat Usulan Baru" @click="bukaForm" no-caps />
    </div>

    <q-table :rows="store.requests" :columns="columns" row-key="id" flat bordered :loading="store.isLoading">
      <template v-slot:body-cell-status="props">
        <q-td :props="props" class="text-center">
          <q-badge :color="getStatusColor(props.value)">{{ props.value }}</q-badge>
        </q-td>
      </template>

      <template v-slot:body-cell-aksi="props">
        <q-td :props="props" class="text-center q-gutter-x-xs">
          <!-- Tombol Verifikasi (jika DRAFT) -->
          <q-btn v-if="props.row.status === 'DRAFT'" dense round size="sm" color="amber" icon="fact_check" @click="bukaAksi(props.row, 'verify')">
            <q-tooltip>Verifikasi KTU</q-tooltip>
          </q-btn>
          
          <!-- Tombol Approve (jika sudah diverifikasi dan belum final) -->
          <q-btn v-if="['VERIFIED_KTU', 'APPROVED_STEP_1'].includes(props.row.status)" dense round size="sm" color="green" icon="check_circle" @click="bukaAksi(props.row, 'approve')">
            <q-tooltip>Approval</q-tooltip>
          </q-btn>

          <!-- Tombol Eksekusi (jika sudah full approved) -->
          <q-btn v-if="props.row.status === 'APPROVED'" dense round size="sm" color="blue" icon="delete_forever" @click="bukaAksi(props.row, 'execute')">
            <q-tooltip>Eksekusi Penghapusan</q-tooltip>
          </q-btn>

          <!-- Tombol PDF (jika sudah dieksekusi) -->
          <q-btn v-if="props.row.status === 'EXECUTED'" dense round size="sm" color="primary" icon="picture_as_pdf" @click="cetakBAP(props.row)">
            <q-tooltip>Cetak BAP</q-tooltip>
          </q-btn>
        </q-td>
      </template>
    </q-table>

    <!-- Dialog Form Usulan -->
    <DisposalFormDialog v-model="showForm" @saved="onSaved" />

    <!-- Dialog Aksi (Verifikasi/Approval/Eksekusi) -->
    <DisposalActionDialog v-model="showAction" :request="selectedRequest" :type="actionType" @saved="onSaved" />
  </q-page>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useQuasar } from 'quasar'
import { useDisposalStore } from '@/stores/sarpras/disposal'
import { useContextStore } from '@/stores/contextStore'
import DisposalFormDialog from '@/components/sarpras/DisposalFormDialog.vue'
import DisposalActionDialog from '@/components/sarpras/DisposalActionDialog.vue'

const $q = useQuasar()
const store = useDisposalStore()
const contextStore = useContextStore()

const showForm = ref(false)
const showAction = ref(false)
const selectedRequest = ref(null)
const actionType = ref('')

const columns = [
  { name: 'request_number', label: 'No. Usulan', field: 'request_number', align: 'left', sortable: true },
  { name: 'proposal_date', label: 'Tanggal', field: 'proposal_date', align: 'center', format: val => new Date(val).toLocaleDateString('id-ID') },
  { name: 'total_assets', label: 'Jml Aset', field: 'total_assets', align: 'center' },
  { name: 'total_value', label: 'Total Nilai', field: 'total_value', align: 'right', format: val => `Rp ${val.toLocaleString('id-ID')}` },
  { name: 'approval_level', label: 'Level Approval', field: 'approval_level', align: 'center' },
  { name: 'status', label: 'Status', field: 'status', align: 'center' },
  { name: 'aksi', label: 'Aksi', align: 'center' }
]

const getStatusColor = (s) => {
  if (s === 'DRAFT') return 'grey'
  if (s === 'REJECTED') return 'red'
  if (s === 'EXECUTED') return 'primary'
  return 'green'
}

onMounted(() => {
  store.fetchRequests()
})

watch(() => contextStore.selectedSemesterId, () => {
  store.fetchRequests()
})

function bukaForm() {
  showForm.value = true
}

function bukaAksi(row, type) {
  selectedRequest.value = row
  actionType.value = type
  showAction.value = true
}

function onSaved() {
  store.fetchRequests()
}

function cetakBAP(row) {
  if (row.bast_pdf_url) {
    window.open(row.bast_pdf_url, '_blank')
  } else {
    $q.notify({ type: 'info', message: 'File BAP belum tersedia atau masih diproses.' })
  }
}
</script>

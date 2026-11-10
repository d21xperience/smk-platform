<template>
  <q-page class="q-pa-md">
    <div class="row q-mb-md items-center justify-between">
      <div class="text-h6 text-weight-bold">Pelayanan Peminjaman Aset</div>
      <q-btn color="primary" icon="add" label="Ajukan Peminjaman" @click="showForm = true" no-caps />
    </div>

    <q-table :rows="store.borrowings" :columns="columns" row-key="id" flat bordered :loading="store.isLoading">
      <template v-slot:body-cell-purpose_type="props">
        <q-td :props="props" class="text-center">
          <q-badge color="blue">{{ props.value }}</q-badge>
        </q-td>
      </template>

      <template v-slot:body-cell-status="props">
        <q-td :props="props" class="text-center">
          <q-badge :color="getStatusColor(props.value)">{{ props.value }}</q-badge>
        </q-td>
      </template>

      <template v-slot:body-cell-aksi="props">
        <q-td :props="props" class="text-center q-gutter-x-xs">
          <!-- Tombol Approve (jika PENDING_APPROVAL) -->
          <q-btn v-if="props.row.status === 'PENDING_APPROVAL'" dense round size="sm" color="amber" icon="fact_check" @click="bukaAksi(props.row, 'approve')">
            <q-tooltip>Approval</q-tooltip>
          </q-btn>
          
          <!-- Tombol Pickup (jika APPROVED) -->
          <q-btn v-if="props.row.status === 'APPROVED'" dense round size="sm" color="positive" icon="assignment_turned_in" @click="bukaAksi(props.row, 'pickup')">
            <q-tooltip>Serah Terima (Pickup)</q-tooltip>
          </q-btn>

          <!-- Tombol Return (jika PICKED_UP) -->
          <q-btn v-if="props.row.status === 'PICKED_UP'" dense round size="sm" color="blue" icon="assignment_return" @click="bukaAksi(props.row, 'return')">
            <q-tooltip>Pengembalian</q-tooltip>
          </q-btn>
        </q-td>
      </template>
    </q-table>

    <!-- Dialog Form Peminjaman -->
    <BorrowingFormDialog v-model="showForm" @saved="onSaved" />

    <!-- Dialog Aksi -->
    <BorrowingActionDialog v-model="showAction" :request="selectedRequest" :type="actionType" @saved="onSaved" />
  </q-page>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
// $q dihapus karena tidak digunakan
import { useBorrowingStore } from '@/stores/sarpras/borrowing'
import { useContextStore } from '@/stores/contextStore'
import BorrowingFormDialog from '@/components/sarpras/BorrowingFormDialog.vue'
import BorrowingActionDialog from '@/components/sarpras/BorrowingActionDialog.vue'

const store = useBorrowingStore()
const contextStore = useContextStore()

const showForm = ref(false)
const showAction = ref(false)
const selectedRequest = ref(null)
const actionType = ref('')

const columns = [
  { name: 'borrowing_number', label: 'No. Peminjaman', field: 'borrowing_number', align: 'left', sortable: true },
  { name: 'borrower_name', label: 'Peminjam', field: 'borrower_name', align: 'left', sortable: true },
  { name: 'asset_name', label: 'Aset', field: 'asset_name', align: 'left' },
  { name: 'borrow_date', label: 'Tgl Pinjam', field: 'borrow_date', align: 'center', format: val => new Date(val).toLocaleDateString('id-ID') },
  { name: 'expected_return_date', label: 'Tgl Kembali', field: 'expected_return_date', align: 'center', format: val => new Date(val).toLocaleDateString('id-ID') },
  { name: 'purpose_type', label: 'Tujuan', field: 'purpose_type', align: 'center' },
  { name: 'status', label: 'Status', field: 'status', align: 'center' },
  { name: 'aksi', label: 'Aksi', align: 'center' }
]

const getStatusColor = (s) => {
  if (s === 'PENDING_APPROVAL') return 'amber'
  if (s === 'REJECTED') return 'red'
  if (s === 'APPROVED') return 'green'
  if (s === 'PICKED_UP') return 'blue'
  if (s === 'RETURNED') return 'grey'
  return 'grey'
}

onMounted(() => {
  store.fetchBorrowings()
})

watch(() => contextStore.selectedSemesterId, (newVal) => {
  if (newVal) store.fetchBorrowings()
})

function bukaAksi(row, type) {
  selectedRequest.value = row
  actionType.value = type
  showAction.value = true
}

function onSaved() {
  store.fetchBorrowings()
}
</script>

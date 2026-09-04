<template>
  <q-page padding>
    <!-- 1. Header & Action Button -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <div class="text-h5 text-weight-bold">Manajemen Tagihan</div>
        <div class="text-grey-7">Kelola tagihan SPP, uang gedung, dan lainnya</div>
      </div>
      <q-btn color="primary" icon="add" label="Buat Tagihan Baru"
        @click="$router.push({ name: 'keuangan-tagihan-create' })" />
    </div>

    <!-- 2. Filter & Search Card -->
    <q-card class="q-mb-md">
      <q-card-section>
        <div class="row q-gutter-md">
          <q-select v-model="localFilters.status" :options="statusOptions" label="Status Tagihan" clearable emit-value
            map-options class="col-12 col-sm-3" @update:model-value="applyFilters" />
          <q-input v-model="localFilters.search" label="Cari ID Siswa / Deskripsi" clearable class="col-12 col-sm-6"
            @update:model-value="debouncedApplyFilters">
            <template v-slot:append>
              <q-icon name="search" />
            </template>
          </q-input>
          <q-btn color="secondary" icon="refresh" label="Reset" class="col-12 col-sm-2" @click="resetFilters" />
        </div>
      </q-card-section>
    </q-card>

    <!-- 3. Data Table -->
    <q-card>
      <q-table :rows="financeStore.invoices" :columns="columns" :loading="financeStore.isLoading"
        :pagination="tablePagination" row-key="invoiceId" @request="onRequest">
        <!-- Custom Cell: Status Badge -->
        <template v-slot:body-cell-status="props">
          <q-td :props="props">
            <q-badge :color="getStatusColor(props.row.status)">
              {{ props.row.status }}
            </q-badge>
          </q-td>
        </template>

        <!-- Custom Cell: Format Rupiah -->
        <template v-slot:body-cell-amount="props">
          <q-td :props="props">
            Rp {{ formatRupiah(props.row.amount) }}
          </q-td>
        </template>

        <template v-slot:body-cell-outstandingAmount="props">
          <q-td :props="props">
            <span :class="{ 'text-negative': props.row.outstandingAmount > 0 }">
              Rp {{ formatRupiah(props.row.outstandingAmount) }}
            </span>
          </q-td>
        </template>

        <!-- Custom Cell: Actions -->
        <template v-slot:body-cell-actions="props">
          <q-td :props="props">
            <q-btn flat dense round icon="visibility" color="primary" @click="viewDetail(props.row.invoiceId)">
              <q-tooltip>Detail</q-tooltip>
            </q-btn>

            <q-btn v-if="props.row.status !== 'PAID'" flat dense round icon="payment" color="positive"
              @click="openPaymentDialog(props.row)">
              <q-tooltip>Bayar</q-tooltip>
            </q-btn>
          </q-td>
        </template>

        <!-- Loading State -->
        <template v-slot:loading>
          <q-inner-loading showing color="primary" />
        </template>
      </q-table>
    </q-card>

    <!-- 4. Payment Dialog (Modal) -->
    <q-dialog v-model="showPaymentDialog" persistent>
      <q-card style="width: 500px">
        <q-card-section>
          <div class="text-h6">Catat Pembayaran</div>
          <div class="text-subtitle2 text-grey-7">Sisa Tagihan: Rp {{ formatRupiah(selectedInvoice?.outstandingAmount ||
            0)
            }}</div>
        </q-card-section>

        <q-card-section>
          <q-form @submit="onSubmitPayment">
            <q-input v-model.number="paymentForm.amount" label="Jumlah Pembayaran" type="number" :rules="[
              val => !!val || 'Jumlah wajib diisi',
              val => val > 0 || 'Jumlah harus lebih dari 0',
              val => val <= (selectedInvoice?.outstandingAmount || 0) || 'Melebihi sisa tagihan'
            ]" class="q-mb-md" />

            <q-select v-model="paymentForm.paymentMethod" :options="paymentMethodOptions" label="Metode Pembayaran"
              emit-value map-options :rules="[val => !!val || 'Metode wajib dipilih']" class="q-mb-md" />

            <q-input v-model="paymentForm.note" label="Catatan (Opsional)" type="textarea" class="q-mb-md" />

            <div class="row q-gutter-sm justify-end q-mt-md">
              <q-btn flat label="Batal" color="grey" @click="showPaymentDialog = false" />
              <q-btn type="submit" label="Proses Pembayaran" color="positive" :loading="financeStore.isLoading" />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- 5. Notification Dialog (Success/Error) -->
    <q-dialog v-model="showNotification">
      <q-card>
        <q-card-section :class="notificationType === 'success' ? 'bg-positive text-white' : 'bg-negative text-white'">
          <div class="text-h6">{{ notificationType === 'success' ? 'Berhasil' : 'Gagal' }}</div>
          <div>{{ notificationMessage }}</div>
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="OK" color="primary" @click="showNotification = false" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useFinanceStore } from '@/stores/finance/financeStore'
// import { useQuasar } from 'quasar'

const router = useRouter()
// const $q = useQuasar()
const financeStore = useFinanceStore()

// === STATE ===
const localFilters = ref({
  status: null,
  search: '',
})

const showPaymentDialog = ref(false)
const showNotification = ref(false)
const notificationMessage = ref('')
const notificationType = ref('success')

const selectedInvoice = ref(null)
const paymentForm = ref({
  amount: 0,
  paymentMethod: 'CASH',
  note: '',
})

// === OPTIONS ===
const statusOptions = [
  { label: 'Belum Lunas', value: 'UNPAID' },
  { label: 'Sebagian', value: 'PARTIAL' },
  { label: 'Lunas', value: 'PAID' },
  { label: 'Jatuh Tempo', value: 'OVERDUE' },
]

const paymentMethodOptions = [
  { label: 'Tunai (Cash)', value: 'CASH' },
  { label: 'Transfer Bank', value: 'TRANSFER' },
  { label: 'QRIS', value: 'QRIS' },
]

// === COLUMNS ===
const columns = [
  { name: 'invoiceId', label: 'ID Tagihan', field: 'invoiceId', align: 'left', sortable: true },
  { name: 'studentId', label: 'ID Siswa', field: 'studentId', align: 'left' },
  { name: 'type', label: 'Jenis', field: 'type', align: 'left' },
  { name: 'amount', label: 'Total Tagihan', field: 'amount', align: 'right', sortable: true },
  { name: 'outstandingAmount', label: 'Sisa', field: 'outstandingAmount', align: 'right', sortable: true },
  { name: 'status', label: 'Status', field: 'status', align: 'center' },
  { name: 'actions', label: 'Aksi', field: 'actions', align: 'center' },
]

// === COMPUTED ===
const tablePagination = computed(() => ({
  page: financeStore.pagination.page,
  rowsPerPage: financeStore.pagination.limit,
  rowsNumber: financeStore.pagination.total,
}))

// === METHODS ===
function formatRupiah(amount) {
  return new Intl.NumberFormat('id-ID').format(amount || 0)
}

function getStatusColor(status) {
  const colors = {
    UNPAID: 'orange',
    PARTIAL: 'blue',
    PAID: 'positive',
    OVERDUE: 'negative',
  }
  return colors[status] || 'grey'
}

function applyFilters() {
  financeStore.updateFilters(localFilters.value)
  financeStore.fetchInvoices({ page: 1 })
}

// Simple debounce for search input
let searchTimeout
function debouncedApplyFilters() {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    applyFilters()
  }, 500)
}

function resetFilters() {
  localFilters.value = { status: null, search: '' }
  applyFilters()
}

function onRequest(props) {
  const { page, rowsPerPage } = props.pagination
  financeStore.fetchInvoices({ page, limit: rowsPerPage })
}

function viewDetail(invoiceId) {
  router.push({ name: 'keuangan-tagihan-detail', params: { id: invoiceId } })
}

function openPaymentDialog(invoice) {
  selectedInvoice.value = invoice
  paymentForm.value = {
    amount: invoice.outstandingAmount, // Default fill with outstanding amount
    paymentMethod: 'CASH',
    note: '',
  }
  showPaymentDialog.value = true
}

async function onSubmitPayment() {
  const result = await financeStore.processPayment({
    invoiceId: selectedInvoice.value.invoiceId,
    amount: paymentForm.value.amount,
    paymentMethod: paymentForm.value.paymentMethod,
    note: paymentForm.value.note,
  })

  if (result) {
    showPaymentDialog.value = false
    notificationMessage.value = 'Pembayaran berhasil dicatat!'
    notificationType.value = 'success'
    showNotification.value = true
    financeStore.fetchInvoices() // Refresh table
  } else {
    notificationMessage.value = financeStore.error?.message || 'Gagal memproses pembayaran'
    notificationType.value = 'error'
    showNotification.value = true
  }
}

// === LIFECYCLE ===
onMounted(() => {
  financeStore.fetchInvoices()
})
</script>

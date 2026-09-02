<template>
  <q-page padding>
    <div class="row items-center justify-between q-mb-md">
      <div class="text-h5">Homeroom Billing</div>
      <q-btn
        flat
        dense
        icon="refresh"
        label="Refresh"
        @click="loadBillingData"
      />
    </div>

    <q-banner v-if="!hasPermission" class="bg-negative text-white q-mb-md">
      <template #avatar>
        <q-icon name="lock" />
      </template>
      Anda tidak memiliki izin untuk mengakses halaman ini.
    </q-banner>

    <template v-else>
      <q-banner v-if="error" class="bg-negative text-white q-mb-md">
        {{ error }}
      </q-banner>

      <q-banner v-if="!isLoading && !hasAssignment" class="bg-warning text-white q-mb-md">
        Anda tidak memiliki homeroom assignment. Hanya wali kelas yang dapat melihat data billing.
      </q-banner>

      <q-card v-if="isLoading" class="q-mb-md">
        <q-card-section class="text-center">
          <q-spinner color="primary" size="2em" />
          <div class="q-mt-sm">Loading...</div>
        </q-card-section>
      </q-card>

      <template v-else-if="hasAssignment && summary">
        <q-card class="q-mb-md">
          <q-card-section>
            <div class="text-h6">
              Kelas {{ summary.className }} - {{ summary.academicYearName }} / {{ summary.semesterName }}
            </div>
            <div class="text-caption text-grey">
              READ-ONLY: Data billing berasal dari tu-core Finance. Anda hanya dapat melihat, tidak dapat mengubah.
            </div>
          </q-card-section>
        </q-card>

        <div class="row q-col-gutter-md q-mb-md">
          <div class="col-12 col-sm-6 col-md-3">
            <q-card>
              <q-card-section class="text-center">
                <q-icon name="groups" color="primary" size="2.5em" />
                <div class="text-h5 q-mt-sm">{{ summary.totalStudents }}</div>
                <div class="text-caption">Total Siswa</div>
              </q-card-section>
            </q-card>
          </div>

          <div class="col-12 col-sm-6 col-md-3">
            <q-card>
              <q-card-section class="text-center">
                <q-icon name="receipt_long" color="info" size="2.5em" />
                <div class="text-h5 q-mt-sm">{{ summary.totalBills }}</div>
                <div class="text-caption">Total Tagihan</div>
              </q-card-section>
            </q-card>
          </div>

          <div class="col-12 col-sm-6 col-md-3">
            <q-card>
              <q-card-section class="text-center">
                <q-icon name="paid" color="positive" size="2.5em" />
                <div class="text-h5 q-mt-sm">{{ formatCurrency(summary.totalPaidAmount) }}</div>
                <div class="text-caption">Total Dibayar</div>
                <q-linear-progress
                  :value="summary.getCollectionRate() / 100"
                  color="positive"
                  class="q-mt-sm"
                />
                <div class="text-caption">{{ summary.getCollectionRate() }}% Collection Rate</div>
              </q-card-section>
            </q-card>
          </div>

          <div class="col-12 col-sm-6 col-md-3">
            <q-card>
              <q-card-section class="text-center">
                <q-icon name="warning" color="negative" size="2.5em" />
                <div class="text-h5 q-mt-sm">{{ formatCurrency(summary.totalOutstandingAmount) }}</div>
                <div class="text-caption">Total Tunggakan</div>
              </q-card-section>
            </q-card>
          </div>
        </div>

        <q-card class="q-mb-md">
          <q-card-section>
            <q-select
              v-model="filterStatus"
              :options="statusOptions"
              label="Filter Status"
              outlined
              dense
              emit-value
              map-options
              clearable
              style="min-width: 200px"
            />
          </q-card-section>
        </q-card>

        <q-card>
          <q-card-section>
            <div class="text-h6 q-mb-md">Daftar Tagihan</div>

            <q-table
              :rows="filteredBills"
              :columns="columns"
              row-key="billId"
              :pagination="{ rowsPerPage: 10 }"
              flat
              bordered
            >
              <template #body-cell-status="props">
                <q-td :props="props">
                  <q-badge :color="getBillStatusColor(props.value)">
                    {{ getBillStatusLabel(props.value) }}
                  </q-badge>
                </q-td>
              </template>

              <template #body-cell-amount="props">
                <q-td :props="props">
                  {{ formatCurrency(props.value) }}
                </q-td>
              </template>

              <template #body-cell-paidAmount="props">
                <q-td :props="props">
                  <span class="text-positive">{{ formatCurrency(props.value) }}</span>
                </q-td>
              </template>

              <template #body-cell-outstandingAmount="props">
                <q-td :props="props">
                  <span :class="props.value > 0 ? 'text-negative' : 'text-grey'">
                    {{ formatCurrency(props.value) }}
                  </span>
                </q-td>
              </template>

              <template #body-cell-paymentProgress="props">
                <q-td :props="props">
                  <div class="row items-center q-gutter-sm">
                    <q-linear-progress
                      :value="props.row.getPaymentPercentage() / 100"
                      :color="props.row.isSettled() ? 'positive' : 'warning'"
                      style="width: 80px"
                    />
                    <span class="text-caption">{{ props.row.getPaymentPercentage() }}%</span>
                  </div>
                </q-td>
              </template>
            </q-table>
          </q-card-section>
        </q-card>
      </template>
    </template>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useHomeroomBilling } from '../../composables/useHomeroomBilling.js'
import { useAuth } from '../../composables/useAuth.js'
import { getBillStatusLabel, getBillStatusColor } from '../../domain/homeroom-billing/models/BillStatus.js'

const {
  summary,
  bills,
  isLoading,
  error,
  hasAssignment,
  loadBillingData
} = useHomeroomBilling()

const { currentUser } = useAuth()

const filterStatus = ref(null)

const hasPermission = computed(() => {
  if (!currentUser.value) return false
  const permissions = currentUser.value.permissions || []
  const role = currentUser.value.role

  const isAdmin = role === 'admin' || role === 'superadmin'
  const hasBillingPermission = permissions.includes('homeroom.billing.read')

  return isAdmin || hasBillingPermission
})

const statusOptions = [
  { label: 'Belum Bayar', value: 'unpaid' },
  { label: 'Bayar Sebagian', value: 'partial' },
  { label: 'Lunas', value: 'paid' },
  { label: 'Terlambat', value: 'overdue' },
  { label: 'Dibebaskan', value: 'waived' }
]

const columns = [
  { name: 'studentName', label: 'Nama Siswa', field: 'studentName', align: 'left', sortable: true },
  { name: 'billName', label: 'Tagihan', field: 'billName', align: 'left', sortable: true },
  { name: 'amount', label: 'Nominal', field: 'amount', align: 'right', sortable: true },
  { name: 'paidAmount', label: 'Dibayar', field: 'paidAmount', align: 'right', sortable: true },
  { name: 'outstandingAmount', label: 'Sisa', field: 'outstandingAmount', align: 'right', sortable: true },
  { name: 'paymentProgress', label: 'Progress', field: 'paidAmount', align: 'center' },
  { name: 'status', label: 'Status', field: 'status', align: 'center', sortable: true },
  { name: 'dueDate', label: 'Jatuh Tempo', field: 'dueDate', align: 'left', sortable: true, format: val => new Date(val).toLocaleDateString('id-ID') }
]

const filteredBills = computed(() => {
  if (!filterStatus.value) {
    return bills.value
  }
  return bills.value.filter(b => b.status === filterStatus.value)
})

onMounted(async () => {
  if (hasPermission.value) {
    await loadBillingData()
  }
})

const formatCurrency = (amount) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0
  }).format(amount)
}
</script>

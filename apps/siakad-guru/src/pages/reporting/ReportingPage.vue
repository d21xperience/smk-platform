<template>
  <q-page padding>
    <div class="text-h5 q-mb-md">Dashboard</div>

    <q-banner v-if="reportingError" class="bg-negative text-white q-mb-md">
      {{ reportingError }}
    </q-banner>

    <q-card v-if="isReportingLoading" class="q-mb-md">
      <q-card-section class="text-center">
        <q-spinner color="primary" size="2em" />
        <div class="q-mt-sm">Loading...</div>
      </q-card-section>
    </q-card>

    <template v-else-if="dashboard">
      <!-- Period Info -->
      <q-card class="q-mb-md">
        <q-card-section>
          <div class="text-h6">Periode Aktif</div>
          <div class="text-caption text-grey">
            Tahun Pelajaran: {{ dashboard.period.academicYearName }} |
            {{ dashboard.period.semesterName }} ({{ dashboard.period.semesterId }})
          </div>
        </q-card-section>
      </q-card>

      <!-- Generic Teacher Summary Cards -->
      <div class="text-h6 q-mb-sm">Ringkasan Akademik</div>
      <div class="row q-col-gutter-md q-mb-md">
        <div class="col-12 col-sm-6 col-md-4">
          <q-card>
            <q-card-section class="text-center">
              <q-icon name="school" color="primary" size="3em" />
              <div class="text-h5 q-mt-sm">
                {{ dashboard.completedTeachingSessions }} / {{ dashboard.totalTeachingSessions }}
              </div>
              <div class="text-caption">Teaching Sessions</div>
            </q-card-section>
          </q-card>
        </div>

        <div class="col-12 col-sm-6 col-md-4">
          <q-card>
            <q-card-section class="text-center">
              <q-icon name="how_to_reg" color="positive" size="3em" />
              <div class="text-h5 q-mt-sm">
                {{ dashboard.submittedAttendanceSessions }} / {{ dashboard.totalAttendanceSessions }}
              </div>
              <div class="text-caption">Attendance Submitted</div>
            </q-card-section>
          </q-card>
        </div>

        <div class="col-12 col-sm-6 col-md-4">
          <q-card>
            <q-card-section class="text-center">
              <q-icon name="grading" color="info" size="3em" />
              <div class="text-h5 q-mt-sm">
                {{ dashboard.finalizedAssessments }} / {{ dashboard.totalAssessments }}
              </div>
              <div class="text-caption">Assessments Finalized</div>
            </q-card-section>
          </q-card>
        </div>

        <div class="col-12 col-sm-6 col-md-4">
          <q-card>
            <q-card-section class="text-center">
              <q-icon name="book" color="warning" size="3em" />
              <div class="text-h5 q-mt-sm">
                {{ dashboard.submittedJournals }} / {{ dashboard.totalJournals }}
              </div>
              <div class="text-caption">Journals Submitted</div>
            </q-card-section>
          </q-card>
        </div>

        <div class="col-12 col-sm-6 col-md-4">
          <q-card>
            <q-card-section class="text-center">
              <q-icon name="trending_up" color="purple" size="3em" />
              <div class="text-h5 q-mt-sm">
                {{ dashboard.totalProgressRecords }}
              </div>
              <div class="text-caption">Progress Records</div>
            </q-card-section>
          </q-card>
        </div>
      </div>

      <!-- Homeroom Billing Section (Conditional) -->
      <template v-if="canViewHomeroomBilling">
        <q-separator class="q-my-md" />

        <div class="row items-center justify-between q-mb-sm">
          <div class="text-h6">Homeroom Billing</div>
          <q-btn flat dense color="primary" label="Lihat Detail" icon="arrow_forward"
            :to="{ name: 'homeroom-billing' }" />
        </div>

        <q-banner v-if="billingError" class="bg-warning text-white q-mb-md">
          {{ billingError }}
        </q-banner>

        <q-card v-if="isBillingLoading" class="q-mb-md">
          <q-card-section class="text-center">
            <q-spinner color="primary" size="2em" />
            <div class="q-mt-sm">Loading billing data...</div>
          </q-card-section>
        </q-card>

        <template v-else-if="billingSummary">
          <q-card class="q-mb-md">
            <q-card-section>
              <div class="text-subtitle1 q-mb-sm">
                Kelas {{ billingSummary.className }} - {{ billingSummary.academicYearName }} / {{
                  billingSummary.semesterName }}
              </div>
              <div class="text-caption text-grey q-mb-md">
                READ-ONLY: Data billing berasal dari tu-core Finance
              </div>

              <div class="row q-col-gutter-md">
                <div class="col-12 col-sm-6 col-md-3">
                  <div class="text-center">
                    <q-icon name="groups" color="primary" size="2em" />
                    <div class="text-h6 q-mt-sm">{{ billingSummary.totalStudents }}</div>
                    <div class="text-caption">Total Siswa</div>
                  </div>
                </div>

                <div class="col-12 col-sm-6 col-md-3">
                  <div class="text-center">
                    <q-icon name="receipt_long" color="info" size="2em" />
                    <div class="text-h6 q-mt-sm">{{ billingSummary.totalBills }}</div>
                    <div class="text-caption">Total Tagihan</div>
                  </div>
                </div>

                <div class="col-12 col-sm-6 col-md-3">
                  <div class="text-center">
                    <q-icon name="paid" color="positive" size="2em" />
                    <div class="text-h6 q-mt-sm">{{ formatCurrency(billingSummary.totalPaidAmount) }}</div>
                    <div class="text-caption">Total Dibayar</div>
                    <q-linear-progress :value="billingSummary.getCollectionRate() / 100" color="positive"
                      class="q-mt-sm" />
                    <div class="text-caption">{{ billingSummary.getCollectionRate() }}%</div>
                  </div>
                </div>

                <div class="col-12 col-sm-6 col-md-3">
                  <div class="text-center">
                    <q-icon name="warning" color="negative" size="2em" />
                    <div class="text-h6 q-mt-sm">{{ formatCurrency(billingSummary.totalOutstandingAmount) }}</div>
                    <div class="text-caption">Total Tunggakan</div>
                  </div>
                </div>
              </div>
            </q-card-section>
          </q-card>

          <!-- Quick Stats: Outstanding Bills -->
          <q-card v-if="outstandingBillsCount > 0">
            <q-card-section>
              <div class="row items-center">
                <q-icon name="notification_important" color="warning" size="2em" class="q-mr-md" />
                <div>
                  <div class="text-subtitle1">
                    {{ outstandingBillsCount }} tagihan belum lunas
                  </div>
                  <div class="text-caption text-grey">
                    {{ overdueBillsCount }} di antaranya terlambat
                  </div>
                </div>
                <q-space />
                <q-btn flat color="warning" label="Lihat Semua" :to="{ name: 'homeroom-billing' }" />
              </div>
            </q-card-section>
          </q-card>
        </template>

        <q-card v-else-if="!hasHomeroomAssignment">
          <q-card-section class="text-center text-grey">
            Anda tidak memiliki homeroom assignment untuk periode ini.
          </q-card-section>
        </q-card>
      </template>
    </template>

    <q-card v-else>
      <q-card-section class="text-center text-grey">
        Data dashboard tidak tersedia.
      </q-card-section>
    </q-card>
  </q-page>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useReporting } from '../../composables/useReporting.js'
import { useHomeroomBilling } from '../../composables/useHomeroomBilling.js'
import { useAuth } from '../../composables/useAuth.js'

const {
  dashboard,
  isLoading: isReportingLoading,
  error: reportingError,
  loadDashboard
} = useReporting()

const {
  summary: billingSummary,
  bills: billingBills,
  isLoading: isBillingLoading,
  error: billingError,
  hasAssignment: hasHomeroomAssignment,
  loadBillingData
} = useHomeroomBilling()

const { currentUser } = useAuth()

const canViewHomeroomBilling = computed(() => {
  if (!currentUser.value) return false
  const permissions = currentUser.value.permissions || []
  const role = currentUser.value.role

  const isAdmin = role === 'admin' || role === 'superadmin'
  const hasPermission = permissions.includes('homeroom.billing.read')

  return isAdmin || hasPermission
})

const outstandingBillsCount = computed(() => {
  if (!billingBills.value) return 0
  return billingBills.value.filter(b => b.isOutstanding()).length
})

const overdueBillsCount = computed(() => {
  if (!billingBills.value) return 0
  return billingBills.value.filter(b => b.status === 'overdue').length
})

onMounted(async () => {
  // Load generic dashboard (all teachers)
  await loadDashboard()

  // Load homeroom billing only if user has permission
  if (canViewHomeroomBilling.value) {
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

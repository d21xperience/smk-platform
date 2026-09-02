<!-- FILE: src/pages/homeroom-attendance/HomeroomAttendancePage.vue -->
<!-- STATUS: NEW -->
<!-- STATUS IMPLEMENTASI: COMPLETE -->

<template>
  <q-page class="q-pa-md">
    <div class="row items-center q-mb-md">
      <div class="text-h5">Rekap Kehadiran Wali Kelas</div>
      <q-badge v-if="className" color="primary" class="q-ml-md" style="font-size: 14px; padding: 4px 12px;">
        Kelas {{ className }}
      </q-badge>
    </div>

    <q-banner v-if="error" class="bg-negative text-white q-mb-md">
      {{ error }}
    </q-banner>

    <div v-if="loading" class="flex flex-center" style="min-height: 200px;">
      <q-spinner-dots size="60px" color="primary" />
    </div>

    <template v-if="!loading">
      <div v-if="summary" class="row q-col-gutter-md q-mb-md">
        <div class="col-6 col-sm-3">
          <q-card flat bordered>
            <q-card-section class="text-center q-py-sm">
              <div class="text-h4 text-primary">{{ summary.totalStudents }}</div>
              <div class="text-caption text-grey-7">Total Siswa</div>
            </q-card-section>
          </q-card>
        </div>
        <div class="col-6 col-sm-3">
          <q-card flat bordered>
            <q-card-section class="text-center q-py-sm">
              <div class="text-h4" :class="getPercentageClass(summary.classAveragePercentage)">
                {{ summary.classAveragePercentage }}%
              </div>
              <div class="text-caption text-grey-7">Rata-rata Kehadiran</div>
            </q-card-section>
          </q-card>
        </div>
        <div class="col-6 col-sm-3">
          <q-card flat bordered>
            <q-card-section class="text-center q-py-sm">
              <div class="text-h4 text-positive">{{ summary.goodCount }}</div>
              <div class="text-caption text-grey-7">Kehadiran Baik</div>
            </q-card-section>
          </q-card>
        </div>
        <div class="col-6 col-sm-3">
          <q-card flat bordered>
            <q-card-section class="text-center q-py-sm">
              <div class="text-h4 text-negative">{{ summary.criticalCount }}</div>
              <div class="text-caption text-grey-7">Perlu Perhatian</div>
            </q-card-section>
          </q-card>
        </div>
      </div>

      <q-table flat bordered :rows="students" :columns="columns" row-key="studentId"
        :pagination="{ rowsPerPage: 15, sortBy: 'seatNumber', ascending: true }">
        <template #body="props">
          <q-tr :props="props">
            <q-td auto-width>
              <q-btn flat dense round size="sm" :icon="props.expand ? 'keyboard_arrow_up' : 'keyboard_arrow_down'"
                @click="props.expand = !props.expand" />
            </q-td>
            <q-td key="seatNumber" :props="props" class="text-center">
              {{ props.row.seatNumber }}
            </q-td>
            <q-td key="studentName" :props="props">
              {{ props.row.studentName }}
            </q-td>
            <q-td key="present" :props="props" class="text-center text-positive">
              {{ props.row.semesterTotal.present }}
            </q-td>
            <q-td key="sick" :props="props" class="text-center text-warning">
              {{ props.row.semesterTotal.sick }}
            </q-td>
            <q-td key="permitted" :props="props" class="text-center text-info">
              {{ props.row.semesterTotal.permitted }}
            </q-td>
            <q-td key="absent" :props="props" class="text-center"
              :class="props.row.semesterTotal.absent > 0 ? 'text-negative' : 'text-grey'">
              {{ props.row.semesterTotal.absent }}
            </q-td>
            <q-td key="percentage" :props="props" class="text-center">
              <strong :class="getPercentageClass(props.row.semesterTotal.percentage)">
                {{ props.row.semesterTotal.percentage }}%
              </strong>
            </q-td>
            <q-td key="status" :props="props" class="text-center">
              <q-badge :color="getStatusColor(props.row.status)">
                {{ getStatusLabel(props.row.status) }}
              </q-badge>
            </q-td>
          </q-tr>

          <q-tr v-show="props.expand" :props="props">
            <q-td colspan="100%">
              <div class="q-pa-md">
                <div class="text-subtitle2 q-mb-sm">
                  Detail Bulanan — {{ props.row.studentName }}
                </div>
                <q-table flat dense bordered :rows="props.row.monthlyRecords" :columns="monthlyColumns"
                  row-key="monthId" :pagination="{ rowsPerPage: 0 }" hide-bottom>
                  <template #body-cell-percentage="monthProps">
                    <q-td :props="monthProps" class="text-center">
                      <strong :class="getPercentageClass(monthProps.row.percentage)">
                        {{ monthProps.row.percentage }}%
                      </strong>
                    </q-td>
                  </template>
                </q-table>
              </div>
            </q-td>
          </q-tr>
        </template>
      </q-table>
    </template>
  </q-page>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useHomeroomAttendance } from '@/composables/useHomeroomAttendance.js'

const {
  className,
  summary,
  students,
  loading,
  error,
  loadAttendanceSummary,
} = useHomeroomAttendance()

const columns = ref([
  { name: 'expand', label: '', field: 'expand', align: 'center' },
  { name: 'seatNumber', label: 'No', field: 'seatNumber', align: 'center', sortable: true },
  { name: 'studentName', label: 'Nama Siswa', field: 'studentName', align: 'left', sortable: true },
  { name: 'present', label: 'Hadir', field: row => row.semesterTotal.present, align: 'center', sortable: true },
  { name: 'sick', label: 'Sakit', field: row => row.semesterTotal.sick, align: 'center', sortable: true },
  { name: 'permitted', label: 'Izin', field: row => row.semesterTotal.permitted, align: 'center', sortable: true },
  { name: 'absent', label: 'Alpha', field: row => row.semesterTotal.absent, align: 'center', sortable: true },
  { name: 'percentage', label: '% Kehadiran', field: row => row.semesterTotal.percentage, align: 'center', sortable: true },
  { name: 'status', label: 'Status', field: 'status', align: 'center', sortable: true },
])

const monthlyColumns = ref([
  { name: 'monthName', label: 'Bulan', field: 'monthName', align: 'left' },
  { name: 'present', label: 'Hadir', field: 'present', align: 'center' },
  { name: 'sick', label: 'Sakit', field: 'sick', align: 'center' },
  { name: 'permitted', label: 'Izin', field: 'permitted', align: 'center' },
  { name: 'absent', label: 'Alpha', field: 'absent', align: 'center' },
  { name: 'totalDays', label: 'Total Hari', field: 'totalDays', align: 'center' },
  { name: 'percentage', label: '% Kehadiran', field: 'percentage', align: 'center' },
])

function getPercentageClass(percentage) {
  if (percentage >= 90) return 'text-positive'
  if (percentage >= 75) return 'text-warning'
  return 'text-negative'
}

function getStatusColor(status) {
  switch (status) {
    case 'GOOD': return 'positive'
    case 'WARNING': return 'warning'
    case 'CRITICAL': return 'negative'
    default: return 'grey'
  }
}

function getStatusLabel(status) {
  switch (status) {
    case 'GOOD': return 'Baik'
    case 'WARNING': return 'Perhatian'
    case 'CRITICAL': return 'Kritis'
    default: return 'Tidak diketahui'
  }
}

onMounted(() => {
  loadAttendanceSummary()
})
</script>

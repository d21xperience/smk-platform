<!-- FILE: src/pages/homeroom-progress/HomeroomProgressPage.vue -->
<!-- STATUS: NEW -->
<!-- STATUS IMPLEMENTASI: COMPLETE -->

<template>
  <q-page class="q-pa-md">
    <div class="row items-center q-mb-md">
      <div class="text-h5">Rekap Nilai Wali Kelas</div>
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
              <div class="text-h4 text-primary">{{ summary.classAverage }}</div>
              <div class="text-caption text-grey-7">Rata-rata Kelas</div>
            </q-card-section>
          </q-card>
        </div>
        <div class="col-6 col-sm-3">
          <q-card flat bordered>
            <q-card-section class="text-center q-py-sm">
              <div class="text-h4" :class="summary.passRate >= 75 ? 'text-positive' : 'text-negative'">
                {{ summary.passRate }}%
              </div>
              <div class="text-caption text-grey-7">Tingkat Ketuntasan</div>
            </q-card-section>
          </q-card>
        </div>
        <div class="col-6 col-sm-3">
          <q-card flat bordered>
            <q-card-section class="text-center q-py-sm">
              <div class="text-h4 text-secondary">{{ summary.highestAverage }}</div>
              <div class="text-caption text-grey-7">Rata-rata Tertinggi</div>
            </q-card-section>
          </q-card>
        </div>
      </div>

      <q-table flat bordered :rows="students" :columns="columns" row-key="studentId"
        :pagination="{ rowsPerPage: 15, sortBy: 'classRank', ascending: true }">
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
            <q-td key="overallAverage" :props="props" class="text-center">
              <span :class="getAverageClass(props.row.overallAverage)">
                {{ props.row.overallAverage }}
              </span>
            </q-td>
            <q-td key="classRank" :props="props" class="text-center">
              <q-badge :color="getRankColor(props.row.classRank)" outline>
                #{{ props.row.classRank }}
              </q-badge>
            </q-td>
            <q-td key="passCount" :props="props" class="text-center text-positive">
              {{ props.row.passCount }}
            </q-td>
            <q-td key="failCount" :props="props" class="text-center"
              :class="props.row.failCount > 0 ? 'text-negative' : 'text-grey'">
              {{ props.row.failCount }}
            </q-td>
          </q-tr>

          <q-tr v-show="props.expand" :props="props">
            <q-td colspan="100%">
              <div class="q-pa-md">
                <div class="text-subtitle2 q-mb-sm">
                  Detail Nilai — {{ props.row.studentName }}
                </div>
                <q-table flat dense bordered :rows="props.row.subjects" :columns="subjectColumns" row-key="subjectId"
                  :pagination="{ rowsPerPage: 0 }" hide-bottom>
                  <template #body-cell-finalScore="subjectProps">
                    <q-td :props="subjectProps" class="text-center">
                      <strong :class="subjectProps.row.status === 'PASS' ? 'text-positive' : 'text-negative'">
                        {{ subjectProps.row.finalScore }}
                      </strong>
                    </q-td>
                  </template>

                  <template #body-cell-status="subjectProps">
                    <q-td :props="subjectProps" class="text-center">
                      <q-badge :color="subjectProps.row.status === 'PASS' ? 'positive' : 'negative'">
                        {{ subjectProps.row.status === 'PASS' ? 'Tuntas' : 'Belum' }}
                      </q-badge>
                    </q-td>
                  </template>

                  <template #body-cell-trend="subjectProps">
                    <q-td :props="subjectProps" class="text-center">
                      <q-icon :name="getTrendIcon(subjectProps.row.trend)"
                        :color="getTrendColor(subjectProps.row.trend)" size="sm">
                        <q-tooltip>{{ getTrendLabel(subjectProps.row.trend) }}</q-tooltip>
                      </q-icon>
                    </q-td>
                  </template>

                  <template #body-cell-predicate="subjectProps">
                    <q-td :props="subjectProps" class="text-center">
                      <q-badge :color="getPredicateColor(subjectProps.row.predicate)" outline>
                        {{ subjectProps.row.predicate }}
                      </q-badge>
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
import { useHomeroomProgress } from '@/composables/useHomeroomProgress.js'

const {
  className,
  summary,
  students,
  loading,
  error,
  loadProgressSummary,
} = useHomeroomProgress()

const columns = ref([
  { name: 'expand', label: '', field: 'expand', align: 'center' },
  { name: 'seatNumber', label: 'No', field: 'seatNumber', align: 'center', sortable: true },
  { name: 'studentName', label: 'Nama Siswa', field: 'studentName', align: 'left', sortable: true },
  { name: 'overallAverage', label: 'Rata-rata', field: 'overallAverage', align: 'center', sortable: true },
  { name: 'classRank', label: 'Rank', field: 'classRank', align: 'center', sortable: true },
  { name: 'passCount', label: 'Tuntas', field: 'passCount', align: 'center' },
  { name: 'failCount', label: 'Belum', field: 'failCount', align: 'center' },
])

const subjectColumns = ref([
  { name: 'subjectName', label: 'Mata Pelajaran', field: 'subjectName', align: 'left' },
  { name: 'formativeScore', label: 'Formatif', field: 'formativeScore', align: 'center' },
  { name: 'summativeScore', label: 'Sumatif', field: 'summativeScore', align: 'center' },
  { name: 'finalScore', label: 'Nilai Akhir', field: 'finalScore', align: 'center' },
  { name: 'kkm', label: 'KKM', field: 'kkm', align: 'center' },
  { name: 'status', label: 'Status', field: 'status', align: 'center' },
  { name: 'trend', label: 'Tren', field: 'trend', align: 'center' },
  { name: 'predicate', label: 'Predikat', field: 'predicate', align: 'center' },
])

function getAverageClass(avg) {
  if (avg >= 90) return 'text-positive text-bold'
  if (avg >= 75) return 'text-primary'
  return 'text-negative text-bold'
}

function getRankColor(rank) {
  if (rank === 1) return 'amber'
  if (rank === 2) return 'grey-6'
  if (rank === 3) return 'brown-4'
  return 'grey'
}

function getTrendIcon(trend) {
  switch (trend) {
    case 'UP': return 'trending_up'
    case 'DOWN': return 'trending_down'
    case 'STABLE': return 'trending_flat'
    case 'NEW': return 'fiber_new'
    default: return 'help_outline'
  }
}

function getTrendColor(trend) {
  switch (trend) {
    case 'UP': return 'positive'
    case 'DOWN': return 'negative'
    case 'STABLE': return 'grey'
    case 'NEW': return 'info'
    default: return 'grey'
  }
}

function getTrendLabel(trend) {
  switch (trend) {
    case 'UP': return 'Naik'
    case 'DOWN': return 'Turun'
    case 'STABLE': return 'Stabil'
    case 'NEW': return 'Baru'
    default: return 'Tidak diketahui'
  }
}

function getPredicateColor(predicate) {
  switch (predicate) {
    case 'A': return 'positive'
    case 'B': return 'primary'
    case 'C': return 'warning'
    case 'D': return 'negative'
    default: return 'grey'
  }
}

onMounted(() => {
  loadProgressSummary()
})
</script>

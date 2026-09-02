<template>
  <div class="q-pa-md">
    <h6>Dashboard & Laporan</h6>

    <!-- Statistik Global -->
    <div v-if="dashboardStats" class="row q-col-gutter-md q-mb-md">
      <div class="col-md-2 col-sm-4 col-xs-6">
        <q-card flat bordered>
          <q-card-section>
            <div class="text-grey-6 text-caption">Total Siswa</div>
            <div class="text-h5">{{ dashboardStats.totalStudents }}</div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-md-2 col-sm-4 col-xs-6">
        <q-card flat bordered>
          <q-card-section>
            <div class="text-grey-6 text-caption">Rata-rata Nilai</div>
            <div class="text-h5">{{ dashboardStats.averageScore.toFixed(1) }}</div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-md-2 col-sm-4 col-xs-6">
        <q-card flat bordered>
          <q-card-section>
            <div class="text-grey-6 text-caption">Tingkat Kelulusan</div>
            <div class="text-h5">{{ dashboardStats.promotionRate.toFixed(1) }}%</div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-md-2 col-sm-4 col-xs-6">
        <q-card flat bordered>
          <q-card-section>
            <div class="text-grey-6 text-caption">Kehadiran</div>
            <div class="text-h5">{{ dashboardStats.attendanceRate.toFixed(1) }}%</div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-md-2 col-sm-4 col-xs-6">
        <q-card flat bordered>
          <q-card-section>
            <div class="text-grey-6 text-caption">Total Guru</div>
            <div class="text-h5">{{ dashboardStats.totalTeachers }}</div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-md-2 col-sm-4 col-xs-6">
        <q-card flat bordered>
          <q-card-section>
            <div class="text-grey-6 text-caption">Total Kelas</div>
            <div class="text-h5">{{ dashboardStats.totalClasses }}</div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Distribusi Nilai -->
    <div v-if="dashboardStats?.gradeDistribution" class="row q-mb-md">
      <div class="col-12">
        <q-card flat bordered>
          <q-card-section>
            <div class="text-subtitle2">Distribusi Nilai</div>
          </q-card-section>
          <q-card-section>
            <div class="row q-col-gutter-sm">
              <div
                v-for="(count, predicate) in dashboardStats.gradeDistribution"
                :key="predicate"
                class="col-auto"
              >
                <q-chip :color="getPredicateColor(predicate)" text-color="white">
                  {{ predicate }}: {{ count }}
                </q-chip>
              </div>
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Laporan Kelas -->
    <q-card flat bordered>
      <q-card-section>
        <div class="row items-center">
          <div class="col">
            <div class="text-subtitle2">Laporan Kelas</div>
          </div>
          <div class="col-auto q-gutter-sm">
            <q-select
              v-model="selectedClassForReport"
              :options="classOptions"
              label="Kelas"
              dense
              outlined
              style="min-width: 150px"
            />
            <q-btn color="primary" label="Generate" @click="generateReport" :loading="reportLoading" />
          </div>
        </div>
      </q-card-section>
      <q-card-section>
        <q-table
          :rows="reportRows"
          :columns="reportColumns"
          row-key="nis"
          flat
          bordered
          :loading="reportLoading"
        />
      </q-card-section>
      <q-card-section v-if="reportRows.length > 0">
        <div class="q-gutter-sm">
          <q-btn color="blue" label="CSV" @click="exportReport('CSV')" />
          <q-btn color="green" label="Excel" @click="exportReport('EXCEL')" />
          <q-btn color="red" label="PDF" @click="exportReport('PDF')" />
        </div>
      </q-card-section>
    </q-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useReportingStore } from '@/stores/reporting.store';
import { useQuasar } from 'quasar';

const $q = useQuasar();
const reportingStore = useReportingStore();

const dashboardStats = computed(() => reportingStore.dashboardStats);
const reportLoading = computed(() => reportingStore.isLoading);
const reportRows = computed(() => reportingStore.currentReport?.rows || []);
const reportColumns = computed(() => {
  const headers = reportingStore.currentReport?.headers || [];
  return headers.map(h => ({ name: h.key, label: h.label, field: h.key, align: 'left' }));
});

const selectedClassForReport = ref('cls-xi-rpl-1');
const classOptions = [
  { label: 'XI RPL 1', value: 'cls-xi-rpl-1' },
  { label: 'XI TKJ 1', value: 'cls-xi-tkj-1' },
];

const loadDashboard = async () => {
  try {
    await reportingStore.loadDashboardStats();
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Gagal load dashboard: ' + err.message });
  }
};

const generateReport = async () => {
  try {
    await reportingStore.generateClassReport(selectedClassForReport.value);
    $q.notify({ type: 'positive', message: 'Laporan berhasil di-generate' });
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Gagal generate laporan: ' + err.message });
  }
};

const exportReport = async (format) => {
  try {
    await reportingStore.exportReport(format, `laporan_${selectedClassForReport.value}.${format.toLowerCase()}`);
    $q.notify({ type: 'positive', message: `Ekspor ${format} berhasil` });
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Gagal ekspor: ' + err.message });
  }
};

const getPredicateColor = (predicate) => {
  const map = { A: 'positive', B: 'blue', C: 'orange', D: 'negative', E: 'grey' };
  return map[predicate] || 'grey';
};

onMounted(loadDashboard);
</script>

<template>
  <div class="q-pa-md">
    <h6>Data Keuangan Siswa</h6>

    <!-- Pilih Kelas -->
    <div class="row q-col-gutter-md q-mb-md">
      <div class="col-md-4 col-sm-6 col-xs-12">
        <q-select v-model="selectedClass" :options="classOptions" label="Pilih Kelas" outlined dense
          @update:model-value="loadData" />
      </div>
      <div class="col-md-4 col-sm-6 col-xs-12">
        <q-btn color="primary" label="Refresh" @click="loadData" :loading="isLoading" />
      </div>
    </div>

    <!-- Stats Cards -->
    <div v-if="classSummaries.length" class="row q-col-gutter-md q-mb-md">
      <div class="col-md-3 col-sm-6 col-xs-12">
        <q-card flat bordered>
          <q-card-section>
            <div class="text-grey-6 text-caption">Total Siswa</div>
            <div class="text-h5">{{ classSummaries.length }}</div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-md-3 col-sm-6 col-xs-12">
        <q-card flat bordered>
          <q-card-section>
            <div class="text-grey-6 text-caption">Lunas</div>
            <div class="text-h5 text-positive">{{ paidStudents.length }}</div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-md-3 col-sm-6 col-xs-12">
        <q-card flat bordered>
          <q-card-section>
            <div class="text-grey-6 text-caption">Tunggakan</div>
            <div class="text-h5 text-negative">{{ overdueStudents.length }}</div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-md-3 col-sm-6 col-xs-12">
        <q-card flat bordered>
          <q-card-section>
            <div class="text-grey-6 text-caption">Total Tunggakan</div>
            <div class="text-h5 text-negative">Rp {{ totalBalance.toLocaleString() }}</div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Table -->
    <q-table :rows="classSummaries" :columns="columns" row-key="studentId" :loading="isLoading" flat bordered
      v-model:pagination="pagination">
      <template v-slot:body-cell-status="props">
        <q-td>
          <q-badge :color="props.row.isPaid() ? 'positive' : 'negative'">
            {{ props.row.status }}
          </q-badge>
        </q-td>
      </template>
      <template v-slot:body-cell-balance="props">
        <q-td>
          Rp {{ props.row.balance.toLocaleString() }}
        </q-td>
      </template>
      <template v-slot:body-cell-action="props">
        <q-td>
          <q-btn color="primary" label="Detail" size="sm" @click="viewDetail(props.row.studentId)" />
        </q-td>
      </template>
    </q-table>

    <!-- Detail Dialog -->
    <q-dialog v-model="showDetailDialog" persistent>
      <q-card style="min-width: 500px">
        <q-card-section>
          <div class="text-h6">Detail Tagihan</div>
          <div class="text-grey-6">{{ currentStudent?.studentName }}</div>
        </q-card-section>
        <q-card-section v-if="currentStudent">
          <div class="row q-col-gutter-sm">
            <div class="col-6">
              <div class="text-caption text-grey-6">Total Tagihan</div>
              <div class="text-subtitle1">Rp {{ currentStudent.totalInvoiced.toLocaleString() }}</div>
            </div>
            <div class="col-6">
              <div class="text-caption text-grey-6">Total Bayar</div>
              <div class="text-subtitle1">Rp {{ currentStudent.totalPaid.toLocaleString() }}</div>
            </div>
          </div>
          <div class="q-mt-sm">
            <div class="text-caption text-grey-6">Status</div>
            <q-badge :color="currentStudent.isPaid() ? 'positive' : 'negative'">
              {{ currentStudent.status }}
            </q-badge>
          </div>
          <q-separator class="q-my-md" />
          <div class="text-caption text-grey-6">Riwayat Tagihan</div>
          <q-list dense>
            <q-item v-for="inv in currentStudent.invoices" :key="inv.id">
              <q-item-section>
                <q-item-label>{{ inv.period }}</q-item-label>
                <q-item-label caption>Rp {{ inv.amount.toLocaleString() }} - Jatuh tempo: {{ inv.dueDate
                  }}</q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-badge :color="inv.isPaid() ? 'positive' : 'negative'">
                  {{ inv.status }}
                </q-badge>
              </q-item-section>
            </q-item>
          </q-list>
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Tutup" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useFinanceStore } from '@/stores/finance.store';
import { useQuasar } from 'quasar';

const $q = useQuasar();
const store = useFinanceStore();

const selectedClass = ref('cls-xi-rpl-1');
const classOptions = [
  { label: 'XI RPL 1', value: 'cls-xi-rpl-1' },
  { label: 'XI TKJ 1', value: 'cls-xi-tkj-1' },
  { label: 'XII MM 1', value: 'cls-xii-mm-1' },
];
const showDetailDialog = ref(false);
const currentStudent = ref(null);
const pagination = ref({ rowsPerPage: 10 });

const isLoading = computed(() => store.isLoading);
const classSummaries = computed(() => store.classSummaries);
const paidStudents = computed(() => store.paidStudents);
const overdueStudents = computed(() => store.overdueStudents);
const totalBalance = computed(() => store.totalBalance);

const columns = [
  { name: 'studentName', label: 'Nama Siswa', field: 'studentName', align: 'left' },
  { name: 'status', label: 'Status', field: 'status', align: 'center' },
  { name: 'balance', label: 'Tunggakan', field: 'balance', align: 'right' },
  { name: 'lastPaymentDate', label: 'Bayar Terakhir', field: 'lastPaymentDate', align: 'center' },
  { name: 'action', label: 'Aksi', field: 'action', align: 'center' },
];

const loadData = async () => {
  try {
    await store.fetchClassFinancialSummary(selectedClass.value);
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Gagal memuat data: ' + err.message });
  }
};

const viewDetail = async (studentId) => {
  try {
    const data = await store.fetchStudentFinancialDetail(studentId);
    currentStudent.value = data;
    showDetailDialog.value = true;
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Gagal memuat detail: ' + err.message });
  }
};

onMounted(loadData);
</script>

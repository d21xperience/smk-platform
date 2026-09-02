<template>
  <q-page padding>
    <!-- Header & Quick Stats -->
    <div class="q-mb-md">
      <h5 class="q-mb-sm text-primary">Data Siswa</h5>
      <div class="row q-gutter-md">
        <q-card class="col">
          <q-card-section>
            <div class="text-h6">{{ stats.total }}</div>
            <div class="text-caption text-grey-7">Total Siswa</div>
          </q-card-section>
        </q-card>
        <q-card class="col">
          <q-card-section>
            <div class="text-h6 text-positive">{{ stats.active }}</div>
            <div class="text-caption text-grey-7">Siswa Aktif</div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Main Data Card -->
    <q-card>
      <!-- Filter Section -->
      <q-card-section>
        <StudentFilterBar :current-search="filters.search" :current-status="filters.status"
          :current-jurusan="filters.jurusan" :current-tingkat="filters.tingkat" @update:search="setSearch"
          @update:status="setStatusFilter" @update:jurusan="setJurusanFilter" @update:tingkat="setTingkatFilter"
          @clear="clearFilters" />
      </q-card-section>

      <q-separator />

      <!-- Table Section -->
      <StudentTable :items="items" :loading="loading" @view="goToDetail" @edit="goToEdit" @mutate="handleMutate" />

      <q-separator />

      <!-- Pagination Section -->
      <StudentPagination :page="pagination.page" :total-pages="pagination.totalPages" :total="pagination.total"
        :limit="pagination.limit" :has-more="pagination.hasMore" @next="nextPage" @prev="prevPage"
        @go-to-page="goToPage" />
    </q-card>

    <!-- Floating Action Button for Create -->
    <q-page-sticky position="bottom-right" :offset="[18, 18]">
      <q-btn fab icon="add" color="primary" @click="goToCreate">
        <q-tooltip>Tambah Siswa</q-tooltip>
      </q-btn>
    </q-page-sticky>

    <!-- Global Loading Overlay -->
    <q-inner-loading :showing="loading" />
  </q-page>
</template>

<script setup>
import { watch } from 'vue';
import { useRouter } from 'vue-router';
import { useQuasar } from 'quasar';
import { useStudentList } from '@/composables/student/useStudentList';

import StudentFilterBar from '@/components/student/StudentFilterBar.vue';
import StudentTable from '@/components/student/StudentTable.vue';
import StudentPagination from '@/components/student/StudentPagination.vue';

const router = useRouter();
const $q = useQuasar();

// Destructure semua kebutuhan dari composable
const {
  items, stats, pagination, filters, loading, error,
  setSearch, setStatusFilter, setJurusanFilter, setTingkatFilter,
  clearFilters, nextPage, prevPage, goToPage, clearError
} = useStudentList();

// Watch error state untuk menampilkan notifikasi global
watch(error, (newError) => {
  if (newError) {
    $q.notify({
      type: 'negative',
      message: newError.message || 'Terjadi kesalahan saat memuat data.'
    });
    clearError();
  }
});

// Navigation Actions
function goToCreate() {
  router.push({ name: 'student-create' });
}

function goToDetail(studentId) {
  router.push({ name: 'student-detail', params: { id: studentId } });
}

function goToEdit(studentId) {
  router.push({ name: 'student-edit', params: { id: studentId } });
}

function handleMutate(studentId) {
  // Trigger notifikasi atau navigasi ke halaman mutasi dengan pre-filled data
  $q.notify({
    type: 'info',
    message: `Memulai proses mutasi untuk siswa ID: ${studentId}`
  });
  // Contoh navigasi jika ada halaman mutasi khusus:
  // router.push({ name: 'manajemen-mutasi-siswa', query: { studentId } });
}
</script>

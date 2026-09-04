<template>
  <q-page padding>
    <!-- Header -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <div class="text-h5 text-weight-bold">Data Siswa</div>
        <div class="text-grey-7">Kelola data siswa aktif, lulus, dan mutasi</div>
      </div>
      <q-btn color="primary" icon="add" label="Tambah Siswa" @click="$router.push({ name: 'student-create' })" />
    </div>

    <!-- Filters -->
    <q-card class="q-mb-md">
      <q-card-section>
        <div class="row q-gutter-md">
          <q-select v-model="localFilters.status" :options="statusOptions" label="Status" clearable emit-value
            map-options class="col-12 col-sm-3" @update:model-value="applyFilters" />
          <q-input v-model="localFilters.search" label="Cari NISN/NIS/Nama" clearable class="col-12 col-sm-6"
            @update:model-value="debouncedApplyFilters">
            <template v-slot:append>
              <q-icon name="search" />
            </template>
          </q-input>
          <q-btn color="secondary" icon="refresh" label="Reset" class="col-12 col-sm-2" @click="resetFilters" />
        </div>
      </q-card-section>
    </q-card>

    <!-- Table -->
    <q-card>
      <q-table :rows="students" :columns="columns" :loading="isLoading" :pagination="tablePagination"
        row-key="studentId" @request="onRequest">
        <template v-slot:body-cell-status="props">
          <q-td :props="props">
            <q-badge :color="getStatusColor(props.row.status)">
              {{ getStatusLabel(props.row.status) }}
            </q-badge>
          </q-td>
        </template>

        <template v-slot:body-cell-actions="props">
          <q-td :props="props">
            <q-btn flat dense round icon="visibility" color="primary" @click="viewDetail(props.row.studentId)">
              <q-tooltip>Detail</q-tooltip>
            </q-btn>
            <q-btn flat dense round icon="delete" color="negative" @click="openDeleteDialog(props.row)">
              <q-tooltip>Hapus</q-tooltip>
            </q-btn>
          </q-td>
        </template>

        <template v-slot:loading>
          <q-inner-loading showing color="primary" />
        </template>
      </q-table>
    </q-card>

    <!-- Delete Confirmation Dialog -->
    <q-dialog v-model="showDeleteDialog" persistent>
      <q-card style="width: 400px">
        <q-card-section>
          <div class="text-h6">Konfirmasi Hapus</div>
          <div class="q-mt-md">
            Apakah Anda yakin ingin menghapus siswa
            <strong>{{ selectedStudent?.fullName?.firstName }} {{ selectedStudent?.fullName?.lastName }}</strong>?
          </div>
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" color="grey" @click="showDeleteDialog = false" />
          <q-btn flat label="Hapus" color="negative" @click="confirmDelete" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Notification Dialog -->
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
import { onMounted } from 'vue'
import { useStudentList } from '@/composables/kesiswaan/useStudentList'

const {
  students,
  isLoading,
  localFilters,
  showDeleteDialog,
  showNotification,
  notificationMessage,
  notificationType,
  selectedStudent,
  statusOptions,
  columns,
  tablePagination,
  getStatusColor,
  getStatusLabel,
  applyFilters,
  debouncedApplyFilters,
  resetFilters,
  onRequest,
  viewDetail,
  openDeleteDialog,
  confirmDelete,
  loadData,
} = useStudentList()

onMounted(() => {
  loadData()
})
</script>

<template>
  <q-page padding>
    <!-- Header -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <div class="text-h5 text-weight-bold">Data Siswa</div>
        <div class="text-grey-7">Kelola data siswa aktif, lulus, dan mutasi</div>
      </div>
      <div class="q-gutter-sm">
        <q-btn color="green-7" icon="table_view" label="XLS" title="Export ke Excel" />
        <q-btn color="deep-orange-7" icon="print" label="PDF" title="Cetak PDF" />
        <!-- <q-btn color="primary" icon="add" label="Siswa" @click="$router.push({ name: 'student-create' })"
          title="Tambah Siswa" /> -->
        <q-btn color="primary" icon="add" label="Siswa" @click="showAddStudents = true" title="Tambah Siswa" />

      </div>
    </div>

    <!-- Filters -->
    <!-- <q-card class="q-mb-md">
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
</q-card> -->
    <!-- Panel Filter & Pencarian -->
    <q-card class="shadow-1 q-mb-md">
      <q-card-section class="q-py-sm">
        <div class="row q-col-gutter-sm items-center">
          <!-- Input Pencarian -->
          <div class="col-12 col-md-4">
            <q-input v-model="filterText" dense outlined placeholder="Cari Nama, NIS, atau NISN..." clearable>
              <template v-slot:prepend>
                <q-icon name="search" />
              </template>
            </q-input>
          </div>

          <!-- Filter Jurusan -->
          <div class="col-12 col-sm-6 col-md-3">
            <q-select v-model="selectedJurusan" :options="opsiJurusan" dense outlined
              label="Kompetensi Keahlian (Jurusan)" emit-value map-options />
          </div>

          <!-- Filter Tingkat -->
          <div class="col-12 col-sm-6 col-md-2">
            <q-select v-model="selectedTingkat" :options="opsiTingkat" dense outlined label="Tingkat Kelas" emit-value
              map-options />
          </div>

          <!-- Tombol Reset Filter -->
          <div class="col-12 col-md-3 flex justify-end">
            <q-btn flat color="grey-7" label="Reset Filter" icon="refresh" @click="resetFilters"
              :disable="!isFilterActive" />
          </div>
        </div>
      </q-card-section>
    </q-card>
    <!-- Table -->
    <q-card>
      <q-table :rows="students" :columns="columns" :loading="isLoading" :pagination="initialPagination"
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

    <!-- Add Students -->
    <q-dialog v-model="showAddStudents">
      <q-card style="min-width: 400px; max-width: 500px; border-radius: 12px;">

        <!-- Bagian Judul Dialog -->
        <q-card-section class="row items-center q-pb-none">
          <div class="text-h6 text-weight-bold text-primary">Metode Tambah Siswa</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <!-- Bagian Isi / Menu Pilihan -->
        <q-card-section class="q-pa-md">
          <q-list bordered separator class="rounded-borders text-grey-9">

            <!-- PILIHAN 1: MANUAL -->
            <q-expansion-item icon="person_add" label="Input Manual" header-class="text-weight-bold text-subtitle1"
              expand-icon-class="text-primary">
              <q-card class="bg-grey-1">
                <q-card-section class="q-py-xs">
                  <q-list dense>
                    <!-- Sub Pilihan: Siswa Baru -->
                    <q-item clickable v-close-popup @click="goToPage({ name: 'student-create' })">
                      <q-item-section avatar>
                        <q-icon name="fiber_new" color="green" />
                      </q-item-section>
                      <q-item-section>
                        <q-item-label class="text-weight-medium">Siswa Baru</q-item-label>
                        <q-item-label caption>Input data siswa baru kelas 10</q-item-label>
                      </q-item-section>
                    </q-item>

                    <!-- Sub Pilihan: Siswa Pindahan / Mutasi Masuk -->
                    <q-item clickable v-close-popup @click="goToPage('/siswa/mutasi-masuk')">
                      <q-item-section avatar>
                        <q-icon name="swap_horiz" color="orange" />
                      </q-item-section>
                      <q-item-section>
                        <q-item-label class="text-weight-medium">Siswa Pindahan</q-item-label>
                        <q-item-label caption>Langsung ke halaman siswa mutasi masuk</q-item-label>
                      </q-item-section>
                    </q-item>
                  </q-list>
                </q-card-section>
              </q-card>
            </q-expansion-item>

            <!-- PILIHAN 2: IMPORT -->
            <q-item clickable v-close-popup @click="goToPage({ name: 'student-import' })">
              <q-item-section avatar>
                <q-icon name="upload_file" color="primary" size="md" />
              </q-item-section>
              <q-item-section>
                <q-item-label class="text-weight-bold text-subtitle1">Import File (Excel / CSV)</q-item-label>
                <q-item-label caption>Unggah data siswa dalam jumlah banyak sekaligus</q-item-label>
              </q-item-section>
            </q-item>

          </q-list>
        </q-card-section>

      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useStudentList } from '@/composables/kesiswaan/useStudentList'
import { useRouter } from 'vue-router'
const showAddStudents = ref(false)
// State untuk Filter & Pencarian
const filterText = ref('')
const selectedJurusan = ref('Semua')
const selectedTingkat = ref('Semua')

const {
  students,
  isLoading,
  // localFilters,
  showDeleteDialog,
  showNotification,
  notificationMessage,
  notificationType,
  selectedStudent,
  // statusOptions,
  // tablePagination,
  getStatusColor,
  getStatusLabel,
  // applyFilters,
  // debouncedApplyFilters,
  resetFilters,
  onRequest,
  viewDetail,
  openDeleteDialog,
  confirmDelete,
  loadData,
} = useStudentList()
const router = useRouter() // Mengaktifkan fungsi navigasi halaman

// Fungsi untuk mengarahkan user ke halaman tujuan sesuai rute (route) Anda
const goToPage = (pageUrl) => {
  router.push(pageUrl)
}

// Konfigurasi Kolom Tabel Quasar
const columns = [
  { name: 'nis', align: 'left', label: 'NIS', field: 'nis', sortable: true },
  { name: 'nisn', align: 'left', label: 'NISN', field: 'nisn', sortable: true },
  { name: 'nama', align: 'left', label: 'Nama Lengkap', field: row => `${row.fullName?.firstName || ''} ${row.fullName?.lastName || ''}`.trim(), sortable: true },
  { name: 'gender', label: 'JK', field: row => row.gender === 'MALE' ? 'L' : 'P', align: 'center' },
  { name: 'jurusan', align: 'left', label: 'Jurusan', field: 'jurusan', sortable: true },
  { name: 'kelas', align: 'left', label: 'Kelas', field: 'kelas', sortable: true },
  { name: 'status', align: 'center', label: 'Status', field: 'status' },
  { name: 'actions', label: 'Aksi', field: 'actions', align: 'center' },

]

// Pagination Awal q-table
const initialPagination = {
  sortBy: 'nama',
  descending: false,
  page: 1,
  rowsPerPage: 10
}


onMounted(() => {
  loadData()
})
</script>

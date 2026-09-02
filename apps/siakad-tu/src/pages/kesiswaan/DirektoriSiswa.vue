<template>
  <q-page class="q-pa-md bg-grey-2">
    <q-inner-loading :showing="loading" />

    <!-- Header Halaman-->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Direktori Siswa</h1>
        <p class="text-caption text-grey-7 q-mb-none">Manajemen data komputerisasi Buku Klaper dan data aktif siswa SMK
        </p>
      </div>
      <div>
        <q-btn color="primary" icon="add" label="Tambah Siswa" @click="openAddDialog" />
        <q-btn color="green-8" icon="file_download" label="Ekspor Excel" class="q-ml-sm" flat @click="exportToExcel" />
      </div>
    </div>

    <!-- Panel Filter & Pencarian-->
    <q-card class="shadow-1 q-mb-md">
      <q-card-section class="q-py-sm">
        <div class="row q-col-gutter-sm items-center">
          <div class="col-12 col-md-4">
            <q-input v-model="localSearch" dense outlined placeholder="Cari Nama, NIS, atau NISN..." clearable
              @update:model-value="debouncedSearch">
              <template v-slot:prepend>
                <q-icon name="search" />
              </template>
            </q-input>
          </div>
          <div class="col-12 col-sm-6 col-md-3">
            <q-select v-model="filters.jurusan" :options="opsiJurusan" dense outlined label="Kompetensi Keahlian"
              emit-value map-options @update:model-value="applyFilter('jurusan', $event)" />
          </div>
          <div class="col-12 col-sm-6 col-md-2">
            <q-select v-model="filters.tingkat" :options="opsiTingkat" dense outlined label="Tingkat Kelas" emit-value
              map-options @update:model-value="applyFilter('tingkat', $event)" />
          </div>
          <div class="col-12 col-md-3 flex justify-end">
            <q-btn flat color="grey-7" label="Reset Filter" icon="refresh" @click="handleClearFilters"
              :disable="!isFilterActive" />
          </div>
        </div>
      </q-card-section>
    </q-card>

    <!-- Tabel Utama Data Siswa-->
    <q-card class="shadow-1">
      <q-card-section class="q-pa-none">
        <q-table flat :rows="items" :columns="columns" row-key="nisn" :loading="loading" selection="multiple"
          v-model:selected="selectedRows" :pagination="initialPagination">
          <template v-slot:body-cell-status="props">
            <q-td :props="props">
              <q-badge :color="props.row.status === 'AKTIF' ? 'green' : 'orange'" :label="props.row.status" />
            </q-td>
          </template>
          <template v-slot:body-cell-aksi="props">
            <q-td :props="props" class="q-gutter-xs text-center">
              <q-btn size="sm" color="blue" icon="visibility" round flat @click="viewDetail(props.row)">
                <q-tooltip>Buku Induk Lengkap</q-tooltip>
              </q-btn>
              <q-btn size="sm" color="amber-9" icon="edit" round flat @click="editSiswa(props.row)">
                <q-tooltip>Ubah Data</q-tooltip>
              </q-btn>
              <q-btn size="sm" color="red" icon="person_remove" round flat @click="confirmDelete(props.row)">
                <q-tooltip>Proses Keluar/ Mutasi</q-tooltip>
              </q-btn>
            </q-td>
          </template>
        </q-table>
      </q-card-section>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useQuasar } from 'quasar';
import { useStudentList } from '@/composables/student/useStudentList';

const $q = useQuasar();
const { items, loading, filters, clearFilters, setJurusanFilter, setTingkatFilter, setSearch } = useStudentList();

const selectedRows = ref([]);
const localSearch = ref('');
let searchTimeout = null;

const opsiJurusan = [
  { label: 'Semua Jurusan', value: 'Semua' },
  { label: 'Rekayasa Perangkat Lunak (RPL)', value: 'RPL' },
  { label: 'Teknik Kendaraan Ringan (TKRO)', value: 'TKRO' },
  { label: 'Akuntansi & Keuangan Lembaga (AKL)', value: 'AKL' }
];

const opsiTingkat = [
  { label: 'Semua Tingkat', value: 'Semua' },
  { label: 'Kelas X', value: 'X' },
  { label: 'Kelas XI', value: 'XI' },
  { label: 'Kelas XII', value: 'XII' }
];

const columns = [
  { name: 'nis', align: 'left', label: 'NIS', field: 'nis', sortable: true },
  { name: 'nisn', align: 'left', label: 'NISN', field: 'nisn', sortable: true },
  { name: 'nama', align: 'left', label: 'Nama Lengkap', field: 'nama', sortable: true },
  { name: 'jk', align: 'center', label: 'JK', field: 'jenisKelamin' },
  { name: 'jurusan', align: 'left', label: 'Jurusan', field: 'jurusan', sortable: true },
  { name: 'kelas', align: 'left', label: 'Kelas', field: 'kelas', sortable: true },
  { name: 'status', align: 'center', label: 'Status', field: 'status' },
  { name: 'aksi', align: 'center', label: 'Aksi' }
];

const initialPagination = { sortBy: 'nama', descending: false, page: 1, rowsPerPage: 10 };

const isFilterActive = computed(() => {
  return localSearch.value !== '' || filters.value.jurusan !== 'Semua' || filters.value.tingkat !== 'Semua';
});

const debouncedSearch = () => {
  clearTimeout(searchTimeout);
  searchTimeout = setTimeout(() => {
    setSearch(localSearch.value);
  }, 500);
};

const applyFilter = (key, value) => {
  if (key === 'jurusan') setJurusanFilter(value);
  if (key === 'tingkat') setTingkatFilter(value);
};

const handleClearFilters = () => {
  localSearch.value = '';
  clearFilters();
};

const openAddDialog = () => {
  $q.notify({ color: 'primary', message: 'Fitur Tambah Siswa Baru (Form PPDB)', icon: 'add' });
};

const viewDetail = (siswa) => {
  $q.notify({ color: 'info', message: `Membuka detail Buku Induk: ${siswa.nama}`, icon: 'visibility' });
};

const editSiswa = (siswa) => {
  $q.notify({ color: 'warning', message: `Ubah data siswa: ${siswa.nama}`, icon: 'edit', textColor: 'dark' });
};

const confirmDelete = (siswa) => {
  $q.dialog({
    title: 'Konfirmasi Mutasi Keluar',
    message: `Apakah Anda yakin ingin memproses status mutasi/keluar untuk siswa ${siswa.nama}?`,
    cancel: true,
    persistent: true
  }).onOk(() => {
    $q.notify({ color: 'negative', message: `${siswa.nama} berhasil diproses keluar.`, icon: 'person_remove' });
  });
};

const exportToExcel = () => {
  $q.notify({ color: 'green-8', message: 'Mengekspor data ke format .xlsx...', icon: 'file_download' });
};
</script>

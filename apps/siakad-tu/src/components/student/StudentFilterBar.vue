<template>
  <div class="row q-col-gutter-sm items-center">
    <!-- Search Input -->
    <div class="col-12 col-md-4">
      <q-input :model-value="currentSearch" @update:model-value="handleSearchInput" dense outlined
        placeholder="Cari Nama, NIS, atau NISN..." clearable @clear="handleClearSearch">
        <template v-slot:prepend>
          <q-icon name="search" />
        </template>
      </q-input>
    </div>

    <!-- Status Filter -->
    <div class="col-12 col-sm-6 col-md-2">
      <q-select :model-value="currentStatus" @update:model-value="$emit('update:status', $event)"
        :options="statusOptions" dense outlined label="Status" emit-value map-options clearable />
    </div>

    <!-- Jurusan Filter -->
    <div class="col-12 col-sm-6 col-md-2">
      <q-select :model-value="currentJurusan" @update:model-value="$emit('update:jurusan', $event)"
        :options="jurusanOptions" dense outlined label="Jurusan" emit-value map-options clearable />
    </div>

    <!-- Tingkat Filter -->
    <div class="col-12 col-sm-6 col-md-2">
      <q-select :model-value="currentTingkat" @update:model-value="$emit('update:tingkat', $event)"
        :options="tingkatOptions" dense outlined label="Tingkat" emit-value map-options clearable />
    </div>

    <!-- Reset Button -->
    <div class="col-12 col-md-2 flex justify-end">
      <q-btn flat color="grey-7" label="Reset" icon="refresh" @click="$emit('clear')" :disable="isDefault" />
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  currentSearch: { type: String, default: '' },
  currentStatus: { type: String, default: '' },
  currentJurusan: { type: String, default: 'Semua' },
  currentTingkat: { type: String, default: 'Semua' }
});

const emit = defineEmits(['update:search', 'update:status', 'update:jurusan', 'update:tingkat', 'clear']);

// Opsi filter didefinisikan di komponen UI karena ini adalah concern presentasi
const statusOptions = [
  { label: 'Semua Status', value: '' },
  { label: 'Aktif', value: 'AKTIF' },
  { label: 'Mutasi Keluar', value: 'MUTASI_KELUAR' },
  { label: 'Lulus', value: 'LULUS' },
  { label: 'DO', value: 'DO' }
];

const jurusanOptions = [
  { label: 'Semua Jurusan', value: 'Semua' },
  { label: 'RPL', value: 'RPL' },
  { label: 'TKRO', value: 'TKRO' },
  { label: 'AKL', value: 'AKL' },
  { label: 'ULP', value: 'ULP' },
  { label: 'DPIB', value: 'DPIB' }
];

const tingkatOptions = [
  { label: 'Semua Tingkat', value: 'Semua' },
  { label: 'Kelas X', value: 'X' },
  { label: 'Kelas XI', value: 'XI' },
  { label: 'Kelas XII', value: 'XII' }
];

const isDefault = computed(() => {
  return !props.currentSearch &&
    !props.currentStatus &&
    props.currentJurusan === 'Semua' &&
    props.currentTingkat === 'Semua';
});

// Debounce search input untuk menghindari spam request ke adapter/backend
let searchTimeout = null;
function handleSearchInput(val) {
  clearTimeout(searchTimeout);
  searchTimeout = setTimeout(() => {
    emit('update:search', val || '');
  }, 400);
}

function handleClearSearch() {
  clearTimeout(searchTimeout);
  emit('update:search', '');
}
</script>

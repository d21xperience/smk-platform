<template>
  <q-table flat :rows="items" :columns="columns" row-key="id" :loading="loading" hide-pagination dense>
    <!-- Custom Column: Jenis Kelamin -->
    <template v-slot:body-cell-jenisKelamin="props">
      <q-td :props="props">
        <q-icon :name="props.row.jenisKelamin === 'L' ? 'male' : 'female'"
          :color="props.row.jenisKelamin === 'L' ? 'blue' : 'pink'" />
      </q-td>
    </template>

    <!-- Custom Column: Status -->
    <template v-slot:body-cell-status="props">
      <q-td :props="props">
        <q-badge :color="getStatusColor(props.row.status)" :label="props.row.status" />
      </q-td>
    </template>

    <!-- Custom Column: Aksi -->
    <template v-slot:body-cell-aksi="props">
      <q-td :props="props" class="q-gutter-xs text-center">
        <q-btn size="sm" color="blue" icon="visibility" round flat @click="$emit('view', props.row.id)">
          <q-tooltip>Detail</q-tooltip>
        </q-btn>
        <q-btn size="sm" color="amber-9" icon="edit" round flat @click="$emit('edit', props.row.id)">
          <q-tooltip>Ubah Data</q-tooltip>
        </q-btn>
        <q-btn size="sm" color="red" icon="person_remove" round flat @click="$emit('mutate', props.row.id)">
          <q-tooltip>Mutasi</q-tooltip>
        </q-btn>
      </q-td>
    </template>

    <!-- Empty State -->
    <template v-slot:no-data>
      <div class="full-width row flex-center text-grey-6 q-py-xl">
        <q-icon size="50px" name="search_off" class="q-mr-md" />
        <div class="text-h6">Tidak ada data siswa yang cocok dengan filter.</div>
      </div>
    </template>
  </q-table>
</template>

<script setup>
defineProps({
  items: { type: Array, required: true },
  loading: { type: Boolean, default: false }
});

defineEmits(['view', 'edit', 'mutate']);

// Definisi kolom standar untuk entitas Student
const columns = [
  { name: 'nis', align: 'left', label: 'NIS', field: 'nis', sortable: true },
  { name: 'nisn', align: 'left', label: 'NISN', field: 'nisn', sortable: true },
  { name: 'nama', align: 'left', label: 'Nama Lengkap', field: 'nama', sortable: true },
  { name: 'jenisKelamin', align: 'center', label: 'JK', field: 'jenisKelamin' },
  { name: 'jurusan', align: 'left', label: 'Jurusan', field: 'jurusan', sortable: true },
  { name: 'kelas', align: 'left', label: 'Kelas', field: 'kelas', sortable: true },
  { name: 'status', align: 'center', label: 'Status', field: 'status' },
  { name: 'aksi', align: 'center', label: 'Aksi', field: 'aksi', sortable: false }
];

function getStatusColor(status) {
  if (status === 'AKTIF') return 'green';
  if (status === 'MUTASI_KELUAR') return 'orange';
  if (status === 'LULUS') return 'blue';
  if (status === 'DO') return 'red';
  return 'grey';
}
</script>

<!-- FILE: src/pages/homeroom-student/HomeroomStudentListPage.vue -->
<!-- STATUS: NEW -->
<!-- STATUS IMPLEMENTASI: COMPLETE -->

<template>
  <q-page class="q-pa-md">
    <div class="text-h5 q-mb-md">Daftar Siswa Wali Kelas</div>

    <q-banner v-if="error" class="bg-negative text-white q-mb-md">
      {{ error }}
    </q-banner>

    <q-table flat bordered :rows="students" :columns="columns" row-key="id" :loading="loading"
      :pagination="{ rowsPerPage: 10 }">
      <template #body-cell-photoUrl="props">
        <q-td :props="props">
          <q-avatar size="40px">
            <img :src="props.value || 'https://via.placeholder.com/150'" alt="Foto Siswa" />
          </q-avatar>
        </q-td>
      </template>

      <template #body-cell-gender="props">
        <q-td :props="props">
          {{ props.value === 'L' ? 'Laki-laki' : 'Perempuan' }}
        </q-td>
      </template>

      <template #body-cell-status="props">
        <q-td :props="props">
          <q-badge :color="getStatusColor(props.value)">
            {{ props.value }}
          </q-badge>
        </q-td>
      </template>

      <template #body-cell-parentContact="props">
        <q-td :props="props">
          <q-btn flat dense icon="call" :label="props.value" :href="`tel:${props.value}`" color="primary" />
        </q-td>
      </template>
    </q-table>
  </q-page>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useHomeroomStudent } from '@/composables/useHomeroomStudent'
import { STUDENT_STATUS } from '@/contracts/homeroomStudentContract'

const { students, loading, error, loadStudents } = useHomeroomStudent()

const columns = ref([
  { name: 'seatNumber', label: 'No Absen', field: 'seatNumber', align: 'center', sortable: true },
  { name: 'photoUrl', label: 'Foto', field: 'photoUrl', align: 'center' },
  { name: 'nisn', label: 'NISN', field: 'nisn', align: 'left' },
  { name: 'nis', label: 'NIS', field: 'nis', align: 'left' },
  { name: 'fullName', label: 'Nama Lengkap', field: 'fullName', align: 'left', sortable: true },
  { name: 'gender', label: 'JK', field: 'gender', align: 'center' },
  { name: 'status', label: 'Status', field: 'status', align: 'center' },
  { name: 'parentContact', label: 'Kontak Wali', field: 'parentContact', align: 'left' }
])

/**
 * Menentukan warna badge berdasarkan status siswa.
 */
function getStatusColor(status) {
  switch (status) {
    case STUDENT_STATUS.ACTIVE: return 'positive'
    case STUDENT_STATUS.INACTIVE: return 'grey'
    case STUDENT_STATUS.TRANSFERRED: return 'warning'
    case STUDENT_STATUS.GRADUATED: return 'info'
    default: return 'primary'
  }
}

onMounted(() => {
  loadStudents()
})
</script>

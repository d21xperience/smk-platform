<template>
  <div class="q-pa-md">
    <!-- Header -->
    <div class="row items-center q-mb-md">
      <div class="col">
        <q-btn icon="arrow_back" flat round @click="$router.push('/')" />
        <span class="text-h6">{{ classDisplayName }}</span>
        <q-chip size="sm" color="primary" text-color="white" class="q-ml-sm">
          {{ totalStudents }} Siswa
        </q-chip>
      </div>
      <div class="col-auto">
        <q-btn icon="refresh" flat round :loading="isLoading" @click="fetchRoster" />
      </div>
    </div>

    <!-- Rekap Cepat (Selector) -->
    <div class="row q-col-gutter-md q-mb-md" v-if="!isLoading && roster.length > 0">
      <div class="col-4">
        <q-card flat bordered class="text-center">
          <q-card-section>
            <div class="text-h5 text-positive">{{ presentCount }}</div>
            <div class="text-caption">Hadir</div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-4">
        <q-card flat bordered class="text-center">
          <q-card-section>
            <div class="text-h5 text-warning">{{ sickCount + permitCount }}</div>
            <div class="text-caption">Sakit / Izin</div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-4">
        <q-card flat bordered class="text-center">
          <q-card-section>
            <div class="text-h5 text-negative">{{ alphaCount }}</div>
            <div class="text-caption">Alpha</div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Loading -->
    <q-spinner v-if="isLoading" color="primary" size="3em" class="q-mt-xl" />

    <!-- Error -->
    <q-banner v-else-if="error" class="bg-negative text-white q-mt-md">
      {{ error }}
    </q-banner>

    <!-- Tabel Daftar Siswa -->
    <q-table
      v-else
      :rows="roster"
      :columns="columns"
      row-key="student.id"
      flat
      bordered
      :pagination="{ rowsPerPage: 0 }"
      hide-pagination
    >
      <!-- Kolom Nama -->
      <template #body-cell-name="props">
        <q-td :props="props">
          <div>
            <div class="text-weight-bold">{{ props.row.student.name }}</div>
            <div class="text-caption text-grey">NISN: {{ props.row.student.nisn }}</div>
          </div>
        </q-td>
      </template>

      <!-- Kolom Status Absensi -->
      <template #body-cell-status="props">
        <q-td :props="props">
          <AttendanceButton
            :model-value="props.row.attendance.status"
            @update:model-value="(val) => handleStatusChange(props.row.student.id, val)"
          />
        </q-td>
      </template>

      <!-- Kolom Waktu (jika perlu) -->
      <template #body-cell-time="props">
        <q-td :props="props">
          <div v-if="props.row.attendance.checkInTime" class="text-caption">
            <q-icon name="login" size="xs" /> {{ props.row.attendance.checkInTime }}
            <span v-if="props.row.attendance.checkOutTime">
              <q-icon name="logout" size="xs" class="q-ml-sm" />
              {{ props.row.attendance.checkOutTime }}
            </span>
          </div>
          <span v-else class="text-grey">-</span>
        </q-td>
      </template>
    </q-table>
  </div>
</template>

<script setup>
import { computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useAttendance } from '../../../composables/useAttendance'
import AttendanceButton from '../../../components/attendance/AttendanceButton.vue'

const route = useRoute()
const classId = route.params.id

const { roster, isLoading, error, loadRoster, submit, presentCount, alphaCount } = useAttendance()

// Selector tambahan (sick & permit) karena kita hanya punya present & alpha di selector, kita tambahkan di sini untuk demo
const sickCount = computed(() => roster.value.filter((r) => r.attendance.status === 'SICK').length)
const permitCount = computed(
  () => roster.value.filter((r) => r.attendance.status === 'PERMIT').length,
)
const totalStudents = computed(() => roster.value.length)

// Nama kelas dari data pertama (jika ada)
const classDisplayName = computed(() => {
  if (roster.value.length > 0 && roster.value[0].student) {
    // Kita asumsikan class name dikirim via props atau kita ambil dari route meta, untuk demo kita hardcode 'X TKJ 1'
    // Lebih baik kita simpan di store atau router. Untuk sekarang, hardcode.
    return 'X TKJ 1'
  }
  return 'Detail Kelas'
})

// Definisi kolom untuk q-table
const columns = [
  {
    name: 'name',
    label: 'Nama Siswa',
    field: (row) => row.student.name,
    align: 'left',
    sortable: true,
  },
  { name: 'status', label: 'Status Kehadiran', align: 'center' },
  { name: 'time', label: 'Jam (In / Out)', align: 'center' },
]

const fetchRoster = () => loadRoster(classId)

// Handler saat tombol status di-klik
const handleStatusChange = async (studentId, newStatus) => {
  // Jika status null (toggle off), kita skip submit? Atau kita set ke null di store.
  // Untuk MVP, jika status null, kita anggap belum diisi, tidak perlu submit.
  if (!newStatus) {
    // Update lokal state ke null (sudah ditangani oleh v-model di komponen)
    // Tapi kita perlu update store agar rekap berubah.
    // Cara cepat: langsung update via store (tapi hati-hati, store action seharusnya yang handle)
    // Karena komponen ini hanya menerima event, kita akan trigger submit hanya jika status ada.
    return
  }

  try {
    // Kirim ke server (mock)
    await submit({
      studentId,
      classId,
      status: newStatus,
      checkInTime: newStatus === 'PRESENT' || newStatus === 'LATE' ? '07:00' : null,
      checkOutTime: newStatus === 'PRESENT' ? '14:30' : null,
      notes: '',
    })
    // Berhasil, state di store sudah diupdate oleh action (updateStudentStatus)
  } catch (err) {
    // Tampilkan notifikasi error
    console.error('Gagal submit absensi:', err)
    // Reload data untuk rollback state
    fetchRoster()
  }
}

// Refresh awal
onMounted(fetchRoster)

// Jika classId berubah (misal pindah halaman), refresh
watch(() => route.params.id, fetchRoster)
</script>

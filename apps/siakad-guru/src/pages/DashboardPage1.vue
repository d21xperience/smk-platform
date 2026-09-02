<template>
  <div class="q-pa-md">
    <h5>Dashboard Guru</h5>
    <p class="text-grey-6">Selamat datang, {{ user?.name || 'Guru' }}</p>

    <!-- Ringkasan -->
    <div class="row q-col-gutter-md q-mb-md">
      <div class="col-md-3 col-sm-6 col-xs-12">
        <q-card flat bordered>
          <q-card-section>
            <div class="text-grey-6 text-caption">Sesi Hari Ini</div>
            <div class="text-h5">{{ todaySchedule.length }}</div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-md-3 col-sm-6 col-xs-12">
        <q-card flat bordered>
          <q-card-section>
            <div class="text-grey-6 text-caption">Sesi Aktif</div>
            <div class="text-h5">{{ hasActiveSession ? '1' : '0' }}</div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-md-3 col-sm-6 col-xs-12">
        <q-card flat bordered>
          <q-card-section>
            <div class="text-grey-6 text-caption">Absensi Hari Ini</div>
            <div class="text-h5">{{ attendanceCount }}</div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Jadwal Hari Ini -->
    <q-card flat bordered>
      <q-card-section>
        <div class="text-h6">Jadwal Hari Ini</div>
      </q-card-section>
      <q-card-section>
        <q-table :rows="todaySchedule" :columns="scheduleColumns" row-key="id" :loading="teachingLoading" flat dense>
          <template v-slot:body-cell-status="props">
            <q-td>
              <q-badge :color="getStatusColor(props.row.status)">
                {{ props.row.status }}
              </q-badge>
            </q-td>
          </template>
          <template v-slot:body-cell-action="props">
            <q-td>
              <q-btn v-if="props.row.status === 'draft'" color="primary" size="sm" label="Mulai"
                @click="startSession(props.row.id)" :loading="sessionLoading" />
              <q-btn v-else-if="props.row.status === 'active'" color="negative" size="sm" label="Akhiri"
                @click="endSession(props.row.id)" :loading="sessionLoading" />
              <q-btn v-else color="grey" size="sm" label="Selesai" disabled />
            </q-td>
          </template>
        </q-table>
      </q-card-section>
    </q-card>

    <!-- Jika ada sesi aktif, tampilkan tombol ke absensi -->
    <div v-if="hasActiveSession" class="q-mt-md">
      <q-btn color="secondary" label="Buka Absensi" icon="checklist" @click="goToAttendance" />
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import { useTeachingStore } from '@/stores/teaching.store';
import { useAttendanceStore } from '@/stores/attendanceStore';
import { useAuth } from '@/composables/useAuth';
// import { useOperationalContext } from '@/composables/context/useOperationalContext';
import { date } from 'quasar';

const router = useRouter();
const { user } = useAuth();
// const { context } = useOperationalContext();
const teachingStore = useTeachingStore();
const attendanceStore = useAttendanceStore();

const sessionLoading = ref(false);
const today = date.formatDate(new Date(), 'YYYY-MM-DD');

// Computed properties dari store
const todaySchedule = computed(() => teachingStore.todaySchedule || []);
const teachingLoading = computed(() => teachingStore.isLoading);
const hasActiveSession = computed(() => teachingStore.hasActiveSession);

// Computed untuk attendance count
const attendanceCount = computed(() => {
  return attendanceStore.records?.length || 0;
});

// Table columns
const scheduleColumns = [
  { name: 'time', label: 'Jam', field: row => `${row.startTime} - ${row.endTime}`, align: 'left' },
  { name: 'class', label: 'Kelas', field: 'classId', align: 'left' },
  { name: 'subject', label: 'Mata Pelajaran', field: 'subjectId', align: 'left' },
  { name: 'status', label: 'Status', field: 'status', align: 'center' },
  { name: 'action', label: 'Aksi', field: 'action', align: 'center' },
];

const getStatusColor = (status) => {
  const map = {
    draft: 'grey',
    active: 'positive',
    completed: 'blue',
    cancelled: 'negative',
  };
  return map[status] || 'grey';
};

// Actions
const startSession = async (sessionId) => {
  sessionLoading.value = true;
  try {
    await teachingStore.startSession(sessionId);
    // Refresh jadwal
    await teachingStore.fetchTodaySchedule(user.value.id, today);
  } catch (err) {
    console.error('Gagal memulai sesi:', err);
  } finally {
    sessionLoading.value = false;
  }
};

const endSession = async (sessionId) => {
  sessionLoading.value = true;
  try {
    await teachingStore.endSession(sessionId);
    await teachingStore.fetchTodaySchedule(user.value.id, today);
  } catch (err) {
    console.error('Gagal mengakhiri sesi:', err);
  } finally {
    sessionLoading.value = false;
  }
};

const goToAttendance = () => {
  const active = teachingStore.activeSession;
  if (active) {
    router.push(`/attendance/${active.id}`);
  }
};

// Lifecycle
onMounted(async () => {
  if (user.value) {
    try {
      await teachingStore.fetchTodaySchedule(user.value.id, today);
    } catch (err) {
      console.error('Gagal memuat jadwal:', err);
    }
  }
});
</script>

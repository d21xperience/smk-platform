<template>
  <q-page class="q-pa-md">
    <!-- HEADER HALAMAN -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <div class="text-h5 text-weight-bold text-grey-9">Jadwal Mengajar</div>
        <div class="text-caption text-grey-7">
          Periode Aktif:
          <q-chip dense outline color="primary" icon="event">
            <!-- {{ academicStore.activePeriodCode }} -->
          </q-chip>
        </div>
      </div>
      <!-- <q-btn flat round color="primary" icon="refresh" :loading="loading" @click="fetchMySchedule">
        <q-tooltip>Segarkan Jadwal</q-tooltip>
      </q-btn> -->
    </div>

    <!-- LOADING STATE -->
    <div v-if="loading && scheduleIsEmpty" class="text-center q-pa-xl">
      <q-spinner-dots color="primary" size="40px" />
      <div class="text-grey-7 q-mt-sm">Memuat jadwal...</div>
    </div>

    <!-- EMPTY STATE -->
    <q-card flat v-else-if="scheduleIsEmpty" class="text-center q-pa-xl">
      <q-icon name="event_busy" size="64px" color="grey-5" />
      <div class="text-subtitle1 text-grey-7 q-mt-md">
        Tidak ada jadwal mengajar untuk periode ini.
      </div>
    </q-card>

    <!-- JADWAL BERDASARKAN HARI -->
    <div v-else class="q-col-gutter-md">
      <div v-for="(classes, day) in schedule" :key="day" v-show="classes.length > 0">
        <!-- Header Hari -->
        <div class="text-subtitle1 text-weight-bold text-primary q-mb-sm row items-center">
          <q-icon name="calendar_today" size="sm" class="q-mr-xs" />
          {{ day }}
          <q-separator class="q-ml-sm col" />
        </div>

        <!-- Kartu Jadwal per Jam Pelajaran -->
        <div class="q-col-gutter-sm">
          <!-- <q-card v-for="(cls, index) in classes" :key="index" flat bordered class="schedule-card">
            <q-card-section class="row items-center no-wrap q-pa-sm"> -->
          <!-- Waktu -->
          <!-- <div class="col-auto text-center q-pr-md" style="min-width: 70px;">
                <div class="text-h6 text-weight-bold text-primary q-mb-none">{{ cls.start_time }}</div>
                <div class="text-caption text-grey-7">s.d. {{ cls.end_time }}</div>
              </div> -->

          <q-separator vertical class="q-mx-sm" />

          <!-- Detail Kelas & Mapel -->
          <!-- <div class="col">
                <div class="text-subtitle1 text-weight-bold text-grey-9">{{ cls.class_name }}</div>
                <div class="text-caption text-grey-7 row items-center q-gutter-x-sm">
                  <span><q-icon name="menu_book" size="xs" /> {{ cls.subject }}</span>
                  <span v-if="cls.room"><q-icon name="meeting_room" size="xs" /> {{ cls.room }}</span>
                </div>
              </div> -->

          <!-- Aksi -->
          <!-- <div class="col-auto">
                <q-btn flat dense color="primary" icon="fact_check" label="Absen"
                  @click="goToAttendance(cls.class_id, cls.class_name)" />
              </div>
            </q-card-section>
          </q-card> -->
        </div>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useTeaching } from '@/composables/useTeaching'
// import { useContext } from '@/composables/useContext'

// const { operational } = useContext()
const {
  schedule,
  loading,
  // sessions,
  loadSchedule,
  loadSessions,
  // startSession,
  // completeSession,
  // activeSessionId
} = useTeaching()

const selectedDate = ref(new Date().toISOString().slice(0, 10))
const scheduleIsEmpty = ref([])
onMounted(async () => {
  const tes = await loadSchedule()
  console.log(tes)
  await loadSessions(selectedDate.value)
})

// function mulaiSesi(id) {
//   startSession(id)
// }

// function akhiriSesi(id) {
//   completeSession(id)
// }


</script>

<style scoped>
.schedule-card {
  transition: transform 0.2s, box-shadow 0.2s;
}

.schedule-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
}
</style>

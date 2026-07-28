<!-- src/pages/siakad/teaching/AttendancePage.vue -->
<template>
  <q-page class="attendance-page">
    <!-- Header -->
    <q-toolbar class="bg-primary text-white">
      <q-toolbar-title>Absensi Kelas</q-toolbar-title>
      <q-btn flat round dense icon="refresh" @click="controller.retry" />
    </q-toolbar>

    <!-- Loading -->
    <div v-if="controller.isLoading.value" class="flex flex-center full-height">
      <q-spinner size="3em" color="primary" />
      <p class="q-mt-md">Memuat halaman...</p>
    </div>

    <!-- Error -->
    <div v-else-if="controller.hasError.value" class="flex flex-center full-height column">
      <q-icon name="error" size="4em" color="negative" />
      <p class="q-mt-md text-negative">{{ controller.errorMessage.value }}</p>
      <q-btn color="primary" label="Coba Lagi" @click="controller.retry" />
    </div>

    <!-- Ready -->
    <div v-else-if="controller.isReady.value" class="q-pa-md">
      <h6>Teaching Session: {{ controller.teachingSession.value?.name }}</h6>
      <!-- Di sini nanti komponen daftar absensi -->
      <q-list bordered>
        <q-item v-for="(item, idx) in controller.attendanceList.value" :key="idx">
          <q-item-section>{{ item }}</q-item-section>
        </q-item>
      </q-list>
      <q-btn
        v-if="!controller.attendanceList.value.length"
        label="Belum ada data absensi"
        color="grey"
        disable
      />
    </div>
  </q-page>
</template>

<script setup>
import { onMounted } from 'vue'
import { useAttendancePage } from './useAttendancePage'

// Controller halaman
const controller = useAttendancePage()

// Satu-satunya lifecycle yang dipanggil dari template
onMounted(() => {
  controller.initialize()
})
</script>

<style scoped>
.full-height {
  min-height: 60vh;
}
</style>

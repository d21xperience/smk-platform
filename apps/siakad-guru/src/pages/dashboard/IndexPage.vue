<template>
  <div class="q-pa-md">
    <div class="row items-center q-mb-md">
      <div class="col">
        <h6 class="q-my-none">📚 Kelas yang Diajar</h6>
      </div>
      <div class="col-auto">
        <q-btn icon="refresh" flat round :loading="isLoading" @click="loadClasses" />
      </div>
    </div>

    <!-- Loading State -->
    <q-spinner v-if="isLoading" color="primary" size="3em" class="q-mt-xl" />

    <!-- Error State -->
    <q-banner v-else-if="error" class="bg-negative text-white q-mt-md">
      {{ error }}
    </q-banner>

    <!-- List Kelas -->
    <q-list v-else bordered separator class="q-mt-sm rounded-borders">
      <q-item v-for="cls in classes" :key="cls.id" clickable v-ripple :to="`/class/${cls.id}`">
        <q-item-section avatar>
          <q-icon name="class" color="primary" />
        </q-item-section>
        <q-item-section>
          <q-item-label class="text-h6">{{ cls.displayName }}</q-item-label>
          <q-item-label caption>
            Wali Kelas: {{ teacherName }} | Ruang: {{ cls.roomId || 'Belum diatur' }}
          </q-item-label>
        </q-item-section>
        <q-item-section side>
          <q-icon name="chevron_right" color="grey" />
        </q-item-section>
      </q-item>
    </q-list>

    <!-- Empty State -->
    <div v-if="!isLoading && classes.length === 0 && !error" class="text-center q-mt-xl">
      <q-icon name="sentiment_dissatisfied" size="4em" color="grey" />
      <div class="text-h6 text-grey">Belum ada kelas yang diajar</div>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useAttendance } from '../../composables/useAttendance'
import { useContextStore } from '../../stores/context.store.js'

// 🔧 Ambil semua yang dibutuhkan dari useAttendance
const { classes, isLoading, error, loadClasses } = useAttendance()
const contextStore = useContextStore()
const teacherName = contextStore.userId || 'Budi Santoso' // Hardcode sementara

// 🔧 Panggil loadClasses saat komponen dimount
onMounted(() => {
  loadClasses()
})
</script>

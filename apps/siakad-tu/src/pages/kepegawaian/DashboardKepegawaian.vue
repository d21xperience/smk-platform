<template>
  <q-page class="q-pa-md bg-grey-2">
    <div class="row justify-between items-center q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Dashboard Kepegawaian</h1>
        <p class="text-caption text-grey-7 q-mb-none">Ringkasan data SDM Sekolah</p>
      </div>
      <q-btn color="primary" icon="refresh" label="Refresh Data" @click="fetchStats" :loading="isLoading" />
    </div>

    <div class="row q-col-gutter-md">
      <!-- Card 1: Total Guru -->
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="bg-primary text-white shadow-2 cursor-pointer" @click="$router.push({ name: 'direktori-gtk' })">
          <q-card-section>
            <div class="row items-center justify-between">
              <div>
                <div class="text-subtitle2 text-grey-4">Total Guru (Aktif)</div>
                <div class="text-h3 text-weight-bold">{{ stats.totalTeachers || 0 }}</div>
              </div>
              <q-icon name="groups" size="48px" class="text-white opacity-50" />
            </div>
          </q-card-section>
        </q-card>
      </div>

      <!-- Card 2: Rata-rata SKP -->
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="bg-green text-white shadow-2 cursor-pointer" @click="$router.push({ name: 'kinerja-gtk' })">
          <q-card-section>
            <div class="row items-center justify-between">
              <div>
                <div class="text-subtitle2 text-grey-4">Rata-rata Nilai SKP</div>
                <div class="text-h3 text-weight-bold">{{ (stats.averageSkp || 0).toFixed(2) }}</div>
              </div>
              <q-icon name="trending_up" size="48px" class="text-white opacity-50" />
            </div>
          </q-card-section>
        </q-card>
      </div>

      <!-- Card 3: Cuti Pending -->
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="bg-orange text-white shadow-2 cursor-pointer" @click="$router.push({ name: 'cuti-kerja' })">
          <q-card-section>
            <div class="row items-center justify-between">
              <div>
                <div class="text-subtitle2 text-grey-4">Cuti Menunggu Persetujuan</div>
                <div class="text-h3 text-weight-bold">{{ stats.pendingLeaves || 0 }}</div>
              </div>
              <q-icon name="pending_actions" size="48px" class="text-white opacity-50" />
            </div>
          </q-card-section>
        </q-card>
      </div>

      <!-- Card 4: Arsip Kadaluarsa -->
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="bg-red text-white shadow-2 cursor-pointer" @click="$router.push({ name: 'arsip-kepegawaian' })">
          <q-card-section>
            <div class="row items-center justify-between">
              <div>
                <div class="text-subtitle2 text-grey-4">Dokumen Mendekati Kadaluarsa</div>
                <div class="text-h3 text-weight-bold">{{ stats.expiringArchives || 0 }}</div>
              </div>
              <q-icon name="warning_amber" size="48px" class="text-white opacity-50" />
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Section Bawah (Optional Info) -->
    <div class="row q-mt-lg">
      <div class="col-12">
        <q-card class="bg-white shadow-1">
          <q-card-section class="row items-center">
            <q-icon name="info" color="primary" size="sm" class="q-mr-sm" />
            <div class="text-subtitle2 text-grey-8">Klik pada kartu di atas untuk langsung menuju ke modul terkait.
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { api } from '@/boot/axios'

const $q = useQuasar()
const isLoading = ref(false)
const stats = ref({
  totalTeachers: 0,
  pendingLeaves: 0,
  averageSkp: 0,
  expiringArchives: 0
})

const fetchStats = async () => {
  isLoading.value = true
  try {
    const res = await api.get('/teachers/dashboard')
    stats.value = res.data
  } catch (err) {
    $q.notify({ type: 'negative', message: `Gagal memuat statistik dashboard ${err}` })
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  fetchStats()
})
</script>

<style scoped>
.opacity-50 {
  opacity: 0.5;
}
</style>

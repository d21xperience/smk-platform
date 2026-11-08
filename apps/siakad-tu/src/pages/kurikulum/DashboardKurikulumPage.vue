<template>
  <q-page class="dashboard-page bg-grey-1">
    <!-- HEADER -->
    <div class="dashboard-header q-px-md q-py-lg">
      <div class="container max-width-center row items-center justify-between">
        <div>
          <div class="text-overline text-secondary text-weight-medium letter-spacing-2">
            SISTEM INFORMASI AKADEMIK
          </div>
          <div class="text-h5 text-weight-bold text-white q-mt-xs">
            Dashboard Kurikulum
          </div>
          <div class="text-caption text-grey-4 q-mt-xs">
            {{ todayLabel }} • Tahun Ajaran Aktif
          </div>
        </div>
      </div>
    </div>

    <div class="container max-width-center q-px-md q-py-lg">
      <!-- KPI CARDS (DATA DINAMIS) -->
      <div class="row q-col-gutter-md q-mb-lg">
        <div class="col-6 col-md-3">
          <q-card flat bordered class="kpi-card">
            <q-card-section class="row items-center no-wrap">
              <q-avatar color="primary" text-color="white" icon="menu_book" size="46px" class="q-mr-md" />
              <div>
                <div class="text-h6 text-weight-bold text-primary">{{ stats.totalMapel }}</div>
                <div class="text-caption text-grey-7">Total Mata Pelajaran</div>
              </div>
            </q-card-section>
          </q-card>
        </div>

        <div class="col-6 col-md-3">
          <q-card flat bordered class="kpi-card">
            <q-card-section class="row items-center no-wrap">
              <q-avatar color="secondary" text-color="white" icon="badge" size="46px" class="q-mr-md" />
              <div>
                <div class="text-h6 text-weight-bold text-primary">{{ stats.guruMengajarAktif }}</div>
                <div class="text-caption text-grey-7">Guru Mengajar Aktif</div>
              </div>
            </q-card-section>
          </q-card>
        </div>

        <div class="col-6 col-md-3">
          <q-card flat bordered class="kpi-card">
            <q-card-section class="row items-center no-wrap">
              <q-avatar color="teal" text-color="white" icon="trending_up" size="46px" class="q-mr-md" />
              <div>
                <div class="text-h6 text-weight-bold text-primary">68%</div>
                <div class="text-caption text-grey-7">Progress Kurikulum</div>
              </div>
            </q-card-section>
          </q-card>
        </div>

        <div class="col-6 col-md-3">
          <q-card flat bordered class="kpi-card">
            <q-card-section class="row items-center no-wrap">
              <q-avatar color="orange-8" text-color="white" icon="pending_actions" size="46px" class="q-mr-md" />
              <div>
                <div class="text-h6 text-weight-bold text-primary">0</div>
                <div class="text-caption text-grey-7">RPP Belum Dikumpulkan</div>
              </div>
            </q-card-section>
          </q-card>
        </div>
      </div>

      <div class="row q-col-gutter-lg">
        <div class="col-12 col-md-8">
          <!-- Progress Materi (Masih Mock) -->
          <q-card flat bordered class="q-mb-lg content-card">
            <q-card-section>
              <div class="row items-center justify-between">
                <div class="text-subtitle1 text-weight-bold text-primary">
                  Progress Penyelesaian Materi per Kelompok Mapel
                </div>
              </div>
            </q-card-section>
            <q-separator />
            <q-card-section>
              <div class="text-caption text-grey-6 text-center q-pa-md">
                Fitur tracking progress per mapel belum tersedia.
              </div>
            </q-card-section>
          </q-card>
        </div>

        <div class="col-12 col-md-4">
          <!-- BEBAN MENGAJAR GURU (DATA DINAMIS) -->
          <q-card flat bordered class="q-mb-lg content-card">
            <q-card-section>
              <div class="text-subtitle1 text-weight-bold text-primary">
                Beban Mengajar Guru (Tugas)
              </div>
            </q-card-section>
            <q-separator />
            <q-list separator>
              <q-item v-for="guru in stats.bebanMengajar" :key="guru.teacherName">
                <q-item-section>
                  <q-item-label class="text-weight-medium">{{ guru.teacherName }}</q-item-label>
                  <q-item-label caption>{{ guru.subjectName }}</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-badge :color="guru.jp > 20 ? 'orange' : 'primary'" text-color="white">
                    {{ guru.jp }} Kelas
                  </q-badge>
                </q-item-section>
              </q-item>
              <q-item v-if="stats.bebanMengajar.length === 0">
                <q-item-section class="text-center text-grey-6">Belum ada data penugasan</q-item-section>
              </q-item>
            </q-list>
          </q-card>
        </div>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { api } from '@/boot/axios'

const todayLabel = computed(() => {
  const today = new Date()
  return today.toLocaleDateString('id-ID', {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  })
})

const stats = ref({
  totalMapel: 0,
  guruMengajarAktif: 0,
  bebanMengajar: []
})

onMounted(async () => {
  try {
    const res = await api.get('/kurikulum/statistics')
    stats.value = res.data
  } catch (err) {
    console.error('Gagal load statistik kurikulum:', err)
  }
})
</script>

<style scoped>
.max-width-center { max-width: 1200px; margin: 0 auto; }
.letter-spacing-2 { letter-spacing: 0.12em; }
.dashboard-header { background: linear-gradient(135deg, #0a192f 0%, #123a63 100%); }
.kpi-card { border-radius: 8px; height: 100%; }
.content-card { border-radius: 8px; }
</style>

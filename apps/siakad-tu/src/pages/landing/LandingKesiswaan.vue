<template>
  <q-page class="dashboard-page bg-grey-1">
    <!-- HEADER (Reusable) -->
    <DashboardHeader
      :title="'Dashboard Kesiswaan'"
      :subtitle="'SISTEM INFORMASI AKADEMIK'"
      :dateLabel="todayLabel"
      :notificationCount="notifCount"
      actionLabel="Input Pelanggaran"
      actionIcon="add"
      @click:action="handleInputPelanggaran"
      @click:notification="handleNotification"
    />

    <div class="container max-width-center q-px-md q-py-lg">
      <!-- Loading State -->
      <div v-if="loading" class="row justify-center q-py-xl">
        <q-spinner color="primary" size="3em" />
        <div class="text-caption text-grey-6 q-mt-sm">Memuat data dashboard...</div>
      </div>

      <!-- Content (only show when data is available and not loading) -->
      <template v-else-if="data">
        <!-- KPI CARDS (Menggunakan komponen KpiCard yang sudah ada) -->
        <div class="row q-col-gutter-md q-mb-lg">
          <div
            class="col-6 col-md-3"
            v-for="kpi in kpiList"
            :key="kpi.key"
          >
            <KpiCard
              :label="kpi.label"
              :value="kpi.value"
              :icon="kpi.icon"
              :color="kpi.color"
              :route-to="kpi.routeTo"
              @click="handleKpiClick(kpi.key)"
            />
          </div>
        </div>

        <div class="row q-col-gutter-lg">
          <!-- LEFT COLUMN -->
          <div class="col-12 col-md-8">
            <!-- SISWA PERLU PERHATIAN KHUSUS -->
            <AttentionList
              v-if="data.siswaPerhatian?.length"
              :items="data.siswaPerhatian"
              title="Siswa Memerlukan Perhatian Khusus"
              icon="report"
              color="red"
              @click:item="handleAttentionItemClick"
            />

            <!-- TREN PELANGGARAN -->
            <TrendChart
              :items="data.trenPelanggaran || []"
              title="Tren Pelanggaran Siswa (6 Bulan Terakhir)"
              :badgeLabel="`${data.totalPelanggaranBulanIni || 0} kasus bulan ini`"
            />

            <!-- DISTRIBUSI POIN -->
            <DistributionBar
              :items="data.poinPerTingkat || []"
              title="Distribusi Poin Pelanggaran per Tingkat Kelas"
            />

            <!-- KASUS TERBARU -->
            <RecentCases
              :items="data.kasusTerbaru || []"
              title="Catatan Kasus Terbaru"
              @click:item="handleKasusItemClick"
            />
          </div>

          <!-- RIGHT COLUMN -->
          <div class="col-12 col-md-4">
            <!-- EKSTRAKURIKULER -->
            <ExtracurricularList
              :items="data.ekstrakurikuler || []"
              title="Peserta Ekstrakurikuler"
              @click:item="handleEkskulItemClick"
            />

            <!-- PRESTASI -->
            <AchievementList
              :items="data.prestasiList || []"
              title="Prestasi Terbaru"
              @click:item="handlePrestasiItemClick"
            />

            <!-- AGENDA -->
            <AgendaList
              :items="data.agendaList || []"
              title="Agenda Kesiswaan"
              @click:item="handleAgendaItemClick"
            />
          </div>
        </div>
      </template>

      <!-- Error State -->
      <div v-else-if="error" class="row justify-center q-py-xl">
        <div class="text-center">
          <q-icon name="error" color="negative" size="3em" />
          <div class="text-h6 text-negative q-mt-sm">Gagal memuat data</div>
          <div class="text-caption text-grey-7">{{ error }}</div>
          <q-btn flat color="primary" label="Coba Lagi" @click="load" class="q-mt-md" />
        </div>
      </div>

      <!-- Empty State (jika data kosong) -->
      <div v-else class="row justify-center q-py-xl">
        <div class="text-center">
          <q-icon name="dashboard" color="grey-5" size="3em" />
          <div class="text-h6 text-grey-7 q-mt-sm">Belum ada data</div>
          <div class="text-caption text-grey-6">Silakan muat ulang halaman</div>
        </div>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useKesiswaanDashboard } from '@/composables/useKesiswaanDashboard'
import DashboardHeader from '@/components/dashboard/DashboardHeader.vue'
import AttentionList from '@/components/dashboard/AttentionList.vue'
import TrendChart from '@/components/dashboard/TrendChart.vue'
import DistributionBar from '@/components/dashboard/DistributionBar.vue'
import RecentCases from '@/components/dashboard/RecentCases.vue'
import ExtracurricularList from '@/components/dashboard/ExtracurricularList.vue'
import AchievementList from '@/components/dashboard/AchievementList.vue'
import AgendaList from '@/components/dashboard/AgendaList.vue'
import KpiCard from '@/components/KpiCard.vue'

// === COMPOSABLE ===
const {
  data,
  loading,
  error,
  load,
  handleKpiClick,
  // helpers (trendHeight, barWidth, kategoriWarna) tidak digunakan langsung di template,
  // karena komponen sudah menangani logika perhitungan sendiri.
} = useKesiswaanDashboard()

// === UI-SPECIFIC STATE (PAGE-ONLY) ===
const todayLabel = computed(() => {
  const today = new Date()
  return today.toLocaleDateString('id-ID', {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  })
})

const notifCount = ref(6)

// === KPI LIST (REACTIVE) ===
const kpiList = computed(() => {
  if (!data.value) return []
  return [
    {
      key: 'totalSiswaAktif',
      label: 'Siswa Aktif',
      value: data.value.totalSiswaAktif || 0,
      icon: 'groups',
      color: 'primary',
      routeTo: { name: 'student-list' },
    },
    {
      key: 'totalPelanggaranBulanIni',
      label: 'Kasus Pelanggaran Bulan Ini',
      value: data.value.totalPelanggaranBulanIni || 0,
      icon: 'gavel',
      color: 'red-7',
      routeTo: { name: 'pelanggaran-list' },
    },
    {
      key: 'totalBerprestasi',
      label: 'Siswa Berprestasi',
      value: data.value.totalBerprestasi || 0,
      icon: 'emoji_events',
      color: 'amber-8',
      routeTo: { name: 'prestasi-list' },
    },
    {
      key: 'totalEkskul',
      label: 'Peserta Ekstrakurikuler',
      value: data.value.totalEkskul || 0,
      icon: 'sports',
      color: 'teal',
      routeTo: { name: 'ekskul-list' },
    },
  ]
})

// === EVENT HANDLERS (meneruskan ke composable atau langsung) ===
function handleInputPelanggaran() {
  // Bisa arahkan ke route atau dialog
  // Misal: router.push({ name: 'input-pelanggaran' })
  // Untuk saat ini, kita serahkan ke handleKpiClick atau requestAccess
  handleKpiClick('totalPelanggaranBulanIni')
}

function handleNotification() {
  // Bisa buka panel notifikasi
  console.log('Notifikasi diklik')
}

function handleAttentionItemClick(item) {
  // Navigasi ke detail siswa
  console.log('Klik siswa perhatian:', item)
  // router.push({ name: 'student-detail', params: { id: item.id } })
}

function handleKasusItemClick(item) {
  // Navigasi ke detail kasus
  console.log('Klik kasus:', item)
  // router.push({ name: 'kasus-detail', params: { id: item.id } })
}

function handleEkskulItemClick(item) {
  // Navigasi ke detail ekstrakurikuler
  console.log('Klik ekstrakurikuler:', item)
  // router.push({ name: 'ekskul-detail', params: { id: item.id } })
}

function handlePrestasiItemClick(item) {
  // Navigasi ke detail prestasi
  console.log('Klik prestasi:', item)
  
  // router.push({ name: 'prestasi-detail', params: { id: item.id } })
}

function handleAgendaItemClick(item) {
  // Navigasi ke detail agenda
  console.log('Klik agenda:', item)
  // router.push({ name: 'agenda-detail', params: { id: item.id } })
}

// === LOAD DATA ON MOUNT ===
onMounted(() => {
  load()
})
</script>

<style scoped>
.max-width-center {
  max-width: 1200px;
  margin: 0 auto;
}

.dashboard-page {
  min-height: 100vh;
}
</style>

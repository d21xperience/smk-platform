<template>
  <q-page class="dashboard-page bg-grey-1">
    <div class="container q-px-md q-py-sm">
      <!-- <div class="row q-col-gutter-lg"> -->
      <ApexBarChart title="Jumlah Siswa Masuk & Lulus per Tahun Ajaran" :labels="tahunList" :lulusData="dataLulus"
        :masukData="dataMasuk" />

      <ApexClassChart title="Data Siswa Per Kelas" :labels="namaKelas" :lakiData="dataLakiLaki"
        :perempuanData="dataPerempuan" />

      <ApexMonthlyCarouselChart :rawMonthlyData="dataSiswaBulanan" title="Jumlah Siswa Per Bulan" />


    </div>

  </q-page>
</template>

<script setup>
import ApexBarChart from '@/components/ApexBarChart.vue'
import { ref } from 'vue'
import ApexClassChart from '@/components/ApexClassChart.vue'

// Data ekstraksi murni dari gambar tabel Anda
const namaKelas = ref([
  'X TKRO', 'X TBSM A', 'X TBSM B', 'X TBSM C',
  'X TKJT A', 'X TKJT B', 'X TKJT C', 'X TKJT D', 'X AP'
])
const dataLakiLaki = ref([23, 35, 32, 4, 18, 13, 24, 21, 21])
const dataPerempuan = ref([0, 3, 6, 16, 23, 6, 11, 0, 12])
// const todayLabel = computed(() => {
//   const today = new Date()
//   return today.toLocaleDateString('id-ID', {
//     weekday: 'long',
//     day: 'numeric',
//     month: 'long',
//     year: 'numeric',
//   })
// })

// const notifCount = ref(6)
// Data dummy yang persis dengan contoh gambar Anda
const tahunList = ref(['2021/2022', '2022/2023', '2023/2024', '2024/2025', '2025/2026', '2026/2027'])
const dataLulus = ref([224, 200, 205, 205, 200, 250])
const dataMasuk = ref([205, 205, 350, 200, 250, 350])
// const kpis = ref([
//   { label: 'Siswa Aktif', value: '853', icon: 'groups', color: 'primary', maleCount: 200, femaleCount: 150 },
//   { label: 'Guru', value: '58', icon: 'groups', color: 'red-7', maleCount: 25, femaleCount: 27 },
//   { label: 'Tata Usaha', value: '11', icon: 'groups', color: 'amber-8' },
//   { label: 'Rombel', value: '27', icon: 'build', color: 'teal' },
// ])


import ApexMonthlyCarouselChart from '@/components/ApexMonthlyCarouselChart.vue'

// Contoh ekstraksi data murni dari tabel semester 1 Anda
const dataSiswaBulanan = ref([
  // --- TINGKAT X ---
  { tingkat: 'X', nama: 'X TKRO', juli: 23, agustus: 32, september: 24, oktober: 28, november: 22, desember: 23 },
  { tingkat: 'X', nama: 'X TBSM A', juli: 38, agustus: 38, september: 39, oktober: 31, november: 35, desember: 26 },
  { tingkat: 'X', nama: 'X TBSM B', juli: 38, agustus: 32, september: 36, oktober: 29, november: 28, desember: 25 },
  { tingkat: 'X', nama: 'X TBSM C', juli: 20, agustus: 31, september: 31, oktober: 39, november: 20, desember: 39 },
  { tingkat: 'X', nama: 'X TKJT A', juli: 41, agustus: 25, september: 30, oktober: 21, november: 40, desember: 20 },
  // ... Tambahkan baris kelas X lainnya dari tabel ...

  // --- TINGKAT XI ---
  { tingkat: 'XI', nama: 'XI TKRO', juli: 36, agustus: 33, september: 25, oktober: 20, november: 40, desember: 24 },
  { tingkat: 'XI', nama: 'XI TBSM A', juli: 28, agustus: 25, september: 37, oktober: 34, november: 24, desember: 22 },
  // ... Tambahkan baris kelas XI lainnya dari tabel ...

  // --- TINGKAT XII ---
  { tingkat: 'XII', nama: 'XII TKRO', juli: 37, agustus: 21, september: 23, oktober: 33, november: 36, desember: 33 },
  { tingkat: 'XII', nama: 'XII TBSM A', juli: 36, agustus: 31, september: 19, oktober: 22, november: 25, desember: 30 }
  // ... Tambahkan baris kelas XII lainnya dari tabel ...
])
</script>

<style scoped>
.max-width-center {
  max-width: 1200px;
  margin: 0 auto;
}

.letter-spacing-2 {
  letter-spacing: 0.12em;
}

.dashboard-header {
  background: linear-gradient(135deg, #0a192f 0%, #123a63 100%);
}

.kpi-card {
  border-radius: 8px;
  height: 100%;
}

.content-card {
  border-radius: 8px;
}

.border-alert {
  border-left: 4px solid #c62828;
}

.bar-track {
  width: 100%;
  height: 10px;
  border-radius: 6px;
  background-color: #eceff1;
  overflow: hidden;
}

.bar-fill {
  height: 100%;
  border-radius: 6px;
  transition: width 0.4s ease;
}

.trend-chart {
  height: 170px;
  gap: 12px;
}

.trend-bar-wrap {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: flex-end;
}

.trend-bar {
  width: 100%;
  max-width: 32px;
  background: linear-gradient(180deg, #123a63 0%, #0a192f 100%);
  border-radius: 4px 4px 0 0;
  transition: height 0.4s ease;
}

.agenda-date {
  min-width: 44px;
}
</style>

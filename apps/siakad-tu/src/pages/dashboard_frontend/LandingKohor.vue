<template>
  <q-page class="" style="height:100vh;border: 15px;overflow: hidden;">
    <div>
      <ApexBarChart title="Jumlah Siswa Masuk & Lulus per Tahun Ajaran" :labels="tahunList" :lulusData="dataLulus"
        :masukData="dataMasuk" />
    </div>
  </q-page>
</template>

<script setup>
import { ref } from 'vue'
import ApexBarChart from '@/components/ApexBarChart.vue'
const tahunList = ref(['2021/2022', '2022/2023', '2023/2024', '2024/2025', '2025/2026', '2026/2027'])
const dataLulus = ref([224, 200, 205, 205, 200, 250])
const dataMasuk = ref([205, 205, 350, 200, 250, 350])
</script>

<style scoped>
/* .max-width-center {
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
} */
</style>

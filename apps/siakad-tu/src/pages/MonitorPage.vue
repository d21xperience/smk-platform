<template>
  <q-page padding>
    <div class="text-h4 q-mb-md">Monitoring Real-time</div>
    <div class="row q-col-gutter-md q-mb-md">
      <div class="col-12 col-md-2" v-for="stat in statusStats" :key="stat.label">
        <q-card :class="`bg-${stat.color} text-white`">
          <q-card-section class="text-center">
            <div class="text-h5">{{ stat.count }}</div>
            <div>{{ stat.label }}</div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <q-table :rows="store.mockStudents" :columns="columns" row-key="id" :rows-per-page-options="[0]" flat bordered>
      <template v-slot:body-cell-status="props">
        <q-td :props="props"><q-chip :color="statusColor(props.value)" text-color="white" size="sm">{{ props.value
        }}</q-chip></q-td>
      </template>
      <template v-slot:body-cell-progress="props">
        <q-td :props="props">
          <q-linear-progress :value="props.value / 100" size="8px" color="primary" class="q-mb-xs" />{{ props.value }}%
        </q-td>
      </template>
      <template v-slot:body-cell-lastHeartbeat="props">
        <q-td :props="props"><span :class="{ 'text-negative': isStale(props.value) }">{{ formatTimeAgo(props.value)
        }}</span></q-td>
      </template>
    </q-table>
  </q-page>
</template>

<script setup>
import { computed } from 'vue'
import { useExamStore } from 'stores/exam-store'

const store = useExamStore()
const columns = [
  { name: 'id', label: 'ID', field: 'id', align: 'left' },
  { name: 'name', label: 'Nama Siswa', field: 'name', align: 'left' },
  { name: 'status', label: 'Status', field: 'status', align: 'center' },
  { name: 'progress', label: 'Progress', field: 'progress', align: 'center' },
  { name: 'lastHeartbeat', label: 'Heartbeat', field: 'lastHeartbeat', align: 'center' }
]

const statusStats = computed(() => {
  const list = ['Belum Login', 'Login', 'Mengerjakan', 'Selesai', 'Terputus']
  const colors = ['grey', 'blue', 'orange', 'green', 'red']
  return list.map((s, i) => ({ label: s, count: store.mockStudents.filter(st => st.status === s).length, color: colors[i] }))
})

const statusColor = (s) => ({ 'Belum Login': 'grey', 'Login': 'blue', 'Mengerjakan': 'orange', 'Selesai': 'green', 'Terputus': 'red' }[s] || 'grey')
const isStale = (ts) => ts ? Date.now() - ts > 30000 : true
const formatTimeAgo = (ts) => {
  if (!ts) return '-'
  const d = Math.floor((Date.now() - ts) / 1000)
  return d < 10 ? 'Baru saja' : `${d}d lalu`
}
</script>

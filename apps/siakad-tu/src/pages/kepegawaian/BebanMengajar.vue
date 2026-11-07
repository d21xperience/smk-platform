<template>
  <q-page class="q-pa-md">
    <div class="row justify-between items-center q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Beban Mengajar</h1>
        <p class="text-caption text-grey-7 q-mb-none">Distribusi Jam Tatap Muka (JTM) Guru per Tahun Ajaran</p>
      </div>
      <q-btn color="primary" icon="refresh" label="Refresh" @click="fetchWorkloads" :loading="isLoading" />
    </div>

    <q-card class="shadow-1">
      <q-table
        :rows="workloads"
        :columns="columns"
        row-key="teacherId"
        flat
        bordered
        dense
        :loading="isLoading"
        :pagination="{ rowsPerPage: 0 }"
      >
        <template v-slot:header="props">
          <q-tr :props="props">
            <q-th auto-width />
            <q-th v-for="col in props.cols" :key="col.name" :props="props">
              {{ col.label }}
            </q-th>
          </q-tr>
        </template>

        <template v-slot:body="props">
          <q-tr :props="props">
            <q-td auto-width>
              <q-btn size="sm" color="secondary" round dense @click="props.expand = !props.expand" :icon="props.expand ? 'remove' : 'add'" />
            </q-td>
            <q-td key="teacherName" :props="props" class="text-weight-medium">{{ props.row.teacherName }}</q-td>
            <q-td key="totalClasses" :props="props">
              <q-badge color="blue-6" outline class="q-pa-sm">{{ props.row.totalClasses }} Kelas</q-badge>
            </q-td>
            <q-td key="totalSubjects" :props="props">
              <q-badge color="green-6" outline class="q-pa-sm">{{ props.row.totalSubjects }} Mapel</q-badge>
            </q-td>
            <q-td key="totalSessions" :props="props">
              <q-badge color="orange-6" outline class="q-pa-sm">{{ props.row.totalSessions }} Sesi</q-badge>
            </q-td>
          </q-tr>
          
          <q-tr v-show="props.expand" :props="props">
            <q-td colspan="100%" class="bg-grey-1 q-pa-md">
              <div v-if="!props.row.schedules || props.row.schedules.length === 0" class="text-grey-6 text-center q-pa-sm">
                <q-icon name="info" /> Belum ada jadwal yang diatur untuk guru ini pada tahun ajaran aktif.
              </div>
              <div v-else>
                <div class="text-subtitle2 text-primary q-mb-sm">Rincian Jadwal:</div>
                <q-table
                  :rows="props.row.schedules"
                  :columns="scheduleColumns"
                  dense
                  flat
                  bordered
                  hide-pagination
                  row-key="id"
                  class="bg-white"
                >
                  <template v-slot:body-cell-day="prop">
                    <q-td :props="prop">
                      <q-badge color="primary">{{ prop.row.day }}</q-badge>
                    </q-td>
                  </template>
                  <template v-slot:body-cell-time="prop">
                    <q-td :props="prop">
                      <q-icon name="schedule" size="xs" class="q-mr-xs"/> {{ prop.row.startTime }} - {{ prop.row.endTime }}
                    </q-td>
                  </template>
                </q-table>
              </div>
            </q-td>
          </q-tr>
        </template>
      </q-table>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Notify } from 'quasar'
import { api } from '@/boot/axios'
import { useContextStore } from '@/stores/contextStore'

const contextStore = useContextStore()
const workloads = ref([])
const isLoading = ref(false)

const columns = [
  { name: 'teacherName', label: 'Nama Guru', field: 'teacherName', align: 'left', sortable: true },
  { name: 'totalClasses', label: 'Total Kelas', field: 'totalClasses', align: 'center', sortable: true },
  { name: 'totalSubjects', label: 'Total Mata Pelajaran', field: 'totalSubjects', align: 'center', sortable: true },
  { name: 'totalSessions', label: 'Total Sesi (JTM)', field: 'totalSessions', align: 'center', sortable: true }
]

const scheduleColumns = [
  { name: 'day', label: 'Hari', field: 'day', align: 'left', sortable: true },
  { name: 'time', label: 'Jam', align: 'left' },
  { name: 'className', label: 'Kelas', field: 'className', align: 'left' },
  { name: 'subjectName', label: 'Mata Pelajaran', field: 'subjectName', align: 'left' }
]

const fetchWorkloads = async () => {
  isLoading.value = true
  try {
    const res = await api.get('/teachers/workload')
    workloads.value = res.data || []
  } catch (err) {
    console.error(err)
    Notify.create({ type: 'negative', message: 'Gagal memuat data beban mengajar.' })
  } finally {
    isLoading.value = false
  }
}

onMounted(async () => {
  if (!contextStore.hasContext) {
    await contextStore.initialize()
  }
  fetchWorkloads()
})
</script>

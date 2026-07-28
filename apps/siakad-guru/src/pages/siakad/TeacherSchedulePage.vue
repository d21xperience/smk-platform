<template>
  <q-page class="q-pa-md">
    <div class="text-h6 q-mb-md">Jadwal Mengajar</div>

    <!-- Info Konteks -->
    <div v-if="!context.isOperationalContextReady.value" class="text-grey q-mb-md">
      Silakan pilih konteks pembelajaran terlebih dahulu.
      <q-btn label="Pilih Konteks" to="/siakad/contex" flat color="primary" />
    </div>

    <div v-else>
      <div class="text-subtitle1 text-grey q-mb-sm">
        {{ context.operational.value.academicYearName }} | {{ context.operational.value.semesterName }}
      </div>

      <!-- Loading -->
      <div v-if="teaching.loading.value.schedule" class="text-center q-pa-lg">
        <q-spinner color="primary" size="3em" />
      </div>

      <!-- Tabel Jadwal -->
      <q-table v-else :rows="groupedSchedule" :columns="columns" row-key="day" flat bordered
        :rows-per-page-options="[0]" hide-bottom>
        <template v-slot:header="props">
          <q-tr :props="props">
            <q-th auto-width>Hari</q-th>
            <q-th>Mata Pelajaran</q-th>
            <q-th>Kelas</q-th>
            <q-th>Jam</q-th>
            <q-th>Ruangan</q-th>
          </q-tr>
        </template>

        <template v-slot:body="props">
          <q-tr :props="props">
            <q-td auto-width>
              <strong>{{ props.row.day }}</strong>
            </q-td>
            <q-td>
              <div v-for="item in props.row.items" :key="item.id" class="q-mb-sm">
                {{ item.subject }}
              </div>
            </q-td>
            <q-td>
              <div v-for="item in props.row.items" :key="item.id" class="q-mb-sm">
                {{ item.className }}
              </div>
            </q-td>
            <q-td>
              <div v-for="item in props.row.items" :key="item.id" class="q-mb-sm">
                {{ item.startTime }} - {{ item.endTime }}
              </div>
            </q-td>
            <q-td>
              <div v-for="item in props.row.items" :key="item.id" class="q-mb-sm">
                {{ item.room }}
              </div>
            </q-td>
          </q-tr>
        </template>
      </q-table>
    </div>
  </q-page>
</template>

<script setup>
import { onMounted, computed } from 'vue'
import { useContext } from '@/composables/useContext'
import { useTeaching } from '@/composables/useTeaching'

const context = useContext()
const teaching = useTeaching()

const columns = [
  { name: 'day', label: 'Hari', align: 'left', field: 'day' },
  { name: 'subject', label: 'Mata Pelajaran', align: 'left', field: 'subject' },
  { name: 'class', label: 'Kelas', align: 'left', field: 'class' },
  { name: 'time', label: 'Jam', align: 'left', field: 'time' },
  { name: 'room', label: 'Ruangan', align: 'left', field: 'room' }
]

const groupedSchedule = computed(() => {
  const days = ['Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat', 'Sabtu']
  return days.map(day => ({
    day,
    items: teaching.schedule.value.filter(item => item.dayOfWeek === day)
  })).filter(day => day.items.length > 0)
})

onMounted(async () => {
  if (context.isOperationalContextReady.value) {
    await teaching.loadSchedule()
  }
})
</script>

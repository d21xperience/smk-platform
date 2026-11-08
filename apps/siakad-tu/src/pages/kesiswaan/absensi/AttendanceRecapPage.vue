<template>
  <q-page class="q-pa-md">
    <div class="text-h5 text-weight-bold q-mb-md">Rekap Kehadiran Siswa</div>

    <q-card flat bordered>
      <q-card-section>
        <div class="row q-col-gutter-md">
          <div class="col-4">
            <q-select v-model="classID" :options="classOptions" label="Kelas" outlined dense emit-value map-options />
          </div>
          <div class="col-3">
            <q-select v-model="month" :options="monthOptions" label="Bulan" outlined dense emit-value map-options />
          </div>
          <div class="col-2">
            <q-input v-model.number="year" label="Tahun" type="number" outlined dense />
          </div>
          <div class="col-3">
            <q-btn color="primary" icon="search" label="Tampilkan" :loading="isLoading"
              class="full-width" style="height: 40px" @click="fetchRecap" />
          </div>
        </div>
      </q-card-section>
    </q-card>

    <q-table
      :rows="recap"
      :columns="columns"
      row-key="StudentID"
      flat bordered
      class="q-mt-md"
      :loading="isLoading"
      no-data-label="Belum ada data absensi pada periode ini"
    >
      <template v-slot:body-cell-percentage="props">
        <q-td :props="props">
          <q-badge :color="badgeColor(props.row.Summary?.AttendancePercentage)">
            {{ formatPct(props.row.Summary?.AttendancePercentage) }}
          </q-badge>
        </q-td>
      </template>
    </q-table>
  </q-page>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { Notify } from 'quasar'
import { api } from '@/boot/axios'
import { useContextStore } from '@/stores/contextStore'

const contextStore = useContextStore()

const classID = ref('')
const month = ref(new Date().getMonth() + 1)
const year = ref(new Date().getFullYear())
const recap = ref([])
const classes = ref([])
const isLoading = ref(false)

const monthOptions = [
  { label: 'Januari', value: 1 }, { label: 'Februari', value: 2 }, { label: 'Maret', value: 3 },
  { label: 'April', value: 4 }, { label: 'Mei', value: 5 }, { label: 'Juni', value: 6 },
  { label: 'Juli', value: 7 }, { label: 'Agustus', value: 8 }, { label: 'September', value: 9 },
  { label: 'Oktober', value: 10 }, { label: 'November', value: 11 }, { label: 'Desember', value: 12 },
]

const classOptions = computed(() => classes.value.map(c => ({
  label: c.Name ?? c.nm_kelas ?? c.Nama ?? c.ID,
  value: c.ID ?? c.rombongan_belajar_id
})))

const columns = [
  { name: 'name', label: 'Nama Siswa', field: r => r.StudentName || r.student_name, align: 'left', sortable: true },
  { name: 'total', label: 'Total Sesi', field: r => r.Summary?.TotalSessions ?? 0, align: 'center' },
  { name: 'present', label: 'Hadir', field: r => r.Summary?.PresentCount ?? 0, align: 'center' },
  { name: 'permission', label: 'Izin', field: r => r.Summary?.PermissionCount ?? 0, align: 'center' },
  { name: 'sick', label: 'Sakit', field: r => r.Summary?.SickCount ?? 0, align: 'center' },
  { name: 'absent', label: 'Alpa', field: r => r.Summary?.AbsentCount ?? 0, align: 'center' },
  { name: 'late', label: 'Terlambat', field: r => r.Summary?.LateCount ?? 0, align: 'center' },
  { name: 'percentage', label: '% Kehadiran', align: 'center' },
]

function formatPct(v) {
  if (typeof v !== 'number' || isNaN(v)) return '0%'
  return v.toFixed(1) + '%'
}

function badgeColor(v) {
  if (typeof v !== 'number') return 'grey'
  if (v >= 90) return 'green'
  if (v >= 75) return 'orange'
  return 'red'
}

async function fetchClasses() {
  try {
    const res = await api.get('/schedules/options')
    classes.value = res.data.classes || []
  } catch (err) {
    Notify.create({ type: 'negative', message: 'Gagal memuat kelas: ' + (err.response?.data?.error || err.message) })
  }
}

async function fetchRecap() {
  if (!classID.value) {
    return Notify.create({ type: 'warning', message: 'Pilih kelas terlebih dahulu' })
  }
  isLoading.value = true
  try {
    const res = await api.get('/attendance/recap', {
      params: { class_id: classID.value, month: month.value, year: year.value }
    })
    recap.value = res.data || []
  } catch (err) {
    Notify.create({ type: 'negative', message: 'Gagal memuat rekap: ' + (err.response?.data?.error || err.message) })
  } finally {
    isLoading.value = false
  }
}

onMounted(async () => {
  if (!contextStore.hasContext) await contextStore.initialize()
  await fetchClasses()
})
</script>

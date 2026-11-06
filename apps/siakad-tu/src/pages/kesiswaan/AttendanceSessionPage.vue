<template>
  <q-page class="q-pa-md">
    <div class="text-h5 text-weight-bold q-mb-md">Input Absensi Harian</div>

    <!-- Step 1: Pilih Jadwal & Tanggal -->
    <q-card flat bordered>
      <q-card-section>
        <div class="text-subtitle1 text-weight-bold q-mb-sm">1. Pilih Jadwal & Tanggal</div>
        <div class="row q-col-gutter-md">
          <div class="col-6">
            <q-select
              v-model="selectedSchedule"
              :options="scheduleOptions"
              label="Pilih Jadwal Pelajaran"
              outlined dense emit-value map-options
              option-label="label" option-value="value"
            />
          </div>
          <div class="col-3">
            <q-input v-model="sessionDate" label="Tanggal Sesi" type="date" outlined dense />
          </div>
          <div class="col-3">
            <q-btn color="primary" icon="groups" label="Muat Roster" :loading="isLoadingRoster"
              :disable="!selectedSchedule" class="full-width" style="height: 40px" @click="loadRoster" />
          </div>
        </div>
      </q-card-section>
    </q-card>

    <!-- Step 2: Input Absensi -->
    <q-card v-if="roster.length > 0" flat bordered class="q-mt-md">
      <q-card-section>
        <div class="text-subtitle1 text-weight-bold q-mb-sm">
          2. Input Absensi — {{ selectedSchedule?.class_name }} | {{ selectedSchedule?.subject_name }}
        </div>
        <q-table :rows="roster" :columns="columns" row-key="id" flat bordered dense>
          <template v-slot:body-cell-status="props">
            <q-td :props="props">
              <q-select
                v-model="props.row.status"
                :options="statusOptions"
                outlined dense emit-value map-options
                style="min-width: 130px"
              />
            </q-td>
          </template>
          <template v-slot:body-cell-note="props">
            <q-td :props="props">
              <q-input v-model="props.row.note" outlined dense placeholder="Opsional" />
            </q-td>
          </template>
        </q-table>

        <div class="row justify-between items-center q-mt-md">
          <div class="text-caption">
            Hadir: <b>{{ statusCount('PRESENT') }}</b> |
            Izin: <b>{{ statusCount('PERMISSION') }}</b> |
            Sakit: <b>{{ statusCount('SICK') }}</b> |
            Alpa: <b>{{ statusCount('ABSENT') }}</b> |
            Terlambat: <b>{{ statusCount('LATE') }}</b>
          </div>
          <q-btn color="primary" icon="save" label="Simpan Absensi"
            :loading="isSubmitting" @click="submitAttendance" />
        </div>
      </q-card-section>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { Notify } from 'quasar'
import { api } from '@/boot/axios'
import { useContextStore } from '@/stores/contextStore'

const contextStore = useContextStore()

const schedules = ref([])
const selectedSchedule = ref(null)
const sessionDate = ref(new Date().toISOString().slice(0, 10))
const roster = ref([])
const isLoadingRoster = ref(false)
const isSubmitting = ref(false)

const statusOptions = [
  { label: 'Hadir', value: 'PRESENT' },
  { label: 'Izin', value: 'PERMISSION' },
  { label: 'Sakit', value: 'SICK' },
  { label: 'Alpa', value: 'ABSENT' },
  { label: 'Terlambat', value: 'LATE' },
]

const scheduleOptions = computed(() => schedules.value.map(s => ({
  label: `${s.day} ${s.start_time?.substring(0,5)} — ${s.class_name} / ${s.subject_name} (${s.teacher_name})`,
  value: s
})))

const columns = [
  { name: 'name', label: 'Nama Siswa', field: r => r.name, align: 'left', sortable: true },
  { name: 'status', label: 'Status', align: 'left' },
  { name: 'note', label: 'Catatan', align: 'left' },
]

async function fetchSchedules() {
  try {
    const res = await api.get('/schedules')
    schedules.value = res.data || []
  } catch (err) {
    Notify.create({ type: 'negative', message: 'Gagal memuat jadwal: ' + (err.response?.data?.error || err.message) })
  }
}

async function loadRoster() {
  if (!selectedSchedule.value) return
  isLoadingRoster.value = true
  try {
    const classID = selectedSchedule.value.class_id
    const res = await api.get('/students', { params: { class_id: classID, limit: 100, status: 'ACTIVE' } })
    const list = Array.isArray(res.data) ? res.data : (res.data.data || [])
    roster.value = list.map(s => ({
      id: s.ID || s.id || s.PesertaDidikID || s.peserta_didik_id,
      name: s.Name || s.name || s.nama || s.Nama || '-',
      status: 'PRESENT',
      note: '',
    }))
    if (roster.value.length === 0) {
      Notify.create({ type: 'warning', message: 'Tidak ada siswa aktif di kelas ini.' })
    }
  } catch (err) {
    Notify.create({ type: 'negative', message: 'Gagal memuat siswa: ' + (err.response?.data?.error || err.message) })
  } finally {
    isLoadingRoster.value = false
  }
}

function statusCount(code) {
  return roster.value.filter(r => r.status === code).length
}

async function submitAttendance() {
  if (!selectedSchedule.value) {
    return Notify.create({ type: 'warning', message: 'Pilih jadwal terlebih dahulu' })
  }
  if (!sessionDate.value) {
    return Notify.create({ type: 'warning', message: 'Isi tanggal sesi' })
  }
  if (roster.value.length === 0) {
    return Notify.create({ type: 'warning', message: 'Roster kosong' })
  }

  isSubmitting.value = true
  try {
    // Step 1: buat session DRAFT
    const s = selectedSchedule.value
    const createPayload = {
      class_id: s.class_id,
      subject_id: s.pembelajaran_id, // pakai pembelajaran_id sebagai referensi sesi
      teacher_id: s.teacher_id,
      session_date: sessionDate.value,
      start_time: (s.start_time || '').substring(0, 5) || '07:00',
      end_time: (s.end_time || '').substring(0, 5) || '08:30',
      notes: '',
    }
    const createRes = await api.post('/attendance/sessions', createPayload)
    const sessionID = createRes.data?.ID || createRes.data?.id
    if (!sessionID) throw new Error('Session ID tidak diterima dari server')

    // Step 2: submit records
    const records = roster.value.map(r => ({
      student_id: r.id,
      status: r.status,
      note: r.note || '',
    }))
    await api.post(`/attendance/sessions/${sessionID}/submit`, { records })

    Notify.create({ type: 'positive', message: 'Absensi berhasil disimpan' })
    roster.value = []
    selectedSchedule.value = null
  } catch (err) {
    Notify.create({ type: 'negative', message: 'Gagal menyimpan absensi: ' + (err.response?.data?.error || err.message) })
  } finally {
    isSubmitting.value = false
  }
}

onMounted(async () => {
  if (!contextStore.hasContext) await contextStore.initialize()
  await fetchSchedules()
})
</script>

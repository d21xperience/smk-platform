<template>
  <q-page class="q-pa-md">
    <!-- PANEL ACTIONS & STATISTIK JADWAL -->
    <div class="row q-col-gutter-sm q-mb-md items-center">
      <div class="col-12 col-sm-4">
        <q-card flat bordered class="bg-indigo-1 text-indigo-9">
          <q-card-section class="row items-center no-wrap q-py-sm">
            <div>
              <div class="text-caption text-grey-7">Total Agenda Perawatan</div>
              <div class="text-h6 text-weight-bold">{{ store.schedules.length }} Rutinitas</div>
            </div>
            <q-space />
            <q-icon name="date_range" size="sm" />
          </q-card-section>
        </q-card>
      </div>

      <div class="col-12 col-sm-4">
        <q-card flat bordered class="bg-deep-orange-1 text-deep-orange-9">
          <q-card-section class="row items-center no-wrap q-py-sm">
            <div>
              <div class="text-caption text-grey-7">Jatuh Tempo Bulan Ini</div>
              <div class="text-h6 text-weight-bold">{{ totalJatuhTempo }} Aset</div>
            </div>
            <q-space />
            <q-icon name="notification_important" size="sm" />
          </q-card-section>
        </q-card>
      </div>

      <div class="col-12 col-sm-4 text-right">
        <q-btn color="primary" icon="add_alarm" label="Tambah Jadwal Servis" class="full-width-xs" @click="bukaDialogTambah" />
      </div>
    </div>

    <!-- TABEL MANAJEMEN AGENDA PERAWATAN BERKALA -->
    <q-table title="Jadwal Perawatan & Servis Rutin" :rows="store.schedules" :columns="columns" row-key="id"
      :filter="filterText" flat bordered :loading="store.isLoading">
      <template v-slot:top-right>
        <q-input borderless dense debounce="300" v-model="filterText" placeholder="Cari nama aset/ruang...">
          <template v-slot:append><q-icon name="search" /></template>
        </q-input>
      </template>

      <template v-slot:body-cell-frequency="props">
        <q-td :props="props">
          <q-badge :color="getFrekuensiColor(props.value)" class="text-weight-bold">{{ props.value }}</q-badge>
        </q-td>
      </template>

      <template v-slot:body-cell-next_service_date="props">
        <q-td :props="props" class="text-center text-weight-medium">
          <div class="row items-center justify-center no-wrap">
            <q-icon name="event" size="xs" class="q-mr-xs text-grey-6" />
            <span>{{ new Date(props.value).toLocaleDateString('id-ID') }}</span>
            <q-badge v-if="apakahJatuhTempo(props.value)" color="deep-orange" text-color="white" floating rounded />
          </div>
        </q-td>
      </template>

      <template v-slot:body-cell-aksi="props">
        <q-td :props="props" class="text-center">
          <q-btn dense size="sm" color="positive" icon="verified" label="Selesai Servis" @click="konfirmasiSelesaiServis(props.row)">
            <q-tooltip>Catat bahwa servis rutin periode ini sudah beres</q-tooltip>
          </q-btn>
        </q-td>
      </template>
    </q-table>

    <!-- DIALOG TAMBAH JADWAL -->
    <q-dialog v-model="dialogTambahTerbuka" persistent>
      <q-card style="width: 500px; max-width: 90vw;">
        <q-card-section class="bg-primary text-white"><div class="text-h6">Tambah Jadwal Perawatan</div></q-card-section>
        <q-card-section class="q-gutter-y-md">
          <q-input outlined dense v-model="formJadwal.asset_name" label="Nama Program / Aset *" />
          <q-input outlined dense v-model="formJadwal.scope_location" label="Area Cakupan" />
          <q-select outlined dense v-model="formJadwal.frequency" :options="['Bulanan', '3 Bulanan', '6 Bulanan', 'Tahunan']" label="Frekuensi *" />
          <q-input outlined dense v-model="formJadwal.next_service_date" type="date" label="Jadwal Berikutnya *" />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn color="primary" label="Simpan" @click="simpanJadwalBaru" :loading="isSubmitting" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useQuasar } from 'quasar'
import { useMaintenanceStore } from '@/stores/sarpras/maintenance'
import { useContextStore } from '@/stores/contextStore'

const $q = useQuasar()
const store = useMaintenanceStore()
const contextStore = useContextStore()

const filterText = ref('')
const dialogTambahTerbuka = ref(false)
const isSubmitting = ref(false)

const formJadwal = ref({
  asset_name: '',
  scope_location: '',
  frequency: 'Bulanan',
  next_service_date: new Date().toISOString().split('T')[0]
})

const columns = [
  { name: 'schedule_code', label: 'ID Jadwal', align: 'left', field: 'schedule_code', sortable: true },
  { name: 'asset_name', label: 'Nama Program / Aset', align: 'left', field: 'asset_name', sortable: true },
  { name: 'scope_location', label: 'Area Cakupan', align: 'left', field: 'scope_location' },
  { name: 'frequency', label: 'Siklus Berkala', align: 'center', field: 'frequency' },
  { name: 'last_service_date', label: 'Servis Terakhir', align: 'center', field: 'last_service_date', format: val => val ? new Date(val).toLocaleDateString('id-ID') : '-' },
  { name: 'next_service_date', label: 'Jadwal Berikutnya', align: 'center', field: 'next_service_date', sortable: true },
  { name: 'aksi', label: 'Status Eksekusi', align: 'center' }
]

function getFrekuensiColor(frekuensi) {
  if (frekuensi === 'Bulanan') return 'teal'
  if (frekuensi === '3 Bulanan') return 'indigo'
  if (frekuensi === '6 Bulanan') return 'blue-grey'
  return 'purple'
}

function apakahJatuhTempo(tanggalTargetTeks) {
  if (!tanggalTargetTeks) return false
  const targetDate = new Date(tanggalTargetTeks)
  const now = new Date()
  const selisihHari = (targetDate - now) / (1000 * 60 * 60 * 24)
  return selisihHari <= 5
}

const totalJatuhTempo = computed(() => {
  return store.schedules.filter(item => apakahJatuhTempo(item.next_service_date)).length
})

onMounted(() => {
  store.fetchSchedules()
})

watch(() => contextStore.selectedSemesterId, (newVal) => {
  if (newVal) store.fetchSchedules()
})

function bukaDialogTambah() {
  formJadwal.value = { asset_name: '', scope_location: '', frequency: 'Bulanan', next_service_date: new Date().toISOString().split('T')[0] }
  dialogTambahTerbuka.value = true
}

async function simpanJadwalBaru() {
  if (!formJadwal.value.asset_name || !formJadwal.value.frequency) {
    $q.notify({ type: 'negative', message: 'Nama aset dan frekuensi wajib diisi' })
    return
  }
  isSubmitting.value = true
  try {
    await store.createSchedule(formJadwal.value)
    $q.notify({ type: 'positive', message: 'Jadwal berhasil ditambahkan' })
    dialogTambahTerbuka.value = false
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || 'Gagal menambahkan jadwal' })
  } finally {
    isSubmitting.value = false
  }
}

function konfirmasiSelesaiServis(row) {
  $q.dialog({
    title: 'Konfirmasi Selesai Perawatan',
    message: `Apakah Anda ingin memperbarui tanggal pengerjaan rutin untuk "${row.asset_name}"?`,
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await store.completeSchedule(row.id, {
        service_date: new Date().toISOString().split('T')[0],
        cost: 0,
        notes: 'Perawatan rutin selesai'
      })
      $q.notify({ type: 'positive', message: 'Jadwal pemeliharaan berhasil dimutakhirkan' })
    } catch (e) {
      $q.notify({ type: 'negative', message: e.response?.data?.error || 'Gagal memperbarui jadwal' })
    }
  })
}
</script>

<style scoped>
@media (max-width: 599px) {
  .full-width-xs { width: 100%; }
}
</style>

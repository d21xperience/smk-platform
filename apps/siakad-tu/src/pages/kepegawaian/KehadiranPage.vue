<template>
  <q-page class="q-pa-md bg-grey-2">
    <div class="row justify-between items-center q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Rekap Kehadiran GTK</h1>
        <p class="text-caption text-grey-7 q-mb-none">Manajemen Kehadiran Harian & Rekapitulasi Bulanan</p>
      </div>
    </div>

    <q-card class="shadow-1">
      <q-tabs v-model="tab" class="text-primary" align="left">
        <q-tab name="harian" icon="today" label="Kehadiran Harian" />
        <q-tab name="rekap" icon="date_range" label="Rekap Bulanan" />
      </q-tabs>
      <q-separator />

      <q-tab-panels v-model="tab" animated class="bg-white">
        <!-- TAB 1: KEHADIRAN HARIAN -->
        <q-tab-panel name="harian">
          <div class="row items-center q-mb-md">
            <q-input v-model="selectedDate" type="date" label="Pilih Tanggal" outlined dense class="col-12 col-md-3 q-mr-md" />
            <q-btn color="secondary" icon="autorenew" label="Auto Collect" @click="collectData" :loading="isCollecting" class="q-mr-sm" />
            <q-btn color="primary" icon="refresh" label="Refresh" @click="fetchHarian" :loading="isLoading" flat />
          </div>

          <q-table :rows="harianList" :columns="harianColumns" row-key="ptkId" flat bordered dense :loading="isLoading">
            <template v-slot:body-cell-status="props">
              <q-td :props="props" class="text-center">
                <q-badge :color="getStatusColor(props.row.status)" class="q-pa-sm">
                  {{ props.row.status }}
                </q-badge>
              </q-td>
            </template>
            <template v-slot:body-cell-jamMasuk="props">
              <q-td :props="props" class="text-center">{{ props.row.jamMasuk || '-' }}</q-td>
            </template>
            <template v-slot:body-cell-jamPulang="props">
              <q-td :props="props" class="text-center">{{ props.row.jamPulang || '-' }}</q-td>
            </template>
            <template v-slot:body-cell-aksi="props">
              <q-td :props="props" class="text-center">
                <q-btn flat round dense icon="edit" color="primary" @click="openManualDialog(props.row)">
                  <q-tooltip>Input Manual</q-tooltip>
                </q-btn>
              </q-td>
            </template>
          </q-table>
        </q-tab-panel>

        <!-- TAB 2: REKAP BULANAN -->
        <q-tab-panel name="rekap">
          <div class="row items-center q-mb-md">
            <q-btn color="primary" icon="add" label="Generate Rekap Baru" @click="openGenerateDialog" />
            <q-space />
            <q-btn color="secondary" icon="refresh" label="Refresh" @click="fetchRekap" :loading="isLoading" flat />
          </div>

          <q-table :rows="rekapList" :columns="rekapColumns" row-key="id" flat bordered dense :loading="isLoading">
            <template v-slot:body-cell-bulan="props">
              <q-td :props="props">{{ getMonthName(props.row.bulan) }} {{ props.row.tahun }}</q-td>
            </template>
            <template v-slot:body-cell-rataPersentase="props">
              <q-td :props="props" class="text-center">
                <q-badge :color="props.row.rataPersentase >= 80 ? 'green-6' : 'orange-6'" class="q-pa-sm">
                  {{ props.row.rataPersentase }}%
                </q-badge>
              </q-td>
            </template>
            <template v-slot:body-cell-status="props">
              <q-td :props="props" class="text-center">
                <q-badge :color="props.row.status === 'COMPLETED' ? 'green-8' : 'blue-6'" class="q-pa-sm">
                  {{ props.row.status }}
                </q-badge>
              </q-td>
            </template>
            <template v-slot:body-cell-aksi="props">
              <q-td :props="props" class="text-center">
                <q-btn flat round dense icon="download" color="secondary" :disable="!props.row.excelUrl" @click="downloadExcel(props.row.excelUrl)">
                  <q-tooltip>Download Excel</q-tooltip>
                </q-btn>
              </q-td>
            </template>
          </q-table>
        </q-tab-panel>
      </q-tab-panels>
    </q-card>

    <!-- Dialog Manual Input -->
    <q-dialog v-model="manualDialog" persistent>
      <q-card style="min-width: 400px">
        <q-card-section class="bg-primary text-white"><div class="text-h6">Input Manual Kehadiran</div></q-card-section>
        <q-card-section class="q-gutter-y-md">
          <q-input v-model="form.namaPegawai" label="Nama Pegawai" outlined dense readonly />
          <q-select v-model="form.status" :options="['HADIR', 'TERLAMBAT', 'IZIN_RESMI', 'SAKIT', 'ALPHA', 'DINAS_LUAR']" label="Status Kehadiran" outlined dense />
          <div class="row q-col-gutter-sm">
            <div class="col-6"><q-input v-model="form.jamMasuk" label="Jam Masuk" type="time" outlined dense /></div>
            <div class="col-6"><q-input v-model="form.jamPulang" label="Jam Pulang" type="time" outlined dense /></div>
          </div>
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" color="negative" v-close-popup />
          <q-btn color="primary" label="Simpan" @click="submitManual" :loading="isSubmitting" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Dialog Generate Rekap -->
    <q-dialog v-model="generateDialog" persistent>
      <q-card style="min-width: 400px">
        <q-card-section class="bg-primary text-white"><div class="text-h6">Generate Rekap Bulanan</div></q-card-section>
        <q-card-section class="q-gutter-y-md">
          <q-select v-model="generateForm.bulan" :options="monthOptions" label="Pilih Bulan" outlined dense emit-value map-options />
          <q-input v-model.number="generateForm.tahun" label="Tahun" type="number" outlined dense />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" color="negative" v-close-popup />
          <q-btn color="primary" label="Generate" @click="submitGenerate" :loading="isSubmitting" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { api } from '@/boot/axios'

const $q = useQuasar()
const tab = ref('harian')
const isLoading = ref(false)
const isCollecting = ref(false)
const isSubmitting = ref(false)

const selectedDate = ref(new Date().toISOString().split('T')[0])
const harianList = ref([])
const rekapList = ref([])

const manualDialog = ref(false)
const generateDialog = ref(false)

const form = ref({ ptkId: '', namaPegawai: '', status: 'HADIR', jamMasuk: '', jamPulang: '' })
const generateForm = ref({ bulan: new Date().getMonth() + 1, tahun: new Date().getFullYear() })

const monthOptions = [
  { label: 'Januari', value: 1 }, { label: 'Februari', value: 2 }, { label: 'Maret', value: 3 },
  { label: 'April', value: 4 }, { label: 'Mei', value: 5 }, { label: 'Juni', value: 6 },
  { label: 'Juli', value: 7 }, { label: 'Agustus', value: 8 }, { label: 'September', value: 9 },
  { label: 'Oktober', value: 10 }, { label: 'November', value: 11 }, { label: 'Desember', value: 12 }
]

const harianColumns = [
  { name: 'namaPegawai', label: 'Nama Pegawai', field: 'namaPegawai', align: 'left', sortable: true },
  { name: 'nip', label: 'NIP', field: 'nip', align: 'left' },
  { name: 'jenisGtk', label: 'Jenis', field: 'jenisGtk', align: 'center' },
  { name: 'status', label: 'Status', align: 'center' },
  { name: 'jamMasuk', label: 'Jam Masuk', align: 'center' },
  { name: 'jamPulang', label: 'Jam Pulang', align: 'center' },
  { name: 'sourceData', label: 'Sumber', field: 'sourceData', align: 'center' },
  { name: 'aksi', label: 'Aksi', align: 'center' }
]

const rekapColumns = [
  { name: 'bulan', label: 'Bulan', align: 'left', sortable: true },
  { name: 'totalPegawai', label: 'Total Pegawai', field: 'totalPegawai', align: 'center' },
  { name: 'rataPersentase', label: 'Rata-rata Kehadiran', align: 'center' },
  { name: 'status', label: 'Status', align: 'center' },
  { name: 'aksi', label: 'Aksi', align: 'center' }
]

const getMonthName = (bulan) => {
  return monthOptions.find(m => m.value === bulan)?.label || bulan
}

const getStatusColor = (status) => {
  switch (status) {
    case 'HADIR': return 'green-6'
    case 'TERLAMBAT': return 'orange-6'
    case 'IZIN_RESMI': case 'SAKIT': return 'blue-6'
    case 'ALPHA': return 'red-6'
    case 'DINAS_LUAR': return 'purple-6'
    default: return 'grey-6'
  }
}

const fetchHarian = async () => {
  isLoading.value = true
  try {
    const res = await api.get(`/kepegawaian/kehadiran/harian?tanggal=${selectedDate.value}`)
    harianList.value = res.data || []
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Gagal memuat data kehadiran harian' })
    console.error(err)
  } finally {
    isLoading.value = false
  }
}

const fetchRekap = async () => {
  isLoading.value = true
  try {
    const res = await api.get('/kepegawaian/kehadiran/rekap-bulanan')
    rekapList.value = res.data || []
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Gagal memuat data rekap bulanan' })
    console.error(err)
  } finally {
    isLoading.value = false
  }
}

const collectData = async () => {
  isCollecting.value = true
  try {
    const res = await api.post(`/kepegawaian/kehadiran/collect?tanggal=${selectedDate.value}`)
    $q.notify({ type: 'positive', message: `Auto-collect selesai. Total: ${res.data.total} pegawai.` })
    await fetchHarian()
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Gagal melakukan auto-collect' })
    console.error(err)
  } finally {
    isCollecting.value = false
  }
}

const openManualDialog = (row) => {
  form.value = {
    ptkId: row.ptkId,
    namaPegawai: row.namaPegawai,
    status: row.status === 'ALPHA' ? 'HADIR' : row.status,
    jamMasuk: row.jamMasuk || '',
    jamPulang: row.jamPulang || ''
  }
  manualDialog.value = true
}

const submitManual = async () => {
  isSubmitting.value = true
  try {
    await api.post('/kepegawaian/kehadiran/manual', {
      ptkId: form.value.ptkId,
      tanggal: selectedDate.value,
      status: form.value.status,
      jamMasuk: form.value.jamMasuk || null,
      jamPulang: form.value.jamPulang || null
    })
    $q.notify({ type: 'positive', message: 'Data kehadiran berhasil disimpan' })
    manualDialog.value = false
    await fetchHarian()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Gagal menyimpan data' })
    console.error(err)
  } finally {
    isSubmitting.value = false
  }
}

const openGenerateDialog = () => {
  generateForm.value = { bulan: new Date().getMonth() + 1, tahun: new Date().getFullYear() }
  generateDialog.value = true
}

const submitGenerate = async () => {
  isSubmitting.value = true
  try {
    await api.post('/kepegawaian/kehadiran/rekap-bulanan/generate', generateForm.value)
    $q.notify({ type: 'positive', message: 'Rekap berhasil dibuat' })
    generateDialog.value = false
    await fetchRekap()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Gagal generate rekap' })
    console.error(err)
  } finally {
    isSubmitting.value = false
  }
}

const downloadExcel = (url) => {
  if (url) window.open(url, '_blank')
}

onMounted(() => {
  fetchHarian()
  fetchRekap()
})
</script>

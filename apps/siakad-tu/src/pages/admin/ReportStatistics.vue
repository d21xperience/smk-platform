<template>
  <q-page class="q-pa-md">
    <div class="text-h5 q-mb-md">Laporan & Statistik</div>

    <q-tabs v-model="tab" dense class="q-mb-md" align="left">
      <q-tab name="overview" label="Laporan Keseluruhan" />
      <q-tab name="per-siswa" label="Per Siswa" />
      <q-tab name="per-kelas" label="Per Kelas" />
    </q-tabs>

    <!-- Filter umum (tahun ajaran, ujian) -->
    <q-card flat bordered class="q-mb-md">
      <q-card-section>
        <div class="row q-col-gutter-md">
          <div class="col-12 col-md-4">
            <q-select v-model="filterUjian" :options="examOptions" label="Pilih Ujian" clearable outlined dense
              @update:model-value="loadData" />
          </div>
          <div class="col-12 col-md-3">
            <q-select v-model="filterTahunAjaran" :options="tahunAjaranOptions" label="Tahun Ajaran" outlined dense
              @update:model-value="loadData" />
          </div>
          <div class="col-12 col-md-3">
            <q-btn color="primary" icon="file_download" label="Export Excel" @click="exportExcel" />
          </div>
        </div>
      </q-card-section>
    </q-card>

    <q-tab-panels v-model="tab" animated>
      <!-- Tab Laporan Keseluruhan -->
      <q-tab-panel name="overview" class="q-pa-none">
        <div class="row q-col-gutter-md q-mb-md">
          <div class="col-12 col-md-3">
            <q-card class="bg-primary text-white">
              <q-card-section>
                <div class="text-h6">Rata-rata Nilai</div>
                <div class="text-h3">{{ overviewStats.rataNilai }}%</div>
              </q-card-section>
            </q-card>
          </div>
          <div class="col-12 col-md-3">
            <q-card class="bg-positive text-white">
              <q-card-section>
                <div class="text-h6">Tingkat Kelulusan</div>
                <div class="text-h3">{{ overviewStats.tingkatKelulusan }}%</div>
              </q-card-section>
            </q-card>
          </div>
          <div class="col-12 col-md-3">
            <q-card class="bg-info text-white">
              <q-card-section>
                <div class="text-h6">Total Peserta</div>
                <div class="text-h3">{{ overviewStats.totalPeserta }}</div>
              </q-card-section>
            </q-card>
          </div>
          <div class="col-12 col-md-3">
            <q-card class="bg-warning text-black">
              <q-card-section>
                <div class="text-h6">Nilai Tertinggi</div>
                <div class="text-h3">{{ overviewStats.nilaiTertinggi }}</div>
              </q-card-section>
            </q-card>
          </div>
        </div>
        <div class="row q-col-gutter-md">
          <div class="col-12 col-md-6">
            <q-card>
              <q-card-section>
                <div class="text-h6">Distribusi Nilai Keseluruhan</div>
                <div ref="chartDistribusi" style="height: 300px;"></div>
              </q-card-section>
            </q-card>
          </div>
          <div class="col-12 col-md-6">
            <q-card>
              <q-card-section>
                <div class="text-h6">Performa per Mata Pelajaran</div>
                <div ref="chartMapel" style="height: 300px;"></div>
              </q-card-section>
            </q-card>
          </div>
        </div>
      </q-tab-panel>

      <!-- Tab Per Siswa -->
      <q-tab-panel name="per-siswa" class="q-pa-none">
        <div class="row q-col-gutter-md q-mb-md">
          <div class="col-12 col-md-4">
            <q-select v-model="filterSiswa" :options="siswaOptions" label="Pilih Siswa" option-value="id"
              option-label="nama" clearable outlined dense @update:model-value="loadPerSiswa" />
          </div>
        </div>
        <q-table :rows="siswaReportRows" :columns="siswaColumns" row-key="id" flat bordered dense :loading="loading" />
      </q-tab-panel>

      <!-- Tab Per Kelas -->
      <q-tab-panel name="per-kelas" class="q-pa-none">
        <div class="row q-col-gutter-md q-mb-md">
          <div class="col-12 col-md-4">
            <q-select v-model="filterKelas" :options="kelasOptions" label="Pilih Kelas" option-value="id"
              option-label="nama" clearable outlined dense @update:model-value="loadPerKelas" />
          </div>
        </div>
        <q-table :rows="kelasReportRows" :columns="kelasColumns" row-key="id" flat bordered dense :loading="loading" />
      </q-tab-panel>
    </q-tab-panels>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useQuasar } from 'quasar'
import * as echarts from 'echarts'
import { utils, writeFile } from 'xlsx'

const $q = useQuasar()

// ==================== DATA DUMMY ====================
// Ujian (dari ExamManagement)
const examList = ref([
  { id: 1, nama: 'UTS Matematika Ganjil', mapel_id: 1, mapel_nama: 'Matematika', tanggal_mulai: '2025-06-10T08:00' },
  { id: 2, nama: 'UAS Fisika', mapel_id: 2, mapel_nama: 'Fisika', tanggal_mulai: '2025-06-15T09:00' }
])

// Siswa (dari UserManagement)
const siswaList = ref([
  { id: 1, nama: 'Ahmad Faizal', kelas_id: 1, kelas_nama: '10 IPA 1' },
  { id: 2, nama: 'Siti Nurhaliza', kelas_id: 1, kelas_nama: '10 IPA 1' },
  { id: 3, nama: 'Budi Santoso', kelas_id: 2, kelas_nama: '10 IPA 2' }
])

// Kelas (dari ClassSubjectManagement)
const kelasList = ref([
  { id: 1, nama: '10 IPA 1' },
  { id: 2, nama: '10 IPA 2' },
  { id: 3, nama: '11 IPA 1' }
])

// Nilai dummy (siswa_id, ujian_id, nilai)
const nilaiList = ref([
  { siswa_id: 1, ujian_id: 1, nilai: 85 },
  { siswa_id: 1, ujian_id: 2, nilai: 78 },
  { siswa_id: 2, ujian_id: 1, nilai: 92 },
  { siswa_id: 2, ujian_id: 2, nilai: 88 },
  { siswa_id: 3, ujian_id: 1, nilai: 70 },
  { siswa_id: 3, ujian_id: 2, nilai: 65 }
])

// Tahun ajaran
const tahunAjaranOptions = ref(['2024/2025', '2025/2026'])
const filterTahunAjaran = ref('2025/2026')
const filterUjian = ref(null)

const tab = ref('overview')
const loading = ref(false)

// Data untuk chart dan stats keseluruhan
const overviewStats = ref({
  rataNilai: 0,
  tingkatKelulusan: 0,
  totalPeserta: 0,
  nilaiTertinggi: 0
})

const examOptions = computed(() => examList.value.map(e => ({ label: e.nama, value: e.id })))
const siswaOptions = computed(() => siswaList.value.map(s => ({ id: s.id, nama: s.nama })))
const kelasOptions = computed(() => kelasList.value)

// Filtered nilai berdasarkan ujian dan tahun ajaran (tahun disini hanya dummy)
const filteredNilai = computed(() => {
  let data = nilaiList.value
  if (filterUjian.value) {
    data = data.filter(n => n.ujian_id === filterUjian.value)
  }
  // asumsi tahun ajaran berdasarkan tanggal ujian, untuk dummy kita abaikan dulu
  return data
})

// Laporan keseluruhan
function computeOverview() {
  const data = filteredNilai.value
  if (data.length === 0) {
    overviewStats.value = { rataNilai: 0, tingkatKelulusan: 0, totalPeserta: 0, nilaiTertinggi: 0 }
    return
  }
  const totalNilai = data.reduce((sum, n) => sum + n.nilai, 0)
  const rata = totalNilai / data.length
  const lulus = data.filter(n => n.nilai >= 70).length
  const tingkat = (lulus / data.length) * 100
  const tertinggi = Math.max(...data.map(n => n.nilai))
  const uniqueSiswa = new Set(data.map(n => n.siswa_id)).size
  overviewStats.value = {
    rataNilai: rata.toFixed(1),
    tingkatKelulusan: tingkat.toFixed(1),
    totalPeserta: uniqueSiswa,
    nilaiTertinggi: tertinggi
  }
  renderCharts()
}

function renderCharts() {
  // Distribusi nilai (histogram)
  const bins = { '0-49': 0, '50-59': 0, '60-69': 0, '70-79': 0, '80-89': 0, '90-100': 0 }
  filteredNilai.value.forEach(n => {
    if (n.nilai < 50) bins['0-49']++
    else if (n.nilai < 60) bins['50-59']++
    else if (n.nilai < 70) bins['60-69']++
    else if (n.nilai < 80) bins['70-79']++
    else if (n.nilai < 90) bins['80-89']++
    else bins['90-100']++
  })
  const chartDom = document.getElementById('chartDistribusi')
  if (chartDom) {
    const chart = echarts.init(chartDom)
    chart.setOption({
      tooltip: { trigger: 'axis' },
      xAxis: { type: 'category', data: Object.keys(bins) },
      yAxis: { type: 'value', name: 'Jumlah Siswa' },
      series: [{ type: 'bar', data: Object.values(bins), itemStyle: { color: '#1976D2' } }]
    })
  }
  // Performa per mapel (rata-rata nilai per mata pelajaran)
  const mapelData = {}
  filteredNilai.value.forEach(n => {
    const exam = examList.value.find(e => e.id === n.ujian_id)
    if (exam) {
      if (!mapelData[exam.mapel_nama]) mapelData[exam.mapel_nama] = { total: 0, count: 0 }
      mapelData[exam.mapel_nama].total += n.nilai
      mapelData[exam.mapel_nama].count++
    }
  })
  const mapelNames = Object.keys(mapelData)
  const avgValues = mapelNames.map(m => (mapelData[m].total / mapelData[m].count).toFixed(1))
  const chartMapelDom = document.getElementById('chartMapel')
  if (chartMapelDom) {
    const chart2 = echarts.init(chartMapelDom)
    chart2.setOption({
      tooltip: { trigger: 'axis' },
      xAxis: { type: 'category', data: mapelNames },
      yAxis: { type: 'value', name: 'Rata-rata Nilai' },
      series: [{ type: 'bar', data: avgValues, itemStyle: { color: '#9C27B0' } }]
    })
  }
}

// Laporan per siswa
const filterSiswa = ref(null)
const siswaReportRows = ref([])
const siswaColumns = [
  { name: 'ujian', label: 'Ujian', field: 'ujian_nama', align: 'left' },
  { name: 'nilai', label: 'Nilai', field: 'nilai', align: 'center' },
  { name: 'status', label: 'Status', field: 'status', align: 'center' }
]
function loadPerSiswa() {
  if (!filterSiswa.value) {
    siswaReportRows.value = []
    return
  }
  loading.value = true
  setTimeout(() => {
    const ujianSiswa = nilaiList.value.filter(n => n.siswa_id === filterSiswa.value.id)
    siswaReportRows.value = ujianSiswa.map(n => {
      const exam = examList.value.find(e => e.id === n.ujian_id)
      return {
        ujian_nama: exam?.nama || '-',
        nilai: n.nilai,
        status: n.nilai >= 70 ? 'Lulus' : 'Tidak Lulus'
      }
    })
    loading.value = false
  }, 300)
}

// Laporan per kelas
const filterKelas = ref(null)
const kelasReportRows = ref([])
const kelasColumns = [
  { name: 'siswa', label: 'Nama Siswa', field: 'siswa_nama', align: 'left' },
  { name: 'rata_nilai', label: 'Rata-rata Nilai', field: 'rata_nilai', align: 'center' },
  { name: 'status', label: 'Status Akhir', field: 'status', align: 'center' }
]
function loadPerKelas() {
  if (!filterKelas.value) {
    kelasReportRows.value = []
    return
  }
  loading.value = true
  setTimeout(() => {
    const siswaDiKelas = siswaList.value.filter(s => s.kelas_id === filterKelas.value.id)
    const rows = siswaDiKelas.map(s => {
      const nilaiSiswa = nilaiList.value.filter(n => n.siswa_id === s.id)
      const rata = nilaiSiswa.length ? nilaiSiswa.reduce((a, b) => a + b.nilai, 0) / nilaiSiswa.length : 0
      const status = rata >= 70 ? 'Lulus' : 'Tidak Lulus'
      return { id: s.id, siswa_nama: s.nama, rata_nilai: rata.toFixed(1), status }
    })
    kelasReportRows.value = rows
    loading.value = false
  }, 300)
}

// Export to Excel
function exportExcel() {
  let dataToExport = []
  if (tab.value === 'overview') {
    dataToExport = filteredNilai.value.map(n => {
      const siswa = siswaList.value.find(s => s.id === n.siswa_id)
      const exam = examList.value.find(e => e.id === n.ujian_id)
      return {
        'Siswa': siswa?.nama || '-',
        'Kelas': siswa?.kelas_nama || '-',
        'Ujian': exam?.nama || '-',
        'Nilai': n.nilai,
        'Status': n.nilai >= 70 ? 'Lulus' : 'Tidak Lulus'
      }
    })
  } else if (tab.value === 'per-siswa' && filterSiswa.value) {
    dataToExport = siswaReportRows.value.map(r => ({
      'Ujian': r.ujian_nama,
      'Nilai': r.nilai,
      'Status': r.status
    }))
  } else if (tab.value === 'per-kelas' && filterKelas.value) {
    dataToExport = kelasReportRows.value.map(r => ({
      'Siswa': r.siswa_nama,
      'Rata-rata Nilai': r.rata_nilai,
      'Status': r.status
    }))
  }
  if (dataToExport.length === 0) {
    $q.notify({ type: 'warning', message: 'Tidak ada data untuk diexport' })
    return
  }
  const ws = utils.json_to_sheet(dataToExport)
  const wb = utils.book_new()
  utils.book_append_sheet(wb, ws, 'Laporan')
  writeFile(wb, `laporan_${tab.value}_${Date.now()}.xlsx`)
}

function loadData() {
  computeOverview()
  if (tab.value === 'per-siswa') loadPerSiswa()
  if (tab.value === 'per-kelas') loadPerKelas()
}

watch(filterUjian, () => loadData())
watch(tab, () => {
  if (tab.value === 'overview') computeOverview()
  if (tab.value === 'per-siswa') loadPerSiswa()
  if (tab.value === 'per-kelas') loadPerKelas()
})

onMounted(() => {
  loadData()
})
</script>

<style scoped></style>

<template>
  <q-page class="q-pa-md">
    <!-- Judul halaman -->
    <div class="text-h5 q-mb-md">Dashboard Admin</div>

    <!-- Grid kartu statistik -->
    <div class="row q-col-gutter-md q-mb-lg">
      <div class="col-xs-12 col-sm-6 col-md-3">
        <q-card class="bg-primary text-white">
          <q-card-section>
            <div class="text-h6">Total Siswa</div>
            <div class="text-h3">{{ stats.siswa }}</div>
            <q-icon name="people" size="32px" class="absolute-bottom-right q-ma-md" />
          </q-card-section>
        </q-card>
      </div>
      <div class="col-xs-12 col-sm-6 col-md-3">
        <q-card class="bg-secondary text-white">
          <q-card-section>
            <div class="text-h6">Total Guru</div>
            <div class="text-h3">{{ stats.guru }}</div>
            <q-icon name="school" size="32px" class="absolute-bottom-right q-ma-md" />
          </q-card-section>
        </q-card>
      </div>
      <div class="col-xs-12 col-sm-6 col-md-3">
        <q-card class="bg-accent text-white">
          <q-card-section>
            <div class="text-h6">Admin Aktif</div>
            <div class="text-h3">{{ stats.admin }}</div>
            <q-icon name="admin_panel_settings" size="32px" class="absolute-bottom-right q-ma-md" />
          </q-card-section>
        </q-card>
      </div>
      <div class="col-xs-12 col-sm-6 col-md-3">
        <q-card class="bg-positive text-white">
          <q-card-section>
            <div class="text-h6">Ujian Aktif</div>
            <div class="text-h3">{{ stats.ujianAktif }}</div>
            <q-icon name="quiz" size="32px" class="absolute-bottom-right q-ma-md" />
          </q-card-section>
        </q-card>
      </div>
      <div class="col-xs-12 col-sm-6 col-md-3">
        <q-card class="bg-info text-white">
          <q-card-section>
            <div class="text-h6">Bank Soal</div>
            <div class="text-h3">{{ stats.totalSoal }}</div>
            <q-icon name="edit_note" size="32px" class="absolute-bottom-right q-ma-md" />
          </q-card-section>
        </q-card>
      </div>
      <div class="col-xs-12 col-sm-6 col-md-3">
        <q-card class="bg-warning text-black">
          <q-card-section>
            <div class="text-h6">Rata-rata Kelulusan</div>
            <div class="text-h3">{{ stats.rataKelulusan }}%</div>
            <q-icon name="trending_up" size="32px" class="absolute-bottom-right q-ma-md" />
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Grafik Tren Ujian per Bulan -->
    <div class="row q-col-gutter-md">
      <div class="col-12 col-md-8">
        <q-card>
          <q-card-section>
            <div class="text-h6">Tren Jumlah Ujian (per bulan)</div>
          </q-card-section>
          <q-card-section>
            <div ref="chartTren" style="height: 300px; width: 100%;"></div>
          </q-card-section>
        </q-card>
      </div>

      <!-- Grafik Distribusi Nilai (Ringkasan) -->
      <div class="col-12 col-md-4">
        <q-card>
          <q-card-section>
            <div class="text-h6">Distribusi Nilai Ujian</div>
          </q-card-section>
          <q-card-section>
            <div ref="chartDistribusi" style="height: 300px; width: 100%;"></div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Tabel Ujian Terbaru -->
    <div class="row q-mt-md">
      <div class="col-12">
        <q-card>
          <q-card-section>
            <div class="text-h6">Ujian Terbaru</div>
          </q-card-section>
          <q-card-section>
            <q-table :rows="ujianTerbaru" :columns="columnsUjian" row-key="id" :pagination="{ rowsPerPage: 5 }" dense
              flat bordered>
              <template v-slot:body-cell-status="props">
                <q-td :props="props">
                  <q-badge :color="props.row.status === 'aktif' ? 'positive' : 'grey'" :label="props.row.status" />
                </q-td>
              </template>
            </q-table>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import * as echarts from 'echarts'

// Data statistik (nanti diganti dengan data dari API)
const stats = ref({
  siswa: 1250,
  guru: 48,
  admin: 5,
  ujianAktif: 12,
  totalSoal: 3420,
  rataKelulusan: 78.5
})

// Data ujian terbaru (contoh)
const ujianTerbaru = ref([
  { id: 1, nama: 'UTS Matematika', kelas: '12 IPA 1', tanggal: '2025-05-10', peserta: 32, status: 'aktif' },
  { id: 2, nama: 'UAS Bahasa Inggris', kelas: '11 IPS 2', tanggal: '2025-05-09', peserta: 28, status: 'aktif' },
  { id: 3, nama: 'Remidi Fisika', kelas: '12 IPA 2', tanggal: '2025-05-07', peserta: 15, status: 'selesai' },
  { id: 4, nama: 'Try Out UN', kelas: '12 IPA/IPS', tanggal: '2025-05-01', peserta: 210, status: 'selesai' },
  { id: 5, nama: 'PTS Kimia', kelas: '11 IPA 1', tanggal: '2025-04-28', peserta: 30, status: 'selesai' },
])

const columnsUjian = [
  { name: 'nama', label: 'Nama Ujian', field: 'nama', align: 'left' },
  { name: 'kelas', label: 'Kelas', field: 'kelas', align: 'left' },
  { name: 'tanggal', label: 'Tanggal', field: 'tanggal', align: 'left' },
  { name: 'peserta', label: 'Peserta', field: 'peserta', align: 'center' },
  { name: 'status', label: 'Status', field: 'status', align: 'center' },
]

// Ref untuk container chart
const chartTren = ref(null)
const chartDistribusi = ref(null)

onMounted(() => {
  // Chart Tren Ujian per Bulan
  const trenChart = echarts.init(chartTren.value)
  trenChart.setOption({
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', data: ['Jan', 'Feb', 'Mar', 'Apr', 'Mei', 'Jun', 'Jul', 'Agu', 'Sep', 'Okt', 'Nov', 'Des'] },
    yAxis: { type: 'value', name: 'Jumlah Ujian' },
    series: [{
      data: [5, 7, 9, 12, 15, 18, 20, 22, 19, 16, 14, 10],
      type: 'line',
      smooth: true,
      lineStyle: { color: '#1976D2', width: 3 },
      areaStyle: { opacity: 0.2, color: '#1976D2' },
      symbol: 'circle',
      symbolSize: 8
    }]
  })

  // Chart Distribusi Nilai (pie chart)
  const distribusiChart = echarts.init(chartDistribusi.value)
  distribusiChart.setOption({
    tooltip: { trigger: 'item' },
    legend: { orient: 'vertical', left: 'left' },
    series: [{
      name: 'Distribusi Nilai',
      type: 'pie',
      radius: '50%',
      data: [
        { value: 15, name: 'A (90-100)' },
        { value: 25, name: 'B (75-89)' },
        { value: 30, name: 'C (60-74)' },
        { value: 20, name: 'D (50-59)' },
        { value: 10, name: 'E (<50)' }
      ],
      emphasis: { scale: true },
      label: { show: true, formatter: '{b}: {d}%' }
    }]
  })

  // Responsif saat window resize
  window.addEventListener('resize', () => {
    trenChart.resize()
    distribusiChart.resize()
  })
})
</script>

<style scoped>
.bg-accent {
  background: #9C27B0;
}
</style>

<template>
  <div>
    <div class="row q-col-gutter-md q-mb-md">
      <div class="col-12 col-md-3">
        <q-card class="bg-primary text-white">
          <q-card-section>
            <div class="text-h6">Total Siswa</div>
            <div class="text-h3">{{ totalSiswa }}</div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-md-3">
        <q-card class="bg-secondary text-white">
          <q-card-section>
            <div class="text-h6">Total Guru</div>
            <div class="text-h3">{{ totalGuru }}</div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-md-3">
        <q-card class="bg-accent text-white">
          <q-card-section>
            <div class="text-h6">Rata-rata Nilai</div>
            <div class="text-h3">{{ rataNilai }}</div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-md-3">
        <q-card class="bg-info text-white">
          <q-card-section>
            <div class="text-h6">Kehadiran Hari Ini</div>
            <div class="text-h3">{{ kehadiran }}%</div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <div class="row q-col-gutter-md">
      <div class="col-12 col-md-6">
        <q-card>
          <q-card-section>
            <div class="text-h6">Rata-rata Nilai per Mata Pelajaran</div>
            <div ref="chartMapel" style="height: 300px;"></div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-md-6">
        <q-card>
          <q-card-section>
            <div class="text-h6">Ujian Terdekat</div>
            <q-list dense>
              <q-item v-for="ujian in ujianTerdekat" :key="ujian.id">
                <q-item-section>
                  <q-item-label>{{ ujian.nama }} - {{ ujian.kelas_nama }}</q-item-label>
                  <q-item-label caption>{{ formatDate(ujian.tanggal_mulai) }}</q-item-label>
                </q-item-section>
              </q-item>
              <q-item v-if="ujianTerdekat.length === 0">
                <q-item-section>Tidak ada ujian terdekat</q-item-section>
              </q-item>
            </q-list>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <div class="row q-col-gutter-md q-mt-md">
      <div class="col-12">
        <q-card>
          <q-card-section>
            <div class="text-h6">Alert Kritis</div>
            <div v-if="kritisList.length === 0" class="text-positive">Tidak ada alert kritis</div>
            <q-list v-else>
              <q-item v-for="alert in kritisList" :key="alert.id" class="bg-negative text-white">
                <q-item-section>
                  <q-item-label>{{ alert.judul }}</q-item-label>
                  <q-item-label caption>{{ alert.deskripsi }}</q-item-label>
                </q-item-section>
              </q-item>
            </q-list>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import * as echarts from 'echarts'
import { useDashboardStore } from 'stores/dashboard'
import { useSiswaStore } from 'stores/kesiswaan/siswa'
import { useGuruStore } from 'stores/kepegawaian/guru'
import { useAbsensiStore } from '@/stores/kesiswaan/absensiStore'
import { useNilaiStore } from 'stores/kurikulum/nilai'
import { useUjianStore } from 'stores/akademik/ujian' // asumsi ada

const dashboardStore = useDashboardStore()
const siswaStore = useSiswaStore()
const guruStore = useGuruStore()
const absensiStore = useAbsensiStore()
const nilaiStore = useNilaiStore()
const ujianStore = useUjianStore()

const totalSiswa = ref(0)
const totalGuru = ref(0)
const rataNilai = ref(0)
const kehadiran = ref(0)
const ujianTerdekat = ref([])
const kritisList = ref([])

function loadData() {
  totalSiswa.value = siswaStore.list.length
  totalGuru.value = guruStore.list.length
  rataNilai.value = dashboardStore.getRataRataNilaiKeseluruhan()
  kehadiran.value = dashboardStore.getPersentaseKehadiranHariIni()
  ujianTerdekat.value = dashboardStore.getUjianTerdekat() || []
  // contoh alert kritis: siswa dengan alpha > 5 dalam sebulan
  const now = new Date()
  const bulanIni = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`
  const alphaSiswa = absensiStore.list.filter(a => a.status === 'alpha' && a.tanggal.startsWith(bulanIni))
  const countPerSiswa = {}
  alphaSiswa.forEach(a => { countPerSiswa[a.siswa_id] = (countPerSiswa[a.siswa_id] || 0) + 1 })
  kritisList.value = Object.entries(countPerSiswa).filter(([, count]) => count >= 5).map(([siswa_id, count]) => {
    const siswa = siswaStore.list.find(s => s.id == siswa_id)
    return { id: siswa_id, judul: `Alpha Tinggi`, deskripsi: `${siswa?.nama} telah alpha ${count} kali bulan ini` }
  })
}

function renderCharts() {
  const data = dashboardStore.getGrafikNilaiPerMapel()
  const chartDom = document.getElementById('chartMapel')
  if (chartDom) {
    const chart = echarts.init(chartDom)
    chart.setOption({
      tooltip: { trigger: 'axis' },
      xAxis: { type: 'category', data: data.map(d => d.mapel) },
      yAxis: { type: 'value', name: 'Rata-rata Nilai' },
      series: [{ type: 'bar', data: data.map(d => d.rata), itemStyle: { color: '#1976D2' } }]
    })
  }
}

function formatDate(iso) {
  return new Date(iso).toLocaleDateString('id-ID')
}

onMounted(() => {
  siswaStore.loadData()
  guruStore.loadData()
  absensiStore.loadData()
  nilaiStore.loadData()
  if (ujianStore.loadData) ujianStore.loadData()
  loadData()
  renderCharts()
})
</script>

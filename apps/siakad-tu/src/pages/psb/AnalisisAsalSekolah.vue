<template>
  <q-page class="q-pa-md bg-grey-2">
    <!-- Header Halaman -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Analisis Asal Sekolah Pendaftar</h1>
        <p class="text-caption text-grey-7 q-mb-none">Metrik kontribusi pendaftar, sebaran wilayah, dan rata-rata nilai
          berdasarkan SMP/MTs asal calon siswa</p>
      </div>
      <div>
        <q-btn color="indigo-9" icon="share" label="Bagikan ke Tim Promosi" @click="shareAnalysis" flat />
      </div>
    </div>

    <!-- Row Visualisasi Grafik & Summary Card -->
    <div class="row q-col-gutter-md q-mb-md">
      <!-- Grafik Batang Kontribusi Top Sekolah Pemasok -->
      <div class="col-12 col-md-7">
        <q-card class="shadow-1 fit">
          <q-card-section>
            <div class="text-subtitle2 text-weight-bold text-grey-9">Top 5 Sekolah Kontributor Terbesar</div>
            <div class="text-caption text-grey-6">Jumlah pendaftar terbanyak berdasarkan SMP/MTs asal</div>
          </q-card-section>
          <q-card-section class="q-pt-none">
            <apexchart type="bar" height="230" :options="chartOptions" :series="chartSeries"></apexchart>
          </q-card-section>
        </q-card>
      </div>

      <!-- Informasi Statistik Rekor Sekolah Pemasok -->
      <div class="col-12 col-md-5">
        <div class="row q-col-gutter-sm fit">
          <div class="col-12">
            <q-card flat class="bg-white text-dark shadow-1">
              <q-card-section class="q-py-md flex justify-between items-center">
                <div>
                  <div class="text-caption text-grey-7 text-weight-bold">SEKOLAH UTAMA (MITRA UTAMA)</div>
                  <div class="text-h6 text-weight-bold text-primary">{{ topSekolah.nama }}</div>
                  <div class="text-caption text-grey-6">{{ topSekolah.jumlah }} Calon Siswa Terdaftar Tahun Ini</div>
                </div>
                <q-avatar color="indigo-1" text-color="indigo-9" icon="workspace_premium" size="44px" />
              </q-card-section>
            </q-card>
          </div>
          <div class="col-12 q-mt-xs">
            <q-card flat class="bg-white text-dark shadow-1">
              <q-card-section class="q-py-md flex justify-between items-center">
                <div>
                  <div class="text-caption text-grey-7 text-weight-bold">RATA-RATA NILAI TERTINGGI</div>
                  <div class="text-h6 text-weight-bold text-green-8">{{ topNilaiSekolah.nama }}</div>
                  <div class="text-caption text-grey-6">Rerata Nilai Rapor/Asesmen SMP: <strong class="text-green-9">{{
                      topNilaiSekolah.nilai }}</strong></div>
                </div>
                <q-avatar color="green-1" text-color="green-8" icon="trending_up" size="44px" />
              </q-card-section>
            </q-card>
          </div>
        </div>
      </div>
    </div>

    <!-- Tabel Analisis Komparasi Seluruh Sekolah Asal -->
    <q-card class="shadow-1">
      <q-card-section class="q-py-sm flex justify-between items-center">
        <div class="text-subtitle2 text-weight-bold text-grey-9">Pemetaan Kontribusi Sekolah Asal</div>
        <q-input v-model="search" dense outlined placeholder="Cari nama SMP / MTs..." style="max-width: 230px;">
          <template v-slot:prepend><q-icon name="search" /></template>
        </q-input>
      </q-card-section>

      <q-card-section class="q-pa-none">
        <q-table flat :rows="filteredAnalisis" :columns="columns" row-key="nama" :pagination="{ rowsPerPage: 5 }">
          <!-- Kustomisasi Kolom Rata-rata Nilai -->
          <template v-slot:body-cell-rerataNilai="props">
            <q-td :props="props" class="text-center text-weight-bold text-grey-9">
              {{ props.row.rerataNilai }}
            </q-td>
          </template>

          <!-- Kustomisasi Kolom Progres Bar Distribusi Kontribusi -->
          <template v-slot:body-cell-progress="props">
            <q-td :props="props">
              <div class="flex items-center no-wrap">
                <q-linear-progress stripe rounded size="10px" :value="props.row.jumlah / 40" color="indigo-9"
                  class="q-mr-sm" style="min-width: 100px;" />
                <span class="text-caption text-weight-medium text-grey-8">{{ Math.round((props.row.jumlah / 25) * 100)
                  }}%</span>
              </div>
            </q-td>
          </template>
        </q-table>
      </q-card-section>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useQuasar } from 'quasar'
import apexchart from 'vue3-apexcharts'

const $q = useQuasar()
const search = ref('')

// Konfigurasi Bar Chart (ApexCharts)
const chartOptions = ref({
  chart: { id: 'asal-sekolah-chart', toolbar: { show: false } },
  colors: ['#303F9F'],
  plotOptions: {
    bar: { borderRadius: 4, horizontal: true, barHeight: '55%' }
  },
  dataLabels: { enabled: true, formatter: (val) => `${val} Siswa` },
  xaxis: {
    categories: ['SMPN 1 Soreang', 'SMPN 2 Bandung', 'MTs As-Syifa', 'SMP Pasundan 1', 'SMP Baiturrahman']
  }
})

const chartSeries = ref([
  { name: 'Jumlah Pendaftar', data: [18, 12, 9, 6, 4] }
])

// Data Ringkasan untuk Summary Cards
const topSekolah = ref({ nama: 'SMP Negeri 1 Soreang', jumlah: 18 })
const topNilaiSekolah = ref({ nama: 'SMP Negeri 2 Bandung', nilai: '87.50' })

// Kolom Struktur Tabel Analisis Komparasi
const columns = [
  { name: 'nama', align: 'left', label: 'Nama Instansi Sekolah Asal', field: 'nama', sortable: true },
  { name: 'jenis', align: 'center', label: 'Jenis', field: 'jenis', sortable: true },
  { name: 'jumlah', align: 'center', label: 'Jumlah Pendaftar', field: 'jumlah', sortable: true },
  { name: 'rerataNilai', align: 'center', label: 'Rerata Nilai Rapor', field: 'rerataNilai', sortable: true },
  { name: 'progress', align: 'left', label: 'Bobot Kontribusi (%)' }
]

// Master Data Analisis Pemasok PPDB
const dataAnalisisSekolah = ref([
  { nama: 'SMP Negeri 1 Soreang', jenis: 'Negeri', jumlah: 18, rerataNilai: 82.45 },
  { nama: 'SMP Negeri 2 Bandung', jenis: 'Negeri', jumlah: 12, rerataNilai: 87.50 },
  { nama: 'MTs As-Syifa', jenis: 'Swasta', jumlah: 9, rerataNilai: 80.12 },
  { nama: 'SMP Pasundan 1', jenis: 'Swasta', jumlah: 6, rerataNilai: 78.35 },
  { nama: 'SMP Baiturrahman', jenis: 'Swasta', jumlah: 4, rerataNilai: 84.20 }
])

// Fungsi Filter Pencarian Tabel
const filteredAnalisis = computed(() => {
  if (!search.value) return dataAnalisisSekolah.value
  return dataAnalisisSekolah.value.filter(item =>
    item.nama.toLowerCase().includes(search.value.toLowerCase())
  )
})

const shareAnalysis = () => {
  $q.notify({ color: 'indigo-9', message: 'Tautan data grafik analisis berhasil disalin untuk tim Humas/Promosi sekolah.', icon: 'share' })
}
</script>

<template>
  <q-page class="q-pa-md bg-grey-2">
    <!-- Header Halaman -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Tracer Study Alumni (BMW)</h1>
        <p class="text-caption text-grey-7 q-mb-none">Pelacakan rekam jejak keterserapan lulusan SMK: Bekerja,
          Melanjutkan Kuliah, atau Wirausaha</p>
      </div>
      <div>
        <q-btn color="primary" icon="analytics" label="Unduh Laporan Kemdikbud" @click="downloadLaporanKementerian" />
      </div>
    </div>

    <!-- Layout Atas: Grafik Visualisasi & Statistik Ringkas -->
    <div class="row q-col-gutter-md q-mb-md">
      <!-- Donut Chart Persentase Keterserapan -->
      <div class="col-12 col-md-5">
        <q-card class="shadow-1 fit">
          <q-card-section>
            <div class="text-subtitle2 text-weight-bold text-grey-9">Persentase Keterserapan Alumni</div>
            <div class="text-caption text-grey-6">Grafik real-time berdasarkan data angkatan terpilih</div>
          </q-card-section>
          <q-card-section class="flex flex-center">
            <apexchart type="donut" width="340" :options="chartOptions" :series="chartSeries"></apexchart>
          </q-card-section>
        </q-card>
      </div>

      <!-- Widget Detail Angka Target BMW -->
      <div class="col-12 col-md-7">
        <div class="row q-col-gutter-sm fit">
          <div class="col-6">
            <q-card class="bg-white text-dark flat shadow-1 q-mb-sm">
              <q-card-section class="q-py-sm flex justify-between items-center">
                <div>
                  <div class="text-caption text-grey-7">B — Bekerja (Industri)</div>
                  <div class="text-h6 text-weight-bold text-blue-9">145 Alumni</div>
                </div>
                <q-avatar color="blue-1" text-color="primary" icon="business" size="36px" />
              </q-card-section>
            </q-card>
          </div>
          <div class="col-6">
            <q-card class="bg-white text-dark flat shadow-1 q-mb-sm">
              <q-card-section class="q-py-sm flex justify-between items-center">
                <div>
                  <div class="text-caption text-grey-7">M — Melanjutkan (Kuliah)</div>
                  <div class="text-h6 text-weight-bold text-green-9">52 Alumni</div>
                </div>
                <q-avatar color="green-1" text-color="green-9" icon="school" size="36px" />
              </q-card-section>
            </q-card>
          </div>
          <div class="col-6">
            <q-card class="bg-white text-dark flat shadow-1">
              <q-card-section class="q-py-sm flex justify-between items-center">
                <div>
                  <div class="text-caption text-grey-7">W — Wirausaha (Mandiri)</div>
                  <div class="text-h6 text-weight-bold text-purple-9">24 Alumni</div>
                </div>
                <q-avatar color="purple-1" text-color="purple-9" icon="storefront" size="36px" />
              </q-card-section>
            </q-card>
          </div>
          <div class="col-6">
            <q-card class="bg-white text-dark flat shadow-1">
              <q-card-section class="q-py-sm flex justify-between items-center">
                <div>
                  <div class="text-caption text-grey-7">Belum Terserap / Mencari</div>
                  <div class="text-h6 text-weight-bold text-red-9">39 Alumni</div>
                </div>
                <q-avatar color="red-1" text-color="red-9" icon="person_search" size="36px" />
              </q-card-section>
            </q-card>
          </div>
          <div class="col-12 q-mt-sm">
            <q-banner class="bg-amber-1 text-amber-10 rounded-borders dense text-caption">
              <q-icon name="info" /> <strong>Indikator Kinerja Utama (IKU):</strong> Target keterserapan minimal
              angkatan lulusan tahun ini adalah sebesar 85% dalam waktu 6 bulan setelah kelulusan resmi.
            </q-banner>
          </div>
        </div>
      </div>
    </div>

    <!-- Panel Filter Pencarian & Tabel Alumni -->
    <q-card class="shadow-1">
      <q-card-section class="q-py-sm">
        <div class="row q-col-gutter-sm items-center">
          <div class="col-12 col-md-4">
            <q-input v-model="search" dense outlined placeholder="Cari Nama Alumni, Instansi, atau NISN..." clearable>
              <template v-slot:prepend><q-icon name="search" /></template>
            </q-input>
          </div>
          <div class="col-12 col-sm-6 col-md-3">
            <q-select v-model="filterStatus"
              :options="['Semua Status', 'Bekerja', 'Melanjutkan Kuliah', 'Wirausaha', 'Belum Bekerja']" dense outlined
              label="Status Keterserapan" />
          </div>
          <div class="col-12 col-sm-6 col-md-3">
            <q-select v-model="filterAngkatan" :options="['Angkatan 2025', 'Angkatan 2024', 'Angkatan 2023']" dense
              outlined label="Tahun Lulus" />
          </div>
          <div class="col-12 col-md-2 flex justify-end">
            <q-btn flat color="grey-7" icon="refresh" label="Reset" @click="resetFilters" />
          </div>
        </div>
      </q-card-section>

      <!-- Tabel Data Rekam Jejak Tracer -->
      <q-card-section class="q-pa-none">
        <q-table flat :rows="filteredTracer" :columns="columns" row-key="id" :pagination="{ rowsPerPage: 10 }">
          <!-- Custom Tampilan Kolom Status BMW -->
          <template v-slot:body-cell-status="props">
            <q-td :props="props" class="text-center">
              <q-badge :color="props.row.status === 'Bekerja' ? 'blue-8' :
                props.row.status === 'Melanjutkan Kuliah' ? 'green-8' :
                  props.row.status === 'Wirausaha' ? 'purple-8' : 'red-8'
                " :label="props.row.status" />
            </q-td>
          </template>

          <!-- Custom Kolom Tampilan Nama Instansi / Tempat Penempatan -->
          <template v-slot:body-cell-instansi="props">
            <q-td :props="props">
              <div class="text-weight-bold text-grey-9">{{ props.row.instansi || '—' }}</div>
              <div class="text-caption text-grey-6">{{ props.row.posisiUsaha || 'Mencari Peluang' }}</div>
            </q-td>
          </template>

          <!-- Custom Kolom Aksi Update Status Tracer -->
          <template v-slot:body-cell-aksi="props">
            <q-td :props="props" class="text-center">
              <q-btn size="sm" color="primary" icon="edit_note" round flat @click="openUpdateStatusDialog(props.row)">
                <q-tooltip>Update Status BMW & Detail Karir</q-tooltip>
              </q-btn>
            </q-td>
          </template>
        </q-table>
      </q-card-section>
    </q-card>

    <!-- DIALOG POPUP: Update Status Penyerapan Alumni -->
    <q-dialog v-model="dialogTracerOpen" persistent>
      <q-card style="min-width: 400px;">
        <q-card-section class="bg-primary text-white">
          <div class="text-h6">Pembaruan Tracer Study</div>
          <div class="text-caption text-blue-2">Alumni: {{ activeAlumni.nama }} ({{ activeAlumni.angkatan }})</div>
        </q-card-section>

        <q-card-section class="q-gutter-md q-pt-md">
          <q-select v-model="formTracer.status"
            :options="['Bekerja', 'Melanjutkan Kuliah', 'Wirausaha', 'Belum Bekerja']"
            label="Status Penyerapan Saat Ini" outlined dense />
          <q-input v-model="formTracer.instansi" label="Nama Tempat Industri / Universitas / Nama Toko" outlined dense
            placeholder="Contoh: PT. Yamaha Motor / Universitas Padjadjaran..." />
          <q-input v-model="formTracer.posisiUsaha" label="Jabatan Kerja / Program Studi / Bidang Bisnis" outlined dense
            placeholder="Contoh: Operator Produksi / S1 Teknik Informatika..." />
        </q-card-section>

        <q-card-actions align="right" class="q-pb-md q-pr-md">
          <q-btn flat label="Batal" v-close-popup color="grey" />
          <q-btn color="primary" label="Simpan Rekam Jejak" @click="saveTracerInfo" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useQuasar } from 'quasar'
import apexchart from 'vue3-apexcharts'

const $q = useQuasar()

// State Komponen
const search = ref('')
const filterStatus = ref('Semua Status')
const filterAngkatan = ref('Angkatan 2025')
const dialogTracerOpen = ref(false)

const activeAlumni = ref({})
const formTracer = ref({ status: '', instansi: '', posisiUsaha: '' })

// Konfigurasi Donut Chart (ApexCharts)
const chartOptions = ref({
  chart: { id: 'donut-tracer-chart' },
  labels: ['Bekerja', 'Melanjutkan Kuliah', 'Wirausaha', 'Belum Terserap'],
  colors: ['#1976D2', '#2E7D32', '#6A1B9A', '#C62828'],
  legend: { position: 'bottom' },
  dataLabels: { enabled: true }
})

// Data Array untuk Donut Chart (Sesuai dengan jumlah widget kanan)
const chartSeries = ref([145, 52, 24, 39])

// Struktur Kolom Tabel q-table Log Tracer Study
const columns = [
  { name: 'nama', align: 'left', label: 'Nama Lengkap Alumni', field: 'nama', sortable: true },
  { name: 'jurusan', align: 'left', label: 'Jurusan', field: 'jurusan' },
  { name: 'status', align: 'center', label: 'Status BMW', field: 'status', sortable: true },
  { name: 'instansi', align: 'left', label: 'Instansi Tempat Berada', field: 'instansi', sortable: true },
  { name: 'aksi', align: 'center', label: 'Aksi' }
]

// Data Mock / Dummy Tracer Lulusan SMK
const dataTracer = ref([
  { id: 1, nama: 'Rian Hidayat', angkatan: '2025', jurusan: 'RPL', status: 'Bekerja', instansi: 'PT. Telekomunikasi Indonesia Tbk', posisiUsaha: 'Frontend Developer' },
  { id: 2, nama: 'Siti Sarah', angkatan: '2025', jurusan: 'AKL', status: 'Melanjutkan Kuliah', instansi: 'Universitas Indonesia', posisiUsaha: 'S1 Akuntansi Perajakan' }, { id: 3, nama: 'Denny Wahyudi', angkatan: '2025', jurusan: 'TKRO', status: 'Wirausaha', instansi: 'Bengkel Wahyudi Motor', posisiUsaha: 'Pemilik Bengkel Mandiri' }, { id: 4, nama: 'Gita Permata', angkatan: '2025', jurusan: 'RPL', status: 'Belum Bekerja', instansi: '', posisiUsaha: '' }])
// Logika Compute Filter Pencarian Alumni di Frontend Sisi Klien
const filteredTracer = computed(() => {
  return dataTracer.value.filter(item => {
    const matchSearch = search.value ? item.nama.toLowerCase().includes(search.value.toLowerCase()) || item.instansi.toLowerCase().includes(search.value.toLowerCase()) : true
    const matchStatus = filterStatus.value !== 'Semua Status' ? item.status === filterStatus.value : true
    const matchAngkatan = filterAngkatan.value ? item.angkatan === filterAngkatan.value.split(' ') // Mengambil angka tahun saja dari string "Angkatan 2025"
      : true
    return matchSearch && matchStatus && matchAngkatan
  })
})
// Fungsi Penanganan Event Aksi Kesiswaan & BKK
const resetFilters = () => {
  search.value = ''
  filterStatus.value = 'Semua Status'
  filterAngkatan.value = 'Angkatan 2025'
}
const openUpdateStatusDialog = (row) => {
  activeAlumni.value = row
  formTracer.value = { status: row.status, instansi: row.instansi, posisiUsaha: row.posisiUsaha }
  dialogTracerOpen.value = true
}
const saveTracerInfo = () => {
  const target = dataTracer.value.find(item => item.id === activeAlumni.value.id)
  if (target) {
    target.status = formTracer.value.status
    target.instansi = formTracer.value.instansi
    target.posisiUsaha = formTracer.value.posisiUsaha
  }
  $q.notify({ color: 'primary', message: `Data penelusuran karir untuk ${activeAlumni.value.nama} berhasil diperbarui!`, icon: 'save' })
  dialogTracerOpen.value = false
}
const downloadLaporanKementerian = () => { $q.notify({ color: 'green-8', message: 'Sukses mengekspor sinkronisasi draf file .xml siap upload ke web Tracer Study Direktorat SMK Kemdikbud.', icon: 'cloud_download' }) }


</script>

<template>
  <q-page class="q-pa-md bg-grey-2">
    <!-- Header Halaman -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Rekapitulasi & Pelaporan Berkala</h1>
        <p class="text-caption text-grey-7 q-mb-none">Penyusunan statistik bulanan siswa aktif, mutasi, dan monitoring
          berkas wali kelas untuk Dinas Pendidikan</p>
      </div>
      <div>
        <q-btn-dropdown color="primary" icon="download" label="Unduh Laporan Resmi">
          <q-list>
            <q-item clickable v-close-popup @click="downloadLaporan('Bulanan Dinas')">
              <q-item-section avatar><q-icon name="picture_as_pdf" color="red" /></q-item-section>
              <q-item-section><q-item-label>Laporan Bulanan Dinas (.pdf)</q-item-label></q-item-section>
            </q-item>
            <q-item clickable v-close-popup @click="downloadLaporan('Statistik Kemdikbud')">
              <q-item-section avatar><q-icon name="table_view" color="green-8" /></q-item-section>
              <q-item-section><q-item-label>Rekap Statistik Siswa (.xlsx)</q-item-label></q-item-section>
            </q-item>
          </q-list>
        </q-btn-dropdown>
      </div>
    </div>

    <!-- Navigasi Tab Modul Pelaporan -->
    <q-card class="shadow-1 q-mb-md">
      <q-tabs v-model="activeReportTab" dense class="text-grey" active-color="primary" indicator-color="primary"
        align="left" narrow-indicator>
        <q-tab name="bulanan" icon="analytics" label="Laporan Bulanan Siswa (Statistik)" />
        <q-tab name="monitoring" icon="fact_check" label="Monitoring Wali Kelas (Absen & Nilai)" />
      </q-tabs>
    </q-card>

    <!-- Konten Tab Pelaporan -->
    <q-tab-panels v-model="activeReportTab" animated class="bg-transparent">

      <!-- TAB 1: LAPORAN BULANAN / STATISTIK JURUSAN -->
      <q-tab-panel name="bulanan" class="q-pa-none">
        <q-card class="shadow-1">
          <q-card-section class="flex justify-between items-center q-py-sm">
            <div class="text-subtitle2 text-weight-bold text-grey-9">Rekapitulasi Keadaan Siswa Bulan: {{
              currentMonthYear }}</div>
            <q-select v-model="selectedPeriodeLaporan" :options="['Juni 2026', 'Mei 2026', 'April 2026']" dense outlined
              style="width: 150px" label="Periode" />
          </q-card-section>

          <q-card-section class="q-pa-none">
            <q-table flat :rows="dataStatistikBulanan" :columns="columnsStatistik" row-key="jurusan"
              :pagination="{ rowsPerPage: 10 }" hide-bottom>
              <!-- Kustomisasi Tampilan Baris Jurusan -->
              <template v-slot:body-cell-jurusan="props">
                <q-td :props="props" class="text-weight-bold text-primary">
                  {{ props.row.jurusan }}
                </q-td>
              </template>
            </q-table>
          </q-card-section>

          <q-card-section class="bg-blue-1 text-primary text-caption">
            <q-icon name="info" /> Rumus Baku Keadaan: <strong>Siswa Akhir = (Siswa Awal + Mutasi Masuk) - Mutasi
              Keluar</strong>.
            Angka ini wajib sinkron dengan total rombel Dapodik lokal.
          </q-card-section>
        </q-card>
      </q-tab-panel>

      <!-- TAB 2: MONITORING WALI KELAS (SETORAN ABSEN & LEGER) -->
      <q-tab-panel name="monitoring" class="q-pa-none">
        <div class="row q-col-gutter-md">
          <!-- Tabel Pemantauan -->
          <div class="col-12 col-md-8">
            <q-card class="shadow-1">
              <q-card-section class="q-py-sm">
                <div class="text-subtitle2 text-weight-bold text-grey-9">Status Pengumpulan Berkas Kelompok Wali Kelas
                </div>
              </q-card-section>
              <q-card-section class="q-pa-none">
                <q-table flat :rows="dataMonitoringWali" :columns="columnsMonitoring" row-key="kelas">
                  <!-- Status Setor Rekap Absensi -->
                  <template v-slot:body-cell-absen="props">
                    <q-td :props="props" class="text-center">
                      <q-icon :name="props.row.absenDisetor ? 'check_circle' : 'cancel'"
                        :color="props.row.absenDisetor ? 'green' : 'red'" size="sm" />
                    </q-td>
                  </template>

                  <!-- Status Setor Leger Nilai -->
                  <template v-slot:body-cell-leger="props">
                    <q-td :props="props" class="text-center">
                      <q-icon :name="props.row.legerDisetor ? 'check_circle' : 'cancel'"
                        :color="props.row.legerDisetor ? 'green' : 'red'" size="sm" />
                    </q-td>
                  </template>

                  <!-- Kolom Tindakan / Kejar Wali Kelas -->
                  <template v-slot:body-cell-aksi="props">
                    <q-td :props="props" class="text-center">
                      <q-btn size="xs" color="orange-9" icon="chat" label="Kirim Reminder"
                        :disable="props.row.absenDisetor && props.row.legerDisetor"
                        @click="sendWaliReminder(props.row)" />
                    </q-td>
                  </template>
                </q-table>
              </q-card-section>
            </q-card>
          </div>

          <!-- Panel Kanan: Log Kelengkapan Berkas Administrasi Ujian -->
          <div class="col-12 col-md-4">
            <q-card class="shadow-1">
              <q-card-section class="bg-teal-7 text-white">
                <div class="text-subtitle2 text-weight-bold">Status Arsip Semesteran</div>
                <div class="text-caption text-teal-1">Kesiapan pelaporan akhir tahun ajaran berjalan</div>
              </q-card-section>
              <q-card-section>
                <q-list dense separator>
                  <q-item class="q-py-sm">
                    <q-item-section><q-item-label>Buku Induk Digital Kelas XII</q-item-label></q-item-section>
                    <q-item-section side><q-badge color="green" label="100% Terisi" /></q-item-section>
                  </q-item>
                  <q-item class="q-py-sm">
                    <q-item-section><q-item-label>Kompilasi Sertifikat PKL Alumni</q-item-label></q-item-section>
                    <q-item-section side><q-badge color="green" label="Lengkap" /></q-item-section>
                  </q-item>
                  <q-item class="q-py-sm">
                    <q-item-section><q-item-label>Verifikasi Buku Klaper Abjad A-Z</q-item-label></q-item-section>
                    <q-item-section side><q-badge color="orange" label="Peninjauan" /></q-item-section>
                  </q-item>
                </q-list>
              </q-card-section>
            </q-card>
          </div>
        </div>
      </q-tab-panel>
    </q-tab-panels>
  </q-page>
</template>

<script setup>
import { ref } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const activeReportTab = ref('bulanan')
const selectedPeriodeLaporan = ref('Juni 2026')
const currentMonthYear = ref('Juni 2026')

// 1. STRUKTUR KOLOM TAB 1: LAPORAN BULANAN (STATISTIK DATA SISWA)
const columnsStatistik = [
  { name: 'jurusan', align: 'left', label: 'Kompetensi Keahlian (SMK)', field: 'jurusan' },
  { name: 'awal', align: 'center', label: 'Siswa Awal Bulan', field: 'awal', sortable: true },
  { name: 'masuk', align: 'center', label: 'Mutasi Masuk (+)', field: 'masuk' },
  { name: 'keluar', align: 'center', label: 'Mutasi Keluar (-)', field: 'keluar' },
  { name: 'akhir', align: 'center', label: 'Siswa Akhir Bulan', field: 'akhir', sortable: true }
]

// Data Mock / Dummy Rekap Negara Statistik Bulanan Per Jurusan SMK
const dataStatistikBulanan = ref([
  { jurusan: 'Rekayasa Perangkat Lunak (RPL)', awal: 345, masuk: 2, keluar: 0, akhir: 347 },
  { jurusan: 'Teknik Kendaraan Ringan Otomotif (TKRO)', awal: 412, masuk: 0, keluar: 1, akhir: 411 },
  { jurusan: 'Akuntansi & Keuangan Lembaga (AKL)', awal: 298, masuk: 1, keluar: 0, akhir: 299 }
])

// 2. STRUKTUR KOLOM TAB 2: MONITORING SETORAN DOKUMEN WALI KELAS
const columnsMonitoring = [
  { name: 'kelas', align: 'left', label: 'Rombongan Belajar (Kelas)', field: 'kelas', sortable: true },
  { name: 'wali', align: 'left', label: 'Nama Wali Kelas', field: 'wali' },
  { name: 'absen', align: 'center', label: 'Rekap Absensi Bulanan', field: 'absenDisetor' },
  { name: 'ledger', align: 'center', label: 'Leger Nilai Rapor', field: 'legerDisetor' },
  { name: 'aksi', align: 'center', label: 'Tindakan' }
]

// Data Mock / Dummy Monitoring Wali Kelas
const dataMonitoringWali = ref([
  { kelas: 'XII RPL 1', wali: 'Drs. H. Mulyana', absenDisetor: true, legerDisetor: true },
  { kelas: 'XII RPL 2', wali: 'Rini Andrini, S.Pd.', absenDisetor: true, legerDisetor: false },
  { kelas: 'XII TKRO 3', wali: 'Agus Setiawan, S.T.', absenDisetor: false, legerDisetor: false },
  { kelas: 'XI AKL 1', wali: 'Sri Wahyuni, S.E.', absenDisetor: true, ledgerDisetor: true }
])

// Fungsi Penanganan Tombol Aksi TU Kesiswaan
const downloadLaporan = (jenisDokumen) => {
  $q.notify({
    color: 'primary',
    message: `Menyusun generator arsip... Mengunduh berkas [${jenisDokumen}] periode ${selectedPeriodeLaporan.value}`,
    icon: 'download'
  })
}

const sendWaliReminder = (row) => {
  $q.notify({
    color: 'orange-9',
    message: `Pemberitahuan peringatan otomatis berhasil dikirim ke nomor WhatsApp Guru Wali Kelas: ${row.wali}(${row.kelas})`, icon: 'send'
  })
}


</script>

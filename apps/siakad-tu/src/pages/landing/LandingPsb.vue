<template>
  <q-page class="dashboard-page bg-grey-1">
    <!-- HEADER -->
    <div class="dashboard-header q-px-md q-py-lg">
      <div class="container max-width-center row items-center justify-between">
        <div>
          <div class="text-overline text-secondary text-weight-medium letter-spacing-2">
            SISTEM INFORMASI AKADEMIK
          </div>
          <div class="text-h5 text-weight-bold text-white q-mt-xs">
            Dashboard Penerimaan Siswa Baru
          </div>
          <div class="text-caption text-grey-4 q-mt-xs">
            {{ todayLabel }} • Tahun Ajaran 2026/2027
          </div>
        </div>
        <div class="row items-center q-gutter-sm">
          <q-btn unelevated color="secondary" text-color="dark" label="Verifikasi Berkas" no-caps icon="fact_check" />
          <q-btn round flat color="white" icon="notifications">
            <q-badge color="secondary" text-color="dark" floating rounded>{{ notifCount }}</q-badge>
          </q-btn>
        </div>
      </div>
    </div>

    <div class="container max-width-center q-px-md q-py-lg">
      <!-- KPI CARDS -->
      <div class="row q-col-gutter-md q-mb-lg">
        <div class="col-6 col-md-3" v-for="kpi in kpis" :key="kpi.label">
          <q-card flat bordered class="kpi-card">
            <q-card-section class="row items-center no-wrap">
              <q-avatar :color="kpi.color" text-color="white" :icon="kpi.icon" size="46px" class="q-mr-md" />
              <div>
                <div class="text-h6 text-weight-bold text-primary">{{ kpi.value }}</div>
                <div class="text-caption text-grey-7">{{ kpi.label }}</div>
              </div>
            </q-card-section>
          </q-card>
        </div>
      </div>

      <!-- PROGRESS KUOTA KESELURUHAN -->
      <q-card flat bordered class="q-mb-lg content-card">
        <q-card-section>
          <div class="row items-center justify-between q-mb-sm">
            <div class="text-subtitle1 text-weight-bold text-primary">
              Progress Kuota Penerimaan Keseluruhan
            </div>
            <div class="text-caption text-grey-7">{{ siswaDiterima }} dari {{ kuotaTotal }} kuota</div>
          </div>
          <div class="bar-track bar-track-lg">
            <div class="bar-fill" :style="{ width: kuotaPct + '%', backgroundColor: '#0a192f' }" />
          </div>
          <div class="text-caption text-grey-6 q-mt-xs">{{ kuotaPct }}% kuota terisi</div>
        </q-card-section>
      </q-card>

      <div class="row q-col-gutter-lg">
        <!-- LEFT COLUMN -->
        <div class="col-12 col-md-8">
          <!-- KUOTA HAMPIR PENUH / KRITIS -->
          <q-card flat bordered class="q-mb-lg content-card border-alert" v-if="jurusanKritis.length">
            <q-card-section>
              <div class="row items-center">
                <q-icon name="report" color="red" size="22px" class="q-mr-sm" />
                <div class="text-subtitle1 text-weight-bold text-red-9">
                  Jurusan dengan Kuota Hampir Penuh
                </div>
              </div>
            </q-card-section>
            <q-separator />
            <q-list separator>
              <q-item v-for="jurusan in jurusanKritis" :key="jurusan.nama">
                <q-item-section avatar>
                  <q-avatar color="red-1" text-color="red-9" icon="school" size="38px" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-medium">{{ jurusan.nama }}</q-item-label>
                  <q-item-label caption>{{ jurusan.diterima }} / {{ jurusan.kuota }} kuota terisi</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-badge color="red" text-color="white" class="q-px-sm">
                    Sisa {{ jurusan.kuota - jurusan.diterima }}
                  </q-badge>
                </q-item-section>
              </q-item>
            </q-list>
          </q-card>

          <!-- TREN PENDAFTARAN HARIAN -->
          <q-card flat bordered class="q-mb-lg content-card">
            <q-card-section>
              <div class="row items-center justify-between">
                <div class="text-subtitle1 text-weight-bold text-primary">
                  Tren Pendaftaran Harian (2 Minggu Terakhir)
                </div>
                <q-badge color="grey-3" text-color="grey-8" class="q-px-sm">
                  {{ pendaftarHariIni }} pendaftar hari ini
                </q-badge>
              </div>
            </q-card-section>
            <q-separator />
            <q-card-section>
              <div class="trend-chart row items-end no-wrap">
                <div
                  v-for="hari in trenPendaftaran"
                  :key="hari.label"
                  class="trend-bar-wrap col text-center"
                >
                  <div class="text-caption text-grey-7 q-mb-xs">{{ hari.jumlah }}</div>
                  <div class="trend-bar" :style="{ height: trendHeight(hari.jumlah) + 'px' }" />
                  <div class="text-caption text-grey-6 q-mt-xs">{{ hari.label }}</div>
                </div>
              </div>
            </q-card-section>
          </q-card>

          <!-- PENDAFTAR PER JURUSAN -->
          <q-card flat bordered class="content-card">
            <q-card-section>
              <div class="text-subtitle1 text-weight-bold text-primary">
                Pendaftar per Kompetensi Keahlian Pilihan
              </div>
            </q-card-section>
            <q-separator />
            <q-card-section>
              <div v-for="jurusan in pendaftarJurusan" :key="jurusan.nama" class="q-mb-md">
                <div class="row items-center justify-between q-mb-xs">
                  <div class="text-body2 text-weight-medium">{{ jurusan.nama }}</div>
                  <div class="text-caption text-grey-7">
                    {{ jurusan.pendaftar }} pendaftar • rasio 1:{{ jurusan.rasio }}
                  </div>
                </div>
                <div class="bar-track">
                  <div
                    class="bar-fill"
                    :style="{ width: barWidth(jurusan.pendaftar) + '%', backgroundColor: jurusan.warna }"
                  />
                </div>
              </div>
            </q-card-section>
          </q-card>
        </div>

        <!-- RIGHT COLUMN -->
        <div class="col-12 col-md-4">
          <!-- JALUR PENDAFTARAN -->
          <q-card flat bordered class="q-mb-lg content-card">
            <q-card-section>
              <div class="text-subtitle1 text-weight-bold text-primary">
                Pendaftar per Jalur
              </div>
            </q-card-section>
            <q-separator />
            <q-list separator>
              <q-item v-for="jalur in jalurPendaftaran" :key="jalur.nama">
                <q-item-section avatar>
                  <q-avatar :color="jalur.warna" text-color="white" :icon="jalur.icon" size="38px" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-medium">{{ jalur.nama }}</q-item-label>
                  <q-item-label caption>Kuota {{ jalur.kuota }}</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <div class="text-weight-bold text-primary">{{ jalur.jumlah }}</div>
                </q-item-section>
              </q-item>
            </q-list>
          </q-card>

          <!-- STATUS VERIFIKASI BERKAS -->
          <q-card flat bordered class="q-mb-lg content-card">
            <q-card-section>
              <div class="text-subtitle1 text-weight-bold text-primary">
                Status Verifikasi Berkas
              </div>
            </q-card-section>
            <q-separator />
            <q-list separator>
              <q-item v-for="berkas in statusVerifikasi" :key="berkas.status">
                <q-item-section>
                  <q-item-label class="text-weight-medium">{{ berkas.status }}</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-badge :color="berkas.warna" text-color="white">
                    {{ berkas.jumlah }}
                  </q-badge>
                </q-item-section>
              </q-item>
            </q-list>
          </q-card>

          <!-- TAHAPAN PPDB -->
          <q-card flat bordered class="content-card">
            <q-card-section>
              <div class="text-subtitle1 text-weight-bold text-primary">
                Jadwal Tahapan PPDB
              </div>
            </q-card-section>
            <q-separator />
            <q-list separator>
              <q-item v-for="tahap in tahapanPpdb" :key="tahap.judul">
                <q-item-section avatar>
                  <div class="agenda-date text-center">
                    <div class="text-weight-bold text-primary">{{ tahap.tanggal }}</div>
                    <div class="text-caption text-grey-7">{{ tahap.bulan }}</div>
                  </div>
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-medium">{{ tahap.judul }}</q-item-label>
                  <q-item-label caption>{{ tahap.status }}</q-item-label>
                </q-item-section>
                <q-item-section side v-if="tahap.selesai">
                  <q-icon name="check_circle" color="teal" size="20px" />
                </q-item-section>
              </q-item>
            </q-list>
          </q-card>
        </div>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, computed } from 'vue'

const todayLabel = computed(() => {
  const today = new Date()
  return today.toLocaleDateString('id-ID', {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  })
})

const notifCount = ref(7)

const kpis = ref([
  { label: 'Total Pendaftar', value: '1.586', icon: 'how_to_reg', color: 'primary' },
  { label: 'Kuota Tersedia', value: '432', icon: 'event_seat', color: 'teal' },
  { label: 'Verifikasi Berkas Pending', value: '58', icon: 'pending_actions', color: 'orange-8' },
  { label: 'Siswa Diterima', value: '374', icon: 'verified', color: 'secondary' },
])

const kuotaTotal = 432
const siswaDiterima = 374
const kuotaPct = Math.round((siswaDiterima / kuotaTotal) * 100)

const jurusanKritis = ref([
  { nama: 'Teknik Komputer & Jaringan', diterima: 108, kuota: 112 },
  { nama: 'Teknik & Bisnis Sepeda Motor (TBSM)', diterima: 128, kuota: 140 },
])

const pendaftarHariIni = 46

const trenPendaftaran = ref([
  { label: '28 Jun', jumlah: 62 },
  { label: '30 Jun', jumlah: 74 },
  { label: '02 Jul', jumlah: 58 },
  { label: '04 Jul', jumlah: 90 },
  { label: '06 Jul', jumlah: 71 },
  { label: '08 Jul', jumlah: 55 },
  { label: '10 Jul', jumlah: 46 },
])

const maxTren = Math.max(...trenPendaftaran.value.map((h) => h.jumlah))
const trendHeight = (jumlah) => Math.max(8, Math.round((jumlah / maxTren) * 120))

const pendaftarJurusan = ref([
  { nama: 'Teknik & Bisnis Sepeda Motor (TBSM)', pendaftar: 520, rasio: 3.7, warna: '#0a192f' },
  { nama: 'Teknik Komputer & Jaringan', pendaftar: 468, rasio: 4.2, warna: '#c9a227' },
  { nama: 'Akuntansi & Keuangan Lembaga', pendaftar: 312, rasio: 2.6, warna: '#00796b' },
  { nama: 'Otomatisasi Tata Kelola Perkantoran', pendaftar: 168, rasio: 1.8, warna: '#5e35b1' },
  { nama: 'Multimedia', pendaftar: 118, rasio: 2.1, warna: '#d84315' },
])

const maxPendaftar = Math.max(...pendaftarJurusan.value.map((j) => j.pendaftar))
const barWidth = (jumlah) => Math.round((jumlah / maxPendaftar) * 100)

const jalurPendaftaran = ref([
  { nama: 'Jalur Reguler', kuota: 300, jumlah: 1024, icon: 'assignment', warna: 'primary' },
  { nama: 'Jalur Prestasi', kuota: 80, jumlah: 342, icon: 'emoji_events', warna: 'amber-8' },
  { nama: 'Jalur Afirmasi (KIP/Disabilitas)', kuota: 52, jumlah: 220, icon: 'volunteer_activism', warna: 'teal' },
])

const statusVerifikasi = ref([
  { status: 'Terverifikasi', jumlah: 1412, warna: 'teal' },
  { status: 'Menunggu Verifikasi', jumlah: 58, warna: 'orange' },
  { status: 'Berkas Ditolak / Perlu Revisi', jumlah: 24, warna: 'red' },
])

const tahapanPpdb = ref([
  { tanggal: '01', bulan: 'Jun', judul: 'Pendaftaran & Unggah Berkas', status: 'Selesai', selesai: true },
  { tanggal: '25', bulan: 'Jun', judul: 'Verifikasi Berkas & Tes Minat Bakat', status: 'Selesai', selesai: true },
  { tanggal: '10', bulan: 'Jul', judul: 'Pengumuman Hasil Seleksi Tahap 1', status: 'Berlangsung', selesai: false },
  { tanggal: '18', bulan: 'Jul', judul: 'Daftar Ulang Siswa Diterima', status: 'Akan Datang', selesai: false },
])
</script>

<style scoped>
.max-width-center {
  max-width: 1200px;
  margin: 0 auto;
}

.letter-spacing-2 {
  letter-spacing: 0.12em;
}

.dashboard-header {
  background: linear-gradient(135deg, #0a192f 0%, #123a63 100%);
}

.kpi-card {
  border-radius: 8px;
  height: 100%;
}

.content-card {
  border-radius: 8px;
}

.border-alert {
  border-left: 4px solid #c62828;
}

.bar-track {
  width: 100%;
  height: 10px;
  border-radius: 6px;
  background-color: #eceff1;
  overflow: hidden;
}

.bar-track-lg {
  height: 14px;
}

.bar-fill {
  height: 100%;
  border-radius: 6px;
  transition: width 0.4s ease;
}

.trend-chart {
  height: 170px;
  gap: 10px;
}

.trend-bar-wrap {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: flex-end;
}

.trend-bar {
  width: 100%;
  max-width: 28px;
  background: linear-gradient(180deg, #123a63 0%, #0a192f 100%);
  border-radius: 4px 4px 0 0;
  transition: height 0.4s ease;
}

.agenda-date {
  min-width: 44px;
}
</style>

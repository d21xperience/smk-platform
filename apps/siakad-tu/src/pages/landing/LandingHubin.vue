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
            Dashboard Hubungan Industri
          </div>
          <div class="text-caption text-grey-4 q-mt-xs">
            {{ todayLabel }}
          </div>
        </div>
        <div class="row items-center q-gutter-sm">
          <q-btn unelevated color="secondary" text-color="dark" label="Tambah Mitra" no-caps icon="add_business" />
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

      <div class="row q-col-gutter-lg">
        <!-- LEFT COLUMN -->
        <div class="col-12 col-md-8">
          <!-- MOU AKAN BERAKHIR -->
          <q-card flat bordered class="q-mb-lg content-card border-alert" v-if="mouPerhatian.length">
            <q-card-section>
              <div class="row items-center">
                <q-icon name="report" color="red" size="22px" class="q-mr-sm" />
                <div class="text-subtitle1 text-weight-bold text-red-9">
                  MoU Akan Berakhir / Perlu Tindak Lanjut
                </div>
              </div>
            </q-card-section>
            <q-separator />
            <q-list separator>
              <q-item v-for="mou in mouPerhatian" :key="mou.mitra">
                <q-item-section avatar>
                  <q-avatar color="red-1" text-color="red-9" icon="handshake" size="38px" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-medium">{{ mou.mitra }}</q-item-label>
                  <q-item-label caption>{{ mou.ruangLingkup }}</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-badge color="red" text-color="white" class="q-px-sm">
                    Berakhir {{ mou.tanggalBerakhir }}
                  </q-badge>
                </q-item-section>
              </q-item>
            </q-list>
          </q-card>

          <!-- DISTRIBUSI SISWA PKL -->
          <q-card flat bordered class="q-mb-lg content-card">
            <q-card-section>
              <div class="row items-center justify-between">
                <div class="text-subtitle1 text-weight-bold text-primary">
                  Distribusi Siswa PKL per Mitra Industri
                </div>
                <q-badge color="grey-3" text-color="grey-8" class="q-px-sm">
                  {{ totalSiswaPkl }} siswa magang
                </q-badge>
              </div>
            </q-card-section>
            <q-separator />
            <q-card-section>
              <div v-for="mitra in distribusiPkl" :key="mitra.nama" class="q-mb-md">
                <div class="row items-center justify-between q-mb-xs">
                  <div class="text-body2 text-weight-medium">{{ mitra.nama }}</div>
                  <div class="text-caption text-grey-7">{{ mitra.jumlah }} siswa</div>
                </div>
                <div class="bar-track">
                  <div class="bar-fill"
                    :style="{ width: barWidth(mitra.jumlah) + '%', backgroundColor: mitra.warna }" />
                </div>
              </div>
            </q-card-section>
          </q-card>

          <!-- MONITORING PKL -->
          <q-card flat bordered class="content-card">
            <q-card-section>
              <div class="row items-center justify-between">
                <div class="text-subtitle1 text-weight-bold text-primary">
                  Monitoring Pelaksanaan PKL
                </div>
                <q-btn flat dense color="primary" label="Lihat Semua" no-caps icon-right="arrow_forward" />
              </div>
            </q-card-section>
            <q-separator />
            <q-list separator>
              <q-item v-for="pkl in monitoringPkl" :key="pkl.id">
                <q-item-section avatar>
                  <q-avatar :color="statusWarna(pkl.status)" text-color="white" icon="work" size="38px" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-medium">{{ pkl.siswa }} • {{ pkl.kelas }}</q-item-label>
                  <q-item-label caption>{{ pkl.mitra }} • Pembimbing: {{ pkl.pembimbing }}</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-chip dense :color="statusWarna(pkl.status)" text-color="white" class="text-caption">
                    {{ pkl.status }}
                  </q-chip>
                </q-item-section>
              </q-item>
            </q-list>
          </q-card>
        </div>

        <!-- RIGHT COLUMN -->
        <div class="col-12 col-md-4">
          <!-- MITRA INDUSTRI BARU -->
          <q-card flat bordered class="q-mb-lg content-card">
            <q-card-section>
              <div class="text-subtitle1 text-weight-bold text-primary">
                Mitra Industri Bergabung Terbaru
              </div>
            </q-card-section>
            <q-separator />
            <q-list separator>
              <q-item v-for="mitra in mitraBaru" :key="mitra.nama">
                <q-item-section avatar>
                  <q-avatar color="primary" text-color="white" icon="business" size="38px" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-medium">{{ mitra.nama }}</q-item-label>
                  <q-item-label caption>{{ mitra.bidang }} • Bergabung {{ mitra.tanggal }}</q-item-label>
                </q-item-section>
              </q-item>
            </q-list>
          </q-card>

          <!-- LOWONGAN BKK -->
          <q-card flat bordered class="q-mb-lg content-card">
            <q-card-section>
              <div class="text-subtitle1 text-weight-bold text-primary">
                Lowongan Aktif (BKK)
              </div>
            </q-card-section>
            <q-separator />
            <q-list separator>
              <q-item v-for="lowongan in lowonganBkk" :key="lowongan.posisi">
                <q-item-section>
                  <q-item-label class="text-weight-medium">{{ lowongan.posisi }}</q-item-label>
                  <q-item-label caption>{{ lowongan.perusahaan }} • Tutup {{ lowongan.batasAkhir }}</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-badge color="teal" text-color="white">
                    {{ lowongan.pelamar }} pelamar
                  </q-badge>
                </q-item-section>
              </q-item>
            </q-list>
          </q-card>

          <!-- ALUMNI TERSERAP -->
          <q-card flat bordered class="q-mb-lg content-card">
            <q-card-section>
              <div class="text-subtitle1 text-weight-bold text-primary">
                Keterserapan Alumni per Jurusan
              </div>
            </q-card-section>
            <q-separator />
            <q-list separator>
              <q-item v-for="jurusan in keterserapanAlumni" :key="jurusan.nama">
                <q-item-section>
                  <q-item-label class="text-weight-medium">{{ jurusan.nama }}</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <div class="text-weight-bold" :style="{ color: jurusan.persen >= 70 ? '#2e7d32' : '#f9a825' }">
                    {{ jurusan.persen }}%
                  </div>
                </q-item-section>
              </q-item>
            </q-list>
          </q-card>

          <!-- AGENDA HUBIN -->
          <q-card flat bordered class="content-card">
            <q-card-section>
              <div class="text-subtitle1 text-weight-bold text-primary">
                Agenda Hubungan Industri
              </div>
            </q-card-section>
            <q-separator />
            <q-list separator>
              <q-item v-for="agenda in agendaList" :key="agenda.judul">
                <q-item-section avatar>
                  <div class="agenda-date text-center">
                    <div class="text-weight-bold text-primary">{{ agenda.tanggal }}</div>
                    <div class="text-caption text-grey-7">{{ agenda.bulan }}</div>
                  </div>
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-medium">{{ agenda.judul }}</q-item-label>
                  <q-item-label caption>{{ agenda.waktu }}</q-item-label>
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

const notifCount = ref(4)

const kpis = ref([
  { label: 'Mitra Industri Aktif', value: '58', icon: 'business', color: 'primary' },
  { label: 'Siswa PKL Berjalan', value: '312', icon: 'work', color: 'teal' },
  { label: 'MoU Aktif', value: '46', icon: 'handshake', color: 'secondary' },
  { label: 'Alumni Terserap (1 Thn)', value: '81%', icon: 'trending_up', color: 'deep-orange' },
])

const totalSiswaPkl = 312

const mouPerhatian = ref([
  { mitra: 'PT Sinar Motor Jaya', ruangLingkup: 'Penempatan PKL & Rekrutmen Alumni TBSM', tanggalBerakhir: '20 Jul 2026' },
  { mitra: 'CV Jaringan Nusantara', ruangLingkup: 'Kerjasama Praktik Kerja TKJ', tanggalBerakhir: '02 Agu 2026' },
])

const distribusiPkl = ref([
  { nama: 'Astra Honda Motor (AHM)', jumlah: 84, warna: '#0a192f' },
  { nama: 'PT Sinar Motor Jaya', jumlah: 62, warna: '#c9a227' },
  { nama: 'CV Jaringan Nusantara', jumlah: 48, warna: '#00796b' },
  { nama: 'Bank Perkreditan Rakyat Sejahtera', jumlah: 40, warna: '#5e35b1' },
  { nama: 'Mitra Lainnya', jumlah: 78, warna: '#d84315' },
])

const maxPkl = Math.max(...distribusiPkl.value.map((m) => m.jumlah))
const barWidth = (jumlah) => Math.round((jumlah / maxPkl) * 100)

const monitoringPkl = ref([
  {
    id: 1,
    siswa: 'Bayu Aji Santoso',
    kelas: 'XII TBSM 1',
    mitra: 'Astra Honda Motor (AHM)',
    pembimbing: 'Bpk. Slamet Riyadi',
    status: 'Berjalan Baik',
  },
  {
    id: 2,
    siswa: 'Nadia Putri',
    kelas: 'XII TKJ 1',
    mitra: 'CV Jaringan Nusantara',
    pembimbing: 'Bpk. Yusuf Hidayat',
    status: 'Berjalan Baik',
  },
  {
    id: 3,
    siswa: 'Fajar Ramadhan',
    kelas: 'XII Akuntansi',
    mitra: 'Bank Perkreditan Rakyat Sejahtera',
    pembimbing: 'Ibu Ratna Sari',
    status: 'Perlu Kunjungan',
  },
  {
    id: 4,
    siswa: 'Rian Firmansyah',
    kelas: 'XII TBSM 2',
    mitra: 'PT Sinar Motor Jaya',
    pembimbing: 'Bpk. Slamet Riyadi',
    status: 'Absensi Bermasalah',
  },
])

const statusWarna = (status) => {
  if (status === 'Absensi Bermasalah') return 'red-7'
  if (status === 'Perlu Kunjungan') return 'orange-7'
  return 'teal-7'
}

const mitraBaru = ref([
  { nama: 'PT Teknologi Jaringan Mandiri', bidang: 'Infrastruktur Jaringan & Cloud', tanggal: '02 Jul 2026' },
  { nama: 'Koperasi Simpan Pinjam Makmur', bidang: 'Jasa Keuangan', tanggal: '25 Jun 2026' },
  { nama: 'Bengkel Resmi Yamaha Cendekia', bidang: 'Perawatan & Perbaikan Motor', tanggal: '18 Jun 2026' },
])

const lowonganBkk = ref([
  { posisi: 'Teknisi Sepeda Motor', perusahaan: 'PT Sinar Motor Jaya', batasAkhir: '25 Jul 2026', pelamar: 14 },
  { posisi: 'Network Support Staff', perusahaan: 'CV Jaringan Nusantara', batasAkhir: '30 Jul 2026', pelamar: 9 },
  { posisi: 'Staf Administrasi Keuangan', perusahaan: 'Bank Perkreditan Rakyat Sejahtera', batasAkhir: '05 Agu 2026', pelamar: 6 },
])

const keterserapanAlumni = ref([
  { nama: 'TBSM', persen: 88 },
  { nama: 'TKJ', persen: 79 },
  { nama: 'Akuntansi', persen: 65 },
  { nama: 'Multimedia', persen: 58 },
])

const agendaList = ref([
  { tanggal: '16', bulan: 'Jul', judul: 'Kunjungan Monitoring PKL ke PT Sinar Motor Jaya', waktu: '09:00 - Lokasi Mitra' },
  { tanggal: '21', bulan: 'Jul', judul: 'Perpanjangan MoU dengan CV Jaringan Nusantara', waktu: '10:00 - Ruang Hubin' },
  { tanggal: '30', bulan: 'Jul', judul: 'Job Matching Fair Bursa Kerja Khusus', waktu: '08:00 - Aula Sekolah' },
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

.bar-fill {
  height: 100%;
  border-radius: 6px;
  transition: width 0.4s ease;
}

.agenda-date {
  min-width: 44px;
}
</style>

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
            Dashboard Sarana &amp; Prasarana
          </div>
          <div class="text-caption text-grey-4 q-mt-xs">
            {{ todayLabel }}
          </div>
        </div>
        <div class="row items-center q-gutter-sm">
          <q-btn unelevated color="secondary" text-color="dark" label="Lapor Kerusakan" no-caps icon="report_problem" />
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
          <!-- ASET KRITIS / PERLU TINDAKAN -->
          <q-card flat bordered class="q-mb-lg content-card border-alert" v-if="asetKritis.length">
            <q-card-section>
              <div class="row items-center">
                <q-icon name="report" color="red" size="22px" class="q-mr-sm" />
                <div class="text-subtitle1 text-weight-bold text-red-9">
                  Aset Memerlukan Tindakan Segera
                </div>
              </div>
            </q-card-section>
            <q-separator />
            <q-list separator>
              <q-item v-for="aset in asetKritis" :key="aset.nama">
                <q-item-section avatar>
                  <q-avatar color="red-1" text-color="red-9" :icon="aset.icon" size="38px" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-medium">{{ aset.nama }}</q-item-label>
                  <q-item-label caption>{{ aset.lokasi }} • {{ aset.catatan }}</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-badge color="red" text-color="white" class="q-px-sm">
                    {{ aset.status }}
                  </q-badge>
                </q-item-section>
              </q-item>
            </q-list>
          </q-card>

          <!-- KONDISI ASET PER KATEGORI -->
          <q-card flat bordered class="q-mb-lg content-card">
            <q-card-section>
              <div class="row items-center justify-between">
                <div class="text-subtitle1 text-weight-bold text-primary">
                  Kondisi Aset per Kategori
                </div>
                <q-badge color="grey-3" text-color="grey-8" class="q-px-sm">
                  {{ totalAset }} unit terdaftar
                </q-badge>
              </div>
            </q-card-section>
            <q-separator />
            <q-card-section>
              <div v-for="kategori in kondisiAset" :key="kategori.nama" class="q-mb-md">
                <div class="row items-center justify-between q-mb-xs">
                  <div class="text-body2 text-weight-medium">{{ kategori.nama }}</div>
                  <div class="text-caption text-grey-7">{{ kategori.total }} unit</div>
                </div>
                <div class="stacked-bar-track">
                  <div class="stacked-segment" style="background-color: #2e7d32"
                    :style="{ width: pct(kategori.baik, kategori.total) + '%', backgroundColor: '#2e7d32' }" />
                  <div class="stacked-segment"
                    :style="{ width: pct(kategori.rusakRingan, kategori.total) + '%', backgroundColor: '#f9a825' }" />
                  <div class="stacked-segment"
                    :style="{ width: pct(kategori.rusakBerat, kategori.total) + '%', backgroundColor: '#c62828' }" />
                </div>
                <div class="row q-mt-xs q-gutter-md">
                  <div class="text-caption text-grey-7">
                    <span class="legend-dot" style="background-color: #2e7d32" /> Baik: {{ kategori.baik }}
                  </div>
                  <div class="text-caption text-grey-7">
                    <span class="legend-dot" style="background-color: #f9a825" /> Rusak Ringan: {{ kategori.rusakRingan
                    }}
                  </div>
                  <div class="text-caption text-grey-7">
                    <span class="legend-dot" style="background-color: #c62828" /> Rusak Berat: {{ kategori.rusakBerat }}
                  </div>
                </div>
              </div>
            </q-card-section>
          </q-card>

          <!-- PERMINTAAN PERBAIKAN -->
          <q-card flat bordered class="content-card">
            <q-card-section>
              <div class="row items-center justify-between">
                <div class="text-subtitle1 text-weight-bold text-primary">
                  Permintaan Perbaikan Terbaru
                </div>
                <q-btn flat dense color="primary" label="Lihat Semua" no-caps icon-right="arrow_forward" />
              </div>
            </q-card-section>
            <q-separator />
            <q-list separator>
              <q-item v-for="tiket in permintaanPerbaikan" :key="tiket.id">
                <q-item-section avatar>
                  <q-avatar :color="statusWarna(tiket.status)" text-color="white" icon="build" size="38px" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-medium">{{ tiket.judul }}</q-item-label>
                  <q-item-label caption>{{ tiket.lokasi }} • Dilaporkan oleh {{ tiket.pelapor }} • {{ tiket.waktu
                    }}</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-chip dense :color="statusWarna(tiket.status)" text-color="white" class="text-caption">
                    {{ tiket.status }}
                  </q-chip>
                </q-item-section>
              </q-item>
            </q-list>
          </q-card>
        </div>

        <!-- RIGHT COLUMN -->
        <div class="col-12 col-md-4">
          <!-- KETERSEDIAAN RUANG -->
          <q-card flat bordered class="q-mb-lg content-card">
            <q-card-section>
              <div class="text-subtitle1 text-weight-bold text-primary">
                Ketersediaan Ruang Hari Ini
              </div>
            </q-card-section>
            <q-separator />
            <q-list separator>
              <q-item v-for="ruang in ketersediaanRuang" :key="ruang.nama">
                <q-item-section avatar>
                  <q-avatar :color="ruang.status === 'Tersedia' ? 'teal-1' : 'orange-1'"
                    :text-color="ruang.status === 'Tersedia' ? 'teal-9' : 'orange-9'" icon="meeting_room" size="38px" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-medium">{{ ruang.nama }}</q-item-label>
                  <q-item-label caption>{{ ruang.keterangan }}</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-badge :color="ruang.status === 'Tersedia' ? 'teal' : 'orange'" text-color="white">
                    {{ ruang.status }}
                  </q-badge>
                </q-item-section>
              </q-item>
            </q-list>
          </q-card>

          <!-- INVENTARIS STOK RENDAH -->
          <q-card flat bordered class="q-mb-lg content-card">
            <q-card-section>
              <div class="text-subtitle1 text-weight-bold text-primary">
                Inventaris Perlu Restok
              </div>
            </q-card-section>
            <q-separator />
            <q-list separator>
              <q-item v-for="item in inventarisRendah" :key="item.nama">
                <q-item-section>
                  <q-item-label class="text-weight-medium">{{ item.nama }}</q-item-label>
                  <q-item-label caption>Stok minimum: {{ item.stokMin }} unit</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-badge color="red" text-color="white">
                    Sisa {{ item.stok }}
                  </q-badge>
                </q-item-section>
              </q-item>
            </q-list>
          </q-card>

          <!-- AGENDA PEMELIHARAAN -->
          <q-card flat bordered class="content-card">
            <q-card-section>
              <div class="text-subtitle1 text-weight-bold text-primary">
                Agenda Pemeliharaan
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

const notifCount = ref(3)

const kpis = ref([
  { label: 'Total Aset Terdaftar', value: '2.146', icon: 'inventory_2', color: 'primary' },
  { label: 'Aset Rusak Berat', value: '14', icon: 'report_problem', color: 'red-7' },
  { label: 'Ruang Tersedia Hari Ini', value: '28 / 34', icon: 'meeting_room', color: 'teal' },
  { label: 'Tiket Perbaikan Terbuka', value: '9', icon: 'build', color: 'orange-8' },
])

const totalAset = 2146

const asetKritis = ref([
  { nama: 'AC Ruang Lab TKJ 2', lokasi: 'Gedung B, Lantai 2', catatan: 'Tidak dingin sejak 3 hari lalu', status: 'Rusak Berat', icon: 'ac_unit' },
  { nama: 'Instalasi Listrik Bengkel TBSM', lokasi: 'Gedung Bengkel Otomotif', catatan: 'Korsleting minor, perlu pengecekan teknisi', status: 'Perlu Cek', icon: 'bolt' },
])

const kondisiAset = ref([
  { nama: 'Peralatan Praktik Bengkel', total: 486, baik: 402, rusakRingan: 58, rusakBerat: 26 },
  { nama: 'Komputer & Perangkat IT', total: 312, baik: 268, rusakRingan: 32, rusakBerat: 12 },
  { nama: 'Furnitur Ruang Kelas', total: 940, baik: 860, rusakRingan: 66, rusakBerat: 14 },
  { nama: 'Alat Kesenian & Olahraga', total: 168, baik: 140, rusakRingan: 22, rusakBerat: 6 },
])

const pct = (jumlah, total) => Math.round((jumlah / total) * 100)

const permintaanPerbaikan = ref([
  {
    id: 1,
    judul: 'Proyektor tidak menyala di Ruang XI TKJ 1',
    lokasi: 'Gedung B, Lantai 1',
    pelapor: 'Bpk. Hendra Gunawan',
    waktu: '1 jam lalu',
    status: 'Baru',
  },
  {
    id: 2,
    judul: 'Kran air di toilet siswa Gedung C bocor',
    lokasi: 'Gedung C',
    pelapor: 'Ibu Wulandari',
    waktu: '5 jam lalu',
    status: 'Diproses',
  },
  {
    id: 3,
    judul: 'Mesin bubut di bengkel TBSM perlu kalibrasi ulang',
    lokasi: 'Bengkel Otomotif',
    pelapor: 'Bpk. Slamet Riyadi',
    waktu: 'Kemarin, 14:20',
    status: 'Diproses',
  },
  {
    id: 4,
    judul: 'Penggantian lampu koridor Gedung A selesai',
    lokasi: 'Gedung A',
    pelapor: 'Satpam Sekolah',
    waktu: '2 hari lalu',
    status: 'Selesai',
  },
])

const statusWarna = (status) => {
  if (status === 'Baru') return 'red-6'
  if (status === 'Diproses') return 'orange-7'
  return 'teal-7'
}

const ketersediaanRuang = ref([
  { nama: 'Aula Utama', keterangan: 'Kapasitas 300 orang', status: 'Tersedia' },
  { nama: 'Ruang Rapat Guru', keterangan: 'Digunakan 09:00 - 11:00', status: 'Terpakai' },
  { nama: 'Lab Komputer TKJ 1', keterangan: 'Kapasitas 36 siswa', status: 'Tersedia' },
  { nama: 'Bengkel Otomotif TBSM', keterangan: 'Praktik kelas XI berlangsung', status: 'Terpakai' },
])

const inventarisRendah = ref([
  { nama: 'Oli Mesin untuk Praktik TBSM', stok: 4, stokMin: 15 },
  { nama: 'Kertas HVS A4', stok: 12, stokMin: 30 },
  { nama: 'Kabel LAN Cat6', stok: 8, stokMin: 20 },
])

const agendaList = ref([
  { tanggal: '15', bulan: 'Jul', judul: 'Servis Berkala AC Seluruh Gedung', waktu: '08:00 - Seluruh Area' },
  { tanggal: '20', bulan: 'Jul', judul: 'Pengecekan Instalasi Listrik Bengkel', waktu: '09:00 - Bengkel Otomotif' },
  { tanggal: '28', bulan: 'Jul', judul: 'Inventarisasi Aset Tahunan', waktu: '07:30 - Seluruh Gedung' },
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

.stacked-bar-track {
  width: 100%;
  height: 10px;
  border-radius: 6px;
  background-color: #eceff1;
  overflow: hidden;
  display: flex;
}

.stacked-segment {
  height: 100%;
  transition: width 0.4s ease;
}

.legend-dot {
  display: inline-block;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  margin-right: 4px;
}

.agenda-date {
  min-width: 44px;
}
</style>

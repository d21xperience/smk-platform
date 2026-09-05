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
            Dashboard Kesiswaan
          </div>
          <div class="text-caption text-grey-4 q-mt-xs">
            {{ todayLabel }}
          </div>
        </div>
        <div class="row items-center q-gutter-sm">
          <q-btn unelevated color="secondary" text-color="dark" label="Input Pelanggaran" no-caps icon="add" />
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
          <!-- SISWA PERLU PERHATIAN KHUSUS -->
          <q-card flat bordered class="q-mb-lg content-card border-alert" v-if="siswaPerhatian.length">
            <q-card-section>
              <div class="row items-center">
                <q-icon name="report" color="red" size="22px" class="q-mr-sm" />
                <div class="text-subtitle1 text-weight-bold text-red-9">
                  Siswa Memerlukan Perhatian Khusus
                </div>
              </div>
            </q-card-section>
            <q-separator />
            <q-list separator>
              <q-item v-for="siswa in siswaPerhatian" :key="siswa.nama">
                <q-item-section avatar>
                  <q-avatar color="red-1" text-color="red-9" icon="person" size="38px" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-medium">{{ siswa.nama }}</q-item-label>
                  <q-item-label caption>{{ siswa.kelas }} • {{ siswa.catatan }}</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-badge color="red" text-color="white" class="q-px-sm">
                    {{ siswa.poin }} poin
                  </q-badge>
                </q-item-section>
              </q-item>
            </q-list>
          </q-card>

          <!-- TREN PELANGGARAN -->
          <q-card flat bordered class="q-mb-lg content-card">
            <q-card-section>
              <div class="row items-center justify-between">
                <div class="text-subtitle1 text-weight-bold text-primary">
                  Tren Pelanggaran Siswa (6 Bulan Terakhir)
                </div>
                <q-badge color="grey-3" text-color="grey-8" class="q-px-sm">
                  {{ totalPelanggaranBulanIni }} kasus bulan ini
                </q-badge>
              </div>
            </q-card-section>
            <q-separator />
            <q-card-section>
              <div class="trend-chart row items-end no-wrap">
                <div v-for="bulan in trenPelanggaran" :key="bulan.label" class="trend-bar-wrap col text-center">
                  <div class="text-caption text-grey-7 q-mb-xs">{{ bulan.jumlah }}</div>
                  <div class="trend-bar" :style="{ height: trendHeight(bulan.jumlah) + 'px' }" />
                  <div class="text-caption text-grey-6 q-mt-xs">{{ bulan.label }}</div>
                </div>
              </div>
            </q-card-section>
          </q-card>

          <!-- KATEGORI PELANGGARAN PER KELAS -->
          <q-card flat bordered class="q-mb-lg content-card">
            <q-card-section>
              <div class="text-subtitle1 text-weight-bold text-primary">
                Distribusi Poin Pelanggaran per Tingkat Kelas
              </div>
            </q-card-section>
            <q-separator />
            <q-card-section>
              <div v-for="tingkat in poinPerTingkat" :key="tingkat.label" class="q-mb-md">
                <div class="row items-center justify-between q-mb-xs">
                  <div class="text-body2 text-weight-medium">{{ tingkat.label }}</div>
                  <div class="text-caption text-grey-7">{{ tingkat.total }} poin akumulasi</div>
                </div>
                <div class="bar-track">
                  <div class="bar-fill"
                    :style="{ width: barWidth(tingkat.total) + '%', backgroundColor: tingkat.warna }" />
                </div>
              </div>
            </q-card-section>
          </q-card>

          <!-- LOG KASUS TERBARU -->
          <q-card flat bordered class="content-card">
            <q-card-section>
              <div class="text-subtitle1 text-weight-bold text-primary">
                Catatan Kasus Terbaru
              </div>
            </q-card-section>
            <q-separator />
            <q-list separator>
              <q-item v-for="kasus in kasusTerbaru" :key="kasus.id">
                <q-item-section avatar>
                  <q-avatar :color="kategoriWarna(kasus.kategori)" text-color="white" size="38px">
                    {{ kasus.kategori.charAt(0) }}
                  </q-avatar>
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-medium">{{ kasus.judul }}</q-item-label>
                  <q-item-label caption>{{ kasus.siswa }} • {{ kasus.kelas }} • {{ kasus.waktu }}</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-chip dense :color="kategoriWarna(kasus.kategori)" text-color="white" class="text-caption">
                    {{ kasus.kategori }}
                  </q-chip>
                </q-item-section>
              </q-item>
            </q-list>
          </q-card>
        </div>

        <!-- RIGHT COLUMN -->
        <div class="col-12 col-md-4">
          <!-- EKSTRAKURIKULER -->
          <q-card flat bordered class="q-mb-lg content-card">
            <q-card-section>
              <div class="text-subtitle1 text-weight-bold text-primary">
                Peserta Ekstrakurikuler
              </div>
            </q-card-section>
            <q-separator />
            <q-list separator>
              <q-item v-for="ekskul in ekstrakurikuler" :key="ekskul.nama">
                <q-item-section avatar>
                  <q-avatar :color="ekskul.warna" text-color="white" :icon="ekskul.icon" size="38px" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-medium">{{ ekskul.nama }}</q-item-label>
                  <q-item-label caption>Pembina: {{ ekskul.pembina }}</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <div class="text-weight-bold text-primary">{{ ekskul.peserta }}</div>
                </q-item-section>
              </q-item>
            </q-list>
          </q-card>

          <!-- PRESTASI TERBARU -->
          <q-card flat bordered class="q-mb-lg content-card">
            <q-card-section>
              <div class="text-subtitle1 text-weight-bold text-primary">
                Prestasi Terbaru
              </div>
            </q-card-section>
            <q-separator />
            <q-list separator>
              <q-item v-for="prestasi in prestasiList" :key="prestasi.judul">
                <q-item-section avatar>
                  <q-avatar color="amber-2" text-color="amber-9" icon="emoji_events" size="38px" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-medium">{{ prestasi.judul }}</q-item-label>
                  <q-item-label caption>{{ prestasi.siswa }} • {{ prestasi.tanggal }}</q-item-label>
                </q-item-section>
              </q-item>
            </q-list>
          </q-card>

          <!-- AGENDA KESISWAAN -->
          <q-card flat bordered class="content-card">
            <q-card-section>
              <div class="text-subtitle1 text-weight-bold text-primary">
                Agenda Kesiswaan
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

const notifCount = ref(6)

const kpis = ref([
  { label: 'Siswa Aktif', value: '1.248', icon: 'groups', color: 'primary' },
  { label: 'Kasus Pelanggaran Bulan Ini', value: '17', icon: 'gavel', color: 'red-7' },
  { label: 'Siswa Berprestasi', value: '32', icon: 'emoji_events', color: 'amber-8' },
  { label: 'Peserta Ekstrakurikuler', value: '486', icon: 'sports', color: 'teal' },
])

const siswaPerhatian = ref([
  { nama: 'Rian Firmansyah', kelas: 'XI TBSM 2', catatan: 'Alpa berulang & poin pelanggaran tinggi', poin: 85 },
  { nama: 'Dedi Kurniawan', kelas: 'X TKJ 1', catatan: 'Terlibat perkelahian di lingkungan sekolah', poin: 70 },
])

const totalPelanggaranBulanIni = 17

const trenPelanggaran = ref([
  { label: 'Feb', jumlah: 12 },
  { label: 'Mar', jumlah: 15 },
  { label: 'Apr', jumlah: 9 },
  { label: 'Mei', jumlah: 14 },
  { label: 'Jun', jumlah: 8 },
  { label: 'Jul', jumlah: 17 },
])

const maxTren = Math.max(...trenPelanggaran.value.map((b) => b.jumlah))
const trendHeight = (jumlah) => Math.max(8, Math.round((jumlah / maxTren) * 120))

const poinPerTingkat = ref([
  { label: 'Kelas X', total: 420, warna: '#0a192f' },
  { label: 'Kelas XI', total: 610, warna: '#c9a227' },
  { label: 'Kelas XII', total: 260, warna: '#00796b' },
])

const maxPoin = Math.max(...poinPerTingkat.value.map((p) => p.total))
const barWidth = (total) => Math.round((total / maxPoin) * 100)

const kasusTerbaru = ref([
  {
    id: 1,
    judul: 'Terlambat masuk sekolah lebih dari 3 kali dalam sepekan',
    siswa: 'Ahmad Fauzan',
    kelas: 'X TKJ 2',
    waktu: '2 jam lalu',
    kategori: 'Ringan',
  },
  {
    id: 2,
    judul: 'Tidak mengenakan atribut seragam lengkap',
    siswa: 'Siti Nurhaliza',
    kelas: 'XI Akuntansi',
    waktu: 'Kemarin, 09:15',
    kategori: 'Ringan',
  },
  {
    id: 3,
    judul: 'Membawa rokok elektrik ke lingkungan sekolah',
    siswa: 'Rian Firmansyah',
    kelas: 'XI TBSM 2',
    waktu: '2 hari lalu',
    kategori: 'Berat',
  },
  {
    id: 4,
    judul: 'Bolos pada jam pelajaran ke-5 dan ke-6',
    siswa: 'Dedi Kurniawan',
    kelas: 'X TKJ 1',
    waktu: '3 hari lalu',
    kategori: 'Sedang',
  },
])

const kategoriWarna = (kategori) => {
  if (kategori === 'Berat') return 'red-8'
  if (kategori === 'Sedang') return 'orange-8'
  return 'blue-grey-6'
}

const ekstrakurikuler = ref([
  { nama: 'Pramuka', pembina: 'Bpk. Sutrisno', peserta: 142, icon: 'terrain', warna: 'green-8' },
  { nama: 'Futsal', pembina: 'Bpk. Andi Wijaya', peserta: 68, icon: 'sports_soccer', warna: 'indigo-8' },
  { nama: 'PMR', pembina: 'Ibu Ratna Sari', peserta: 54, icon: 'medical_services', warna: 'red-6' },
  { nama: 'Robotika', pembina: 'Bpk. Yusuf Hidayat', peserta: 38, icon: 'precision_manufacturing', warna: 'deep-purple-6' },
])

const prestasiList = ref([
  { judul: 'Juara 1 Mekanik Edukasi Tingkat Provinsi', siswa: 'Bayu Aji Santoso - XII TBSM 1', tanggal: '30 Mei 2026' },
  { judul: 'Juara 2 LKS Bidang IT Network Systems Administration', siswa: 'Nadia Putri - XI TKJ 1', tanggal: '18 Mei 2026' },
  { judul: 'Juara 3 Lomba Debat Bahasa Inggris Kabupaten', siswa: 'Fajar Ramadhan - XI Akuntansi', tanggal: '05 Mei 2026' },
])

const agendaList = ref([
  { tanggal: '13', bulan: 'Jul', judul: 'Pembinaan Karakter & Kedisiplinan Siswa Kelas X', waktu: '08:00 - Aula Sekolah' },
  { tanggal: '19', bulan: 'Jul', judul: 'Seleksi Ketua OSIS Periode 2026/2027', waktu: '09:00 - Ruang OSIS' },
  { tanggal: '25', bulan: 'Jul', judul: 'Rapat Evaluasi Tata Tertib Semester Ganjil', waktu: '13:00 - Ruang Guru' },
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

.trend-chart {
  height: 170px;
  gap: 12px;
}

.trend-bar-wrap {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: flex-end;
}

.trend-bar {
  width: 100%;
  max-width: 32px;
  background: linear-gradient(180deg, #123a63 0%, #0a192f 100%);
  border-radius: 4px 4px 0 0;
  transition: height 0.4s ease;
}

.agenda-date {
  min-width: 44px;
}
</style>

<template>
  <q-page class="dashboard-page bg-grey-1">
    <!-- HEADER -->
    <div class="container max-width-center q-px-md q-py-lg">
      <!-- KPI CARDS -->
      <!-- <div class="row q-col-gutter-md q-mb-lg">
        <div class="col-6 col-md-3" v-for="kpi in kpis" :key="kpi.label">
          <q-card flat bordered class="kpi-card">
            <q-card-section class="row items-center no-wrap">
              <q-avatar :color="kpi.color" text-color="white" :icon="kpi.icon" size="46px" class="q-mr-md" />
              <div>
                <div class="text-h6 text-weight-bold text-primary">{{ kpi.value }}</div>
                <div class="text-caption text-grey-7">{{ kpi.label }}</div>
              </div>
            </q-card-section>
            <div class="kpi-trend q-px-md q-pb-sm" v-if="kpi.trend">
              <q-icon :name="kpi.trend > 0 ? 'trending_up' : 'trending_down'"
                :color="kpi.trend > 0 ? 'positive' : 'negative'" size="16px" />
              <span class="text-caption q-ml-xs" :class="kpi.trend > 0 ? 'text-positive' : 'text-negative'">
                {{ Math.abs(kpi.trend) }}% dari bulan lalu
              </span>
            </div>
          </q-card>
        </div>
      </div> -->





      <!-- KEHADIRAN HARI INI -->
      <q-card flat bordered class="q-mb-lg content-card">
        <q-card-section>
          <div class="text-subtitle1 text-weight-bold text-primary">
            Rekap Kehadiran Hari Ini
          </div>
        </q-card-section>
        <q-separator />
        <q-card-section class="row q-col-gutter-md">
          <div class="col-6 col-sm-3 text-center" v-for="item in kehadiran" :key="item.label">
            <div class="text-h5 text-weight-bold" :style="{ color: item.warna }">
              {{ item.persen }}%
            </div>
            <div class="text-caption text-grey-7 q-mt-xs">{{ item.label }}</div>
            <div class="text-caption text-grey-5">{{ item.jumlah }} siswa</div>
          </div>
        </q-card-section>
      </q-card>




      <!-- KPI CARDS (MENGGUNAKAN KOMPONEN TERPISAH) -->
      <div class="row q-col-gutter-md q-mb-lg">
        <div class="col-6 col-md-3" v-for="(kpi, index) in kpis" :key="index">
          <KpiCard :label="kpi.label" :value="kpi.value" :icon="kpi.icon" :color="kpi.color" :trend="kpi.trend"
            :route-to="kpi.routeTo" @click="handleKpiClick(kpi.routeTo)" />
        </div>
      </div>
      <div class="row q-col-gutter-lg">
        <!-- LEFT COLUMN -->
        <div class="col-12 col-md-8">
          <!-- DISTRIBUSI SISWA PER JURUSAN -->
          <q-card flat bordered class="q-mb-lg content-card">
            <q-card-section>
              <div class="row items-center justify-between">
                <div class="text-subtitle1 text-weight-bold text-primary">
                  Distribusi Siswa per Kompetensi Keahlian
                </div>
                <q-badge color="grey-3" text-color="grey-8" class="q-px-sm">
                  Total {{ totalSiswa }} Siswa
                </q-badge>
              </div>
            </q-card-section>
            <q-separator />
            <q-card-section>
              <div v-for="jurusan in jurusanDistribusi" :key="jurusan.nama" class="q-mb-md">
                <div class="row items-center justify-between q-mb-xs">
                  <div class="text-body2 text-weight-medium">{{ jurusan.nama }}</div>
                  <div class="text-caption text-grey-7">{{ jurusan.jumlah }} siswa</div>
                </div>
                <div class="bar-track">
                  <div class="bar-fill"
                    :style="{ width: barWidth(jurusan.jumlah) + '%', backgroundColor: jurusan.warna }" />
                </div>
              </div>
            </q-card-section>
          </q-card>



          <!-- AKTIVITAS TERBARU -->
          <q-card flat bordered class="content-card">
            <q-card-section>
              <div class="text-subtitle1 text-weight-bold text-primary">
                Aktivitas Terbaru
              </div>
            </q-card-section>
            <q-separator />
            <q-list separator>
              <q-item v-for="log in aktivitas" :key="log.id">
                <q-item-section avatar>
                  <q-avatar :color="log.warna" text-color="white" :icon="log.icon" size="38px" />
                </q-item-section>
                <q-item-section>
                  <q-item-label>{{ log.judul }}</q-item-label>
                  <q-item-label caption>{{ log.waktu }}</q-item-label>
                </q-item-section>
              </q-item>
            </q-list>
          </q-card>
        </div>

        <!-- RIGHT COLUMN -->
        <div class="col-12 col-md-4">
          <!-- RINGKASAN TENAGA PENDIDIK -->
          <q-card flat bordered class="q-mb-lg content-card">
            <q-card-section>
              <div class="text-subtitle1 text-weight-bold text-primary">
                Tenaga Pendidik & Kependidikan
              </div>
            </q-card-section>
            <q-separator />
            <q-list separator>
              <q-item v-for="staff in ringkasanStaff" :key="staff.label">
                <q-item-section>
                  <q-item-label>{{ staff.label }}</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <div class="text-weight-bold text-primary">{{ staff.jumlah }}</div>
                </q-item-section>
              </q-item>
            </q-list>
          </q-card>

          <!-- AGENDA -->
          <q-card flat bordered class="q-mb-lg content-card">
            <q-card-section>
              <div class="text-subtitle1 text-weight-bold text-primary">
                Agenda Mendatang
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

          <!-- PENGUMUMAN -->
          <q-card flat bordered class="content-card">
            <q-card-section>
              <div class="text-subtitle1 text-weight-bold text-primary">
                Pengumuman
              </div>
            </q-card-section>
            <q-separator />
            <q-card-section>
              <q-banner v-for="pengumuman in pengumumanList" :key="pengumuman.judul" dense
                class="q-mb-sm pengumuman-banner" :class="pengumuman.urgent ? 'bg-red-1' : 'bg-blue-1'">
                <template v-slot:avatar>
                  <q-icon :name="pengumuman.urgent ? 'priority_high' : 'campaign'"
                    :color="pengumuman.urgent ? 'red' : 'primary'" />
                </template>
                <div class="text-body2 text-weight-medium">{{ pengumuman.judul }}</div>
                <div class="text-caption text-grey-7">{{ pengumuman.tanggal }}</div>
              </q-banner>
            </q-card-section>
          </q-card>
        </div>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import KpiCard from '@/components/KpiCard.vue'
import { useProtectedNavigation } from '@/composables/guard/useProtectedNavigation'
import { ref } from 'vue'

const { navigateTo } = useProtectedNavigation()
const kpis = ref([
  {
    label: 'Total Siswa',
    value: '1.248',
    icon: 'groups',
    color: 'primary',
    trend: 3.2,
    routeTo: '/kesiswaan/direktori-siswa' // <-- FITUR BARU: Arahkan ke route direktori
  },
  // { label: 'Total Siswa', value: '1.248', icon: 'groups', color: 'primary', trend: 3.2 },
  { label: 'Total Guru & Staf', value: '86', icon: 'badge', color: 'secondary', trend: 1.1, routeTo: '/kepegawaian/data-guru' },
  { label: 'Total Rombel', value: '42', icon: 'meeting_room', color: 'teal', trend: 0 },
  { label: 'Kompetensi Keahlian', value: '6', icon: 'school', color: 'deep-orange', trend: null },
])

const totalSiswa = 1248

const jurusanDistribusi = ref([
  { nama: 'Teknik & Bisnis Sepeda Motor (TBSM)', jumlah: 412, warna: '#0a192f' },
  { nama: 'Teknik Komputer & Jaringan', jumlah: 356, warna: '#c9a227' },
  { nama: 'Akuntansi & Keuangan Lembaga', jumlah: 248, warna: '#00796b' },
  { nama: 'Otomatisasi Tata Kelola Perkantoran', jumlah: 132, warna: '#5e35b1' },
  { nama: 'Multimedia', jumlah: 100, warna: '#d84315' },
])

const maxJurusan = Math.max(...jurusanDistribusi.value.map((j) => j.jumlah))
const barWidth = (jumlah) => Math.round((jumlah / maxJurusan) * 100)

const kehadiran = ref([
  { label: 'Hadir', persen: 92, jumlah: 1148, warna: '#2e7d32' },
  { label: 'Izin', persen: 3, jumlah: 37, warna: '#1976d2' },
  { label: 'Sakit', persen: 3, jumlah: 32, warna: '#f9a825' },
  { label: 'Alpa', persen: 2, jumlah: 31, warna: '#c62828' },
])

const ringkasanStaff = ref([
  { label: 'Guru Tetap (PNS/Yayasan)', jumlah: 54 },
  { label: 'Guru Honorer', jumlah: 22 },
  { label: 'Tenaga Kependidikan', jumlah: 8 },
  { label: 'Wali Kelas Aktif', jumlah: 42 },
])

const aktivitas = ref([
  {
    id: 1,
    judul: 'Nilai UAS Semester Genap kelas XII TBSM 1 telah diinput',
    waktu: '10 menit lalu',
    icon: 'grading',
    warna: 'primary',
  },
  {
    id: 2,
    judul: 'Guru baru a.n. Rizka Amalia bergabung di jurusan TKJ',
    waktu: '1 jam lalu',
    icon: 'person_add',
    warna: 'teal',
  },
  {
    id: 3,
    judul: 'Jadwal Ujian Praktik Kejuruan (UPK) telah diperbarui',
    waktu: '3 jam lalu',
    icon: 'event_note',
    warna: 'deep-orange',
  },
  {
    id: 4,
    judul: 'Absensi otomatis melalui RFID diaktifkan untuk seluruh kelas X',
    waktu: 'Kemarin, 16:45',
    icon: 'fingerprint',
    warna: 'secondary',
  },
])

const agendaList = ref([
  { tanggal: '14', bulan: 'Jul', judul: 'Rapat Koordinasi Wali Kelas', waktu: '09:00 - Ruang Rapat Utama' },
  { tanggal: '17', bulan: 'Jul', judul: 'Ujian Praktik Kejuruan TBSM', waktu: '08:00 - Bengkel Otomotif' },
  { tanggal: '22', bulan: 'Jul', judul: 'Kunjungan Industri Mitra AHM', waktu: '07:30 - Aula Sekolah' },
])

const pengumumanList = ref([
  {
    judul: 'Batas akhir input nilai rapor semester genap: 20 Juli 2026',
    tanggal: '10 Juli 2026',
    urgent: true,
  },
  {
    judul: 'Pemeliharaan sistem SIAKAD akan dilakukan Sabtu malam',
    tanggal: '9 Juli 2026',
    urgent: false,
  },
])

// Fungsi handler menjadi sangat deklaratif dan bersih
const handleKpiClick = (targetRoute) => {
  if (targetRoute) {
    // Cukup kirim rute dan pesan opsional.
    // Role akan otomatis dibaca dari router.resolve(targetRoute).meta.allowedRoles
    navigateTo(
      targetRoute,
      'Data bersifat rahasia. Silakan login dengan akun yang sah.'
    )
  }
}





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

.kpi-trend {
  display: flex;
  align-items: center;
}

.content-card {
  border-radius: 8px;
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

.pengumuman-banner {
  border-radius: 6px;
}
</style>

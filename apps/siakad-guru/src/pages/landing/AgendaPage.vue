<template>
  <q-page class="bg-grey-1 q-py-xl">
    <div class="container q-px-md max-width-center">

      <!-- 1. HEADER HALAMAN -->
      <div class="text-center q-mb-xl">
        <div class="text-h4 text-sm-h3 text-weight-bold text-primary">
          Kalender Akademik & <span class="text-amber">Agenda</span>
        </div>
        <p class="text-body1 text-grey-7 q-mt-sm">
          Informasi resmi mengenai jadwal kegiatan sekolah, pelaksanaan ujian, dan hari libur nasional.
        </p>
      </div>

      <div class="row q-col-gutter-lg">
        <!-- 2. KOLOM KIRI: DAFTAR AGENDA & KEGIATAN MENDATANG (7/12 DESKTOP) -->
        <div class="col-12 col-md-7">
          <q-card flat bordered class="bg-white rounded-borders shadow-1">
            <q-card-section class="bg-primary text-white q-py-md text-weight-bold text-subtitle1 row items-center">
              <q-icon name="event" class="q-mr-sm" /> Kegiatan Sekolah & Jadwal Ujian
            </q-card-section>

            <q-card-section class="q-pa-none">
              <q-list separator>
                <q-item v-for="(agenda, index) in agendaList" :key="index" class="q-py-md">
                  <!-- Bagian Tanggal Bergaya Kalender Mini -->
                  <q-item-section avatar class="flex flex-center">
                    <div class="date-box text-center rounded-borders overflow-hidden shadow-1">
                      <div class="date-month bg-primary text-white text-weight-bold text-caption uppercase">{{
                        agenda.month }}</div>
                      <div class="date-day bg-grey-2 text-grey-9 text-weight-bolder text-h6">{{ agenda.day }}</div>
                    </div>
                  </q-item-section>

                  <!-- Detail Informasi Kegiatan -->
                  <q-item-section>
                    <div class="row items-center q-gutter-xs">
                      <q-badge :color="getCategoryColor(agenda.category)" class="text-weight-bold text-caption q-px-sm">
                        {{ agenda.category }}
                      </q-badge>
                    </div>
                    <q-item-label class="text-subtitle1 text-weight-bold text-grey-9 q-mt-xs">
                      {{ agenda.title }}
                    </q-item-label>
                    <q-item-label caption class="text-grey-7 row items-center q-mt-xs">
                      <q-icon name="schedule" size="xs" class="q-mr-xs" /> Waktu: {{ agenda.time }}
                      <q-icon name="place" size="xs" class="q-ml-md q-mr-xs" /> Lokasi: {{ agenda.location }}
                    </q-item-label>
                  </q-item-section>
                </q-item>
              </q-list>
            </q-card-section>
          </q-card>
        </div>

        <!-- 3. KOLOM KANAN: DAFTAR HARI LIBUR RESMI (5/12 DESKTOP) -->
        <div class="col-12 col-md-5">
          <q-card flat bordered class="bg-white rounded-borders shadow-1">
            <q-card-section class="bg-negative text-white q-py-md text-weight-bold text-subtitle1 row items-center">
              <q-icon name="no_accounts" class="q-mr-sm" /> Pengumuman Hari Libur
            </q-card-section>

            <q-card-section class="q-pa-none">
              <q-list separator>
                <q-item v-for="(holiday, idx) in holidayList" :key="idx" class="q-py-md bg-red-1">
                  <q-item-section avatar>
                    <q-avatar color="negative" text-color="white" icon="hotel_class" size="md" />
                  </q-item-section>

                  <q-item-section>
                    <q-item-label class="text-subtitle2 text-weight-bold text-grey-9">
                      {{ holiday.title }}
                    </q-item-label>
                    <q-item-label caption class="text-negative text-weight-medium q-mt-xs">
                      <q-icon name="calendar_today" size="xs" class="q-mr-xs" /> {{ holiday.dateRange }}
                    </q-item-label>
                  </q-item-section>
                </q-item>
              </q-list>
            </q-card-section>
          </q-card>

          <!-- PENGUMUMAN DARURAT SEMENTARA (SNEAK PEEK ALERT) -->
          <q-banner rounded class="bg-amber-2 text-grey-9 q-mt-lg shadow-1 border-amber">
            <template v-slot:avatar>
              <q-icon name="warning" color="amber-9" />
            </template>
            <div class="text-weight-bold text-caption text-sm-body2">
              Pemberitahuan: Selama Libur Akhir Semester, pelayanan administrasi Tata Usaha (TU) fisik di sekolah tetap
              buka
              terbatas pada jam 08.00 - 12.00 WIB.
            </div>
          </q-banner>
        </div>
      </div>

    </div>
  </q-page>
</template>

<script setup>
import { ref } from 'vue'

// Data Kegiatan / Agenda Terdekat Sekolah
const agendaList = ref([
  {
    day: '08',
    month: 'Jun',
    category: 'Ujian',
    title: 'Pelaksanaan Penilaian Akhir Semester (PAS) Genap',
    time: '07.30 - 12.30 WIB',
    location: 'Ruang Kelas Teori & Lab Komputer'
  },
  {
    day: '19',
    month: 'Jun',
    category: 'Kegiatan',
    title: 'Classmeeting & Turnamen Futsal OSIS Cup',
    time: '08.00 - 15.00 WIB',
    location: 'Gedung Olahraga (GOR) Indoor'
  },
  {
    day: '26',
    month: 'Jun',
    category: 'Akademik',
    title: 'Pembagian Rapor Hasil Belajar Semester Genap',
    time: '08.00 - 11.00 WIB (Didampingi Wali Murid)',
    location: 'Ruang Kelas Masing-masing'
  }
])

// Data Hari Libur Sekolah Terdekat
const holidayList = ref([
  {
    title: 'Hari Libur Nasional: Idul Adha 1447 H',
    dateRange: '27 - 28 Juni 2026'
  },
  {
    title: 'Libur Akhir Semester Genap & Tahun Ajaran',
    dateRange: '29 Juni - 12 Juli 2026'
  }
])

function getCategoryColor(category) {
  if (category === 'Ujian') return 'orange-8'
  if (category === 'Kegiatan') return 'blue-7'
  return 'green-7'
}
</script>

<style scoped>
.max-width-center {
  max-width: 1100px;
  margin: 0 auto;
}

/* Styling Mini Calendar Badge */
.date-box {
  width: 55px;
  border: 1px solid #e0e0e0;
}

.date-month {
  padding: 2px 0;
  font-size: 10px;
  letter-spacing: 1px;
}

.date-day {
  padding: 4px 0;
}

.border-amber {
  border: 1px solid #ffe082;
}
</style>

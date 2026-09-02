<!--
  KalenderPendidikanCetak.vue
  Contoh pemakaian MonthCalendarGrid pada konteks berbeda: kalender pendidikan bergaya cetak resmi
  (3 bulan per baris, sel tanggal diberi blok warna solid, keterangan kegiatan di bawah tiap grid).
  Ini membuktikan komponennya reusable - halaman dashboard interaktif (KalenderAkademik.vue) dan
  halaman cetak ini memakai komponen & sumber data yang sama, hanya beda prop `variant`.
-->
<template>
  <q-page class="print-calendar-page q-pa-md">
    <!-- HEADER (disembunyikan saat cetak) -->
    <div class="row items-center justify-between q-mb-md wrap q-gutter-sm no-print">
      <div>
        <div class="text-h6 text-weight-bold text-grey-9">Kalender Pendidikan</div>
        <div class="text-caption text-grey-6">
          Tahun Ajaran {{ academicYear }} &middot; Tampilan Cetak
        </div>
      </div>
      <div class="row q-gutter-sm">
        <q-btn
          no-caps
          outline
          color="primary"
          icon="calendar_month"
          label="Tampilan Interaktif"
          :to="{ name: 'kalender-akademik' }"
        />
        <q-btn no-caps unelevated color="primary" icon="print" label="Cetak" @click="handlePrint" />
      </div>
    </div>

    <!-- KOP KALENDER -->
    <div class="print-header q-mb-md">
      <div class="text-h6 text-weight-bold text-primary">KALENDER PENDIDIKAN</div>
      <div class="text-subtitle2 text-grey-8">
        SMK Pasundan Jatinangor &middot; Tahun Ajaran {{ academicYear }}
      </div>
    </div>

    <!-- SEMESTER 1 -->
    <div class="text-subtitle1 text-weight-bold text-grey-8 q-mb-sm">
      Semester 1 ({{ semesterData['1'].label }})
    </div>
    <div class="row q-col-gutter-md q-mb-lg">
      <div
        class="col-12 col-sm-6 col-md-4"
        v-for="m in semester1Months"
        :key="`${m.year}-${m.month}`"
      >
        <q-card flat bordered class="month-card">
          <q-card-section>
            <MonthCalendarGrid
              :year="m.year"
              :month="m.month"
              :events="allEvents"
              variant="compact"
              week-start="sunday"
              show-legend
            />
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- SEMESTER 2 -->
    <div class="text-subtitle1 text-weight-bold text-grey-8 q-mb-sm">
      Semester 2 ({{ semesterData['2'].label }})
    </div>
    <div class="row q-col-gutter-md q-mb-lg">
      <div
        class="col-12 col-sm-6 col-md-4"
        v-for="m in semester2Months"
        :key="`${m.year}-${m.month}`"
      >
        <q-card flat bordered class="month-card">
          <q-card-section>
            <MonthCalendarGrid
              :year="m.year"
              :month="m.month"
              :events="allEvents"
              variant="compact"
              week-start="sunday"
              show-legend
            />
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- LEGENDA KATEGORI -->
    <q-card flat bordered class="q-pa-md">
      <div class="text-subtitle2 text-weight-bold q-mb-sm">Keterangan Warna</div>
      <div class="row q-gutter-md">
        <div v-for="cat in categories" :key="cat.value" class="row items-center q-gutter-xs">
          <div
            class="legend-swatch"
            :style="{ backgroundColor: cat.fill, borderColor: cat.text }"
          />
          <span class="text-caption text-grey-8">{{ cat.label }}</span>
        </div>
      </div>
    </q-card>
  </q-page>
</template>

<script setup>
//import { computed } from 'vue'
import MonthCalendarGrid from '@/components/MonthCalendarGrid.vue'
import {
  academicYear,
  CALENDAR_CATEGORIES,
  semesterData,
  getAllEvents,
} from '@/data/academicCalendarData'

const categories = CALENDAR_CATEGORIES
const allEvents = getAllEvents()

// TODO: sesuaikan rentang bulan bila tahun ajaran sekolah kamu tidak mulai Juli / berakhir Juni
const semester1Months = [
  { year: 2026, month: 6 }, // Juli
  { year: 2026, month: 7 }, // Agustus
  { year: 2026, month: 8 }, // September
  { year: 2026, month: 9 }, // Oktober
  { year: 2026, month: 10 }, // November
  { year: 2026, month: 11 }, // Desember
]

const semester2Months = [
  { year: 2027, month: 0 }, // Januari
  { year: 2027, month: 1 }, // Februari
  { year: 2027, month: 2 }, // Maret
  { year: 2027, month: 3 }, // April
  { year: 2027, month: 4 }, // Mei
  { year: 2027, month: 5 }, // Juni
]

function handlePrint() {
  // TODO: bila butuh PDF terpisah (bukan print dialog browser), ganti dengan
  // permintaan ke Go backend, misalnya GET /api/kalender-akademik/export?format=pdf
  window.print()
}
</script>

<style scoped>
.print-calendar-page {
  max-width: 1300px;
  margin: 0 auto;
}

.print-header {
  text-align: center;
  border-bottom: 2px solid #1a4a7a;
  padding-bottom: 8px;
}

.month-card {
  border-radius: 8px;
}

.legend-swatch {
  width: 16px;
  height: 16px;
  border-radius: 3px;
  border: 1px solid;
}

/* Saat dicetak: sembunyikan tombol aksi dan biarkan grid mengikuti tata letak halaman */
@media print {
  .no-print {
    display: none !important;
  }
  .month-card {
    break-inside: avoid;
  }
}
</style>

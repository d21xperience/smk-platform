<!--
  MonthCalendarGrid.vue
  Komponen kalender bulanan yang reusable (tanpa library eksternal, custom CSS grid).

  Dua variant:
  - "interactive" : dipakai pada dashboard (KalenderAkademik.vue). Sel ada efek hover, highlight
                     "hari ini", chip event berwarna dengan teks di dalam sel, dan bisa diklik.
  - "compact"      : dipakai pada tampilan cetak (KalenderPendidikanCetak.vue), meniru kalender
                     pendidikan resmi: sel tanggal diberi warna latar solid untuk rentang kegiatan,
                     tanggal Minggu & libur nasional ditulis merah tanpa blok latar, tanpa teks
                     kegiatan di dalam sel (keterangan ditampilkan sebagai daftar di bawah grid).

  Props utama: year, month (0-based), events, variant, weekStart ('monday'|'sunday'), showLegend.
  Emits: 'day-click' saat sebuah sel diklik (berguna untuk membuka dialog detail hari).
-->
<template>
  <div class="mcg" :class="`mcg--${variant}`">
    <div class="mcg__header">{{ monthLabel }}</div>

    <div class="mcg__weekdays">
      <div
        v-for="label in weekdayLabels"
        :key="label"
        class="mcg__weekday"
        :class="{ 'mcg__weekday--sunday': label === 'Minggu' }"
      >
        {{ label }}
      </div>
    </div>

    <div class="mcg__grid">
      <div
        v-for="cell in cells"
        :key="cell.iso"
        class="mcg__cell"
        :class="{
          'mcg__cell--muted': !cell.inCurrentMonth,
          'mcg__cell--today': cell.iso === today,
          'mcg__cell--clickable': variant === 'interactive',
        }"
        :style="cellFillStyle(cell)"
        @click="$emit('day-click', cell)"
      >
        <div class="mcg__date" :style="dateTextStyle(cell)">{{ cell.date }}</div>

        <div v-if="variant === 'interactive'" class="mcg__events">
          <div
            v-for="ev in cell.events.slice(0, 2)"
            :key="ev.id"
            class="mcg__event-chip"
            :style="{ backgroundColor: colorFor(ev.kategori).chip }"
          >
            <span class="mcg__event-text">{{ ev.nama }}</span>
          </div>
          <div v-if="cell.events.length > 2" class="mcg__event-more">
            +{{ cell.events.length - 2 }} lainnya
          </div>
        </div>
      </div>
    </div>

    <div class="mcg__legend" v-if="showLegend && legendItems.length">
      <div v-for="item in legendItems" :key="item.id" class="mcg__legend-item">
        <span class="mcg__legend-date">{{ item.dateLabel }}</span>
        <span class="mcg__legend-text">{{ item.nama }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import {
  MONTHS_ID,
  WEEKDAYS_MONDAY_FIRST,
  WEEKDAYS_SUNDAY_FIRST,
  todayISO,
  formatRange,
  monthBoundsISO,
  buildMonthGrid,
} from '@/composables/useCalendarDate'

const props = defineProps({
  year: { type: Number, required: true },
  month: { type: Number, required: true }, // 0-based, Januari = 0
  events: { type: Array, default: () => [] },
  variant: { type: String, default: 'interactive' }, // 'interactive' | 'compact'
  weekStart: { type: String, default: 'monday' }, // 'monday' | 'sunday'
  showLegend: { type: Boolean, default: false },
  categoryColors: { type: Array, default: null }, // override CALENDAR_CATEGORIES bila perlu
})
defineEmits(['day-click'])

const DEFAULT_COLORS = {
  libur: { fill: '#fde3e8', text: '#c10015', chip: '#c10015' },
  kbm: { fill: '#e3f2fd', text: '#1a4a7a', chip: '#1a4a7a' },
  ujian: { fill: '#ede7f6', text: '#5e35b1', chip: '#5e35b1' },
  kegiatan: { fill: '#e0f2f1', text: '#00695c', chip: '#00695c' },
  rapor: { fill: '#e8f5e9', text: '#2e7d32', chip: '#21ba45' },
}

function colorFor(kategori) {
  if (props.categoryColors) {
    const found = props.categoryColors.find((c) => c.value === kategori)
    if (found) return { fill: found.fill, text: found.text, chip: found.chip }
  }
  return DEFAULT_COLORS[kategori] || { fill: '#f5f5f5', text: '#616161', chip: '#9e9e9e' }
}

const today = todayISO()

const monthLabel = computed(() => `${MONTHS_ID[props.month]} ${props.year}`)
const weekdayLabels = computed(() =>
  props.weekStart === 'sunday' ? WEEKDAYS_SUNDAY_FIRST : WEEKDAYS_MONDAY_FIRST,
)

const cells = computed(() => buildMonthGrid(props.year, props.month, props.events, props.weekStart))

// Mode compact: sel diberi warna latar solid untuk event dengan fillStyle 'solid' (default)
function cellFillStyle(cell) {
  if (props.variant !== 'compact') return {}
  const fillEvent = cell.events.find((e) => (e.fillStyle || 'solid') === 'solid')
  if (fillEvent) return { backgroundColor: colorFor(fillEvent.kategori).fill }
  return {}
}

// Tanggal Minggu & libur nasional (fillStyle 'text-only') ditulis merah tebal, meniru kalender resmi
function dateTextStyle(cell) {
  const textOnlyEvent = cell.events.find((e) => e.fillStyle === 'text-only')
  if (textOnlyEvent) {
    return { color: colorFor(textOnlyEvent.kategori).text, fontWeight: 700 }
  }
  if (cell.jsWeekday === 0) {
    return { color: '#c10015' }
  }
  return {}
}

// Daftar kegiatan yang beririsan dengan bulan ini, dipakai sebagai keterangan di bawah grid (mode compact)
const legendItems = computed(() => {
  const [monthStart, monthEnd] = monthBoundsISO(props.year, props.month)
  return props.events
    .filter((e) => e.endDate >= monthStart && e.startDate <= monthEnd)
    .slice()
    .sort((a, b) => (a.startDate < b.startDate ? -1 : 1))
    .map((e) => ({
      id: e.id,
      dateLabel: formatRange(e.startDate, e.endDate),
      nama: e.keterangan ? `${e.nama} (${e.keterangan})` : e.nama,
    }))
})
</script>

<style scoped>
.mcg {
  width: 100%;
}

.mcg__header {
  text-align: center;
  font-weight: 700;
  color: #1a4a7a;
  margin-bottom: 6px;
}

.mcg--compact .mcg__header {
  font-size: 1rem;
  letter-spacing: 0.5px;
  text-transform: uppercase;
}

.mcg__weekdays {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  margin-bottom: 2px;
}

.mcg__weekday {
  text-align: center;
  font-size: 0.7rem;
  font-weight: 600;
  color: #607d8b;
  padding: 2px 0;
}

.mcg__weekday--sunday {
  color: #c10015;
}

.mcg--compact .mcg__weekday {
  background-color: #eceff1;
  border-radius: 3px;
  padding: 2px 0;
}

.mcg__grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 2px;
}

.mcg__cell {
  border: 1px solid #ebebeb;
  border-radius: 4px;
  min-height: 92px;
  padding: 4px;
  transition: background-color 0.15s ease;
}

.mcg__cell--clickable {
  cursor: pointer;
}

.mcg__cell--clickable:hover {
  background-color: #f5f7fa;
}

.mcg__cell--muted {
  background-color: #fafafa;
}

.mcg__cell--muted .mcg__date {
  color: #bdbdbd;
}

.mcg__cell--today {
  border-color: #1a4a7a;
}

.mcg__date {
  font-size: 0.8rem;
  font-weight: 600;
  color: #37474f;
  width: 22px;
  height: 22px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
}

.mcg__cell--today .mcg__date {
  background-color: #1a4a7a;
  color: white;
}

.mcg__events {
  margin-top: 2px;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.mcg__event-text {
  width: 100%;
  display: -webkit-box;
  -webkit-box-orient: vertical;

  /* ⬇️ Menentukan jumlah maksimal baris sebelum teks dipotong menjadi '...' */
  -webkit-line-clamp: 2;

  line-clamp: 2;
  /* Standar modern untuk kompabilitas */
  overflow: hidden;
  text-overflow: ellipsis;
  /* Memunculkan efek '...' */
  word-break: break-word;
  /* Memaksa kata panjang patah ke bawah */
}

.mcg__event-chip {
  display: flex;
  align-items: center;
  max-width: 100%;
  /* Memastikan lebar chip tidak melebihi kolom parent */
  overflow: hidden;
  /* Menyembunyikan konten yang meluap keluar chip */
  font-size: 0.68rem;
  line-height: 1.3;
  padding: 1px 4px;
  border-radius: 3px;
  /* white-space: nowrap; */
  /* overflow: hidden; */
  /* text-overflow: ellipsis; */
  color: white;
}

.mcg__event-more {
  font-size: 0.65rem;
  color: #757575;
  padding-left: 4px;
}

/* ===== Mode compact (cetak) ===== */
.mcg--compact .mcg__cell {
  min-height: 30px;
  padding: 2px;
  border-color: #cfd8dc;
}

.mcg--compact .mcg__date {
  font-size: 0.72rem;
  width: 20px;
  height: 20px;
}

.mcg--compact .mcg__cell--today {
  border: 1.5px solid #1a4a7a;
}

.mcg__legend {
  margin-top: 8px;
  border-top: 1px dashed #cfd8dc;
  padding-top: 6px;
}

.mcg__legend-item {
  font-size: 0.68rem;
  line-height: 1.5;
  color: #37474f;
}

.mcg__legend-date {
  font-weight: 700;
  text-decoration: underline;
  margin-right: 4px;
}

/* Layar kecil: perkecil sel kalender interaktif supaya tetap muat */
@media (max-width: 599px) {
  .mcg--interactive .mcg__cell {
    min-height: 64px;
    padding: 2px;
  }

  .mcg--interactive .mcg__event-chip {
    font-size: 0.6rem;
  }

  .mcg--interactive .mcg__weekday {
    font-size: 0.65rem;
  }
}
</style>

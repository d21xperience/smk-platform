<template>
  <q-page class="calendar-page q-pa-md">
    <!-- HEADER -->
    <div class="row items-center justify-between q-mb-md wrap q-gutter-sm">
      <div>
        <div class="text-h6 text-weight-bold text-grey-9">Kalender Akademik</div>
        <div class="text-caption text-grey-6">Tahun Ajaran {{ academicYear }}</div>
      </div>
      <div class="row items-center q-gutter-sm">
        <q-btn-toggle
          v-model="viewMode"
          no-caps
          unelevated
          toggle-color="primary"
          color="white"
          text-color="grey-8"
          :options="[
            { label: 'Kalender', value: 'grid', icon: 'calendar_month' },
            { label: 'Daftar', value: 'list', icon: 'view_list' }
          ]"
          class="view-toggle"
        />
        <q-btn no-caps outline color="primary" icon="download" label="Unduh" @click="handleDownloadCalendar" />
        <q-btn
          no-caps
          flat
          color="primary"
          icon="print"
          label="Kalender Cetak"
          :to="{ name: 'kalender-pendidikan-cetak' }"
        />
      </div>
    </div>

    <!-- TAB SEMESTER -->
    <q-card flat bordered class="q-mb-md">
      <q-tabs
        v-model="activeSemester"
        class="text-grey-7 semester-tabs"
        active-color="primary"
        indicator-color="primary"
        no-caps
        align="justify"
      >
        <q-tab name="1" :label="`Semester 1 (${semesterData['1'].label})`" />
        <q-tab name="2" :label="`Semester 2 (${semesterData['2'].label})`" />
      </q-tabs>
    </q-card>

    <q-tab-panels v-model="activeSemester" animated class="bg-transparent">
      <q-tab-panel v-for="sem in ['1', '2']" :key="sem" :name="sem" class="q-pa-none">
        <!-- RINGKASAN SEMESTER -->
        <div class="row q-col-gutter-md q-mb-md">
          <div class="col-6 col-md-3" v-for="stat in summaryStats(sem)" :key="stat.label">
            <q-card flat bordered class="stat-card">
              <q-card-section class="row items-center no-wrap">
                <div class="stat-icon q-mr-md" :class="`bg-${stat.color}-1`">
                  <q-icon :name="stat.icon" :color="stat.color" size="22px" />
                </div>
                <div>
                  <div class="text-h6 text-weight-bold">{{ stat.value }}</div>
                  <div class="text-caption text-grey-7">{{ stat.label }}</div>
                </div>
              </q-card-section>
            </q-card>
          </div>
        </div>

        <div class="row q-col-gutter-md">
          <!-- KONTEN UTAMA: KALENDER / DAFTAR -->
          <div class="col-12 col-md-8">
            <q-card flat bordered>
              <q-card-section>
                <!-- FILTER (berlaku untuk kedua mode tampilan) -->
                <div class="row q-col-gutter-sm q-mb-md">
                  <div class="col-12 col-sm-6">
                    <q-input dense outlined v-model="searchQuery" placeholder="Cari kegiatan..." clearable>
                      <template v-slot:prepend>
                        <q-icon name="search" />
                      </template>
                    </q-input>
                  </div>
                  <div class="col-12 col-sm-6">
                    <div class="row q-gutter-xs items-center category-filter">
                      <q-chip
                        v-for="cat in categories"
                        :key="cat.value"
                        clickable
                        square
                        :outline="filterCategory !== cat.value"
                        :color="cat.color"
                        :text-color="filterCategory === cat.value ? 'white' : cat.color"
                        class="q-px-sm"
                        @click="toggleCategoryFilter(cat.value)"
                      >
                        {{ cat.label }}
                      </q-chip>
                    </div>
                  </div>
                </div>

                <!-- ===== TAMPILAN KALENDER (memakai komponen reusable) ===== -->
                <template v-if="viewMode === 'grid'">
                  <div class="row items-center justify-between q-mb-sm wrap">
                    <div class="row items-center q-gutter-xs">
                      <q-btn flat round dense icon="chevron_left" @click="changeMonth(sem, -1)" />
                      <div class="text-subtitle1 text-weight-bold month-label">{{ monthLabel(sem) }}</div>
                      <q-btn flat round dense icon="chevron_right" @click="changeMonth(sem, 1)" />
                    </div>
                    <q-btn no-caps flat dense color="primary" label="Hari Ini" @click="goToday(sem)" />
                  </div>

                  <MonthCalendarGrid
                    :year="currentMonth[sem].year"
                    :month="currentMonth[sem].month"
                    :events="filteredEvents(sem)"
                    variant="interactive"
                    week-start="monday"
                    @day-click="openDay"
                  />
                </template>

                <!-- ===== TAMPILAN DAFTAR (PER BULAN) ===== -->
                <template v-else>
                  <div v-if="groupedEvents(sem).length">
                    <div v-for="group in groupedEvents(sem)" :key="group.month" class="month-group q-mb-md">
                      <div class="month-group__label text-weight-bold text-primary q-mb-xs">
                        {{ group.month }}
                      </div>
                      <q-list separator bordered class="rounded-borders">
                        <q-item v-for="event in group.events" :key="event.id">
                          <q-item-section avatar top>
                            <div class="date-chip" :class="`date-chip--${event.kategori}`">
                              <div class="date-chip__day">{{ event.tanggalAwal }}</div>
                              <div class="date-chip__month">{{ event.bulanSingkat }}</div>
                            </div>
                          </q-item-section>

                          <q-item-section>
                            <q-item-label class="text-weight-medium">{{ event.nama }}</q-item-label>
                            <q-item-label caption>
                              {{ event.tanggalLengkap }}
                              <span v-if="event.keterangan"> &middot; {{ event.keterangan }}</span>
                            </q-item-label>
                          </q-item-section>

                          <q-item-section side top>
                            <q-badge :color="categoryColor(event.kategori)" outline class="q-px-sm">
                              {{ categoryLabel(event.kategori) }}
                            </q-badge>
                          </q-item-section>
                        </q-item>
                      </q-list>
                    </div>
                  </div>

                  <div v-else class="text-center text-grey-6 q-pa-xl">
                    <q-icon name="event_busy" size="40px" class="q-mb-sm" />
                    <div>Tidak ada kegiatan yang cocok dengan filter</div>
                  </div>
                </template>
              </q-card-section>
            </q-card>
          </div>

          <!-- SIDEBAR: LEGENDA & KEGIATAN TERDEKAT -->
          <div class="col-12 col-md-4">
            <q-card flat bordered class="q-mb-md">
              <q-card-section>
                <div class="text-subtitle1 text-weight-bold q-mb-sm">Keterangan</div>
                <div class="row q-gutter-sm">
                  <div v-for="cat in categories" :key="cat.value" class="row items-center legend-item">
                    <div class="legend-dot q-mr-xs" :class="`bg-${cat.color}`" />
                    <span class="text-caption text-grey-8">{{ cat.label }}</span>
                  </div>
                </div>
              </q-card-section>
            </q-card>

            <q-card flat bordered>
              <q-card-section>
                <div class="text-subtitle1 text-weight-bold q-mb-sm">Kegiatan Terdekat</div>
                <q-list v-if="upcomingEvents(sem).length" separator>
                  <q-item v-for="event in upcomingEvents(sem)" :key="event.id" class="q-px-none">
                    <q-item-section avatar>
                      <div class="date-chip date-chip--sm" :class="`date-chip--${event.kategori}`">
                        <div class="date-chip__day">{{ event.tanggalAwal }}</div>
                        <div class="date-chip__month">{{ event.bulanSingkat }}</div>
                      </div>
                    </q-item-section>
                    <q-item-section>
                      <q-item-label class="text-body2 text-weight-medium">{{ event.nama }}</q-item-label>
                      <q-item-label caption>{{ event.tanggalLengkap }}</q-item-label>
                    </q-item-section>
                  </q-item>
                </q-list>
                <div v-else class="text-caption text-grey-6">
                  Tidak ada kegiatan mendatang pada semester ini.
                </div>
              </q-card-section>
            </q-card>
          </div>
        </div>
      </q-tab-panel>
    </q-tab-panels>

    <!-- DIALOG DETAIL HARI (saat sel kalender diklik) -->
    <q-dialog v-model="dayDialogOpen">
      <q-card style="min-width: 320px; max-width: 420px" class="full-width">
        <q-card-section class="row items-center justify-between">
          <div class="text-subtitle1 text-weight-bold">{{ selectedDayLabel }}</div>
          <q-btn flat round dense icon="close" v-close-popup />
        </q-card-section>
        <q-separator />
        <q-card-section v-if="selectedDayEvents.length" class="q-gutter-sm">
          <div v-for="event in selectedDayEvents" :key="event.id" class="row items-start no-wrap q-gutter-sm">
            <q-badge :color="categoryColor(event.kategori)" class="q-mt-xs">&nbsp;</q-badge>
            <div>
              <div class="text-weight-medium">{{ event.nama }}</div>
              <div class="text-caption text-grey-7">
                {{ event.tanggalLengkap }}
                <span v-if="event.keterangan"> &middot; {{ event.keterangan }}</span>
              </div>
            </div>
          </div>
        </q-card-section>
        <q-card-section v-else class="text-grey-6 text-center">
          Tidak ada kegiatan pada tanggal ini.
        </q-card-section>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useQuasar } from 'quasar'
import MonthCalendarGrid from '@/components/MonthCalendarGrid.vue'
import { MONTHS_ID, MONTHS_SHORT_ID, parseISO, formatRange, monthLabelOf, todayISO } from '@/composables/useCalendarDate'
import { academicYear, CALENDAR_CATEGORIES, semesterData, categoryLabel, categoryMeta } from '@/data/academicCalendarData'

const $q = useQuasar()

const activeSemester = ref('1')
const viewMode = ref('grid') // 'grid' | 'list'
const searchQuery = ref('')
const filterCategory = ref(null)

const dayDialogOpen = ref(false)
const selectedDayIso = ref(null)

const categories = CALENDAR_CATEGORIES
function categoryColor(value) {
  return categoryMeta(value).color
}

const today = todayISO()

// Bulan yang sedang ditampilkan pada kalender grid, per semester (0-based month index)
const currentMonth = reactive({
  '1': { year: 2026, month: 6 }, // Juli 2025
  '2': { year: 2027, month: 0 }  // Januari 2026
})

function toggleCategoryFilter(value) {
  filterCategory.value = filterCategory.value === value ? null : value
}

// Kegiatan hasil filter pencarian & kategori (dipakai baik oleh grid maupun daftar)
function filteredEvents(sem) {
  const query = (searchQuery.value || '').toLowerCase().trim()
  return semesterData[sem].events
    .filter(event => {
      const matchQuery = !query || event.nama.toLowerCase().includes(query)
      const matchCategory = !filterCategory.value || event.kategori === filterCategory.value
      return matchQuery && matchCategory
    })
    .slice()
    .sort((a, b) => (a.startDate < b.startDate ? -1 : 1))
}

// ---------- tampilan daftar (dikelompokkan per bulan) ----------
function groupedEvents(sem) {
  const events = filteredEvents(sem)
  const months = []
  events.forEach(event => {
    const label = monthLabelOf(event.startDate)
    let group = months.find(m => m.month === label)
    if (!group) {
      group = { month: label, events: [] }
      months.push(group)
    }
    group.events.push({
      ...event,
      tanggalAwal: parseISO(event.startDate).getDate(),
      bulanSingkat: MONTHS_SHORT_ID[parseISO(event.startDate).getMonth()],
      tanggalLengkap: formatRange(event.startDate, event.endDate)
    })
  })
  return months
}

// ---------- navigasi kalender grid ----------
function monthLabel(sem) {
  const { year, month } = currentMonth[sem]
  return `${MONTHS_ID[month]} ${year}`
}

function changeMonth(sem, delta) {
  const state = currentMonth[sem]
  let newMonth = state.month + delta
  let newYear = state.year
  if (newMonth < 0) { newMonth = 11; newYear -= 1 }
  if (newMonth > 11) { newMonth = 0; newYear += 1 }
  currentMonth[sem] = { year: newYear, month: newMonth }
}

function goToday(sem) {
  const now = new Date()
  currentMonth[sem] = { year: now.getFullYear(), month: now.getMonth() }
}

function openDay(cell) {
  selectedDayIso.value = cell.iso
  dayDialogOpen.value = true
}

const selectedDayLabel = computed(() => {
  if (!selectedDayIso.value) return ''
  const d = parseISO(selectedDayIso.value)
  const names = ['Minggu', 'Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat', 'Sabtu']
  return `${names[d.getDay()]}, ${d.getDate()} ${MONTHS_ID[d.getMonth()]} ${d.getFullYear()}`
})

const selectedDayEvents = computed(() => {
  if (!selectedDayIso.value) return []
  return filteredEvents(activeSemester.value)
    .filter(e => e.startDate <= selectedDayIso.value && selectedDayIso.value <= e.endDate)
    .map(e => ({ ...e, tanggalLengkap: formatRange(e.startDate, e.endDate) }))
})

// ---------- ringkasan & sidebar ----------
// TODO: ganti dengan perhitungan riil berdasarkan hari kalender & hari libur dari backend
function summaryStats(sem) {
  const events = semesterData[sem].events
  const totalUjian = events.filter(e => e.kategori === 'ujian').length
  const totalLibur = events.filter(e => e.kategori === 'libur').length
  const totalKegiatan = events.filter(e => e.kategori === 'kegiatan').length

  return [
    { label: 'Total Kegiatan', value: events.length, icon: 'event', color: 'primary' },
    { label: 'Ujian/Penilaian', value: totalUjian, icon: 'assignment', color: 'orange' },
    { label: 'Periode Libur', value: totalLibur, icon: 'beach_access', color: 'red' },
    { label: 'Kegiatan Sekolah', value: totalKegiatan, icon: 'celebration', color: 'teal' }
  ]
}

function upcomingEvents(sem) {
  return groupedEvents(sem)
    .flatMap(g => g.events)
    .filter(e => e.endDate >= today)
    .slice(0, 3)
}

function handleDownloadCalendar() {
  // TODO: ganti dengan permintaan unduh kalender akademik (PDF) ke Go backend,
  // misalnya GET /api/kalender-akademik/export?semester=1|2&format=pdf
  $q.notify({
    type: 'info',
    message: 'Fitur unduh kalender akademik akan segera tersedia.',
    position: 'top',
    timeout: 2000
  })
}
</script>

<style scoped>
.calendar-page {
  max-width: 1300px;
  margin: 0 auto;
}

.view-toggle {
  border: 1px solid #e0e0e0;
  border-radius: 8px;
}

.semester-tabs {
  border-radius: 8px;
}

.stat-card {
  border-radius: 8px;
  height: 100%;
}

.stat-icon {
  width: 44px;
  height: 44px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.month-label {
  min-width: 140px;
  text-align: center;
}

.month-group__label {
  padding-left: 4px;
}

.date-chip {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background-color: #eceff1;
  color: #37474f;
}

.date-chip--sm {
  width: 40px;
  height: 40px;
}

.date-chip__day {
  font-weight: 700;
  font-size: 1rem;
  line-height: 1;
}

.date-chip__month {
  font-size: 0.65rem;
  text-transform: uppercase;
  line-height: 1.2;
}

.date-chip--libur { background-color: #ffebee; color: #c10015; }
.date-chip--kbm { background-color: #e3f2fd; color: #1a4a7a; }
.date-chip--ujian { background-color: #fff3e0; color: #e65100; }
.date-chip--kegiatan { background-color: #e0f2f1; color: #00695c; }
.date-chip--rapor { background-color: #e8f5e9; color: #21ba45; }

.legend-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
}

.legend-item {
  padding: 2px 0;
}

.category-filter {
  flex-wrap: wrap;
}
</style>

// src/composables/useCalendarDate.js
// Utilitas tanggal bersama untuk fitur kalender akademik.
// Sengaja tanpa library eksternal (date-fns/moment) - hanya Date bawaan JS.
// Semua tanggal disimpan/diteruskan dalam format ISO 'YYYY-MM-DD'.

export const MONTHS_ID = [
  'Januari',
  'Februari',
  'Maret',
  'April',
  'Mei',
  'Juni',
  'Juli',
  'Agustus',
  'September',
  'Oktober',
  'November',
  'Desember',
]

export const MONTHS_SHORT_ID = [
  'Jan',
  'Feb',
  'Mar',
  'Apr',
  'Mei',
  'Jun',
  'Jul',
  'Agu',
  'Sep',
  'Okt',
  'Nov',
  'Des',
]

export const WEEKDAYS_MONDAY_FIRST = [
  'Senin',
  'Selasa',
  'Rabu',
  'Kamis',
  'Jumat',
  'Sabtu',
  'Minggu',
]
export const WEEKDAYS_SUNDAY_FIRST = [
  'Minggu',
  'Senin',
  'Selasa',
  'Rabu',
  'Kamis',
  'Jumat',
  'Sabtu',
]

/** Parse 'YYYY-MM-DD' menjadi objek Date lokal (hindari pergeseran timezone). */
export function parseISO(iso) {
  const [y, m, d] = iso.split('-').map(Number)
  return new Date(y, m - 1, d)
}

/** Ubah objek Date menjadi string ISO 'YYYY-MM-DD'. */
export function toISO(dateObj) {
  const y = dateObj.getFullYear()
  const m = String(dateObj.getMonth() + 1).padStart(2, '0')
  const d = String(dateObj.getDate()).padStart(2, '0')
  return `${y}-${m}-${d}`
}

export function todayISO() {
  return toISO(new Date())
}

export function formatFullDate(iso) {
  const d = parseISO(iso)
  return `${d.getDate()} ${MONTHS_ID[d.getMonth()]} ${d.getFullYear()}`
}

/** Format rentang tanggal secara ringkas, mis. "4 - 7 Agustus 2025" atau "14 Juli 2025" jika satu hari. */
export function formatRange(startISO, endISO) {
  if (startISO === endISO) return formatFullDate(startISO)
  const start = parseISO(startISO)
  const end = parseISO(endISO)
  if (start.getMonth() === end.getMonth() && start.getFullYear() === end.getFullYear()) {
    return `${start.getDate()} - ${end.getDate()} ${MONTHS_ID[start.getMonth()]} ${start.getFullYear()}`
  }
  return `${formatFullDate(startISO)} - ${formatFullDate(endISO)}`
}

export function monthLabelOf(iso) {
  const d = parseISO(iso)
  return `${MONTHS_ID[d.getMonth()]} ${d.getFullYear()}`
}

/** Rentang ISO awal & akhir dari sebuah bulan kalender (month 0-based). */
export function monthBoundsISO(year, month) {
  const start = toISO(new Date(year, month, 1))
  const end = toISO(new Date(year, month + 1, 0))
  return [start, end]
}

/**
 * Membangun grid 42 sel (6 minggu x 7 hari) untuk sebuah bulan.
 * weekStart: 'monday' | 'sunday'
 * events: [{ id, startDate, endDate, ... }]
 * Mengembalikan array sel: { iso, date, inCurrentMonth, jsWeekday, events }
 */
export function buildMonthGrid(year, month, events = [], weekStart = 'monday') {
  const firstOfMonth = new Date(year, month, 1)
  const jsDay = firstOfMonth.getDay() // 0=Minggu ... 6=Sabtu
  const startOffset = weekStart === 'sunday' ? jsDay : (jsDay + 6) % 7
  const gridStart = new Date(year, month, 1 - startOffset)

  const cells = []
  for (let i = 0; i < 42; i++) {
    const d = new Date(gridStart.getFullYear(), gridStart.getMonth(), gridStart.getDate() + i)
    const iso = toISO(d)
    const dayEvents = events.filter((e) => e.startDate <= iso && iso <= e.endDate)
    cells.push({
      iso,
      date: d.getDate(),
      inCurrentMonth: d.getMonth() === month,
      jsWeekday: d.getDay(),
      events: dayEvents,
    })
  }
  return cells
}

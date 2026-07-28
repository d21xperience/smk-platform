// src/data/academicCalendarData.js
// Sumber data tunggal untuk kalender akademik, dipakai oleh KalenderAkademik.vue (tampilan interaktif)
// maupun KalenderPendidikanCetak.vue (tampilan cetak ala kalender pendidikan resmi).
//
// TODO: ganti seluruh isi berkas ini dengan hasil GET /api/kalender-akademik?tahun_ajaran=... (Go backend).
// Field tiap kegiatan:
//   startDate, endDate : ISO 'YYYY-MM-DD'
//   kategori           : salah satu key pada CALENDAR_CATEGORIES di bawah
//   nama               : judul kegiatan
//   keterangan         : detail tambahan (opsional)
//   fillStyle          : 'solid' (default, sel diberi warna latar) atau 'text-only'
//                         (dipakai untuk libur nasional satu hari - tanggal cukup ditebalkan
//                         warna merah tanpa blok latar, meniru gaya kalender pendidikan resmi)

export const academicYear = '2026/2027'

// Kategori kegiatan & warna terkait. Dipakai bersama oleh MonthCalendarGrid, daftar kegiatan, dan legenda.
export const CALENDAR_CATEGORIES = [
  {
    value: 'libur',
    label: 'Libur',
    color: 'red',
    fill: '#fde3e8',
    text: '#c10015',
    chip: '#c10015',
  },
  {
    value: 'kbm',
    label: 'KBM Efektif',
    color: 'primary',
    fill: '#e3f2fd',
    text: '#1a4a7a',
    chip: '#1a4a7a',
  },
  {
    value: 'ujian',
    label: 'Ujian/Asesmen',
    color: 'deep-purple',
    fill: '#ede7f6',
    text: '#5e35b1',
    chip: '#5e35b1',
  },
  {
    value: 'kegiatan',
    label: 'Kegiatan Sekolah',
    color: 'teal',
    fill: '#e0f2f1',
    text: '#00695c',
    chip: '#00695c',
  },
  {
    value: 'rapor',
    label: 'Pembagian Rapor',
    color: 'positive',
    fill: '#e8f5e9',
    text: '#2e7d32',
    chip: '#21ba45',
  },
]

export const semesterData = {
  1: {
    label: 'Ganjil',
    periode: 'Juli - Desember 2026',
    events: [
      {
        id: 's1-1',
        startDate: '2026-07-15',
        endDate: '2026-07-17',
        nama: 'Masa Pengenalan Lingkungan Sekolah (MPLS)',
        kategori: 'kegiatan',
        keterangan: 'Khusus siswa baru kelas X',
      },
      {
        id: 's1-2',
        startDate: '2026-07-20',
        endDate: '2026-07-20',
        nama: 'Hari Pertama KBM Efektif',
        kategori: 'kbm',
        keterangan: '',
      },
      {
        id: 's1-3',
        startDate: '2026-08-04',
        endDate: '2026-08-07',
        nama: 'Pelaksanaan Asesmen Nasional SMA/SMK/SMALB dan yang Sederajat',
        kategori: 'ujian',
        keterangan: '',
      },
      {
        id: 's1-4',
        startDate: '2026-08-09',
        endDate: '2026-08-10',
        nama: 'Pelaksanaan AN Paket C/PKPPS Ulya',
        kategori: 'ujian',
        keterangan: '',
      },
      {
        id: 's1-5',
        startDate: '2026-08-17',
        endDate: '2026-08-17',
        nama: 'Libur HUT Kemerdekaan RI',
        kategori: 'libur',
        keterangan: 'Libur nasional',
        fillStyle: 'text-only',
      },
      {
        id: 's1-6',
        startDate: '2026-08-25',
        endDate: '2026-08-25',
        nama: 'Libur Maulid Nabi Muhammad SAW',
        kategori: 'libur',
        keterangan: 'Libur nasional',
        fillStyle: 'text-only',
      },
      {
        id: 's1-7',
        startDate: '2026-09-09',
        endDate: '2026-09-12',
        nama: 'Pelaksanaan Asesmen Nasional Paket B/PKPPS Wustha',
        kategori: 'ujian',
        keterangan: '',
      },
      {
        id: 's1-8',
        startDate: '2026-09-15',
        endDate: '2026-09-20',
        nama: 'Penilaian Tengah Semester (PTS) Ganjil',
        kategori: 'ujian',
        keterangan: '',
      },
      {
        id: 's1-9',
        startDate: '2026-12-01',
        endDate: '2026-12-10',
        nama: 'Penilaian Akhir Semester (PAS) Ganjil',
        kategori: 'ujian',
        keterangan: '',
      },
      {
        id: 's1-10',
        startDate: '2026-12-19',
        endDate: '2026-12-19',
        nama: 'Pembagian Rapor Semester Ganjil',
        kategori: 'rapor',
        keterangan: '',
      },
      {
        id: 's1-11',
        startDate: '2026-12-22',
        endDate: '2027-01-03',
        nama: 'Libur Akhir Semester Ganjil',
        kategori: 'libur',
        keterangan: '',
      },
    ],
  },
  2: {
    label: 'Genap',
    periode: 'Januari - Juni 2027',
    events: [
      {
        id: 's2-1',
        startDate: '2027-01-05',
        endDate: '2027-01-05',
        nama: 'Hari Pertama KBM Semester Genap',
        kategori: 'kbm',
        keterangan: '',
      },
      {
        id: 's2-2',
        startDate: '2027-03-09',
        endDate: '2027-03-14',
        nama: 'Penilaian Tengah Semester (PTS) Genap',
        kategori: 'ujian',
        keterangan: '',
      },
      {
        id: 's2-3',
        startDate: '2027-04-20',
        endDate: '2027-04-30',
        nama: 'Ujian Kompetensi Keahlian (UKK) Kelas XII',
        kategori: 'ujian',
        keterangan: 'Wajib bagi kelas XII',
      },
      {
        id: 's2-4',
        startDate: '2027-06-01',
        endDate: '2027-06-10',
        nama: 'Penilaian Akhir Tahun (PAT)',
        kategori: 'ujian',
        keterangan: '',
      },
      {
        id: 's2-5',
        startDate: '2027-06-19',
        endDate: '2027-06-19',
        nama: 'Pembagian Rapor Semester Genap & Kenaikan Kelas',
        kategori: 'rapor',
        keterangan: '',
      },
      {
        id: 's2-6',
        startDate: '2027-06-22',
        endDate: '2027-07-12',
        nama: 'Libur Akhir Tahun Ajaran',
        kategori: 'libur',
        keterangan: '',
      },
    ],
  },
}

/** Gabungan semua kegiatan semester 1 & 2, dipakai halaman cetak yang menampilkan satu tahun penuh. */
export function getAllEvents() {
  return [...semesterData['1'].events, ...semesterData['2'].events]
}

export function categoryLabel(value) {
  return CALENDAR_CATEGORIES.find((c) => c.value === value)?.label || value
}
export function categoryMeta(value) {
  return CALENDAR_CATEGORIES.find((c) => c.value === value) || CALENDAR_CATEGORIES[0]
}

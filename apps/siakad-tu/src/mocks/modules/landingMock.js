export function registerLandingMocks(mock) {
  // 1. GET /dashboard/kpi-metrics
  mock.onGet('/dashboard/kpi-metrics').reply(200, {
    meta: { status: 'success', code: 200 },
    data: [
      { label: 'Total Siswa', value: '1.248', icon: 'groups', color: 'primary', trend: 3.2 },
      { label: 'Total Guru & Staf', value: '86', icon: 'badge', color: 'secondary', trend: 1.1 },
      { label: 'Total Rombel', value: '42', icon: 'meeting_room', color: 'teal', trend: 0 },
      {
        label: 'Kompetensi Keahlian',
        value: '6',
        icon: 'school',
        color: 'deep-orange',
        trend: null,
      },
    ],
  })

  // 2. GET /dashboard/student-distribution
  mock.onGet('/dashboard/student-distribution').reply(200, {
    meta: { status: 'success', code: 200 },
    data: {
      total_siswa: 1248,
      distribusi: [
        { nama: 'Teknik & Bisnis Sepeda Motor (TBSM)', jumlah: 412, warna: '#0a192f' },
        { nama: 'Teknik Komputer & Jaringan', jumlah: 356, warna: '#c9a227' },
        { nama: 'Otomatisasi Tata Kelola Perkantoran', jumlah: 132, warna: '#5e35b1' },
        { nama: 'Akuntansi & Keuangan Lembaga', jumlah: 248, warna: '#00796b' },
        { nama: 'Multimedia', jumlah: 100, warna: '#d84315' },
      ],
    },
  })

  // 3. GET /dashboard/recent-activities
  mock.onGet('/dashboard/recent-activities').reply(200, {
    data: [
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
    ],
  })
}

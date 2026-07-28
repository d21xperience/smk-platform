// src/mocks/siakadMock.js

import { dummySchedule, mockStudentsDB } from './data/schoolData'

export function setupSiakadMocks(mock) {
  // 1. MOCK: Ambil Jadwal Mengajar Guru
  mock.onGet('/siakad/attendance/schedule').reply((config) => {
    const { date } = config.params
    console.log(`[MOCK] Fetching schedule for date: ${date}`)

    // Simulasi: Jika hari Minggu (misal), tidak ada jadwal
    const day = new Date(date).getDay()
    if (day === 0) {
      return [200, []]
    }

    // Data jadwal statis untuk demo
    return [
      200,
      [
        { class_id: 101, class_name: 'X RPL 1', hour_start: '07:00', hour_end: '08:30' },
        { class_id: 102, class_name: 'XI TKR 2', hour_start: '09:00', hour_end: '10:30' },
      ],
    ]
  })

  // 2. MOCK: Ambil Daftar Siswa per Kelas
  mock.onGet('/siakad/attendance/students').reply((config) => {
    const { class_id, date } = config.params
    console.log(`[MOCK] Fetching students for class_id: ${class_id}, date: ${date}`)

    const students = mockStudentsDB[class_id] || []

    // Simulasi: Berikan status acak untuk demo (sebagian H, sebagian S/I/A)
    // Di backend asli, ini akan diambil dari database absensi hari itu
    const studentsWithStatus = students.map((s) => ({
      ...s,
      status: s.id % 5 === 0 ? 'S' : s.id % 7 === 0 ? 'I' : 'H',
    }))

    return [200, studentsWithStatus]
  })

  // 3. MOCK: Submit Absensi
  mock.onPost('/siakad/attendance/submit').reply((config) => {
    const payload = JSON.parse(config.data)
    console.log('[MOCK] Submitting attendance payload:', payload)

    // Simulasi delay jaringan agar tombol loading terlihat
    return new Promise((resolve) => {
      setTimeout(() => {
        resolve([
          200,
          {
            message: 'Absensi berhasil disimpan ke sistem.',
            saved_count: payload.attendances.length,
          },
        ])
      }, 800) // Delay 800ms
    })
  })

  // Di dalam setupMockInterceptor(api)

  mock.onGet('/siakad/dashboard/overview').reply((config) => {
    const { periode } = config.params
    console.log(`[MOCK] Fetching Dashboard Overview | Periode: ${periode}`)

    return [
      200,
      {
        stats: [
          { label: 'Total Siswa Diampu', value: '186', icon: 'groups', color: 'primary' },
          { label: 'Kehadiran Hari Ini', value: '94%', icon: 'fact_check', color: 'positive' },
          { label: 'Jadwal Hari Ini', value: '4', icon: 'schedule', color: 'orange' },
          { label: 'Nilai Belum Diisi', value: '12', icon: 'assignment_late', color: 'red' },
        ],
        weeklyAttendance: [
          { day: 'Sen', hadir: 96, alpha: 4 },
          { day: 'Sel', hadir: 92, alpha: 8 },
          { day: 'Rab', hadir: 98, alpha: 2 },
          { day: 'Kam', hadir: 90, alpha: 10 },
          { day: 'Jum', hadir: 94, alpha: 6 },
        ],
        todaySchedule: [
          {
            id: 1,
            jam: '07:00',
            mapel: 'Pemrograman Web',
            kelas: 'XI RPL 1',
            ruang: 'Lab Komputer 2',
            status: 'selesai',
          },
          {
            id: 2,
            jam: '09:15',
            mapel: 'Basis Data',
            kelas: 'XI RPL 2',
            ruang: 'Lab Komputer 1',
            status: 'selesai',
          },
          {
            id: 3,
            jam: '13:00',
            mapel: 'Pemrograman Web',
            kelas: 'XII RPL 1',
            ruang: 'Lab Komputer 2',
            status: 'belum',
          },
        ],
        criticalAlerts: [
          {
            id: 1,
            icon: 'warning',
            title: 'Rafi Ardiansyah (XI RPL 1) alpha 3 hari',
            description: 'Perlu tindak lanjut wali kelas / BK',
          },
          {
            id: 2,
            icon: 'assignment_late',
            title: '12 siswa belum memiliki nilai UTS',
            description: 'Batas input nilai: 20 Juli 2026',
          },
        ],
      },
    ]
  })

  // Di dalam setupMockInterceptor(api)

  mock.onGet('/siakad/attendance/schedule').reply((config) => {
    const { periode, date } = config.params
    console.log(`[MOCK] GET Schedule | Periode: ${periode}, Date: ${date}`)
    return [
      200,
      [
        {
          class_id: 101,
          class_name: 'X RPL 1',
          subject: 'Pemrograman Dasar',
          start_time: '07:00',
          end_time: '08:30',
          room: 'Lab 1',
        },
        {
          class_id: 102,
          class_name: 'XI TKR 2',
          subject: 'Basis Data',
          start_time: '09:00',
          end_time: '10:30',
          room: 'Ruang 204',
        },
      ],
    ]
  })

  mock.onGet('/siakad/attendance/students').reply((config) => {
    const { periode, class_id, date } = config.params
    console.log(`[MOCK] GET Students | Periode: ${periode}, Class: ${class_id}, Date: ${date}`)

    const studentsDB = {
      101: [
        { id: 1001, nis: '2627001', name: 'Ahmad Fauzi' },
        { id: 1002, nis: '2627002', name: 'Siti Aminah' },
        { id: 1003, nis: '2627003', name: 'Budi Santoso' },
      ],
      102: [
        { id: 2001, nis: '2527001', name: 'Gilang Ramadhan' },
        { id: 2002, nis: '2527002', name: 'Hana Pertiwi' },
      ],
    }

    const students = studentsDB[class_id] || []
    // Simulasi status acak
    return [200, students.map((s, i) => ({ ...s, status: i === 1 ? 'S' : 'H' }))]
  })

  mock.onPost('/siakad/attendance/submit').reply((config) => {
    const payload = JSON.parse(config.data)
    console.log(`[MOCK] POST Attendance | Periode: ${payload.periode}`, payload)

    return new Promise((resolve) => {
      setTimeout(() => {
        resolve([200, { message: `Absensi kelas ${payload.class_id} berhasil disimpan.` }])
      }, 800)
    })
  })

  mock.onGet('/siakad/schedule/teacher').reply((config) => {
    const { periode } = config.params
    console.log(`[MOCK] Fetching teacher schedule | Periode: ${periode}`)

    return [200, dummySchedule]
  })

  // ==========================================
  // 1. MOCK: LOGIN SIAKAD (Email Based)
  // ==========================================
  mock.onPost('/siakad/auth/login').reply((config) => {
    // Parse payload dari request
    const { identifier, password, remember_me } = JSON.parse(config.data)

    console.log(`[MOCK] Login SIAKAD | Email: ${identifier}, Remember: ${remember_me}`)

    // Kredensial valid untuk testing
    const VALID_EMAIL = 'd21xperience@gmail.com'
    const VALID_PASSWORD = 'pass123'

    // Simulasi delay jaringan agar tombol loading terlihat realistis
    return new Promise((resolve) => {
      setTimeout(() => {
        if (identifier === VALID_EMAIL && password === VALID_PASSWORD) {
          // ✅ LOGIN BERHASIL
          resolve([
            200,
            {
              token: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.mock_siakad_jwt_token_12345',
              user: {
                id: 101,
                name: 'Deden Moh Jaenudin, S.H.',
                email: VALID_EMAIL,
                role: 'siakad_guru', // Opsi lain: 'siakad_wali_kelas', 'siakad_admin'
                avatar_url: null,
              },
            },
          ])
        } else {
          // ❌ LOGIN GAGAL
          resolve([
            401,
            {
              message: 'Email atau kata sandi yang Anda masukkan salah.',
            },
          ])
        }
      }, 800) // Delay 800ms
    })
  })
}

export const MOCK_TEACHERS = {
  // Guru 1: Budi (Mengajar RPL)
  'guru@smkpasjat.my.id': {
    id: 101,
    name: 'Budi Santoso, S.Pd., M.Kom.',
    email: 'guru@smkpasjat.my.id',
    role: 'siakad_guru',
    token: 'mock_jwt_token_budi_12345',
    schedule: [
      {
        day: 'Senin',
        class_id: 101,
        class_name: 'X RPL 1',
        subject: 'Pemrograman Dasar',
        start_time: '07:00',
        end_time: '08:30',
        room: 'Lab Komputer 1',
      },
      {
        day: 'Rabu',
        class_id: 102,
        class_name: 'XI RPL 2',
        subject: 'Basis Data',
        start_time: '09:00',
        end_time: '10:30',
        room: 'Lab Komputer 2',
      },
    ],
    students: {
      101: [
        // Siswa X RPL 1
        { id: 1001, nis: '2627001', name: 'Ahmad Fauzi' },
        { id: 1002, nis: '2627002', name: 'Siti Aminah' },
        { id: 1003, nis: '2627003', name: 'Rudi Hermawan' },
      ],
      102: [
        // Siswa XI RPL 2
        { id: 2001, nis: '2527010', name: 'Dewi Lestari' },
        { id: 2002, nis: '2527011', name: 'Eko Prasetyo' },
      ],
    },
    dashboardStats: [
      { label: 'Total Siswa Diampu', value: '5', icon: 'groups', color: 'primary' },
      { label: 'Kehadiran Hari Ini', value: '100%', icon: 'fact_check', color: 'positive' },
      { label: 'Jadwal Hari Ini', value: '2', icon: 'schedule', color: 'orange' },
      { label: 'Nilai Belum Diisi', value: '0', icon: 'assignment_late', color: 'red' },
    ],
  },

  // Guru 2: Siti (Wali Kelas, Mengajar TKR)
  'wali@smkpasjat.my.id': {
    id: 102,
    name: 'Siti Aminah, S.Pd.',
    email: 'wali@smkpasjat.my.id',
    role: 'siakad_wali_kelas',
    token: 'mock_jwt_token_siti_67890',
    schedule: [
      {
        day: 'Selasa',
        class_id: 201,
        class_name: 'X TKR 1',
        subject: 'Dasar Otomotif',
        start_time: '07:00',
        end_time: '09:00',
        room: 'Bengkel A',
      },
      {
        day: 'Kamis',
        class_id: 201,
        class_name: 'X TKR 1',
        subject: 'Wali Kelas',
        start_time: '09:00',
        end_time: '10:00',
        room: 'Ruang Kelas',
      },
    ],
    students: {
      201: [
        // Siswa X TKR 1
        { id: 3001, nis: '2628001', name: 'Gilang Ramadhan' },
        { id: 3002, nis: '2628002', name: 'Hana Pertiwi' },
        { id: 3003, nis: '2628003', name: 'Indra Lesmana' },
        { id: 3004, nis: '2628004', name: 'Joko Anwar' },
      ],
    },
    dashboardStats: [
      { label: 'Total Siswa Diampu', value: '4', icon: 'groups', color: 'primary' },
      { label: 'Kehadiran Hari Ini', value: '75%', icon: 'fact_check', color: 'warning' },
      { label: 'Jadwal Hari Ini', value: '1', icon: 'schedule', color: 'orange' },
      { label: 'Nilai Belum Diisi', value: '4', icon: 'assignment_late', color: 'red' },
    ],
  },
}

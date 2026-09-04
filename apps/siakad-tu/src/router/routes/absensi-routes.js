const absensiRoutes = {
  path: '/absensi',
  meta: {
    requiresAuth: true,
    allowedRoles: ['admin', 'tu', 'guru', 'kesiswaan'],
  },
  children: [
    {
      path: '',
      redirect: '/absensi/sesi',
    },
    {
      path: 'sesi',
      name: 'absensi-sesi',
      component: () => import('@/pages/absensi/SessionListPage.vue'),
    },
    {
      path: 'sesi/create',
      name: 'absensi-sesi-create',
      component: () => import('@/pages/absensi/SessionFormPage.vue'),
    },
    {
      path: 'sesi/:id/records',
      name: 'absensi-records',
      component: () => import('@/pages/absensi/AttendanceRecordPage.vue'),
      props: true,
    },
    {
      path: 'rekap/siswa/:studentId',
      name: 'absensi-rekap-siswa',
      component: () => import('@/pages/absensi/StudentAttendanceSummaryPage.vue'),
      props: true,
    },
    {
      path: 'rekap/kelas/:classId',
      name: 'absensi-rekap-kelas',
      component: () => import('@/pages/absensi/ClassAttendanceSummaryPage.vue'),
      props: true,
    },
  ],
}

export default absensiRoutes

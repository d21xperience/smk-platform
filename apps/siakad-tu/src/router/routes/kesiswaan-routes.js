export default {
  path: 'kesiswaan',
  children: [
    {
      path: '',
      redirect: '/kesiswaan/siswa'
    },
    {
      path: 'siswa',
      name: 'kesiswaan-siswa',
      component: () => import('@/pages/kesiswaan/siswa/StudentListPage.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: 'siswa/:id',
      name: 'kesiswaan-siswa-detail',
      component: () => import('@/pages/kesiswaan/siswa/StudentDetailPage.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: 'siswa/edit/:id',
      name: 'kesiswaan-siswa-edit',
      component: () => import('@/pages/kesiswaan/siswa/StudentEditPage.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: 'absensi-harian',
      name: 'kesiswaan-absensi-harian',
      component: () => import('@/pages/kesiswaan/absensi/AttendanceSessionPage.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: 'rekap-kehadiran',
      name: 'kesiswaan-rekap-kehadiran',
      component: () => import('@/pages/kesiswaan/absensi/AttendanceRecapPage.vue'),
      meta: { requiresAuth: true },
    },
  ],
}

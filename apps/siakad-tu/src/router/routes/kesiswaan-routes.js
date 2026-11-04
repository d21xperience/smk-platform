const kesiswaanRoutes = {
  path: '/kesiswaan',
  children: [
    {
      path: '',
      redirect: '/kesiswaan/dashboard',
    },
    {
      path: 'dashboard',
      name: 'kesiswaan-dashboard',
      component: () => import('@/pages/kesiswaan/DashboardPage.vue'),
    },
    // === RUTE BARU UNTUK MODUL KELAS ===
    {
      path: 'kelas',
      name: 'kesiswaan-kelas',
      component: () => import('@/pages/kesiswaan/DataKelas.vue'),
    },
    // ===================================
    {
      path: 'siswa',
      children: [
        {
          path: '',
          name: 'student-list',
          component: () => import('@/pages/kesiswaan/StudentListPage.vue'),
        },
        {
          path: 'create',
          name: 'student-create',
          component: () => import('@/pages/kesiswaan/StudentCreatePage.vue'),
        },
        {
          path: 'import',
          name: 'student-import',
          component: () => import('@/pages/kesiswaan/StudentImportPage.vue'),
        },
        {
          path: ':id',
          name: 'student-detail',
          component: () => import('@/pages/kesiswaan/StudentDetailPage.vue'),
        },
        {
          path: ':id/edit',
          name: 'student-edit',
          component: () => import('@/pages/kesiswaan/StudentEditPage.vue'),
        },
      ],
    },
    {
      path: 'mutasi-siswa',
      children: [
        {
          path: '',
          name: 'student-mutation',
          component: () => import('@/pages/kesiswaan/StudentMutationPage.vue'),
        },
      ],
    },
    {
      path: 'direktori-siswa',
      name: 'direktori-siswa',
      component: () => import('@/pages/kesiswaan/DirektoriSiswa.vue'),
    },
    {
      path: 'buku-induk-siswa',
      name: 'buku-induk-siswa',
      component: () => import('@/pages/kesiswaan/BukuIndukKlaper.vue'),
    },
    {
      path: 'beasiswa-pip-siswa',
      name: 'beasiswa-pip-siswa',
      component: () => import('@/pages/kesiswaan/BeasiswaPip.vue'),
    },
  ],
  meta: {
    requiresAuth: true,
    allowedRoles: ['admin', 'kesiswaan'],
  },
}

export default kesiswaanRoutes

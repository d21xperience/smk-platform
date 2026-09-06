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
    {
      path: 'siswa',
      children: [
        {
          path: '',
          name: 'student-list',
          component: () => import('@/pages/kesiswaan/StudentListPage.vue'),
          // component: () => import('@/pages/kesiswaan/DirektoriSiswa.vue'),
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
    // {
    //   path: 'psb-siswa-baru',
    //   name: 'psb-siswa-baru',
    //   component: () => import('@/pages/kesiswaan/SiswaBaruPage.vue'),
    // },
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

    // {
    //   path: 'surat-aktif-siswa',
    //   name: 'surat-aktif-siswa',
    //   component: () => import('@/pages/kesiswaan/SuratAktif.vue'),
    // },
    {
      path: 'beasiswa-pip-siswa',
      name: 'beasiswa-pip-siswa',
      component: () => import('@/pages/kesiswaan/BeasiswaPip.vue'),
    },
    // {
    //   path: 'legalisir-alumni',
    //   name: 'legalisir-alumni',
    //   component: () => import('@/pages/kesiswaan/LegalisirAlumni.vue'),
    // },

    // {
    //   path: 'rekap-pelaporan',
    //   name: 'rekap-pelaporan',
    //   component: () => import('@/pages/kesiswaan/RekapPelaporan.vue'),
    // },
    // // -------------------------
    // {
    //   path: 'data-siswa',
    //   name: 'data-siswa',
    //   component: () => import('@/pages/kesiswaan/KesiswaanSiswa.vue'),
    // },
    // {
    //   path: 'monitoring-absensi',
    //   name: 'monitoring-absensi',
    //   component: () => import('@/pages/kesiswaan/MonitoringAbsensi.vue'),
    // },
    // {
    //   path: 'pendaftaran',
    //   // component: () => import('@/pages/kesiswaan/MonitoringAbsensi.vue'),
    //   children: [
    //     {
    //       path: 'siswa-baru',
    //       name: 'pendaftaran-siswa-baru',
    //       component: () => import('@/pages/kesiswaan/PendaftaranSiswaBaru.vue'),
    //     },
    // {
    //   path: 'manajemen-mutasi-siswa',
    //   name: 'manajemen-mutasi-siswa',
    //   component: () => import('@/pages/kesiswaan/MutasiPage.vue'),
    // },
    //   ],
    // },
  ],
  meta: {
    requiresAuth: true,
    allowedRoles: ['admin', 'kesiswaan'], // <-- Definisikan role yang diizinkan di sini
  },
}

export default kesiswaanRoutes

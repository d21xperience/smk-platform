const kepegawaian = {
  path: '/kepegawaian',
  children: [
    { path: '', redirect: '/kepegawaian/dashboard' },
    {
      path: 'dashboard',
      name: 'dashboard-kepegawaian',
      component: () => import('@/pages/kepegawaian/DashboardKepegawaian.vue'),
    },
    {
      path: 'statistik',
      name: 'statistik-keberadaan',
      component: () => import('@/pages/kepegawaian/StatistikKeberadaanPage.vue'),
    },
    {
      path: 'direktori-gtk',
      name: 'direktori-gtk',
      component: () => import('@/pages/kepegawaian/DataGuru.vue'),
    },
    {
      path: 'buku-induk',
      name: 'buku-induk',
      component: () => import('@/pages/kepegawaian/BukuIndukPage.vue'),
    },
    {
      path: 'buku-induk/:id',
      name: 'buku-induk-detail',
      component: () => import('@/pages/kepegawaian/BukuIndukDetailPage.vue'),
    },
    {
      path: 'sk-pembagian-tugas',
      name: 'sk-pembagian-tugas',
      component: () => import('@/pages/kepegawaian/SKPembagianTugasPage.vue'),
    },
    {
      path: 'sk-pembagian-tugas/create',
      name: 'sk-pembagian-tugas-create',
      component: () => import('@/pages/kepegawaian/SKWizardPage.vue'),
    },
    {
      path: 'skp',
      name: 'skp-periode',
      component: () => import('@/pages/kepegawaian/SKPPeriodePage.vue'),
    },
    {
      path: 'skp/periode/:id',
      name: 'skp-penilaian',
      component: () => import('@/pages/kepegawaian/SKPPenilaianPage.vue'),
    },
    {
      path: 'skp/penilaian/:id',
      name: 'skp-form',
      component: () => import('@/pages/kepegawaian/SKPFormPage.vue'),
    },
    {
      path: 'buku-kontrol',
      name: 'buku-kontrol',
      component: () => import('@/pages/kepegawaian/BukuKontrolPage.vue'),
    },
    {
      path: 'usulan-kenaikan',
      name: 'usulan-kenaikan',
      component: () => import('@/pages/kepegawaian/UsulanKenaikanPage.vue'),
    },
    {
      path: 'kehadiran',
      name: 'kehadiran-pegawai',
      component: () => import('@/pages/kepegawaian/KehadiranPage.vue'),
    },
    {
      path: 'beban-mengajar',
      name: 'beban-mengajar',
      component: () => import('@/pages/kepegawaian/BebanMengajar.vue'),
    },
    {
      path: 'cuti-kerja',
      name: 'cuti-kerja',
      component: () => import('@/pages/kepegawaian/CutiMengajar.vue'),
    },
    {
      path: 'arsip-kepegawaian',
      name: 'arsip-kepegawaian',
      component: () => import('@/pages/kepegawaian/ArsipDokumen.vue'),
    },
  ],
  meta: {
    requiresAuth: true,
    allowedRoles: ['admin', 'kepegawaian'],
  },
}

export default kepegawaian

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
      path: 'kinerja-gtk',
      name: 'kinerja-gtk',
      component: () => import('@/pages/kepegawaian/KinerjaGTK.vue'),
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

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
    allowedRoles: ['admin', 'kepegawaian'], // <-- Definisikan role yang diizinkan di sini
  },
}

export default kepegawaian

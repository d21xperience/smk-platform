const kepegawaian = {
  path: '/kepegawaian',
  children: [
    { path: '', redirect: '/kepegawaian/data-guru' },
    {
      path: 'data-guru',
      component: () => import('@/pages/kepegawaian/DataGuru.vue'),
      name: 'main',
    },
    {
      path: 'beban-mengajar',
      component: () => import('@/pages/kepegawaian/BebanMengajar.vue'),
      name: 'bebanMengajar',
    },
    {
      path: 'cuti-kerja',
      component: () => import('@/pages/kepegawaian/CutiMengajar.vue'),
      name: 'cuti',
    },
    {
      path: 'arsip',
      component: () => import('@/pages/kepegawaian/ArsipDokumen.vue'),
      name: 'cuti',
    },
  ],
  meta: {
    requiresAuth: true,
    allowedRoles: ['admin', 'kepegawaian'], // <-- Definisikan role yang diizinkan di sini
  },
}

export default kepegawaian

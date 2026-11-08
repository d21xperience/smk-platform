const kurikulumRoutes = {
  path: '/kurikulum',
  children: [
    {
      path: '',
      redirect: '/kurikulum/mapel',
    },
    {
      path: 'mapel',
      name: 'kurikulum-mapel',
      component: () => import('@/pages/kurikulum/DataMapel.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: 'kelas',
      name: 'kurikulum-kelas',
      component: () => import('@/pages/kurikulum/DataKelas.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: 'pembelajaran',
      name: 'kurikulum-pembelajaran',
      component: () => import('@/pages/kurikulum/DataPembelajaran.vue'),
      meta: { requiresAuth: true },
    },
  ],
}

export default kurikulumRoutes

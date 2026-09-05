const previewDev = {
  path: '/preview',
  children: [
    {
      path: '',
      name: 'landing-home',
      component: () => import('@/pages/landing/LandingIndex.vue'),
    },
    {
      path: 'kesiswaan',
      name: 'landing-kesiswaan',
      component: () => import('@/pages/landing/LandingKesiswaan.vue'),
    },
    {
      path: 'sarpras',
      name: 'landing-sarpras',
      component: () => import('@/pages/landing/LandingSarpras.vue'),
    },
    {
      path: 'kurikulum',
      name: 'landing-kurikulum',

      component: () => import('@/pages/landing/LandingKurikulum.vue'),
    },
    {
      path: 'hubin',
      name: 'landing-hubin',
      component: () => import('@/pages/landing/LandingHubin.vue'),
    },
  ],
}

export default previewDev

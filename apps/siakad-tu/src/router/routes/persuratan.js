const persuratan =
  // SURAT - MENYURAT
  {
    path: '/persuratan',
    children: [
      {
        path: '',
        redirect: '/persuratan/dashboard',
      },
      {
        path: 'dashboard',
        name: 'surat-menyurat-dashboard',
        component: () => import('@/pages/surat-menyurat/SuratMenyuratDashboardPage.vue'),
      },
      {
        path: 'surat-keluar',
        name: 'surat-keluar',
        component: () => import('@/pages/surat-menyurat/SuratMenyuratSuratKeluarPage.vue'),
      },
      {
        path: 'surat-masuk',
        name: 'surat-masuk',
        component: () => import('@/pages/surat-menyurat/SuratMenyuratSuratMasukPage.vue'),
      },
      {
        path: 'disposisi',
        name: 'disposisi',
        component: () => import('@/pages/surat-menyurat/SuratMenyuratDisposisiPage.vue'),
      },
      {
        path: 'klasifikasi',
        name: 'klasifikasi',
        component: () => import('@/pages/surat-menyurat/SuratMenyuratKlasifikasiArsipPage.vue'),
      },
    ],
  }

export default persuratan

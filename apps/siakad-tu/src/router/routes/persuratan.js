const persuratan = {
  path: '/persuratan',
  children: [
    { path: '', redirect: '/persuratan/dashboard' },
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
    {
      path: 'master-jenis-surat',
      name: 'master-jenis-surat',
      component: () => import('@/pages/surat-menyurat/MasterJenisSuratPage.vue'),
      meta: { title: 'Master Jenis Surat' },
    },
    {
      path: 'buku-agenda',
      name: 'buku-agenda',
      component: () => import('@/pages/surat-menyurat/BukuAgendaPage.vue'),
      meta: { title: 'Buku Agenda' },
    },
  ],
}

export default persuratan

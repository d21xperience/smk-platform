// import hubin from './routes/hubin.js'
// import kepegawaian from './routes/kepegawaian.js'
// import kesiswaanRoutes from './routes/kesiswaan-routes.js'
import kesiswaanRoutes from './routes/kesiswaan-routes.js'
import persuratan from './routes/persuratan.js'
// import persuratan from './routes/persuratan.js'
// import psb from './routes/psb.js'
// import sarpras from './routes/sarpras.js'

const routes = [
  {
    path: '/',
    component: () => import('@/layouts/LandingLayout.vue'),
    children: [
      { path: '', name: 'landing-home', component: () => import('@/pages/landing/IndexPage.vue') },
      {
        path: 'preview-kurikulum',
        name: 'landing-kurikulum',
        component: () => import('@/pages/landing/KurikulumPage.vue'),
      },
      {
        path: 'preview-kesiswaan',
        name: 'landing-kesiswaan',
        component: () => import('@/pages/landing/LandingKesiswaan.vue'),
      },
      {
        path: 'preview-sarpras',
        name: 'landing-sarpras',
        component: () => import('@/pages/landing/LandingSarpras.vue'),
      },
      {
        path: 'preview-hubin',
        name: 'landing-hubin',
        component: () => import('@/pages/landing/LandingHubin.vue'),
      },
      {
        path: 'preview-psb',
        name: 'landing-psb',
        component: () => import('@/pages/landing/LandingPsb.vue'),
      },
      // kepegawaian,
      persuratan,

      kesiswaanRoutes,
      // hubin,
      // sarpras,
      // psb,
    ],
  },

  {
    path: '/kohor',
    // name: 'landing-kesiswaan',
    component: () => import('@/layouts/KohorLayout.vue'),
    children: [
      {
        path: '',
        name: 'dashboard-rotator',
        component: () => import('@/pages/dashboard_frontend/DashboardRotator.vue'),
      },
    ],
  },
  {
    path: '/antrian',
    // name: 'landing-kesiswaan',
    component: () => import('@/layouts/KohorLayout.vue'),
    children: [
      {
        path: '',
        name: 'landing-antrian',
        component: () => import('@/pages/QueueDashboard.vue'),
      },
    ],
  },
  // OPERATOR SEKOLAH
  // {
  //   path: '/ops',
  //   component: () => import('@/layouts/OperatorSekolahLayout.vue'),
  //   children: [
  //     {
  //       path: 'dashboard',
  //       name: 'ops-dashboard',
  //       component: () => import('@/pages/operator-sekolah/OpsDashboardPage.vue'),
  //     },
  //     {
  //       path: 'antrean-tugas',
  //       name: 'ops-antrean-tugas',
  //       component: () => import('@/pages/operator-sekolah/OpsAntreanTugasPage.vue'),
  //     },
  //     {
  //       path: 'sinkronisasi',
  //       name: 'ops-sinkronisasi',
  //       component: () => import('@/pages/operator-sekolah/OpsSinkronisasiPage.vue'),
  //     },
  //     {
  //       path: 'riwayat',
  //       name: 'ops-riwayat',
  //       component: () => import('@/pages/operator-sekolah/OpsRiwayatPage.vue'),
  //     },
  //   ],
  // },
  // KEUANGAN
  // {
  //   path: '/keuangan',
  //   component: () => import('@/layouts/KeuanganLayout.vue'),
  //   children: [
  //     {
  //       path: '',
  //       children: [{ path: '', redirect: '/keuangan/dashboard' }],
  //     },
  //     {
  //       path: 'dashboard',
  //       component: () => import('@/pages/keuangan/DashboardKeuanganPage.vue'),
  //     },
  //     {
  //       path: 'pembayaran-lain',
  //       component: () => import('@/pages/keuangan/PembayaranLainPage.vue'),
  //     },
  //     { path: 'tagihan', component: () => import('@/pages/keuangan/TagihanPage.vue') },
  //     { path: 'pembayaran', component: () => import('@/pages/keuangan/PembayaranPage.vue') },
  //     {
  //       path: 'pengeluaran',
  //       component: () => import('@/pages/keuangan/PengeluaranPage.vue'),
  //     },
  //   ],
  // },
  // {
  //   path: '/auth/login',
  //   name: 'auth-login',
  //   component: () => import('@/layouts/LoginLayout.vue'),
  // },
  // {
  //   path: '/admin1',
  //   component: () => import('@/layouts/AdminLayout.vue'),
  // },
  // Penerimaan Siswa Baru

  // ADMIN
  // {
  //   path: '/admin',
  //   component: () => import('@/layouts/AdminLayout.vue'),
  //   children: [
  //     {
  //       path: 'kepegawaian/data-guru',
  //       name: 'admin-data-guru',
  //       component: () => import('@/pages/kepegawaian/DataGuru.vue'),
  //     },
  //     {
  //       path: 'kepegawaian/beban-mengajar',
  //       name: 'admin-beban-mengajar',
  //       component: () => import('@/pages/kepegawaian/BebanMengajar.vue'),
  //     },
  //     {
  //       path: 'kepegawaian/cuti',
  //       name: 'admin-cuti-mengajar',
  //       component: () => import('@/pages/kepegawaian/CutiMengajar.vue'),
  //     },
  //     {
  //       path: 'kepegawaian/arsip',
  //       name: 'admin-arsip-dokumen',
  //       component: () => import('@/pages/kepegawaian/ArsipDokumen.vue'),
  //     },
  //     {
  //       path: 'pengaturan',
  //       name: 'admin-pengaturan',
  //       component: () => import('@/pages/pengaturan/UserManagement.vue'),
  //     },
  //   ],
  // },
  // KURIKULUM
  // Bertanggung jawab atas manajemen akademik utama, kalender akademik, penjadwalan, dan penilaian (Rapor).
  // {
  //   path: '/kurikulum',
  //   children: [
  //     {
  //       path: 'dashboard',
  //     },
  //     {
  //       path: 'tahun-akademik',
  //       // Pengaturan tahun ajaran aktif dan semester.
  //     },
  //     {
  //       path: 'mata-pelajaran',
  //       // Kelola daftar mapel (Nasional, Kewilayahan, Kejuruan)
  //     },
  //     {
  //       path: 'jadwal',
  //       // Pembuatan jadwal pelajaran per kelas dan guru.
  //     },
  //     {
  //       path: 'kelas',
  //       // Pembentukan kelas baru (Struktur rombel)
  //     },
  //     {
  //       path: 'nilai',
  //       // Pembentukan kelas baru (Struktur rombel)
  //       children: [
  //         {
  //           path: 'bobot',
  //           // Pengaturan bobot penilaian (Teori vs Praktik)
  //         },
  //       ],
  //     },
  //   ],
  // },

  // Always leave this as last one,
  // but you can also remove it
  {
    path: '/:catchAll(.*)*',
    component: () => import('@/pages/ErrorNotFound.vue'),
  },
]

export default routes

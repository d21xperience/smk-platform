import { siakadGuard } from '../siakad-guards'

const siakadAdmin = {
  path: 'admin',
  component: () => import('@/layouts/SiakadLayout.vue'),
  beforeEnter: siakadGuard,
  children: [
    // {
    //   path: '/admin',
    //   component: () => import('@/layouts/AdminLayout.vue'), // Kerangka khusus admin panel
    //   children: [
    //     { path: 'dashboard', component: () => import('@/pages/admin/DashboardPage.vue') },
    //     { path: 'articles', component: () => import('@/pages/admin/ManageArticlesPage.vue') },
    //   ],
    // },

    {
      path: '',
      redirect: '/admin/dashboard',
    },
    {
      path: 'dashboard',
      name: 'siakad-dashboard',
      component: () => import('@/pages/siakad/DashboardPage.vue'),
    },
    {
      path: 'contex',
      name: 'contex',
      component: () => import('@/pages/siakad/SelectContextPage.vue'),
    },
    {
      path: 'jadwal-pelajaran',
      name: 'jadwal-pelajaran',
      component: () => import('@/pages/siakad/TeacherSchedulePage.vue'),
    },
    {
      path: 'input-absensi',
      name: 'input-absensi',
      component: () => import('@/pages/siakad/AttendancePage.vue'),
    },
    {
      path: 'profil',
      name: 'profil',
      component: () => import('@/pages/siakad/ProfilPage.vue'),
    },
    {
      path: 'kelender-akademik',
      name: 'kelender-akademik',
      component: () => import('@/pages/siakad/KalenderAkademik.vue'),
    },
    {
      path: 'catatan-siswa',
      name: 'catatan-siswa',
      component: () => import('@/pages/siakad/RiwayatCatatanSiswa.vue'),
    },
    // {
    //   path: 'cek-page',
    //   name: 'cek-page',
    //   component: () => import('@/pages/siakad/CekPage.vue'),
    // },
  ],
}

export default siakadAdmin

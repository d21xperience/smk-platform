const siakadSiswa = {
  path: '/siswa-portal',
  component: () => import('@/layouts/StudentLayout.vue'),
  children: [
    {
      path: 'dashboard-ppdb',
      component: () => import('@/pages/psb/DashboardPpdbSiswa.vue'),
    },
    {
      path: 'biodata',
      component: () => import('@/pages/psb/BiodataMandiriSiswa.vue'),
    },
    {
      path: 'dokumen',
      component: () => import('@/pages/psb/UploadBerkasSiswa.vue'),
    },
    {
      path: 'pembayaran',
      component: () => import('@/pages/psb/BiayaPendaftaranSiswa.vue'),
    },
  ],
}
export default siakadSiswa

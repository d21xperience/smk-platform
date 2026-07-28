const siakadSiswa = {
  path: 'siswa',
  component: () => import('@/layouts/StudentLayout.vue'),
  children: [
    {
      path: 'dashboard-ppdb',
      name: 'siswa-dash-ppdb',
      component: () => import('@/pages/psb/DashboardPpdbSiswa.vue'),
    },
    {
      path: 'biodata',
      name: 'siswa-biodata',
      component: () => import('@/pages/psb/BiodataMandiriSiswa.vue'),
    },
    {
      path: 'dokumen',
      name: 'siswa-dokumen',
      component: () => import('@/pages/psb/UploadBerkasSiswa.vue'),
    },
    {
      path: 'pembayaran',
      name: 'siswa-pembayaran',
      component: () => import('@/pages/psb/BiayaPendaftaranSiswa.vue'),
    },
  ],
}
export default siakadSiswa

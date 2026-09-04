const psb = {
  path: '/psb',
  children: [
    {
      path: 'siswa/baru',
      name: 'tambah-siswa',
      component: () => import('@/pages/psb/SiswaBaruPage.vue'),
    },
    // {
    //   path: 'siswa/ppdb/edit/:id',
    //   name: 'edit-siswa',
    //   component: () => import('pages/psb/FormEditPSB.vue'),
    // },
    // {
    //   path: 'analisis-asal-sekolah',
    //   name: 'analisis-asal-sekolah',
    //   component: () => import('pages/psb/AnalisisAsalSekolah.vue'),
    // },
    // {
    //   path: 'generate-kelas',
    //   name: 'generate-kelas',
    //   component: () => import('pages/psb/PemetaanKelas.vue'),
    // },
  ],
}
export default psb

const hubin = {
  path: '/hubin',
  children: [
    {
      path: 'mitra',
      component: () => import('@/pages/pkl/PklMitra.vue'),
      name: 'pkl-mitra',
    },
    {
      path: 'penempatan',
      component: () => import('@/pages/pkl/PklPenempatan.vue'),
      name: 'pkl-penempatan',
    },
    {
      path: 'logbook',
      component: () => import('@/pages/pkl/PklLogbook.vue'),
      name: 'pkl-logbook',
    },
    {
      path: 'penilaian',
      component: () => import('@/pages/pkl/PklPenilaian.vue'),
      name: 'pkl-penilaian',
    },
    {
      path: 'sertifikat',
      component: () => import('@/pages/pkl/PklSertifikat.vue'),
      name: 'pkl-sertifikat',
    },
    {
      path: 'ploting-berkas-siswa',
      name: 'ploting-berkas-siswa',
      component: () => import('@/pages/kesiswaan/PlottingBerkas.vue'),
    },
    {
      path: 'sertifikat-ukk-siswa',
      name: 'sertifikat-ukk-siswa',
      component: () => import('@/pages/kesiswaan/SertifikasiUkk.vue'),
    },
    {
      path: 'asuransi-kerja-siswa',
      name: 'asuransi-kerja-siswa',
      component: () => import('@/pages/kesiswaan/AsuransiKerja.vue'),
    },
    {
      path: 'tracer-study-alumni',
      name: 'tracer-study-alumni',
      component: () => import('@/pages/kesiswaan/TracerStudy.vue'),
    },
    {
      path: 'loker',
      name: 'loker',
      component: () => import('@/pages/kesiswaan/LowonganKerja.vue'),
    },
  ],
}
export default hubin

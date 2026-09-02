const hubin = {
  path: '/hubin',
  children: [
    {
      path: 'mitra',
      component: () => import('pages/pkl/PklMitra.vue'),
      name: 'pkl-mitra',
    },
    {
      path: 'penempatan',
      component: () => import('pages/pkl/PklPenempatan.vue'),
      name: 'pkl-penempatan',
    },
    {
      path: 'logbook',
      component: () => import('pages/pkl/PklLogbook.vue'),
      name: 'pkl-logbook',
    },
    {
      path: 'penilaian',
      component: () => import('pages/pkl/PklPenilaian.vue'),
      name: 'pkl-penilaian',
    },
    {
      path: 'sertifikat',
      component: () => import('pages/pkl/PklSertifikat.vue'),
      name: 'pkl-sertifikat',
    },
  ],
}
export default hubin

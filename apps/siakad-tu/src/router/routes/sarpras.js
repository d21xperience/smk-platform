const sarpras = {
  path: '/sarpras',
  children: [
    {
      path: '',
      redirect: '/sarpras/dashboard',
    },
    {
      path: 'dashboard',
      component: () => import('pages/sarpras/SarPrasDashboard.vue'),
    },
    {
      path: 'inventaris/tambah',
      component: () => import('pages/sarpras/SarPrasAssetForm.vue'),
    },
    {
      path: 'inventaris/logistik',
      component: () => import('pages/sarpras/SarPrasKatalogLogistik.vue'),
    },
    {
      path: 'inventaris/ruangan',
      component: () => import('pages/sarpras/SarPrasKartuInventarisRuang.vue'),
    },
    {
      path: 'inventaris/cetak-qr',
      component: () => import('pages/sarpras/SarPrasCetakQrAset.vue'),
    },
    {
      path: 'servis/tiket',
      component: () => import('pages/sarpras/SarPrasTiketServis.vue'),
    },
    {
      path: 'servis/jadwal',
      component: () => import('pages/sarpras/SarPrasPerawatanBerkala.vue'),
    },
    {
      path: 'servis/vendor',
      component: () => import('pages/sarpras/SarPrasDataVendor.vue'),
    },
    {
      path: 'sirkulasi/pinjam',
      component: () => import('pages/sarpras/SarPrasPeminjamanAset.vue'),
    },
    {
      path: 'sirkulasi/minta-atk',
      component: () => import('pages/sarpras/SarPrasPermintaanAtk.vue'),
    },
    {
      path: 'mutasi/log',
      component: () => import('pages/sarpras/SarPrasLogMutasi.vue'),
    },
    {
      path: 'mutasi/hapus',
      component: () => import('pages/sarpras/SarPrasPenghapusanAset.vue'),
    },
    {
      path: 'laporan',
      component: () => import('pages/sarpras/SarPrasLaporanAkuntansi.vue'),
    },
  ],
}
export default sarpras

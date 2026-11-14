const financeRoutes = {
  path: '/keuangan',
  component: () => import('@/layouts/finance/FinanceLayout.vue'),
  meta: {
    requiresAuth: true,
    allowedRoles: ['bendahara', 'kepala_sekolah', 'admin'],
  },
  children: [
    { path: '', redirect: '/keuangan/dashboard' },

    {
      path: 'dashboard',
      name: 'finance-dashboard',
      component: () => import('@/pages/finance/DashboardFinancePage.vue'),
    },

    // Penerimaan
    {
      path: 'penerimaan/input',
      name: 'finance-penerimaan-input',
      component: () => import('@/pages/finance/penerimaan/InputPenerimaanPage.vue'),
      meta: { title: 'Input Penerimaan Iuran' },
    },
    {
      path: 'penerimaan/rekap-harian',
      name: 'finance-penerimaan-rekap-harian',
      component: () => import('@/pages/finance/penerimaan/RekapHarianPage.vue'),
      meta: { title: 'Rekap Harian' },
    },
    {
      path: 'penerimaan/rekap-kelas',
      name: 'finance-penerimaan-rekap-kelas',
      component: () => import('@/pages/finance/penerimaan/RekapKelasPage.vue'),
      meta: { title: 'Rekap Per Kelas' },
    },

    // Tagihan
    {
      path: 'tagihan',
      name: 'finance-tagihan',
      component: () => import('@/pages/finance/TagihanPage.vue'),
      meta: { title: 'Daftar Tagihan' },
    },

    // BKU
    {
      path: 'bku',
      name: 'finance-bku',
      component: () => import('@/pages/finance/BkuPage.vue'),
      meta: { title: 'Buku Kas Umum' },
    },

    // Pengeluaran
    {
      path: 'pengeluaran/ajukan',
      name: 'finance-pengeluaran-ajukan',
      component: () => import('@/pages/finance/pengeluaran/AjukanPengeluaranPage.vue'),
      meta: { title: 'Ajukan Pengeluaran' },
    },
    {
      path: 'pengeluaran/verifikasi',
      name: 'finance-pengeluaran-verifikasi',
      component: () => import('@/pages/finance/pengeluaran/VerifikasiBendaharaPage.vue'),
      meta: { title: 'Verifikasi Bendahara' },
    },
    {
      path: 'pengeluaran/persetujuan',
      name: 'finance-pengeluaran-persetujuan',
      component: () => import('@/pages/finance/pengeluaran/PersetujuanKepsekPage.vue'),
      meta: { title: 'Persetujuan Kepala Sekolah' },
    },
    {
      path: 'pengeluaran',
      name: 'finance-pengeluaran',
      component: () => import('@/pages/finance/pengeluaran/PengeluaranListPage.vue'),
      meta: { title: 'Daftar Pengeluaran' },
    },

    // RAB
    {
      path: 'rab',
      name: 'finance-rab',
      component: () => import('@/pages/finance/RabFormPage.vue'),
      meta: { title: 'RAB Tahunan' },
    },
    {
      path: 'rab/realisasi',
      name: 'finance-rab-realisasi',
      component: () => import('@/pages/finance/rab/RealisasiPage.vue'),
      meta: { title: 'Realisasi Anggaran' },
    },

    // Rekonsiliasi
    {
      path: 'rekonsiliasi',
      name: 'finance-rekonsiliasi',
      component: () => import('@/pages/finance/RekonsiliasiPage.vue'),
      meta: { title: 'Rekonsiliasi Bank' },
    },

    // Piutang
    {
      path: 'piutang',
      name: 'finance-piutang',
      component: () => import('@/pages/finance/PiutangPage.vue'),
      meta: { title: 'Buku Piutang' },
    },

    // Penutupan
    {
      path: 'penutupan',
      name: 'finance-penutupan',
      component: () => import('@/pages/finance/PenutupanPage.vue'),
      meta: { title: 'Penutupan Buku Kas' },
    },

    // Laporan
    {
      path: 'laporan/neraca',
      name: 'finance-laporan-neraca',
      component: () => import('@/pages/finance/laporan/NeracaPage.vue'),
      meta: { title: 'Neraca' },
    },
    {
      path: 'laporan/lra',
      name: 'finance-laporan-lra',
      component: () => import('@/pages/finance/laporan/LraPage.vue'),
      meta: { title: 'LRA' },
    },
    {
      path: 'laporan/ringkasan-bulan',
      name: 'finance-laporan-ringkasan-bulan',
      component: () => import('@/pages/finance/laporan/RingkasanBulanPage.vue'),
      meta: { title: 'Ringkasan Bulan' },
    },

    // Master
    {
      path: 'master/jenis-iuran',
      name: 'finance-master-jenis-iuran',
      component: () => import('@/pages/finance/master/JenisIuranPage.vue'),
      meta: { title: 'Jenis Iuran' },
    },
    {
      path: 'master/kode-rekening',
      name: 'finance-master-kode-rekening',
      component: () => import('@/pages/finance/master/KodeRekeningPage.vue'),
      meta: { title: 'Kode Rekening (COA)' },
    },
    {
      path: 'master/rekening-bank',
      name: 'finance-master-rekening-bank',
      component: () => import('@/pages/finance/master/RekeningBankPage.vue'),
      meta: { title: 'Rekening Bank' },
    },
    {
      path: 'master/tarif-iuran',
      name: 'finance-master-tarif-iuran',
      component: () => import('@/pages/finance/master/TarifIuranPage.vue'),
      meta: { title: 'Tarif Iuran' },
    },

    // Pengaturan
    {
      path: 'pengaturan/pin-ttd',
      name: 'finance-pengaturan-pin-ttd',
      component: () => import('@/pages/finance/PinTTDPage.vue'),
      meta: { title: 'PIN TTD Digital' },
    },
  ],
}

export default financeRoutes

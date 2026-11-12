export default {
  path: 'kesiswaan',
  children: [
    {
      path: '',
      redirect: '/kesiswaan/dashboard',
    },
    {
      path: 'dashboard',
      name: 'kesiswaan-dashboard',
      component: () => import('@/pages/kesiswaan/DashboardPage.vue'),
    },
    {
      path: 'surat-aktif-siswa',
      name: 'surat-aktif-siswa',
      component: () => import('@/pages/kesiswaan/pelaporan/SuratAktif.vue'),
    },
    // MUTASI
    {
      path: 'mutasi-siswa',
      name: 'kesiswaan-mutasi',
      component: () => import('@/pages/kesiswaan/pelaporan/ManajemenMutasi.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: 'pengajuan-mutasi',
      name: 'kesiswaan-pengajuan-mutasi',
      component: () => import('@/pages/kesiswaan/pelaporan/PengajuanMutasi.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: 'siswa',
      name: 'kesiswaan-siswa',
      component: () => import('@/pages/kesiswaan/siswa/StudentListPage.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: 'siswa/:id',
      name: 'kesiswaan-siswa-detail',
      component: () => import('@/pages/kesiswaan/siswa/StudentDetailPage.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: 'siswa/edit/:id',
      name: 'kesiswaan-siswa-edit',
      component: () => import('@/pages/kesiswaan/siswa/StudentEditPage.vue'),
      meta: { requiresAuth: true },
    },
    // ABSENSI
    {
      path: 'absensi-harian',
      name: 'kesiswaan-absensi-harian',
      component: () => import('@/pages/kesiswaan/absensi/AttendanceSessionPage.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: 'rekap-kehadiran',
      name: 'kesiswaan-rekap-kehadiran',
      component: () => import('@/pages/kesiswaan/absensi/AttendanceRecapPage.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: 'beasiswa-pip-siswa',
      name: 'kesiswaan-beasiswa-pip',
      component: () => import('@/pages/kesiswaan/kesejahteraan/BeasiswaPip.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: 'rekap-pelaporan',
      name: 'rekap-pelaporan',
      component: () => import('@/pages/kesiswaan/pelaporan/RekapPelaporan.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: 'legalisir-alumni',
      name: 'legalisir-alumni',
      component: () => import('@/pages/kesiswaan/alumni/LegalisirAlumni.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: 'tracer-study',
      name: 'tracer-study',
      component: () => import('@/pages/kesiswaan/alumni/TracerStudy.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: 'ploting-berkas-siswa',
      name: 'ploting-berkas-siswa',
      component: () => import('@/pages/kesiswaan/pelaporan/PlottingBerkas.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: 'buku-induk-siswa',
      name: 'buku-induk-siswa',
      component: () => import('@/pages/kesiswaan/siswa/BukuIndukSiswa.vue'),
      meta: { requiresAuth: true },
    },
  ],
}

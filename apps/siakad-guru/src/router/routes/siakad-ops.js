// import { siakadGuard } from '../siakad-guards'

const siakadOps = {
  path: 'ops',
  component: () => import('@/layouts/OperatorSekolahLayout.vue'),
  // beforeEnter: siakadGuard,
  children: [
    {
      path: '',
      redirect: '/siakad/ops/dashboard',
    },
    {
      path: 'dashboard',
      name: 'ops-dashboard',
      component: () => import('@/pages/operator-sekolah/OpsDashboardPage.vue'),
    },
    {
      path: 'antrean-tugas',
      name: 'ops-antrean-tugas',
      component: () => import('@/pages/operator-sekolah/OpsAntreanTugasPage.vue'),
    },
    {
      path: 'sinkronisasi',
      name: 'ops-sinkronisasi',
      component: () => import('@/pages/operator-sekolah/OpsSinkronisasiPage.vue'),
    },
    {
      path: 'riwayat',
      name: 'ops-riwayat',
      component: () => import('@/pages/operator-sekolah/OpsRiwayatPage.vue'),
    },
  ],
}

export default siakadOps

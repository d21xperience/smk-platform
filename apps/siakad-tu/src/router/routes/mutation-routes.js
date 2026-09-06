// apps/siakad-tu/src/router/routes/mutation.js

export default {
  path: '/mutation',
  component: () => import('layouts/MainLayout.vue'),
  meta: {
    requiresAuth: true,
    breadcrumb: 'Mutasi Siswa',
  },
  children: [
    {
      path: '',
      name: 'mutation-list',
      component: () => import('pages/mutation/MutationListPage.vue'),
      meta: {
        title: 'Mutasi Siswa',
        permission: 'mutation.view',
        breadcrumb: 'Daftar Mutasi',
      },
    },
    {
      path: 'request-in',
      name: 'mutation-request-in',
      component: () => import('pages/mutation/MutationRequestPage.vue'),
      meta: {
        title: 'Permohonan Mutasi Masuk',
        permission: 'mutation.create',
        breadcrumb: 'Mutasi Masuk',
      },
      props: { type: 'IN' },
    },
    {
      path: 'request-out',
      name: 'mutation-request-out',
      component: () => import('pages/mutation/MutationRequestPage.vue'),
      meta: {
        title: 'Permohonan Mutasi Keluar',
        permission: 'mutation.create',
        breadcrumb: 'Mutasi Keluar',
      },
      props: { type: 'OUT' },
    },
    {
      path: ':id',
      name: 'mutation-detail',
      component: () => import('pages/mutation/MutationDetailPage.vue'),
      meta: {
        title: 'Detail Mutasi',
        permission: 'mutation.view',
        breadcrumb: 'Detail',
      },
      props: true,
    },
  ],
}

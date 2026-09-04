// apps/siakad-tu/src/router/routes/registration.js

export default {
  path: '/registration',
  component: () => import('layouts/MainLayout.vue'),
  meta: {
    requiresAuth: true,
    breadcrumb: 'PPDB',
  },
  children: [
    {
      path: '',
      name: 'registration-list',
      component: () => import('pages/registration/RegistrationListPage.vue'),
      meta: {
        title: 'Penerimaan Peserta Didik Baru',
        permission: 'registration.view',
        breadcrumb: 'Daftar Pendaftaran',
      },
    },
    {
      path: 'create',
      name: 'registration-create',
      component: () => import('pages/registration/RegistrationFormPage.vue'),
      meta: {
        title: 'Pendaftaran Baru',
        permission: 'registration.create',
        breadcrumb: 'Pendaftaran Baru',
      },
    },
    {
      path: ':id',
      name: 'registration-detail',
      component: () => import('pages/registration/RegistrationDetailPage.vue'),
      meta: {
        title: 'Detail Pendaftaran',
        permission: 'registration.view',
        breadcrumb: 'Detail',
      },
      props: true,
    },
    {
      path: ':id/edit',
      name: 'registration-edit',
      component: () => import('pages/registration/RegistrationFormPage.vue'),
      meta: {
        title: 'Edit Pendaftaran',
        permission: 'registration.edit',
        breadcrumb: 'Edit',
      },
      props: true,
    },
  ],
}

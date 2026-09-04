const keuanganRoutes = {
  path: '/keuangan',
  meta: {
    requiresAuth: true,
    allowedRoles: ['admin', 'tu', 'keuangan'],
  },
  children: [
    {
      path: '',
      redirect: '/keuangan/dashboard',
    },
    {
      path: 'dashboard',
      name: 'keuangan-dashboard',
      component: () => import('@/pages/keuangan/DashboardKeuanganPage.vue'),
    },
    {
      path: 'tagihan',
      name: 'keuangan-tagihan',
      component: () => import('@/pages/keuangan/InvoiceListPage.vue'),
    },
    {
      path: 'tagihan/create',
      name: 'keuangan-tagihan-create',
      component: () => import('@/pages/keuangan/InvoiceFormPage.vue'),
    },
    {
      path: 'tagihan/:id',
      name: 'keuangan-tagihan-detail',
      component: () => import('@/pages/keuangan/InvoiceDetailPage.vue'),
      props: true,
    },
    {
      path: 'tagihan/:id/payment',
      name: 'keuangan-pembayaran',
      component: () => import('@/pages/keuangan/PaymentFormPage.vue'),
      props: true,
    },
    {
      path: 'tunggakan',
      name: 'keuangan-tunggakan',
      component: () => import('@/pages/keuangan/OutstandingListPage.vue'),
    },
  ],
}

export default keuanganRoutes

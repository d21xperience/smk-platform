const routes = [
  {
    path: '/',
    component: () => import('@/layouts/MainLayout.vue'),
    meta: { public: true },
    children: [
      { path: '', component: () => import('@/pages/landing/IndexPage.vue') },
      { path: 'mading', component: () => import('@/pages/landing/MadingPage.vue') },
      { path: 'faq', component: () => import('@/pages/landing/PaqPage.vue') },
      { path: 'jurusan/:id', component: () => import('@/pages/landing/JurusanDetailPage.vue') },
      { path: 'berita', component: () => import('@/pages/landing/ArticlesPage.vue') },
      { path: 'simulasi', component: () => import('@/pages/landing/CostSimulationPage.vue') },
      { path: 'fasilitas', component: () => import('@/pages/landing/FacilitiesPage.vue') },
      { path: 'agenda', component: () => import('@/pages/landing/AgendaPage.vue') },
      { path: 'psb', component: () => import('@/pages/landing/PSBPage.vue') },
    ],
  },
  {
    path: '/auth',
    component: () => import('@/layouts/LoginLayout.vue'),
    meta: { public: true },
    children: [
      { path: '', redirect: '/auth/login' },
      { path: 'login', component: () => import('@/pages/LoginPage.vue'), name: 'webLogin' },
    ],
  },
  {
    path: '/admin',
    component: () => import('@/layouts/AdminLayout.vue'), // Kerangka khusus admin panel
    children: [
      { path: 'dashboard', component: () => import('@/pages/admin/DashboardPage.vue') },
      { path: 'articles', component: () => import('@/pages/admin/ManageArticlesPage.vue') },
    ],
  },

  {
    path: '/unauthorized',
    component: () => import('@/pages/ErrorUnauthorized.vue'), // 🚀 Hubungkan ke file komponen baru
    meta: { title: 'Akses Ditolak' },
  },

  // Always leave this as last one,
  // but you can also remove it
  {
    path: '/:catchAll(.*)*',
    component: () => import('@/pages/ErrorNotFound.vue'),
  },
]

export default routes

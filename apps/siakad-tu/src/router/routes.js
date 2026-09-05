// apps/siakad-tu/src/router/routes.js

import kesiswaanRoutes from './routes/kesiswaan-routes.js'
import persuratan from './routes/persuratan.js'
import previewDev from './routes/preview-routes.js'
import psb from './routes/psb.js'

// Import routes baru
// import keuanganRoutes from './routes/keuangan-routes.js'
// import absensiRoutes from './routes/absensi-routes.js'
// import penilaianRoutes from './routes/penilaian-routes.js'

const routes = [
  {
    path: '/',
    component: () => import('@/layouts/LandingLayout.vue'), // <-- Gunakan layout yang sudah ada
    children: [
      previewDev,
      // Existing routes
      persuratan,
      kesiswaanRoutes,
      psb,

      // New routes (Frontend Integration)
      // keuanganRoutes,
      // absensiRoutes,
      // penilaianRoutes,
    ],
  },
  // === AUTH ROUTES (Tanpa Layout Landing) ===
  {
    path: '/auth',
    children: [
      {
        path: 'login',
        name: 'auth-login',
        component: () => import('@/pages/auth/LoginPage.vue'),
        meta: {
          guestOnly: true, // Hanya untuk user yang belum login
        },
      },
    ],
  },
  {
    path: '/kohor',
    component: () => import('@/layouts/KohorLayout.vue'),
    children: [
      {
        path: '',
        name: 'dashboard-rotator',
        component: () => import('@/pages/dashboard_frontend/DashboardRotator.vue'),
      },
    ],
  },
  {
    path: '/:catchAll(.*)*',
    component: () => import('@/pages/ErrorNotFound.vue'),
  },
]

export default routes

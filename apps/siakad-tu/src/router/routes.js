// import hubin from './routes/hubin-routes.js'
import kepegawaian from './routes/kepegawaian-routes.js'
import kesiswaanRoutes from './routes/kesiswaan-routes.js'
import kurikulumRoutes from './routes/kurikulum-routes.js'
import persuratan from './routes/persuratan.js'
import previewDev from './routes/preview-routes.js'
import psb from './routes/psb-routes.js'
import sarpras from './routes/sarpras.js'

const routes = [
  {
    path: '/',
    component: () => import('@/layouts/LandingLayout.vue'),
    children: [
      previewDev,
      kesiswaanRoutes,
      kurikulumRoutes,
      persuratan,
      psb,
      // hubin,
      kepegawaian,
      sarpras,
      {
        path: 'pengaturan/tahun-ajaran',
        name: 'pengaturan-tahun-ajaran',
        component: () => import('@/pages/pengaturan/AcademicPeriodPage.vue'),
        meta: { requiresAuth: true },
      },
      {
        path: 'pengaturan/pengguna',
        name: 'pengaturan-pengguna',
        component: () => import('@/pages/pengaturan/UserManagementPage.vue'),
        meta: { requiresAuth: true },
      },
      {
        path: 'pengaturan/sekolah',
        name: 'pengaturan-sekolah',
        component: () => import('@/pages/pengaturan/SchoolProfilePage.vue'),
        meta: { requiresAuth: true },
      },
    ],
  },
  {
    path: '/auth',
    children: [
      {
        path: 'login',
        name: 'auth-login',
        component: () => import('@/pages/auth/LoginPage.vue'),
        meta: { guestOnly: true },
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

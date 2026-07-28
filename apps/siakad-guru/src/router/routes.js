import siakadGuru from './routes/siakad-guru'
import siakadOps from './routes/siakad-ops'

const routes = [
  {
    path: '/',
    meta: { public: true },
    children: [{ path: '', redirect: '/auth' }],
  },
  {
    path: '/auth',
    component: () => import('@/layouts/LoginLayout.vue'),
    meta: { public: true },
    children: [
      {
        path: '',
        component: () => import('@/pages/SiakadLoginPage.vue'),
        name: 'siakadLogin',
      },
    ],
  },
  //
  {
    path: '/siakad',
    children: [siakadGuru, siakadOps],
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

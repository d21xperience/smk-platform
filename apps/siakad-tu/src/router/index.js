import { route } from 'quasar/wrappers'
import {
  createRouter,
  createMemoryHistory,
  createWebHistory,
  createWebHashHistory,
} from 'vue-router'
import routes from './routes.js'
import { guestGuard } from './guards.js'
import { useAuthStore } from '@/stores/authStore.js'

export default route(function (/* { store, ssrContext } */) {
  // MENGGUNAKAN FORMAT VITE (import.meta.env) YANG AMAN
  const createHistory = import.meta.env.SSR
    ? createMemoryHistory
    : import.meta.env.VUE_ROUTER_MODE === 'history'
      ? createWebHistory
      : createWebHashHistory

  const Router = createRouter({
    scrollBehavior: () => ({ left: 0, top: 0 }),
    routes,
    // MENGGUNAKAN FORMAT VITE (import.meta.env) YANG AMAN
    history: createHistory(import.meta.env.VUE_ROUTER_BASE),
  })

  Router.beforeEach(async (to) => {
    const authStore = useAuthStore()
    if (to.meta.public) {
      const guestRedirect = guestGuard(to) // ðŸš€ Jalankan guest guard
      if (guestRedirect) return guestRedirect
      return true
    }
    // Cek apakah rute ini butuh login
    if (to.meta.requiresAuth && !authStore.isAuthenticated) {
      // 1. Batalkan navigasi ke halaman tujuan
      // 2. Munculkan pop-up login
      // 3. Simpan URL tujuan agar setelah login bisa diarahkan ke sana
      authStore.triggerLogin(to.fullPath)
      return false
    }
  })

  return Router
})

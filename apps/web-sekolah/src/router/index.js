import { route } from 'quasar/wrappers'
import {
  createRouter,
  createMemoryHistory,
  createWebHistory,
  createWebHashHistory,
} from 'vue-router'
import routes from './routes'
import { guestGuard } from './guards'

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
    if (to.meta.public) {
      const guestRedirect = guestGuard(to) // 🚀 Jalankan guest guard
      if (guestRedirect) return guestRedirect
      return true
    }
  })

  return Router
})

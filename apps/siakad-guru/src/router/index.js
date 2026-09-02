import { route } from 'quasar/wrappers'
import {
  createRouter,
  createMemoryHistory,
  createWebHistory,
  createWebHashHistory,
} from 'vue-router'
import routes from './routes'

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

  Router.onError((error) => {
    console.error('❌Terjadi error saat navigasi rute:', error)
    // console.log(Router.currentRoute.value)
    // if (error) {
    //   Router.push({
    //     name: 'ErrorServer',
    //     query: { message: error.message || 'Terjadi kesalahan sistem' },
    //   })
    // }
    // Contoh penanganan: Jika error karena file chunk gagal dimuat (misal setelah deploy baru)
    // if (error.message.includes('Failed to fetch dynamically imported module')) {
    //   window.location.reload()
    // }
  })
  return Router
})

import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/authStore'
// import { useQuasar } from 'quasar'

export function useProtectedNavigation() {
  const router = useRouter()
  const authStore = useAuthStore()
  // const $q = useQuasar()

  /**
   * Menavigasi ke rute dengan otomatis mengecek hak akses dari konfigurasi Router.
   * @param {string|object} targetRoute - Tujuan rute (path string atau object route)
   * @param {string} [customMessage] - Pesan opsional untuk dialog login
   */
  const navigateTo = (targetRoute, customMessage = null) => {
    // 1. "Intip" konfigurasi rute tujuan
    const resolvedRoute = router.resolve(targetRoute)
    const meta = resolvedRoute.meta || {}

    const requiresAuth = meta.requiresAuth
    const allowedRoles = meta.allowedRoles // Baca dari guard.js / routes.js

    const isAuthenticated = authStore.isAuthenticated
    const userRole = authStore.user?.role

    // 2. Logika Pengecekan
    // Skenario A: Belum login, tapi rute butuh auth atau punya batasan role
    if (!isAuthenticated && (requiresAuth || allowedRoles)) {
      const msg =
        customMessage || 'Halaman ini memerlukan autentikasi. Silakan login untuk melanjutkan.'
      authStore.triggerLogin(targetRoute, msg)
      return
    }

    // Skenario B: Sudah login, tapi role tidak diizinkan
    if (allowedRoles && !allowedRoles.includes(userRole)) {
      // UX yang lebih baik: Beri tahu user bahwa mereka sudah login, tapi salah role
      const msg =
        customMessage || `Akses ditolak. Halaman ini khusus untuk role: ${allowedRoles.join(', ')}.`

      // Kita tetap pakai triggerLogin agar user bisa switch account,
      // tapi dengan pesan yang sangat spesifik.
      authStore.triggerLogin(targetRoute, msg)

      // Opsional: Bisa juga langsung tolak dengan notifikasi tanpa buka dialog login:
      // $q.notify({ color: 'negative', message: msg, icon: 'block', position: 'top' })
      // return
    }

    // 3. Jika semua aman, lakukan navigasi
    router.push(targetRoute)
  }

  return {
    navigateTo,
  }
}

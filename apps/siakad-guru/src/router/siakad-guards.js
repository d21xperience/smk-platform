// src/router/siakad-guards.js
import { useAuthStore } from '@/stores/auth.store'
import { useContextStore } from '@/stores/context.store'

// Fungsi helper menunggu hingga store siap
function waitForAuthInit() {
  const auth = useAuthStore()
  if (auth.initialized) return Promise.resolve()
  return new Promise((resolve) => {
    const watcher = setInterval(() => {
      if (auth.initialized) {
        clearInterval(watcher)
        resolve()
      }
    }, 50) // cek setiap 50ms
  })
}
export async function authInitGuard() {
  await waitForAuthInit()
  return true
}

export function siakadGuard(to) {
  const auth = useAuthStore()
  if (auth.isAuthenticated) {
    return true // Izinkan akses halaman internal lainnya
  } else {
    // Jalankan pengalihan rute dengan menyertakan halaman asal
    return {
      path: '/auth/siakad',
      query: { redirect: to.fullPath },
    }
  }
}

// Guard untuk halaman transaksi: harus ada operational context
export function operationalContextGuard() {
  const ctx = useContextStore()
  if (!ctx.operational.academicYearId) {
    // Redirect ke halaman pemilihan konteks
    return '/select-context'
  } else {
    return false
  }
}

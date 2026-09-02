// src/router/guards.js
import { useAuthStore } from '@/stores/authStore'

async function waitForAuthInit() {
  const auth = useAuthStore()
  if (auth.initialized) return

  const MAX_WAIT = 5000 // 5 detik
  const start = Date.now()

  // eslint-disable-next-line no-unused-vars
  return new Promise((resolve, reject) => {
    const interval = setInterval(() => {
      if (auth.initialized) {
        clearInterval(interval)
        resolve()
      } else if (Date.now() - start > MAX_WAIT) {
        clearInterval(interval)
        // Force initialized agar tidak menggantung
        auth.initialized = true
        resolve()
      }
    }, 100)
  })
}

// authInitGuard: tunggu inisialisasi, lalu kembalikan true
export async function authInitGuard() {
  await waitForAuthInit()
  return true
}

// authGuard: cek autentikasi, kembalikan path redirect jika belum login
export function authGuard() {
  const auth = useAuthStore()
  if (!auth.isAuthenticated) {
    return '/auth/login'
  }
  return true
}

export function guestGuard(to) {
  const auth = useAuthStore()
  // Jika user sudah login tapi nekat buka halaman tamu (seperti Login)
  if (auth.isAuthenticated) {
    // Jika login pada aplikasi siakad
    if (to.path === '/auth/siakad') return { path: '/siakad/dashboard' } // ðŸš€ Alihkan ke halaman utama internal
    // tambahkan yang lainnya di bawah
  }

  return true // Jika belum login, silakan akses halaman login
}

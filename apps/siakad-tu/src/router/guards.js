import { useAuthStore } from '@/stores/authStore'
import { useContextStore } from '@/stores/contextStore'

export function guestGuard(to, from, next) {
  const authStore = useAuthStore()
  if (authStore.isAuthenticated) {
    return next({ name: 'landing-home' })
  }
  next()
}

export function authGuard(to, from, next) {
  const authStore = useAuthStore()
  if (!authStore.isAuthenticated) {
    // Redirect ke landing, lalu biarkan layout yang memunculkan pop-up login
    // (Karena aplikasi Anda menggunakan pop-up, bukan halaman /login terpisah)
    return next({ name: 'landing-home' })
  }
  next()
}

export function contextGuard(to, from, next) {
  const contextStore = useContextStore()
  // PERBAIKAN: Gunakan .current bukan .currentContext
  if (!contextStore.current) {
    return next({ name: 'landing-home' })
  }
  next()
}

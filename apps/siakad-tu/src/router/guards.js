import { useAuthStore } from '@/stores/authStore'
import { useContextStore } from '@/stores/contextStore'
/**
Guest Guard - Redirect authenticated users away from guest-only pages
(e.g., login, register pages)
*/
export function guestGuard(to, from, next) {
  const authStore = useAuthStore()
  if (authStore.isAuthenticated) {
    // User sudah login, redirect ke dashboard
    return next({ name: 'landing-home' })
  }
  next()
}
/**
Auth Guard - Require authentication for protected routes
*/
export function authGuard(to, from, next) {
  const authStore = useAuthStore()
  if (!authStore.isAuthenticated) {
    // User belum login, redirect ke login page
    return next({ name: 'auth-login' })
  }
  next()
}
/**
Context Guard - Require operational context to be loaded
*/
export function contextGuard(to, from, next) {
  const contextStore = useContextStore()
  if (!contextStore.currentContext) {
    // Context belum dimuat, redirect ke select context page
    return next({ name: 'select-context' })
  }
  next()
}

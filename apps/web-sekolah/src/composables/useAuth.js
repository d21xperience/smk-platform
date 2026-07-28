import { computed } from 'vue'
import { useAuthStore } from '@/stores/auth.store'

export function useAuth() {
  const store = useAuthStore()

  return {
    // State (readonly via computed)
    user: computed(() => store.user),
    isAuthenticated: computed(() => store.isAuthenticated),
    userRole: computed(() => store.userRole),
    loading: computed(() => store.loading),
    initialized: computed(() => store.initialized),
    error: computed(() => store.error),

    // Actions
    login: store.login,
    logout: store.logout,
    checkAuth: store.checkAuth,
  }
}

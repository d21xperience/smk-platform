// FILE: src/router/guards.js
// STATUS: MODIFY
// STATUS IMPLEMENTASI: COMPLETE

import { useAuthStore } from '@/stores/authStore.js'
import { useContextStore } from '@/stores/contextStore.js'

export function createAuthGuard() {
  return async (to) => {
    const authStore = useAuthStore()
    const contextStore = useContextStore()

    if (!authStore.initialized) {
      await authStore.initialize()
    }

    if (authStore.isExpired) {
      if (to.name !== 'login') {
        return { name: 'login', query: { expired: 'true' } }
      }

      return
    }

    const requiresAuth = to.matched.some((record) => record.meta.requiresAuth !== false)
    const publicPages = ['login', 'ErrorServer']
    const isPublicPage = publicPages.includes(to.name)

    if (requiresAuth && !authStore.isAuthenticated) {
      return { name: 'login' }
    }

    if (to.name === 'login' && authStore.isAuthenticated) {
      return { name: 'dashboard' }
    }

    if (to.name === 'context-selection' && authStore.isAuthenticated) {
      return { name: 'dashboard' }
    }

    if (authStore.isAuthenticated && to.meta.requiredPermission) {
      const userPermissions = authStore.currentUser?.permissions || []
      const userRole = authStore.currentUser?.role

      const isAdmin = userRole === 'admin' || userRole === 'superadmin'
      const hasPermission = userPermissions.includes(to.meta.requiredPermission)

      if (!isAdmin && !hasPermission) {
        return { name: 'unauthorized' }
      }
    }

    if (authStore.isAuthenticated && !isPublicPage) {
      if (!contextStore.isContextReady) {
        try {
          await contextStore.ensureCurrentContext(authStore.currentUser)
        } catch (err) {
          return {
            name: 'ErrorServer',
            query: {
              reason: 'context',
              message: err.message,
            },
          }
        }
      }

      if (!contextStore.isContextReady) {
        return {
          name: 'ErrorServer',
          query: {
            reason: 'context-empty',
          },
        }
      }
    }

    return
  }
}

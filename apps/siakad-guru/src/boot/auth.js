import { boot } from 'quasar/wrappers'
import { useAuthStore } from '../stores/authStore.js'
import { createAuthGuard } from '../router/guards.js'

export default boot(async ({ router }) => {
  const authStore = useAuthStore()
  // Initialize authentication (session restoration on app boot)
  await authStore.initialize()

  // Register auth + context guard
  const authGuard = createAuthGuard()
  router.beforeEach(authGuard)
})

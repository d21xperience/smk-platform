import { boot } from 'quasar/wrappers'
import { useContextStore } from '../stores/contextStore.js'
import { useAuthStore } from '../stores/authStore.js'
import { eventDispatcher } from './services.js'
import { EVENT_CONTEXT_SELECTED } from '../events/Type.js'
export default boot(({ app }) => {
  const contextStore = useContextStore()
  const authStore = useAuthStore()
  // Subscribe to context events for coordination
  eventDispatcher.subscribe(EVENT_CONTEXT_SELECTED, () => {
    // Context selected - additional coordination if needed
  })
  // Provide context store globally if needed
  app.config.globalProperties.$contextStore = contextStore
  app.config.globalProperties.$authStore = authStore
})

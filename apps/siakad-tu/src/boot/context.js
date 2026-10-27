import { boot } from 'quasar/wrappers'
import { useContextStore } from '@/stores/contextStore'
import { ContextEngine } from '@/domain/context/engine/ContextEngine'

export default boot(() => {
  const store = useContextStore()
  const engine = new ContextEngine()

  try {
    const saved = localStorage.getItem('operational_context')
    if (saved) {
      const data = JSON.parse(saved)
      const { context } = engine.establishContext(data)
      store.setContext(context)
    }
  } catch (error) {
    console.error('[Boot Context] Gagal restore context:', error)
    localStorage.removeItem('operational_context')
    store.clearContext()
  }
})

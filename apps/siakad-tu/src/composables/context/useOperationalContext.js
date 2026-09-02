// apps/siakad-tu/src/composables/context/useOperationalContext.js

import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useContextStore } from '@/stores/contextStore'
import { ContextEngine } from '@/engine/context/ContextEngine'
import { eventDispatcher } from '@/events/eventDispatcher' // Akan dibuat

/**
 * Composable untuk mengakses dan mengelola Operational Context di Vue components.
 * Ini adalah satu-satunya pintu masuk UI ke context.
 */
export function useOperationalContext() {
  const store = useContextStore()
  const engine = new ContextEngine()
  const router = useRouter()

  /**
   * Reaktif: context saat ini
   */
  const current = computed(() => store.current)
  const isReady = computed(() => store.isContextReady)
  const displayLabel = computed(() => store.displayLabel)

  /**
   * Pilih context baru (setelah login atau ganti periode)
   * @param {Object} contextData - { schoolId, schoolName, academicYear, semester }
   */
  async function selectContext(contextData) {
    try {
      // 1. Engine validasi dan bangun context (Pure JS)
      const { context, events } = engine.establishContext(contextData)

      // 2. Store simpan state
      store.setContext(context)

      // 3. Persist ke localStorage agar survive refresh
      localStorage.setItem('operational_context', JSON.stringify(context.toJSON()))

      // 4. Dispatch events ke seluruh modul
      events.forEach((event) => eventDispatcher.dispatch(event))

      return { success: true }
    } catch (error) {
      return { success: false, error: error.message }
    }
  }

  /**
   * Restore context dari localStorage (saat app init / refresh)
   */
  function restoreContext() {
    try {
      const saved = localStorage.getItem('operational_context')
      if (!saved) return false

      const data = JSON.parse(saved)
      const { context } = engine.establishContext(data)
      store.setContext(context)
      return true
    } catch {
      store.clearContext()
      localStorage.removeItem('operational_context')
      return false
    }
  }

  /**
   * Clear context (logout)
   */
  function clearContext() {
    const { events } = engine.clearContext()
    store.clearContext()
    localStorage.removeItem('operational_context')
    events.forEach((event) => eventDispatcher.dispatch(event))
  }

  /**
   * Guard: Pastikan context sudah dipilih sebelum akses modul
   * Gunakan di router guard atau di awal composable modul
   */
  function requireContext() {
    const validation = engine.validateForTransaction(store.current)
    if (!validation.valid) {
      router.push({ name: 'context-selection' }) // Halaman pilih context
      throw new Error(validation.reason)
    }
    return store.current
  }

  return {
    current,
    isReady,
    displayLabel,
    selectContext,
    restoreContext,
    clearContext,
    requireContext,
  }
}

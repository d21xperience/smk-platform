// FILE: src/composables/useReadOnlyContext.js
// STATUS: NEW
// STATUS IMPLEMENTASI: COMPLETE

import { computed } from 'vue'
import { useContext } from './useContext.js'

export function useReadOnlyContext() {
  const { isReadOnlyContext, isHistoryModeActive, displayContext } = useContext()

  const readOnlyMessage = computed(() => {
    if (!isReadOnlyContext.value) {
      return ''
    }

    return 'Anda sedang melihat data historis. Mode ini read-only.'
  })

  return {
    isReadOnlyContext,
    isHistoryModeActive,
    displayContext,
    readOnlyMessage,
  }
}

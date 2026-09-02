// FILE: src/composables/useReadOnlyLoader.js
// STATUS: NEW
// STATUS IMPLEMENTASI: COMPLETE

import { ref } from 'vue'
import { useContext } from './useContext.js'

export function useReadOnlyLoader() {
  const { isReadOnlyContext } = useContext()

  const historicalDataNotFound = ref(false)
  const isLoaderReady = ref(false)

  const loadWithHistoryGuard = async (loadFn, options = {}) => {
    const { onHistoricalNotFound = null, onCurrentModeLoad = null } = options

    historicalDataNotFound.value = false
    isLoaderReady.value = false

    try {
      if (isReadOnlyContext.value) {
        // Historical mode: try to load existing data only.
        // If the loader function is loadOrCreate, the Axios guard
        // will block the create part. For mock adapters, if it
        // creates empty data in memory, we detect and handle it.
        const result = await loadFn()

        // If the result indicates no existing data was found
        // (e.g., empty/new record), mark as not found
        if (!result || isEmptyHistoricalData(result)) {
          historicalDataNotFound.value = true

          if (onHistoricalNotFound) {
            await onHistoricalNotFound()
          }
        }

        isLoaderReady.value = true
        return result
      }

      // Current mode: normal loadOrCreate behavior
      const result = await loadFn()

      if (onCurrentModeLoad) {
        await onCurrentModeLoad(result)
      }

      isLoaderReady.value = true
      return result
    } catch (err) {
      isLoaderReady.value = true
      throw err
    }
  }

  const isEmptyHistoricalData = (data) => {
    if (!data) return true

    // Check common patterns for empty/new records
    if (data.status === 'draft' && !data.hasExistingData) {
      return true
    }

    return false
  }

  return {
    isReadOnlyContext,
    historicalDataNotFound,
    isLoaderReady,
    loadWithHistoryGuard,
  }
}

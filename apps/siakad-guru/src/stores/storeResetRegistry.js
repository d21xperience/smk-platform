const resetHandlers = new Map()

export function registerStoreReset(storeName, resetFn) {
  if (typeof resetFn !== 'function') {
    throw new Error(`Reset handler for "${storeName}" must be a function`)
  }
  resetHandlers.set(storeName, resetFn)
}

export function unregisterStoreReset(storeName) {
  resetHandlers.delete(storeName)
}

export function resetAllStores() {
  resetHandlers.forEach((resetFn, storeName) => {
    try {
      resetFn()
    } catch (e) {
      console.error(`[StoreResetRegistry] Failed to reset store "${storeName}":`, e)
    }
  })
}

export function getRegisteredStores() {
  return Array.from(resetHandlers.keys())
}

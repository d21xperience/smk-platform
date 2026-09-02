// FILE: src/stores/contextStore.js
// STATUS: MODIFY
// STATUS IMPLEMENTASI: COMPLETE

import { defineStore } from 'pinia'
import { contextService } from '../boot/services.js'
import { registerStoreReset } from './storeResetRegistry.js'
import {
  setRequestContext,
  clearRequestContext,
} from '../utils/requestContextRegistry.js'

const CONTEXT_STORAGE_KEY = 'sdp_context'

function persistCurrentContext(context) {
  if (!context) return

  localStorage.setItem(CONTEXT_STORAGE_KEY, JSON.stringify(context))
}

function readPersistedCurrentContext() {
  const saved = localStorage.getItem(CONTEXT_STORAGE_KEY)

  if (!saved) {
    return null
  }

  try {
    return JSON.parse(saved)
  } catch {
    return null
  }
}

function isStoredContextValid(context, user) {
  if (!context) return false

  const hasRequiredFields = Boolean(
    context.schoolId &&
    context.academicYearId &&
    context.semesterId &&
    context.userId &&
    context.role
  )

  if (!hasRequiredFields) return false

  if (user && context.userId !== user.id) {
    return false
  }

  return true
}

function syncStoreRequestContext(store) {
  if (!store.currentContext) {
    clearRequestContext()
    return
  }

  try {
    const requestContext = contextService.resolveRequestContext({
      isHistoryMode: store.isHistoryMode,
      currentContext: store.currentContext,
      historicalContext: store.historicalContext,
    })

    setRequestContext(requestContext)
  } catch {
    clearRequestContext()
  }
}

export const useContextStore = defineStore('context', {
  state: () => ({
    availableContexts: null,
    availableHistory: null,

    currentContext: null,
    historicalContext: null,

    isHistoryMode: false,

    loading: false,
    loadingHistory: false,

    error: null,
    historyError: null,
  }),

  getters: {
    isContextSelected: (state) => Boolean(state.currentContext),
    isContextReady: (state) => Boolean(state.currentContext),

    isHistoryModeActive: (state) => state.isHistoryMode && Boolean(state.historicalContext),
    isReadOnlyContext: (state) => state.isHistoryMode,

    displayContext: (state) => {
      if (state.isHistoryMode && state.historicalContext) {
        return state.historicalContext
      }

      return state.currentContext
    },
  },

  actions: {
    async loadAvailableContexts() {
      this.loading = true
      this.error = null

      try {
        this.availableContexts = await contextService.loadAvailableContexts()
      } catch (err) {
        this.error = err.message
      } finally {
        this.loading = false
      }
    },

    async selectContext({ schoolId, academicYearId, semesterId, user }) {
      if (!user) {
        throw new Error('User tidak tersedia untuk memilih context')
      }

      this.loading = true
      this.error = null

      try {
        const context = await contextService.selectContext({
          schoolId,
          academicYearId,
          semesterId,
          user,
        })

        this.currentContext = context
        this.historicalContext = null
        this.isHistoryMode = false

        persistCurrentContext(context)
        syncStoreRequestContext(this)

        return context
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.loading = false
      }
    },

    async ensureCurrentContext(user) {
      if (!user) {
        throw new Error('User tidak tersedia untuk memuat context')
      }

      if (this.isContextReady && this.currentContext.userId === user.id) {
        syncStoreRequestContext(this)
        return this.currentContext
      }

      const restoredContext = readPersistedCurrentContext()

      if (isStoredContextValid(restoredContext, user)) {
        this.currentContext = restoredContext
        this.historicalContext = null
        this.isHistoryMode = false

        syncStoreRequestContext(this)

        return this.currentContext
      }

      return await this.loadCurrentContext(user)
    },

    async loadCurrentContext(user) {
      if (!user) {
        throw new Error('User tidak tersedia untuk memuat current context')
      }

      this.loading = true
      this.error = null

      try {
        const context = await contextService.loadCurrentOperationalContext({ user })

        this.currentContext = context
        this.historicalContext = null
        this.isHistoryMode = false
        this.availableHistory = null

        persistCurrentContext(context)
        syncStoreRequestContext(this)

        return context
      } catch (err) {
        this.currentContext = null
        this.error = err.message
        clearRequestContext()
        throw err
      } finally {
        this.loading = false
      }
    },

    async loadAvailableHistory(user) {
      if (!user) {
        throw new Error('User tidak tersedia untuk memuat riwayat context')
      }

      this.loadingHistory = true
      this.historyError = null

      try {
        this.availableHistory = await contextService.loadAvailableHistory({ user })

        return this.availableHistory
      } catch (err) {
        this.historyError = err.message
        throw err
      } finally {
        this.loadingHistory = false
      }
    },

    async enterHistoricalContext({ user, academicYearId, semesterId }) {
      if (!user) {
        throw new Error('User tidak tersedia untuk masuk ke historical context')
      }

      this.loadingHistory = true
      this.historyError = null

      try {
        const historicalContext = await contextService.enterHistoricalContext({
          user,
          academicYearId,
          semesterId,
        })

        this.historicalContext = historicalContext
        this.isHistoryMode = true

        syncStoreRequestContext(this)

        return historicalContext
      } catch (err) {
        this.historyError = err.message
        throw err
      } finally {
        this.loadingHistory = false
      }
    },

    exitHistoricalContext() {
      this.historicalContext = null
      this.isHistoryMode = false
      this.historyError = null

      syncStoreRequestContext(this)
    },

    restoreContext() {
      const restoredContext = readPersistedCurrentContext()

      if (isStoredContextValid(restoredContext, null)) {
        this.currentContext = restoredContext
        this.historicalContext = null
        this.isHistoryMode = false

        syncStoreRequestContext(this)

        return
      }

      this.clearContext()
    },

    clearContext() {
      this.currentContext = null
      this.historicalContext = null
      this.isHistoryMode = false

      localStorage.removeItem(CONTEXT_STORAGE_KEY)
      clearRequestContext()
    },

    resetContextState() {
      this.availableContexts = null
      this.availableHistory = null

      this.currentContext = null
      this.historicalContext = null

      this.isHistoryMode = false

      this.loading = false
      this.loadingHistory = false

      this.error = null
      this.historyError = null

      localStorage.removeItem(CONTEXT_STORAGE_KEY)
      clearRequestContext()
    },
  },
})

registerStoreReset('context', () => {
  localStorage.removeItem(CONTEXT_STORAGE_KEY)

  const store = useContextStore()
  store.$reset()

  clearRequestContext()
})

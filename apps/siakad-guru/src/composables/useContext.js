// FILE: src/composables/useContext.js
// STATUS: MODIFY
// STATUS IMPLEMENTASI: COMPLETE

import { computed } from 'vue'
import { useContextStore } from '../stores/contextStore.js'
import { useAuth } from './useAuth.js'

export function useContext() {
  const contextStore = useContextStore()
  const { currentUser, isAuthenticated } = useAuth()

  const isContextReady = computed(() => contextStore.isContextReady)
  const isContextSelected = computed(() => contextStore.isContextSelected)

  const currentContext = computed(() => contextStore.currentContext)
  const availableContexts = computed(() => contextStore.availableContexts)

  const historicalContext = computed(() => contextStore.historicalContext)
  const availableHistory = computed(() => contextStore.availableHistory)
  const isHistoryMode = computed(() => contextStore.isHistoryMode)
  const isHistoryModeActive = computed(() => contextStore.isHistoryModeActive)
  const isReadOnlyContext = computed(() => contextStore.isReadOnlyContext)
  const displayContext = computed(() => contextStore.displayContext)

  const isLoading = computed(() => contextStore.loading)
  const isLoadingHistory = computed(() => contextStore.loadingHistory)

  const error = computed(() => contextStore.error)
  const historyError = computed(() => contextStore.historyError)

  const canViewHistory = computed(() => {
    return isAuthenticated.value && Boolean(currentUser.value)
  })

  const loadAvailableContexts = async () => {
    if (!contextStore.availableContexts) {
      await contextStore.loadAvailableContexts()
    }
  }

  const restoreContext = () => {
    contextStore.restoreContext()
  }

  const selectContext = async ({ schoolId, academicYearId, semesterId }) => {
    await contextStore.selectContext({
      schoolId,
      academicYearId,
      semesterId,
      user: currentUser.value,
    })
  }

  const clearContext = () => {
    contextStore.clearContext()
  }

  const ensureCurrentContext = async () => {
    await contextStore.ensureCurrentContext(currentUser.value)
  }

  const loadCurrentContext = async () => {
    await contextStore.loadCurrentContext(currentUser.value)
  }

  const loadAvailableHistory = async () => {
    if (!contextStore.availableHistory) {
      await contextStore.loadAvailableHistory(currentUser.value)
    }
  }

  const enterHistoricalContext = async ({ academicYearId, semesterId }) => {
    await contextStore.enterHistoricalContext({
      user: currentUser.value,
      academicYearId,
      semesterId,
    })
  }

  const exitHistoricalContext = () => {
    contextStore.exitHistoricalContext()
  }

  return {
    isContextReady,
    isContextSelected,

    currentContext,
    availableContexts,

    historicalContext,
    availableHistory,
    isHistoryMode,
    isHistoryModeActive,
    isReadOnlyContext,
    displayContext,

    isLoading,
    isLoadingHistory,

    error,
    historyError,

    canViewHistory,

    loadAvailableContexts,
    restoreContext,
    selectContext,
    clearContext,

    ensureCurrentContext,
    loadCurrentContext,
    loadAvailableHistory,
    enterHistoricalContext,
    exitHistoricalContext,
  }
}

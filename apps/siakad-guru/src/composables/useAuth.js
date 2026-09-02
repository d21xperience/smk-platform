// FILE: src/composables/useAuth.js
// STATUS: MODIFY
// STATUS IMPLEMENTASI: COMPLETE

import { computed } from 'vue'
import { useAuthStore } from '../stores/authStore.js'
import { useRouter } from 'vue-router'
import { SESSION_STATE } from '../domain/auth/models/SessionState.js'
import { CREDENTIAL_TYPE } from '../domain/auth/models/CredentialType.js'

function checkPasskeyAvailability() {
  if (typeof window === 'undefined') return false
  if (!window.PublicKeyCredential) return false
  if (!navigator.credentials) return false

  return true
}

export function useAuth() {
  const authStore = useAuthStore()
  const router = useRouter()

  const sessionState = computed(() => authStore.sessionState)
  const isAuthenticated = computed(() => authStore.isAuthenticated)
  const isAnonymous = computed(() => authStore.isAnonymous)
  const isAuthenticating = computed(() => authStore.isAuthenticating)
  const isRestoring = computed(() => authStore.isRestoring)
  const isExpired = computed(() => authStore.isExpired)
  const isLoggingOut = computed(() => authStore.isLoggingOut)
  const isSessionTransitioning = computed(() => authStore.isSessionTransitioning)

  const currentUser = computed(() => authStore.currentUser)
  const isLoading = computed(() => authStore.loading)
  const error = computed(() => authStore.error)
  const initialized = computed(() => authStore.initialized)

  const isPasskeyAvailable = computed(() => checkPasskeyAvailability())

  const initialize = async () => {
    await authStore.initialize()
  }

  const login = async (username, password) => {
    await authStore.login(username, password)

    if (authStore.isAuthenticated) {
      router.push({ name: 'dashboard' })
    }
  }

  const loginWithPasskey = async (challenge) => {
    await authStore.loginWithPasskey(challenge)

    if (authStore.isAuthenticated) {
      router.push({ name: 'dashboard' })
    }
  }

  const logout = async () => {
    await authStore.logout()

    router.push({ name: 'login' })
  }

  const handleSessionExpired = () => {
    authStore.handleSessionExpired()

    router.push({ name: 'login' })
  }

  return {
    sessionState,
    isAuthenticated,
    isAnonymous,
    isAuthenticating,
    isRestoring,
    isExpired,
    isLoggingOut,
    isSessionTransitioning,

    currentUser,
    isLoading,
    error,
    initialized,

    isPasskeyAvailable,

    initialize,
    login,
    loginWithPasskey,
    logout,
    handleSessionExpired,

    SESSION_STATE,
    CREDENTIAL_TYPE,
  }
}

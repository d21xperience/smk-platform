// FILE: src/stores/authStore.js
// STATUS: MODIFY
// STATUS IMPLEMENTASI: COMPLETE

import { defineStore } from 'pinia'
import { authService, passkeyAuthService } from '../boot/services.js'
import { SESSION_STATE } from '../domain/auth/models/SessionState.js'
import { resetAllStores } from './storeResetRegistry.js'
import { setAuthToken, clearAuthToken } from '../utils/requestContextRegistry.js'

function extractAuthToken(session) {
  if (!session) return null

  return (
    session.token || session.accessToken || session.credential?.token || session.user?.token || null
  )
}

function syncAuthSession(session) {
  const token = extractAuthToken(session)

  if (token) {
    setAuthToken(token)
    return
  }

  clearAuthToken()
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    session: null,
    sessionState: SESSION_STATE.ANONYMOUS,
    initialized: false,
    loading: false,
    error: null,
  }),

  getters: {
    isAuthenticated: (state) => state.sessionState === SESSION_STATE.AUTHENTICATED,
    isAnonymous: (state) => state.sessionState === SESSION_STATE.ANONYMOUS,
    isAuthenticating: (state) => state.sessionState === SESSION_STATE.AUTHENTICATING,
    isRestoring: (state) => state.sessionState === SESSION_STATE.RESTORING,
    isExpired: (state) => state.sessionState === SESSION_STATE.EXPIRED,
    isLoggingOut: (state) => state.sessionState === SESSION_STATE.LOGGING_OUT,

    isSessionTransitioning: (state) => {
      return (
        state.sessionState === SESSION_STATE.AUTHENTICATING ||
        state.sessionState === SESSION_STATE.RESTORING ||
        state.sessionState === SESSION_STATE.LOGGING_OUT
      )
    },

    currentUser: (state) => {
      if (state.sessionState !== SESSION_STATE.AUTHENTICATED) {
        return null
      }

      return state.session ? state.session.user : null
    },
  },

  actions: {
    async initialize() {
      if (this.initialized) return

      this.sessionState = SESSION_STATE.RESTORING
      this.loading = true
      this.error = null

      try {
        const restoredSession = await authService.restoreSession()

        if (restoredSession) {
          this.session = restoredSession
          this.sessionState = SESSION_STATE.AUTHENTICATED
          syncAuthSession(restoredSession)
        } else {
          this.session = null
          this.sessionState = SESSION_STATE.ANONYMOUS
          clearAuthToken()
        }
      } catch (err) {
        this.session = null
        this.sessionState = SESSION_STATE.ANONYMOUS
        this.error = err.message
        clearAuthToken()
      } finally {
        this.loading = false
        this.initialized = true
      }
    },

    async login(username, password) {
      this.sessionState = SESSION_STATE.AUTHENTICATING
      this.loading = true
      this.error = null

      try {
        this.session = await authService.login(username, password)
        this.sessionState = SESSION_STATE.AUTHENTICATED

        syncAuthSession(this.session)

        return this.session
      } catch (err) {
        this.session = null
        this.sessionState = SESSION_STATE.ANONYMOUS
        this.error = err.message

        clearAuthToken()

        throw err
      } finally {
        this.loading = false
      }
    },

    async loginWithPasskey(challenge) {
      this.sessionState = SESSION_STATE.AUTHENTICATING
      this.loading = true
      this.error = null

      try {
        this.session = await passkeyAuthService.loginWithPasskey(challenge)
        this.sessionState = SESSION_STATE.AUTHENTICATED

        syncAuthSession(this.session)

        return this.session
      } catch (err) {
        this.session = null
        this.sessionState = SESSION_STATE.ANONYMOUS
        this.error = err.message

        clearAuthToken()

        throw err
      } finally {
        this.loading = false
      }
    },

    async logout() {
      this.sessionState = SESSION_STATE.LOGGING_OUT
      this.loading = true
      this.error = null

      try {
        const userId = this.session ? this.session.user.id : null

        await authService.logout(userId)

        this.session = null
        this.sessionState = SESSION_STATE.ANONYMOUS

        clearAuthToken()
        resetAllStores()
      } catch (err) {
        this.error = err.message

        this.session = null
        this.sessionState = SESSION_STATE.ANONYMOUS

        clearAuthToken()
        resetAllStores()
      } finally {
        this.loading = false
      }
    },

    handleSessionExpired() {
      this.session = null
      this.sessionState = SESSION_STATE.EXPIRED
      this.error = 'Session has expired. Please login again.'

      clearAuthToken()
    },

    resetAuthState() {
      this.session = null
      this.sessionState = SESSION_STATE.ANONYMOUS
      this.initialized = false
      this.loading = false
      this.error = null

      clearAuthToken()
    },
  },
})

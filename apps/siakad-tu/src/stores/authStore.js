// src/stores/authStore.js
// import { defineStore } from 'pinia'
// import { ref, computed } from 'vue'
// // import { loginApi, logoutApi } from 'src/services/authService' // <-- Nanti diaktifkan untuk Golang
// // import { mockLogin } from '@/services/mocks/authMock' // <-- Aktif saat ini untuk mock

// export const useAuthStore = defineStore('auth', () => {
//   const user = ref(JSON.parse(localStorage.getItem('user')) || null)
//   const token = ref(localStorage.getItem('token') || null)

//   const isAuthenticated = computed(() => !!token.value)
//   const showLoginDialog = ref(false)
//   const redirectUrl = ref(null)
//   const loginPromptMessage = ref(null)
//   const triggerLogin = (targetRoute = null, message = null) => {
//     if (targetRoute) redirectUrl.value = targetRoute
//     if (message) loginPromptMessage.value = message
//     showLoginDialog.value = true
//   }

//   const closeLoginDialog = () => {
//     showLoginDialog.value = false
//     redirectUrl.value = null // Reset juga saat dibatalkan
//     loginPromptMessage.value = null // <-- Reset pesan saat dialog ditutup
//   }

//   // const login = async (username, password) => {
//   //   try {
//   //     // --- FASE MOCK (Sekarang) ---
//   //     const response = await mockLogin(username, password)

//   //     // --- FASE BACKEND GOLANG (Nanti) ---
//   //     // const response = await loginApi({ username, password })

//   //     if (response.success) {
//   //       token.value = response.token
//   //       user.value = response.user

//   //       // Simpan ke localStorage agar persisten saat refresh
//   //       localStorage.setItem('token', response.token)
//   //       localStorage.setItem('user', JSON.stringify(response.user))

//   //       return true
//   //     }
//   //     return false
//   //   } catch (error) {
//   //     console.error('Login error:', error)
//   //     return false
//   //   }
//   // }

//   const logout = async () => {
//     // --- FASE BACKEND GOLANG (Nanti) ---
//     // await logoutApi()

//     token.value = null
//     user.value = null
//     localStorage.removeItem('token')
//     localStorage.removeItem('user')
//   }
//   // Fungsi untuk menyimpan dan membersihkan redirect URL
//   const setRedirectUrl = (url) => {
//     redirectUrl.value = url
//   }

//   const clearRedirectUrl = () => {
//     redirectUrl.value = null
//   }
//   return {
//     user,
//     token,
//     isAuthenticated,
//     showLoginDialog, // <-- Export state
//     redirectUrl, // <-- Export state
//     loginPromptMessage,
//     triggerLogin, // <-- Export action
//     closeLoginDialog, // <-- Export action
//     // login,
//     logout,
//     setRedirectUrl,
//     clearRedirectUrl,
//   }
// })

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authService } from '@/services/auth/serviceFactory.js'

/**
 * authStore — Single Source of Truth untuk authentication state.
 *
 * TANGGUNG JAWAB:
 * - Menyimpan token, user, isAuthenticated
 * - Mengorkestrasi login/logout melalui AuthService
 * - Persistensi ke localStorage
 * - Menyediakan computed properties untuk UI
 *
 * DILARANG:
 * - Axios/API call langsung
 * - Hardcoded mock data
 * - setTimeout untuk simulasi network
 * - Business logic (itu tugas Service/Engine)
 * - Mengetahui detail HTTP/infrastructure
 *
 * Flow:
 * Page → Composable → authStore → AuthService → AuthAdapter → Mock/Real
 */

// === CONSTANTS ===
const STORAGE_KEY_TOKEN = 'auth_token'
const STORAGE_KEY_USER = 'auth_user'

export const useAuthStore = defineStore('auth', () => {
  // === STATE ===
  const token = ref(localStorage.getItem(STORAGE_KEY_TOKEN) || null)
  const user = ref(JSON.parse(localStorage.getItem(STORAGE_KEY_USER) || 'null'))
  const isLoading = ref(false)
  const error = ref(null)

  // === COMPUTED ===
  const isAuthenticated = computed(() => !!token.value && !!user.value)
  const userRole = computed(() => user.value?.role || null)
  const userName = computed(() => user.value?.name || 'Guest')
  const userPermissions = computed(() => user.value?.permissions || [])

  // === ACTIONS ===

  /**
   * Persist token & user ke localStorage
   * @private
   */
  function _persistToStorage() {
    if (token.value) {
      localStorage.setItem(STORAGE_KEY_TOKEN, token.value)
    } else {
      localStorage.removeItem(STORAGE_KEY_TOKEN)
    }

    if (user.value) {
      localStorage.setItem(STORAGE_KEY_USER, JSON.stringify(user.value))
    } else {
      localStorage.removeItem(STORAGE_KEY_USER)
    }
  }

  /**
   * Login dengan username & password
   *
   * @param {Object} credentials
   * @param {string} credentials.username
   * @param {string} credentials.password
   * @returns {Promise<boolean>} true jika berhasil, false jika gagal
   */
  async function login(credentials) {
    isLoading.value = true
    error.value = null

    try {
      // Panggil AuthService (bukan langsung adapter)
      const result = await authService.login(credentials)

      if (result.success) {
        token.value = result.data.token
        user.value = result.data.user
        _persistToStorage()
        return true
      } else {
        error.value = result.error
        return false
      }
    } catch (err) {
      error.value = {
        code: 'UNEXPECTED_ERROR',
        message: err.message || 'Terjadi kesalahan saat login.',
        details: null,
      }
      return false
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Logout
   *
   * @returns {Promise<void>}
   */
  async function logout() {
    isLoading.value = true
    error.value = null

    try {
      // Panggil AuthService (bukan langsung adapter)
      await authService.logout()
    } catch (err) {
      // Logout tetap dilanjutkan meskipun Service error
      console.warn('[authStore] Logout Service error:', err)
    } finally {
      token.value = null
      user.value = null
      _persistToStorage()
      isLoading.value = false
    }
  }

  /**
   * Cek apakah user memiliki permission tertentu
   *
   * @param {string} permission
   * @returns {boolean}
   */
  function hasPermission(permission) {
    if (!user.value || !user.value.permissions) return false

    // Admin memiliki semua permission
    if (user.value.permissions.includes('*')) return true

    return user.value.permissions.includes(permission)
  }

  /**
   * Refresh token (untuk handle expired token)
   *
   * @returns {Promise<boolean>}
   */
  async function refreshToken() {
    if (!token.value) return false

    isLoading.value = true
    error.value = null

    try {
      // Panggil AuthService (bukan langsung adapter)
      const result = await authService.refreshToken(token.value)

      if (result.success) {
        token.value = result.data.token
        _persistToStorage()
        return true
      } else {
        // Jika refresh gagal, logout
        await logout()
        return false
      }
    } catch (err) {
      console.log(err)
      await logout()
      return false
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Load user dari token (saat app di-reload)
   *
   * @returns {Promise<boolean>}
   */
  async function loadUserFromToken() {
    if (!token.value) return false

    isLoading.value = true
    error.value = null

    try {
      // Panggil AuthService (bukan langsung adapter)
      const result = await authService.getCurrentUser(token.value)

      if (result.success) {
        user.value = result.data
        _persistToStorage()
        return true
      } else {
        // Token invalid, logout
        await logout()
        return false
      }
    } catch (err) {
      console.log(err)
      await logout()
      return false
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Clear error state
   */
  function clearError() {
    error.value = null
  }

  /**
   * Reset store (untuk testing)
   */
  function reset() {
    token.value = null
    user.value = null
    isLoading.value = false
    error.value = null
    _persistToStorage()
  }

  return {
    // State
    token,
    user,
    isLoading,
    error,

    // Computed
    isAuthenticated,
    userRole,
    userName,
    userPermissions,

    // Actions
    login,
    logout,
    hasPermission,
    refreshToken,
    loadUserFromToken,
    clearError,
    reset,
  }
})

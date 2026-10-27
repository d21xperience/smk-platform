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
  const showLoginDialog = ref(false)
  const redirectUrl = ref(null)
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
  function triggerLogin(targetRoute = null) {
    if (targetRoute) redirectUrl.value = targetRoute
    showLoginDialog.value = true
  }

  function closeLoginDialog() {
    showLoginDialog.value = false
    redirectUrl.value = null
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
    closeLoginDialog,
    triggerLogin,
    showLoginDialog,
    redirectUrl,
  }
})

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authService } from '@/services/auth/serviceFactory.js'

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

  // Operational context fields (diisi dari response /auth/login)
  const userSchoolId = computed(() => user.value?.school_id || '')
  const userAcademicPeriodId = computed(() => user.value?.academic_period_id || '')
  const userAcademicYear = computed(() => user.value?.academic_year || '')
  const userSemester = computed(() => user.value?.semester || 1)

  // === ACTIONS ===
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

  async function login(credentials) {
    isLoading.value = true
    error.value = null

    try {
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

  async function logout() {
    isLoading.value = true
    error.value = null

    try {
      await authService.logout()
    } catch (err) {
      console.warn('[authStore] Logout Service error:', err)
    } finally {
      token.value = null
      user.value = null
      _persistToStorage()
      isLoading.value = false
    }
  }

  function hasPermission(permission) {
    if (!user.value || !user.value.permissions) return false
    if (user.value.permissions.includes('*')) return true
    return user.value.permissions.includes(permission)
  }

  async function refreshToken() {
    if (!token.value) return false

    isLoading.value = true
    error.value = null

    try {
      const result = await authService.refreshToken(token.value)

      if (result.success) {
        token.value = result.data.token
        _persistToStorage()
        return true
      } else {
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

  async function loadUserFromToken() {
    if (!token.value) return false

    isLoading.value = true
    error.value = null

    try {
      const result = await authService.getCurrentUser(token.value)

      if (result.success) {
        user.value = result.data
        _persistToStorage()
        return true
      } else {
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
    userSchoolId,
    userAcademicPeriodId,
    userAcademicYear,
    userSemester,

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

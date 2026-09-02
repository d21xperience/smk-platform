// src/stores/authStore.js
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
// import { loginApi, logoutApi } from 'src/services/authService' // <-- Nanti diaktifkan untuk Golang
// import { mockLogin } from '@/services/mocks/authMock' // <-- Aktif saat ini untuk mock

export const useAuthStore = defineStore('auth', () => {
  const user = ref(JSON.parse(localStorage.getItem('user')) || null)
  const token = ref(localStorage.getItem('token') || null)

  const isAuthenticated = computed(() => !!token.value)
  const showLoginDialog = ref(false)
  const redirectUrl = ref(null)
  const loginPromptMessage = ref(null)
  const triggerLogin = (targetRoute = null, message = null) => {
    if (targetRoute) redirectUrl.value = targetRoute
    if (message) loginPromptMessage.value = message
    showLoginDialog.value = true
  }

  const closeLoginDialog = () => {
    showLoginDialog.value = false
    redirectUrl.value = null // Reset juga saat dibatalkan
    loginPromptMessage.value = null // <-- Reset pesan saat dialog ditutup
  }

  // const login = async (username, password) => {
  //   try {
  //     // --- FASE MOCK (Sekarang) ---
  //     const response = await mockLogin(username, password)

  //     // --- FASE BACKEND GOLANG (Nanti) ---
  //     // const response = await loginApi({ username, password })

  //     if (response.success) {
  //       token.value = response.token
  //       user.value = response.user

  //       // Simpan ke localStorage agar persisten saat refresh
  //       localStorage.setItem('token', response.token)
  //       localStorage.setItem('user', JSON.stringify(response.user))

  //       return true
  //     }
  //     return false
  //   } catch (error) {
  //     console.error('Login error:', error)
  //     return false
  //   }
  // }

  const logout = async () => {
    // --- FASE BACKEND GOLANG (Nanti) ---
    // await logoutApi()

    token.value = null
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }
  // Fungsi untuk menyimpan dan membersihkan redirect URL
  const setRedirectUrl = (url) => {
    redirectUrl.value = url
  }

  const clearRedirectUrl = () => {
    redirectUrl.value = null
  }
  return {
    user,
    token,
    isAuthenticated,
    showLoginDialog, // <-- Export state
    redirectUrl, // <-- Export state
    loginPromptMessage,
    triggerLogin, // <-- Export action
    closeLoginDialog, // <-- Export action
    // login,
    logout,
    setRedirectUrl,
    clearRedirectUrl,
  }
})

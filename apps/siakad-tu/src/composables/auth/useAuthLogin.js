// apps/siakad-tu/src/composables/auth/useAuthLogin.js

import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/authStore'

/**
 * Composable untuk halaman Login.
 *
 * Tanggung jawab:
 * - Mengelola form login (username, password)
 * - Menangani proses autentikasi via AuthStore
 * - Mengelola redirect setelah login berhasil
 * - Membaca query parameter `redirectTo` dari URL
 *
 * Alur:
 * 1. User datang dari landing page dengan query ?redirectTo=student-list
 * 2. User input kredensial
 * 3. Composable memanggil authStore.login()
 * 4. Jika sukses → redirect ke halaman yang dituju
 * 5. Jika gagal → tampilkan error
 */
export function useAuthLogin() {
  const router = useRouter()
  const route = useRoute()
  const authStore = useAuthStore()

  // === STATE ===
  const formData = ref({
    username: '',
    password: '',
  })

  const showPassword = ref(false)
  const isLoading = ref(false)
  const error = ref(null)
  const redirectTo = ref(null)

  // === COMPUTED ===
  const isAuthenticated = computed(() => authStore.isAuthenticated)
  const isFormValid = computed(() => {
    return formData.value.username.trim() !== '' &&
           formData.value.password.trim() !== ''
  })

  // === ACTIONS ===

  /**
   * Baca query parameter `redirectTo` dari URL.
   * Dipanggil saat component mounted.
   */
  function readRedirectParam() {
    const query = route.query
    if (query.redirectTo) {
      redirectTo.value = {
        name: query.redirectTo,
        params: {},
      }

      // Kumpulkan parameter tambahan dari query (jika ada)
      const reservedKeys = ['redirectTo']
      Object.keys(query).forEach(key => {
        if (!reservedKeys.includes(key)) {
          redirectTo.value.params[key] = query[key]
        }
      })

      console.log('[useAuthLogin] Redirect target:', redirectTo.value)
    }
  }

  /**
   * Handle submit form login.
   */
  async function handleLogin() {
    if (!isFormValid.value) return

    isLoading.value = true
    error.value = null

    try {
      // Panggil authStore untuk login
      const result = await authStore.login({
        username: formData.value.username,
        password: formData.value.password,
      })

      if (result) {
        console.log('[useAuthLogin] Login berhasil')

        // Redirect ke halaman yang dituju (jika ada)
        if (redirectTo.value) {
          console.log('[useAuthLogin] Redirecting to:', redirectTo.value)
          router.push(redirectTo.value)
        } else {
          // Default: redirect ke dashboard kesiswaan
          router.push({ name: 'landing-home' })
        }
      } else {
        // Login gagal
        error.value = authStore.error || {
          code: 'LOGIN_FAILED',
          message: 'Username atau password salah.',
        }
      }
    } catch (err) {
      console.error('[useAuthLogin] Login error:', err)
      error.value = {
        code: 'LOGIN_ERROR',
        message: err.message || 'Terjadi kesalahan saat login.',
      }
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Toggle visibility password.
   */
  function togglePassword() {
    showPassword.value = !showPassword.value
  }

  /**
   * Reset form.
   */
  function resetForm() {
    formData.value = { username: '', password: '' }
    error.value = null
  }

  /**
   * Navigasi ke halaman landing.
   */
  function goToLanding() {
    router.push({ name: 'landing-home' })
  }

  // === LIFECYCLE ===
  onMounted(() => {
    // 1. Cek apakah sudah login
    if (isAuthenticated.value) {
      // Sudah login, langsung redirect
      if (route.query.redirectTo) {
        router.push({ name: route.query.redirectTo })
      } else {
        router.push({ name: 'landing-home' })
      }
      return
    }

    // 2. Baca parameter redirect
    readRedirectParam()
  })

  // === RETURN ===
  return {
    // State
    formData,
    showPassword,
    isLoading,
    error,
    redirectTo,

    // Computed
    isAuthenticated,
    isFormValid,

    // Actions
    handleLogin,
    togglePassword,
    resetForm,
    goToLanding,
  }
}

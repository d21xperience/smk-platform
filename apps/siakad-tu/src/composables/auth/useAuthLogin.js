import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/authStore'
import { useContextStore } from '@/stores/contextStore'

/**
 * Composable untuk halaman Login.
 *
 * Alur:
 *   1. Baca query ?redirectTo=... (opsional).
 *   2. Submit → authStore.login().
 *   3. Jika sukses → contextStore.initialize() (mengisi periods + semester aktif).
 *   4. Redirect ke halaman tujuan / landing-home.
 */
export function useAuthLogin() {
  const router = useRouter()
  const route = useRoute()
  const authStore = useAuthStore()
  const contextStore = useContextStore()

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
  const isFormValid = computed(
    () => formData.value.username.trim() !== '' && formData.value.password.trim() !== '',
  )

  // === ACTIONS ===
  function readRedirectParam() {
    const query = route.query
    if (query.redirectTo) {
      redirectTo.value = {
        name: query.redirectTo,
        params: {},
      }

      const reservedKeys = ['redirectTo']
      Object.keys(query).forEach((key) => {
        if (!reservedKeys.includes(key)) {
          redirectTo.value.params[key] = query[key]
        }
      })

      console.log('[useAuthLogin] Redirect target:', redirectTo.value)
    }
  }

  async function handleLogin() {
    if (!isFormValid.value) return

    isLoading.value = true
    error.value = null

    try {
      const result = await authStore.login({
        username: formData.value.username,
        password: formData.value.password,
      })

      if (result) {
        console.log('[useAuthLogin] Login berhasil')

        // Isi operational context SEBELUM redirect. Jika gagal (mis. network),
        // tetap lanjut redirect; halaman tujuan bisa initialize ulang.
        try {
          await contextStore.initialize()
        } catch (ctxErr) {
          console.warn('[useAuthLogin] Context init warning:', ctxErr)
        }

        if (redirectTo.value) {
          router.push(redirectTo.value)
        } else {
          router.push({ name: 'landing-home' })
        }
      } else {
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

  function togglePassword() {
    showPassword.value = !showPassword.value
  }

  function resetForm() {
    formData.value = { username: '', password: '' }
    error.value = null
  }

  function goToLanding() {
    router.push({ name: 'landing-home' })
  }

  // === LIFECYCLE ===
  onMounted(() => {
    if (isAuthenticated.value) {
      if (route.query.redirectTo) {
        router.push({ name: route.query.redirectTo })
      } else {
        router.push({ name: 'landing-home' })
      }
      return
    }
    readRedirectParam()
  })

  return {
    formData,
    showPassword,
    isLoading,
    error,
    redirectTo,
    isAuthenticated,
    isFormValid,
    handleLogin,
    togglePassword,
    resetForm,
    goToLanding,
  }
}

import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null, // Instance User
    token: localStorage.getItem('token') || null,
    loading: false,
    error: null,
    initialized: false,
  }),

  getters: {
    isAuthenticated: (state) => !!state.token && !!state.user,
    userRole: (state) => state.user?.role || null,
    userName: (state) => state.user?.name || '',
  },

  actions: {
    /**
     * Injeksi AuthService dari luar.
     * Dipanggil saat boot aplikasi.
     * @param {import('src/services/AuthService').AuthService} service
     */
    setService(service) {
      this.authService = service
    },

    /**
     * Login dengan kredensial.
     * @param {{ username: string, password: string }} credentials
     */
    async login(credentials) {
      if (!this.authService) throw new Error('AuthService belum diinisialisasi')

      this.loading = true
      this.error = null
      try {
        const { user, token } = await this.authService.login(credentials)
        this.user = user
        this.token = token
        localStorage.setItem('token', token)
        return true
      } catch (err) {
        this.error = err.message || 'Login gagal'
        throw err
      } finally {
        this.loading = false
      }
    },

    /**
     * Cek apakah token masih valid (untuk guard atau refresh halaman).
     */
    async checkAuth() {
      if (!this.token) {
        this.initialized = true
        return
      }
      this.loading = true
      try {
        const user = await this.authService.me(this.token)
        this.user = user
      } catch {
        this.logout()
      } finally {
        this.loading = false
        this.initialized = true
      }
    },

    /**
     * Logout, hapus semua state dan token.
     */
    logout() {
      this.user = null
      this.token = null
      localStorage.removeItem('token')
    },
  },
})

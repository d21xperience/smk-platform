import { api } from '../../boot/axios.js'

export const authRealAdapter = {
  async login(credentials) {
    try {
      // Memanggil POST /api/v1/auth/login
      const res = await api.post('/auth/login', credentials)
      return {
        success: true,
        data: res.data.data
      }
    } catch (err) {
      console.error("Login Error Detail:", err.response || err)
      return {
        success: false,
        error: {
          code: 'LOGIN_FAILED',
          message: err.response?.data?.error?.message || 'Username atau password salah.',
        }
      }
    }
  },
  async logout() { return { success: true } },
  async getCurrentUser() { return { success: false } },
  async refreshToken() { return { success: false } }
}

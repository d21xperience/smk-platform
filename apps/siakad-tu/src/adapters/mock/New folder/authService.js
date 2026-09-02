import { api } from 'boot/axios'

export const authService = {
  /**
   * Mengirim kredensial login ke backend
   * @param {Object} credentials - { username, password }
   */
  async login(credentials) {
    const response = await api.post('/auth/login', credentials)
    return response.data // Mengembalikan objek wrapper { meta, data: { token, user } }
  },

  /**
   * Memvalidasi token aktif saat aplikasi pertama kali dimuat
   */
  async verifyCurrentSession() {
    const response = await api.get('/auth/me')
    return response.data.data // Mengembalikan data profil user terbaru
  },

  /**
   * Logout dari server jika diperlukan pencatatan log audit di backend Golang
   */
  async logout() {
    return await api.post('/auth/logout')
  },
}

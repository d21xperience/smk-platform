import { api as apiClient } from '@/boot/axios'
import { User } from '@/models/User'

export const authApiAdapter = {
  async login({ username, password }) {
    const { data } = await apiClient.post('/auth/login', { username, password })
    return {
      user: new User(data.user),
      token: data.token,
    }
  },

  async me(token) {
    const { data } = await apiClient.get('/auth/me', {
      headers: { Authorization: `Bearer ${token}` },
    })
    return new User(data)
  },
}

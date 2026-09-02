// @/adapters/mock/auth.mock.js
import { User } from '@/models/User'

const dummyUser = new User({ id: 1, name: 'Deden Moh Jaenudin', role: 'guru', nip: '123' })

export const authMockAdapter = {
  async login({ username, password }) {
    // Simulasi delay
    await new Promise((resolve) => setTimeout(resolve, 500))
    if (username === 'denmailforprogramming@gmail.com' && password === '123') {
      return {
        user: dummyUser,
        token: 'mock-jwt-token',
      }
    }
    throw new Error('Username atau password salah')
  },

  async me(token) {
    await new Promise((resolve) => setTimeout(resolve, 200))
    if (token === 'mock-jwt-token') return dummyUser
    throw new Error('Unauthorized')
  },

  async cek() {},
}

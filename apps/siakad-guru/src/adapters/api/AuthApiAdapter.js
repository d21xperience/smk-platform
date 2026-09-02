import { BaseApiAdapter } from './BaseApiAdapter.js'

export class AuthApiAdapter extends BaseApiAdapter {
  constructor() {
    super('/auth')
  }

  async login(username, password) {
    const data = await this.post('/login', { username, password })
    return {
      token: data.token,
      expiresAt: data.expiresAt,
      user: {
        id: data.user.id,
        name: data.user.name,
        email: data.user.email,
        role: data.user.role,
        permissions: data.user.permissions,
      },
    }
  }

  async logout() {
    await this.post('/logout')
    return { success: true }
  }
}

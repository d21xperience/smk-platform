export class AuthService {
  /**
   * @param {Object} adapter - Objek dengan method login() dan me()
   */
  constructor(adapter) {
    this.adapter = adapter
  }

  /**
   * Login dengan kredensial
   * @param {{ username: string, password: string }} credentials
   * @returns {Promise<{ user: import('src/models/User').User, token: string }>}
   */
  async login(credentials) {
    // Bisa tambahkan validasi, logging, atau transformasi di sini
    return this.adapter.login(credentials)
  }

  /**
   * Dapatkan data user dari token
   * @param {string} token
   * @returns {Promise<import('src/models/User').User>}
   */
  async me(token) {
    return this.adapter.me(token)
  }
}

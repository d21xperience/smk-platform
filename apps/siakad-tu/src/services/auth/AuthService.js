/**
 * AuthService — Application/Use-Case layer untuk authentication.
 *
 * TANGGUNG JAWAB:
 * - Orkestrasi use case authentication (login, logout, refresh, getCurrentUser)
 * - Memanggil Adapter untuk operasi data
 * - Melakukan application-level transformation jika diperlukan
 * - Menjaga agar Store tidak mengetahui detail infrastructure
 *
 * DILARANG:
 * - Menggunakan Vue/Pinia
 * - Mengakses component/DOM
 * - Mengelola UI notification
 * - Bergantung langsung pada implementasi Mock/API tertentu
 * - Melakukan HTTP request langsung (itu tugas Adapter)
 *
 * Flow:
 * Store → AuthService → AuthAdapter → Mock/Real Implementation
 */

export class AuthService {
  /**
   * @param {Object} dependencies
   * @param {Object} dependencies.adapter - AuthAdapter (Mock atau Real)
   */
  constructor({ adapter }) {
    if (!adapter) {
      throw new Error('AuthService: adapter wajib di-inject.')
    }

    this.adapter = adapter
  }

  /**
   * Helper: format error response yang konsisten
   * @private
   */
  _formatError(source, error) {
    return {
      success: false,
      data: null,
      error: {
        code: error.code || `${source.toUpperCase()}_ERROR`,
        message: error.message || 'Terjadi kesalahan.',
        source,
        details: error.details || null,
      },
    }
  }

  /**
   * USE CASE: Login
   *
   * Alur:
   * 1. Validasi input dasar (optional, bisa juga di Store)
   * 2. Panggil adapter.login()
   * 3. Return hasil dengan format konsisten
   *
   * @param {Object} credentials
   * @param {string} credentials.username
   * @param {string} credentials.password
   * @returns {Promise<{ success, data: { token, user }, error }>}
   */
  async login(credentials) {
    // Validasi input dasar (application-level)
    if (!credentials || !credentials.username || !credentials.password) {
      return this._formatError('validation', {
        code: 'CREDENTIALS_REQUIRED',
        message: 'Username dan password wajib diisi.',
        details: { missingFields: [] },
      })
    }

    // Panggil adapter
    const adapterResult = await this.adapter.login(credentials)

    if (!adapterResult.success) {
      return this._formatError('adapter', adapterResult.error)
    }

    // Application-level transformation (jika diperlukan)
    // Contoh: normalize user data, add default values, dll
    const { token, user } = adapterResult.data

    return {
      success: true,
      data: {
        token,
        user: {
          ...user,
          // Tambahkan field default jika perlu
          lastLogin: new Date().toISOString(),
        },
      },
      error: null,
    }
  }

  /**
   * USE CASE: Logout
   *
   * @returns {Promise<{ success, data: { message }, error }>}
   */
  async logout() {
    const adapterResult = await this.adapter.logout()

    if (!adapterResult.success) {
      return this._formatError('adapter', adapterResult.error)
    }

    return {
      success: true,
      data: adapterResult.data,
      error: null,
    }
  }

  /**
   * USE CASE: Get Current User
   *
   * @param {string} token
   * @returns {Promise<{ success, data: user, error }>}
   */
  async getCurrentUser(token) {
    if (!token) {
      return this._formatError('validation', {
        code: 'TOKEN_REQUIRED',
        message: 'Token wajib diisi.',
      })
    }

    const adapterResult = await this.adapter.getCurrentUser(token)

    if (!adapterResult.success) {
      return this._formatError('adapter', adapterResult.error)
    }

    return {
      success: true,
      data: adapterResult.data,
      error: null,
    }
  }

  /**
   * USE CASE: Refresh Token
   *
   * @param {string} token
   * @returns {Promise<{ success, data: { token }, error }>}
   */
  async refreshToken(token) {
    if (!token) {
      return this._formatError('validation', {
        code: 'TOKEN_REQUIRED',
        message: 'Token wajib diisi.',
      })
    }

    const adapterResult = await this.adapter.refreshToken(token)

    if (!adapterResult.success) {
      return this._formatError('adapter', adapterResult.error)
    }

    return {
      success: true,
      data: adapterResult.data,
      error: null,
    }
  }

  /**
   * USE CASE: Check Permission
   *
   * @param {string} token
   * @param {string} permission
   * @returns {Promise<{ success, data: { hasPermission }, error }>}
   */
  async checkPermission(token, permission) {
    if (!token || !permission) {
      return this._formatError('validation', {
        code: 'PARAMS_REQUIRED',
        message: 'Token dan permission wajib diisi.',
      })
    }

    const adapterResult = await this.adapter.checkPermission(token, permission)

    if (!adapterResult.success) {
      return this._formatError('adapter', adapterResult.error)
    }

    return {
      success: true,
      data: adapterResult.data,
      error: null,
    }
  }
}

// === MOCK DATA ===
const mockUsers = [
  {
    userId: 'user-001',
    username: 'admin',
    password: 'admin123',
    name: 'Administrator',
    email: 'admin@smk.local',
    role: 'admin',
    permissions: ['*'],
  },
  {
    userId: 'user-002',
    username: 'tu',
    password: 'tu123',
    name: 'Staff Tata Usaha',
    email: 'tu@smk.local',
    role: 'tu',
    permissions: [
      'student.view',
      'student.create',
      'student.edit',
      'student.delete',
      'mutation.view',
      'mutation.create',
      'mutation.approve',
      'finance.view',
      'finance.create',
      'finance.payment',
      'attendance.view',
      'attendance.create',
      'assessment.view',
      'assessment.create',
    ],
  },
  {
    userId: 'user-003',
    username: 'guru',
    password: 'guru123',
    name: 'Guru Mata Pelajaran',
    email: 'guru@smk.local',
    role: 'guru',
    permissions: ['attendance.view', 'attendance.create', 'assessment.view', 'assessment.create'],
  },
]
// In-memory token store (simulasi JWT)
const activeTokens = new Map()
// === HELPER METHODS ===
/**
Simulasi latency network
@private
*/
async function _simulateLatency(ms = 50) {
  return new Promise((resolve) => setTimeout(resolve, ms))
}
/**
Helper: buat error response
@private
*/
function _error(code, message, details = null) {
  return {
    success: false,
    data: null,
    error: { code, message, details },
  }
}
/**
Helper: buat success response
@private
*/
function _success(data) {
  return { success: true, data, error: null }
}
/**
Helper: strip password dari user object
@private
*/
function _sanitizeUser(user) {
  if (!user) return null
  const { password, ...sanitized } = user
  console.log(password)
  return sanitized
}
/**
Helper: generate mock JWT token
@private
*/
function _generateToken(userId) {
  const timestamp = Date.now()
  const random = Math.random().toString(36).substring(2, 15)
  return `mock-jwt-${userId}-${timestamp}-${random}`
}
// === CONTRACT IMPLEMENTATION ===
export const authMockAdapter = {
  /**
Login dengan username & password
@param {Object} credentials
@param {string} credentials.username
@param {string} credentials.password
@returns {Promise<{ success, data: { token, user }, error }>}
*/
  async login(credentials) {
    await _simulateLatency(300)
    const { username, password } = credentials || {}

    // Validasi input
    if (!username || !password) {
      return _error('VALIDATION_ERROR', 'Username dan password wajib diisi.', { missingFields: [] })
    }

    // Cari user
    const user = mockUsers.find((u) => u.username === username && u.password === password)

    if (!user) {
      return _error('INVALID_CREDENTIALS', 'Username atau password salah.', null)
    }

    // Generate token
    const token = _generateToken(user.userId)

    // Simpan token ke active store
    activeTokens.set(token, {
      userId: user.userId,
      issuedAt: Date.now(),
      expiresAt: Date.now() + 24 * 60 * 60 * 1000, // 24 jam
    })

    return _success({
      token,
      user: _sanitizeUser(user),
    })
  },
  /**
Logout
@returns {Promise<{ success, data: { message }, error }>}
*/
  async logout() {
    await _simulateLatency(100)
    // Clear semua token (dalam real implementation, hanya clear token spesifik)
    activeTokens.clear()

    return _success({
      message: 'Logout berhasil',
    })
  },
  /**
Get current user berdasarkan token
@param {string} token
@returns {Promise<{ success, data: user, error }>}
*/
  async getCurrentUser(token) {
    await _simulateLatency(50)
    if (!token) {
      return _error('UNAUTHORIZED', 'Token tidak valid.', null)
    }

    const tokenData = activeTokens.get(token)
    if (!tokenData) {
      return _error('INVALID_TOKEN', 'Token tidak ditemukan atau sudah expired.', null)
    }

    // Cek expired
    if (Date.now() > tokenData.expiresAt) {
      activeTokens.delete(token)
      return _error('TOKEN_EXPIRED', 'Token sudah expired.', null)
    }

    const user = mockUsers.find((u) => u.userId === tokenData.userId)
    if (!user) {
      return _error('USER_NOT_FOUND', 'User tidak ditemukan.', null)
    }

    return _success(_sanitizeUser(user))
  },
  /**
Refresh token
@param {string} token
@returns {Promise<{ success, data: { token }, error }>}
*/
  async refreshToken(token) {
    await _simulateLatency(100)
    if (!token) {
      return _error('UNAUTHORIZED', 'Token tidak valid.', null)
    }

    const tokenData = activeTokens.get(token)
    if (!tokenData) {
      return _error('INVALID_TOKEN', 'Token tidak ditemukan.', null)
    }

    // Hapus token lama
    activeTokens.delete(token)

    // Generate token baru
    const newToken = _generateToken(tokenData.userId)
    activeTokens.set(newToken, {
      userId: tokenData.userId,
      issuedAt: Date.now(),
      expiresAt: Date.now() + 24 * 60 * 60 * 1000,
    })

    return _success({ token: newToken })
  },
  /**
Check apakah user memiliki permission tertentu
@param {string} token
@param {string} permission
@returns {Promise<{ success, data: { hasPermission }, error }>}
*/
  async checkPermission(token, permission) {
    await _simulateLatency(30)
    if (!token) {
      return _error('UNAUTHORIZED', 'Token tidak valid.', null)
    }

    const tokenData = activeTokens.get(token)
    if (!tokenData) {
      return _error('INVALID_TOKEN', 'Token tidak ditemukan.', null)
    }

    const user = mockUsers.find((u) => u.userId === tokenData.userId)
    if (!user) {
      return _error('USER_NOT_FOUND', 'User tidak ditemukan.', null)
    }

    // Admin memiliki semua permission
    const hasPermission = user.permissions.includes('*') || user.permissions.includes(permission)

    return _success({ hasPermission })
  },
  /**
Reset mock data (khusus development/testing)
*/
  async __resetMockData() {
    activeTokens.clear()
    return _success({ reset: true })
  },
}

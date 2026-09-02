export class AuthMockAdapter {
  async login(username, password) {
    await new Promise((resolve) => setTimeout(resolve, 800)) // Simulate network latency

    // Mock validation logic
    if (username === 'guru' && password === 'password') {
      return {
        token: 'mock-jwt-token-12345-teacher',
        expiresAt: new Date(Date.now() + 24 * 60 * 60 * 1000).toISOString(), // 24 hours
        user: {
          id: 'USR-001',
          name: 'Budi Santoso, S.Pd.',
          email: 'budi.santoso@schooldigital.id',
          role: 'teacher',
          permissions: [
            'attendance.create',
            'attendance.submit',
            'assessment.create',
            'journal.write',
            'inventory.report',
          ],
        },
      }
    }

    if (username === 'wali' && password === 'password') {
      return {
        token: 'mock-jwt-token-67890-homeroom',
        expiresAt: new Date(Date.now() + 24 * 60 * 60 * 1000).toISOString(),
        user: {
          id: 'USR-002',
          name: 'Siti Aminah, M.Pd.',
          email: 'siti.aminah@schooldigital.id',
          role: 'homeroom',
          permissions: [
            'attendance.create',
            'attendance.submit',
            'assessment.create',
            'journal.write',
            'inventory.report',
            'student.progress.view',
            'homeroom.billing.read',
            'homeroom.student.read',
            'homeroom.progress.read',
            'homeroom.attendance.read',
            'homeroom.journal.read',
            'homeroom.communication.read',
            'homeroom.report.print',
          ],
        },
      }
    }

    throw new Error('Username atau password salah')
  }
  async loginWithPasskey(challenge) {
    await new Promise((resolve) => setTimeout(resolve, 1000))

    // Mock passkey authentication
    // challenge is used to simulate unique auth request
    const challengeSuffix = challenge ? challenge.substring(0, 8) : 'default'

    return {
      token: `mock-passkey-token-${challengeSuffix}-${Date.now()}`,
      expiresAt: new Date(Date.now() + 24 * 60 * 60 * 1000).toISOString(),
      user: {
        id: 'USR-001',
        name: 'Budi Santoso, S.Pd.',
        email: 'budi.santoso@schooldigital.id',
        role: 'teacher',
        permissions: [
          'attendance.create',
          'attendance.submit',
          'assessment.create',
          'journal.write',
          'inventory.report',
          'teacher_document.manage_own',
        ],
      },
    }
  }

  async validateToken(token) {
    await new Promise((resolve) => setTimeout(resolve, 200))

    if (!token) {
      return { valid: false, reason: 'No token provided' }
    }

    if (token.startsWith('mock-')) {
      return { valid: true, reason: null }
    }

    return { valid: false, reason: 'Invalid token' }
  }
  async logout() {
    await new Promise((resolve) => setTimeout(resolve, 300))
    return { success: true }
  }
}

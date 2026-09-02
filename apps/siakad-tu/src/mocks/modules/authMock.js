export function registerAuthMocks(mock) {
  // 1. POST /auth/login
  mock.onPost('/auth/login').reply((config) => {
    const { username, password } = JSON.parse(config.data)

    // Simulasi akun admin & guru sesuai kebutuhan sistem
    if (username === 'admin' && password === 'password123') {
      return [
        200,
        {
          meta: { status: 'success', message: 'Login berhasil' },
          data: {
            token: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.mockTokenAdmin123',
            user: { id: 1, name: 'Deden M.J.', role: 'admin', email: 'admin@siakad.sch.id' },
          },
        },
      ]
    }

    if (username === 'guru' && password === 'password123') {
      return [
        200,
        {
          meta: { status: 'success', message: 'Login berhasil' },
          data: {
            token: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.mockTokenGuru456',
            user: { id: 2, name: 'Rizka Amalia', role: 'guru', email: 'rizka@siakad.sch.id' },
          },
        },
      ]
    }

    // Gagal login: skema respons error standar Golang
    return [
      401,
      {
        meta: { status: 'error', message: 'Kombinasi username atau password salah.' },
      },
    ]
  })

  // 2. GET /auth/me (Verifikasi Token Aktif)
  mock.onGet('/auth/me').reply((config) => {
    const authHeader = config.headers.Authorization

    if (authHeader === 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.mockTokenAdmin123') {
      return [
        200,
        {
          data: { id: 1, name: 'Deden M.J.', role: 'admin', email: 'admin@siakad.sch.id' },
        },
      ]
    }

    return [403, { meta: { status: 'error', message: 'Sesi kedaluwarsa atau token tidak valid' } }]
  })
}

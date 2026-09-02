// export const mockAuthService = {
//   async login(credentials) {
//     // Simulasi delay network backend Golang
//     await new Promise((resolve) => setTimeout(resolve, 800))

//     // Simulasi validasi sederhana
//     if (credentials.username === 'admin' && credentials.password === 'admin123') {
//       return {
//         status: 200,
//         data: {
//           token: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...',
//           user: {
//             name: 'Budi Santoso, S.Pd.',
//             role: 'Administrator',
//             school: 'SMA Negeri 1 Jakarta',
//           },
//         },
//       }
//     } else {
//       throw { response: { status: 401, data: { message: 'Username atau password salah!' } } }
//     }
//   },
// }
// src/services/mocks/authMock.js

// Fungsi helper untuk meniru delay jaringan (network latency)
const delay = (ms) => new Promise((resolve) => setTimeout(resolve, ms))

export const mockLogin = async (username, password) => {
  // Simulasi delay 800ms agar terasa seperti request ke server asli
  await delay(800)

  // Mock validasi kredensial (Nanti logika ini ada di Golang Backend)
  if (username === 'admin' && password === 'admin123') {
    return {
      success: true,
      token: 'mock-jwt-token-admin-xyz123',
      user: {
        id: 1,
        name: 'Administrator Sekolah',
        role: 'admin',
        email: 'admin@smk.sch.id',
        avatar: 'https://cdn.quasar.dev/img/boy-avatar.png',
      },
    }
  }

  if (username === 'guru' && password === 'guru123') {
    return {
      success: true,
      token: 'mock-jwt-token-guru-abc456',
      user: {
        id: 2,
        name: 'Budi Santoso, S.Pd',
        role: 'guru',
        email: 'guru@smk.sch.id',
        avatar: 'https://cdn.quasar.dev/img/boy-avatar.png',
      },
    }
  }

  // Jika kredensial salah
  return {
    success: false,
    message: 'Username atau password tidak valid.',
  }
}

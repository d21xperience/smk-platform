import MockAdapter from 'axios-mock-adapter'
import { registerAuthMocks } from './modules/authMock.js'
import { registerLandingMocks } from './modules/landingMock.js'

export function initAxiosMocks(axiosInstance) {
  // delayResponse: 500 mensimulasikan latency jaringan riil backend Golang selama 500ms
  const mock = new MockAdapter(axiosInstance, { delayResponse: 500 })

  // Registrasi sub-modul mock berdasarkan domain bisnis
  registerAuthMocks(mock)
  registerLandingMocks(mock)

  // Biarkan request eksternal yang tidak cocok dengan mock tetap lolos (misal: fetch asset/maps)
  mock.onAny().passThrough()
}

import MockAdapter from 'axios-mock-adapter'
import { setupSiakadMocks } from './siakadMock'

export const setupMockInterceptor = (api) => {
  const mock = new MockAdapter(api, { delayResponse: 300 }) // delay global opsional

  setupSiakadMocks(mock)
  // (Opsional) fallback untuk route yang tidak terdefinisi
  mock.onAny().passThrough()
}

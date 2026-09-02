// FILE: src/boot/axios.js
// STATUS: MODIFY
// STATUS IMPLEMENTASI: COMPLETE

import axios from 'axios'
import { boot } from 'quasar/wrappers'
import { getRequestHeaders, getRequestMeta } from '../utils/requestContextRegistry.js'
import { ensureRequestAllowed } from '../adapters/http/RequestGuard.js'

const api = axios.create({
  baseURL: import.meta.env.QCLI_API_BASE_URL || '/api/v1',
  timeout: 10000,
})

api.interceptors.request.use((config) => {
  const requestMeta = getRequestMeta()

  try {
    ensureRequestAllowed({
      method: config.method,
      url: config.url,
      mode: requestMeta.mode,
      context: requestMeta.context,
    })
  } catch (err) {
    return Promise.reject(err)
  }

  const headers = getRequestHeaders()

  if (!config.headers) {
    config.headers = {}
  }

  Object.entries(headers).forEach(([key, value]) => {
    if (typeof config.headers.set === 'function') {
      config.headers.set(key, value)
      return
    }

    config.headers[key] = value
  })

  return config
})

api.interceptors.response.use(
  (response) => response,
  (error) => Promise.reject(error),
)

export { api }

export default boot(({ app }) => {
  app.config.globalProperties.$api = api
})

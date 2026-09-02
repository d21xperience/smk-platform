// src/composables/useApi.js
import { ref } from 'vue'
import { api } from 'boot/axios'

export function useApi() {
  const loading = ref(false)
  const error = ref(null)

  // =====================
  // Generic request handler
  // =====================
  const request = async (method, url, payload = null) => {
    loading.value = true
    error.value = null

    try {
      const config = {
        method,
        url,
        ...(payload && ['post', 'put', 'patch'].includes(method)
          ? { data: payload }
          : { params: payload }),
      }

      const { data } = await api(config)
      return { data, success: true }
    } catch (err) {
      // Normalisasi pesan error dari server
      const message =
        err.response?.data?.message ||
        err.response?.data?.error ||
        err.message ||
        'Terjadi kesalahan.'

      error.value = message
      return { data: null, success: false, message }
    } finally {
      loading.value = false
    }
  }

  // Shorthand methods
  const get = (url, params) => request('get', url, params)
  const post = (url, payload) => request('post', url, payload)
  const put = (url, payload) => request('put', url, payload)
  const patch = (url, payload) => request('patch', url, payload)
  const remove = (url) => request('delete', url)

  return { loading, error, get, post, put, patch, remove }
}

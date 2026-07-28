// src/composables/useApi.js
import { api } from '@/boot/axios'
import { useQuasar } from 'quasar'
import { ref } from 'vue'

export function useApi() {
  const $q = useQuasar()
  const loading = ref(false)
  const error = ref(null)

  /**
   * Fungsi wrapper untuk memanggil API.
   * Menangani state loading, error, dan notifikasi secara terpusat.
   *
   * @param {Function} apiCall - Fungsi yang mengembalikan Promise axios (misal: () => api.get('/url'))
   * @param {Object} options - Opsi tambahan (misal: { silent: true } untuk menyembunyikan notifikasi error)
   * @returns {Promise<any>} - Mengembalikan data dari response (response.data)
   */
  const request = async (apiCall, options = {}) => {
    loading.value = true
    error.value = null

    try {
      const response = await apiCall()
      // Mengembalikan hanya data payload, bukan seluruh objek response axios
      return response.data
    } catch (err) {
      error.value = err

      // Tampilkan notifikasi error kecuali jika opsi 'silent' diaktifkan
      if (!options.silent) {
        const errorMessage = err.response?.data?.message || 'Terjadi kesalahan pada server.'
        $q.notify({
          type: 'negative',
          message: errorMessage,
          position: 'top',
        })
      }

      // Lempar error agar bisa ditangkap (catch) di komponen/composable pemanggil jika diperlukan
      throw err
    } finally {
      loading.value = false
    }
  }

  // WAJIB: Pastikan 'request' di-return di sini agar bisa dipakai di useAttendance.js
  return {
    api,
    loading,
    error,
    request,
  }
}

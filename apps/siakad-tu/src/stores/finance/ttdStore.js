import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useAuthStore } from '@/stores/authStore'
import { ttdAdapter } from '@/adapters/api/keuangan/ttdAdapter'

export const useTTDStore = defineStore('keuangan-ttd', () => {
  const authStore = useAuthStore()
  const currentUser = computed(() => authStore.user || authStore.currentUser || {})
  const userID = computed(() => currentUser.value.id || currentUser.value.user_id || '')

  const pinStatus = ref(null)
  const loading = ref(false)
  const error = ref(null)

  // Status PIN: { is_setup, is_locked, current_attempt, max_attempt, locked_until? }
  async function getPinStatus() {
    loading.value = true
    error.value = null
    try {
      const res = await ttdAdapter.getPinStatus()
      pinStatus.value = res.data
      return res
    } catch (e) {
      error.value = e.response?.data?.error || e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  async function setupPin(newPin) {
    loading.value = true
    error.value = null
    try {
      const res = await ttdAdapter.setupPin({ new_pin: newPin })
      await getPinStatus()
      return res
    } catch (e) {
      error.value = e.response?.data?.error || e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  /**
   * Verifikasi PIN
   * @param {string} pin — 6 digit
   * @returns {Object} — { valid, current_attempt, remaining_attempt, locked_until?, message? }
   */
  async function verifyPin(pin) {
    loading.value = true
    error.value = null
    try {
      const res = await ttdAdapter.verifyPin({ pin })
      await getPinStatus()
      return res.data
    } catch (e) {
      error.value = e.response?.data?.error || e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  async function changePin(oldPin, newPin) {
    loading.value = true
    error.value = null
    try {
      const res = await ttdAdapter.changePin({ old_pin: oldPin, new_pin: newPin })
      await getPinStatus()
      return res
    } catch (e) {
      error.value = e.response?.data?.error || e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  async function unlockPin(targetUserID = null) {
    loading.value = true
    error.value = null
    try {
      const res = await ttdAdapter.unlockPin(targetUserID || userID.value)
      await getPinStatus()
      return res
    } catch (e) {
      error.value = e.response?.data?.error || e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  function reset() {
    pinStatus.value = null
    error.value = null
  }

  return {
    userID,
    pinStatus, loading, error,
    getPinStatus, setupPin, verifyPin, changePin, unlockPin,
    reset,
  }
})

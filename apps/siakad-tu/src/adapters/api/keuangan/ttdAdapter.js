import { api } from '@/boot/axios'

const pick = (res) => res.data

// Header X-User-ID: di-set global via axios interceptor (boot/axios.js)
export const ttdAdapter = {
  // ============ PIN TTD ============
  async setupPin(payload) {
    return pick(await api.post('/keuangan/ttd/pin/setup', payload))
  },
  async verifyPin(payload) {
    return pick(await api.post('/keuangan/ttd/pin/verify', payload))
  },
  async changePin(payload) {
    return pick(await api.put('/keuangan/ttd/pin/change', payload))
  },
  async unlockPin(userID) {
    return pick(await api.post('/keuangan/ttd/pin/unlock', { user_id: userID }))
  },
  async getPinStatus() {
    return pick(await api.get('/keuangan/ttd/pin/status'))
  },

  // ============ VERIFY PUBLIC (QR Code) ============
  async verifyByHash(hash) {
    return pick(await api.get(`/keuangan/verify/${hash}`))
  },
}

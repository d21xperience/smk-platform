import { api } from '@/boot/axios'

const pick = (res) => res.data

export const transaksiAdapter = {
  // ============ TAGIHAN ============
  async listTagihan(params) {
    return pick(await api.get('/keuangan/tagihan', { params }))
  },
  async getTagihan(id, params) {
    return pick(await api.get(`/keuangan/tagihan/${id}`, { params }))
  },

  // ============ PENERIMAAN IURAN (auto BKM) ============
  async listPenerimaan(params) {
    return pick(await api.get('/keuangan/penerimaan', { params }))
  },
  async getPenerimaan(noBKM, params) {
    return pick(await api.get(`/keuangan/penerimaan/${noBKM}`, { params }))
  },
  async createPenerimaan(payload, params) {
    return pick(await api.post('/keuangan/penerimaan', payload, { params }))
  },

  // ============ BUKU KAS UMUM ============
  async listBKU(params) {
    return pick(await api.get('/keuangan/bku', { params }))
  },
  async getBKUSaldo(params) {
    return pick(await api.get('/keuangan/bku/saldo', { params }))
  },
}

import { api } from '@/boot/axios'

const pick = (res) => res.data

export const pengeluaranAdapter = {
  // ============ PENGAJUAN ============
  async listPengajuan(params) {
    return pick(await api.get('/keuangan/pengajuan', { params }))
  },
  async getPengajuan(id, params) {
    return pick(await api.get(`/keuangan/pengajuan/${id}`, { params }))
  },
  async createPengajuan(payload, params) {
    return pick(await api.post('/keuangan/pengajuan', payload, { params }))
  },
  async verifikasiBendahara(id, payload, params) {
    return pick(await api.put(`/keuangan/pengajuan/${id}/verifikasi`, payload, { params }))
  },
  async approvalKepsek(id, payload, params) {
    return pick(await api.put(`/keuangan/pengajuan/${id}/approval`, payload, { params }))
  },

  // ============ PENGELUARAN (auto BKK) ============
  async listPengeluaran(params) {
    return pick(await api.get('/keuangan/pengeluaran', { params }))
  },
  async getPengeluaran(noBKK, params) {
    return pick(await api.get(`/keuangan/pengeluaran/${noBKK}`, { params }))
  },
}

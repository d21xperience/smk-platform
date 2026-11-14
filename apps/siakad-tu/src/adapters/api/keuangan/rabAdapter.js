import { api } from '@/boot/axios'

const pick = (res) => res.data

export const rabAdapter = {
  async listRab(params) {
    return pick(await api.get('/keuangan/rab', { params }))
  },
  async getRab(id, params) {
    return pick(await api.get(`/keuangan/rab/${id}`, { params }))
  },
  async createRab(payload, params) {
    return pick(await api.post('/keuangan/rab', payload, { params }))
  },
  async updateRab(id, payload, params) {
    return pick(await api.put(`/keuangan/rab/${id}`, payload, { params }))
  },
  async sahkanRab(id, payload, params) {
    return pick(await api.put(`/keuangan/rab/${id}/sahkan`, payload, { params }))
  },
  async deleteRab(id, params) {
    return pick(await api.delete(`/keuangan/rab/${id}`, { params }))
  },
  async getRealisasi(tahun, params) {
    return pick(await api.get(`/keuangan/rab/tahun/${tahun}/realisasi`, { params }))
  },
}

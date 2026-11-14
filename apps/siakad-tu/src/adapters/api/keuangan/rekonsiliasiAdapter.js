import { api } from '@/boot/axios'

const pick = (res) => res.data

export const rekonsiliasiAdapter = {
  async listRekonsiliasi(params) {
    return pick(await api.get('/keuangan/rekonsiliasi', { params }))
  },
  async getRekonsiliasi(id, params) {
    return pick(await api.get(`/keuangan/rekonsiliasi/${id}`, { params }))
  },
  async jalankanRekonsiliasi(payload, params) {
    return pick(await api.post('/keuangan/rekonsiliasi', payload, { params }))
  },
  async deleteRekonsiliasi(id, params) {
    return pick(await api.delete(`/keuangan/rekonsiliasi/${id}`, { params }))
  },
}

import { api } from '@/boot/axios'

const pick = (res) => res.data

export const penutupanAdapter = {
  async listPenutupan(params) {
    return pick(await api.get('/keuangan/penutupan', { params }))
  },
  async getPenutupan(id, params) {
    return pick(await api.get(`/keuangan/penutupan/${id}`, { params }))
  },
  async previewSnapshot(params) {
    return pick(await api.get('/keuangan/penutupan/preview', { params }))
  },
  async tutupBuku(payload, params) {
    return pick(await api.post('/keuangan/penutupan', payload, { params }))
  },
  async signPenutupan(id, payload, params) {
    return pick(await api.post(`/keuangan/penutupan/${id}/sign`, payload, { params }))
  },
}

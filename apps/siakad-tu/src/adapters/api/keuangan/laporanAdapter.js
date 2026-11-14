import { api } from '@/boot/axios'

const pick = (res) => res.data

export const laporanAdapter = {
  async getNeraca(params) {
    return pick(await api.get('/keuangan/laporan/neraca', { params }))
  },
  async getPosisiKeuangan(params) {
    return pick(await api.get('/keuangan/laporan/posisi-keuangan', { params }))
  },
  async getLRA(params) {
    return pick(await api.get('/keuangan/laporan/lra', { params }))
  },
  async getDashboard(params) {
    return pick(await api.get('/keuangan/laporan/dashboard', { params }))
  },
  async getRingkasanBulan(params) {
    return pick(await api.get('/keuangan/laporan/ringkasan-bulan', { params }))
  },
}

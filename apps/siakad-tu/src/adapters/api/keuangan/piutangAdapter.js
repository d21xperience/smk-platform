import { api } from '@/boot/axios'

const pick = (res) => res.data

export const piutangAdapter = {
  async listPiutang(params) {
    return pick(await api.get('/keuangan/piutang', { params }))
  },
  async listPiutangPerSiswa(params) {
    return pick(await api.get('/keuangan/piutang/per-siswa', { params }))
  },
  async getPiutangByNIS(nis, params) {
    return pick(await api.get(`/keuangan/piutang/siswa/${nis}`, { params }))
  },
  async getAgingReport(params) {
    return pick(await api.get('/keuangan/piutang/aging', { params }))
  },
}

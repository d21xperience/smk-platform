import { api } from '@/boot/axios'

// Helper: extract data dari response backend { success, data, pagination }
const pick = (res) => res.data

export const masterAdapter = {
  // ============ JENIS IURAN ============
  async listJenisIuran(params) {
    return pick(await api.get('/keuangan/master/jenis-iuran', { params }))
  },
  async getJenisIuran(id, params) {
    return pick(await api.get(`/keuangan/master/jenis-iuran/${id}`, { params }))
  },
  async createJenisIuran(payload, params) {
    return pick(await api.post('/keuangan/master/jenis-iuran', payload, { params }))
  },
  async updateJenisIuran(id, payload, params) {
    return pick(await api.put(`/keuangan/master/jenis-iuran/${id}`, payload, { params }))
  },
  async deleteJenisIuran(id, params) {
    return pick(await api.delete(`/keuangan/master/jenis-iuran/${id}`, { params }))
  },

  // ============ TARIF IURAN ============
  async listTarifIuran(params) {
    return pick(await api.get('/keuangan/master/tarif-iuran', { params }))
  },
  async getTarifIuran(id, params) {
    return pick(await api.get(`/keuangan/master/tarif-iuran/${id}`, { params }))
  },
  async createTarifIuran(payload, params) {
    return pick(await api.post('/keuangan/master/tarif-iuran', payload, { params }))
  },
  async updateTarifIuran(id, payload, params) {
    return pick(await api.put(`/keuangan/master/tarif-iuran/${id}`, payload, { params }))
  },
  async deleteTarifIuran(id, params) {
    return pick(await api.delete(`/keuangan/master/tarif-iuran/${id}`, { params }))
  },

  // ============ KODE REKENING (COA) ============
  async listKodeRekening(params) {
    return pick(await api.get('/keuangan/master/kode-rekening', { params }))
  },
  async getKodeRekening(id, params) {
    return pick(await api.get(`/keuangan/master/kode-rekening/${id}`, { params }))
  },
  async createKodeRekening(payload, params) {
    return pick(await api.post('/keuangan/master/kode-rekening', payload, { params }))
  },
  async updateKodeRekening(id, payload, params) {
    return pick(await api.put(`/keuangan/master/kode-rekening/${id}`, payload, { params }))
  },
  async deleteKodeRekening(id, params) {
    return pick(await api.delete(`/keuangan/master/kode-rekening/${id}`, { params }))
  },

  // ============ REKENING BANK ============
  async listRekeningBank(params) {
    return pick(await api.get('/keuangan/master/rekening-bank', { params }))
  },
  async getRekeningBank(id, params) {
    return pick(await api.get(`/keuangan/master/rekening-bank/${id}`, { params }))
  },
  async createRekeningBank(payload, params) {
    return pick(await api.post('/keuangan/master/rekening-bank', payload, { params }))
  },
  async updateRekeningBank(id, payload, params) {
    return pick(await api.put(`/keuangan/master/rekening-bank/${id}`, payload, { params }))
  },
  async deleteRekeningBank(id, params) {
    return pick(await api.delete(`/keuangan/master/rekening-bank/${id}`, { params }))
  },
}

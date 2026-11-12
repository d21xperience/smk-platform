// apps/siakad-tu/src/adapters/api/mutationRealAdapter.js

import { api } from '@/boot/axios.js'

const mapBackendToMutation = (d) => {
  if (!d) return null
  return {
    id: d.ID,
    schoolId: d.SchoolID,
    academicPeriodId: d.AcademicPeriodID,
    studentId: d.StudentID,
    mutationType: d.MutationType, // 'MASUK' | 'KELUAR'
    nis: d.NIS,
    nisn: d.NISN,
    nama: d.Nama,
    jenisKelamin: d.JenisKelamin,
    tempatLahir: d.TempatLahir,
    tanggalLahir: d.TanggalLahir,
    asalSekolah: d.AsalSekolah,
    tujuanSekolah: d.TujuanSekolah,
    alasan: d.Alasan,
    dokumenUrl: d.DokumenURL,
    status: d.Status, // 'PENDING' | 'APPROVED' | 'REJECTED'
    catatan: d.Catatan,
    createdAt: d.CreatedAt,
    updatedAt: d.UpdatedAt,
  }
}

const handleApiError = (error) => ({
  success: false,
  data: null,
  error: {
    code: error.response?.status || 'API_ERROR',
    message: error.response?.data?.error || error.message || 'Terjadi kesalahan.',
  },
})

const handleSuccess = (data) => ({ success: true, data, error: null })

export const mutationRealAdapter = {
  async getMutations(context, filters = {}) {
    try {
      const params = {}
      if (filters.status) params.status = filters.status
      if (filters.type) params.type = filters.type

      const response = await api.get('/mutations', { params })
      const items = response.data.data || []
      return handleSuccess({
        items: items.map(mapBackendToMutation),
        total: response.data.total || items.length,
      })
    } catch (error) { return handleApiError(error) }
  },

  async getMutationById(mutationId) {
    try {
      const response = await api.get(`/mutations/${mutationId}`)
      return handleSuccess(mapBackendToMutation(response.data))
    } catch (error) { return handleApiError(error) }
  },

  async createMutation(payload) {
    try {
      const response = await api.post('/mutations', payload)
      return handleSuccess(mapBackendToMutation(response.data))
    } catch (error) { return handleApiError(error) }
  },

  async approveMutation(mutationId, catatan) {
    try {
      const response = await api.post(`/mutations/${mutationId}/approve`, { catatan })
      return handleSuccess(response.data)
    } catch (error) { return handleApiError(error) }
  },

  async rejectMutation(mutationId, catatan) {
    try {
      const response = await api.post(`/mutations/${mutationId}/reject`, { catatan })
      return handleSuccess(response.data)
    } catch (error) { return handleApiError(error) }
  },
}

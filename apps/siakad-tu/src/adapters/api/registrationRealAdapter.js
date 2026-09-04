// apps/siakad-tu/src/adapters/api/registrationRealAdapter.js

import { api } from '@/boot/axios.js'

/**
 * Helper: Map Backend Registration DTO ke Frontend Contract
 */
const mapBackendToRegistration = (backendDto) => {
  if (!backendDto) return null

  return {
    registrationId: backendDto.registration_id || backendDto.id,
    studentId: backendDto.student_id || null, // Null jika belum di-approve
    nisn: backendDto.nisn,
    nis: backendDto.nis,
    fullName: backendDto.full_name || backendDto.fullName,
    gender: backendDto.gender,
    birthDate: backendDto.birth_date || backendDto.birthDate,
    status: backendDto.status, // 'DRAFT', 'SUBMITTED', 'VERIFIED', 'APPROVED', 'REJECTED'
    schoolId: backendDto.school_id,
    periodId: backendDto.period_id,
    createdAt: backendDto.created_at || backendDto.createdAt,
    updatedAt: backendDto.updated_at || backendDto.updatedAt,
    rejectionReason: backendDto.rejection_reason || null,
  }
}

/**
 * Helper: Konversi ApplicationError ke format Service yang konsisten
 */
const handleApiError = (error) => ({
  success: false,
  data: null,
  error: {
    code: error.code || 'API_ERROR',
    message: error.message || 'Terjadi kesalahan pada server.',
    details: error.details || null,
  },
})

const handleSuccess = (data) => ({
  success: true,
  data,
  error: null,
})

export const registrationRealAdapter = {
  /**
   * Buat draft registrasi baru
   */
  async createDraft(command, context) {
    try {
      const payload = { ...command, schoolId: context.schoolId, periodId: context.periodId }
      const response = await api.post('/registrations/draft', payload)
      return handleSuccess(mapBackendToRegistration(response.data))
    } catch (error) {
      return handleApiError(error)
    }
  },

  /**
   * Submit draft untuk diverifikasi
   */
  async submitRegistration(registrationId) {
    try {
      const response = await api.post(`/registrations/${registrationId}/submit`)
      return handleSuccess(mapBackendToRegistration(response.data))
    } catch (error) {
      return handleApiError(error)
    }
  },

  /**
   * Verifikasi data registrasi (oleh admin/tu)
   */
  async verifyRegistration(registrationId, command) {
    try {
      const response = await api.post(`/registrations/${registrationId}/verify`, command)
      return handleSuccess(mapBackendToRegistration(response.data))
    } catch (error) {
      return handleApiError(error)
    }
  },

  /**
   * Approve registrasi (Trigger auto-create Student)
   */
  async approveRegistration(registrationId, command) {
    try {
      const payload = {
        classId: command.classId,
        enrollmentDate: command.enrollmentDate || new Date().toISOString(),
      }
      const response = await api.post(`/registrations/${registrationId}/approve`, payload)
      return handleSuccess(mapBackendToRegistration(response.data))
    } catch (error) {
      return handleApiError(error)
    }
  },

  /**
   * Reject registrasi
   */
  async rejectRegistration(registrationId, command) {
    try {
      const payload = { reason: command.reason }
      const response = await api.post(`/registrations/${registrationId}/reject`, payload)
      return handleSuccess(mapBackendToRegistration(response.data))
    } catch (error) {
      return handleApiError(error)
    }
  },

  /**
   * Ambil detail registrasi by ID
   */
  async getRegistrationById(registrationId) {
    try {
      const response = await api.get(`/registrations/${registrationId}`)
      return handleSuccess(mapBackendToRegistration(response.data))
    } catch (error) {
      return handleApiError(error)
    }
  },

  /**
   * Ambil daftar registrasi dengan filter
   */
  async getRegistrations(context, filters = {}) {
    try {
      const params = { schoolId: context.schoolId, periodId: context.periodId, ...filters }
      const response = await api.get('/registrations', { params })

      const mappedData = response.data.items
        ? { ...response.data, items: response.data.items.map(mapBackendToRegistration) }
        : response.data

      return handleSuccess(mappedData)
    } catch (error) {
      return handleApiError(error)
    }
  },
}

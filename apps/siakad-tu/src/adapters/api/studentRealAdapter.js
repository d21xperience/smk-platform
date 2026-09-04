// apps/siakad-tu/src/adapters/api/studentRealAdapter.js

import { api } from '@/boot/axios.js'
// import {
//   ApplicationError,
//   NotFoundError,
//   ConflictError,
//   ValidationError,
//   BusinessRuleError,
//   NetworkError,
//   TimeoutError,
// } from '@/domain/errors/index.js'

/**
 * Helper: Map Backend DTO ke Frontend Contract (Anti-Corruption Layer)
 * Pastikan mapping ini sesuai dengan struktur DTO yang dikembalikan TU-Core.
 */
const mapBackendToStudent = (backendDto) => {
  if (!backendDto) return null

  return {
    studentId: backendDto.student_id || backendDto.id || backendDto.studentId,
    nisn: backendDto.nisn,
    nis: backendDto.nis,
    fullName: backendDto.full_name || backendDto.fullName,
    gender: backendDto.gender,
    status: backendDto.status,
    enrollments: backendDto.enrollments || [],
    createdAt: backendDto.created_at || backendDto.createdAt,
    updatedAt: backendDto.updated_at || backendDto.updatedAt,
  }
}

/**
 * Helper: Konversi ApplicationError (dari Axios Interceptor)
 * ke format error yang diharapkan oleh Command/Query Service.
 */
const handleApiError = (error) => {
  return {
    success: false,
    data: null,
    error: {
      code: error.code || 'API_ERROR',
      message: error.message || 'Terjadi kesalahan pada server.',
      details: error.details || null,
    },
  }
}

const handleSuccess = (data) => ({
  success: true,
  data,
  error: null,
})

export const studentRealAdapter = {
  async registerStudent(command, context) {
    try {
      const payload = { ...command, schoolId: context.schoolId, periodId: context.periodId }
      const response = await api.post('/students', payload)
      return handleSuccess(mapBackendToStudent(response.data))
    } catch (error) {
      return handleApiError(error)
    }
  },

  async enrollStudent(command, context) {
    try {
      const payload = {
        enrollmentId: command.enrollmentId,
        classId: command.classId,
        enrollmentDate: command.enrollmentDate,
        periodId: context.periodId,
      }
      const response = await api.post(`/students/${command.studentId}/enroll`, payload)
      return handleSuccess(mapBackendToStudent(response.data))
    } catch (error) {
      return handleApiError(error)
    }
  },

  async updateStudentProfile(command) {
    try {
      const payload = {
        address: command.address,
        contactInfo: command.contactInfo,
        guardianInfo: command.guardianInfo,
      }
      const response = await api.put(`/students/${command.studentId}`, payload)
      return handleSuccess(mapBackendToStudent(response.data))
    } catch (error) {
      return handleApiError(error)
    }
  },

  async graduateStudent(command) {
    try {
      const payload = { graduationDate: command.graduationDate }
      const response = await api.post(`/students/${command.studentId}/graduate`, payload)
      return handleSuccess(mapBackendToStudent(response.data))
    } catch (error) {
      return handleApiError(error)
    }
  },

  async transferStudent(command) {
    try {
      const payload = {
        targetSchool: command.targetSchool,
        reason: command.reason,
        transferDate: command.transferDate,
      }
      const response = await api.post(`/students/${command.studentId}/transfer`, payload)
      return handleSuccess(mapBackendToStudent(response.data))
    } catch (error) {
      return handleApiError(error)
    }
  },

  async getStudentById(studentId) {
    try {
      const response = await api.get(`/students/${studentId}`)
      return handleSuccess(mapBackendToStudent(response.data))
    } catch (error) {
      return handleApiError(error)
    }
  },

  async getStudents(context, filters = {}) {
    try {
      console.log('[StudentRealAdapter] getStudents called with:', { context, filters })

      // Pastikan kita mengirim snake_case sesuai ekspektasi backend Go
      const params = {
        school_id: context.schoolId,
        period_id: context.periodId,
      }

      // Hanya tambahkan filter jika nilainya ada (hindari 'search=' yang kosong)
      if (filters.status) {
        params.status = filters.status
      }
      if (filters.search && filters.search.trim() !== '') {
        params.search = filters.search.trim()
      }
      if (filters.page) {
        params.page = filters.page
      }
      if (filters.limit) {
        params.limit = filters.limit
      }

      console.log('[StudentRealAdapter] Sending params to backend:', params)

      const response = await api.get('/students', { params })

      const responseData = response.data || response
      const items = Array.isArray(responseData)
        ? responseData
        : responseData.data || responseData.items || []
      const total = responseData.meta?.total || responseData.total || items.length
      const page = responseData.meta?.page || filters.page || 1
      const limit = responseData.meta?.limit || filters.limit || 20

      console.log('[StudentRealAdapter] getStudents success:', { items: items.length, total })

      return handleSuccess({
        items: items.map(mapBackendToStudent),
        total,
        page,
        limit,
        totalPages: Math.ceil(total / limit),
      })
    } catch (error) {
      console.error('[StudentRealAdapter] getStudents error:', error)
      return handleApiError(error)
    }
  },

  async checkNisnAvailability(nisn, context) {
    try {
      const response = await api.get('/students/check-nisn', {
        params: { nisn, schoolId: context.schoolId },
      })
      return handleSuccess(response.data)
    } catch (error) {
      return handleApiError(error)
    }
  },

  async checkNisAvailability(nis, context) {
    try {
      const response = await api.get('/students/check-nis', {
        params: { nis, schoolId: context.schoolId },
      })
      return handleSuccess(response.data)
    } catch (error) {
      return handleApiError(error)
    }
  },
}

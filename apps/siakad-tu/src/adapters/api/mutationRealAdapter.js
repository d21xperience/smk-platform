// apps/siakad-tu/src/adapters/api/mutationRealAdapter.js

import { api } from '@/boot/axios.js'

/**
 * Helper: Map Backend Mutation DTO ke Frontend Contract
 */
const mapBackendToMutation = (backendDto) => {
  if (!backendDto) return null

  return {
    mutationId: backendDto.mutation_id || backendDto.id,
    studentId: backendDto.student_id,
    type: backendDto.type, // 'IN' | 'OUT'
    originSchool: backendDto.origin_school,
    destinationSchool: backendDto.destination_school,
    reason: backendDto.reason,
    effectiveDate: backendDto.effective_date || backendDto.effectiveDate,
    status: backendDto.status, // 'PENDING', 'APPROVED', 'REJECTED'
    approvalNote: backendDto.approval_note || null,
    createdAt: backendDto.created_at || backendDto.createdAt,
    updatedAt: backendDto.updated_at || backendDto.updatedAt,
    // Jika backend mengembalikan data siswa yang sudah diupdate sebagai side-effect
    updatedStudent: backendDto.updated_student
      ? {
          studentId: backendDto.updated_student.id,
          status: backendDto.updated_student.status,
          currentClass: backendDto.updated_student.current_class,
        }
      : null,
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

export const mutationRealAdapter = {
  /**
   * Ajukan Mutasi Masuk (Siswa dari luar masuk ke sekolah ini)
   */
  async requestMutationIn(command, context) {
    try {
      const payload = {
        ...command,
        schoolId: context.schoolId,
        periodId: context.periodId,
      }
      const response = await api.post('/mutations/in', payload)
      return handleSuccess(mapBackendToMutation(response.data))
    } catch (error) {
      return handleApiError(error)
    }
  },

  /**
   * Ajukan Mutasi Keluar (Siswa di sekolah ini pindah ke sekolah lain)
   */
  async requestMutationOut(command) {
    try {
      const payload = {
        studentId: command.studentId,
        destinationSchool: command.destinationSchool,
        reason: command.reason,
        effectiveDate: command.effectiveDate,
      }
      const response = await api.post('/mutations/out', payload)
      return handleSuccess(mapBackendToMutation(response.data))
    } catch (error) {
      return handleApiError(error)
    }
  },

  /**
   * Approve mutasi (Trigger side-effect: update status siswa, tutup enrollment lama)
   */
  async approveMutation(mutationId, command) {
    try {
      const payload = {
        approvalNote: command.approvalNote,
        targetClassId: command.targetClassId, // Khusus untuk mutasi masuk: masukkan ke kelas mana
      }
      const response = await api.post(`/mutations/${mutationId}/approve`, payload)
      return handleSuccess(mapBackendToMutation(response.data))
    } catch (error) {
      return handleApiError(error)
    }
  },

  /**
   * Reject mutasi
   */
  async rejectMutation(mutationId, command) {
    try {
      const payload = { rejectionReason: command.rejectionReason }
      const response = await api.post(`/mutations/${mutationId}/reject`, payload)
      return handleSuccess(mapBackendToMutation(response.data))
    } catch (error) {
      return handleApiError(error)
    }
  },

  /**
   * Ambil detail mutasi by ID
   */
  async getMutationById(mutationId) {
    try {
      const response = await api.get(`/mutations/${mutationId}`)
      return handleSuccess(mapBackendToMutation(response.data))
    } catch (error) {
      return handleApiError(error)
    }
  },

  /**
   * Ambil daftar mutasi dengan filter
   */
  async getMutations(context, filters = {}) {
    try {
      const params = { schoolId: context.schoolId, periodId: context.periodId, ...filters }
      const response = await api.get('/mutations', { params })

      const mappedData = response.data.items
        ? { ...response.data, items: response.data.items.map(mapBackendToMutation) }
        : response.data

      return handleSuccess(mappedData)
    } catch (error) {
      return handleApiError(error)
    }
  },
}

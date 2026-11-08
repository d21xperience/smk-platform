// apps/siakad-tu/src/adapters/api/studentRealAdapter.js

import { api } from '@/boot/axios.js'

const mapBackendToStudent = (backendDto) => {
  if (!backendDto) return null
  let mappedGender = backendDto.gender
  if (mappedGender === 'MALE') mappedGender = 'L'
  if (mappedGender === 'FEMALE') mappedGender = 'P'

  return {
    id: backendDto.id || backendDto.student_id,
    nisn: backendDto.nisn,
    nis: backendDto.nis,
    fullName: backendDto.name,
    gender: mappedGender,
    status: backendDto.status,
    className: backendDto.className,
    birthPlace: backendDto.birthPlace,
    birthDate: backendDto.birthDate,
    address: backendDto.address,
    phone: backendDto.phone,
    parentName: backendDto.parentName,
    parentPhone: backendDto.parentPhone,
    motherName: backendDto.motherName,
    enrollments: backendDto.enrollments || [],
  }
}

const handleApiError = (error) => ({
  success: false,
  data: null,
  error: {
    code: error.response?.status || 'API_ERROR',
    message: error.response?.data?.error || error.message || 'Terjadi kesalahan.',
    details: null,
  },
})

const handleSuccess = (data) => ({ success: true, data, error: null })

export const studentRealAdapter = {
  async getStudents(context, filters = {}) {
    try {
      const params = {}
      if (filters.status) params.status = filters.status
      if (filters.search) params.search = filters.search.trim()
      
      const response = await api.get('/students', { params })
      const responseData = response.data
      
      const items = responseData.data || []
      const total = responseData.total || items.length

      return handleSuccess({
        items: items.map(mapBackendToStudent),
        total,
        page: responseData.page || 1,
        limit: responseData.limit || 20,
        totalPages: responseData.totalPages || 1,
      })
    } catch (error) { return handleApiError(error) }
  },

  async getStudentById(studentId) {
    try {
      const response = await api.get(`/students/${studentId}`)
      return handleSuccess(mapBackendToStudent(response.data))
    } catch (error) { return handleApiError(error) }
  },

  // ... (fungsi registerStudent, updateStudentProfile, deleteStudent biarkan sama) ...

  async getStudentDocuments(studentId) {
    try {
      const response = await api.get(`/students/${studentId}/documents`)
      return handleSuccess(response.data.data || [])
    } catch (error) { return handleApiError(error) }
  },

  async uploadStudentDocument(studentId, formData) {
    try {
      const response = await api.post(`/students/${studentId}/documents`, formData, {
        headers: { 'Content-Type': 'multipart/form-data' }
      })
      return handleSuccess(response.data)
    } catch (error) { return handleApiError(error) }
  },

  async deleteStudentDocument(docId) {
    try {
      await api.delete(`/students/documents/${docId}`)
      return handleSuccess({ id: docId })
    } catch (error) { return handleApiError(error) }
  }
}

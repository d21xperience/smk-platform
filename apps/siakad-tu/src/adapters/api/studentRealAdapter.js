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
    parentPhone: backendDto.parentPhone, // <-- Tambahan
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

  async registerStudent(command) {
    try {
      const payload = {
        nis: command.nis || '',
        nisn: command.nisn || '',
        name: command.nama || command.fullName || '',
        gender: command.jenisKelamin === 'L' ? 'MALE' : 'FEMALE',
        birthPlace: command.tempatLahir || '',
        address: command.alamat || '',
        phone: command.noHp || '',
        email: '',
        parentName: command.parentName || 'Belum Diisi',
        parentPhone: command.parentPhone || '',
        parentEmail: '',
        classId: command.kelas || ''
      }

      if (command.tanggalLahir && command.tanggalLahir.trim() !== '') {
        payload.birthDate = command.tanggalLahir;
      }

      const response = await api.post('/students', payload)
      return handleSuccess(mapBackendToStudent(response.data))
    } catch (error) { return handleApiError(error) }
  },

  async updateStudentProfile(command) {
    try {
      const payload = {
        name: command.name,
        nis: command.nis,
        nisn: command.nisn,
        classId: command.classId,
        gender: command.gender,
        birthPlace: command.birthPlace,
        address: command.address,
        phone: command.phone,
        parentName: command.parentName,   // <-- Tambahan
        parentPhone: command.parentPhone, // <-- Tambahan
      }
      
      if (command.birthDate && command.birthDate.trim() !== '') {
        payload.birthDate = command.birthDate;
      }

      const response = await api.put(`/students/${command.studentId}`, payload)
      return handleSuccess(mapBackendToStudent(response.data))
    } catch (error) { return handleApiError(error) }
  },

  async deleteStudent(studentId) {
    try {
      await api.delete(`/students/${studentId}`)
      return handleSuccess({ id: studentId })
    } catch (error) { return handleApiError(error) }
  },
}

// apps/siakad-tu/src/adapters/mock/registrationMockAdapter.js

import { mockStorage } from './mockStorage.js'
import { idGenerator } from '../utils/idGenerator.js'

const _simulateLatency = (ms = 50) => new Promise((resolve) => setTimeout(resolve, ms))
const _error = (code, message, details = null) => ({
  success: false,
  data: null,
  error: { code, message, details },
})
const _success = (data) => ({ success: true, data, error: null })

export const registrationMockAdapter = {
  async createDraft(command, context) {
    await _simulateLatency()
    const registrationId = idGenerator.registrationId()

    const draftData = {
      registrationId,
      studentId: null,
      ...command,
      schoolId: context.schoolId,
      periodId: context.periodId,
      status: 'DRAFT',
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    }

    const registrations = mockStorage.read(context.schoolId, 'registrations')
    registrations.push(draftData)
    mockStorage.write(context.schoolId, 'registrations', registrations)

    return _success(draftData)
  },

  async submitRegistration(registrationId, context) {
    await _simulateLatency(100)
    const registrations = mockStorage.read(context.schoolId, 'registrations')
    const index = registrations.findIndex((r) => r.registrationId === registrationId)

    if (index === -1) return _error('NOT_FOUND', 'Data registrasi tidak ditemukan.')
    if (registrations[index].status !== 'DRAFT')
      return _error('INVALID_STATE', 'Hanya status DRAFT yang bisa di-submit.')

    registrations[index].status = 'SUBMITTED'
    registrations[index].updatedAt = new Date().toISOString()
    mockStorage.write(context.schoolId, 'registrations', registrations)

    return _success(registrations[index])
  },

  async verifyRegistration(registrationId, command, context) {
    await _simulateLatency(100)
    const registrations = mockStorage.read(context.schoolId, 'registrations')
    const index = registrations.findIndex((r) => r.registrationId === registrationId)

    if (index === -1) return _error('NOT_FOUND', 'Data registrasi tidak ditemukan.')
    if (registrations[index].status !== 'SUBMITTED')
      return _error('INVALID_STATE', 'Hanya status SUBMITTED yang bisa diverifikasi.')

    registrations[index].status = 'VERIFIED'
    registrations[index].updatedAt = new Date().toISOString()
    mockStorage.write(context.schoolId, 'registrations', registrations)

    return _success(registrations[index])
  },

  async approveRegistration(registrationId, command, context) {
    await _simulateLatency(200) // Simulasi proses auto-create student yang lebih lama
    const registrations = mockStorage.read(context.schoolId, 'registrations')
    const index = registrations.findIndex((r) => r.registrationId === registrationId)

    if (index === -1) return _error('NOT_FOUND', 'Data registrasi tidak ditemukan.')
    if (registrations[index].status !== 'VERIFIED')
      return _error('INVALID_STATE', 'Hanya status VERIFIED yang bisa di-approve.')

    // Simulasi side-effect: Auto-create Student ID
    const newStudentId = idGenerator.studentId()

    registrations[index].status = 'APPROVED'
    registrations[index].studentId = newStudentId
    registrations[index].updatedAt = new Date().toISOString()
    mockStorage.write(context.schoolId, 'registrations', registrations)

    return _success(registrations[index])
  },

  async rejectRegistration(registrationId, command, context) {
    await _simulateLatency(100)
    const registrations = mockStorage.read(context.schoolId, 'registrations')
    const index = registrations.findIndex((r) => r.registrationId === registrationId)

    if (index === -1) return _error('NOT_FOUND', 'Data registrasi tidak ditemukan.')

    registrations[index].status = 'REJECTED'
    registrations[index].rejectionReason = command.reason || 'Tidak memenuhi syarat.'
    registrations[index].updatedAt = new Date().toISOString()
    mockStorage.write(context.schoolId, 'registrations', registrations)

    return _success(registrations[index])
  },

  async getRegistrationById(registrationId, context) {
    await _simulateLatency(30)
    const registrations = mockStorage.read(context.schoolId, 'registrations')
    const data = registrations.find((r) => r.registrationId === registrationId)

    if (!data) return _error('NOT_FOUND', 'Data registrasi tidak ditemukan.')
    return _success(data)
  },

  async getRegistrations(context, filters = {}) {
    await _simulateLatency(100)
    let registrations = mockStorage.read(context.schoolId, 'registrations')

    if (filters.status) {
      registrations = registrations.filter((r) => r.status === filters.status)
    }
    if (filters.search) {
      const term = filters.search.toLowerCase()
      registrations = registrations.filter(
        (r) => r.fullName.toLowerCase().includes(term) || r.nisn.includes(term),
      )
    }

    const page = filters.page || 1
    const limit = filters.limit || 20
    const total = registrations.length
    const items = registrations.slice((page - 1) * limit, page * limit)

    return _success({ items, total, page, limit, totalPages: Math.ceil(total / limit) })
  },
}

// apps/siakad-tu/src/adapters/mock/mutationMockAdapter.js

import { mockStorage } from './mockStorage.js'
import { idGenerator } from '../utils/idGenerator.js'

const _simulateLatency = (ms = 50) => new Promise((resolve) => setTimeout(resolve, ms))
const _error = (code, message, details = null) => ({
  success: false,
  data: null,
  error: { code, message, details },
})
const _success = (data) => ({ success: true, data, error: null })

export const mutationMockAdapter = {
  async requestMutationIn(command, context) {
    await _simulateLatency(100)

    const mutationData = {
      mutationId: idGenerator.mutationId(),
      studentId: command.studentId || idGenerator.studentId(), // Jika siswa belum ada ID-nya
      type: 'IN',
      originSchool: command.originSchool,
      destinationSchool: context.schoolName, // Sekolah tujuan adalah sekolah saat ini
      reason: command.reason,
      effectiveDate: command.effectiveDate,
      status: 'PENDING',
      schoolId: context.schoolId,
      periodId: context.periodId,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    }

    const mutations = mockStorage.read(context.schoolId, 'mutations')
    mutations.push(mutationData)
    mockStorage.write(context.schoolId, 'mutations', mutations)

    return _success(mutationData)
  },

  async requestMutationOut(command, context) {
    await _simulateLatency(100)

    // Validasi: Siswa harus ada dan statusnya active
    const students = mockStorage.read(context.schoolId, 'students')
    const student = students.find((s) => s.studentId === command.studentId)

    if (!student) return _error('STUDENT_NOT_FOUND', 'Siswa tidak ditemukan.')
    if (student.status !== 'ACTIVE')
      return _error('INVALID_STATE', 'Hanya siswa aktif yang bisa dimutasi keluar.')

    const mutationData = {
      mutationId: idGenerator.mutationId(),
      studentId: command.studentId,
      type: 'OUT',
      originSchool: context.schoolName,
      destinationSchool: command.destinationSchool,
      reason: command.reason,
      effectiveDate: command.effectiveDate,
      status: 'PENDING',
      schoolId: context.schoolId,
      periodId: context.periodId,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    }

    const mutations = mockStorage.read(context.schoolId, 'mutations')
    mutations.push(mutationData)
    mockStorage.write(context.schoolId, 'mutations', mutations)

    return _success(mutationData)
  },

  async approveMutation(mutationId, command, context) {
    await _simulateLatency(200) // Simulasi proses side-effect yang lebih berat

    const mutations = mockStorage.read(context.schoolId, 'mutations')
    const index = mutations.findIndex((m) => m.mutationId === mutationId)

    if (index === -1) return _error('NOT_FOUND', 'Data mutasi tidak ditemukan.')
    if (mutations[index].status !== 'PENDING')
      return _error('INVALID_STATE', 'Hanya mutasi PENDING yang bisa di-approve.')

    // === SIMULASI SIDE-EFFECT ===
    const students = mockStorage.read(context.schoolId, 'students')

    if (mutations[index].type === 'IN') {
      // Side-effect Mutasi Masuk: Buat/Tambahkan siswa ke sekolah ini
      const newStudent = {
        studentId: mutations[index].studentId,
        fullName: { firstName: 'Siswa', lastName: 'Mutasi' }, // Simplified
        status: 'ACTIVE',
        enrollments: [
          {
            enrollmentId: idGenerator.enrollmentId(),
            classId: command.targetClassId,
            periodId: context.periodId,
            status: 'active',
            enrollmentDate: mutations[index].effectiveDate,
          },
        ],
      }
      students.push(newStudent)
    } else if (mutations[index].type === 'OUT') {
      // Side-effect Mutasi Keluar: Update status siswa menjadi TRANSFERRED
      const studentIndex = students.findIndex((s) => s.studentId === mutations[index].studentId)
      if (studentIndex !== -1) {
        // Tutup enrollment aktif
        students[studentIndex].enrollments?.forEach((e) => {
          if (e.status === 'active') e.status = 'completed'
        })
        students[studentIndex].status = 'TRANSFERRED'
        students[studentIndex].transferInfo = {
          targetSchool: mutations[index].destinationSchool,
          transferDate: mutations[index].effectiveDate,
        }
      }
    }

    mockStorage.write(context.schoolId, 'students', students)
    // ============================

    mutations[index].status = 'APPROVED'
    mutations[index].approvalNote = command.approvalNote || 'Disetujui'
    mutations[index].updatedAt = new Date().toISOString()
    mockStorage.write(context.schoolId, 'mutations', mutations)

    return _success(mutations[index])
  },

  async rejectMutation(mutationId, command, context) {
    await _simulateLatency(100)
    const mutations = mockStorage.read(context.schoolId, 'mutations')
    const index = mutations.findIndex((m) => m.mutationId === mutationId)

    if (index === -1) return _error('NOT_FOUND', 'Data mutasi tidak ditemukan.')
    if (mutations[index].status !== 'PENDING')
      return _error('INVALID_STATE', 'Hanya mutasi PENDING yang bisa di-reject.')

    mutations[index].status = 'REJECTED'
    mutations[index].rejectionReason = command.rejectionReason || 'Tidak memenuhi syarat mutasi.'
    mutations[index].updatedAt = new Date().toISOString()
    mockStorage.write(context.schoolId, 'mutations', mutations)

    return _success(mutations[index])
  },

  async getMutationById(mutationId, context) {
    await _simulateLatency(30)
    const mutations = mockStorage.read(context.schoolId, 'mutations')
    const data = mutations.find((m) => m.mutationId === mutationId)

    if (!data) return _error('NOT_FOUND', 'Data mutasi tidak ditemukan.')
    return _success(data)
  },

  async getMutations(context, filters = {}) {
    await _simulateLatency(100)
    let mutations = mockStorage.read(context.schoolId, 'mutations')

    if (filters.type) {
      mutations = mutations.filter((m) => m.type === filters.type)
    }
    if (filters.status) {
      mutations = mutations.filter((m) => m.status === filters.status)
    }
    if (filters.search) {
      const term = filters.search.toLowerCase()
      mutations = mutations.filter(
        (m) =>
          m.originSchool?.toLowerCase().includes(term) ||
          m.destinationSchool?.toLowerCase().includes(term),
      )
    }

    const page = filters.page || 1
    const limit = filters.limit || 20
    const total = mutations.length
    const items = mutations.slice((page - 1) * limit, page * limit)

    return _success({ items, total, page, limit, totalPages: Math.ceil(total / limit) })
  },
}

// apps/siakad-tu/src/adapters/mock/studentMockAdapter.js

import { mockStorage } from './mockStorage.js'
import { idGenerator } from '../utils/idGenerator.js'

/**
 * StudentMockAdapter — Implementasi mock dari StudentApiContract.
 *
 * Menggunakan localStorage sebagai "database".
 * Mensimulasikan:
 * - Persist data
 * - Validasi keunikan NISN/NIS
 * - Query dengan filter
 * - Error scenarios
 *
 * Catatan: Adapter TIDAK mengandung logika bisnis.
 * Logika bisnis ada di Engine. Adapter hanya menyimpan dan mengambil data.
 */
export const studentMockAdapter = {
  // === HELPER METHODS ===

  /**
   * Simulasi latency network (opsional, untuk realisme)
   * @private
   */
  async _simulateLatency(ms = 50) {
    return new Promise((resolve) => setTimeout(resolve, ms))
  },

  /**
   * Helper: buat error response
   * @private
   */
  _error(code, message, details = null) {
    return {
      success: false,
      data: null,
      error: { code, message, details },
    }
  },

  /**
   * Helper: buat success response
   * @private
   */
  _success(data) {
    return { success: true, data, error: null }
  },

  /**
   * Helper: cari siswa by ID di storage
   * @private
   */
  _findStudent(schoolId, studentId) {
    const students = mockStorage.read(schoolId, 'students')
    return students.find((s) => s.studentId === studentId) || null
  },

  /**
   * Helper: update siswa di storage
   * @private
   */
  _updateStudent(schoolId, updatedStudent) {
    const students = mockStorage.read(schoolId, 'students')
    const index = students.findIndex((s) => s.studentId === updatedStudent.studentId)
    if (index === -1) {
      students.push(updatedStudent)
    } else {
      students[index] = updatedStudent
    }
    mockStorage.write(schoolId, 'students', students)
  },

  // === CONTRACT IMPLEMENTATION ===

  /**
   * Register siswa baru
   */
  async registerStudent(command, context) {
    await this._simulateLatency()

    const { schoolId } = context
    const students = mockStorage.read(schoolId, 'students')

    // Validasi keunikan NISN
    const nisnExists = students.find((s) => s.nisn === command.nisn)
    if (nisnExists) {
      return this._error(
        'NISN_DUPLICATE',
        `NISN ${command.nisn} sudah terdaftar atas nama ${nisnExists.fullName.firstName} ${nisnExists.fullName.lastName || ''}`.trim(),
        { existingStudentId: nisnExists.studentId },
      )
    }

    // Validasi keunikan NIS
    const nisExists = students.find((s) => s.nis === command.nis)
    if (nisExists) {
      return this._error('NIS_DUPLICATE', `NIS ${command.nis} sudah digunakan.`, {
        existingStudentId: nisExists.studentId,
      })
    }

    // Simpan data (Engine sudah validasi format, Adapter hanya simpan)
    const studentData = {
      ...command,
      schoolId, // dari context
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    }

    this._updateStudent(schoolId, studentData)

    return this._success(studentData)
  },

  /**
   * Enroll siswa ke kelas
   */
  async enrollStudent(command, context) {
    await this._simulateLatency()

    const { schoolId, periodId } = context
    const student = this._findStudent(schoolId, command.studentId)

    if (!student) {
      return this._error('STUDENT_NOT_FOUND', `Siswa ${command.studentId} tidak ditemukan.`)
    }

    if (student.schoolId !== schoolId) {
      return this._error('SCHOOL_MISMATCH', 'Siswa terdaftar di sekolah lain.')
    }

    // Tambahkan enrollment baru
    const newEnrollment = {
      enrollmentId: command.enrollmentId || idGenerator.enrollmentId(),
      schoolId,
      periodId,
      classId: command.classId,
      enrollmentDate: command.enrollmentDate || new Date().toISOString(),
      status: 'active',
    }

    student.enrollments = student.enrollments || []
    student.enrollments.push(newEnrollment)
    student.updatedAt = new Date().toISOString()

    this._updateStudent(schoolId, student)

    return this._success(student)
  },

  /**
   * Update profil siswa
   */
  async updateStudentProfile(command, context) {
    await this._simulateLatency()

    const { schoolId } = context
    const student = this._findStudent(schoolId, command.studentId)

    if (!student) {
      return this._error('STUDENT_NOT_FOUND', `Siswa ${command.studentId} tidak ditemukan.`)
    }

    if (student.schoolId !== schoolId) {
      return this._error('SCHOOL_MISMATCH', 'Siswa terdaftar di sekolah lain.')
    }

    // Partial update
    if (command.address) student.address = command.address
    if (command.contactInfo) student.contactInfo = command.contactInfo
    if (command.guardianInfo) student.guardianInfo = command.guardianInfo
    student.updatedAt = new Date().toISOString()

    this._updateStudent(schoolId, student)

    return this._success(student)
  },

  /**
   * Luluskan siswa
   */
  async graduateStudent(command, context) {
    await this._simulateLatency()

    const { schoolId } = context
    const student = this._findStudent(schoolId, command.studentId)

    if (!student) {
      return this._error('STUDENT_NOT_FOUND', `Siswa ${command.studentId} tidak ditemukan.`)
    }

    if (student.schoolId !== schoolId) {
      return this._error('SCHOOL_MISMATCH', 'Siswa terdaftar di sekolah lain.')
    }

    // Complete semua enrollment aktif
    if (student.enrollments) {
      student.enrollments.forEach((e) => {
        if (e.status === 'active') e.status = 'completed'
      })
    }

    student.status = 'GRADUATED'
    student.graduationDate = command.graduationDate || new Date().toISOString()
    student.updatedAt = new Date().toISOString()

    this._updateStudent(schoolId, student)

    return this._success(student)
  },

  /**
   * Pindahkan siswa
   */
  async transferStudent(command, context) {
    await this._simulateLatency()

    const { schoolId } = context
    const student = this._findStudent(schoolId, command.studentId)

    if (!student) {
      return this._error('STUDENT_NOT_FOUND', `Siswa ${command.studentId} tidak ditemukan.`)
    }

    if (student.schoolId !== schoolId) {
      return this._error('SCHOOL_MISMATCH', 'Siswa terdaftar di sekolah lain.')
    }

    // Complete semua enrollment aktif
    if (student.enrollments) {
      student.enrollments.forEach((e) => {
        if (e.status === 'active') e.status = 'completed'
      })
    }

    student.status = 'TRANSFERRED'
    student.transferInfo = {
      targetSchool: command.targetSchool,
      reason: command.reason || '',
      transferDate: command.transferDate || new Date().toISOString(),
    }
    student.updatedAt = new Date().toISOString()

    this._updateStudent(schoolId, student)

    return this._success(student)
  },

  /**
   * Ambil siswa by ID
   */
  async getStudentById(studentId, context) {
    await this._simulateLatency(30)

    const { schoolId } = context
    const student = this._findStudent(schoolId, studentId)

    if (!student) {
      return this._error('STUDENT_NOT_FOUND', `Siswa ${studentId} tidak ditemukan.`)
    }

    if (student.schoolId !== schoolId) {
      return this._error('SCHOOL_MISMATCH', 'Siswa terdaftar di sekolah lain.')
    }

    return this._success(student)
  },

  /**
   * Ambil daftar siswa dengan filter
   */
  async getStudents(context, filters = {}) {
    await this._simulateLatency(100)

    const { schoolId, periodId } = context
    let students = mockStorage.read(schoolId, 'students')

    // Filter by status
    if (filters.status) {
      students = students.filter((s) => s.status === filters.status)
    }

    // Filter by class (dari enrollment aktif di periodId)
    if (filters.classId) {
      students = students.filter((s) =>
        s.enrollments?.some(
          (e) => e.periodId === periodId && e.classId === filters.classId && e.status === 'active',
        ),
      )
    }

    // Search by name/NISN/NIS
    if (filters.search) {
      const term = filters.search.toLowerCase()
      students = students.filter(
        (s) =>
          `${s.fullName.firstName} ${s.fullName.middleName || ''} ${s.fullName.lastName || ''}`
            .toLowerCase()
            .includes(term) ||
          s.nisn.toLowerCase().includes(term) ||
          s.nis.toLowerCase().includes(term),
      )
    }

    // Pagination
    const page = filters.page || 1
    const limit = filters.limit || 20
    const total = students.length
    const startIndex = (page - 1) * limit
    const items = students.slice(startIndex, startIndex + limit)

    return this._success({
      items,
      total,
      page,
      limit,
      totalPages: Math.ceil(total / limit),
    })
  },

  /**
   * Cek ketersediaan NISN
   */
  async checkNisnAvailability(nisn, context) {
    await this._simulateLatency(30)

    const { schoolId } = context
    const students = mockStorage.read(schoolId, 'students')
    const existing = students.find((s) => s.nisn === nisn)

    return this._success({
      available: !existing,
      existingStudentId: existing?.studentId || null,
    })
  },

  /**
   * Cek ketersediaan NIS
   */
  async checkNisAvailability(nis, context) {
    await this._simulateLatency(30)

    const { schoolId } = context
    const students = mockStorage.read(schoolId, 'students')
    const existing = students.find((s) => s.nis === nis)

    return this._success({
      available: !existing,
      existingStudentId: existing?.studentId || null,
    })
  },

  /**
   * Reset data mock (khusus development/testing)
   */
  async __resetMockData(schoolId = null) {
    mockStorage.clearSchool(schoolId)
    return this._success({ reset: true })
  },
}

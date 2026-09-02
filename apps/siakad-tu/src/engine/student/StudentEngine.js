// apps/siakad-tu/src/engine/student/StudentEngine.js

import { Student } from '../../domain/student/Student.js'
import { OperationalContext } from '../../domain/context/OperationalContext.js'
import { Result } from '../Result.js'
import { RegisterStudentCommand } from './commands/RegisterStudentCommand.js'
import { EnrollStudentCommand } from './commands/EnrollStudentCommand.js'
import { UpdateStudentProfileCommand } from './commands/UpdateStudentProfileCommand.js'
import { GraduateStudentCommand } from './commands/GraduateStudentCommand.js'
import { TransferStudentCommand } from './commands/TransferStudentCommand.js'

/**
 * StudentEngine — Pure JavaScript Business Engine.
 *
 * TIDAK BOLEH import: vue, pinia, axios, router, localStorage, database.
 *
 * Tanggung jawab:
 * 1. Menerima Command + OperationalContext
 * 2. Validasi command dan context
 * 3. Rehydrate Aggregate dari currentData (jika perlu)
 * 4. Execute domain behavior
 * 5. Return Result { success, data, events, errors }
 *
 * Engine TIDAK melakukan:
 * - Persist ke database
 * - Dispatch event ke listener
 * - Akses global state
 */
export class StudentEngine {
  /**
   * Validasi Operational Context
   * @private
   */
  _validateContext(context) {
    if (!(context instanceof OperationalContext)) {
      return Result.failure([
        {
          code: 'INVALID_CONTEXT',
          message: 'OperationalContext tidak valid.',
        },
      ])
    }

    if (!context.isValid()) {
      return Result.failure([
        {
          code: 'INCOMPLETE_CONTEXT',
          message:
            'OperationalContext tidak lengkap. Pastikan schoolId dan periodId sudah dipilih.',
        },
      ])
    }

    return null // no error
  }

  /**
   * Rehydrate Student Aggregate dari data mentah
   * @private
   */
  _rehydrateStudent(currentData) {
    if (!currentData) {
      return Result.failure([
        {
          code: 'STUDENT_NOT_FOUND',
          message: 'Data siswa tidak ditemukan.',
        },
      ])
    }

    try {
      const student = currentData instanceof Student ? currentData : Student.fromJSON(currentData)
      return { success: true, student }
    } catch (error) {
      return Result.failure([
        {
          code: 'REHYDRATION_FAILED',
          message: `Gagal memuat data siswa: ${error.message}`,
        },
      ])
    }
  }

  /**
   * USE CASE 1: Register Student
   *
   * Mendaftarkan siswa baru di sekolah saat ini.
   * Siswa otomatis berstatus ACTIVE.
   *
   * @param {RegisterStudentCommand|Object} commandData
   * @param {OperationalContext} context
   * @returns {Result} { success, data: { student }, events, errors }
   */
  registerStudent(commandData, context) {
    // 1. Validasi context
    const contextError = this._validateContext(context)
    if (contextError) return contextError

    // 2. Bangun command (akan throw jika invalid)
    let command
    try {
      command =
        commandData instanceof RegisterStudentCommand
          ? commandData
          : new RegisterStudentCommand(commandData)
    } catch (error) {
      return Result.failure(error.errors || [{ code: 'INVALID_COMMAND', message: error.message }])
    }

    // 3. Execute domain behavior
    try {
      const student = Student.register({
        studentId: command.studentId,
        schoolId: context.schoolId, // dari context, bukan dari command
        nisn: command.nisn,
        nis: command.nis,
        fullName: command.fullName,
        birthDate: command.birthDate,
        gender: command.gender,
        address: command.address,
        contactInfo: command.contactInfo,
        guardianInfo: command.guardianInfo,
      })

      return Result.success({ student: student.toJSON() }, student.getUncommittedEvents())
    } catch (error) {
      return Result.failure([
        {
          code: 'REGISTRATION_FAILED',
          message: error.message,
        },
      ])
    }
  }

  /**
   * USE CASE 2: Enroll Student
   *
   * Mendaftarkan siswa ke kelas pada periode akademik dari context.
   *
   * Invariants:
   * - Siswa harus ACTIVE
   * - Tidak boleh ada enrollment aktif di periode yang sama
   *
   * @param {EnrollStudentCommand|Object} commandData
   * @param {OperationalContext} context
   * @param {Object} currentData - Data siswa saat ini (dari projection/store)
   * @returns {Result}
   */
  enrollStudent(commandData, context, currentData) {
    // 1. Validasi context
    const contextError = this._validateContext(context)
    if (contextError) return contextError

    // 2. Bangun command
    let command
    try {
      command =
        commandData instanceof EnrollStudentCommand
          ? commandData
          : new EnrollStudentCommand(commandData)
    } catch (error) {
      return Result.failure(error.errors || [{ code: 'INVALID_COMMAND', message: error.message }])
    }

    // 3. Rehydrate aggregate
    const rehydrateResult = this._rehydrateStudent(currentData)
    if (!rehydrateResult.success) return rehydrateResult
    const student = rehydrateResult.student

    // 4. Validasi: siswa harus dari sekolah yang sama dengan context
    if (student.schoolId !== context.schoolId) {
      return Result.failure([
        {
          code: 'SCHOOL_MISMATCH',
          message: `Siswa terdaftar di sekolah lain (${student.schoolId}).`,
        },
      ])
    }

    // 5. Execute domain behavior
    try {
      student.enroll({
        enrollmentId: command.enrollmentId,
        periodId: context.periodId, // dari context, bukan dari command
        classId: command.classId,
        enrollmentDate: command.enrollmentDate,
      })

      return Result.success({ student: student.toJSON() }, student.getUncommittedEvents())
    } catch (error) {
      return Result.failure([
        {
          code: 'ENROLLMENT_FAILED',
          message: error.message,
        },
      ])
    }
  }

  /**
   * USE CASE 3: Update Student Profile
   *
   * Update address, contactInfo, atau guardianInfo.
   * Partial update — field yang tidak diisi diabaikan.
   *
   * @param {UpdateStudentProfileCommand|Object} commandData
   * @param {OperationalContext} context
   * @param {Object} currentData
   * @returns {Result}
   */
  updateStudentProfile(commandData, context, currentData) {
    const contextError = this._validateContext(context)
    if (contextError) return contextError

    let command
    try {
      command =
        commandData instanceof UpdateStudentProfileCommand
          ? commandData
          : new UpdateStudentProfileCommand(commandData)
    } catch (error) {
      return Result.failure(error.errors || [{ code: 'INVALID_COMMAND', message: error.message }])
    }

    const rehydrateResult = this._rehydrateStudent(currentData)
    if (!rehydrateResult.success) return rehydrateResult
    const student = rehydrateResult.student

    if (student.schoolId !== context.schoolId) {
      return Result.failure([
        {
          code: 'SCHOOL_MISMATCH',
          message: `Siswa terdaftar di sekolah lain.`,
        },
      ])
    }

    try {
      student.updateProfile({
        address: command.address,
        contactInfo: command.contactInfo,
        guardianInfo: command.guardianInfo,
      })

      return Result.success({ student: student.toJSON() }, student.getUncommittedEvents())
    } catch (error) {
      return Result.failure([
        {
          code: 'UPDATE_FAILED',
          message: error.message,
        },
      ])
    }
  }

  /**
   * USE CASE 4: Graduate Student
   *
   * Meluluskan siswa. Status berubah jadi GRADUATED.
   * Semua enrollment aktif di-complete.
   *
   * @param {GraduateStudentCommand|Object} commandData
   * @param {OperationalContext} context
   * @param {Object} currentData
   * @returns {Result}
   */
  graduateStudent(commandData, context, currentData) {
    const contextError = this._validateContext(context)
    if (contextError) return contextError

    let command
    try {
      command =
        commandData instanceof GraduateStudentCommand
          ? commandData
          : new GraduateStudentCommand(commandData)
      console.log('from command', command)
    } catch (error) {
      return Result.failure(error.errors || [{ code: 'INVALID_COMMAND', message: error.message }])
    }

    const rehydrateResult = this._rehydrateStudent(currentData)
    if (!rehydrateResult.success) return rehydrateResult
    const student = rehydrateResult.student

    if (student.schoolId !== context.schoolId) {
      return Result.failure([
        {
          code: 'SCHOOL_MISMATCH',
          message: `Siswa terdaftar di sekolah lain.`,
        },
      ])
    }

    try {
      student.graduate()

      return Result.success({ student: student.toJSON() }, student.getUncommittedEvents())
    } catch (error) {
      return Result.failure([
        {
          code: 'GRADUATION_FAILED',
          message: error.message,
        },
      ])
    }
  }

  /**
   * USE CASE 5: Transfer Student
   *
   * Memindahkan siswa ke sekolah lain.
   * Status berubah jadi TRANSFERRED.
   *
   * @param {TransferStudentCommand|Object} commandData
   * @param {OperationalContext} context
   * @param {Object} currentData
   * @returns {Result}
   */
  transferStudent(commandData, context, currentData) {
    const contextError = this._validateContext(context)
    if (contextError) return contextError

    let command
    try {
      command =
        commandData instanceof TransferStudentCommand
          ? commandData
          : new TransferStudentCommand(commandData)
    } catch (error) {
      return Result.failure(error.errors || [{ code: 'INVALID_COMMAND', message: error.message }])
    }

    const rehydrateResult = this._rehydrateStudent(currentData)
    if (!rehydrateResult.success) return rehydrateResult
    const student = rehydrateResult.student

    if (student.schoolId !== context.schoolId) {
      return Result.failure([
        {
          code: 'SCHOOL_MISMATCH',
          message: `Siswa terdaftar di sekolah lain.`,
        },
      ])
    }

    try {
      student.transfer({
        reason: command.reason,
        targetSchool: command.targetSchool,
      })

      return Result.success({ student: student.toJSON() }, student.getUncommittedEvents())
    } catch (error) {
      return Result.failure([
        {
          code: 'TRANSFER_FAILED',
          message: error.message,
        },
      ])
    }
  }
}

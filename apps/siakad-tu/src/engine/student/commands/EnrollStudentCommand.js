// apps/siakad-tu/src/engine/student/commands/EnrollStudentCommand.js

/**
 * Command: EnrollStudent
 *
 * Intent: Mendaftarkan siswa ke kelas pada periode akademik tertentu.
 */
export class EnrollStudentCommand {
  constructor(data) {
    this._validate(data)

    this.studentId = data.studentId
    this.enrollmentId = data.enrollmentId
    this.classId = data.classId
    this.enrollmentDate = data.enrollmentDate || new Date()

    Object.freeze(this)
  }

  _validate(data) {
    const errors = []

    if (!data.studentId)
      errors.push({ code: 'MISSING_FIELD', message: 'studentId wajib.', field: 'studentId' })
    if (!data.enrollmentId)
      errors.push({ code: 'MISSING_FIELD', message: 'enrollmentId wajib.', field: 'enrollmentId' })
    if (!data.classId)
      errors.push({ code: 'MISSING_FIELD', message: 'classId wajib.', field: 'classId' })

    if (errors.length > 0) {
      const error = new Error('EnrollStudentCommand validation failed')
      error.errors = errors
      throw error
    }
  }
}

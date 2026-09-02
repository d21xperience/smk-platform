// apps/siakad-tu/src/engine/student/commands/TransferStudentCommand.js

/**
 * Command: TransferStudent
 *
 * Intent: Memindahkan siswa ke sekolah lain.
 */
export class TransferStudentCommand {
  constructor(data) {
    this._validate(data)

    this.studentId = data.studentId
    this.reason = data.reason || ''
    this.targetSchool = data.targetSchool
    this.transferDate = data.transferDate || new Date()

    Object.freeze(this)
  }

  _validate(data) {
    const errors = []

    if (!data.studentId)
      errors.push({ code: 'MISSING_FIELD', message: 'studentId wajib.', field: 'studentId' })
    if (!data.targetSchool)
      errors.push({ code: 'MISSING_FIELD', message: 'targetSchool wajib.', field: 'targetSchool' })

    if (errors.length > 0) {
      const error = new Error('TransferStudentCommand validation failed')
      error.errors = errors
      throw error
    }
  }
}

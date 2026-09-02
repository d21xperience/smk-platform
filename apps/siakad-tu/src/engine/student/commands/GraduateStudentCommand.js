// apps/siakad-tu/src/engine/student/commands/GraduateStudentCommand.js

/**
 * Command: GraduateStudent
 *
 * Intent: Meluluskan siswa. Semua enrollment aktif akan di-complete.
 */
export class GraduateStudentCommand {
  constructor(data) {
    this._validate(data)

    this.studentId = data.studentId
    this.graduationDate = data.graduationDate || new Date()

    Object.freeze(this)
  }

  _validate(data) {
    if (!data.studentId) {
      const error = new Error('GraduateStudentCommand validation failed')
      error.errors = [{ code: 'MISSING_FIELD', message: 'studentId wajib.', field: 'studentId' }]
      throw error
    }
  }
}

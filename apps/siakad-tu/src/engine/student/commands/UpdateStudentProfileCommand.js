// apps/siakad-tu/src/engine/student/commands/UpdateStudentProfileCommand.js

/**
 * Command: UpdateStudentProfile
 *
 * Intent: Mengupdate data profil siswa (address, contact, guardian).
 * Field yang tidak diisi akan diabaikan (partial update).
 */
export class UpdateStudentProfileCommand {
  constructor(data) {
    this._validate(data)

    this.studentId = data.studentId
    this.address = data.address ? Object.freeze({ ...data.address }) : null
    this.contactInfo = data.contactInfo ? Object.freeze({ ...data.contactInfo }) : null
    this.guardianInfo = data.guardianInfo ? Object.freeze({ ...data.guardianInfo }) : null

    Object.freeze(this)
  }

  _validate(data) {
    const errors = []

    if (!data.studentId)
      errors.push({ code: 'MISSING_FIELD', message: 'studentId wajib.', field: 'studentId' })

    // Minimal satu field harus diupdate
    if (!data.address && !data.contactInfo && !data.guardianInfo) {
      errors.push({ code: 'NO_UPDATE', message: 'Minimal satu field harus diupdate.' })
    }

    if (errors.length > 0) {
      const error = new Error('UpdateStudentProfileCommand validation failed')
      error.errors = errors
      throw error
    }
  }
}

// apps/siakad-tu/src/engine/student/commands/RegisterStudentCommand.js

/**
 * Command: RegisterStudent
 *
 * Merepresentasikan intent untuk mendaftarkan siswa baru.
 * Immutable setelah dibuat. Self-validating.
 */
export class RegisterStudentCommand {
  constructor(data) {
    this._validate(data)

    this.studentId = data.studentId
    this.nisn = data.nisn
    this.nis = data.nis
    this.fullName = Object.freeze({ ...data.fullName })
    this.birthDate = data.birthDate
    this.gender = data.gender
    this.address = Object.freeze({ ...data.address })
    this.contactInfo = Object.freeze({ ...data.contactInfo })
    this.guardianInfo = Object.freeze({ ...data.guardianInfo })

    Object.freeze(this)
  }

  _validate(data) {
    const errors = []

    if (!data.studentId)
      errors.push({ code: 'MISSING_FIELD', message: 'studentId wajib diisi.', field: 'studentId' })
    if (!data.nisn)
      errors.push({ code: 'MISSING_FIELD', message: 'NISN wajib diisi.', field: 'nisn' })
    if (!data.nis) errors.push({ code: 'MISSING_FIELD', message: 'NIS wajib diisi.', field: 'nis' })
    if (!data.fullName?.firstName)
      errors.push({
        code: 'MISSING_FIELD',
        message: 'Nama depan wajib diisi.',
        field: 'fullName.firstName',
      })
    if (!data.birthDate)
      errors.push({
        code: 'MISSING_FIELD',
        message: 'Tanggal lahir wajib diisi.',
        field: 'birthDate',
      })
    if (!data.gender)
      errors.push({ code: 'MISSING_FIELD', message: 'Jenis kelamin wajib diisi.', field: 'gender' })
    if (!data.address?.street)
      errors.push({
        code: 'MISSING_FIELD',
        message: 'Alamat wajib diisi.',
        field: 'address.street',
      })
    if (!data.address?.village)
      errors.push({
        code: 'MISSING_FIELD',
        message: 'Kelurahan wajib diisi.',
        field: 'address.village',
      })
    if (!data.address?.city)
      errors.push({ code: 'MISSING_FIELD', message: 'Kota wajib diisi.', field: 'address.city' })
    if (!data.guardianInfo?.name)
      errors.push({
        code: 'MISSING_FIELD',
        message: 'Nama wali wajib diisi.',
        field: 'guardianInfo.name',
      })
    if (!data.guardianInfo?.relation)
      errors.push({
        code: 'MISSING_FIELD',
        message: 'Relasi wali wajib diisi.',
        field: 'guardianInfo.relation',
      })

    // Validasi minimal kontak
    if (!data.contactInfo?.phone && !data.contactInfo?.email) {
      errors.push({
        code: 'MISSING_CONTACT',
        message: 'Minimal phone atau email harus diisi.',
        field: 'contactInfo',
      })
    }

    if (errors.length > 0) {
      const error = new Error('RegisterStudentCommand validation failed')
      error.errors = errors
      throw error
    }
  }
}

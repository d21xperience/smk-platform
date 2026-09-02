// apps/siakad-tu/src/composables/student/helpers/studentValidators.js

/**
 * Student Validators — Pure functions untuk UI-level validation.
 *
 * PENTING: Ini BUKAN business validation.
 * - UI validation: format NISN (10 digit), required fields, format email
 * - Business validation: NISN unik, siswa tidak bisa enroll 2x → di Engine
 *
 * UI validation memberikan feedback INSTAN di form,
 * sebelum data dikirim ke Service/Engine.
 */

/**
 * Validasi format NISN (10 digit angka)
 * @param {string} nisn
 * @returns {{ valid: boolean, message: string }}
 */
export function validateNisnFormat(nisn) {
  if (!nisn || String(nisn).trim() === '') {
    return { valid: false, message: 'NISN wajib diisi.' }
  }
  if (!/^\d{10}$/.test(String(nisn).trim())) {
    return { valid: false, message: 'NISN harus 10 digit angka.' }
  }
  return { valid: true, message: '' }
}

/**
 * Validasi format NIS (tidak kosong)
 * @param {string} nis
 * @returns {{ valid: boolean, message: string }}
 */
export function validateNisFormat(nis) {
  if (!nis || String(nis).trim() === '') {
    return { valid: false, message: 'NIS wajib diisi.' }
  }
  return { valid: true, message: '' }
}

/**
 * Validasi nama depan (required)
 * @param {string} firstName
 * @returns {{ valid: boolean, message: string }}
 */
export function validateFirstName(firstName) {
  if (!firstName || String(firstName).trim() === '') {
    return { valid: false, message: 'Nama depan wajib diisi.' }
  }
  return { valid: true, message: '' }
}

/**
 * Validasi tanggal lahir
 * @param {string} birthDate - YYYY-MM-DD format
 * @returns {{ valid: boolean, message: string }}
 */
export function validateBirthDate(birthDate) {
  if (!birthDate) {
    return { valid: false, message: 'Tanggal lahir wajib diisi.' }
  }
  const date = new Date(birthDate)
  if (isNaN(date.getTime())) {
    return { valid: false, message: 'Format tanggal tidak valid.' }
  }
  if (date > new Date()) {
    return { valid: false, message: 'Tanggal lahir tidak boleh di masa depan.' }
  }
  return { valid: true, message: '' }
}

/**
 * Validasi jenis kelamin
 * @param {string} gender
 * @returns {{ valid: boolean, message: string }}
 */
export function validateGender(gender) {
  if (!gender || !['MALE', 'FEMALE'].includes(gender)) {
    return { valid: false, message: 'Jenis kelamin wajib dipilih.' }
  }
  return { valid: true, message: '' }
}

/**
 * Validasi alamat (minimal street, village, city)
 * @param {Object} address
 * @returns {{ valid: boolean, errors: Object }}
 */
export function validateAddress(address) {
  const errors = {}
  if (!address?.street?.trim()) errors.street = 'Jalan wajib diisi.'
  if (!address?.village?.trim()) errors.village = 'Kelurahan wajib diisi.'
  if (!address?.city?.trim()) errors.city = 'Kota wajib diisi.'

  return {
    valid: Object.keys(errors).length === 0,
    errors,
  }
}

/**
 * Validasi contact info (minimal phone atau email)
 * @param {Object} contactInfo
 * @returns {{ valid: boolean, errors: Object }}
 */
export function validateContactInfo(contactInfo) {
  const errors = {}

  if (!contactInfo?.phone?.trim() && !contactInfo?.email?.trim()) {
    errors._form = 'Minimal phone atau email harus diisi.'
    return { valid: false, errors }
  }

  if (contactInfo?.phone?.trim()) {
    if (!/^[\d\s\-+()]{8,20}$/.test(contactInfo.phone)) {
      errors.phone = 'Format nomor telepon tidak valid.'
    }
  }

  if (contactInfo?.email?.trim()) {
    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(contactInfo.email)) {
      errors.email = 'Format email tidak valid.'
    }
  }

  return {
    valid: Object.keys(errors).length === 0,
    errors,
  }
}

/**
 * Validasi guardian info
 * @param {Object} guardianInfo
 * @returns {{ valid: boolean, errors: Object }}
 */
export function validateGuardianInfo(guardianInfo) {
  const errors = {}
  if (!guardianInfo?.name?.trim()) errors.name = 'Nama wali wajib diisi.'
  if (!guardianInfo?.relation) {
    errors.relation = 'Relasi wali wajib dipilih.'
  } else if (!['father', 'mother', 'guardian'].includes(guardianInfo.relation)) {
    errors.relation = 'Relasi tidak valid.'
  }

  return {
    valid: Object.keys(errors).length === 0,
    errors,
  }
}

/**
 * Validasi form lengkap (semua field)
 * @param {Object} formData
 * @returns {{ valid: boolean, errors: Object }}
 */
export function validateStudentForm(formData) {
  const errors = {}

  // NISN
  const nisnResult = validateNisnFormat(formData.nisn)
  if (!nisnResult.valid) errors.nisn = nisnResult.message

  // NIS
  const nisResult = validateNisFormat(formData.nis)
  if (!nisResult.valid) errors.nis = nisResult.message

  // Nama depan
  const nameResult = validateFirstName(formData.fullName?.firstName)
  if (!nameResult.valid) errors['fullName.firstName'] = nameResult.message

  // Tanggal lahir
  const birthResult = validateBirthDate(formData.birthDate)
  if (!birthResult.valid) errors.birthDate = birthResult.message

  // Gender
  const genderResult = validateGender(formData.gender)
  if (!genderResult.valid) errors.gender = genderResult.message

  // Address
  const addressResult = validateAddress(formData.address)
  if (!addressResult.valid) {
    Object.keys(addressResult.errors).forEach((key) => {
      errors[`address.${key}`] = addressResult.errors[key]
    })
  }

  // Contact
  const contactResult = validateContactInfo(formData.contactInfo)
  if (!contactResult.valid) {
    Object.keys(contactResult.errors).forEach((key) => {
      errors[`contactInfo.${key}`] = contactResult.errors[key]
    })
  }

  // Guardian
  const guardianResult = validateGuardianInfo(formData.guardianInfo)
  if (!guardianResult.valid) {
    Object.keys(guardianResult.errors).forEach((key) => {
      errors[`guardianInfo.${key}`] = guardianResult.errors[key]
    })
  }

  return {
    valid: Object.keys(errors).length === 0,
    errors,
  }
}

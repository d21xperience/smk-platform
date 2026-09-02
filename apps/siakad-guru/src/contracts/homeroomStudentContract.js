// FILE: src/contracts/homeroomStudentContract.js
// STATUS: NEW
// STATUS IMPLEMENTASI: COMPLETE

/**
 * Homeroom Student Contract
 * Menentukan bentuk data yang dibutuhkan oleh UI Daftar Siswa Wali Kelas.
 * Frontend-First: Hanya field yang diperlukan UI yang dipetakan.
 */

export const STUDENT_STATUS = {
  ACTIVE: 'ACTIVE',
  INACTIVE: 'INACTIVE',
  TRANSFERRED: 'TRANSFERRED',
  GRADUATED: 'GRADUATED'
}

export const GENDER = {
  MALE: 'L',
  FEMALE: 'P'
}

/**
 * Memetakan data mentah dari adapter ke kontrak frontend.
 * @param {Object} raw - Data siswa mentah dari adapter/backend.
 * @returns {Object} Data siswa yang telah dipetakan sesuai kontrak.
 */
export function mapHomeroomStudent(raw) {
  if (!raw) return null

  return {
    id: String(raw.id || ''),
    nisn: String(raw.nisn || ''),
    nis: String(raw.nis || ''),
    fullName: String(raw.fullName || raw.name || ''),
    gender: raw.gender === GENDER.FEMALE ? GENDER.FEMALE : GENDER.MALE,
    photoUrl: raw.photoUrl || null,
    status: Object.values(STUDENT_STATUS).includes(raw.status) ? raw.status : STUDENT_STATUS.ACTIVE,
    parentContact: String(raw.parentContact || raw.parentPhone || ''),
    seatNumber: Number(raw.seatNumber || raw.absen || 0)
  }
}

/**
 * Memetakan daftar data mentah menjadi array sesuai kontrak.
 * @param {Array} rawList - Array data siswa mentah.
 * @returns {Array} Array data siswa yang telah dipetakan.
 */
export function mapHomeroomStudents(rawList) {
  if (!Array.isArray(rawList)) return []
  return rawList.map(mapHomeroomStudent).filter(s => s && s.id)
}

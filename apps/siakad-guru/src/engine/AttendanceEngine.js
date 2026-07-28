/**
 * Attendance Engine
 * Pure business logic untuk absensi siswa.
 * Tidak bergantung pada framework apa pun.
 */

export const AttendanceStatus = {
  HADIR: 'hadir',
  IZIN: 'izin',
  SAKIT: 'sakit',
  ALPHA: 'alpha',
  UNMARKED: 'unmarked', // belum ditandai
}

export class AttendanceEngine {
  /**
   * Menghitung ringkasan absensi dari array siswa.
   * @param {Array<{ status: string }>} students
   * @returns {{ total: number, hadir: number, izin: number, sakit: number, alpha: number, unmarked: number }}
   */
  static calculateSummary(students) {
    const total = students.length
    let hadir = 0,
      izin = 0,
      sakit = 0,
      alpha = 0,
      unmarked = 0
    console.log('AttendanceService',students)
    for (const s of students) {
      switch (s.status) {
        case AttendanceStatus.HADIR:
          hadir++
          break
        case AttendanceStatus.IZIN:
          izin++
          break
        case AttendanceStatus.SAKIT:
          sakit++
          break
        case AttendanceStatus.ALPHA:
          alpha++
          break
        default:
          unmarked++
      }
    }

    return { total, hadir, izin, sakit, alpha, unmarked }
  }

  /**
   * Memeriksa apakah semua siswa sudah ditandai (tidak ada 'unmarked').
   * @param {Array<{ status: string }>} students
   * @returns {boolean}
   */
  static isAllMarked(students) {
    return students.every((s) => s.status !== AttendanceStatus.UNMARKED)
  }

  /**
   * Mendapatkan daftar siswa yang belum ditandai.
   * @param {Array<{ studentId: number|string, status: string }>} students
   * @returns {Array<{ studentId: number|string }>}
   */
  static getUnmarkedStudents(students) {
    return students
      .filter((s) => s.status === AttendanceStatus.UNMARKED)
      .map((s) => ({ studentId: s.studentId }))
  }

  /**
   * Memvalidasi apakah record absensi bisa disubmit.
   * Syarat:
   * - Semua siswa harus ditandai (tidak boleh unmarked).
   * - Status record harus 'draft' (tidak bisa submit dua kali).
   *
   * @param {Object} record - { status: string, students: Array<{ status: string }> }
   * @returns {{ valid: boolean, errors: string[] }}
   */
  static validateForSubmit(record) {
    const errors = []

    if (record.status !== 'draft') {
      errors.push('Absensi sudah disubmit sebelumnya.')
    }

    if (!AttendanceEngine.isAllMarked(record.students)) {
      const unmarked = AttendanceEngine.getUnmarkedStudents(record.students)
      errors.push(`Masih ada ${unmarked.length} siswa yang belum ditandai.`)
    }

    return {
      valid: errors.length === 0,
      errors,
    }
  }

  /**
   * Menentukan apakah suatu record bisa diedit.
   * Hanya record dengan status 'draft' yang bisa diedit.
   * @param {string} status
   * @returns {boolean}
   */
  static canEdit(status) {
    return status === 'draft'
  }

  /**
   * Membuat objek siswa baru dengan status default 'unmarked'.
   * @param {{ id: number|string, name: string, nis: string }} studentData
   * @returns {{ studentId: number|string, studentName: string, nis: string, status: string, note: string }}
   */
  static createStudentItem(studentData) {
    console.assert('AttendanceEngine!', studentData)
    return {
      studentId: studentData.studentId,
      studentName: studentData.studentName,
      nis: studentData.nis,
      status: AttendanceStatus.HADIR,
      note: '',
    }
  }

  /**
   * Mengubah status absensi seorang siswa dengan aturan:
   * - Jika di-click status yang sama → kembali ke 'unmarked'
   * - Jika berbeda → ganti ke status baru
   *
   * @param {Object} student - { status: string }
   * @param {string} newStatus
   * @returns {string} status akhir
   */
  static toggleStudentStatus(student, newStatus) {
    console.log('AttendanceEngine', newStatus)
    console.log('AttendanceEngine', student.status === newStatus)
    if (student.status === newStatus) {
      return AttendanceStatus.UNMARKED
    }
    return newStatus
  }

  /**
   * Menandai semua siswa yang belum ditandai dengan status tertentu.
   * @param {Array<{ status: string }>} students
   * @param {string} status
   * @returns {Array<{ status: string }>} array status baru
   */
  static markAllUnmarkedAs(students, status) {
    return students.map((s) => {
      if (s.status === AttendanceStatus.UNMARKED) {
        return { ...s, status }
      }
      return s
    })
  }
}

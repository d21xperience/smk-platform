// FILE: src/domain/homeroom-attendance/HomeroomAttendanceEngine.js
// STATUS: NEW
// STATUS IMPLEMENTASI: COMPLETE

/**
 * Homeroom Attendance Engine
 * Pure JavaScript Domain Engine untuk business rules rekap kehadiran wali kelas.
 * TIDAK BOLEH import: Vue, Pinia, Axios, Quasar, HTTP, Database, Browser API.
 *
 * Business Rules:
 * - Persentase kehadiran: (present / totalDays) × 100
 * - Status: GOOD >= 90%, WARNING >= 75%, CRITICAL < 75%
 * - Agregasi bulanan per siswa
 * - Agregasi semester (total seluruh bulan)
 * - Validasi: count tidak boleh negatif, totalDays harus > 0
 */

const GOOD_THRESHOLD = 90
const WARNING_THRESHOLD = 75

export const ATTENDANCE_STATUS = {
  GOOD: 'GOOD',
  WARNING: 'WARNING',
  CRITICAL: 'CRITICAL',
}

export class HomeroomAttendanceEngine {
  /**
   * Validasi jumlah (tidak boleh negatif).
   * Jika null/undefined/NaN/negatif → return 0.
   * @param {number} value - Nilai mentah.
   * @returns {number} Nilai valid (>= 0).
   */
  validateCount(value) {
    if (value === null || value === undefined) return 0
    const num = Number(value)
    if (Number.isNaN(num)) return 0
    return Math.max(0, Math.floor(num))
  }

  /**
   * Hitung persentase kehadiran.
   * Formula: (present / totalDays) × 100
   * @param {number} present - Jumlah hari hadir.
   * @param {number} totalDays - Total hari efektif.
   * @returns {number} Persentase (0-100, dibulatkan 1 desimal).
   */
  calculatePercentage(present, totalDays) {
    const validPresent = this.validateCount(present)
    const validTotalDays = this.validateCount(totalDays)

    if (validTotalDays === 0) return 0

    const percentage = (validPresent / validTotalDays) * 100
    return Math.round(percentage * 10) / 10
  }

  /**
   * Tentukan status kehadiran berdasarkan persentase.
   * GOOD: >= 90%
   * WARNING: >= 75% dan < 90%
   * CRITICAL: < 75%
   * @param {number} percentage - Persentase kehadiran.
   * @returns {string} Status ('GOOD', 'WARNING', 'CRITICAL').
   */
  determineStatus(percentage) {
    if (percentage >= GOOD_THRESHOLD) return ATTENDANCE_STATUS.GOOD
    if (percentage >= WARNING_THRESHOLD) return ATTENDANCE_STATUS.WARNING
    return ATTENDANCE_STATUS.CRITICAL
  }

  /**
   * Proses data kehadiran satu bulan untuk satu siswa.
   * @param {Object} monthInfo - Informasi bulan (monthId, monthName, totalDays).
   * @param {Object} attendanceData - Data kehadiran mentah untuk bulan tersebut.
   * @returns {Object} Rekap kehadiran bulanan yang telah diproses.
   */
  processMonthlyRecord(monthInfo, attendanceData) {
    const present = this.validateCount(attendanceData.present)
    const sick = this.validateCount(attendanceData.sick)
    const permitted = this.validateCount(attendanceData.permitted)
    const absent = this.validateCount(attendanceData.absent)
    const totalDays = this.validateCount(monthInfo.totalDays)
    const percentage = this.calculatePercentage(present, totalDays)

    return {
      monthId: monthInfo.monthId,
      monthName: monthInfo.monthName,
      present: present,
      sick: sick,
      permitted: permitted,
      absent: absent,
      totalDays: totalDays,
      percentage: percentage,
    }
  }

  /**
   * Agregasi seluruh data bulanan menjadi total semester.
   * @param {Array<Object>} monthlyRecords - Array rekap bulanan yang telah diproses.
   * @returns {Object} Total semester.
   */
  aggregateSemester(monthlyRecords) {
    if (!Array.isArray(monthlyRecords) || monthlyRecords.length === 0) {
      return {
        present: 0,
        sick: 0,
        permitted: 0,
        absent: 0,
        totalDays: 0,
        percentage: 0,
      }
    }

    let totalPresent = 0
    let totalSick = 0
    let totalPermitted = 0
    let totalAbsent = 0
    let totalDays = 0

    for (const record of monthlyRecords) {
      totalPresent += record.present
      totalSick += record.sick
      totalPermitted += record.permitted
      totalAbsent += record.absent
      totalDays += record.totalDays
    }

    const percentage = this.calculatePercentage(totalPresent, totalDays)

    return {
      present: totalPresent,
      sick: totalSick,
      permitted: totalPermitted,
      absent: totalAbsent,
      totalDays: totalDays,
      percentage: percentage,
    }
  }

  /**
   * Proses data mentah satu siswa menjadi rekap lengkap.
   * @param {Object} rawStudent - Data mentah siswa dari adapter.
   * @param {Array<Object>} months - Array informasi bulan.
   * @returns {Object} Rekap kehadiran siswa yang telah diproses.
   */
  processStudent(rawStudent, months) {
    const monthlyRecords = []

    for (const monthInfo of months) {
      const attendanceData = rawStudent.attendance[monthInfo.monthId] || {
        present: 0,
        sick: 0,
        permitted: 0,
        absent: 0,
      }
      const record = this.processMonthlyRecord(monthInfo, attendanceData)
      monthlyRecords.push(record)
    }

    const semesterTotal = this.aggregateSemester(monthlyRecords)
    const status = this.determineStatus(semesterTotal.percentage)

    return {
      studentId: rawStudent.studentId,
      studentName: rawStudent.studentName,
      seatNumber: rawStudent.seatNumber || 0,
      monthlyRecords: monthlyRecords,
      semesterTotal: semesterTotal,
      status: status,
    }
  }

  /**
   * Hitung summary kelas dari seluruh siswa yang telah diproses.
   * @param {Array<Object>} processedStudents - Array siswa yang telah diproses.
   * @returns {Object} Summary kelas.
   */
  calculateClassSummary(processedStudents) {
    if (!Array.isArray(processedStudents) || processedStudents.length === 0) {
      return {
        totalStudents: 0,
        classAveragePercentage: 0,
        goodCount: 0,
        warningCount: 0,
        criticalCount: 0,
      }
    }

    let totalPercentage = 0
    let goodCount = 0
    let warningCount = 0
    let criticalCount = 0

    for (const student of processedStudents) {
      totalPercentage += student.semesterTotal.percentage

      if (student.status === ATTENDANCE_STATUS.GOOD) {
        goodCount++
      } else if (student.status === ATTENDANCE_STATUS.WARNING) {
        warningCount++
      } else {
        criticalCount++
      }
    }

    const classAveragePercentage =
      Math.round((totalPercentage / processedStudents.length) * 10) / 10

    return {
      totalStudents: processedStudents.length,
      classAveragePercentage: classAveragePercentage,
      goodCount: goodCount,
      warningCount: warningCount,
      criticalCount: criticalCount,
    }
  }

  /**
   * Proses seluruh data mentah kelas menjadi rekap lengkap dengan summary.
   * @param {Object} rawData - Data mentah dari adapter.
   * @param {string} rawData.classId - ID kelas.
   * @param {string} rawData.className - Nama kelas.
   * @param {Array<Object>} rawData.months - Array informasi bulan.
   * @param {Array<Object>} rawData.students - Array data siswa mentah.
   * @returns {Object} Rekap lengkap kehadiran kelas.
   */
  processClassAttendance(rawData) {
    const { classId, className, months, students } = rawData

    if (!Array.isArray(months) || !Array.isArray(students)) {
      return {
        classId: classId || null,
        className: className || null,
        months: [],
        summary: {
          totalStudents: 0,
          classAveragePercentage: 0,
          goodCount: 0,
          warningCount: 0,
          criticalCount: 0,
        },
        students: [],
      }
    }

    const processedStudents = students.map((student) => this.processStudent(student, months))

    const summary = this.calculateClassSummary(processedStudents)

    return {
      classId: classId || null,
      className: className || null,
      months: months,
      summary: summary,
      students: processedStudents,
    }
  }
}

// FILE: src/contracts/homeroomAttendanceContract.js
// STATUS: NEW
// STATUS IMPLEMENTASI: COMPLETE

/**
 * Homeroom Attendance Contract
 * Menentukan bentuk data yang dibutuhkan oleh UI Rekap Kehadiran Wali Kelas.
 * Frontend-First: hanya field yang diperlukan UI yang dipetakan.
 */

export const ATTENDANCE_STATUS = {
  GOOD: 'GOOD',
  WARNING: 'WARNING',
  CRITICAL: 'CRITICAL',
}

/**
 * Memetakan hasil olahan Domain Engine ke kontrak UI.
 * Engine sudah memproses, contract hanya memastikan shape konsisten.
 * @param {Object} engineResult - Hasil dari HomeroomAttendanceEngine.processClassAttendance().
 * @returns {Object} Data sesuai kontrak UI.
 */
export function mapHomeroomAttendance(engineResult) {
  if (!engineResult) {
    return {
      classId: null,
      className: null,
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

  return {
    classId: engineResult.classId || null,
    className: engineResult.className || null,
    months: Array.isArray(engineResult.months) ? engineResult.months.map(mapMonthInfo) : [],
    summary: {
      totalStudents: engineResult.summary?.totalStudents || 0,
      classAveragePercentage: engineResult.summary?.classAveragePercentage || 0,
      goodCount: engineResult.summary?.goodCount || 0,
      warningCount: engineResult.summary?.warningCount || 0,
      criticalCount: engineResult.summary?.criticalCount || 0,
    },
    students: Array.isArray(engineResult.students)
      ? engineResult.students.map(mapStudentAttendance)
      : [],
  }
}

/**
 * Memetakan informasi bulan.
 * @param {Object} month - Data bulan dari engine.
 * @returns {Object} Data bulan sesuai kontrak UI.
 */
function mapMonthInfo(month) {
  return {
    monthId: month.monthId || '',
    monthName: month.monthName || '',
    totalDays: month.totalDays || 0,
  }
}

/**
 * Memetakan data satu siswa dari hasil engine ke kontrak UI.
 * @param {Object} student - Data siswa dari engine.
 * @returns {Object} Data siswa sesuai kontrak UI.
 */
function mapStudentAttendance(student) {
  return {
    studentId: student.studentId || '',
    studentName: student.studentName || '',
    seatNumber: student.seatNumber || 0,
    monthlyRecords: Array.isArray(student.monthlyRecords)
      ? student.monthlyRecords.map(mapMonthlyRecord)
      : [],
    semesterTotal: {
      present: student.semesterTotal?.present || 0,
      sick: student.semesterTotal?.sick || 0,
      permitted: student.semesterTotal?.permitted || 0,
      absent: student.semesterTotal?.absent || 0,
      totalDays: student.semesterTotal?.totalDays || 0,
      percentage: student.semesterTotal?.percentage || 0,
    },
    status: student.status || ATTENDANCE_STATUS.CRITICAL,
  }
}

/**
 * Memetakan rekap bulanan dari hasil engine ke kontrak UI.
 * @param {Object} record - Data rekap bulanan dari engine.
 * @returns {Object} Data rekap bulanan sesuai kontrak UI.
 */
function mapMonthlyRecord(record) {
  return {
    monthId: record.monthId || '',
    monthName: record.monthName || '',
    present: record.present || 0,
    sick: record.sick || 0,
    permitted: record.permitted || 0,
    absent: record.absent || 0,
    totalDays: record.totalDays || 0,
    percentage: record.percentage || 0,
  }
}

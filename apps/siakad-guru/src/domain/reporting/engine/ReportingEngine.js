export class ReportingEngine {
  /**
   * Agregasi statistik kelas dari data siswa
   * @param {Array} students - Array of StudentProgress
   * @param {Array} gradeScale - [{ min, max, predicate }]
   * @returns {Object} { total, average, highest, lowest, distribution, attendanceRate, promotionRate }
   */
  static aggregateClassStatistics(
    students,
    gradeScale = [
      { min: 85, max: 100, predicate: 'A' },
      { min: 70, max: 84.99, predicate: 'B' },
      { min: 60, max: 69.99, predicate: 'C' },
      { min: 0, max: 59.99, predicate: 'D' },
    ],
  ) {
    if (!students || students.length === 0) {
      return {
        total: 0,
        average: 0,
        highest: 0,
        lowest: 0,
        distribution: {},
        attendanceRate: 0,
        promotionRate: 0,
      }
    }

    const total = students.length
    let sumScores = 0
    let highest = -Infinity
    let lowest = Infinity
    const distribution = {}

    // Inisialisasi distribution berdasarkan gradeScale
    for (const scale of gradeScale) {
      distribution[scale.predicate] = 0
    }

    let totalAttendance = 0
    let totalPossibleAttendance = 0
    let naikCount = 0

    for (const student of students) {
      // Academic scores
      const score = student.averageScore || 0
      sumScores += score
      if (score > highest) highest = score
      if (score < lowest) lowest = score

      // Distribution
      let assigned = false
      for (const scale of gradeScale) {
        if (score >= scale.min && score <= scale.max) {
          distribution[scale.predicate] = (distribution[scale.predicate] || 0) + 1
          assigned = true
          break
        }
      }
      if (!assigned) {
        distribution['E'] = (distribution['E'] || 0) + 1
      }

      // Attendance
      // const totalRecords = (student.totalPresent || 0) + (student.totalSick || 0) +
      //                      (student.totalPermit || 0) + (student.totalAbsent || 0) +
      //                      (student.totalLate || 0);
      totalPossibleAttendance += student.totalSessions || 0 // Asumsikan ada field totalSessions
      totalAttendance += (student.totalPresent || 0) + (student.totalLate || 0) // Hadir + Terlambat dianggap hadir

      // Promotion
      if (student.promotionStatus === 'NAIK') naikCount++
    }

    const average = total > 0 ? sumScores / total : 0
    const attendanceRate =
      totalPossibleAttendance > 0 ? (totalAttendance / totalPossibleAttendance) * 100 : 0
    const promotionRate = total > 0 ? (naikCount / total) * 100 : 0

    return {
      total,
      average: Math.round(average * 100) / 100,
      highest: highest === -Infinity ? 0 : Math.round(highest * 100) / 100,
      lowest: lowest === Infinity ? 0 : Math.round(lowest * 100) / 100,
      distribution,
      attendanceRate: Math.round(attendanceRate * 100) / 100,
      promotionRate: Math.round(promotionRate * 100) / 100,
      naikCount,
      tidakNaikCount: total - naikCount,
    }
  }

  /**
   * Generate data untuk laporan kelas (rows)
   */
  static generateClassReportRows(students, columns) {
    return students.map((student) => {
      const row = {}
      for (const col of columns) {
        if (col.key === 'no') {
          // akan diisi nanti
        } else if (col.key === 'name') {
          row[col.key] = student.studentName || ''
        } else if (col.key === 'nis') {
          row[col.key] = student.nis || ''
        } else if (col.key === 'averageScore') {
          row[col.key] = student.averageScore || 0
        } else if (col.key === 'predicate') {
          row[col.key] = student.predicate || ''
        } else if (col.key === 'status') {
          row[col.key] = student.promotionStatus || 'UNDEFINED'
        } else if (col.key === 'attendance') {
          const total =
            (student.totalPresent || 0) +
            (student.totalSick || 0) +
            (student.totalPermit || 0) +
            (student.totalAbsent || 0) +
            (student.totalLate || 0)
          row[col.key] = total
        } else {
          row[col.key] = student[col.key] || ''
        }
      }
      return row
    })
  }

  /**
   * Agregasi dashboard guru
   */
  static aggregateTeacherDashboard(teachingSessions, attendanceSummaries, assessmentSummaries) {
    const totalSessions = teachingSessions?.length || 0
    const totalClasses = new Set(teachingSessions?.map((s) => s.classId) || []).size
    const totalSubjects = new Set(teachingSessions?.map((s) => s.subjectId) || []).size

    // Rata-rata kehadiran
    let totalAttendanceRate = 0
    let count = 0
    if (attendanceSummaries) {
      for (const summary of attendanceSummaries) {
        totalAttendanceRate += summary.attendanceRate || 0
        count++
      }
    }
    const avgAttendanceRate = count > 0 ? totalAttendanceRate / count : 0

    // Rata-rata nilai
    let totalAvgScore = 0
    let scoreCount = 0
    if (assessmentSummaries) {
      for (const summary of assessmentSummaries) {
        totalAvgScore += summary.averageScore || 0
        scoreCount++
      }
    }
    const avgScore = scoreCount > 0 ? totalAvgScore / scoreCount : 0

    return {
      totalSessions,
      totalClasses,
      totalSubjects,
      avgAttendanceRate: Math.round(avgAttendanceRate * 100) / 100,
      avgScore: Math.round(avgScore * 100) / 100,
    }
  }

  /**
   * Siapkan data untuk ekspor (CSV/Excel/PDF)
   * @param {Array} rows - Array of row objects
   * @param {Array} columns - [{ key, label }]
   * @param {Object} summary - Summary data
   * @returns {Object} { headers: [], rows: [], summary: {} }
   */
  static prepareExportData(rows, columns, summary = {}) {
    const headers = columns.map((col) => ({ key: col.key, label: col.label || col.key }))
    const dataRows = rows.map((row) => {
      const newRow = {}
      for (const col of columns) {
        newRow[col.label || col.key] = row[col.key] !== undefined ? row[col.key] : ''
      }
      return newRow
    })

    return {
      headers,
      rows: dataRows,
      summary,
      totalRows: rows.length,
    }
  }

  /**
   * Hitung distribusi nilai berdasarkan gradeScale
   */
  static calculateDistribution(
    scores,
    gradeScale = [
      { min: 85, max: 100, predicate: 'A' },
      { min: 70, max: 84.99, predicate: 'B' },
      { min: 60, max: 69.99, predicate: 'C' },
      { min: 0, max: 59.99, predicate: 'D' },
    ],
  ) {
    const distribution = {}
    for (const scale of gradeScale) {
      distribution[scale.predicate] = 0
    }

    for (const score of scores) {
      let assigned = false
      for (const scale of gradeScale) {
        if (score >= scale.min && score <= scale.max) {
          distribution[scale.predicate] = (distribution[scale.predicate] || 0) + 1
          assigned = true
          break
        }
      }
      if (!assigned) {
        distribution['E'] = (distribution['E'] || 0) + 1
      }
    }
    return distribution
  }
}

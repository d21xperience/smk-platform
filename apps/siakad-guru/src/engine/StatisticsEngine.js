export class StatisticsEngine {
  /**
   * Menghitung ringkasan kehadiran dari kumpulan record absensi.
   * @param {Array<{ students: Array<{ status: string }> }>} attendanceRecords
   * @returns {{ totalSessions: number, totalStudents: number, hadir: number, izin: number, sakit: number, alpha: number, persentaseKehadiran: number }}
   */
  static calculateAttendanceStats(attendanceRecords) {
    let totalStudents = 0,
      hadir = 0,
      izin = 0,
      sakit = 0,
      alpha = 0
    for (const rec of attendanceRecords) {
      for (const s of rec.students) {
        totalStudents++
        if (s.status === 'hadir') hadir++
        else if (s.status === 'izin') izin++
        else if (s.status === 'sakit') sakit++
        else if (s.status === 'alpha') alpha++
      }
    }
    const persentaseKehadiran = totalStudents > 0 ? (hadir / totalStudents) * 100 : 0
    return {
      totalSessions: attendanceRecords.length,
      totalStudents,
      hadir,
      izin,
      sakit,
      alpha,
      persentaseKehadiran: Math.round(persentaseKehadiran * 100) / 100,
    }
  }

  /**
   * Menghitung statistik kehadiran guru berdasarkan data presence.
   * @param {Array<{ status: string, checkInTime: string }>} teacherPresences - status 'on_time' | 'late' | 'absent'
   * @returns {{ total: number, onTime: number, late: number, absent: number, persentaseTepatWaktu: number }}
   */
  static calculateTeacherPresenceStats(teacherPresences) {
    let onTime = 0,
      late = 0,
      absent = 0
    for (const p of teacherPresences) {
      if (p.status === 'on_time') onTime++
      else if (p.status === 'late') late++
      else if (p.status === 'absent') absent++
    }
    const total = onTime + late + absent
    const persentaseTepatWaktu = total > 0 ? (onTime / total) * 100 : 0
    return {
      total,
      onTime,
      late,
      absent,
      persentaseTepatWaktu: Math.round(persentaseTepatWaktu * 100) / 100,
    }
  }

  /**
   * Menghitung durasi mengajar total dan rata-rata dari sesi.
   * @param {Array<{ startTime: string, endTime: string }>} sessions - startTime/endTime format 'HH:mm'
   * @returns {{ totalMinutes: number, totalSessions: number, averageMinutes: number }}
   */
  static calculateTeachingDuration(sessions) {
    let totalMinutes = 0
    for (const s of sessions) {
      if (s.startTime && s.endTime) {
        const [sh, sm] = s.startTime.split(':').map(Number)
        const [eh, em] = s.endTime.split(':').map(Number)
        totalMinutes += eh * 60 + em - (sh * 60 + sm)
      }
    }
    return {
      totalMinutes,
      totalSessions: sessions.length,
      averageMinutes: sessions.length > 0 ? totalMinutes / sessions.length : 0,
    }
  }

  /**
   * Menghitung persentase kehadiran per siswa dari kumpulan record.
   * @param {Array<{ studentId: number|string, studentName: string }>} students
   * @param {Array<{ students: Array<{ studentId: number|string, status: string }> }>} records
   * @returns {Array<{ studentId, studentName, total, hadir, izin, sakit, alpha, persentase }>}
   */
  static calculatePerStudentAttendance(students, records) {
    const map = new Map()
    students.forEach((s) =>
      map.set(s.studentId, { ...s, total: 0, hadir: 0, izin: 0, sakit: 0, alpha: 0 }),
    )
    for (const rec of records) {
      for (const s of rec.students) {
        const entry = map.get(s.studentId)
        if (entry) {
          entry.total++
          if (s.status === 'hadir') entry.hadir++
          else if (s.status === 'izin') entry.izin++
          else if (s.status === 'sakit') entry.sakit++
          else if (s.status === 'alpha') entry.alpha++
        }
      }
    }
    return [...map.values()].map((entry) => ({
      ...entry,
      persentase: entry.total > 0 ? Math.round((entry.hadir / entry.total) * 10000) / 100 : 0,
    }))
  }

  /**
   * Mendeteksi guru paling disiplin / sering terlambat dari data presence.
   * @param {Array<{ teacherId: number|string, teacherName: string, status: string }>} presences
   * @returns {{ mostOnTime: { teacherId, teacherName, onTimeCount }, mostLate: { teacherId, teacherName, lateCount } }}
   */
  static findTeacherExtremes(presences) {
    const stats = new Map()
    presences.forEach((p) => {
      if (!stats.has(p.teacherId)) {
        stats.set(p.teacherId, {
          teacherId: p.teacherId,
          teacherName: p.teacherName,
          onTime: 0,
          late: 0,
        })
      }
      const entry = stats.get(p.teacherId)
      if (p.status === 'on_time') entry.onTime++
      else if (p.status === 'late') entry.late++
    })
    let mostOnTime = { teacherId: null, teacherName: '', onTimeCount: 0 }
    let mostLate = { teacherId: null, teacherName: '', lateCount: 0 }
    stats.forEach((entry) => {
      if (entry.onTime > mostOnTime.onTimeCount) {
        mostOnTime = {
          teacherId: entry.teacherId,
          teacherName: entry.teacherName,
          onTimeCount: entry.onTime,
        }
      }
      if (entry.late > mostLate.lateCount) {
        mostLate = {
          teacherId: entry.teacherId,
          teacherName: entry.teacherName,
          lateCount: entry.late,
        }
      }
    })
    return { mostOnTime, mostLate }
  }

  /**
   * Menghitung efektivitas jam mengajar: jam terjadwal vs jam aktual.
   * @param {Array<{ scheduledMinutes: number, actualMinutes: number }>} sessions
   * @returns {{ totalScheduled: number, totalActual: number, efektivitasPersen: number }}
   */
  static calculateTeachingEffectiveness(sessions) {
    let totalScheduled = 0,
      totalActual = 0
    sessions.forEach((s) => {
      totalScheduled += s.scheduledMinutes
      totalActual += s.actualMinutes
    })
    const efektivitas = totalScheduled > 0 ? (totalActual / totalScheduled) * 100 : 0
    return {
      totalScheduled,
      totalActual,
      efektivitasPersen: Math.round(efektivitas * 100) / 100,
    }
  }

  /**
   * Membuat ringkasan dari audit trail.
   * @param {Array<{ action: string }>} auditEntries
   * @returns {Object} count per action
   */
  static summarizeAudit(auditEntries) {
    const summary = {}
    auditEntries.forEach((entry) => {
      summary[entry.action] = (summary[entry.action] || 0) + 1
    })
    return summary
  }
}

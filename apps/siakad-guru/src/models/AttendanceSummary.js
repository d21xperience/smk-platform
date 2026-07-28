// src/models/AttendanceSummary.js

export class AttendanceSummary {
  constructor({ totalSessions, hadir, izin, sakit, alpha, persentaseKehadiran }) {
    this.totalSessions = totalSessions
    this.hadir = hadir
    this.izin = izin
    this.sakit = sakit
    this.alpha = alpha
    this.persentaseKehadiran = persentaseKehadiran
  }
}

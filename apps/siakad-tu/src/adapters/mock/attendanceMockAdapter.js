// apps/siakad-tu/src/adapters/mock/attendanceMockAdapter.js

import { mockStorage } from './mockStorage.js'
import { idGenerator } from '../utils/idGenerator.js'

const _simulateLatency = (ms = 50) => new Promise((resolve) => setTimeout(resolve, ms))
const _error = (code, message, details = null) => ({ success: false, data: null, error: { code, message, details } })
const _success = (data) => ({ success: true, data, error: null })

export const attendanceMockAdapter = {
  async openSession(command, context) {
    await _simulateLatency(100)

    const sessionData = {
      sessionId: idGenerator.sessionId(),
      classId: command.classId,
      subjectId: command.subjectId,
      periodId: context.periodId,
      date: command.date || new Date().toISOString().split('T')[0],
      status: 'OPEN',
      schoolId: context.schoolId,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    }

    const sessions = mockStorage.read(context.schoolId, 'attendanceSessions')
    sessions.push(sessionData)
    mockStorage.write(context.schoolId, 'attendanceSessions', sessions)

    return _success(sessionData)
  },

  async closeSession(sessionId, context) {
    await _simulateLatency(100)

    const sessions = mockStorage.read(context.schoolId, 'attendanceSessions')
    const index = sessions.findIndex(s => s.sessionId === sessionId)

    if (index === -1) return _error('NOT_FOUND', 'Sesi absensi tidak ditemukan.')
    if (sessions[index].status === 'CLOSED') return _error('ALREADY_CLOSED', 'Sesi sudah ditutup.')

    sessions[index].status = 'CLOSED'
    sessions[index].updatedAt = new Date().toISOString()
    mockStorage.write(context.schoolId, 'attendanceSessions', sessions)

    return _success(sessions[index])
  },

  async submitBulkAttendance(sessionId, command, context) {
    await _simulateLatency(150)

    const sessions = mockStorage.read(context.schoolId, 'attendanceSessions')
    const session = sessions.find(s => s.sessionId === sessionId)

    if (!session) return _error('NOT_FOUND', 'Sesi absensi tidak ditemukan.')
    if (session.status === 'CLOSED') return _error('SESSION_CLOSED', 'Sesi sudah ditutup, tidak dapat menambah absensi.')

    const records = mockStorage.read(context.schoolId, 'attendanceRecords')

    const newRecords = command.records.map(r => ({
      recordId: idGenerator.recordId(),
      sessionId: sessionId,
      studentId: r.studentId,
      status: r.status,
      note: r.note || null,
      recordedAt: new Date().toISOString(),
      recordedBy: context.userId || 'system',
    }))

    records.push(...newRecords)
    mockStorage.write(context.schoolId, 'attendanceRecords', records)

    return _success(newRecords)
  },

  async updateAttendanceRecord(recordId, command, context) {
    await _simulateLatency(100)

    const records = mockStorage.read(context.schoolId, 'attendanceRecords')
    const index = records.findIndex(r => r.recordId === recordId)

    if (index === -1) return _error('NOT_FOUND', 'Record absensi tidak ditemukan.')

    records[index].status = command.status
    records[index].note = command.note
    records[index].recordedAt = new Date().toISOString()
    records[index].recordedBy = context.userId || 'system'

    mockStorage.write(context.schoolId, 'attendanceRecords', records)

    return _success(records[index])
  },

  async getSessionById(sessionId, context) {
    await _simulateLatency(30)
    const sessions = mockStorage.read(context.schoolId, 'attendanceSessions')
    const data = sessions.find(s => s.sessionId === sessionId)

    if (!data) return _error('NOT_FOUND', 'Sesi absensi tidak ditemukan.')
    return _success(data)
  },

  async getSessions(context, filters = {}) {
    await _simulateLatency(100)
    let sessions = mockStorage.read(context.schoolId, 'attendanceSessions')

    // Filter by classId
    if (filters.classId) {
      sessions = sessions.filter(s => s.classId === filters.classId)
    }

    // Filter by subjectId
    if (filters.subjectId) {
      sessions = sessions.filter(s => s.subjectId === filters.subjectId)
    }

    // Filter by status
    if (filters.status) {
      sessions = sessions.filter(s => s.status === filters.status)
    }

    // Filter by date
    if (filters.date) {
      sessions = sessions.filter(s => s.date === filters.date)
    }

    // Filter by periodId
    if (filters.periodId) {
      sessions = sessions.filter(s => s.periodId === filters.periodId)
    }

    // Sort by date descending
    sessions.sort((a, b) => new Date(b.date) - new Date(a.date))

    const page = filters.page || 1
    const limit = filters.limit || 20
    const total = sessions.length
    const items = sessions.slice((page - 1) * limit, page * limit)

    return _success({ items, total, page, limit, totalPages: Math.ceil(total / limit) })
  },

  async getAttendanceRecords(sessionId, context) {
    await _simulateLatency(50)
    const records = mockStorage.read(context.schoolId, 'attendanceRecords')
    const sessionRecords = records.filter(r => r.sessionId === sessionId)

    return _success(sessionRecords)
  },

  async getAttendanceSummary(studentId, context) {
    await _simulateLatency(50)
    const records = mockStorage.read(context.schoolId, 'attendanceRecords')
    const sessions = mockStorage.read(context.schoolId, 'attendanceSessions')

    // Filter records untuk siswa ini di periodId tertentu
    const studentRecords = records.filter(r => {
      const session = sessions.find(s => s.sessionId === r.sessionId)
      return r.studentId === studentId && session?.periodId === context.periodId
    })

    const summary = {
      studentId,
      periodId: context.periodId,
      totalSessions: studentRecords.length,
      present: studentRecords.filter(r => r.status === 'PRESENT').length,
      absent: studentRecords.filter(r => r.status === 'ABSENT').length,
      sick: studentRecords.filter(r => r.status === 'SICK').length,
      permission: studentRecords.filter(r => r.status === 'PERMISSION').length,
      late: studentRecords.filter(r => r.status === 'LATE').length,
      attendanceRate: studentRecords.length > 0
        ? ((studentRecords.filter(r => r.status === 'PRESENT' || r.status === 'LATE').length / studentRecords.length) * 100).toFixed(2)
        : 0,
    }

    return _success(summary)
  },

  async getClassAttendanceSummary(classId, context) {
    await _simulateLatency(100)
    const records = mockStorage.read(context.schoolId, 'attendanceRecords')
    const sessions = mockStorage.read(context.schoolId, 'attendanceSessions')
    // const students = mockStorage.read(context.schoolId, 'students')

    // Filter sessions untuk kelas ini
    const classSessions = sessions.filter(s => s.classId === classId && s.periodId === context.periodId)

    // Dapatkan semua studentId yang unik dari records di sessions kelas ini
    const classSessionIds = classSessions.map(s => s.sessionId)
    const classRecords = records.filter(r => classSessionIds.includes(r.sessionId))

    // Group by studentId
    const studentSummaries = {}
    classRecords.forEach(r => {
      if (!studentSummaries[r.studentId]) {
        studentSummaries[r.studentId] = {
          studentId: r.studentId,
          periodId: context.periodId,
          totalSessions: 0,
          present: 0,
          absent: 0,
          sick: 0,
          permission: 0,
          late: 0,
        }
      }

      studentSummaries[r.studentId].totalSessions++
      if (r.status === 'PRESENT') studentSummaries[r.studentId].present++
      else if (r.status === 'ABSENT') studentSummaries[r.studentId].absent++
      else if (r.status === 'SICK') studentSummaries[r.studentId].sick++
      else if (r.status === 'PERMISSION') studentSummaries[r.studentId].permission++
      else if (r.status === 'LATE') studentSummaries[r.studentId].late++
    })

    // Hitung attendance rate
    const summaries = Object.values(studentSummaries).map(s => ({
      ...s,
      attendanceRate: s.totalSessions > 0
        ? ((s.present + s.late) / s.totalSessions * 100).toFixed(2)
        : 0,
    }))

    return _success(summaries)
  },
}

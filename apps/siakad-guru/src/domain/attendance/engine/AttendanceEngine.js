import { AttendanceSession } from '../models/AttendanceSession.js'
import { AttendanceRecord } from '../models/AttendanceRecord.js'
import {
  SESSION_STATUS,
  isValidAttendanceStatus,
  isValidSessionStatus,
  canSubmitSession
} from '../models/AttendanceStatus.js'

export class AttendanceEngine {
  static validateSession(session) {
    if (!session || !(session instanceof AttendanceSession)) {
      throw new Error('Invalid attendance session object')
    }
    if (!session.isValid()) {
      throw new Error('Attendance session is missing required fields')
    }
    if (!isValidSessionStatus(session.status)) {
      throw new Error(`Invalid session status: ${session.status}`)
    }
    return true
  }

  static validateRecord(record) {
    if (!record || !(record instanceof AttendanceRecord)) {
      throw new Error('Invalid attendance record object')
    }
    if (!record.isValid()) {
      throw new Error('Attendance record is missing required fields')
    }
    if (!isValidAttendanceStatus(record.status)) {
      throw new Error(`Invalid attendance status: ${record.status}`)
    }
    return true
  }

  static validateRecords(records) {
    if (!Array.isArray(records)) {
      throw new Error('Records must be an array')
    }
    if (records.length === 0) {
      throw new Error('At least one attendance record is required')
    }
    records.forEach(record => {
      AttendanceEngine.validateRecord(record)
    })
    return true
  }

  static submitSession(session) {
    AttendanceEngine.validateSession(session)

    if (!canSubmitSession(session.status)) {
      throw new Error(`Cannot submit session with status: ${session.status}`)
    }

    if (session.records.length === 0) {
      throw new Error('Cannot submit session without attendance records')
    }

    // Validate all records before submitting
    AttendanceEngine.validateRecords(session.records)

    session.status = SESSION_STATUS.SUBMITTED
    return session
  }

  static updateRecord(session, studentId, status, note) {
    AttendanceEngine.validateSession(session)

    if (session.isSubmitted()) {
      throw new Error('Cannot update records of a submitted session')
    }

    if (!isValidAttendanceStatus(status)) {
      throw new Error(`Invalid attendance status: ${status}`)
    }

    const record = session.records.find(r => r.studentId === studentId)
    if (!record) {
      throw new Error(`Student not found: ${studentId}`)
    }

    record.status = status
    if (note !== undefined) {
      record.note = note
    }

    return session
  }

  static calculateSummary(session) {
    AttendanceEngine.validateSession(session)

    return {
      totalStudents: session.getTotalStudents(),
      present: session.getPresentCount(),
      sick: session.getSickCount(),
      permission: session.getPermissionCount(),
      absent: session.getAbsentCount()
    }
  }
}

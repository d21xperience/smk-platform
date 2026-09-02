import { AttendanceEngine } from '../domain/attendance/engine/AttendanceEngine.js'
import { AttendanceSession } from '../domain/attendance/models/AttendanceSession.js'
import { AttendanceRecord } from '../domain/attendance/models/AttendanceRecord.js'
import { AttendanceDraftedEvent } from '../domain/attendance/events/AttendanceDrafted.js'
import { AttendanceSubmittedEvent } from '../domain/attendance/events/AttendanceSubmitted.js'

export class AttendanceService {
  constructor({ attendanceAdapter, eventDispatcher }) {
    this.attendanceAdapter = attendanceAdapter
    this.eventDispatcher = eventDispatcher
  }

  async loadStudentsByClass({ classId }) {
    return await this.attendanceAdapter.fetchStudentsByClass({ classId })
  }

  async loadAttendanceSession({ teachingSessionId }) {
    const rawData = await this.attendanceAdapter.loadAttendanceSession({ teachingSessionId })

    if (!rawData) {
      return null
    }

    const records = rawData.records.map(r => new AttendanceRecord(r))
    const session = new AttendanceSession({
      ...rawData,
      records
    })

    AttendanceEngine.validateSession(session)
    return session
  }

  async createAttendanceSession({
    teachingSessionId,
    classId,
    className,
    subjectId,
    subjectName,
    date,
    schoolId,
    academicYearId,
    semesterId
  }) {
    const students = await this.attendanceAdapter.fetchStudentsByClass({ classId })

    if (students.length === 0) {
      throw new Error('No students found for this class')
    }

    const rawData = await this.attendanceAdapter.createAttendanceSession({
      teachingSessionId,
      classId,
      className,
      subjectId,
      subjectName,
      date,
      students,
      schoolId,
      academicYearId,
      semesterId
    })

    const records = rawData.records.map(r => new AttendanceRecord(r))
    const session = new AttendanceSession({
      ...rawData,
      records
    })

    AttendanceEngine.validateSession(session)

    if (this.eventDispatcher) {
      this.eventDispatcher.dispatch(new AttendanceDraftedEvent(session))
    }

    return session
  }

  async saveDraft({ sessionId, records }) {
    const domainRecords = records.map(r => new AttendanceRecord(r))
    AttendanceEngine.validateRecords(domainRecords)

    const rawData = await this.attendanceAdapter.saveAttendanceDraft({
      sessionId,
      records: records.map(r => ({
        studentId: r.studentId,
        status: r.status,
        note: r.note
      }))
    })

    const updatedRecords = rawData.records.map(r => new AttendanceRecord(r))
    const session = new AttendanceSession({
      ...rawData,
      records: updatedRecords
    })

    AttendanceEngine.validateSession(session)

    if (this.eventDispatcher) {
      this.eventDispatcher.dispatch(new AttendanceDraftedEvent(session))
    }

    return session
  }

  async submitAttendance({ sessionId }) {
    const submittedData = await this.attendanceAdapter.submitAttendance({ sessionId })

    const records = submittedData.records.map(r => new AttendanceRecord(r))
    const session = new AttendanceSession({
      ...submittedData,
      records
    })

    AttendanceEngine.validateSession(session)
    AttendanceEngine.submitSession(session)

    if (this.eventDispatcher) {
      this.eventDispatcher.dispatch(new AttendanceSubmittedEvent(session))
    }

    return session
  }
}

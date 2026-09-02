import { EVENT_ATTENDANCE_DRAFTED } from '../../../events/Type.js'

export class AttendanceDraftedEvent {
  constructor(session) {
    this.type = EVENT_ATTENDANCE_DRAFTED
    this.payload = {
      sessionId: session.id,
      teachingSessionId: session.teachingSessionId,
      classId: session.classId,
      subjectId: session.subjectId,
      date: session.date,
      totalRecords: session.records.length,
      academicYearId: session.academicYearId,
      semesterId: session.semesterId,
      schoolId: session.schoolId,
      timestamp: new Date().toISOString()
    }
  }
}

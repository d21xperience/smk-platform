import { EVENT_ATTENDANCE_SUBMITTED } from '../../../events/Type.js'

export class AttendanceSubmittedEvent {
  constructor(session) {
    this.type = EVENT_ATTENDANCE_SUBMITTED
    this.payload = {
      sessionId: session.id,
      teachingSessionId: session.teachingSessionId,
      classId: session.classId,
      subjectId: session.subjectId,
      date: session.date,
      totalStudents: session.getTotalStudents(),
      present: session.getPresentCount(),
      sick: session.getSickCount(),
      permission: session.getPermissionCount(),
      absent: session.getAbsentCount(),
      academicYearId: session.academicYearId,
      semesterId: session.semesterId,
      schoolId: session.schoolId,
      timestamp: new Date().toISOString()
    }
  }
}

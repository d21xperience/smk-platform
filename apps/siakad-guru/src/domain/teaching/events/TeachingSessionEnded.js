import { EVENT_TEACHING_SESSION_ENDED } from '../../../events/Type.js'

export class TeachingSessionEndedEvent {
  constructor(session) {
    this.type = EVENT_TEACHING_SESSION_ENDED
    this.payload = {
      sessionId: session.id,
      scheduleId: session.scheduleId,
      classId: session.classId,
      subjectId: session.subjectId,
      teacherId: session.teacherId,
      teacherPresence: session.teacherPresence,
      notes: session.notes,
      academicYearId: session.academicYearId,
      semesterId: session.semesterId,
      schoolId: session.schoolId,
      timestamp: new Date().toISOString(),
    }
  }
}

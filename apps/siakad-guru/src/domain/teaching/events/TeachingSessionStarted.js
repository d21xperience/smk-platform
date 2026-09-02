import { EVENT_TEACHING_SESSION_STARTED } from '../../../events/Type.js'

export class TeachingSessionStartedEvent {
  constructor(session) {
    this.type = EVENT_TEACHING_SESSION_STARTED
    this.payload = {
      sessionId: session.id,
      scheduleId: session.scheduleId,
      classId: session.classId,
      subjectId: session.subjectId,
      teacherId: session.teacherId,
      teacherPresence: session.teacherPresence,
      academicYearId: session.academicYearId,
      semesterId: session.semesterId,
      schoolId: session.schoolId,
      timestamp: new Date().toISOString(),
    }
  }
}

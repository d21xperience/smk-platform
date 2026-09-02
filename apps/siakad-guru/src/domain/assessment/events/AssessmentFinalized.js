import { EVENT_ASSESSMENT_FINALIZED } from '../../../events/Type.js'

export class AssessmentFinalizedEvent {
  constructor(session, result) {
    this.type = EVENT_ASSESSMENT_FINALIZED
    this.payload = {
      sessionId: session.id,
      classId: session.classId,
      subjectId: session.subjectId,
      date: session.date,
      finalScore: result.finalScore,
      grade: result.grade,
      predicate: result.predicate,
      academicYearId: session.academicYearId,
      semesterId: session.semesterId,
      schoolId: session.schoolId,
      timestamp: new Date().toISOString(),
    }
  }
}

// src/models/AssessmentSession.js
export class AssessmentSession {
  constructor({
    id,
    teachingSessionId, // FK ke TeachingSession
    className,
    subject,
    academicYearId,
    semesterId,
    components = [], // array of AssessmentComponent
    status = 'draft', // 'draft' | 'submitted' | 'locked'
  }) {
    this.id = id
    this.teachingSessionId = teachingSessionId
    this.className = className
    this.subject = subject
    this.academicYearId = academicYearId
    this.semesterId = semesterId
    this.components = components
    this.status = status
  }
}

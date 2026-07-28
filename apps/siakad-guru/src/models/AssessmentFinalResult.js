// src/models/AssessmentFinalResult.js
export class AssessmentFinalResult {
  constructor({
    sessionId,
    className,
    subject,
    grades = [], // array of Grade
    generatedAt = null,
  }) {
    this.sessionId = sessionId
    this.className = className
    this.subject = subject
    this.grades = grades
    this.generatedAt = generatedAt
  }
}

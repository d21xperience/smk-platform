// src/models/ScoreEntry.js
export class ScoreEntry {
  constructor({
    id,
    componentId,
    studentId,
    studentName,
    score, // nilai numerik (0 - maxScore)
    notes = '', // catatan guru
  }) {
    this.id = id
    this.componentId = componentId
    this.studentId = studentId
    this.studentName = studentName
    this.score = score
    this.notes = notes
  }
}

export class CounselingRecord {
  constructor({
    id,
    studentId,
    counselorId,
    counselorName,
    sessionDate,
    topic,
    notes,
    followUp,
    createdAt = new Date().toISOString(),
  }) {
    this.id = id
    this.studentId = studentId
    this.counselorId = counselorId
    this.counselorName = counselorName
    this.sessionDate = sessionDate
    this.topic = topic
    this.notes = notes
    this.followUp = followUp
    this.createdAt = createdAt
  }
}

export class Violation {
  constructor({
    id,
    studentId,
    type, // 'ringan', 'sedang', 'berat'
    description,
    sanction,
    date,
    recordedBy,
    recordedAt = new Date().toISOString(),
  }) {
    this.id = id
    this.studentId = studentId
    this.type = type
    this.description = description
    this.sanction = sanction
    this.date = date
    this.recordedBy = recordedBy
    this.recordedAt = recordedAt
  }
}

export class Achievement {
  constructor({
    id,
    studentId,
    title,
    description,
    category, // 'akademik', 'non-akademik', 'organisasi', 'lainnya'
    date,
    recordedBy,
    recordedAt = new Date().toISOString(),
  }) {
    this.id = id
    this.studentId = studentId
    this.title = title
    this.description = description
    this.category = category
    this.date = date
    this.recordedBy = recordedBy
    this.recordedAt = recordedAt
  }
}

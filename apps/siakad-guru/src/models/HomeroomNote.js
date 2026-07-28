export class HomeroomNote {
  constructor({
    id,
    studentId,
    homeroomTeacherId,
    homeroomTeacherName,
    note,
    semesterId,
    createdAt = new Date().toISOString(),
  }) {
    this.id = id
    this.studentId = studentId
    this.homeroomTeacherId = homeroomTeacherId
    this.homeroomTeacherName = homeroomTeacherName
    this.note = note
    this.semesterId = semesterId
    this.createdAt = createdAt
  }
}

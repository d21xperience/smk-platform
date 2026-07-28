export class TeacherNote {
  constructor({
    id,
    studentId,
    teacherId,
    teacherName,
    subject,
    sessionId, // TeachingSession terkait
    note,
    type = 'umum', // 'umum', 'sikap', 'tugas', 'ujian'
    createdAt = new Date().toISOString(),
  }) {
    this.id = id
    this.studentId = studentId
    this.teacherId = teacherId
    this.teacherName = teacherName
    this.subject = subject
    this.sessionId = sessionId
    this.note = note
    this.type = type
    this.createdAt = createdAt
  }
}

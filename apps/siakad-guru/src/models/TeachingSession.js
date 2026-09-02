// apps/siakad-guru/src/models/TeachingSession.js

export class TeachingSession {
  constructor(data) {
    this.id = data.id
    this.teacherId = data.teacherId
    this.classId = data.classId
    this.subjectId = data.subjectId
    this.date = data.date
    this.startTime = data.startTime
    this.endTime = data.endTime
    this.status = data.status || 'draft' // draft, active, completed, cancelled
    this.material = data.material || ''
    this.attendanceSubmitted = data.attendanceSubmitted || false
    this.createdAt = data.createdAt
    this.updatedAt = data.updatedAt
    this.version = data.version || 1
  }

  // Method untuk mengecek apakah sesi dapat dimulai
  canStart() {
    return this.status === 'draft' || this.status === 'DRAFT'
  }
  canEnd() {
    return this.status === 'draft' || this.status === 'DRAFT'
  }

  // Method untuk mengecek apakah sesi sedang aktif
  isActive() {
    return this.status === 'active' || this.status === 'ACTIVE'
  }

  // Method untuk mengecek apakah sesi sudah selesai
  isCompleted() {
    return this.status === 'completed' || this.status === 'COMPLETED'
  }

  toJSON() {
    return {
      id: this.id,
      teacherId: this.teacherId,
      classId: this.classId,
      subjectId: this.subjectId,
      date: this.date,
      startTime: this.startTime,
      endTime: this.endTime,
      status: this.status,
      material: this.material,
      attendanceSubmitted: this.attendanceSubmitted,
      createdAt: this.createdAt,
      updatedAt: this.updatedAt,
      version: this.version,
    }
  }
}

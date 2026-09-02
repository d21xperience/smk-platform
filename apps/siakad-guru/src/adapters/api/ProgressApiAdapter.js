import { BaseApiAdapter } from './BaseApiAdapter.js'

export class ProgressApiAdapter extends BaseApiAdapter {
  constructor() {
    super('/progress')
  }

  async fetchStudentsByClass({ classId }) {
    return await this.get('/students', { classId })
  }

  async loadProgressByStudent({ studentId, academicYearId, semesterId }) {
    return await this.get(`/students/${studentId}/records`, {
      academicYearId,
      semesterId,
    })
  }

  async createProgress({
    studentId,
    studentName,
    classId,
    className,
    type,
    title,
    description,
    date,
    teacherId,
    teacherName,
    schoolId,
    academicYearId,
    semesterId,
  }) {
    return await this.post('/records', {
      studentId,
      studentName,
      classId,
      className,
      type,
      title,
      description,
      date,
      teacherId,
      teacherName,
      schoolId,
      academicYearId,
      semesterId,
    })
  }
}

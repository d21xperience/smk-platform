import { BaseApiAdapter } from './BaseApiAdapter.js'

export class AttendanceApiAdapter extends BaseApiAdapter {
  constructor() {
    super('/attendance')
  }

  async fetchStudentsByClass({ classId }) {
    return await this.get('/students', { classId })
  }

  async loadAttendanceSession({ teachingSessionId }) {
    return await this.get(`/sessions/${teachingSessionId}`)
  }

  async createAttendanceSession({
    teachingSessionId,
    classId,
    className,
    subjectId,
    subjectName,
    date,
    students,
    schoolId,
    academicYearId,
    semesterId,
  }) {
    return await this.post('/sessions', {
      teachingSessionId,
      classId,
      className,
      subjectId,
      subjectName,
      date,
      students,
      schoolId,
      academicYearId,
      semesterId,
    })
  }

  async saveAttendanceDraft({ sessionId, records }) {
    return await this.put(`/sessions/${sessionId}/draft`, {
      records,
    })
  }

  async submitAttendance({ sessionId }) {
    return await this.post(`/sessions/${sessionId}/submit`)
  }
}

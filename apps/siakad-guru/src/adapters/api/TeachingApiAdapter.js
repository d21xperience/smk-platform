import { BaseApiAdapter } from './BaseApiAdapter.js'

export class TeachingApiAdapter extends BaseApiAdapter {
  constructor() {
    super('/teaching')
  }

  async fetchSchedulesByContext({ schoolId, academicYearId, semesterId, teacherId }) {
    return await this.get('/schedules', {
      schoolId,
      academicYearId,
      semesterId,
      teacherId,
    })
  }

  async fetchSessionsByDate({ schoolId, academicYearId, semesterId, teacherId, date }) {
    return await this.get('/sessions', {
      schoolId,
      academicYearId,
      semesterId,
      teacherId,
      date,
    })
  }

  async createSession({ scheduleId, date, schoolId, academicYearId, semesterId, teacherId }) {
    return await this.post('/sessions', {
      scheduleId,
      date,
      schoolId,
      academicYearId,
      semesterId,
      teacherId,
    })
  }

  async startSession({ sessionId, teacherPresence }) {
    return await this.post(`/sessions/${sessionId}/start`, {
      teacherPresence,
    })
  }

  async endSession({ sessionId, notes }) {
    return await this.post(`/sessions/${sessionId}/end`, {
      notes,
    })
  }

  async cancelSession({ sessionId, reason }) {
    return await this.post(`/sessions/${sessionId}/cancel`, {
      reason,
    })
  }
}

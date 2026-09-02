import { BaseApiAdapter } from './BaseApiAdapter.js'

export class JournalApiAdapter extends BaseApiAdapter {
  constructor() {
    super('/journal')
  }

  async loadJournalByTeachingSession({ teachingSessionId }) {
    return await this.get(`/sessions/${teachingSessionId}`)
  }

  async createJournal({
    teachingSessionId,
    classId,
    className,
    subjectId,
    subjectName,
    date,
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
      schoolId,
      academicYearId,
      semesterId,
    })
  }

  async updateJournal({ journalId, material, activities, reflection }) {
    return await this.put(`/sessions/${journalId}`, {
      material,
      activities,
      reflection,
    })
  }

  async submitJournal({ journalId }) {
    return await this.post(`/sessions/${journalId}/submit`)
  }
}

import { api as apiClient } from '@/boot/axios'
import { TeacherSchedule } from '@/models/TeacherSchedule'
import { TeachingSession } from '@/models/TeachingSession'
import { JournalEntry } from '@/models/JournalEntry'

export const teachingApiAdapter = {
  async getTeacherSchedule(academicYearId, semesterId) {
    const { data } = await apiClient.get('/teaching/schedule', {
      params: { academicYearId, semesterId },
    })
    return data.map((item) => new TeacherSchedule(item))
  },

  async getTeachingSessions(date) {
    const { data } = await apiClient.get('/teaching/sessions', {
      params: { date },
    })
    return data.map((item) => new TeachingSession(item))
  },

  async startSession(sessionId) {
    const { data } = await apiClient.patch(`/teaching/sessions/${sessionId}/start`)
    return new TeachingSession(data)
  },

  async completeSession(sessionId) {
    const { data } = await apiClient.patch(`/teaching/sessions/${sessionId}/complete`)
    return new TeachingSession(data)
  },

  async saveJournal(sessionId, journalData) {
    const { data } = await apiClient.post(`/teaching/sessions/${sessionId}/journal`, journalData)
    return new JournalEntry(data)
  },

  async getJournalBySession(sessionId) {
    const { data } = await apiClient.get(`/teaching/sessions/${sessionId}/journal`)
    return data ? new JournalEntry(data) : null
  },
}

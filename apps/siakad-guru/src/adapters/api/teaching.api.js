import { api } from '@/boot/axios'

export const TeachingApi = {
  getTodaySchedule: (params) => api.get('/teaching/schedule', { params }),
  getSessionDetail: (sessionId) => api.get(`/teaching/session/${sessionId}`),
  getActiveSessions: (params) => api.get('/teaching/active', { params }),
  startSession: (sessionId) => api.post(`/teaching/session/${sessionId}/start`),
  endSession: (sessionId) => api.post(`/teaching/session/${sessionId}/end`),
  loadOrCreateJournal: (sessionId) => api.get(`/teaching/session/${sessionId}/journal`),
  getJournal: (journalId) => api.get(`/teaching/journal/${journalId}`),
  saveJournalDraft: (journalId, data) => api.post(`/teaching/journal/${journalId}/draft`, data),
  finalizeJournal: (journalId) => api.post(`/teaching/journal/${journalId}/finalize`),
}

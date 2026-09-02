// apps/siakad-guru/src/adapters/api/attendance.api.js
import { api } from '@/boot/axios'

export const AttendanceApi = {
  // Method untuk mendapatkan atau membuat sesi absensi
  getOrCreateSession: (params) => {
    // params: { classId, subjectId, date, sessionId? }
    return api.get('/attendance/session', { params })
  },
  getSession: (sessionId) => api.get(`/attendance/session/${sessionId}`),
  loadOrCreate: (params) => api.get('/attendance/session', { params }),
  upsertRecord: (sessionId, data) => api.post(`/attendance/session/${sessionId}/record`, data),
  submit: (sessionId) => api.post(`/attendance/session/${sessionId}/submit`),
  getSummary: (sessionId) => api.get(`/attendance/session/${sessionId}/summary`),
  saveDraft: (data) => api.post('/attendance/session/draft', data),
}

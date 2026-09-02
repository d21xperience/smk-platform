import { api } from '@/boot/axios'

export const AssessmentApi = {
  loadOrCreate: (params) => api.get('/assessment/session', { params }),
  getSession: (sessionId) => api.get(`/assessment/session/${sessionId}`),
  updateResults: (sessionId, results) =>
    api.post(`/assessment/session/${sessionId}/results`, { results }),
  finalize: (sessionId) => api.post(`/assessment/session/${sessionId}/finalize`),
  getComponents: (sessionId) => api.get(`/assessment/session/${sessionId}/components`),
}

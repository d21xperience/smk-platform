import { api as apiClient } from '@/boot/axios'
import { AssessmentSession } from '@/models/AssessmentSession'
import { AssessmentComponent } from '@/models/AssessmentComponent'
import { AssessmentFinalResult } from '@/models/AssessmentFinalResult'

export const assessmentApiAdapter = {
  async fetchOrCreateAssessmentSession(teachingSessionId, options) {
    const { data } = await apiClient.post('/assessment/sessions', { teachingSessionId, ...options })
    return new AssessmentSession({
      ...data,
      components: data.components.map((c) => new AssessmentComponent(c)),
    })
  },

  async saveAssessmentDraft(sessionData) {
    const { data } = await apiClient.put(`/assessment/sessions/${sessionData.id}`, sessionData)
    return new AssessmentSession({
      ...data,
      components: data.components.map((c) => new AssessmentComponent(c)),
    })
  },

  async submitAssessment(sessionId) {
    const { data } = await apiClient.post(`/assessment/sessions/${sessionId}/submit`)
    return new AssessmentSession({
      ...data,
      components: data.components.map((c) => new AssessmentComponent(c)),
    })
  },

  async getAssessmentComponents(sessionId) {
    const { data } = await apiClient.get(`/assessment/sessions/${sessionId}/components`)
    return data.map((c) => new AssessmentComponent(c))
  },

  async saveComponentScores(componentId, scores) {
    const { data } = await apiClient.put(`/assessment/components/${componentId}/scores`, { scores })
    return new AssessmentComponent(data)
  },

  async getAssessmentFinalResult(sessionId) {
    const { data } = await apiClient.get(`/assessment/sessions/${sessionId}/result`)
    return data ? new AssessmentFinalResult(data) : null
  },

  async getStudentsByClass(className) {
    const { data } = await apiClient.get('/students', { params: { className } })
    return data.map((s) => ({ studentId: s.id, studentName: s.name }))
  },
}

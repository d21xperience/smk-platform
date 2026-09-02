import { BaseApiAdapter } from './BaseApiAdapter.js'

export class AssessmentApiAdapter extends BaseApiAdapter {
  constructor() {
    super('/assessment')
  }

  async loadAssessmentByClass({ classId, subjectId, academicYearId, semesterId }) {
    return await this.get('/sessions', {
      classId,
      subjectId,
      academicYearId,
      semesterId,
    })
  }

  async loadAssessmentById({ assessmentId }) {
    return await this.get(`/sessions/${assessmentId}`)
  }

  async createAssessment({
    classId,
    className,
    subjectId,
    subjectName,
    date,
    components,
    schoolId,
    academicYearId,
    semesterId,
  }) {
    return await this.post('/sessions', {
      classId,
      className,
      subjectId,
      subjectName,
      date,
      components,
      schoolId,
      academicYearId,
      semesterId,
    })
  }

  async updateScore({ assessmentId, componentId, score }) {
    return await this.put(`/sessions/${assessmentId}/scores`, {
      componentId,
      score,
    })
  }

  async finalizeAssessment({ assessmentId }) {
    return await this.post(`/sessions/${assessmentId}/finalize`)
  }
}

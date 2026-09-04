// apps/siakad-tu/src/services/assessment/AssessmentQueryService.js

export class AssessmentQueryService {
  constructor({ adapter }) {
    if (!adapter) throw new Error('AssessmentQueryService: adapter wajib.')
    this.adapter = adapter
  }

  _formatError(error) {
    return {
      success: false,
      data: null,
      error: {
        code: error.code || 'QUERY_ERROR',
        message: error.message || 'Gagal mengambil data.',
        source: 'adapter',
        details: error.details || null,
      },
    }
  }

  async getAssessmentById(assessmentId, context) {
    const result = await this.adapter.getAssessmentById(assessmentId, context)
    if (!result.success) return this._formatError(result.error)
    return { success: true, data: result.data, error: null }
  }

  async getAssessments(context, filters = {}) {
    const result = await this.adapter.getAssessments(context, filters)
    if (!result.success) return this._formatError(result.error)
    return { success: true, data: result.data, error: null }
  }

  async getGradeSummary(studentId, subjectId, context) {
    const result = await this.adapter.getGradeSummary(studentId, subjectId, context)
    if (!result.success) return this._formatError(result.error)
    return { success: true, data: result.data, error: null }
  }

  async getStudentReportCard(studentId, context) {
    const result = await this.adapter.getStudentReportCard(studentId, context)
    if (!result.success) return this._formatError(result.error)
    return { success: true, data: result.data, error: null }
  }

  async getClassGradeSummary(classId, subjectId, context) {
    const result = await this.adapter.getClassGradeSummary(classId, subjectId, context)
    if (!result.success) return this._formatError(result.error)
    return { success: true, data: result.data, error: null }
  }
}

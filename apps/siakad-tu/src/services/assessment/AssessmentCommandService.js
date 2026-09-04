// apps/siakad-tu/src/services/assessment/AssessmentCommandService.js

export class AssessmentCommandService {
  constructor({ adapter }) {
    if (!adapter) throw new Error('AssessmentCommandService: adapter wajib.')
    this.adapter = adapter
  }

  _formatError(source, error) {
    return {
      success: false,
      data: null,
      error: {
        code: error.code || `${source}_ERROR`,
        message: error.message || 'Terjadi kesalahan.',
        source,
        details: error.details || null,
      },
    }
  }

  async inputScore(commandData, context) {
    if (!commandData.studentId) {
      return this._formatError('validation', {
        code: 'STUDENT_REQUIRED',
        message: 'Siswa wajib dipilih.',
      })
    }

    if (!commandData.subjectId) {
      return this._formatError('validation', {
        code: 'SUBJECT_REQUIRED',
        message: 'Mata pelajaran wajib dipilih.',
      })
    }

    if (commandData.score === undefined || commandData.score === null) {
      return this._formatError('validation', {
        code: 'SCORE_REQUIRED',
        message: 'Nilai wajib diisi.',
      })
    }

    const maxScore = commandData.maxScore || 100
    if (commandData.score < 0 || commandData.score > maxScore) {
      return this._formatError('validation', {
        code: 'INVALID_SCORE_RANGE',
        message: `Nilai harus antara 0 dan ${maxScore}.`,
      })
    }

    const result = await this.adapter.createAssessment(commandData, context)
    if (!result.success) return this._formatError('adapter', result.error)
    return { success: true, data: result.data, error: null }
  }

  async correctScore(commandData, context) {
    if (!commandData.assessmentId) {
      return this._formatError('validation', {
        code: 'ASSESSMENT_REQUIRED',
        message: 'Data penilaian wajib dipilih.',
      })
    }

    const maxScore = commandData.maxScore || 100
    if (commandData.score < 0 || commandData.score > maxScore) {
      return this._formatError('validation', {
        code: 'INVALID_SCORE_RANGE',
        message: `Nilai harus antara 0 dan ${maxScore}.`,
      })
    }

    const result = await this.adapter.updateAssessment(
      commandData.assessmentId,
      commandData,
      context,
    )
    if (!result.success) return this._formatError('adapter', result.error)
    return { success: true, data: result.data, error: null }
  }

  async deleteAssessment(assessmentId, context) {
    const result = await this.adapter.deleteAssessment(assessmentId, context)
    if (!result.success) return this._formatError('adapter', result.error)
    return { success: true, data: result.data, error: null }
  }
}

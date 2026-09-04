// apps/siakad-tu/src/services/attendance/AttendanceQueryService.js

export class AttendanceQueryService {
  constructor({ adapter }) {
    if (!adapter) throw new Error('AttendanceQueryService: adapter wajib.')
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

  async getSessionById(sessionId, context) {
    const result = await this.adapter.getSessionById(sessionId, context)
    if (!result.success) return this._formatError(result.error)
    return { success: true, data: result.data, error: null }
  }

  async getSessions(context, filters = {}) {
    const result = await this.adapter.getSessions(context, filters)
    if (!result.success) return this._formatError(result.error)
    return { success: true, data: result.data, error: null }
  }

  async getAttendanceRecords(sessionId, context) {
    const result = await this.adapter.getAttendanceRecords(sessionId, context)
    if (!result.success) return this._formatError(result.error)
    return { success: true, data: result.data, error: null }
  }

  async getAttendanceSummary(studentId, context) {
    const result = await this.adapter.getAttendanceSummary(studentId, context)
    if (!result.success) return this._formatError(result.error)
    return { success: true, data: result.data, error: null }
  }

  async getClassAttendanceSummary(classId, context) {
    const result = await this.adapter.getClassAttendanceSummary(classId, context)
    if (!result.success) return this._formatError(result.error)
    return { success: true, data: result.data, error: null }
  }
}

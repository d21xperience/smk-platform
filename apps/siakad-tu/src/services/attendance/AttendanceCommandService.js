// apps/siakad-tu/src/services/attendance/AttendanceCommandService.js

export class AttendanceCommandService {
  constructor({ adapter }) {
    if (!adapter) throw new Error('AttendanceCommandService: adapter wajib.')
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

  async openSession(commandData, context) {
    if (!commandData.classId) {
      return this._formatError('validation', {
        code: 'CLASS_REQUIRED',
        message: 'Kelas wajib dipilih.',
      })
    }

    if (!commandData.subjectId) {
      return this._formatError('validation', {
        code: 'SUBJECT_REQUIRED',
        message: 'Mata pelajaran wajib dipilih.',
      })
    }

    const result = await this.adapter.openSession(commandData, context)
    if (!result.success) return this._formatError('adapter', result.error)
    return { success: true, data: result.data, error: null }
  }

  async closeSession(sessionId, context) {
    const result = await this.adapter.closeSession(sessionId, context)
    if (!result.success) return this._formatError('adapter', result.error)
    return { success: true, data: result.data, error: null }
  }

  async submitBulkAttendance(commandData, context) {
    if (!commandData.records || commandData.records.length === 0) {
      return this._formatError('validation', {
        code: 'NO_RECORDS',
        message: 'Tidak ada data absensi yang dikirim.',
      })
    }

    const result = await this.adapter.submitBulkAttendance(commandData.sessionId, commandData, context)
    if (!result.success) return this._formatError('adapter', result.error)
    return { success: true, data: result.data, error: null }
  }

  async updateAttendanceRecord(commandData, context) {
    const result = await this.adapter.updateAttendanceRecord(commandData.recordId, commandData, context)
    if (!result.success) return this._formatError('adapter', result.error)
    return { success: true, data: result.data, error: null }
  }
}

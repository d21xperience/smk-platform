// apps/siakad-tu/src/services/registration/RegistrationCommandService.js

export class RegistrationCommandService {
  constructor({ adapter }) {
    if (!adapter) throw new Error('RegistrationCommandService: adapter wajib.')
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

  async createRegistration(commandData) {
    if (!commandData.calonNisn) {
      return this._formatError('validation', { code: 'NISN_REQUIRED', message: 'NISN wajib diisi.' })
    }
    if (!commandData.calonNama) {
      return this._formatError('validation', { code: 'NAME_REQUIRED', message: 'Nama lengkap wajib diisi.' })
    }

    const result = await this.adapter.createRegistration(commandData)
    if (!result.success) {
      return this._formatError('adapter', result.error)
    }
    return { success: true, data: result.data, error: null }
  }

  async generateNisMassal(prefix) {
    if (!prefix) {
      return this._formatError('validation', { code: 'PREFIX_REQUIRED', message: 'Prefix NIS wajib diisi.' })
    }
    const result = await this.adapter.generateNisMassal(prefix)
    if (!result.success) {
      return this._formatError('adapter', result.error)
    }
    return { success: true, data: result.data, error: null }
  }

  async activateRegistration(regId) {
    const result = await this.adapter.activateRegistration(regId)
    if (!result.success) {
      return this._formatError('adapter', result.error)
    }
    return { success: true, data: result.data, error: null }
  }
}

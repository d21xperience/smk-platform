// apps/siakad-tu/src/services/registration/RegistrationCommandService.js

/**
 * RegistrationCommandService — Orkestrator untuk write operations pada Registration Domain.
 *
 * TANGGUNG JAWAB:
 * 1. Validasi business rules sebelum eksekusi
 * 2. Memanggil Adapter untuk persist data
 * 3. Return hasil yang konsisten
 *
 * TIDAK BOLEH:
 * - Mengandung logika presentasi (itu tugas Store/Page)
 * - Akses database langsung (itu tugas Adapter)
 * - Akses global state (context di-inject sebagai parameter)
 */
export class RegistrationCommandService {
  constructor({ adapter }) {
    if (!adapter) throw new Error('RegistrationCommandService: adapter wajib.')
    this.adapter = adapter
  }

  /**
   * Helper: format error response yang konsisten
   */
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

  /**
   * USE CASE 1: Buat Draft Registrasi
   */
  async createDraft(commandData, context) {
    // Validasi: NISN wajib ada
    if (!commandData.nisn) {
      return this._formatError('validation', {
        code: 'NISN_REQUIRED',
        message: 'NISN wajib diisi.',
      })
    }

    // Validasi: Nama wajib ada
    if (!commandData.fullName) {
      return this._formatError('validation', {
        code: 'NAME_REQUIRED',
        message: 'Nama lengkap wajib diisi.',
      })
    }

    const result = await this.adapter.createDraft(commandData, context)

    if (!result.success) {
      return this._formatError('adapter', result.error)
    }

    return { success: true, data: result.data, error: null }
  }

  /**
   * USE CASE 2: Submit Draft untuk Verifikasi
   */
  async submitRegistration(registrationId, context) {
    const result = await this.adapter.submitRegistration(registrationId, context)

    if (!result.success) {
      return this._formatError('adapter', result.error)
    }

    return { success: true, data: result.data, error: null }
  }

  /**
   * USE CASE 3: Approve Registrasi (Trigger auto-create Student)
   */
  async approveRegistration(registrationId, commandData, context) {
    // Validasi: classId wajib untuk mutasi masuk
    if (!commandData.classId) {
      return this._formatError('validation', {
        code: 'CLASS_REQUIRED',
        message: 'Kelas tujuan wajib dipilih.',
      })
    }

    const result = await this.adapter.approveRegistration(registrationId, commandData, context)

    if (!result.success) {
      return this._formatError('adapter', result.error)
    }

    return { success: true, data: result.data, error: null }
  }

  /**
   * USE CASE 4: Reject Registrasi
   */
  async rejectRegistration(registrationId, commandData, context) {
    // Validasi: Alasan wajib diisi
    if (!commandData.reason || commandData.reason.trim() === '') {
      return this._formatError('validation', {
        code: 'REASON_REQUIRED',
        message: 'Alasan penolakan wajib diisi.',
      })
    }

    const result = await this.adapter.rejectRegistration(registrationId, commandData, context)

    if (!result.success) {
      return this._formatError('adapter', result.error)
    }

    return { success: true, data: result.data, error: null }
  }
}

// apps/siakad-tu/src/services/student/StudentCommandService.js

/**
 * StudentCommandService — Orkestrator untuk write operations.
 *
 * TANGGUNG JAWAB:
 * 1. Memanggil Engine untuk validasi & business logic
 * 2. Memanggil Adapter untuk persist data
 * 3. Dispatch Domain Events setelah persist berhasil
 * 4. Return hasil yang konsisten
 *
 * TIDAK BOLEH:
 * - Mengandung logika bisnis (itu tugas Engine)
 * - Akses database langsung (itu tugas Adapter)
 * - Akses global state (context di-inject sebagai parameter)
 *
 * Pola: Engine → Adapter → Event Dispatch → Return
 */
export class StudentCommandService {
  /**
   * @param {Object} dependencies
   * @param {import('@/engine/student/StudentEngine').StudentEngine} dependencies.engine
   * @param {Object} dependencies.adapter - studentApi (mock atau real)
   * @param {import('@/events/eventDispatcher').EventDispatcher} dependencies.eventDispatcher
   */
  constructor({ engine, adapter, eventDispatcher }) {
    if (!engine) throw new Error('StudentCommandService: engine wajib.')
    if (!adapter) throw new Error('StudentCommandService: adapter wajib.')
    if (!eventDispatcher) throw new Error('StudentCommandService: eventDispatcher wajib.')

    this.engine = engine
    this.adapter = adapter
    this.eventDispatcher = eventDispatcher
  }

  /**
   * Helper: dispatch semua events dari Engine
   * @private
   */
  _dispatchEvents(events) {
    if (!events || events.length === 0) return
    events.forEach((event) => {
      try {
        this.eventDispatcher.dispatch(event)
      } catch (error) {
        // Event dispatch error tidak boleh menggagalkan operasi utama
        console.error('[StudentCommandService] Error dispatching event:', error)
      }
    })
  }

  /**
   * Helper: format error response yang konsisten
   * @private
   */
  _formatError(source, error) {
    return {
      success: false,
      data: null,
      error: {
        code: error.code || `${source}_ERROR`,
        message: error.message || 'Terjadi kesalahan.',
        source, // 'engine' | 'adapter'
        details: error.details || null,
      },
    }
  }

  /**
   * USE CASE 1: Register Student
   *
   * Alur:
   * 1. Engine validasi & bangun Student aggregate
   * 2. Adapter persist ke storage (dengan validasi keunikan NISN/NIS)
   * 3. Dispatch StudentRegistered event
   */
  async registerStudent(commandData, context) {
    // 1. Engine: validasi & business logic
    const engineResult = this.engine.registerStudent(commandData, context)
    if (!engineResult.success) {
      return this._formatError('engine', engineResult.errors[0])
    }

    // 2. Adapter: persist data
    const adapterResult = await this.adapter.registerStudent(engineResult.data.student, context)
    if (!adapterResult.success) {
      return this._formatError('adapter', adapterResult.error)
    }

    // 3. Dispatch events (hanya jika persist berhasil)
    this._dispatchEvents(engineResult.events)

    return {
      success: true,
      data: adapterResult.data,
      error: null,
    }
  }

  /**
   * USE CASE 2: Enroll Student
   *
   * Alur:
   * 1. Fetch currentData dari Adapter (siswa harus ada)
   * 2. Engine validasi & tambahkan enrollment
   * 3. Adapter persist perubahan
   * 4. Dispatch StudentEnrolled event
   */
  async enrollStudent(commandData, context) {
    // 1. Fetch current data
    const currentDataResult = await this.adapter.getStudentById(commandData.studentId, context)
    if (!currentDataResult.success) {
      return this._formatError('adapter', currentDataResult.error)
    }

    // 2. Engine: validasi & business logic
    const engineResult = this.engine.enrollStudent(commandData, context, currentDataResult.data)
    if (!engineResult.success) {
      return this._formatError('engine', engineResult.errors[0])
    }

    // 3. Adapter: persist
    const adapterResult = await this.adapter.enrollStudent(
      {
        studentId: commandData.studentId,
        enrollmentId: engineResult.data.student.enrollments.at(-1).enrollmentId,
        classId: commandData.classId,
      },
      context,
    )
    if (!adapterResult.success) {
      return this._formatError('adapter', adapterResult.error)
    }

    // 4. Dispatch events
    this._dispatchEvents(engineResult.events)

    return {
      success: true,
      data: adapterResult.data,
      error: null,
    }
  }

  /**
   * USE CASE 3: Update Student Profile
   *
   * Alur:
   * 1. Fetch currentData
   * 2. Engine validasi & update profile
   * 3. Adapter persist
   * 4. Dispatch StudentProfileUpdated event
   */
  async updateStudentProfile(commandData, context) {
    const currentDataResult = await this.adapter.getStudentById(commandData.studentId, context)
    if (!currentDataResult.success) {
      return this._formatError('adapter', currentDataResult.error)
    }

    const engineResult = this.engine.updateStudentProfile(
      commandData,
      context,
      currentDataResult.data,
    )
    if (!engineResult.success) {
      return this._formatError('engine', engineResult.errors[0])
    }

    const adapterResult = await this.adapter.updateStudentProfile(commandData, context)
    if (!adapterResult.success) {
      return this._formatError('adapter', adapterResult.error)
    }

    this._dispatchEvents(engineResult.events)

    return {
      success: true,
      data: adapterResult.data,
      error: null,
    }
  }

  /**
   * USE CASE 4: Graduate Student
   *
   * Alur:
   * 1. Fetch currentData
   * 2. Engine validasi & graduate
   * 3. Adapter persist
   * 4. Dispatch StudentGraduated event
   */
  async graduateStudent(commandData, context) {
    const currentDataResult = await this.adapter.getStudentById(commandData.studentId, context)
    if (!currentDataResult.success) {
      return this._formatError('adapter', currentDataResult.error)
    }

    const engineResult = this.engine.graduateStudent(commandData, context, currentDataResult.data)
    if (!engineResult.success) {
      return this._formatError('engine', engineResult.errors[0])
    }

    const adapterResult = await this.adapter.graduateStudent(commandData, context)
    if (!adapterResult.success) {
      return this._formatError('adapter', adapterResult.error)
    }

    this._dispatchEvents(engineResult.events)

    return {
      success: true,
      data: adapterResult.data,
      error: null,
    }
  }

  /**
   * USE CASE 5: Transfer Student
   *
   * Alur:
   * 1. Fetch currentData
   * 2. Engine validasi & transfer
   * 3. Adapter persist
   * 4. Dispatch StudentTransferred event
   */
  async transferStudent(commandData, context) {
    const currentDataResult = await this.adapter.getStudentById(commandData.studentId, context)
    if (!currentDataResult.success) {
      return this._formatError('adapter', currentDataResult.error)
    }

    const engineResult = this.engine.transferStudent(commandData, context, currentDataResult.data)
    if (!engineResult.success) {
      return this._formatError('engine', engineResult.errors[0])
    }

    const adapterResult = await this.adapter.transferStudent(commandData, context)
    if (!adapterResult.success) {
      return this._formatError('adapter', adapterResult.error)
    }

    this._dispatchEvents(engineResult.events)

    return {
      success: true,
      data: adapterResult.data,
      error: null,
    }
  }
}

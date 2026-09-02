// FILE: src/domain/context/engine/ContextEngine.js
// STATUS: MODIFY
// STATUS IMPLEMENTASI: COMPLETE

import { OperationalContext } from '../models/OperationalContext.js'

export const CONTEXT_MODE = Object.freeze({
  CURRENT: 'current',
  HISTORICAL: 'historical',
})

export const HISTORY_MUTATION_PERMISSION = 'context.history.write'

export class ContextEngine {
  static validate(context) {
    this.validateContext(context)
  }

  static validateContext(context) {
    if (!context) {
      throw new Error('Context tidak boleh kosong')
    }

    const isValid =
      typeof context.isValid === 'function'
        ? context.isValid()
        : this.isValidContextPayload(context)

    if (!isValid) {
      throw new Error('Operational context tidak valid')
    }
  }

  static isValidContextPayload(context) {
    return Boolean(
      context &&
      context.schoolId &&
      context.academicYearId &&
      context.semesterId &&
      context.userId &&
      context.role,
    )
  }

  static validateUser(user) {
    if (!user || !user.id || !user.role) {
      throw new Error('User tidak valid untuk membangun context')
    }
  }

  static validateAcademicPeriod({ academicYearId, semesterId }) {
    if (!academicYearId || !semesterId) {
      throw new Error('Tahun pelajaran dan semester wajib dipilih')
    }
  }

  static buildCurrentContext({ schoolId, academicYearId, semesterId, user }) {
    this.validateUser(user)
    this.validateAcademicPeriod({ academicYearId, semesterId })

    if (!schoolId) {
      throw new Error('School context tidak valid')
    }

    return new OperationalContext({
      schoolId,
      academicYearId,
      semesterId,
      userId: user.id,
      role: user.role,
      permissions: Array.isArray(user.permissions) ? user.permissions : [],
    })
  }

  static buildHistoricalContext({ academicYearId, semesterId }) {
    this.validateAcademicPeriod({ academicYearId, semesterId })

    return Object.freeze({
      academicYearId,
      semesterId,
      mode: CONTEXT_MODE.HISTORICAL,
    })
  }

  static validateHistoricalContext(historicalContext) {
    if (!historicalContext) {
      throw new Error('Historical context tidak boleh kosong')
    }

    this.validateAcademicPeriod(historicalContext)
  }

  static isSameAcademicPeriod(left, right) {
    return Boolean(
      left &&
      right &&
      left.academicYearId === right.academicYearId &&
      left.semesterId === right.semesterId,
    )
  }

  static isHistoricalSameAsCurrent({ currentContext, historicalContext }) {
    return this.isSameAcademicPeriod(currentContext, historicalContext)
  }

  static resolveRequestContext({ isHistoryMode, currentContext, historicalContext }) {
    if (isHistoryMode && historicalContext) {
      this.validateHistoricalContext(historicalContext)

      return Object.freeze({
        mode: CONTEXT_MODE.HISTORICAL,
        context: Object.freeze({
          schoolId: currentContext && currentContext.schoolId ? currentContext.schoolId : null,
          academicYearId: historicalContext.academicYearId,
          semesterId: historicalContext.semesterId,
        }),
      })
    }

    this.validateContext(currentContext)

    return Object.freeze({
      mode: CONTEXT_MODE.CURRENT,
      context: Object.freeze({
        schoolId: currentContext.schoolId,
        academicYearId: currentContext.academicYearId,
        semesterId: currentContext.semesterId,
        userId: currentContext.userId,
        role: currentContext.role,
      }),
    })
  }

  static canViewHistory({ user }) {
    return Boolean(user && user.id)
  }

  static canMutateHistoricalContext({ user }) {
    if (!user) return false

    const isAdmin = user.role === 'admin' || user.role === 'superadmin'
    const permissions = Array.isArray(user.permissions) ? user.permissions : []

    return isAdmin && permissions.includes(HISTORY_MUTATION_PERMISSION)
  }

  static canMutateContext({ isHistoryMode, user }) {
    if (!user || !user.id) return false

    if (!isHistoryMode) {
      return true
    }

    return this.canMutateHistoricalContext({ user })
  }

  static assertMutationAllowed({ isHistoryMode, user }) {
    if (!this.canMutateContext({ isHistoryMode, user })) {
      throw new Error('Data historis bersifat read-only untuk context ini')
    }
  }
}

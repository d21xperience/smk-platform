// FILE: src/services/ContextService.js
// STATUS: MODIFY
// STATUS IMPLEMENTASI: COMPLETE

import { ContextEngine } from '@/domain/context/engine/ContextEngine.js'
import { ContextSelectedEvent } from '@/domain/context/events/ContextSelected.js'
import { normalizeAvailableHistory } from '@/contracts/contextContract.js'

export class ContextService {
  constructor({ contextAdapter, eventDispatcher }) {
    this.contextAdapter = contextAdapter
    this.eventDispatcher = eventDispatcher
  }

  async loadAvailableContexts() {
    return await this.contextAdapter.fetchAvailableContexts()
  }

  async selectContext({ schoolId, academicYearId, semesterId, user }) {
    const context = await this.contextAdapter.selectContext({
      schoolId,
      academicYearId,
      semesterId,
      user,
    })

    ContextEngine.validate(context)
    this.dispatchContextSelected(context)

    return context
  }

  async loadCurrentOperationalContext({ user }) {
    const context = await this.contextAdapter.fetchActiveOperationalContext({ user })

    ContextEngine.validate(context)
    this.dispatchContextSelected(context)

    return context
  }

  async loadAvailableHistory({ user }) {
    const history = await this.contextAdapter.fetchAvailableHistory({ user })

    return normalizeAvailableHistory(history)
  }

  async enterHistoricalContext({ user, academicYearId, semesterId }) {
    const historicalContext = await this.contextAdapter.fetchHistoricalContext({
      user,
      academicYearId,
      semesterId,
    })

    ContextEngine.validateHistoricalContext(historicalContext)

    return historicalContext
  }

  resolveRequestContext({ isHistoryMode, currentContext, historicalContext }) {
    return ContextEngine.resolveRequestContext({
      isHistoryMode,
      currentContext,
      historicalContext,
    })
  }

  assertMutationAllowed({ isHistoryMode, user }) {
    ContextEngine.assertMutationAllowed({ isHistoryMode, user })
  }

  dispatchContextSelected(context) {
    if (!this.eventDispatcher) return

    this.eventDispatcher.dispatch(new ContextSelectedEvent(context))
  }
}

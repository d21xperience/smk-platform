import { mutationApi } from '@/adapters/api/mutationApi.js'
// import { mutationEngine } from '@/engine/mutation/MutationEngine.js'
// import { eventDispatcher } from '@/events/eventDispatcher.js'

export class MutationCommandService {
  constructor({ engine, adapter = mutationApi, eventDispatcher }) {
    this.engine = engine
    this.adapter = adapter
    this.eventDispatcher = eventDispatcher
  }

  async approveMutation(commandData, context) {
    // 1. Engine: Validasi business rule (misal: targetClassId harus ada untuk mutasi masuk)
    // const engineResult = this.engine.validateApproval(commandData, context)
    // if (!engineResult.success) return { success: false, error: engineResult.error }

    // 2. Adapter: Eksekusi API (Mock atau Real).
    // Backend akan menangani side-effect (update student status) secara atomik.
    const adapterResult = await this.adapter.approveMutation(
      commandData.mutationId,
      commandData,
      context,
    )

    if (!adapterResult.success) {
      return {
        success: false,
        error: {
          code: adapterResult.error.code,
          message: adapterResult.error.message,
          source: 'adapter',
        },
      }
    }

    // 3. Dispatch Event (misal: 'MUTATION_APPROVED' untuk trigger notifikasi/email)
    // this.eventDispatcher.dispatch({ type: 'MUTATION_APPROVED', payload: adapterResult.data })

    return { success: true, data: adapterResult.data, error: null }
  }
}

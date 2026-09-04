// apps/siakad-tu/src/services/mutation/MutationQueryService.js

export class MutationQueryService {
  constructor({ adapter }) {
    if (!adapter) throw new Error('MutationQueryService: adapter wajib.')
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

  async getMutationById(mutationId, context) {
    const result = await this.adapter.getMutationById(mutationId, context)
    if (!result.success) return this._formatError(result.error)
    return { success: true, data: result.data, error: null }
  }

  async getMutations(context, filters = {}) {
    const result = await this.adapter.getMutations(context, filters)
    if (!result.success) return this._formatError(result.error)
    return { success: true, data: result.data, error: null }
  }
}

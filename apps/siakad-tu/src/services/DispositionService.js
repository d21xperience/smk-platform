export function createDispositionService({ adapter }) {
  return {
    async getPendingDispositions(context) {
      return adapter.listPendingForKepsek(context)
    },

    async getMyTasks(context, userId) {
      return adapter.listTasksForStaf(context, userId)
    },

    async createDisposition(context, payload) {
      return adapter.create(context, payload)
    },

    async submitFollowUp(context, dispositionId, payload) {
      return adapter.submitFollowUp(context, dispositionId, payload)
    },

    async getStafOptions(context) {
      return adapter.getStafOptions(context)
    },
  }
}

export function createIncomingLetterService({ adapter }) {
  return {
    async listIncomingLetters(context, query) {
      return adapter.list(context, query)
    },

    async getIncomingLetterDetail(context, id) {
      return adapter.getById(context, id)
    },

    async createIncomingLetter(context, payload) {
      // Bisa ditambahkan validasi domain di sini jika perlu
      return adapter.create(context, payload)
    },

    async updateIncomingLetterStatus(context, id, newStatus) {
      return adapter.updateStatus(context, id, newStatus)
    },
  }
}

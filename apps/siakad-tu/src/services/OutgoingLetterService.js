export function createOutgoingLetterService({ adapter }) {
  return {
    async listOutgoingLetters(context, kategori, query) {
      return adapter.listByCategory(context, kategori, query)
    },

    async getOutgoingLetterDetail(context, id) {
      return adapter.getById(context, id)
    },

    async createOutgoingLetter(context, payload) {
      return adapter.create(context, payload)
    },

    async requestNextLetterNumber(context, kategori, jenis) {
      // Mencegah race condition: nomor di-generate dan di-lock di sisi "backend" (adapter)
      return adapter.getNextLetterNumber(context, kategori, jenis)
    },

    async getStudentOptions(context) {
      return adapter.getStudentOptions(context)
    },

    async getDudiOptions(context) {
      return adapter.getDudiOptions(context)
    },
  }
}

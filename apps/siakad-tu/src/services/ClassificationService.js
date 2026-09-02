export function createClassificationService({ adapter }) {
  return {
    async getClassificationTree(context) {
      return adapter.getTree(context)
    },

    async createClassificationNode(context, payload) {
      return adapter.createNode(context, payload)
    },

    async updateClassificationNode(context, id, payload) {
      return adapter.updateNode(context, id, payload)
    },

    async deleteClassificationNode(context, id) {
      return adapter.deleteNode(context, id)
    },

    async listArchiveLocations(context) {
      return adapter.listLocations(context)
    },

    async createArchiveLocation(context, payload) {
      return adapter.createLocation(context, payload)
    },

    // BARU
    async updateArchiveLocation(context, id, payload) {
      return adapter.updateLocation(context, id, payload)
    },

    // BARU
    async deleteArchiveLocation(context, id) {
      return adapter.deleteLocation(context, id)
    },

    async getKlasifikasiOptions(context) {
      return adapter.getKlasifikasiOptions(context)
    }
  }
}

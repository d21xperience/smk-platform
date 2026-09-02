import { defineStore } from 'pinia'
import { createClassificationService } from '../services/ClassificationService.js'
import { ClassificationMockAdapter } from '../adapters/mock/ClassificationMockAdapter.js'

const service = createClassificationService({ adapter: ClassificationMockAdapter })

export const useClassificationStore = defineStore('classification', {
  state: () => ({
    tree: [],
    locations: [],
    klasifikasiOptions: [],
    isLoading: false,
    error: null,
  }),
  actions: {
    async loadTree(context) {
      this.isLoading = true
      try {
        this.tree = await service.getClassificationTree(context)
      } catch (e) {
        this.error = e.message
      } finally {
        this.isLoading = false
      }
    },
    async createNode(context, payload) {
      this.isLoading = true
      try {
        await service.createClassificationNode(context, payload)
        await this.loadTree(context)
      } catch (e) {
        this.error = e.message
        throw e
      } finally {
        this.isLoading = false
      }
    },
    async updateNode(context, id, payload) {
      this.isLoading = true
      try {
        await service.updateClassificationNode(context, id, payload)
        await this.loadTree(context)
      } catch (e) {
        this.error = e.message
        throw e
      } finally {
        this.isLoading = false
      }
    },
    async deleteNode(context, id) {
      this.isLoading = true
      try {
        await service.deleteClassificationNode(context, id)
        await this.loadTree(context)
      } catch (e) {
        this.error = e.message
        throw e
      } finally {
        this.isLoading = false
      }
    },
    async loadLocations(context) {
      this.isLoading = true
      try {
        this.locations = await service.listArchiveLocations(context)
      } catch (e) {
        this.error = e.message
      } finally {
        this.isLoading = false
      }
    },
    async createLocation(context, payload) {
      this.isLoading = true
      try {
        await service.createArchiveLocation(context, payload)
        await this.loadLocations(context)
      } catch (e) {
        this.error = e.message
        throw e
      } finally {
        this.isLoading = false
      }
    },
    async updateLocation(context, id, payload) {
      this.isLoading = true
      try {
        await service.updateArchiveLocation(context, id, payload)
        await this.loadLocations(context)
      } catch (e) {
        this.error = e.message
        throw e
      } finally {
        this.isLoading = false
      }
    },
    async deleteLocation(context, id) {
      this.isLoading = true
      try {
        await service.deleteArchiveLocation(context, id)
        await this.loadLocations(context)
      } catch (e) {
        this.error = e.message
        throw e
      } finally {
        this.isLoading = false
      }
    },
    async loadKlasifikasiOptions(context) {
      try {
        this.klasifikasiOptions = await service.getKlasifikasiOptions(context)
      } catch (e) {
        console.log(e)
        this.klasifikasiOptions = []
      }
    },
  },
})

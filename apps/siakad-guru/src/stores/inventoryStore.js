import { defineStore } from 'pinia'
import { inventoryService } from '../boot/services.js'

export const useInventoryStore = defineStore('inventory', {
  state: () => ({
    items: [],
    selectedItem: null,
    loading: false,
    error: null,
  }),

  getters: {
    hasItems: (state) => state.items.length > 0,
    hasSelectedItem: (state) => !!state.selectedItem,
  },

  actions: {
    async loadItems({ classId }) {
      this.loading = true
      this.error = null
      try {
        this.items = await inventoryService.loadItemsByClass({ classId })
      } catch (err) {
        this.error = err.message
      } finally {
        this.loading = false
      }
    },

    async loadItemById({ itemId }) {
      this.loading = true
      this.error = null
      try {
        this.selectedItem = await inventoryService.loadItemById({ itemId })
      } catch (err) {
        this.error = err.message
      } finally {
        this.loading = false
      }
    },

    async submitDamageReport({
      itemId,
      itemName,
      condition,
      description,
      reportMethod,
      date,
      teacherId,
      teacherName,
      classId,
      className,
      schoolId,
      academicYearId,
      semesterId,
    }) {
      this.loading = true
      this.error = null
      try {
        const report = await inventoryService.submitDamageReport({
          itemId,
          itemName,
          condition,
          description,
          reportMethod,
          date,
          teacherId,
          teacherName,
          classId,
          className,
          schoolId,
          academicYearId,
          semesterId,
        })
        return report
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.loading = false
      }
    },
  },
})

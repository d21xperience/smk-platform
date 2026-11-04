// apps/siakad-tu/src/stores/student/studentDetailStore.js

import { defineStore } from 'pinia'
import { useStudentQueryService } from '@/services/student/serviceFactory.js'
import { useContextStore } from '@/stores/contextStore.js'

export const useStudentDetailStore = defineStore('studentDetail', {
  state: () => ({
    cache: new Map(),
    loadingIds: new Set(),
    errorMap: new Map(),
  }),

  getters: {
    getById: (state) => (studentId) => {
      return state.cache.get(studentId) || null
    },
    isLoading: (state) => (studentId) => {
      return state.loadingIds.has(studentId)
    },
    getError: (state) => (studentId) => {
      return state.errorMap.get(studentId) || null
    },
    cacheSize: (state) => state.cache.size,
  },

  actions: {
    setItem(studentId, data) {
      this.cache = new Map(this.cache)
      this.cache.set(studentId, data)
      this.loadingIds = new Set(this.loadingIds)
      this.loadingIds.delete(studentId)
      this.errorMap = new Map(this.errorMap)
      this.errorMap.delete(studentId)
    },

    setLoading(studentId, loading) {
      this.loadingIds = new Set(this.loadingIds)
      if (loading) {
        this.loadingIds.add(studentId)
      } else {
        this.loadingIds.delete(studentId)
      }
    },

    setError(studentId, error) {
      this.errorMap = new Map(this.errorMap)
      this.errorMap.set(studentId, error)
      this.loadingIds = new Set(this.loadingIds)
      this.loadingIds.delete(studentId)
    },

    /**
     * Invalidasi cache untuk 1 siswa (misal: setelah selesai Edit)
     */
    invalidate(studentId) {
      this.cache = new Map(this.cache)
      this.cache.delete(studentId)
    },

    /**
     * Invalidasi SELURUH cache.
     * Digunakan saat Tahun Ajaran/Semester berubah, atau menerima event dari NAT.
     */
    reset() {
      this.cache = new Map()
      this.loadingIds = new Set()
      this.errorMap = new Map()
      console.log('[StudentDetailStore] Cache dibersihkan (Invalidated).')
    },

    async fetchStudentById(studentId) {
      // CACHE-FIRST STRATEGY: Jika sudah ada di cache, langsung return.
      if (this.cache.has(studentId)) {
        console.log(`[StudentDetailStore] Mengambil data ${studentId} dari Cache.`)
        return
      }

      console.log(`[StudentDetailStore] Memanggil API untuk ${studentId}...`)
      this.setLoading(studentId, true)
      
      try {
        const queryService = useStudentQueryService()
        const contextStore = useContextStore()
        const context = contextStore.currentContext

        const result = await queryService.getStudentById(studentId, context)
        
        if (result.success) {
          this.setItem(studentId, result.data)
        } else {
          this.setError(studentId, result.error)
        }
      } catch (err) {
        this.setError(studentId, { message: err.message })
      }
    }
  },
})

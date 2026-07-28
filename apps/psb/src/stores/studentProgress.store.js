// src/stores/studentProgress.store.js
import { defineStore } from 'pinia'
import { getEventDispatcher } from '@/boot/events'
import { createInitialStudentProgressCache } from '@/cache/StudentProgressCache'
import { studentProgressProjection } from '@/projection/StudentProgressProjection'

export const useStudentProgressStore = defineStore('studentProgress', {
  state: () => ({
    cache: createInitialStudentProgressCache(),
  }),
  getters: {
    getProgress: (state) => (studentId) => state.cache.students[studentId] || null,
    getAllProgress: (state) => state.cache.students,
    getByClass: (state) => (className) => {
      return Object.values(state.cache.students).filter((s) => s.className === className)
    },
  },
  actions: {
    init() {
      const dispatcher = getEventDispatcher()
      dispatcher.on('*', (event) => {
        this.cache = studentProgressProjection(this.cache, event)
      })
    },
  },
})

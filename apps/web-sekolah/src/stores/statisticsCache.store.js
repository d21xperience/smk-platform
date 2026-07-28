// @/stores/statisticsCache.store.js
import { defineStore } from 'pinia'
import { getEventDispatcher } from '@/boot/events'
import { createInitialCache } from '@/cache/StatisticsCache'
import { classAttendanceProjection } from '@/projection/ClassAttendanceProjection'
import { teacherDashboardProjection } from '@/projection/TeacherDashboardProjection'
import { studentSummaryProjection } from '@/projection/StudentSummaryProjection'
import { homeroomSummaryProjection } from '@/projection/HomeroomSummaryProjection'

export const useStatisticsCacheStore = defineStore('statisticsCache', {
  state: () => ({
    cache: createInitialCache(),
  }),
  getters: {
    classAttendance: (state) => state.cache.classAttendance,
    teacherDashboard: (state) => (teacherId) => state.cache.teacherDashboard[teacherId] || null,
    studentSummary: (state) => (studentId) => state.cache.studentSummary[studentId] || null,
    homeroomSummary: (state) => (className) => state.cache.homeroomSummary[className] || null,
    globalStats: (state) => state.cache.globalStats,
  },
  actions: {
    init() {
      const dispatcher = getEventDispatcher()
      // Daftarkan semua proyeksi
      dispatcher.on('*', (event) => {
        this.applyProjections(event)
      })
    },
    applyProjections(event) {
      let newState = this.cache
      newState = classAttendanceProjection(newState, event)
      newState = teacherDashboardProjection(newState, event)
      newState = studentSummaryProjection(newState, event)
      newState = homeroomSummaryProjection(newState, event)
      this.cache = newState
    },
  },
})

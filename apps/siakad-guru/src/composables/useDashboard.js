// src/composables/useDashboard.js
import { computed } from 'vue'
import { useStatisticsCacheStore } from '@/stores/statisticsCache.store'

// src/composables/useDashboard.js (update)
export function useDashboard() {
  const cacheStore = useStatisticsCacheStore()

  return {
    classAttendance: computed(() => cacheStore.classAttendance),
    getTeacherDashboard: (teacherId) => computed(() => cacheStore.teacherDashboard(teacherId)),
    getStudentSummary: (studentId) => computed(() => cacheStore.studentSummary(studentId)),
    globalStats: computed(() => cacheStore.globalStats),
    assessment: computed(() => cacheStore.cache.assessment),
    getClassAssessment: (className, subject) =>
      computed(() => cacheStore.cache.assessment[`${className}_${subject}`] || null),
  }
}

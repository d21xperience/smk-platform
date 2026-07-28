import { computed } from 'vue'
import { useAttendanceStore } from '@/stores/attendance.store'

export function useAttendance() {
  const store = useAttendanceStore()

  return {
    // State (readonly)
    currentRecord: computed(() => store.currentRecord),
    students: computed(() => store.students),
    lastAutoSave: computed(() => store.lastAutoSave),
    loading: computed(() => store.loading),
    saving: computed(() => store.saving),
    error: computed(() => store.error),

    // Getter
    canEdit: computed(() => store.canEdit),
    isAllMarked: computed(() => store.isAllMarked),
    summary: computed(() => store.summary),
    submitValidation: computed(() => store.submitValidation),
    recordStatus: computed(() => store.recordStatus),

    // Actions
    loadOrCreateAttendance: store.loadOrCreateAttendance,
    updateStudentStatus: store.updateStudentStatus,
    markAllAs: store.markAllAs,
    updateStudentNote: store.updateStudentNote,
    saveDraftManually: store.saveDraftManually,
    submitAttendance: store.submitAttendance,
    loadSummary: store.loadSummary,
    reset: store.reset,
  }
}

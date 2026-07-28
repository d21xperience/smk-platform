import { computed } from 'vue'
import { useAssessmentStore } from '@/stores/assessment.store'

export function useAssessment() {
  const store = useAssessmentStore()

  return {
    // State
    session: computed(() => store.session),
    components: computed(() => store.components),
    selectedComponent: computed(() => store.selectedComponent),
    finalResult: computed(() => store.finalResult),
    loading: computed(() => store.loading),
    error: computed(() => store.error),

    // Getter
    isAllComponentsScored: computed(() => store.isAllComponentsScored),
    previewGrades: computed(() => store.previewGrades),
    classStats: computed(() => store.classStats),
    canEdit: computed(() => store.canEdit),
    isSubmitted: computed(() => store.isSubmitted),

    // Actions
    loadOrCreateSession: store.loadOrCreateSession,
    selectComponent: store.selectComponent,
    addComponent: store.addComponent,
    removeComponent: store.removeComponent,
    updateScore: store.updateScore,
    updateScores: store.updateScores,
    saveDraft: store.saveDraft,
    submitAssessment: store.submitAssessment,
    reset: store.reset,
  }
}

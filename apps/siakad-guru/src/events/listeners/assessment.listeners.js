import { EVENTS } from '../EventCatalog'
// import { useStudentProgressStore } from 'src/stores/studentProgress.store'; // nanti di Sprint H

export function registerAssessmentListeners(eventBus) {
  eventBus.on(EVENTS.ASSESSMENT_FINALIZED, (payload) => {
    console.log('[AssessmentListener] Assessment finalized:', payload.id)
    // Update student progress projection (akan diimplementasikan di Sprint H)
    // const progressStore = useStudentProgressStore();
    // progressStore.updateFromAssessment(payload);
  })
}

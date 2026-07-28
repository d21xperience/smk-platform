import { computed } from 'vue'
import { useTeachingStore } from '@/stores/teaching.store'

export function useTeaching() {
  const store = useTeachingStore()

  return {
    // --- State (readonly) ---
    schedule: computed(() => store.schedule),
    sessions: computed(() => store.sessions),
    activeSession: computed(() => store.activeSession),
    activeSessionId: computed(() => store.activeSessionId),
    hasOngoingSession: computed(() => store.hasOngoingSession),
    completedSessions: computed(() => store.completedSessions),
    currentJournal: computed(() => store.currentJournal),
    loading: computed(() => store.loading),
    error: computed(() => store.error),

    // --- Actions ---
    loadSchedule: store.loadSchedule,
    loadSessions: store.loadSessions,
    startSession: store.startSession,
    completeSession: store.completeSession,
    loadJournal: store.loadJournal,
    saveJournal: store.saveJournal,
    reset: store.reset,

    // --- Getter dengan parameter (tetap gunakan function) ---
    scheduleByDay: store.scheduleByDay,
  }
}

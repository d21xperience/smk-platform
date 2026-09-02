// FILE: src/composables/useJournal.js
// STATUS: MODIFY
// STATUS IMPLEMENTASI: COMPLETE

import { computed } from 'vue'
import { useJournalStore } from '../stores/journalStore.js'
import { useTeachingStore } from '../stores/teachingStore.js'
import { useRouter } from 'vue-router'
import { useContext } from './useContext.js'

export function useJournal() {
  const journalStore = useJournalStore()
  const teachingStore = useTeachingStore()
  const router = useRouter()

  const { currentContext, isHistoryMode } = useContext()

  const journal = computed(() => journalStore.journal)
  const isLoading = computed(() => journalStore.loading)
  const isSaving = computed(() => journalStore.saving)
  const error = computed(() => journalStore.error)
  const isDraft = computed(() => journalStore.isDraft)
  const isSubmitted = computed(() => journalStore.isSubmitted)

  const loadJournalOnly = async (teachingSessionId) => {
    if (!currentContext.value) {
      return null
    }

    await journalStore.loadJournal({ teachingSessionId })

    return journalStore.journal
  }

  const loadOrCreateJournal = async (teachingSessionId) => {
    if (!currentContext.value) {
      return null
    }

    await journalStore.loadJournal({ teachingSessionId })

    if (journalStore.journal) {
      return journalStore.journal
    }

    if (isHistoryMode.value) {
      return null
    }

    const teachingSession = teachingStore.sessions.find((s) => s.id === teachingSessionId)

    if (!teachingSession) {
      throw new Error('Teaching session not found')
    }

    await journalStore.createJournal({
      teachingSessionId,
      classId: teachingSession.classId,
      className: teachingSession.className,
      subjectId: teachingSession.subjectId,
      subjectName: teachingSession.subjectName,
      date: teachingSession.date,
      schoolId: currentContext.value.schoolId,
      academicYearId: currentContext.value.academicYearId,
      semesterId: currentContext.value.semesterId,
    })

    return journalStore.journal
  }

  const updateJournal = async ({ material, activities, reflection }) => {
    if (isHistoryMode.value) return

    if (!journalStore.journal || journalStore.isSubmitted) return

    await journalStore.updateJournal({
      journalId: journalStore.journal.id,
      material,
      activities,
      reflection,
    })
  }

  const submitJournal = async () => {
    if (isHistoryMode.value) return

    if (!journalStore.journal) return

    await journalStore.submitJournal({
      journalId: journalStore.journal.id,
    })
  }

  const navigateToSession = (teachingSessionId) => {
    router.push({
      name: 'journal-session',
      params: { teachingSessionId },
    })
  }

  return {
    journal,
    isLoading,
    isSaving,
    error,
    isDraft,
    isSubmitted,
    loadJournalOnly,
    loadOrCreateJournal,
    updateJournal,
    submitJournal,
    navigateToSession,
  }
}

// FILE: src/composables/useAttendance.js
// STATUS: MODIFY
// STATUS IMPLEMENTASI: COMPLETE

import { computed } from 'vue'
import { useAttendanceStore } from '../stores/attendanceStore.js'
import { useTeachingStore } from '../stores/teachingStore.js'
import { useRouter } from 'vue-router'
import { useContext } from './useContext.js'

export function useAttendance() {
  const attendanceStore = useAttendanceStore()
  const teachingStore = useTeachingStore()
  const router = useRouter()

  const {
    currentContext,
    isHistoryMode,
  } = useContext()

  const session = computed(() => attendanceStore.session)
  const students = computed(() => attendanceStore.students)
  const isLoading = computed(() => attendanceStore.loading)
  const isSaving = computed(() => attendanceStore.saving)
  const error = computed(() => attendanceStore.error)
  const isDraft = computed(() => attendanceStore.isDraft)
  const isSubmitted = computed(() => attendanceStore.isSubmitted)
  const summary = computed(() => attendanceStore.getSummary)

  const loadSession = async (teachingSessionId) => {
    if (!currentContext.value) {
      return null
    }

    await attendanceStore.loadSession({ teachingSessionId })

    return attendanceStore.session
  }

  const loadOrCreateSession = async (teachingSessionId) => {
    if (!currentContext.value) {
      return null
    }

    await attendanceStore.loadSession({ teachingSessionId })

    if (attendanceStore.session) {
      return attendanceStore.session
    }

    if (isHistoryMode.value) {
      return null
    }

    const teachingSession = teachingStore.sessions.find((s) => s.id === teachingSessionId)

    if (!teachingSession) {
      throw new Error('Teaching session not found')
    }

    await attendanceStore.createSession({
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

    return attendanceStore.session
  }

  const saveDraft = async () => {
    if (isHistoryMode.value) return

    if (!attendanceStore.session || attendanceStore.isSubmitted) return

    await attendanceStore.saveDraft({
      sessionId: attendanceStore.session.id,
      records: attendanceStore.students.map((r) => ({
        studentId: r.studentId,
        status: r.status,
        note: r.note,
      })),
    })
  }

  const updateRecordStatus = async (studentId, status) => {
    if (isHistoryMode.value) return

    attendanceStore.updateRecordStatus(studentId, status)
    await saveDraft()
  }

  const updateRecordNote = (studentId, note) => {
    if (isHistoryMode.value) return

    attendanceStore.updateRecordNote(studentId, note)
  }

  const submitAttendance = async () => {
    if (isHistoryMode.value) return

    if (!attendanceStore.session) return

    await attendanceStore.submitAttendance({
      sessionId: attendanceStore.session.id,
    })
  }

  const navigateToSession = (teachingSessionId) => {
    router.push({
      name: 'attendance-session',
      params: { teachingSessionId },
    })
  }

  return {
    session,
    students,
    isLoading,
    isSaving,
    error,
    isDraft,
    isSubmitted,
    summary,
    loadSession,
    loadOrCreateSession,
    updateRecordStatus,
    updateRecordNote,
    saveDraft,
    submitAttendance,
    navigateToSession,
  }
}

import { computed, ref } from 'vue'
import { useTeachingStore } from '../stores/teachingStore.js'
import { useContextStore } from '../stores/contextStore.js'
import { useRouter } from 'vue-router'

export function useTeaching() {
  const teachingStore = useTeachingStore()
  const contextStore = useContextStore()
  const router = useRouter()

  const selectedDate = ref(new Date().toISOString().split('T')[0])

  const currentContext = computed(() => contextStore.currentContext)
  const schedules = computed(() => teachingStore.schedules)
  const sessions = computed(() => teachingStore.sessions)
  const activeSession = computed(() => teachingStore.activeSession)
  const isLoading = computed(() => teachingStore.loading)
  const error = computed(() => teachingStore.error)

  const loadSchedules = async () => {
    if (!currentContext.value) return
    await teachingStore.loadSchedules({
      schoolId: currentContext.value.schoolId,
      academicYearId: currentContext.value.academicYearId,
      semesterId: currentContext.value.semesterId,
      teacherId: currentContext.value.userId,
    })
  }

  const loadSessions = async (date) => {
    if (!currentContext.value) return
    const targetDate = date || selectedDate.value
    await teachingStore.loadSessions({
      schoolId: currentContext.value.schoolId,
      academicYearId: currentContext.value.academicYearId,
      semesterId: currentContext.value.semesterId,
      teacherId: currentContext.value.userId,
      date: targetDate,
    })
  }

  const createSession = async (scheduleId) => {
    if (!currentContext.value) return
    const session = await teachingStore.createSession({
      scheduleId,
      date: selectedDate.value,
      schoolId: currentContext.value.schoolId,
      academicYearId: currentContext.value.academicYearId,
      semesterId: currentContext.value.semesterId,
      teacherId: currentContext.value.userId,
    })
    return session
  }

  const startSession = async (sessionId, teacherPresence) => {
    if (!currentContext.value) return
    const session = await teachingStore.startSession({
      sessionId,
      teacherPresence,
      schoolId: currentContext.value.schoolId,
      academicYearId: currentContext.value.academicYearId,
      semesterId: currentContext.value.semesterId,
    })
    return session
  }

  const endSession = async (sessionId, notes) => {
    if (!currentContext.value) return
    const session = await teachingStore.endSession({
      sessionId,
      notes,
      schoolId: currentContext.value.schoolId,
      academicYearId: currentContext.value.academicYearId,
      semesterId: currentContext.value.semesterId,
    })
    return session
  }

  const cancelSession = async (sessionId, reason) => {
    if (!currentContext.value) return
    const session = await teachingStore.cancelSession({
      sessionId,
      reason,
      schoolId: currentContext.value.schoolId,
      academicYearId: currentContext.value.academicYearId,
      semesterId: currentContext.value.semesterId,
    })
    return session
  }

  const navigateToSession = (sessionId) => {
    router.push({ name: 'teaching-session', params: { sessionId } })
  }

  // Filter schedules that match the selected date's day of week
  const getAvailableSchedules = computed(() => {
    if (!schedules.value.length) return []

    const dateObj = new Date(selectedDate.value)
    const dayOfWeek = dateObj.getDay() === 0 ? 7 : dateObj.getDay() // Convert: 0=Sunday -> 7
    // Get scheduleIds that already have sessions on this date
    const existingScheduleIds = new Set(sessions.value.map((s) => s.scheduleId))
    return schedules.value.filter(
      (schedule) => schedule.dayOfWeek === dayOfWeek && !existingScheduleIds.has(schedule.id),
    )
  })

  return {
    selectedDate,
    schedules,
    sessions,
    activeSession,
    isLoading,
    error,
    getAvailableSchedules,
    loadSchedules,
    loadSessions,
    createSession,
    startSession,
    endSession,
    cancelSession,
    navigateToSession,
  }
}

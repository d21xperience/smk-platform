import { useAttendanceStore } from './attendanceStore.js'
import { computed } from 'vue'

export const useAttendanceSelector = () => {
  const store = useAttendanceStore()

  const totalStudents = computed(() => store.roster.length)
  const sickCount = computed(
    () => store.roster.filter((r) => r.attendance.status === 'SICK').length,
  )
  const permitCount = computed(
    () => store.roster.filter((r) => r.attendance.status === 'PERMIT').length,
  )
  const presentCount = computed(
    () => store.roster.filter((r) => r.attendance.status === 'PRESENT').length,
  )
  const alphaCount = computed(
    () => store.roster.filter((r) => r.attendance.status === 'ALPHA').length,
  )

  return { totalStudents, presentCount, alphaCount, sickCount, permitCount }
}

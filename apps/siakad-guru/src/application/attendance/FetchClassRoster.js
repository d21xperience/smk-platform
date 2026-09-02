// application/attendance/FetchClassRoster.js
import { classMock } from '@/adapters/mock/classMock.js'
import { useContextStore } from '@/stores/context.store.js'

export const FetchClassRosterUseCase = {
  execute: async (classId) => {
    const contextStore = useContextStore()
    const { semesterId, academicYearId } = contextStore

    // Panggil mock untuk mendapatkan siswa + status absensi
    const roster = await classMock.getStudentsByClassId(classId, semesterId, academicYearId)

    // Urutkan berdasarkan nama (opsional)
    return roster.sort((a, b) => a.student.name.localeCompare(b.student.name))
  },
}

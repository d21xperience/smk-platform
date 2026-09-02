// stores/attendance.action.js
import { useAttendanceStore } from './attendanceStore.js'
import { attendanceService } from '../services/AttendanceService.js'

export const useAttendanceAction = () => {
  const store = useAttendanceStore()

  const loadClasses = async () => {
    console.log('[Action] loadClasses called')
    store.setLoading(true)
    store.setError(null)
    try {
      console.log('[Action] Calling attendanceService.fetchClasses...')
      const data = await attendanceService.fetchClasses()
      console.log('[Action] Data returned from service:', data)
      store.setClasses(data)
      console.log('[Action] Store updated. classes.value =', store.classes.value)
    } catch (err) {
      console.error('[Action] Error in loadClasses:', err)
      store.setError(err.message || 'Gagal memuat kelas')
    } finally {
      store.setLoading(false)
    }
  }

  const loadRoster = async (classId) => {
    // Validasi input
    if (!classId) {
      store.setError('ID Kelas tidak valid')
      return
    }

    store.setLoading(true)
    store.setError(null)
    try {
      const data = await attendanceService.fetchRoster(classId)
      store.setRoster(data)
    } catch (err) {
      store.setError(err.message || 'Gagal memuat daftar siswa')
    } finally {
      store.setLoading(false)
    }
  }

  const submitAttendance = async (payload) => {
    // Validasi payload
    if (!payload.studentId || !payload.classId) {
      throw new Error('Data siswa atau kelas tidak lengkap')
    }

    store.setLoading(true)
    store.setError(null)
    try {
      const result = await attendanceService.submit(payload)
      // Update lokal state setelah sukses
      store.updateStudentStatus(payload.studentId, result.status)
      return result
    } catch (err) {
      store.setError(err.message || 'Gagal menyimpan absensi')
      throw err
    } finally {
      store.setLoading(false)
    }
  }

  return { loadClasses, loadRoster, submitAttendance }
}

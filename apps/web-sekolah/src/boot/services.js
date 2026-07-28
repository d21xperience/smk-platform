import { useAuthStore } from '@/stores/auth.store'
import { useContextStore } from '@/stores/context.store'
import { useTeachingStore } from '@/stores/teaching.store' // tambahan
import { useAttendanceStore } from '@/stores/attendance.store'
import { useAssessmentStore } from '@/stores/assessment.store'
// import { useReportingStore } from '@/stores/reporting.store'

import { AuthService } from '@/services/AuthService'
import { ContextService } from '@/services/ContextService'
import { TeachingService } from '@/services/TeachingService' // tambahan
import { AttendanceService } from '@/services/AttendanceService'
import { AssessmentService } from '@/services/AssessmentService'

// Import adapter mock
import { authMockAdapter } from '@/adapters/mock/auth.mock'
import { academicMockAdapter } from '@/adapters/mock/academic.mock'
import { teachingMockAdapter } from '@/adapters/mock/teaching.mock' // tambahan
import { attendanceMockAdapter } from '@/adapters/mock/attendance.mock'
import { assessmentMockAdapter } from '@/adapters/mock/assessment.mock'

// Import adapter API
import { authApiAdapter } from '@/adapters/api/auth.api'
import { academicApiAdapter } from '@/adapters/api/academic.api'
import { teachingApiAdapter } from '@/adapters/api/teaching.api' // tambahan
import { attendanceApiAdapter } from '@/adapters/api/attendance.api'
import { assessmentApiAdapter } from '@/adapters/api/assessment.api'

// eslint-disable-next-line no-unused-vars
export default async ({ router }) => {
  // Tentukan adapter berdasarkan env variable (default: pakai mock saat dev)
  const useMock = import.meta.env.QCLI_MOCK_MODE === true //|| process.env.DEV
  if (useMock) {
    console.warn('🚫 MOCK MODE ACTIVE - Menggunakan data dummy')
  } else {
    console.log('✔ MOCK MODE INACTIVE - Koneksi ke backend asli')
  }
  // Pilih adapter
  const authAdapter = useMock ? authMockAdapter : authApiAdapter
  const academicAdapter = useMock ? academicMockAdapter : academicApiAdapter
  const teachingAdapter = useMock ? teachingMockAdapter : teachingApiAdapter // tambahan
  const attendanceAdapter = useMock ? attendanceMockAdapter : attendanceApiAdapter
  const assessmentAdapter = useMock ? assessmentMockAdapter : assessmentApiAdapter

  // Buat service
  const authService = new AuthService(authAdapter)
  const contextService = new ContextService(academicAdapter)
  const teachingService = new TeachingService(teachingAdapter) // tambahan
  const attendanceService = new AttendanceService(attendanceAdapter)
  const assessmentService = new AssessmentService(assessmentAdapter)

  // Inject ke store
  // const authStore = useAuthStore()
  // const contextStore = useContextStore()
  // const teachingStore = useTeachingStore() // tambahan

  // authStore.setService(authService)
  // contextStore.setService(contextService)
  // teachingStore.setService(teachingService) // tambahan
  // Restore session jika token masih ada
  useAuthStore().setService(authService)
  useContextStore().setService(contextService)
  useTeachingStore().setService(teachingService)
  useAttendanceStore().setService(attendanceService)
  useAssessmentStore().setService(assessmentService)

  if (useAuthStore().token) {
    try {
      await useAuthStore().checkAuth()
    } catch {
      // Token invalid, sudah di-handle di store
    }
  } else {
    useAuthStore().initialized = true
  }

  // Opsional: preload data akademik untuk dropdown (bisa juga dilakukan di halaman)
  // contextStore.loadAcademicYears()
}

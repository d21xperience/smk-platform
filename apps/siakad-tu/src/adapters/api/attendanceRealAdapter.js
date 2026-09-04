// apps/siakad-tu/src/adapters/api/attendanceRealAdapter.js

import { api } from '@/boot/axios.js'

/**
 * Helper: Map Backend Attendance Session DTO ke Frontend Contract
 */
const mapBackendToSession = (backendDto) => {
  if (!backendDto) return null

  return {
    sessionId: backendDto.session_id || backendDto.id,
    classId: backendDto.class_id,
    subjectId: backendDto.subject_id,
    periodId: backendDto.period_id,
    date: backendDto.date,
    status: backendDto.status, // 'OPEN', 'CLOSED'
    createdAt: backendDto.created_at || backendDto.createdAt,
    updatedAt: backendDto.updated_at || backendDto.updatedAt,
  }
}

/**
 * Helper: Map Backend Attendance Record DTO ke Frontend Contract
 */
const mapBackendToAttendanceRecord = (backendDto) => {
  if (!backendDto) return null

  return {
    recordId: backendDto.record_id || backendDto.id,
    sessionId: backendDto.session_id,
    studentId: backendDto.student_id,
    status: backendDto.status, // 'PRESENT', 'ABSENT', 'SICK', 'PERMISSION', 'LATE'
    note: backendDto.note,
    recordedAt: backendDto.recorded_at || backendDto.recordedAt,
    recordedBy: backendDto.recorded_by || backendDto.recordedBy,
  }
}

/**
 * Helper: Map Backend Attendance Summary DTO ke Frontend Contract
 */
const mapBackendToAttendanceSummary = (backendDto) => {
  if (!backendDto) return null

  return {
    studentId: backendDto.student_id,
    periodId: backendDto.period_id,
    totalSessions: backendDto.total_sessions || 0,
    present: backendDto.present || 0,
    absent: backendDto.absent || 0,
    sick: backendDto.sick || 0,
    permission: backendDto.permission || 0,
    late: backendDto.late || 0,
    attendanceRate: backendDto.attendance_rate || 0, // Persentase kehadiran
  }
}

/**
 * Helper: Konversi ApplicationError ke format Service yang konsisten
 */
const handleApiError = (error) => ({
  success: false,
  data: null,
  error: {
    code: error.code || 'API_ERROR',
    message: error.message || 'Terjadi kesalahan pada server.',
    details: error.details || null,
  },
})

const handleSuccess = (data) => ({
  success: true,
  data,
  error: null,
})

export const attendanceRealAdapter = {
  /**
   * Buka sesi absensi baru
   */
  async openSession(command) {
    try {
      const payload = {
        classId: command.classId,
        subjectId: command.subjectId,
        date: command.date,
      }
      const response = await api.post('/attendance/sessions', payload)
      return handleSuccess(mapBackendToSession(response.data))
    } catch (error) {
      return handleApiError(error)
    }
  },

  /**
   * Tutup sesi absensi
   */
  async closeSession(sessionId) {
    try {
      const response = await api.post(`/attendance/sessions/${sessionId}/close`)
      return handleSuccess(mapBackendToSession(response.data))
    } catch (error) {
      return handleApiError(error)
    }
  },

  /**
   * Catat absensi bulk (banyak siswa sekaligus)
   */
  async submitBulkAttendance(sessionId, command) {
    try {
      const payload = {
        records: command.records.map((r) => ({
          studentId: r.studentId,
          status: r.status,
          note: r.note,
        })),
      }
      const response = await api.post(`/attendance/sessions/${sessionId}/records`, payload)
      return handleSuccess(response.data.records.map(mapBackendToAttendanceRecord))
    } catch (error) {
      return handleApiError(error)
    }
  },

  /**
   * Update absensi individual (koreksi)
   */
  async updateAttendanceRecord(recordId, command) {
    try {
      const payload = {
        status: command.status,
        note: command.note,
      }
      const response = await api.put(`/attendance/records/${recordId}`, payload)
      return handleSuccess(mapBackendToAttendanceRecord(response.data))
    } catch (error) {
      return handleApiError(error)
    }
  },

  /**
   * Ambil detail sesi absensi by ID
   */
  async getSessionById(sessionId) {
    try {
      const response = await api.get(`/attendance/sessions/${sessionId}`)
      return handleSuccess(mapBackendToSession(response.data))
    } catch (error) {
      return handleApiError(error)
    }
  },

  /**
   * Ambil daftar sesi absensi dengan filter
   */
  async getSessions(context, filters = {}) {
    try {
      const params = { schoolId: context.schoolId, periodId: context.periodId, ...filters }
      const response = await api.get('/attendance/sessions', { params })

      const mappedData = response.data.items
        ? { ...response.data, items: response.data.items.map(mapBackendToSession) }
        : response.data

      return handleSuccess(mappedData)
    } catch (error) {
      return handleApiError(error)
    }
  },

  /**
   * Ambil daftar record absensi dalam sesi tertentu
   */
  async getAttendanceRecords(sessionId) {
    try {
      const response = await api.get(`/attendance/sessions/${sessionId}/records`)
      return handleSuccess(response.data.map(mapBackendToAttendanceRecord))
    } catch (error) {
      return handleApiError(error)
    }
  },

  /**
   * Ambil ringkasan kehadiran siswa
   */
  async getAttendanceSummary(studentId, context) {
    try {
      const response = await api.get(`/students/${studentId}/attendance-summary`, {
        params: { periodId: context.periodId },
      })
      return handleSuccess(mapBackendToAttendanceSummary(response.data))
    } catch (error) {
      return handleApiError(error)
    }
  },

  /**
   * Ambil rekap absensi per kelas
   */
  async getClassAttendanceSummary(classId, context) {
    try {
      const response = await api.get(`/classes/${classId}/attendance-summary`, {
        params: { periodId: context.periodId },
      })
      return handleSuccess(response.data.map(mapBackendToAttendanceSummary))
    } catch (error) {
      return handleApiError(error)
    }
  },
}

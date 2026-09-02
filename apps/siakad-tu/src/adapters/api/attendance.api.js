import { api as apiClient } from '@/boot/axios'
import { AttendanceRecord } from '@/models/Attendance'

export const attendanceApiAdapter = {
  // eslint-disable-next-line no-unused-vars
  async fetchStudentsForSession(sessionId, className) {
    const { data } = await apiClient.get(`/teaching/sessions/${sessionId}/students`)
    return data.map((s) => ({
      studentId: s.id,
      studentName: s.name,
      nis: s.nis,
      status: s.status || 'unmarked',
      note: s.note || '',
    }))
  },

  async fetchAttendanceRecord(sessionId) {
    const { data } = await apiClient.get(`/attendance/sessions/${sessionId}`)
    if (!data) return null
    return new AttendanceRecord(data)
  },

  async saveDraft(session, students) {
    const { data } = await apiClient.post(`/attendance/sessions/${session.id}/draft`, {
      students,
    })
    return new AttendanceRecord(data)
  },

  async submitAttendance(sessionId) {
    const { data } = await apiClient.post(`/attendance/sessions/${sessionId}/submit`)
    return new AttendanceRecord(data)
  },

  async getAttendanceSummary(sessionId) {
    const { data } = await apiClient.get(`/attendance/sessions/${sessionId}/summary`)
    if (!data) return null
    return {
      sessionId: data.sessionId,
      date: data.date,
      className: data.className,
      subject: data.subject,
      students: data.students,
      status: data.status,
      lastSavedAt: data.lastSavedAt,
    }
  },
}

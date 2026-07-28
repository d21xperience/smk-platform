// src/projection/StudentProgressProjection.js
import { DomainEventType } from '@/events/DomainEventType'

export function studentProgressProjection(state, event) {
  switch (event.type) {
    case DomainEventType.ATTENDANCE_SUBMITTED: {
      const { students, className } = event.payload
      if (!students) return state
      const newStudents = { ...state.students }
      students.forEach((s) => {
        const existing = newStudents[s.studentId] || {
          studentId: s.studentId,
          studentName: s.studentName || `Siswa ${s.studentId}`,
          className: className || '',
          attendanceSummary: null,
          assessmentSummaries: [],
          teacherNotes: [],
          homeroomNotes: [],
          achievements: [],
          violations: [],
          counselingRecords: [],
          lastUpdated: event.timestamp,
        }
        // Update attendance summary
        existing.attendanceSummary = {
          totalSessions: (existing.attendanceSummary?.totalSessions || 0) + 1,
          hadir: (existing.attendanceSummary?.hadir || 0) + (s.status === 'hadir' ? 1 : 0),
          izin: (existing.attendanceSummary?.izin || 0) + (s.status === 'izin' ? 1 : 0),
          sakit: (existing.attendanceSummary?.sakit || 0) + (s.status === 'sakit' ? 1 : 0),
          alpha: (existing.attendanceSummary?.alpha || 0) + (s.status === 'alpha' ? 1 : 0),
        }
        existing.attendanceSummary.persentaseKehadiran =
          existing.attendanceSummary.totalSessions > 0
            ? Math.round(
                (existing.attendanceSummary.hadir / existing.attendanceSummary.totalSessions) * 100,
              )
            : 0
        existing.lastUpdated = event.timestamp
        newStudents[s.studentId] = existing
      })
      return { ...state, students: newStudents }
    }

    case DomainEventType.ASSESSMENT_SUBMITTED: {
      const { grades } = event.payload
      if (!grades) return state
      const newStudents = { ...state.students }
      grades.forEach((grade) => {
        const existing = newStudents[grade.studentId]
        if (existing) {
          const existingAssessments = existing.assessmentSummaries || []
          const idx = existingAssessments.findIndex((a) => a.subject === event.payload.subject)
          const summary = {
            subject: event.payload.subject,
            finalScore: grade.finalScore,
            predicate: grade.predicate,
            componentDetails: grade.componentDetails || [],
          }
          if (idx >= 0) {
            existingAssessments[idx] = summary
          } else {
            existingAssessments.push(summary)
          }
          existing.assessmentSummaries = existingAssessments
          existing.lastUpdated = event.timestamp
          newStudents[grade.studentId] = existing
        }
      })
      return { ...state, students: newStudents }
    }

    default:
      return state
  }
}

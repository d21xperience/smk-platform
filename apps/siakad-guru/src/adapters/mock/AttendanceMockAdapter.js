import { SESSION_STATUS, ATTENDANCE_STATUS } from '../../domain/attendance/models/AttendanceStatus.js'

// Mock students data
const mockStudents = {
  'CLS-X-A': [
    { studentId: 'STU-001', studentName: 'Ahmad Fauzi', classId: 'CLS-X-A', className: 'X-A' },
    { studentId: 'STU-002', studentName: 'Budi Hartono', classId: 'CLS-X-A', className: 'X-A' },
    { studentId: 'STU-003', studentName: 'Citra Dewi', classId: 'CLS-X-A', className: 'X-A' },
    { studentId: 'STU-004', studentName: 'Diana Putri', classId: 'CLS-X-A', className: 'X-A' },
    { studentId: 'STU-005', studentName: 'Eko Prasetyo', classId: 'CLS-X-A', className: 'X-A' },
    { studentId: 'STU-006', studentName: 'Fitri Handayani', classId: 'CLS-X-A', className: 'X-A' },
    { studentId: 'STU-007', studentName: 'Gunawan Wijaya', classId: 'CLS-X-A', className: 'X-A' },
    { studentId: 'STU-008', studentName: 'Hana Safitri', classId: 'CLS-X-A', className: 'X-A' },
    { studentId: 'STU-009', studentName: 'Irfan Hakim', classId: 'CLS-X-A', className: 'X-A' },
    { studentId: 'STU-010', studentName: 'Joko Susilo', classId: 'CLS-X-A', className: 'X-A' }
  ],
  'CLS-XI-B': [
    { studentId: 'STU-011', studentName: 'Kartika Sari', classId: 'CLS-XI-B', className: 'XI-B' },
    { studentId: 'STU-012', studentName: 'Lukman Hakim', classId: 'CLS-XI-B', className: 'XI-B' },
    { studentId: 'STU-013', studentName: 'Maya Anggraini', classId: 'CLS-XI-B', className: 'XI-B' },
    { studentId: 'STU-014', studentName: 'Nanda Pratama', classId: 'CLS-XI-B', className: 'XI-B' },
    { studentId: 'STU-015', studentName: 'Oscar Firmansyah', classId: 'CLS-XI-B', className: 'XI-B' }
  ]
}

// Mock attendance sessions storage (keyed by teachingSessionId)
let mockAttendanceSessions = {}
let sessionCounter = 1

export class AttendanceMockAdapter {
  async fetchStudentsByClass({ classId }) {
    await new Promise(resolve => setTimeout(resolve, 500))
    return mockStudents[classId] || []
  }

  async loadAttendanceSession({ teachingSessionId }) {
    await new Promise(resolve => setTimeout(resolve, 400))
    return mockAttendanceSessions[teachingSessionId] || null
  }

  async createAttendanceSession({
    teachingSessionId,
    classId,
    className,
    subjectId,
    subjectName,
    date,
    students,
    schoolId,
    academicYearId,
    semesterId
  }) {
    await new Promise(resolve => setTimeout(resolve, 400))

    if (mockAttendanceSessions[teachingSessionId]) {
      throw new Error('Attendance session already exists for this teaching session')
    }

    const records = students.map(student => ({
      studentId: student.studentId,
      studentName: student.studentName,
      status: ATTENDANCE_STATUS.PRESENT,
      note: ''
    }))

    const newSession = {
      id: `ATT-${String(sessionCounter++).padStart(3, '0')}`,
      teachingSessionId,
      classId,
      className,
      subjectId,
      subjectName,
      date,
      status: SESSION_STATUS.DRAFT,
      records,
      academicYearId,
      semesterId,
      schoolId
    }

    mockAttendanceSessions[teachingSessionId] = newSession
    return newSession
  }

  async saveAttendanceDraft({ sessionId, records }) {
    await new Promise(resolve => setTimeout(resolve, 300))

    let targetSession = null
    for (const key in mockAttendanceSessions) {
      if (mockAttendanceSessions[key].id === sessionId) {
        targetSession = mockAttendanceSessions[key]
        break
      }
    }

    if (!targetSession) {
      throw new Error(`Attendance session not found: ${sessionId}`)
    }

    if (targetSession.status !== SESSION_STATUS.DRAFT) {
      throw new Error('Cannot update a submitted attendance session')
    }

    records.forEach(updateRecord => {
      const existingRecord = targetSession.records.find(
        r => r.studentId === updateRecord.studentId
      )
      if (existingRecord) {
        existingRecord.status = updateRecord.status
        if (updateRecord.note !== undefined) {
          existingRecord.note = updateRecord.note
        }
      }
    })

    return targetSession
  }

  async submitAttendance({ sessionId }) {
    await new Promise(resolve => setTimeout(resolve, 400))

    let targetSession = null
    for (const key in mockAttendanceSessions) {
      if (mockAttendanceSessions[key].id === sessionId) {
        targetSession = mockAttendanceSessions[key]
        break
      }
    }

    if (!targetSession) {
      throw new Error(`Attendance session not found: ${sessionId}`)
    }

    if (targetSession.status !== SESSION_STATUS.DRAFT) {
      throw new Error('Attendance session has already been submitted')
    }

    if (targetSession.records.length === 0) {
      throw new Error('Cannot submit attendance without records')
    }

    targetSession.status = SESSION_STATUS.SUBMITTED
    return targetSession
  }
}

import { TEACHING_STATUS, TEACHER_PRESENCE } from '../../domain/teaching/models/TeachingStatus.js'

// Mock data store (simulates database)
let mockSessions = [
  {
    id: 'TS-001',
    scheduleId: 'SCH-001',
    classId: 'CLS-X-A',
    className: 'X-A',
    subjectId: 'SUB-MATH',
    subjectName: 'Matematika',
    date: '2026-08-15',
    startTime: '08:00',
    endTime: '09:30',
    status: TEACHING_STATUS.COMPLETED,
    teacherId: 'USR-001',
    teacherName: 'Budi Santoso, S.Pd.',
    teacherPresence: TEACHER_PRESENCE.PRESENT,
    notes: 'Materi Aljabar Linear',
    academicYearId: 'AY-2026',
    semesterId: '20261',
    schoolId: 'SCH-001'
  },
  {
    id: 'TS-002',
    scheduleId: 'SCH-002',
    classId: 'CLS-XI-B',
    className: 'XI-B',
    subjectId: 'SUB-PHYS',
    subjectName: 'Fisika',
    date: '2026-08-15',
    startTime: '10:00',
    endTime: '11:30',
    status: TEACHING_STATUS.SCHEDULED,
    teacherId: 'USR-001',
    teacherName: 'Budi Santoso, S.Pd.',
    teacherPresence: null,
    notes: '',
    academicYearId: 'AY-2026',
    semesterId: '20261',
    schoolId: 'SCH-001'
  }
]

let sessionCounter = 3

export class TeachingMockAdapter {
  async fetchSchedulesByContext({ schoolId, academicYearId, semesterId, teacherId }) {
    await new Promise(resolve => setTimeout(resolve, 500))

    return [
      {
        id: 'SCH-001',
        classId: 'CLS-X-A',
        className: 'X-A',
        subjectId: 'SUB-MATH',
        subjectName: 'Matematika',
        dayOfWeek: 1,
        startTime: '08:00',
        endTime: '09:30',
        teacherId,
        teacherName: 'Budi Santoso, S.Pd.',
        academicYearId,
        semesterId,
        schoolId
      },
      {
        id: 'SCH-002',
        classId: 'CLS-XI-B',
        className: 'XI-B',
        subjectId: 'SUB-PHYS',
        subjectName: 'Fisika',
        dayOfWeek: 1,
        startTime: '10:00',
        endTime: '11:30',
        teacherId,
        teacherName: 'Budi Santoso, S.Pd.',
        academicYearId,
        semesterId,
        schoolId
      },
      {
        id: 'SCH-003',
        classId: 'CLS-XII-C',
        className: 'XII-C',
        subjectId: 'SUB-MATH',
        subjectName: 'Matematika',
        dayOfWeek: 2,
        startTime: '08:00',
        endTime: '09:30',
        teacherId,
        teacherName: 'Budi Santoso, S.Pd.',
        academicYearId,
        semesterId,
        schoolId
      }
    ]
  }

  async fetchSessionsByDate({ schoolId, academicYearId, semesterId, teacherId, date }) {
    await new Promise(resolve => setTimeout(resolve, 500))

    return mockSessions.filter(session =>
      session.schoolId === schoolId &&
      session.academicYearId === academicYearId &&
      session.semesterId === semesterId &&
      session.teacherId === teacherId &&
      session.date === date
    )
  }

  async createSession({ scheduleId, date, schoolId, academicYearId, semesterId, teacherId }) {
    await new Promise(resolve => setTimeout(resolve, 400))

    const schedules = await this.fetchSchedulesByContext({
      schoolId,
      academicYearId,
      semesterId,
      teacherId
    })

    const schedule = schedules.find(s => s.id === scheduleId)
    if (!schedule) {
      throw new Error(`Schedule not found: ${scheduleId}`)
    }

    const existingSession = mockSessions.find(
      s => s.scheduleId === scheduleId && s.date === date
    )
    if (existingSession) {
      throw new Error('Session already exists for this schedule and date')
    }

    const newSession = {
      id: `TS-${String(sessionCounter++).padStart(3, '0')}`,
      scheduleId: schedule.id,
      classId: schedule.classId,
      className: schedule.className,
      subjectId: schedule.subjectId,
      subjectName: schedule.subjectName,
      date,
      startTime: schedule.startTime,
      endTime: schedule.endTime,
      status: TEACHING_STATUS.SCHEDULED,
      teacherId: schedule.teacherId,
      teacherName: schedule.teacherName,
      teacherPresence: null,
      notes: '',
      academicYearId,
      semesterId,
      schoolId
    }

    mockSessions.push(newSession)
    return newSession
  }

  async startSession({ sessionId, teacherPresence }) {
    await new Promise(resolve => setTimeout(resolve, 300))

    const session = mockSessions.find(s => s.id === sessionId)
    if (!session) {
      throw new Error(`Session not found: ${sessionId}`)
    }

    if (session.status !== TEACHING_STATUS.SCHEDULED) {
      throw new Error(`Cannot start session with status: ${session.status}`)
    }

    session.status = TEACHING_STATUS.ACTIVE
    session.teacherPresence = teacherPresence

    return session
  }

  async endSession({ sessionId, notes }) {
    await new Promise(resolve => setTimeout(resolve, 300))

    const session = mockSessions.find(s => s.id === sessionId)
    if (!session) {
      throw new Error(`Session not found: ${sessionId}`)
    }

    if (session.status !== TEACHING_STATUS.ACTIVE) {
      throw new Error(`Cannot end session with status: ${session.status}`)
    }

    session.status = TEACHING_STATUS.COMPLETED
    session.notes = notes || ''

    return session
  }

  async cancelSession({ sessionId, reason }) {
    await new Promise(resolve => setTimeout(resolve, 300))

    const session = mockSessions.find(s => s.id === sessionId)
    if (!session) {
      throw new Error(`Session not found: ${sessionId}`)
    }

    if (session.status !== TEACHING_STATUS.SCHEDULED && session.status !== TEACHING_STATUS.ACTIVE) {
      throw new Error(`Cannot cancel session with status: ${session.status}`)
    }

    session.status = TEACHING_STATUS.CANCELLED
    session.notes = reason || ''

    return session
  }
}

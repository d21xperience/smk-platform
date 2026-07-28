import { TeacherSchedule } from '@/models/TeacherSchedule'
import { TeachingSession } from '@/models/TeachingSession'
import { TeacherPresence } from '@/models/TeacherPresence'
import { AuditEntry } from '@/models/AuditEntry'
import { JournalEntry } from '@/models/JournalEntry'
import { TeachingEngine, SessionStatus, AuditAction } from '@/engine/TeachingEngine'

// --- Jadwal (statis) ---
const scheduleList = [
  new TeacherSchedule({
    id: 1,
    dayOfWeek: 'Senin',
    startTime: '07:00',
    endTime: '09:30',
    subject: 'Matematika',
    className: 'XII RPL 1',
    room: 'R.201',
  }),
  new TeacherSchedule({
    id: 2,
    dayOfWeek: 'Senin',
    startTime: '10:00',
    endTime: '11:30',
    subject: 'Fisika',
    className: 'XII RPL 2',
    room: 'Lab IPA',
  }),
  new TeacherSchedule({
    id: 3,
    dayOfWeek: 'Selasa',
    startTime: '07:00',
    endTime: '09:30',
    subject: 'Matematika',
    className: 'XII RPL 1',
    room: 'R.201',
  }),
  new TeacherSchedule({
    id: 4,
    dayOfWeek: 'Rabu',
    startTime: '10:00',
    endTime: '11:30',
    subject: 'Fisika',
    className: 'XII RPL 2',
    room: 'Lab IPA',
  }),
]

// --- State in-memory ---
let sessions = [
  new TeachingSession({
    id: 1,
    scheduleId: 1,
    date: '2025-07-18',
    startTime: '07:00',
    endTime: '09:30',
    subject: 'Matematika',
    className: 'XII RPL 1',
    status: SessionStatus.SCHEDULED,
    teacherPresence: null,
    auditTrail: [],
  }),
  new TeachingSession({
    id: 2,
    scheduleId: 2,
    date: '2026-07-19',
    startTime: '10:00',
    endTime: '11:30',
    subject: 'Fisika',
    className: 'XII RPL 2',
    status: SessionStatus.STARTED,
    teacherPresence: new TeacherPresence({
      id: 1,
      sessionId: 2,
      teacherId: 1,
      checkInTime: '2025-07-18T09:55:00Z',
      checkOutTime: null,
      status: 'on_time',
    }),
    auditTrail: [
      new AuditEntry({
        id: 1,
        sessionId: 2,
        timestamp: '2025-07-18T09:55:00Z',
        action: AuditAction.SESSION_STARTED,
        details: 'Sesi dimulai',
      }),
      new AuditEntry({
        id: 2,
        sessionId: 2,
        timestamp: '2025-07-18T09:55:10Z',
        action: AuditAction.TEACHER_CHECK_IN,
        details: 'Guru check-in',
      }),
    ],
  }),
  new TeachingSession({
    id: 3,
    scheduleId: 3,
    date: '2025-07-21',
    startTime: '07:00',
    endTime: '09:30',
    subject: 'Matematika',
    className: 'XII RPL 1',
    status: SessionStatus.COMPLETED,
    teacherPresence: new TeacherPresence({
      id: 2,
      sessionId: 3,
      teacherId: 1,
      checkInTime: '2025-07-21T06:50:00Z',
      checkOutTime: '2025-07-21T09:30:00Z',
      status: 'on_time',
    }),
    auditTrail: [
      new AuditEntry({
        id: 3,
        sessionId: 3,
        timestamp: '2025-07-21T06:50:00Z',
        action: AuditAction.SESSION_STARTED,
        details: 'Sesi dimulai',
      }),
      new AuditEntry({
        id: 4,
        sessionId: 3,
        timestamp: '2025-07-21T06:50:00Z',
        action: AuditAction.TEACHER_CHECK_IN,
        details: 'Guru check-in',
      }),
      new AuditEntry({
        id: 5,
        sessionId: 3,
        timestamp: '2025-07-21T09:30:00Z',
        action: AuditAction.SESSION_COMPLETED,
        details: 'Sesi selesai',
      }),
    ],
  }),
]

// Jurnal tetap di array terpisah (dari sprint B)
let journals = [
  new JournalEntry({
    id: 1,
    sessionId: 3,
    date: '2025-07-21',
    subject: 'Matematika',
    className: 'XII RPL 1',
    material: 'Integral parsial',
    attendance: '28/30',
    obstacles: 'Beberapa siswa belum paham substitusi',
    notes: 'Berikan latihan tambahan',
  }),
]

// --- Helper untuk menghasilkan ID ---
let nextPresenceId = 3
let nextAuditId = 6

export const teachingMockAdapter = {
  // --- Jadwal (tidak berubah) ---
  // eslint-disable-next-line no-unused-vars
  async getTeacherSchedule(academicYearId, semesterId) {
    await new Promise((resolve) => setTimeout(resolve, 400))
    return [...scheduleList]
  },

  // --- Sesi (dengan lifecycle) ---
  async getTeachingSessions(date) {
    await new Promise((resolve) => setTimeout(resolve, 300))
    return sessions.filter((s) => s.date === date)
    // return sessions.filter((s) => s.date === date)
  },

  async startSession(sessionId) {
    await new Promise((resolve) => setTimeout(resolve, 200))
    const session = sessions.find((s) => s.id === sessionId)
    if (!session) throw new Error('Sesi tidak ditemukan')
    TeachingEngine.transition(session.status, SessionStatus.STARTED)
    session.status = SessionStatus.STARTED
    // Tambah audit
    const audit = new AuditEntry({
      id: nextAuditId++,
      sessionId,
      timestamp: new Date().toISOString(),
      action: AuditAction.SESSION_STARTED,
      details: 'Sesi dimulai',
    })
    session.auditTrail.push(audit)
    return new TeachingSession({
      ...session,
      teacherPresence: session.teacherPresence
        ? new TeacherPresence({ ...session.teacherPresence })
        : null,
      auditTrail: [...session.auditTrail],
    })
  },

  async beginProgress(sessionId) {
    await new Promise((resolve) => setTimeout(resolve, 200))
    const session = sessions.find((s) => s.id === sessionId)
    if (!session) throw new Error('Sesi tidak ditemukan')
    TeachingEngine.transition(session.status, SessionStatus.IN_PROGRESS)
    session.status = SessionStatus.IN_PROGRESS
    const audit = new AuditEntry({
      id: nextAuditId++,
      sessionId,
      timestamp: new Date().toISOString(),
      action: AuditAction.SESSION_IN_PROGRESS,
      details: 'Pembelajaran dimulai',
    })
    session.auditTrail.push(audit)
    return new TeachingSession({
      ...session,
      teacherPresence: session.teacherPresence
        ? new TeacherPresence({ ...session.teacherPresence })
        : null,
      auditTrail: [...session.auditTrail],
    })
  },

  async completeSession(sessionId) {
    await new Promise((resolve) => setTimeout(resolve, 200))
    const session = sessions.find((s) => s.id === sessionId)
    if (!session) throw new Error('Sesi tidak ditemukan')
    TeachingEngine.transition(session.status, SessionStatus.COMPLETED)
    session.status = SessionStatus.COMPLETED
    const audit = new AuditEntry({
      id: nextAuditId++,
      sessionId,
      timestamp: new Date().toISOString(),
      action: AuditAction.SESSION_COMPLETED,
      details: 'Sesi selesai',
    })
    session.auditTrail.push(audit)
    // Check-out otomatis jika belum
    if (session.teacherPresence && !session.teacherPresence.checkOutTime) {
      session.teacherPresence.checkOutTime = new Date().toISOString()
      session.auditTrail.push(
        new AuditEntry({
          id: nextAuditId++,
          sessionId,
          timestamp: new Date().toISOString(),
          action: AuditAction.TEACHER_CHECK_OUT,
          details: 'Check-out otomatis',
        }),
      )
    }
    return new TeachingSession({
      ...session,
      teacherPresence: session.teacherPresence
        ? new TeacherPresence({ ...session.teacherPresence })
        : null,
      auditTrail: [...session.auditTrail],
    })
  },

  async lockSession(sessionId) {
    await new Promise((resolve) => setTimeout(resolve, 200))
    const session = sessions.find((s) => s.id === sessionId)
    if (!session) throw new Error('Sesi tidak ditemukan')
    TeachingEngine.transition(session.status, SessionStatus.LOCKED)
    session.status = SessionStatus.LOCKED
    const audit = new AuditEntry({
      id: nextAuditId++,
      sessionId,
      timestamp: new Date().toISOString(),
      action: AuditAction.SESSION_LOCKED,
      details: 'Sesi dikunci',
    })
    session.auditTrail.push(audit)
    return new TeachingSession({
      ...session,
      teacherPresence: session.teacherPresence
        ? new TeacherPresence({ ...session.teacherPresence })
        : null,
      auditTrail: [...session.auditTrail],
    })
  },

  // --- Teacher Presence ---
  async checkInTeacher(sessionId, teacherId, checkInTime) {
    await new Promise((resolve) => setTimeout(resolve, 200))
    const session = sessions.find((s) => s.id === sessionId)
    if (!session) throw new Error('Sesi tidak ditemukan')
    if (session.teacherPresence) throw new Error('Guru sudah check-in')
    const scheduledStart = new Date(`${session.date}T${session.startTime}:00`)
    const status = TeachingEngine.calculatePresenceStatus(checkInTime, scheduledStart.toISOString())
    const presence = new TeacherPresence({
      id: nextPresenceId++,
      sessionId,
      teacherId,
      checkInTime,
      checkOutTime: null,
      status,
    })
    session.teacherPresence = presence
    const audit = new AuditEntry({
      id: nextAuditId++,
      sessionId,
      timestamp: checkInTime,
      action: AuditAction.TEACHER_CHECK_IN,
      details: `Guru hadir (${status})`,
    })
    session.auditTrail.push(audit)
    return new TeacherPresence({ ...presence })
  },

  async checkOutTeacher(sessionId, checkOutTime) {
    await new Promise((resolve) => setTimeout(resolve, 200))
    const session = sessions.find((s) => s.id === sessionId)
    if (!session) throw new Error('Sesi tidak ditemukan')
    if (!session.teacherPresence) throw new Error('Belum check-in')
    if (session.teacherPresence.checkOutTime) throw new Error('Sudah check-out')
    session.teacherPresence.checkOutTime = checkOutTime
    const audit = new AuditEntry({
      id: nextAuditId++,
      sessionId,
      timestamp: checkOutTime,
      action: AuditAction.TEACHER_CHECK_OUT,
      details: 'Guru check-out',
    })
    session.auditTrail.push(audit)
    return new TeacherPresence({ ...session.teacherPresence })
  },

  // --- Audit Trail ---
  async getAuditTrail(sessionId) {
    await new Promise((resolve) => setTimeout(resolve, 150))
    const session = sessions.find((s) => s.id === sessionId)
    return session ? session.auditTrail.map((e) => new AuditEntry({ ...e })) : []
  },

  async addAuditEntry(sessionId, action, details) {
    const session = sessions.find((s) => s.id === sessionId)
    if (!session) throw new Error('Sesi tidak ditemukan')
    const entry = new AuditEntry({
      id: nextAuditId++,
      sessionId,
      timestamp: new Date().toISOString(),
      action,
      details,
    })
    session.auditTrail.push(entry)
    return new AuditEntry({ ...entry })
  },

  // --- Jurnal (dari Sprint B) ---
  async saveJournal(sessionId, journalData) {
    await new Promise((resolve) => setTimeout(resolve, 300))
    const existing = journals.find((j) => j.sessionId === sessionId)
    if (existing) {
      Object.assign(existing, journalData)
      return new JournalEntry({ ...existing })
    } else {
      const newJournal = new JournalEntry({
        id: journals.length + 1,
        sessionId,
        date: '',
        subject: '',
        className: '',
        ...journalData,
      })
      journals.push(newJournal)
      return newJournal
    }
  },

  async getJournalBySession(sessionId) {
    await new Promise((resolve) => setTimeout(resolve, 200))
    return journals.find((j) => j.sessionId === sessionId) || null
  },
}

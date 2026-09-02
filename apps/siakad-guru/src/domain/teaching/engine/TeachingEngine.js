import { TeachingSession } from '../models/TeachingSession.js'
import { TeachingSchedule } from '../models/TeachingSchedule.js'
import {
  TEACHING_STATUS,
  isValidTeachingStatus,
  isValidTeacherPresence,
  canStartSession,
  canEndSession,
  canCancelSession,
} from '../models/TeachingStatus.js'

export class TeachingEngine {
  static validateSession(session) {
    if (!session || !(session instanceof TeachingSession)) {
      throw new Error('Invalid teaching session object')
    }
    if (!session.isValid()) {
      throw new Error('Teaching session is missing required fields')
    }
    if (!isValidTeachingStatus(session.status)) {
      throw new Error(`Invalid teaching status: ${session.status}`)
    }
    return true
  }

  static validateSchedule(schedule) {
    if (!schedule || !(schedule instanceof TeachingSchedule)) {
      throw new Error('Invalid teaching schedule object')
    }
    if (!schedule.isValid()) {
      throw new Error('Teaching schedule is missing required fields')
    }
    return true
  }

  static startSession(session, teacherPresence) {
    TeachingEngine.validateSession(session)

    if (!canStartSession(session.status)) {
      throw new Error(`Cannot start session with status: ${session.status}`)
    }

    if (!isValidTeacherPresence(teacherPresence)) {
      throw new Error(`Invalid teacher presence: ${teacherPresence}`)
    }

    session.status = TEACHING_STATUS.ACTIVE
    session.teacherPresence = teacherPresence

    return session
  }

  static endSession(session, notes) {
    TeachingEngine.validateSession(session)

    if (!canEndSession(session.status)) {
      throw new Error(`Cannot end session with status: ${session.status}`)
    }

    if (!session.teacherPresence) {
      throw new Error('Teacher presence must be set before ending session')
    }

    session.status = TEACHING_STATUS.COMPLETED
    if (notes) {
      session.notes = notes
    }

    return session
  }

  static cancelSession(session, reason) {
    TeachingEngine.validateSession(session)

    if (!canCancelSession(session.status)) {
      throw new Error(`Cannot cancel session with status: ${session.status}`)
    }

    session.status = TEACHING_STATUS.CANCELLED
    if (reason) {
      session.notes = reason
    }

    return session
  }

  static validateTeacherPresence(presence) {
    if (!isValidTeacherPresence(presence)) {
      throw new Error(`Invalid teacher presence: ${presence}`)
    }
    return true
  }
}

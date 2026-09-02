import { TeachingEngine } from '../domain/teaching/engine/TeachingEngine.js'
import { TeachingSchedule } from '../domain/teaching/models/TeachingSchedule.js'
import { TeachingSession } from '../domain/teaching/models/TeachingSession.js'
import { TeachingSessionStartedEvent } from '../domain/teaching/events/TeachingSessionStarted.js'
import { TeachingSessionEndedEvent } from '../domain/teaching/events/TeachingSessionEnded.js'

export class TeachingService {
  constructor({ teachingAdapter, eventDispatcher }) {
    this.teachingAdapter = teachingAdapter
    this.eventDispatcher = eventDispatcher
  }

  async loadSchedulesByContext({ schoolId, academicYearId, semesterId, teacherId }) {
    const rawData = await this.teachingAdapter.fetchSchedulesByContext({
      schoolId,
      academicYearId,
      semesterId,
      teacherId,
    })

    return rawData.map((item) => new TeachingSchedule(item))
  }

  async loadSessionsByDate({ schoolId, academicYearId, semesterId, teacherId, date }) {
    const rawData = await this.teachingAdapter.fetchSessionsByDate({
      schoolId,
      academicYearId,
      semesterId,
      teacherId,
      date,
    })

    return rawData.map((item) => new TeachingSession(item))
  }

  async createSession({ scheduleId, date, schoolId, academicYearId, semesterId, teacherId }) {
    const rawData = await this.teachingAdapter.createSession({
      scheduleId,
      date,
      schoolId,
      academicYearId,
      semesterId,
      teacherId,
    })

    const session = new TeachingSession(rawData)
    TeachingEngine.validateSession(session)

    return session
  }

  // eslint-disable-next-line no-unused-vars
  async startSession({ sessionId, teacherPresence, schoolId, academicYearId, semesterId }) {
    // Validate teacher presence using domain engine
    TeachingEngine.validateTeacherPresence(teacherPresence)

    const rawData = await this.teachingAdapter.startSession({
      sessionId,
      teacherPresence,
    })

    const session = new TeachingSession(rawData)
    TeachingEngine.validateSession(session)

    // Dispatch event
    if (this.eventDispatcher) {
      this.eventDispatcher.dispatch(new TeachingSessionStartedEvent(session))
    }

    return session
  }

  async endSession({ sessionId, notes, schoolId, academicYearId, semesterId }) {
    const rawData = await this.teachingAdapter.endSession({
      sessionId,
      notes,
      schoolId,
      academicYearId,
      semesterId,
    })

    const session = new TeachingSession(rawData)
    TeachingEngine.validateSession(session)

    // Dispatch event
    if (this.eventDispatcher) {
      this.eventDispatcher.dispatch(new TeachingSessionEndedEvent(session))
    }

    return session
  }

  async cancelSession({ sessionId, reason, schoolId, academicYearId, semesterId }) {
    const rawData = await this.teachingAdapter.cancelSession({
      sessionId,
      reason,
      schoolId,
      academicYearId,
      semesterId,
    })

    const session = new TeachingSession(rawData)
    TeachingEngine.validateSession(session)

    return session
  }
}

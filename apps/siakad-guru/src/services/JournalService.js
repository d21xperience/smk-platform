import { JournalEngine } from '../domain/journal/engine/JournalEngine.js'
import { TeachingJournal } from '../domain/journal/models/TeachingJournal.js'
import { JournalCreatedEvent } from '../domain/journal/events/JournalCreated.js'
import { JournalSubmittedEvent } from '../domain/journal/events/JournalSubmitted.js'

export class JournalService {
  constructor({ journalAdapter, eventDispatcher }) {
    this.journalAdapter = journalAdapter
    this.eventDispatcher = eventDispatcher
  }

  async loadJournal({ teachingSessionId }) {
    const rawData = await this.journalAdapter.loadJournalByTeachingSession({ teachingSessionId })

    if (!rawData) {
      return null
    }

    const journal = new TeachingJournal(rawData)
    JournalEngine.validateJournal(journal)
    return journal
  }

  async createJournal({
    teachingSessionId,
    classId,
    className,
    subjectId,
    subjectName,
    date,
    schoolId,
    academicYearId,
    semesterId
  }) {
    const rawData = await this.journalAdapter.createJournal({
      teachingSessionId,
      classId,
      className,
      subjectId,
      subjectName,
      date,
      schoolId,
      academicYearId,
      semesterId
    })

    const journal = new TeachingJournal(rawData)
    JournalEngine.validateJournal(journal)

    if (this.eventDispatcher) {
      this.eventDispatcher.dispatch(new JournalCreatedEvent(journal))
    }

    return journal
  }

  async updateJournal({ journalId, material, activities, reflection }) {
    const rawData = await this.journalAdapter.updateJournal({
      journalId,
      material,
      activities,
      reflection
    })

    const journal = new TeachingJournal(rawData)
    JournalEngine.validateJournal(journal)

    return journal
  }

  async submitJournal({ journalId }) {
    const rawData = await this.journalAdapter.submitJournal({ journalId })

    const journal = new TeachingJournal(rawData)
    JournalEngine.validateJournal(journal)
    JournalEngine.submitJournal(journal)

    if (this.eventDispatcher) {
      this.eventDispatcher.dispatch(new JournalSubmittedEvent(journal))
    }

    return journal
  }
}

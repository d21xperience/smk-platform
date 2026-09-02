import { EVENT_JOURNAL_CREATED } from '../../../events/Type.js'

export class JournalCreatedEvent {
  constructor(journal) {
    this.type = EVENT_JOURNAL_CREATED
    this.payload = {
      journalId: journal.id,
      teachingSessionId: journal.teachingSessionId,
      classId: journal.classId,
      subjectId: journal.subjectId,
      date: journal.date,
      academicYearId: journal.academicYearId,
      semesterId: journal.semesterId,
      schoolId: journal.schoolId,
      timestamp: new Date().toISOString()
    }
  }
}

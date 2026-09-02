import { JOURNAL_STATUS } from '../../domain/journal/models/JournalStatus.js'

// Mock journals storage (keyed by teachingSessionId)
let mockJournals = {}
let journalCounter = 1

export class JournalMockAdapter {
  async loadJournalByTeachingSession({ teachingSessionId }) {
    await new Promise(resolve => setTimeout(resolve, 400))
    return mockJournals[teachingSessionId] || null
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
    await new Promise(resolve => setTimeout(resolve, 400))

    if (mockJournals[teachingSessionId]) {
      throw new Error('Journal already exists for this teaching session')
    }

    const newJournal = {
      id: `JRN-${String(journalCounter++).padStart(3, '0')}`,
      teachingSessionId,
      classId,
      className,
      subjectId,
      subjectName,
      date,
      material: '',
      activities: '',
      reflection: '',
      status: JOURNAL_STATUS.DRAFT,
      academicYearId,
      semesterId,
      schoolId
    }

    mockJournals[teachingSessionId] = newJournal
    return newJournal
  }

  async updateJournal({ journalId, material, activities, reflection }) {
    await new Promise(resolve => setTimeout(resolve, 300))

    let targetJournal = null
    for (const key in mockJournals) {
      if (mockJournals[key].id === journalId) {
        targetJournal = mockJournals[key]
        break
      }
    }

    if (!targetJournal) {
      throw new Error(`Journal not found: ${journalId}`)
    }

    if (targetJournal.status !== JOURNAL_STATUS.DRAFT) {
      throw new Error('Cannot update a submitted journal')
    }

    if (material !== undefined) {
      targetJournal.material = material
    }
    if (activities !== undefined) {
      targetJournal.activities = activities
    }
    if (reflection !== undefined) {
      targetJournal.reflection = reflection
    }

    return targetJournal
  }

  async submitJournal({ journalId }) {
    await new Promise(resolve => setTimeout(resolve, 400))

    let targetJournal = null
    for (const key in mockJournals) {
      if (mockJournals[key].id === journalId) {
        targetJournal = mockJournals[key]
        break
      }
    }

    if (!targetJournal) {
      throw new Error(`Journal not found: ${journalId}`)
    }

    if (targetJournal.status !== JOURNAL_STATUS.DRAFT) {
      throw new Error('Journal has already been submitted')
    }

    if (!targetJournal.material.trim() || !targetJournal.activities.trim()) {
      throw new Error('Journal must have material and activities filled before submitting')
    }

    targetJournal.status = JOURNAL_STATUS.SUBMITTED
    return targetJournal
  }
}

import { TeachingJournal } from '../models/TeachingJournal.js'
import { JOURNAL_STATUS, isValidJournalStatus, canSubmitJournal } from '../models/JournalStatus.js'

export class JournalEngine {
  static validateJournal(journal) {
    if (!journal || !(journal instanceof TeachingJournal)) {
      throw new Error('Invalid teaching journal object')
    }
    if (!journal.isValid()) {
      throw new Error('Teaching journal is missing required fields')
    }
    if (!isValidJournalStatus(journal.status)) {
      throw new Error(`Invalid journal status: ${journal.status}`)
    }
    return true
  }

  static submitJournal(journal) {
    JournalEngine.validateJournal(journal)

    if (!canSubmitJournal(journal.status)) {
      throw new Error(`Cannot submit journal with status: ${journal.status}`)
    }

    if (!journal.isComplete()) {
      throw new Error('Journal must have material and activities filled before submitting')
    }

    journal.status = JOURNAL_STATUS.SUBMITTED
    return journal
  }

  static updateJournal(journal, { material, activities, reflection }) {
    JournalEngine.validateJournal(journal)

    if (journal.isSubmitted()) {
      throw new Error('Cannot update a submitted journal')
    }

    if (material !== undefined) {
      journal.material = material
    }
    if (activities !== undefined) {
      journal.activities = activities
    }
    if (reflection !== undefined) {
      journal.reflection = reflection
    }

    return journal
  }
}

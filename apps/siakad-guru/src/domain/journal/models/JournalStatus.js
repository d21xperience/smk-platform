export const JOURNAL_STATUS = {
  DRAFT: 'draft',
  SUBMITTED: 'submitted',
}

export function isValidJournalStatus(status) {
  return Object.values(JOURNAL_STATUS).includes(status)
}

export function canSubmitJournal(status) {
  return status === JOURNAL_STATUS.DRAFT
}

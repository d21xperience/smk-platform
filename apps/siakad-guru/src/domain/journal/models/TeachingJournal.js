import { JOURNAL_STATUS } from './JournalStatus.js'

export class TeachingJournal {
  constructor({
    id,
    teachingSessionId,
    classId,
    className,
    subjectId,
    subjectName,
    date,
    material,
    activities,
    reflection,
    status,
    academicYearId,
    semesterId,
    schoolId
  }) {
    this.id = id
    this.teachingSessionId = teachingSessionId
    this.classId = classId
    this.className = className
    this.subjectId = subjectId
    this.subjectName = subjectName
    this.date = date
    this.material = material || ''
    this.activities = activities || ''
    this.reflection = reflection || ''
    this.status = status || JOURNAL_STATUS.DRAFT
    this.academicYearId = academicYearId
    this.semesterId = semesterId
    this.schoolId = schoolId
  }

  isDraft() {
    return this.status === JOURNAL_STATUS.DRAFT
  }

  isSubmitted() {
    return this.status === JOURNAL_STATUS.SUBMITTED
  }

  isComplete() {
    return !!(
      this.material.trim() &&
      this.activities.trim()
    )
  }

  isValid() {
    return !!(
      this.id &&
      this.teachingSessionId &&
      this.classId &&
      this.subjectId &&
      this.date &&
      this.status &&
      this.academicYearId &&
      this.semesterId &&
      this.schoolId
    )
  }
}

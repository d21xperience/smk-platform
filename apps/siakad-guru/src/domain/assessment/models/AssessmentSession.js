import { ASSESSMENT_STATUS } from './AssessmentStatus.js'

export class AssessmentSession {
  constructor({
    id,
    classId,
    className,
    subjectId,
    subjectName,
    date,
    components,
    scores,
    status,
    academicYearId,
    semesterId,
    schoolId,
  }) {
    this.id = id
    this.classId = classId
    this.className = className
    this.subjectId = subjectId
    this.subjectName = subjectName
    this.date = date
    this.components = components || []
    this.scores = scores || {}
    this.status = status || ASSESSMENT_STATUS.DRAFT
    this.academicYearId = academicYearId
    this.semesterId = semesterId
    this.schoolId = schoolId
  }

  isDraft() {
    return this.status === ASSESSMENT_STATUS.DRAFT
  }

  isFinalized() {
    return this.status === ASSESSMENT_STATUS.FINALIZED
  }

  getTotalWeight() {
    return this.components.reduce((sum, comp) => sum + comp.weight, 0)
  }

  isValid() {
    return !!(
      this.id &&
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

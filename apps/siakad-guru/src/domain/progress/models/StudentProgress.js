import { PROGRESS_TYPE } from './ProgressType.js'

export class StudentProgress {
  constructor({
    id,
    studentId,
    studentName,
    classId,
    className,
    type,
    title,
    description,
    date,
    teacherId,
    teacherName,
    academicYearId,
    semesterId,
    schoolId
  }) {
    this.id = id
    this.studentId = studentId
    this.studentName = studentName
    this.classId = classId
    this.className = className
    this.type = type
    this.title = title || ''
    this.description = description || ''
    this.date = date
    this.teacherId = teacherId
    this.teacherName = teacherName
    this.academicYearId = academicYearId
    this.semesterId = semesterId
    this.schoolId = schoolId
  }

  isNote() {
    return this.type === PROGRESS_TYPE.NOTE
  }

  isAchievement() {
    return this.type === PROGRESS_TYPE.ACHIEVEMENT
  }

  isViolation() {
    return this.type === PROGRESS_TYPE.VIOLATION
  }

  isCounseling() {
    return this.type === PROGRESS_TYPE.COUNSELING
  }

  isValid() {
    return !!(
      this.id &&
      this.studentId &&
      this.studentName &&
      this.classId &&
      this.type &&
      this.date &&
      this.teacherId &&
      this.academicYearId &&
      this.semesterId &&
      this.schoolId
    )
  }
}

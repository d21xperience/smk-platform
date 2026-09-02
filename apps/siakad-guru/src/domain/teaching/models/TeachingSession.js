import { TEACHING_STATUS } from './TeachingStatus.js'

export class TeachingSession {
  constructor({
    id,
    scheduleId,
    classId,
    className,
    subjectId,
    subjectName,
    date,
    startTime,
    endTime,
    status,
    teacherId,
    teacherName,
    teacherPresence,
    notes,
    academicYearId,
    semesterId,
    schoolId,
  }) {
    this.id = id
    this.scheduleId = scheduleId
    this.classId = classId
    this.className = className
    this.subjectId = subjectId
    this.subjectName = subjectName
    this.date = date
    this.startTime = startTime
    this.endTime = endTime
    this.status = status || TEACHING_STATUS.SCHEDULED
    this.teacherId = teacherId
    this.teacherName = teacherName
    this.teacherPresence = teacherPresence || null
    this.notes = notes || ''
    this.academicYearId = academicYearId
    this.semesterId = semesterId
    this.schoolId = schoolId
  }

  isActive() {
    return this.status === TEACHING_STATUS.ACTIVE
  }

  isCompleted() {
    return this.status === TEACHING_STATUS.COMPLETED
  }

  isCancelled() {
    return this.status === TEACHING_STATUS.CANCELLED
  }

  isValid() {
    return !!(
      this.id &&
      this.scheduleId &&
      this.classId &&
      this.subjectId &&
      this.date &&
      this.startTime &&
      this.endTime &&
      this.status &&
      this.teacherId &&
      this.academicYearId &&
      this.semesterId &&
      this.schoolId
    )
  }
}

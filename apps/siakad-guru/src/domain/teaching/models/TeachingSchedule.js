export class TeachingSchedule {
  constructor({
    id,
    classId,
    className,
    subjectId,
    subjectName,
    dayOfWeek,
    startTime,
    endTime,
    teacherId,
    teacherName,
    academicYearId,
    semesterId,
    schoolId
  }) {
    this.id = id
    this.classId = classId
    this.className = className
    this.subjectId = subjectId
    this.subjectName = subjectName
    this.dayOfWeek = dayOfWeek
    this.startTime = startTime
    this.endTime = endTime
    this.teacherId = teacherId
    this.teacherName = teacherName
    this.academicYearId = academicYearId
    this.semesterId = semesterId
    this.schoolId = schoolId
  }

  isValid() {
    return !!(
      this.id &&
      this.classId &&
      this.subjectId &&
      this.dayOfWeek &&
      this.startTime &&
      this.endTime &&
      this.teacherId &&
      this.academicYearId &&
      this.semesterId &&
      this.schoolId
    )
  }
}

import { ProgressEngine } from '../domain/progress/engine/ProgressEngine.js'
import { StudentProgress } from '../domain/progress/models/StudentProgress.js'
import { ProgressRecordedEvent } from '../domain/progress/events/ProgressRecorded.js'

export class ProgressService {
  constructor({ progressAdapter, eventDispatcher }) {
    this.progressAdapter = progressAdapter
    this.eventDispatcher = eventDispatcher
  }

  async loadStudentsByClass({ classId }) {
    return await this.progressAdapter.fetchStudentsByClass({ classId })
  }

  async loadProgressByStudent({ studentId, academicYearId, semesterId }) {
    const rawData = await this.progressAdapter.loadProgressByStudent({
      studentId,
      academicYearId,
      semesterId,
    })

    return rawData.map((item) => new StudentProgress(item))
  }

  async createProgress({
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
    schoolId,
    academicYearId,
    semesterId,
  }) {
    // Validate input using domain engine
    ProgressEngine.validateProgressInput({ type, title, description })

    const rawData = await this.progressAdapter.createProgress({
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
      schoolId,
      academicYearId,
      semesterId,
    })

    const progress = new StudentProgress(rawData)
    ProgressEngine.validateProgress(progress)

    if (this.eventDispatcher) {
      this.eventDispatcher.dispatch(new ProgressRecordedEvent(progress))
    }

    return progress
  }
}

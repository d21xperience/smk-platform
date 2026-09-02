import { EVENT_PROGRESS_RECORDED } from '../../../events/Type.js'

export class ProgressRecordedEvent {
  constructor(progress) {
    this.type = EVENT_PROGRESS_RECORDED
    this.payload = {
      progressId: progress.id,
      studentId: progress.studentId,
      classId: progress.classId,
      progressType: progress.type,
      date: progress.date,
      teacherId: progress.teacherId,
      academicYearId: progress.academicYearId,
      semesterId: progress.semesterId,
      schoolId: progress.schoolId,
      timestamp: new Date().toISOString()
    }
  }
}

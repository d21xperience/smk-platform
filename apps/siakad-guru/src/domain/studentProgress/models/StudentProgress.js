import { Achievement } from '@/models/Achievement'
import { CounselingRecord } from '@/models/CounselingRecord'
import { TeacherNote } from '@/models/TeacherNote'
import { Violation } from '@/models/Violation'
import { BaseEntity } from '@/domain/shared/BaseEntity'

export class StudentProgress extends BaseEntity {
  constructor(props) {
    super(props)
    this.studentId = props.studentId
    this.classId = props.classId
    this.academicYearId = props.academicYearId
    this.semesterId = props.semesterId
    this.studentName = props.studentName || ''
    this.nis = props.nis || ''
    // Data akademik
    this.averageScore = props.averageScore || 0
    this.coreSubjectsGrades = props.coreSubjectsGrades || {} // { subjectId: score }
    this.promotionStatus = props.promotionStatus || 'UNDEFINED' // 'NAIK' | 'TIDAK_NAIK' | 'PERHATIAN'
    // Poin
    this.achievementPoints = props.achievementPoints || 0
    this.violationPoints = props.violationPoints || 0
    // Kehadiran (dari attendance)
    this.totalPresent = props.totalPresent || 0
    this.totalSick = props.totalSick || 0
    this.totalPermit = props.totalPermit || 0
    this.totalAbsent = props.totalAbsent || 0
    this.totalLate = props.totalLate || 0
    // Catatan
    this.teacherNotes = (props.teacherNotes || []).map((n) => new TeacherNote(n))
    this.achievements = (props.achievements || []).map((a) => new Achievement(a))
    this.violations = (props.violations || []).map((v) => new Violation(v))
    this.counselingRecords = (props.counselingRecords || []).map((c) => new CounselingRecord(c))
    // Metadata
    this.lastUpdated = props.lastUpdated || new Date().toISOString()
  }

  addTeacherNote(note) {
    this.teacherNotes.push(new TeacherNote({ ...note, studentId: this.studentId }))
    this.lastUpdated = new Date().toISOString()
    this.updateTimestamps()
  }

  addAchievement(achievement) {
    this.achievements.push(new Achievement({ ...achievement, studentId: this.studentId }))
    this.achievementPoints += achievement.points || 0
    this.lastUpdated = new Date().toISOString()
    this.updateTimestamps()
  }

  addViolation(violation) {
    this.violations.push(new Violation({ ...violation, studentId: this.studentId }))
    this.violationPoints += violation.points || 0
    this.lastUpdated = new Date().toISOString()
    this.updateTimestamps()
  }

  updateAcademicData(averageScore, coreSubjectsGrades) {
    this.averageScore = averageScore
    this.coreSubjectsGrades = coreSubjectsGrades
    this.lastUpdated = new Date().toISOString()
    this.updateTimestamps()
  }

  updateAttendanceStats(stats) {
    this.totalPresent = stats.present || 0
    this.totalSick = stats.sick || 0
    this.totalPermit = stats.permit || 0
    this.totalAbsent = stats.absent || 0
    this.totalLate = stats.late || 0
    this.lastUpdated = new Date().toISOString()
    this.updateTimestamps()
  }

  toJSON() {
    return {
      ...super.toJSON(),
      studentId: this.studentId,
      classId: this.classId,
      academicYearId: this.academicYearId,
      semesterId: this.semesterId,
      studentName: this.studentName,
      nis: this.nis,
      averageScore: this.averageScore,
      coreSubjectsGrades: this.coreSubjectsGrades,
      promotionStatus: this.promotionStatus,
      achievementPoints: this.achievementPoints,
      violationPoints: this.violationPoints,
      totalPresent: this.totalPresent,
      totalSick: this.totalSick,
      totalPermit: this.totalPermit,
      totalAbsent: this.totalAbsent,
      totalLate: this.totalLate,
      teacherNotes: this.teacherNotes.map((n) => n.toJSON()),
      achievements: this.achievements.map((a) => a.toJSON()),
      violations: this.violations.map((v) => v.toJSON()),
      counselingRecords: this.counselingRecords.map((c) => c.toJSON()),
      lastUpdated: this.lastUpdated,
    }
  }
}

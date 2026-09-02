import { BaseEntity } from '@/domain/shared/BaseEntity'
import { Material } from './Material'
import { LearningActivity } from './LearningActivity'
import { Reflection } from './Reflection'

export const JOURNAL_STATUS = {
  DRAFT: 'draft',
  FINALIZED: 'finalized',
}

export class TeachingJournal extends BaseEntity {
  constructor(props) {
    super(props)
    this.sessionId = props.sessionId
    this.teacherId = props.teacherId
    this.classId = props.classId
    this.subjectId = props.subjectId
    this.academicYearId = props.academicYearId
    this.semesterId = props.semesterId
    this.date = props.date // YYYY-MM-DD
    this.status = props.status || JOURNAL_STATUS.DRAFT
    this.material = props.material ? new Material(props.material) : new Material({})
    this.learningActivities = (props.learningActivities || []).map((a) =>
      a instanceof LearningActivity ? a : LearningActivity.fromJSON(a),
    )
    this.reflection = props.reflection ? new Reflection(props.reflection) : new Reflection({})
    this.finalizedAt = props.finalizedAt || null
    this.attendanceSubmitted = props.attendanceSubmitted || false
  }

  isDraft() {
    return this.status === JOURNAL_STATUS.DRAFT
  }

  isFinalized() {
    return this.status === JOURNAL_STATUS.FINALIZED
  }

  // Menambahkan aktivitas belajar
  addActivity(activity) {
    if (this.isFinalized()) {
      throw new Error('Tidak dapat menambahkan aktivitas ke jurnal yang sudah di-finalisasi.')
    }
    const newActivity =
      activity instanceof LearningActivity ? activity : new LearningActivity(activity)
    newActivity.order = this.learningActivities.length
    this.learningActivities.push(newActivity)
    this.updateTimestamps()
  }

  // Menghapus aktivitas
  removeActivity(index) {
    if (this.isFinalized()) {
      throw new Error('Tidak dapat menghapus aktivitas dari jurnal yang sudah di-finalisasi.')
    }
    if (index >= 0 && index < this.learningActivities.length) {
      this.learningActivities.splice(index, 1)
      this.updateTimestamps()
    }
  }

  // Memperbarui materi
  updateMaterial(material) {
    if (this.isFinalized()) {
      throw new Error('Tidak dapat mengubah materi jurnal yang sudah di-finalisasi.')
    }
    this.material = material instanceof Material ? material : new Material(material)
    this.updateTimestamps()
  }

  // Memperbarui refleksi
  updateReflection(reflection) {
    if (this.isFinalized()) {
      throw new Error('Tidak dapat mengubah refleksi jurnal yang sudah di-finalisasi.')
    }
    this.reflection = reflection instanceof Reflection ? reflection : new Reflection(reflection)
    this.updateTimestamps()
  }

  // Finalisasi jurnal
  finalize() {
    if (this.isFinalized()) {
      throw new Error('Jurnal sudah di-finalisasi.')
    }
    this.status = JOURNAL_STATUS.FINALIZED
    this.finalizedAt = new Date().toISOString()
    this.updateTimestamps()
  }

  // Menandai bahwa absensi sudah di-submit untuk sesi ini (sinkronisasi dengan Attendance)
  markAttendanceSubmitted() {
    this.attendanceSubmitted = true
    this.updateTimestamps()
  }

  // Mengecek apakah jurnal lengkap (untuk validasi final)
  isComplete() {
    return (
      !this.material.isEmpty() && this.learningActivities.length > 0 && !this.reflection.isEmpty()
    )
  }

  toJSON() {
    return {
      ...super.toJSON(),
      sessionId: this.sessionId,
      teacherId: this.teacherId,
      classId: this.classId,
      subjectId: this.subjectId,
      academicYearId: this.academicYearId,
      semesterId: this.semesterId,
      date: this.date,
      status: this.status,
      material: this.material.toJSON(),
      learningActivities: this.learningActivities.map((a) => a.toJSON()),
      reflection: this.reflection.toJSON(),
      finalizedAt: this.finalizedAt,
      attendanceSubmitted: this.attendanceSubmitted,
    }
  }
}

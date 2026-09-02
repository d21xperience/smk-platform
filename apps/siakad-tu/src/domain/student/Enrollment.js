// apps/siakad-tu/src/domain/student/Enrollment.js

/**
 * Entity: Enrollment
 *
 * Merepresentasikan riwayat pendaftaran siswa di suatu kelas pada periode tertentu.
 * Enrollment adalah bagian dari Student Aggregate.
 *
 * Identity: enrollmentId
 * Lifecycle: terikat pada Student
 */
export class Enrollment {
  #enrollmentId
  #schoolId // dari OperationalContext
  #periodId // dari OperationalContext
  #classId
  #enrollmentDate
  #status // 'active', 'completed', 'transferred'

  constructor({
    enrollmentId,
    schoolId,
    periodId,
    classId,
    enrollmentDate = new Date(),
    status = 'active',
  }) {
    if (!enrollmentId) throw new Error('Enrollment: enrollmentId wajib.')
    if (!schoolId) throw new Error('Enrollment: schoolId wajib.')
    if (!periodId) throw new Error('Enrollment: periodId wajib.')
    if (!classId) throw new Error('Enrollment: classId wajib.')

    this.#enrollmentId = enrollmentId
    this.#schoolId = schoolId
    this.#periodId = periodId
    this.#classId = classId
    this.#enrollmentDate =
      enrollmentDate instanceof Date ? enrollmentDate : new Date(enrollmentDate)
    this.#status = status
  }

  get enrollmentId() {
    return this.#enrollmentId
  }
  get schoolId() {
    return this.#schoolId
  }
  get periodId() {
    return this.#periodId
  }
  get classId() {
    return this.#classId
  }
  get enrollmentDate() {
    return this.#enrollmentDate
  }
  get status() {
    return this.#status
  }

  isActive() {
    return this.#status === 'active'
  }

  /**
   * Cek apakah enrollment ini untuk periode yang sama dengan yang diberikan
   */
  isForPeriod(periodId) {
    return this.#periodId === periodId
  }

  /**
   * Complete enrollment (saat siswa lulus atau pindah)
   */
  complete() {
    if (this.#status !== 'active') {
      throw new Error(
        `Enrollment ${this.#enrollmentId} tidak bisa di-complete. Status saat ini: ${this.#status}`,
      )
    }
    this.#status = 'completed'
  }

  toJSON() {
    return {
      enrollmentId: this.#enrollmentId,
      schoolId: this.#schoolId,
      periodId: this.#periodId,
      classId: this.#classId,
      enrollmentDate: this.#enrollmentDate.toISOString(),
      status: this.#status,
    }
  }

  static fromJSON(data) {
    return new Enrollment(data)
  }
}

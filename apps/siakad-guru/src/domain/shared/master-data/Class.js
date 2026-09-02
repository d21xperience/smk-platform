// domain/shared/master-data/Class.js
import { Entity } from '../Entity.js'

export class Class extends Entity {
  constructor({
    id,
    academicYearId,
    grade,
    major,
    parallelClass,
    homeroomTeacherId,
    roomId,
    ...props
  }) {
    super({ id, ...props })
    this.academicYearId = academicYearId
    this.grade = grade // 10, 11, 12 (numeric)
    this.major = major // 'TKJ', 'RPL', 'AKL'
    this.parallelClass = parallelClass // 1, 2, 3 (numeric)
    this.homeroomTeacherId = homeroomTeacherId
    this.roomId = roomId
  }

  // Display name via formatter (dipanggil di UI, bukan di sini)
  get displayName() {
    return `${this.grade} ${this.major} ${this.parallelClass}` // Format: "X TKJ 1"
  }
}

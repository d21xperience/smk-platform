// apps/siakad-tu/src/domain/student/value-objects/StudentStatus.js

/**
 * Value Object: StudentStatus
 * Enum status siswa
 */
export class StudentStatus {
  static ACTIVE = new StudentStatus('ACTIVE', 'Aktif');
  static GRADUATED = new StudentStatus('GRADUATED', 'Lulus');
  static TRANSFERRED = new StudentStatus('TRANSFERRED', 'Pindah');
  static DROPPED = new StudentStatus('DROPPED', 'Keluar');
  static ALUMNI = new StudentStatus('ALUMNI', 'Alumni');

  static VALUES = [
    StudentStatus.ACTIVE,
    StudentStatus.GRADUATED,
    StudentStatus.TRANSFERRED,
    StudentStatus.DROPPED,
    StudentStatus.ALUMNI
  ];

  #code;
  #label;

  constructor(code, label) {
    this.#code = code;
    this.#label = label;
    Object.freeze(this);
  }

  get code() { return this.#code; }
  get label() { return this.#label; }

  isActive() {
    return this.#code === 'ACTIVE';
  }

  canBeEnrolled() {
    // Hanya ACTIVE yang bisa di-enroll
    return this.#code === 'ACTIVE';
  }

  equals(other) {
    if (!(other instanceof StudentStatus)) return false;
    return this.#code === other.#code;
  }

  toString() {
    return this.#code;
  }

  toJSON() {
    return this.#code;
  }

  static fromJSON(code) {
    const found = StudentStatus.VALUES.find(s => s.#code === code);
    if (!found) {
      throw new Error(`StudentStatus tidak valid: "${code}".`);
    }
    return found;
  }
}

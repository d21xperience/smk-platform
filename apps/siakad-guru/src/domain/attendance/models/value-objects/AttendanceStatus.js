// PURE JS - Value Object
export const ATTENDANCE_STATUS = {
  PRESENT: 'present',
  SICK: 'sick',
  PERMIT: 'permit',
  ABSENT: 'absent',
  LATE: 'late',
}

export class AttendanceStatus {
  constructor(value) {
    const validValues = Object.values(ATTENDANCE_STATUS)
    if (!validValues.includes(value)) {
      throw new Error(
        `Status kehadiran tidak valid. Harus salah satu dari: ${validValues.join(', ')}`,
      )
    }
    this.value = value
  }

  isPresent() {
    return this.value === ATTENDANCE_STATUS.PRESENT
  }

  isAbsent() {
    return this.value === ATTENDANCE_STATUS.ABSENT
  }

  isLate() {
    return this.value === ATTENDANCE_STATUS.LATE
  }

  isExcused() {
    return this.value === ATTENDANCE_STATUS.SICK || this.value === ATTENDANCE_STATUS.PERMIT
  }

  equals(other) {
    if (!other) return false
    return this.value === other.value
  }

  toJSON() {
    return this.value
  }

  toString() {
    return this.value
  }
}

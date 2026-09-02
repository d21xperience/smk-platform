export const ATTENDANCE_STATUS = {
  PRESENT: 'present',
  SICK: 'sick',
  PERMISSION: 'permission',
  ABSENT: 'absent'
}

export const SESSION_STATUS = {
  DRAFT: 'draft',
  SUBMITTED: 'submitted'
}

export function isValidAttendanceStatus(status) {
  return Object.values(ATTENDANCE_STATUS).includes(status)
}

export function isValidSessionStatus(status) {
  return Object.values(SESSION_STATUS).includes(status)
}

export function isStudentPresent(status) {
  return status === ATTENDANCE_STATUS.PRESENT
}

export function isStudentAbsent(status) {
  return status === ATTENDANCE_STATUS.ABSENT
}

export function canSubmitSession(status) {
  return status === SESSION_STATUS.DRAFT
}

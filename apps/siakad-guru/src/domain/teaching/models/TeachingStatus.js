export const TEACHING_STATUS = {
  SCHEDULED: 'scheduled',
  ACTIVE: 'active',
  COMPLETED: 'completed',
  CANCELLED: 'cancelled',
}

export const TEACHER_PRESENCE = {
  PRESENT: 'present',
  ABSENT: 'absent',
  LATE: 'late',
  DELEGATED: 'delegated',
}

export function isValidTeachingStatus(status) {
  return Object.values(TEACHING_STATUS).includes(status)
}

export function isValidTeacherPresence(presence) {
  return Object.values(TEACHER_PRESENCE).includes(presence)
}

export function canStartSession(status) {
  return status === TEACHING_STATUS.SCHEDULED
}

export function canEndSession(status) {
  return status === TEACHING_STATUS.ACTIVE
}

export function canCancelSession(status) {
  return status === TEACHING_STATUS.SCHEDULED || status === TEACHING_STATUS.ACTIVE
}

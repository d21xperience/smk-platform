export const ASSESSMENT_STATUS = {
  DRAFT: 'draft',
  FINALIZED: 'finalized'
}

export const GRADE = {
  A: 'A',
  B: 'B',
  C: 'C',
  D: 'D',
  E: 'E'
}

export const PREDICATE = {
  EXCELLENT: 'Sangat Baik',
  GOOD: 'Baik',
  SUFFICIENT: 'Cukup',
  POOR: 'Kurang',
  VERY_POOR: 'Sangat Kurang'
}

export function isValidAssessmentStatus(status) {
  return Object.values(ASSESSMENT_STATUS).includes(status)
}

export function canFinalize(status) {
  return status === ASSESSMENT_STATUS.DRAFT
}

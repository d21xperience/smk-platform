export const INVENTORY_CONDITION = {
  GOOD: 'good',
  DAMAGED: 'damaged',
  BROKEN: 'broken',
  MISSING: 'missing',
}

export const REPORT_METHOD = {
  QR: 'qr',
  MANUAL: 'manual',
}

export function isValidCondition(condition) {
  return Object.values(INVENTORY_CONDITION).includes(condition)
}

export function isValidReportMethod(method) {
  return Object.values(REPORT_METHOD).includes(method)
}

export function getConditionLabel(condition) {
  const labels = {
    [INVENTORY_CONDITION.GOOD]: 'Baik',
    [INVENTORY_CONDITION.DAMAGED]: 'Rusak Ringan',
    [INVENTORY_CONDITION.BROKEN]: 'Rusak Berat',
    [INVENTORY_CONDITION.MISSING]: 'Hilang',
  }
  return labels[condition] || condition
}

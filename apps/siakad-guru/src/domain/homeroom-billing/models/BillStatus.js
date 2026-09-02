export const BILL_STATUS = {
  UNPAID: 'unpaid',
  PARTIAL: 'partial',
  PAID: 'paid',
  OVERDUE: 'overdue',
  WAIVED: 'waived',
}

const BILL_STATUS_LABELS = {
  [BILL_STATUS.UNPAID]: 'Belum Bayar',
  [BILL_STATUS.PARTIAL]: 'Bayar Sebagian',
  [BILL_STATUS.PAID]: 'Lunas',
  [BILL_STATUS.OVERDUE]: 'Terlambat',
  [BILL_STATUS.WAIVED]: 'Dibebaskan',
}

const BILL_STATUS_COLORS = {
  [BILL_STATUS.UNPAID]: 'warning',
  [BILL_STATUS.PARTIAL]: 'info',
  [BILL_STATUS.PAID]: 'positive',
  [BILL_STATUS.OVERDUE]: 'negative',
  [BILL_STATUS.WAIVED]: 'grey',
}

export function isValidBillStatus(status) {
  return Object.values(BILL_STATUS).includes(status)
}

export function getBillStatusLabel(status) {
  return BILL_STATUS_LABELS[status] || status
}

export function getBillStatusColor(status) {
  return BILL_STATUS_COLORS[status] || 'grey'
}

export function isBillOutstanding(status) {
  return (
    status === BILL_STATUS.UNPAID ||
    status === BILL_STATUS.PARTIAL ||
    status === BILL_STATUS.OVERDUE
  )
}

export function isBillSettled(status) {
  return status === BILL_STATUS.PAID || status === BILL_STATUS.WAIVED
}

import { isBillOutstanding, isBillSettled } from './BillStatus.js'

export class StudentBill {
  constructor({
    studentId,
    studentName,
    billId,
    billType,
    billName,
    amount,
    paidAmount,
    outstandingAmount,
    status,
    dueDate,
  }) {
    this.studentId = studentId
    this.studentName = studentName
    this.billId = billId
    this.billType = billType
    this.billName = billName
    this.amount = amount
    this.paidAmount = paidAmount
    this.outstandingAmount = outstandingAmount
    this.status = status
    this.dueDate = dueDate
  }

  isOutstanding() {
    return isBillOutstanding(this.status)
  }

  isSettled() {
    return isBillSettled(this.status)
  }

  getPaymentPercentage() {
    if (this.amount === 0) return 100
    return Math.round((this.paidAmount / this.amount) * 100)
  }
}

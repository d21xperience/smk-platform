import { Invoice } from './Invoice'

export class StudentFinancialSummary {
  constructor(props) {
    this.studentId = props.studentId
    this.studentName = props.studentName
    this.classId = props.classId
    this.totalInvoiced = props.totalInvoiced || 0
    this.totalPaid = props.totalPaid || 0
    this.balance = props.balance || 0 // totalInvoiced - totalPaid
    this.lastPaymentDate = props.lastPaymentDate || null
    this.status = props.status || 'LUNAS' // LUNAS, TUNGGAKAN, BELUM_BAYAR
    this.invoices = (props.invoices || []).map((i) => new Invoice(i))
  }

  isPaid() {
    return this.status === 'LUNAS'
  }

  isOverdue() {
    return this.status === 'TUNGGAKAN'
  }

  toJSON() {
    return {
      studentId: this.studentId,
      studentName: this.studentName,
      classId: this.classId,
      totalInvoiced: this.totalInvoiced,
      totalPaid: this.totalPaid,
      balance: this.balance,
      lastPaymentDate: this.lastPaymentDate,
      status: this.status,
      invoices: this.invoices.map((i) => i.toJSON()),
    }
  }
}

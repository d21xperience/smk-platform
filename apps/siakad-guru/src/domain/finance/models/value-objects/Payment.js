export class Payment {
  constructor(props) {
    this.id = props.id
    this.studentId = props.studentId
    this.invoiceId = props.invoiceId
    this.amount = props.amount
    this.paymentDate = props.paymentDate
    this.method = props.method // cash, transfer, dll.
    this.reference = props.reference || ''
  }

  toJSON() {
    return {
      id: this.id,
      studentId: this.studentId,
      invoiceId: this.invoiceId,
      amount: this.amount,
      paymentDate: this.paymentDate,
      method: this.method,
      reference: this.reference,
    }
  }
}

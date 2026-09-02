export class Invoice {
  constructor(props) {
    this.id = props.id;
    this.studentId = props.studentId;
    this.period = props.period; // '2026-08'
    this.amount = props.amount;
    this.status = props.status; // 'paid', 'unpaid', 'partial'
    this.dueDate = props.dueDate;
    this.paidDate = props.paidDate || null;
  }

  isPaid() {
    return this.status === 'paid';
  }

  isUnpaid() {
    return this.status === 'unpaid';
  }

  isOverdue(today = new Date()) {
    return this.isUnpaid() && new Date(this.dueDate) < today;
  }

  toJSON() {
    return {
      id: this.id,
      studentId: this.studentId,
      period: this.period,
      amount: this.amount,
      status: this.status,
      dueDate: this.dueDate,
      paidDate: this.paidDate,
    };
  }
}

import { BaseEntity } from '@/domain/shared/BaseEntity'

export const ISSUE_STATUS = {
  BORROWED: 'borrowed',
  RETURNED: 'returned',
  OVERDUE: 'overdue',
}

export class AssetIssue extends BaseEntity {
  constructor(props) {
    super(props)
    this.assetId = props.assetId
    this.borrowerName = props.borrowerName
    this.borrowerType = props.borrowerType || 'teacher' // teacher, student, staff
    this.purpose = props.purpose || ''
    this.borrowedDate = props.borrowedDate || new Date().toISOString().split('T')[0]
    this.expectedReturnDate = props.expectedReturnDate || ''
    this.returnedDate = props.returnedDate || null
    this.status = props.status || ISSUE_STATUS.BORROWED
    this.notes = props.notes || ''
  }

  isBorrowed() {
    return this.status === ISSUE_STATUS.BORROWED
  }
  isReturned() {
    return this.status === ISSUE_STATUS.RETURNED
  }

  returnItem(returnDate = new Date().toISOString().split('T')[0]) {
    if (this.isReturned()) throw new Error('Aset sudah dikembalikan.')
    this.status = ISSUE_STATUS.RETURNED
    this.returnedDate = returnDate
    this.updateTimestamps()
  }

  isOverdue(today = new Date().toISOString().split('T')[0]) {
    return this.isBorrowed() && this.expectedReturnDate < today
  }

  toJSON() {
    return {
      ...super.toJSON(),
      assetId: this.assetId,
      borrowerName: this.borrowerName,
      borrowerType: this.borrowerType,
      purpose: this.purpose,
      borrowedDate: this.borrowedDate,
      expectedReturnDate: this.expectedReturnDate,
      returnedDate: this.returnedDate,
      status: this.status,
      notes: this.notes,
    }
  }
}

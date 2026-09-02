import { StudentBill } from '../domain/homeroom-billing/models/StudentBill.js'
import { HomeroomBillingSummary } from '../domain/homeroom-billing/models/HomeroomBillingSummary.js'
import { BILL_STATUS } from '../domain/homeroom-billing/models/BillStatus.js'

export class HomeroomBillingProjection {
  static buildBills(rawBills) {
    if (!Array.isArray(rawBills)) {
      return []
    }
    return rawBills.map((bill) => new StudentBill(bill))
  }

  static buildSummary({ classInfo, bills }) {
    const totalAmount = bills.reduce((sum, b) => sum + b.amount, 0)
    const totalPaidAmount = bills.reduce((sum, b) => sum + b.paidAmount, 0)
    const totalOutstandingAmount = bills.reduce((sum, b) => sum + b.outstandingAmount, 0)

    const uniqueStudents = new Set(bills.map((b) => b.studentId))

    return new HomeroomBillingSummary({
      classId: classInfo.classId,
      className: classInfo.className,
      academicYearId: classInfo.academicYearId,
      academicYearName: classInfo.academicYearName,
      semesterId: classInfo.semesterId,
      semesterName: classInfo.semesterName,
      totalStudents: uniqueStudents.size,
      totalBills: bills.length,
      totalAmount,
      totalPaidAmount,
      totalOutstandingAmount,
      paidCount: bills.filter((b) => b.status === BILL_STATUS.PAID).length,
      partialCount: bills.filter((b) => b.status === BILL_STATUS.PARTIAL).length,
      unpaidCount: bills.filter((b) => b.status === BILL_STATUS.UNPAID).length,
      overdueCount: bills.filter((b) => b.status === BILL_STATUS.OVERDUE).length,
      waivedCount: bills.filter((b) => b.status === BILL_STATUS.WAIVED).length,
    })
  }

  static buildReadModel({ classInfo, rawBills }) {
    const bills = HomeroomBillingProjection.buildBills(rawBills)
    const summary = HomeroomBillingProjection.buildSummary({ classInfo, bills })
    return { bills, summary }
  }
}

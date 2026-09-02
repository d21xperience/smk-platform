export class HomeroomBillingSummary {
  constructor({
    classId,
    className,
    academicYearId,
    academicYearName,
    semesterId,
    semesterName,
    totalStudents,
    totalBills,
    totalAmount,
    totalPaidAmount,
    totalOutstandingAmount,
    paidCount,
    partialCount,
    unpaidCount,
    overdueCount,
    waivedCount
  }) {
    this.classId = classId
    this.className = className
    this.academicYearId = academicYearId
    this.academicYearName = academicYearName
    this.semesterId = semesterId
    this.semesterName = semesterName
    this.totalStudents = totalStudents || 0
    this.totalBills = totalBills || 0
    this.totalAmount = totalAmount || 0
    this.totalPaidAmount = totalPaidAmount || 0
    this.totalOutstandingAmount = totalOutstandingAmount || 0
    this.paidCount = paidCount || 0
    this.partialCount = partialCount || 0
    this.unpaidCount = unpaidCount || 0
    this.overdueCount = overdueCount || 0
    this.waivedCount = waivedCount || 0
  }

  getCollectionRate() {
    if (this.totalAmount === 0) return 100
    return Math.round((this.totalPaidAmount / this.totalAmount) * 100)
  }

  getOutstandingRate() {
    if (this.totalAmount === 0) return 0
    return Math.round((this.totalOutstandingAmount / this.totalAmount) * 100)
  }
}

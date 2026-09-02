import { REPORT_METHOD } from './InventoryCondition.js'

export class DamageReport {
  constructor({
    id,
    itemId,
    itemName,
    condition,
    description,
    reportMethod,
    date,
    teacherId,
    teacherName,
    classId,
    className,
    schoolId,
    academicYearId,
    semesterId,
  }) {
    this.id = id
    this.itemId = itemId
    this.itemName = itemName
    this.condition = condition
    this.description = description || ''
    this.reportMethod = reportMethod || REPORT_METHOD.MANUAL
    this.date = date
    this.teacherId = teacherId
    this.teacherName = teacherName
    this.classId = classId
    this.className = className
    this.schoolId = schoolId
    this.academicYearId = academicYearId
    this.semesterId = semesterId
  }

  isQRReport() {
    return this.reportMethod === REPORT_METHOD.QR
  }

  isManualReport() {
    return this.reportMethod === REPORT_METHOD.MANUAL
  }

  isValid() {
    return !!(
      this.id &&
      this.itemId &&
      this.itemName &&
      this.condition &&
      this.date &&
      this.teacherId &&
      this.classId &&
      this.schoolId &&
      this.academicYearId &&
      this.semesterId
    )
  }
}

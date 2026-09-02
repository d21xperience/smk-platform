export class InventoryItem {
  constructor({ id, name, category, location, classId, className, hasQRCode, schoolId }) {
    this.id = id
    this.name = name
    this.category = category
    this.location = location || ''
    this.classId = classId
    this.className = className
    this.hasQRCode = hasQRCode || false
    this.schoolId = schoolId
  }

  canScanQR() {
    return this.hasQRCode
  }

  isValid() {
    return !!(this.id && this.name && this.category && this.classId && this.schoolId)
  }
}

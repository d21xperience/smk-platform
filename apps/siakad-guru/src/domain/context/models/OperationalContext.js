export class OperationalContext {
  constructor({ schoolId, academicYearId, semesterId, userId, role, permissions = [] }) {
    this.schoolId = schoolId
    this.academicYearId = academicYearId
    this.semesterId = semesterId
    this.userId = userId
    this.role = role
    this.permissions = permissions
  }

  isValid() {
    return !!(this.schoolId && this.academicYearId && this.semesterId && this.userId && this.role)
  }
}

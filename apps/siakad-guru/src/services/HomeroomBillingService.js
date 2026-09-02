import { HomeroomBillingProjection } from '../projection/HomeroomBillingProjection.js'
export class HomeroomBillingService {
  constructor({ homeroomBillingAdapter }) {
    this.homeroomBillingAdapter = homeroomBillingAdapter
  }
  async getHomeroomAssignment({ teacherId }) {
    console.log(teacherId)
    if (!teacherId) {
      throw new Error('Teacher ID is required')
    }
    return await this.homeroomBillingAdapter.fetchHomeroomAssignment({ teacherId })
  }
  async getBillingData({ classId, academicYearId, semesterId, teacherId }) {
    if (!classId) {
      throw new Error('Class ID is required')
    }
    if (!academicYearId) {
      throw new Error('Academic Year ID is required')
    }
    if (!semesterId) {
      throw new Error('Semester ID is required')
    }
    if (!teacherId) {
      throw new Error('Teacher ID is required')
    }
    const rawData = await this.homeroomBillingAdapter.fetchBillingData({
      classId,
      academicYearId,
      semesterId,
      teacherId,
    })
    const readModel = HomeroomBillingProjection.buildReadModel({
      classInfo: rawData.classInfo,
      rawBills: rawData.bills,
    })
    return readModel
  }
}

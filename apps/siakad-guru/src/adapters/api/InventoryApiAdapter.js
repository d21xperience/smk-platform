import { BaseApiAdapter } from './BaseApiAdapter.js'

export class InventoryApiAdapter extends BaseApiAdapter {
  constructor() {
    super('/inventory')
  }

  async fetchItemsByClass({ classId }) {
    return await this.get('/items', { classId })
  }

  async fetchItemById({ itemId }) {
    return await this.get(`/items/${itemId}`)
  }

  async submitDamageReport({
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
    return await this.post('/reports', {
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
    })
  }
}

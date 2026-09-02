import { BaseApiAdapter } from './BaseApiAdapter.js'

export class ReportingApiAdapter extends BaseApiAdapter {
  constructor() {
    super('/reporting')
  }

  async fetchDashboardSummary({ schoolId, academicYearId, semesterId, teacherId }) {
    return await this.get('/dashboard', {
      schoolId,
      academicYearId,
      semesterId,
      teacherId,
    })
  }
}

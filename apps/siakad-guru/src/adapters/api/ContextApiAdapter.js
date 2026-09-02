import { BaseApiAdapter } from './BaseApiAdapter.js'

export class ContextApiAdapter extends BaseApiAdapter {
  constructor() {
    super('/context')
  }

  async fetchAvailableContexts() {
    const data = await this.get('/available')
    return {
      school: data.school,
      academicYears: data.academicYears,
      user: data.user,
    }
  }

  async selectContext({ schoolId, academicYearId, semesterId, user }) {
    const data = await this.post('/select', {
      schoolId,
      academicYearId,
      semesterId,
    })
    return {
      schoolId: data.schoolId,
      academicYearId: data.academicYearId,
      semesterId: data.semesterId,
      userId: user.id,
      role: user.role,
      permissions: user.permissions,
    }
  }
}

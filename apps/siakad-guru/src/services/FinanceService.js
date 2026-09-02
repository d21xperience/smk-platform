import { FinanceApi } from '@/adapters/api/finance.api'
import { useContextStore } from '@/stores/context.store'
import { BaseService } from './BaseService'

class FinanceService extends BaseService {
  constructor() {
    super(useContextStore())
  }

  async getClassFinancialSummary(classId) {
    const context = this.getContext()
    if (!context || !context.isValid()) throw new Error('Context tidak tersedia.')
    const response = await FinanceApi.getClassFinancialSummary({
      classId,
      semesterId: context.semesterId,
      academicYearId: context.academicYearId,
    })
    return response.data
  }

  async getStudentFinancialDetail(studentId) {
    const context = this.getContext()
    if (!context || !context.isValid()) throw new Error('Context tidak tersedia.')
    const response = await FinanceApi.getStudentFinancialDetail({
      studentId,
      semesterId: context.semesterId,
      academicYearId: context.academicYearId,
    })
    return response.data
  }
}

export const financeService = new FinanceService()

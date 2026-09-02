import MockAdapter from 'axios-mock-adapter'
import { setupAuthMock } from './handlers/auth.mock'
import { setupContextMock } from './handlers/context.mock'
import { setupTeachingMock } from './handlers/teaching.mock'
import { setupAttendanceMock } from './handlers/attendance.mock'
import { setupAssessmentMock } from './handlers/assessment.mock'
import { setupInventoryMock } from './handlers/inventory.mock'
import { setupFinanceMock } from './handlers/finance.mock'
import { setupStudentProgressMock } from './handlers/studentProgress.mock'
import { setupReportingMock } from './handlers/reporting.mock'
import { setupStudentsMock } from '@/adapters/mock/handlers/students.mock.js'
let mockInstance = null

export function setupMocks(axiosInstance) {
  if (mockInstance) {
    mockInstance.restore()
  }
  mockInstance = new MockAdapter(axiosInstance, { delayResponse: 300 })

  // Register all mock handlers
  setupAuthMock(mockInstance)
  setupContextMock(mockInstance)
  setupTeachingMock(mockInstance)
  setupAttendanceMock(mockInstance)
  setupAssessmentMock(mockInstance)
  setupInventoryMock(mockInstance)
  setupFinanceMock(mockInstance)
  setupStudentProgressMock(mockInstance)
  setupReportingMock(mockInstance)
  setupStudentsMock(mockInstance)
  console.log('[Mock] All mocks registered successfully.')
  return mockInstance
}

export function getMockInstance() {
  return mockInstance
}

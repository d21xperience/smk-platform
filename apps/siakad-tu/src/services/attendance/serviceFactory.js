// apps/siakad-tu/src/services/attendance/serviceFactory.js

import { attendanceApi } from '@/adapters/api/attendanceApi.js'

export function createAttendanceCommandService(dependencies = {}) {
  const adapter = dependencies.adapter || attendanceApi
  const { AttendanceCommandService } = require('./AttendanceCommandService.js')
  return new AttendanceCommandService({ adapter })
}

export function createAttendanceQueryService(dependencies = {}) {
  const adapter = dependencies.adapter || attendanceApi
  const { AttendanceQueryService } = require('./AttendanceQueryService.js')
  return new AttendanceQueryService({ adapter })
}

export const useAttendanceCommandService = () => createAttendanceCommandService()
export const useAttendanceQueryService = () => createAttendanceQueryService()

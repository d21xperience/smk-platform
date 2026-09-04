// apps/siakad-tu/src/adapters/api/attendanceApi.js

import { attendanceMockAdapter } from '../mock/attendanceMockAdapter.js'
import { attendanceRealAdapter } from './attendanceRealAdapter.js'

const useApi = import.meta.env?.QCLI_USE_API === 'true'
const forceMock = import.meta.env?.QCLI_MOCK_MODE === 'true'

export const attendanceApi = forceMock || !useApi ? attendanceMockAdapter : attendanceRealAdapter

export function __forceMockAttendanceAdapter() {
  return attendanceMockAdapter
}

export function __forceRealAttendanceAdapter() {
  return attendanceRealAdapter
}

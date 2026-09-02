/* eslint-disable no-unused-vars */
import { EVENTS } from '../EventCatalog'
// import { reportingService } from 'src/services/ReportingService'

export function registerReportingListeners(eventBus) {
  // Ketika assessment difinalisasi, update dashboard stats cache
  eventBus.on(EVENTS.ASSESSMENT_FINALIZED, async (payload) => {
    console.log('[ReportingListener] Assessment finalized, invalidating dashboard cache...')
    // Di aplikasi nyata, kita bisa panggil API untuk refresh cache di backend.
    // Atau di frontend, kita reload dashboard stats.
    // Untuk mock, kita hanya log.
    try {
      // await reportingService.getDashboardStats(); // ini akan memuat ulang
      // Tapi kita tidak mau terlalu agresif, kita trigger event ke UI untuk refresh.
      eventBus.emit(EVENTS.DASHBOARD_DATA_UPDATED, { source: 'ASSESSMENT_FINALIZED' })
    } catch (error) {
      console.error('[ReportingListener] Error refreshing dashboard:', error)
    }
  })

  // Ketika attendance disubmit
  eventBus.on(EVENTS.ATTENDANCE_SUBMITTED, () => {
    console.log('[ReportingListener] Attendance submitted, updating attendance stats...')
    // eventBus.emit(EVENTS.DASHBOARD_DATA_UPDATED, { source: 'ATTENDANCE_SUBMITTED' });
  })

  // Ketika student progress berubah (achievement/violation)
  eventBus.on(EVENTS.ACHIEVEMENT_ADDED, () => {
    console.log('[ReportingListener] Achievement added, updating progress stats...')
    // eventBus.emit(EVENTS.DASHBOARD_DATA_UPDATED, { source: 'ACHIEVEMENT_ADDED' });
  })
}

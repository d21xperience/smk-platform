// src/boot/events.js
import { EventDispatcher } from '@/events/EventDispatcher'
import { useTeachingStore } from '@/stores/teaching.store'
import { DomainEventType } from '@/events/DomainEventType'
import { AuditAction } from '@/engine/TeachingEngine'
import { useStudentProgressStore } from '@/stores/studentProgress.store'
import { useStatisticsCacheStore } from '@/stores/statisticsCache.store'
import { useReportingStore } from '@/stores/reporting.store'
let dispatcher = null

export function getEventDispatcher() {
  if (!dispatcher) {
    dispatcher = new EventDispatcher()
    setupInternalListeners(dispatcher)
  }
  return dispatcher
}

function setupInternalListeners(disp) {
  disp.on(DomainEventType.TEACHING_SESSION_STARTED, (event) => {
    const teaching = useTeachingStore()
    teaching.addAuditEntry(AuditAction.SESSION_STARTED, `Sesi ${event.payload.className} dimulai`)
  })
  // ... listener lain sesuai kebutuhan
}

export default ({ app }) => {
  const dispatcherInstance = getEventDispatcher()
  app.config.globalProperties.$eventBus = dispatcherInstance
  app.provide('eventBus', dispatcherInstance)

  // Inisialisasi store yang bergantung pada event
  const studentProgressStore = useStudentProgressStore()
  useStatisticsCacheStore().init()
  useStudentProgressStore().init()
  studentProgressStore.init()

  const cacheStore = useStatisticsCacheStore()
  cacheStore.init()
}

const reportingStore = useReportingStore()
reportingStore.init({
  // Map projection name ke getter store
  studentProgress: () => useStudentProgressStore().getAllProgress,
  classAttendance: () => useStatisticsCacheStore().classAttendance,
  teacherDashboard: () => {
    const cache = useStatisticsCacheStore()
    // Ubah objek teacherDashboard ke array dengan teacherId & teacherName
    return Object.entries(cache.cache.teacherDashboard).map(([teacherId, data]) => ({
      teacherId,
      teacherName: `Guru ${teacherId}`, // placeholder, nanti dari auth
      ...data,
    }))
  },
  assessment: () => useStatisticsCacheStore().cache.assessment,
})

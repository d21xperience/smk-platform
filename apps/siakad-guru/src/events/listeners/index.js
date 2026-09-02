import { eventBus } from '../EventBus'
import { registerTeachingListeners } from './teaching.listeners'
import { registerAttendanceListeners } from './attendance.listeners'
import { registerLoggingListeners } from './logging.listeners'

let isRegistered = false

export function registerAllListeners() {
  if (isRegistered) {
    console.warn('[EventBus] Listeners already registered. Skipping.')
    return
  }

  console.log('[EventBus] Registering all listeners...')

  registerTeachingListeners(eventBus)
  registerAttendanceListeners(eventBus)
  registerLoggingListeners(eventBus)

  isRegistered = true
  console.log('[EventBus] All listeners registered successfully.')
}

export function resetListeners() {
  eventBus.removeAllListeners()
  isRegistered = false
}

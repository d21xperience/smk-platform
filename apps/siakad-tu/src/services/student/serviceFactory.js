import { StudentEngine } from '@/engine/student/StudentEngine.js'
import { studentApi } from '@/adapters/api/studentApi.js'
import { eventDispatcher } from '@/events/eventDispatcher.js'
import { StudentCommandService } from './StudentCommandService.js'
import { StudentQueryService } from './StudentQueryService.js'
/**
Service Factory — Dependency Injection container untuk Student Domain.
Mendukung CQRS pattern:
CommandService: untuk write operations (create, update, delete, transfer, graduate)
QueryService: untuk read operations (get, list, check availability)
Keuntungan DI:
Mudah di-test (bisa inject mock dependencies)
Mudah di-swap (mock → real adapter)
Singleton per application lifecycle
Dependencies eksplisit, tidak tersembunyi
*/
// === SINGLETON INSTANCES ===
let commandServiceInstance = null
let queryServiceInstance = null
/**
Buat instance StudentCommandService dengan dependencies custom
(untuk testing atau environment khusus)
*/
export function createStudentCommandService(dependencies = {}) {
const engine = dependencies.engine || new StudentEngine()
const adapter = dependencies.adapter || studentApi
const eventDispatcherInstance = dependencies.eventDispatcher || eventDispatcher
return new StudentCommandService({
engine,
adapter,
eventDispatcher: eventDispatcherInstance,
})
}
/**
Buat instance StudentQueryService dengan dependencies custom
*/
export function createStudentQueryService(dependencies = {}) {
const adapter = dependencies.adapter || studentApi
return new StudentQueryService({ adapter })
}
/**
Singleton instance untuk StudentCommandService.
Dipanggil sekali, dipakai di seluruh aplikasi.
*/
export function useStudentCommandService() {
if (!commandServiceInstance) {
commandServiceInstance = createStudentCommandService()
}
return commandServiceInstance
}
/**
Singleton instance untuk StudentQueryService.
*/
export function useStudentQueryService() {
if (!queryServiceInstance) {
queryServiceInstance = createStudentQueryService()
}
return queryServiceInstance
}
/**
Backward compatibility: studentService (legacy)
Untuk kode lama yang masih menggunakan studentService langsung
*/
export const studentService = createStudentCommandService()
/**
Reset instances (untuk testing)
*/
export function __resetServiceInstances() {
commandServiceInstance = null
queryServiceInstance = null
}

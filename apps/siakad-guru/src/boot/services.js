// FILE: src/boot/services.js
// STATUS: MODIFY
// STATUS IMPLEMENTASI: COMPLETE

import { EventDispatcher } from '../events/Dispatcher.js'

// Storage Adapters
import { LocalStorageCredentialStorage } from '../adapters/storage/LocalStorageCredentialStorage.js'

// Credential Providers
import { PasswordCredentialProvider } from '../adapters/credential/PasswordCredentialProvider.js'
import { PasskeyCredentialProvider } from '../adapters/credential/PasskeyCredentialProvider.js'

// Mock Adapters (Mock First)
import { ContextMockAdapter } from '../adapters/mock/ContextMockAdapter.js'
import { AuthMockAdapter } from '../adapters/mock/AuthMockAdapter.js'
import { TeachingMockAdapter } from '../adapters/mock/TeachingMockAdapter.js'
import { AttendanceMockAdapter } from '../adapters/mock/AttendanceMockAdapter.js'
import { JournalMockAdapter } from '../adapters/mock/JournalMockAdapter.js'
import { AssessmentMockAdapter } from '../adapters/mock/AssessmentMockAdapter.js'
import { ProgressMockAdapter } from '../adapters/mock/ProgressMockAdapter.js'
import { ReportingMockAdapter } from '../adapters/mock/ReportingMockAdapter.js'
import { InventoryMockAdapter } from '../adapters/mock/InventoryMockAdapter.js'
import { TeacherDocumentMockAdapter } from '../adapters/mock/TeacherDocumentMockAdapter.js'
import { HomeroomBillingMockAdapter } from '../adapters/mock/HomeroomBillingMockAdapter.js'
import { HomeroomStudentMockAdapter } from '../adapters/mock/HomeroomStudentMockAdapter.js'
import { HomeroomProgressMockAdapter } from '../adapters/mock/HomeroomProgressMockAdapter.js'
import { HomeroomAttendanceMockAdapter } from '../adapters/mock/HomeroomAttendanceMockAdapter.js'
import { HomeroomJournalMockAdapter } from '../adapters/mock/HomeroomJournalMockAdapter.js'
import { HomeroomCommunicationMockAdapter } from '../adapters/mock/HomeroomCommunicationMockAdapter.js'

// API Adapters (Real Backend)
import { ContextApiAdapter } from '../adapters/api/ContextApiAdapter.js'
import { AuthApiAdapter } from '../adapters/api/AuthApiAdapter.js'
import { TeachingApiAdapter } from '../adapters/api/TeachingApiAdapter.js'
import { AttendanceApiAdapter } from '../adapters/api/AttendanceApiAdapter.js'
import { JournalApiAdapter } from '../adapters/api/JournalApiAdapter.js'
import { AssessmentApiAdapter } from '../adapters/api/AssessmentApiAdapter.js'
import { ProgressApiAdapter } from '../adapters/api/ProgressApiAdapter.js'
import { ReportingApiAdapter } from '../adapters/api/ReportingApiAdapter.js'
import { InventoryApiAdapter } from '../adapters/api/InventoryApiAdapter.js'

// Services
import { ContextService } from '../services/ContextService.js'
import { AuthService } from '../services/AuthService.js'
import { TeachingService } from '../services/TeachingService.js'
import { AttendanceService } from '../services/AttendanceService.js'
import { JournalService } from '../services/JournalService.js'
import { AssessmentService } from '../services/AssessmentService.js'
import { ProgressService } from '../services/ProgressService.js'
import { ReportingService } from '../services/ReportingService.js'
import { InventoryService } from '../services/InventoryService.js'
import { TeacherDocumentService } from '../services/TeacherDocumentService.js'
import { HomeroomBillingService } from '../services/HomeroomBillingService.js'
import { HomeroomStudentService } from '../services/HomeroomStudentService.js'
import { HomeroomProgressService } from '../services/HomeroomProgressService.js'
import { HomeroomAttendanceService } from '../services/HomeroomAttendanceService.js'
import { HomeroomJournalService } from '../services/HomeroomJournalService.js'
import { HomeroomCommunicationService } from '../services/HomeroomCommunicationService.js'

// 1. Initialize shared infrastructure
const eventDispatcher = new EventDispatcher()

// 2. Initialize Credential Storage
const credentialStorage = new LocalStorageCredentialStorage()

// 3. Determine adapter mode based on environment
// Default tetap Mock First.
// API hanya aktif jika QCLI_USE_API=true dan QCLI_MOCK_MODE tidak sedang true.
const USE_API = import.meta.env.QCLI_USE_API === 'true' && import.meta.env.QCLI_MOCK_MODE !== 'true'

// 4. Initialize Adapters (Mock First, switch to API when ready)
const contextAdapter = USE_API ? new ContextApiAdapter() : new ContextMockAdapter()
const authAdapter = USE_API ? new AuthApiAdapter() : new AuthMockAdapter()
const teachingAdapter = USE_API ? new TeachingApiAdapter() : new TeachingMockAdapter()
const attendanceAdapter = USE_API ? new AttendanceApiAdapter() : new AttendanceMockAdapter()
const journalAdapter = USE_API ? new JournalApiAdapter() : new JournalMockAdapter()
const assessmentAdapter = USE_API ? new AssessmentApiAdapter() : new AssessmentMockAdapter()
const progressAdapter = USE_API ? new ProgressApiAdapter() : new ProgressMockAdapter()
const reportingAdapter = USE_API ? new ReportingApiAdapter() : new ReportingMockAdapter()
const inventoryAdapter = USE_API ? new InventoryApiAdapter() : new InventoryMockAdapter()
const teacherDocumentAdapter = new TeacherDocumentMockAdapter()
const homeroomBillingAdapter = new HomeroomBillingMockAdapter()
const homeroomStudentAdapter = new HomeroomStudentMockAdapter()
const homeroomProgressAdapter = new HomeroomProgressMockAdapter()
const homeroomAttendanceAdapter = new HomeroomAttendanceMockAdapter()
const homeroomJournalAdapter = new HomeroomJournalMockAdapter()
const homeroomCommunicationAdapter = new HomeroomCommunicationMockAdapter()

// 5. Initialize Credential Providers
const passwordCredentialProvider = new PasswordCredentialProvider({ authAdapter })
const passkeyCredentialProvider = new PasskeyCredentialProvider({ authAdapter })

// 6. Initialize and export Services
export const contextService = new ContextService({
  contextAdapter,
  eventDispatcher,
})

export const authService = new AuthService({
  authAdapter,
  credentialStorage,
  credentialProvider: passwordCredentialProvider,
  eventDispatcher,
})

export const passkeyAuthService = new AuthService({
  authAdapter,
  credentialStorage,
  credentialProvider: passkeyCredentialProvider,
  eventDispatcher,
})

export const teachingService = new TeachingService({
  teachingAdapter,
  eventDispatcher,
})

export const attendanceService = new AttendanceService({
  attendanceAdapter,
  eventDispatcher,
})

export const journalService = new JournalService({
  journalAdapter,
  eventDispatcher,
})

export const assessmentService = new AssessmentService({
  assessmentAdapter,
  eventDispatcher,
})

export const progressService = new ProgressService({
  progressAdapter,
  eventDispatcher,
})

export const reportingService = new ReportingService({
  reportingAdapter,
})

export const inventoryService = new InventoryService({
  inventoryAdapter,
  eventDispatcher,
})

export const teacherDocumentService = new TeacherDocumentService({
  teacherDocumentAdapter,
  eventDispatcher,
})

export const homeroomBillingService = new HomeroomBillingService({
  homeroomBillingAdapter,
})

export const homeroomStudentService = new HomeroomStudentService({
  homeroomStudentAdapter,
})

export const homeroomProgressService = new HomeroomProgressService({
  homeroomProgressAdapter,
})

export const homeroomAttendanceService = new HomeroomAttendanceService({
  homeroomAttendanceAdapter,
})

export const homeroomJournalService = new HomeroomJournalService({
  homeroomJournalAdapter,
  eventDispatcher,
})

export const homeroomCommunicationService = new HomeroomCommunicationService({
  homeroomCommunicationAdapter,
  eventDispatcher,
})

// Export shared infrastructure for boot files
export { eventDispatcher, credentialStorage }

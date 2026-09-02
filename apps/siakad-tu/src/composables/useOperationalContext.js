import { reactive } from 'vue'

// Mock operational context.
// In production, this will be replaced by the Auth/Session store.
export const operationalContext = reactive({
  schoolId: 'SCH-1',
  academicYearId: 'AY-2026',
  semesterId: 'SEM-1',
  userId: 'USER-TU-1',
  role: 'tu',
  permissions: [
    'correspondence:view',
    'correspondence:create',
    'correspondence:update',
    'correspondence:delete',
    'correspondence:approve',
  ],
})

export function useOperationalContext() {
  return operationalContext
}

// Helper to change context for testing/development
export function setOperationalContext(newContext) {
  Object.assign(operationalContext, newContext)
}

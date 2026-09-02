/*
 * @sdp/context — Public Entry Point
 *
 * School Digital Platform — Shared Operational Context Package
 *
 * Menyediakan capability Operational Context yang reusable
 * untuk seluruh aplikasi SDP:
 * - siakad-tu
 * - siakad-guru
 * - psb
 * - website
 * - aplikasi masa depan
 *
 * PENGGUNAAN:
 * import { OperationalContext, ContextEngine, School } from '@sdp/context'
 *
 * DILARANG mengimport file internal secara langsung:
 * import { OperationalContext } from '@sdp/context/src/models/OperationalContext.js' ❌
 */
// Models
export { School } from './models/School.js'
export { AcademicYear } from './models/AcademicYear.js'
export { Semester } from './models/Semester.js'
export { OperationalContext } from './models/OperationalContext.js'
// Engine
export { ContextEngine } from './engine/ContextEngine.js'
// Contracts & Serialization
export {
  ContextContract,
  CONTEXT_STORAGE_KEY,
  serializeContext,
  isValidSerializedContext
} from './contracts/contextContract.js'

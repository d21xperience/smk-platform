/*
 * @sdp/storage — Public Entry Point
 *
 * School Digital Platform — Shared Storage Abstraction Package
 *
 * Menyediakan storage abstraction yang reusable untuk seluruh aplikasi SDP:
 * - siakad-tu
 * - siakad-guru
 * - psb
 * - website
 * - aplikasi masa depan
 *
 * PENGGUNAAN:
 * import { LocalStorageAdapter, buildContextAwareKey } from '@sdp/storage'
 *
 * DILARANG mengimport file internal secara langsung:
 * import { LocalStorageAdapter } from '@sdp/storage/src/adapters/LocalStorageAdapter.js' ❌
 */
// Abstractions
export { StorageInterface } from './abstractions/StorageInterface.js'
// Adapters
export { LocalStorageAdapter } from './adapters/LocalStorageAdapter.js'
export { SessionStorageAdapter } from './adapters/SessionStorageAdapter.js'
export { MemoryStorageAdapter } from './adapters/MemoryStorageAdapter.js'
// Utilities
export {
  buildStorageKey,
  buildContextAwareKey,
  buildUserScopedKey,
  parseStorageKey,
  hasKeyPrefix,
  KEY_SEPARATOR,
  DEFAULT_PREFIX
} from './utils/storageKeyBuilder.js'
/*
 * Factory: Membuat storage adapter berdasarkan type.
 *
 * @param {string} type - 'local'|'session'|'memory'
 * @param {Object} options - Adapter options
 * @returns {import('./abstractions/StorageInterface.js').StorageInterface}
 *
 * @example
 * const storage = createStorage('local', { prefix: 'sdp' })
 * storage.setItem('context', { schoolId: 'school-001' })
 * const context = storage.getItem('context')
 */
export function createStorage(type = 'local', options = {}) {
  switch (type) {
    case 'local':
      return new LocalStorageAdapter(options)
    case 'session':
      return new SessionStorageAdapter(options)
    case 'memory':
      return new MemoryStorageAdapter(options)
    default:
      throw new Error(`Unknown storage type: "${type}". Must be "local", "session", or "memory".`)
  }
}
/*
 * Factory: Membuat storage adapter dengan auto-detect.
 * Menggunakan localStorage jika tersedia, fallback ke memory.
 *
 * @param {Object} options - Adapter options
 * @returns {import('./abstractions/StorageInterface.js').StorageInterface}
 */
export function createAutoStorage(options = {}) {
  try {
    const testKey = '__sdp_auto_detect__'
    window.localStorage.setItem(testKey, 'test')
    window.localStorage.removeItem(testKey)
    return new LocalStorageAdapter(options)
  } catch (e) {
    console.warn('[createAutoStorage] localStorage not available, falling back to memory storage')
    return new MemoryStorageAdapter(options)
  }
}

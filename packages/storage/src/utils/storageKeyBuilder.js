/*
 * @sdp/storage — Utility: storageKeyBuilder
 *
 * Membangun storage keys yang context-aware.
 * Mencegah kebocoran data antar sekolah, tahun ajaran, dan semester.
 *
 * Contoh key:
 * - "sdp:context:school01:ay-2026:sem-20261"
 * - "sdp:inventory:school01:ay-2026:sem-20261:cache"
 * - "sdp:auth:token:user123"
 *
 * Pure JavaScript — TIDAK bergantung pada Vue, Pinia, Quasar, Axios.
 */
/*
 * Default separator untuk key segments.
 */
export const KEY_SEPARATOR = ":";
/*
 * Default prefix untuk SDP storage keys.
 */
export const DEFAULT_PREFIX = "sdp";
/*
 * Membangun storage key dari segments.
 *
 * @param {Array<string>} segments - Array of key segments
 * @param {Object} options
 * @param {string} options.prefix - Prefix (default: 'sdp')
 * @param {string} options.separator - Separator (default: ':')
 * @returns {string} Full storage key
 *
 * @example
 * buildStorageKey(['context', 'school01', 'ay-2026'])
 * // → "sdp:context:school01:ay-2026"
 */
export function buildStorageKey(
  segments,
  { prefix = DEFAULT_PREFIX, separator = KEY_SEPARATOR } = {},
) {
  if (!segments || !Array.isArray(segments) || segments.length === 0) {
    throw new Error(
      "Storage key segments are required and must be a non-empty array",
    );
  }
  // Filter out null/undefined/empty segments
  const validSegments = segments.filter(
    (s) => s !== null && s !== undefined && s !== "",
  );
  if (validSegments.length === 0) {
    throw new Error("Storage key must have at least one valid segment");
  }
  const parts = prefix ? [prefix, ...validSegments] : validSegments;
  return parts.join(separator);
}
/*
 * Membangun context-aware storage key.
 * Mencakup SchoolID, AcademicYearID, dan SemesterID.
 *
 * @param {Object} params
 * @param {string} params.domain - Domain name (e.g., 'inventory', 'attendance')
 * @param {string} params.schoolId - School ID
 * @param {string} params.academicYearId - Academic Year ID
 * @param {string} params.semesterId - Semester ID (optional)
 * @param {string|null} params.suffix - Additional suffix (optional)
 * @param {string} params.prefix - Key prefix (default: 'sdp')
 * @returns {string} Context-aware storage key
 *
 * @example
 * buildContextAwareKey({
 *   domain: 'inventory',
 *   schoolId: 'school-001',
 *   academicYearId: 'ay-2026',
 *   semesterId: 'sem-20261',
 *   suffix: 'cache'
 * })
 * // → "sdp:inventory:school-001:ay-2026:sem-20261:cache"
 */
export function buildContextAwareKey({
  domain,
  schoolId,
  academicYearId,
  semesterId,
  suffix,
  prefix = DEFAULT_PREFIX,
}) {
  if (!domain) {
    throw new Error("Domain is required for context-aware storage key");
  }
  if (!schoolId) {
    throw new Error("SchoolID is required for context-aware storage key");
  }
  if (!academicYearId) {
    throw new Error("AcademicYearID is required for context-aware storage key");
  }
  const segments = [domain, schoolId, academicYearId];
  if (semesterId) {
    segments.push(semesterId);
  }
  if (suffix) {
    segments.push(suffix);
  }
  return buildStorageKey(segments, { prefix });
}
/*
 * Membangun user-scoped storage key.
 *
 * @param {Object} params
 * @param {string} params.domain - Domain name
 * @param {string} params.userId - User ID
 * @param {string|null} params.suffix - Additional suffix
 * @param {string} params.prefix - Key prefix
 * @returns {string}
 *
 * @example
 * buildUserScopedKey({ domain: 'preferences', userId: 'user-001', suffix: 'theme' })
 * // → "sdp:preferences:user-001:theme"
 */
export function buildUserScopedKey({
  domain,
  userId,
  suffix,
  prefix = DEFAULT_PREFIX,
}) {
  if (!domain) {
    throw new Error("Domain is required for user-scoped storage key");
  }
  if (!userId) {
    throw new Error("UserID is required for user-scoped storage key");
  }
  const segments = [domain, userId];
  if (suffix) {
    segments.push(suffix);
  }
  return buildStorageKey(segments, { prefix });
}
/*
 * Parse storage key menjadi segments.
 *
 * @param {string} key - Full storage key
 * @param {Object} options
 * @param {string} options.separator - Separator (default: ':')
 * @returns {Array<string>} Array of segments
 *
 * @example
 * parseStorageKey('sdp:context:school01:ay-2026')
 * // → ['sdp', 'context', 'school01', 'ay-2026']
 */
export function parseStorageKey(key, { separator = KEY_SEPARATOR } = {}) {
  if (!key || typeof key !== "string") {
    return [];
  }
  return key.split(separator);
}
/*
 * Memeriksa apakah key memiliki prefix tertentu.
 *
 * @param {string} key - Full storage key
 * @param {string} prefix - Prefix to check
 * @param {string} separator - Separator (default: ':')
 * @returns {boolean}
 */
export function hasKeyPrefix(key, prefix, { separator = KEY_SEPARATOR } = {}) {
  if (!key || !prefix) return false;
  return key.startsWith(`${prefix}${separator}`);
}

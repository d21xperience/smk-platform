// FILE: src/adapters/http/RequestGuard.js
// STATUS: NEW
// STATUS IMPLEMENTASI: COMPLETE

import { CONTEXT_MODE, ContextEngine } from '../../domain/context/engine/ContextEngine.js'

const READ_ONLY_ALLOWED_METHODS = Object.freeze(['get', 'head', 'options'])

const READ_ONLY_ALLOWED_URL_MARKERS = Object.freeze(['/export', '/print'])

function normalizeMethod(method) {
  if (!method) {
    return 'get'
  }

  return String(method).toLowerCase()
}

function normalizeUrl(url) {
  if (!url) {
    return ''
  }

  return String(url).toLowerCase()
}

function isReadOnlyAllowedRequest({ method, url }) {
  const normalizedMethod = normalizeMethod(method)
  const normalizedUrl = normalizeUrl(url)

  if (READ_ONLY_ALLOWED_METHODS.includes(normalizedMethod)) {
    return true
  }

  return READ_ONLY_ALLOWED_URL_MARKERS.some((marker) => {
    return normalizedUrl.includes(marker)
  })
}

export function ensureRequestAllowed({ method, url, mode, context }) {
  if (mode !== CONTEXT_MODE.HISTORICAL) {
    return
  }

  if (isReadOnlyAllowedRequest({ method, url })) {
    return
  }

  const user = {
    id: context && context.userId ? context.userId : null,
    role: context && context.role ? context.role : null,
    permissions: context && Array.isArray(context.permissions) ? context.permissions : [],
  }

  ContextEngine.assertMutationAllowed({
    isHistoryMode: true,
    user,
  })
}

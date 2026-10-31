// FILE: src/utils/requestContextRegistry.js
// STATUS: MODIFY
// STATUS IMPLEMENTASI: COMPLETE

const requestState = {
  authToken: null,
  mode: null,
  context: null,
}

export function setAuthToken(token) {
  requestState.authToken = token ? String(token) : null
}

export function clearAuthToken() {
  requestState.authToken = null
}

export function setRequestContext({ mode, context }) {
  requestState.mode = mode || null
  requestState.context = context ? { ...context } : null
}

export function clearRequestContext() {
  requestState.mode = null
  requestState.context = null
}

export function getRequestMeta() {
  return {
    mode: requestState.mode,
    context: requestState.context ? { ...requestState.context } : null,
  }
}

export function getRequestHeaders() {
  const headers = {}

  if (requestState.authToken) {
    headers['Authorization'] = `Bearer ${requestState.authToken}`
  }

  if (requestState.mode) {
    headers['X-Context-Mode'] = requestState.mode
  }

  const context = requestState.context

  if (context) {
    if (context.schoolId) {
      headers['x-school-id'] = String(context.schoolId)
    }

    if (context.academicYearId) {
      headers['x-academic-year-id'] = String(context.academicYearId)
    }

    if (context.semesterId) {
      headers['x-semester-id'] = String(context.semesterId)
    }

    if (context.userId) {
      headers['x-user-id'] = String(context.userId)
    }

    if (context.role) {
      headers['x-user-role'] = String(context.role)
    }
  }

  return headers
}

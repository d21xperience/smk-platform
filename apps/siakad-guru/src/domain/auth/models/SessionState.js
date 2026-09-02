export const SESSION_STATE = {
  ANONYMOUS: 'anonymous',
  AUTHENTICATING: 'authenticating',
  AUTHENTICATED: 'authenticated',
  RESTORING: 'restoring',
  EXPIRED: 'expired',
  LOGGING_OUT: 'logging_out',
}

export function isValidSessionState(state) {
  return Object.values(SESSION_STATE).includes(state)
}

export function isSessionActive(state) {
  return state === SESSION_STATE.AUTHENTICATED
}

export function isSessionTransitioning(state) {
  return (
    state === SESSION_STATE.AUTHENTICATING ||
    state === SESSION_STATE.RESTORING ||
    state === SESSION_STATE.LOGGING_OUT
  )
}

export function canAttemptLogin(state) {
  return state === SESSION_STATE.ANONYMOUS || state === SESSION_STATE.EXPIRED
}

export function canAccessApplication(state) {
  return state === SESSION_STATE.AUTHENTICATED
}

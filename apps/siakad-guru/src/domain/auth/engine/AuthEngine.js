import { User } from '../models/User.js'
import { AuthSession } from '../models/AuthSession.js'
import { SESSION_STATE } from '../models/SessionState.js'

export class AuthEngine {
  static validateCredentials(username, password) {
    if (!username || typeof username !== 'string' || username.trim().length === 0) {
      throw new Error('Username is required')
    }
    if (!password || typeof password !== 'string' || password.length < 6) {
      throw new Error('Password must be at least 6 characters')
    }
    return true
  }

  static validateSession(session) {
    if (!session || !(session instanceof AuthSession)) {
      throw new Error('Invalid session object')
    }
    if (!session.isValid()) {
      throw new Error('Session is invalid or expired')
    }
    return true
  }

  static checkSessionExpiry(session) {
    if (!session) return SESSION_STATE.ANONYMOUS
    if (session.isExpired()) return SESSION_STATE.EXPIRED
    if (session.isValid()) return SESSION_STATE.AUTHENTICATED
    return SESSION_STATE.ANONYMOUS
  }

  static checkPermission(session, requiredPermission) {
    AuthEngine.validateSession(session)
    if (!(session.user instanceof User)) {
      throw new Error('Session does not contain a valid User object')
    }

    const hasAccess = session.user.hasPermission(requiredPermission)
    if (!hasAccess) {
      throw new Error(`Permission denied: ${requiredPermission}`)
    }

    return true
  }
}

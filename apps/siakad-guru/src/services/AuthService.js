import { AuthEngine } from '../domain/auth/engine/AuthEngine.js'
import { User } from '../domain/auth/models/User.js'
import { AuthSession } from '../domain/auth/models/AuthSession.js'
import { SESSION_STATE } from '../domain/auth/models/SessionState.js'
import { UserLoggedInEvent } from '../domain/auth/events/UserLoggedIn.js'
import { UserLoggedOutEvent } from '../domain/auth/events/UserLoggedOut.js'
import { SessionRestoredEvent } from '../domain/auth/events/SessionRestored.js'
import { SessionExpiredEvent } from '../domain/auth/events/SessionExpired.js'
import { AuthenticationFailedEvent } from '../domain/auth/events/AuthenticationFailed.js'

export class AuthService {
  constructor({ authAdapter, credentialStorage, credentialProvider, eventDispatcher }) {
    this.authAdapter = authAdapter
    this.credentialStorage = credentialStorage
    this.credentialProvider = credentialProvider
    this.eventDispatcher = eventDispatcher
  }

  async login(username, password) {
    AuthEngine.validateCredentials(username, password)

    try {
      const credential = await this.credentialProvider.authenticate(username, password)

      const user = new User(credential.user)
      const session = new AuthSession({
        user,
        token: credential.token,
        expiresAt: credential.expiresAt,
        credentialType: credential.credentialType
      })

      AuthEngine.validateSession(session)

      this.credentialStorage.save(session)

      if (this.eventDispatcher) {
        this.eventDispatcher.dispatch(new UserLoggedInEvent(user))
      }

      return session
    } catch (err) {
      if (this.eventDispatcher) {
        this.eventDispatcher.dispatch(new AuthenticationFailedEvent(username, err.message))
      }
      throw err
    }
  }

  async loginWithPasskey(challenge) {
    try {
      const credential = await this.credentialProvider.authenticate(challenge)

      const user = new User(credential.user)
      const session = new AuthSession({
        user,
        token: credential.token,
        expiresAt: credential.expiresAt,
        credentialType: credential.credentialType
      })

      AuthEngine.validateSession(session)

      this.credentialStorage.save(session)

      if (this.eventDispatcher) {
        this.eventDispatcher.dispatch(new UserLoggedInEvent(user))
      }

      return session
    } catch (err) {
      if (this.eventDispatcher) {
        this.eventDispatcher.dispatch(new AuthenticationFailedEvent(null, err.message))
      }
      throw err
    }
  }

  async restoreSession() {
    if (!this.credentialStorage.exists()) {
      return null
    }

    const session = this.credentialStorage.load()
    if (!session) {
      return null
    }

    const sessionState = AuthEngine.checkSessionExpiry(session)

    if (sessionState === SESSION_STATE.EXPIRED) {
      this.credentialStorage.clear()
      if (this.eventDispatcher) {
        this.eventDispatcher.dispatch(new SessionExpiredEvent(session.user ? session.user.id : null))
      }
      return null
    }

    if (sessionState === SESSION_STATE.AUTHENTICATED) {
      // Optionally validate token with backend
      // For mock, we trust the session if not expired
      if (this.eventDispatcher) {
        this.eventDispatcher.dispatch(new SessionRestoredEvent(session))
      }
      return session
    }

    this.credentialStorage.clear()
    return null
  }

  async validateCurrentSession(session) {
    if (!session) return SESSION_STATE.ANONYMOUS

    const sessionState = AuthEngine.checkSessionExpiry(session)

    if (sessionState === SESSION_STATE.EXPIRED) {
      this.credentialStorage.clear()
      if (this.eventDispatcher) {
        this.eventDispatcher.dispatch(new SessionExpiredEvent(session.user ? session.user.id : null))
      }
      return SESSION_STATE.EXPIRED
    }

    return sessionState
  }

  async logout(userId) {
    await this.authAdapter.logout()

    this.credentialStorage.clear()

    if (this.eventDispatcher && userId) {
      this.eventDispatcher.dispatch(new UserLoggedOutEvent(userId))
    }

    return true
  }
}

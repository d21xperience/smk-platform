// import { SESSION_STATE } from './SessionState.js'

export class AuthSession {
  constructor({ user, token, expiresAt, credentialType }) {
    this.user = user
    this.token = token
    this.expiresAt = expiresAt ? new Date(expiresAt) : null
    this.credentialType = credentialType || null
  }

  isExpired() {
    if (!this.expiresAt) return false
    return new Date() > this.expiresAt
  }

  isValid() {
    return !!(this.user && this.token && !this.isExpired())
  }

  toPersistable() {
    return {
      user: this.user,
      token: this.token,
      expiresAt: this.expiresAt ? this.expiresAt.toISOString() : null,
      credentialType: this.credentialType,
    }
  }

  static fromPersistable(data) {
    if (!data) return null
    return new AuthSession({
      user: data.user,
      token: data.token,
      expiresAt: data.expiresAt,
      credentialType: data.credentialType,
    })
  }
}

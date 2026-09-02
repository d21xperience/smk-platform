import { EVENT_AUTHENTICATION_FAILED } from '../../../events/Type.js'

export class AuthenticationFailedEvent {
  constructor(username, reason) {
    this.type = EVENT_AUTHENTICATION_FAILED
    this.payload = {
      username: username || null,
      reason: reason || 'unknown',
      timestamp: new Date().toISOString(),
    }
  }
}

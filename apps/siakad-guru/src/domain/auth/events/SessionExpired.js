import { EVENT_SESSION_EXPIRED } from '../../../events/Type.js'

export class SessionExpiredEvent {
  constructor(userId) {
    this.type = EVENT_SESSION_EXPIRED
    this.payload = {
      userId: userId || null,
      timestamp: new Date().toISOString(),
    }
  }
}

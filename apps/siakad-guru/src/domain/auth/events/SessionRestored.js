import { EVENT_SESSION_RESTORED } from '../../../events/Type.js'

export class SessionRestoredEvent {
  constructor(session) {
    this.type = EVENT_SESSION_RESTORED
    this.payload = {
      userId: session.user.id,
      role: session.user.role,
      credentialType: session.credentialType,
      timestamp: new Date().toISOString(),
    }
  }
}

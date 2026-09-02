import { EVENT_USER_LOGGED_OUT } from '../../../events/Type.js'

export class UserLoggedOutEvent {
  constructor(userId) {
    this.type = EVENT_USER_LOGGED_OUT
    this.payload = {
      userId: userId,
      timestamp: new Date().toISOString()
    }
  }
}

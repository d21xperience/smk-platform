import { EVENT_USER_LOGGED_IN } from '../../../events/Type.js'

export class UserLoggedInEvent {
  constructor(user) {
    this.type = EVENT_USER_LOGGED_IN
    this.payload = {
      userId: user.id,
      role: user.role,
      timestamp: new Date().toISOString(),
    }
  }
}

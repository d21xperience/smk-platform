import { CREDENTIAL_TYPE } from '../../domain/auth/models/CredentialType.js'

export class PasswordCredentialProvider {
  constructor({ authAdapter }) {
    this.authAdapter = authAdapter
    this.type = CREDENTIAL_TYPE.PASSWORD
  }

  isAvailable() {
    return true
  }

  async authenticate(username, password) {
    if (!username || !password) {
      throw new Error('Username and password are required')
    }

    const responseData = await this.authAdapter.login(username, password)

    return {
      credentialType: this.type,
      token: responseData.token,
      expiresAt: responseData.expiresAt,
      user: responseData.user,
    }
  }

  async validate(session) {
    if (!session || !session.token) {
      return false
    }
    return !session.isExpired()
  }
}

import { CREDENTIAL_TYPE } from '../../domain/auth/models/CredentialType.js'

export class PasskeyCredentialProvider {
  constructor({ authAdapter }) {
    this.authAdapter = authAdapter
    this.type = CREDENTIAL_TYPE.PASSKEY
  }

  isAvailable() {
    if (typeof window === 'undefined') return false
    if (!window.PublicKeyCredential) return false
    if (!navigator.credentials) return false
    return true
  }

  async authenticate(challenge) {
    if (!this.isAvailable()) {
      throw new Error('Passkey/WebAuthn is not available in this browser')
    }

    const mockChallenge = challenge || 'mock-challenge-' + Date.now()

    // Simulate WebAuthn authentication delay
    await new Promise(resolve => setTimeout(resolve, 800))

    // Mock: use challenge to simulate unique auth flow
    const responseData = await this.authAdapter.loginWithPasskey(mockChallenge)

    return {
      credentialType: this.type,
      token: responseData.token,
      expiresAt: responseData.expiresAt,
      user: responseData.user
    }
  }

  async register() {
    if (!this.isAvailable()) {
      throw new Error('Passkey/WebAuthn is not available in this browser')
    }

    // Mock implementation for development
    // In production, this would call navigator.credentials.create()
    await new Promise(resolve => setTimeout(resolve, 800))

    return {
      success: true,
      credentialId: 'mock-credential-id-' + Date.now(),
      message: 'Passkey registered successfully (mock)'
    }
  }

  async validate(session) {
    if (!session || !session.token) {
      return false
    }
    return !session.isExpired()
  }
}

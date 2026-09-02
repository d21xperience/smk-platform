export const CREDENTIAL_TYPE = {
  PASSWORD: 'password',
  PASSKEY: 'passkey',
}

export function isValidCredentialType(type) {
  return Object.values(CREDENTIAL_TYPE).includes(type)
}

export function getCredentialTypeLabel(type) {
  const labels = {
    [CREDENTIAL_TYPE.PASSWORD]: 'Username & Password',
    [CREDENTIAL_TYPE.PASSKEY]: 'Passkey / WebAuthn',
  }
  return labels[type] || type
}

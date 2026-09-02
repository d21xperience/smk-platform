export class CredentialStorage {
  save() {
    throw new Error('CredentialStorage.save() must be implemented')
  }

  load() {
    throw new Error('CredentialStorage.load() must be implemented')
  }

  clear() {
    throw new Error('CredentialStorage.clear() must be implemented')
  }

  exists() {
    throw new Error('CredentialStorage.exists() must be implemented')
  }
}

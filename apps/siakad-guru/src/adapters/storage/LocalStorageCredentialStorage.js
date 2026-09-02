import { CredentialStorage } from './CredentialStorage.js'
import { AuthSession } from '../../domain/auth/models/AuthSession.js'

const STORAGE_KEY = 'sdp_auth_session'

export class LocalStorageCredentialStorage extends CredentialStorage {
  save(session) {
    if (!session || !session.isValid()) {
      throw new Error('Cannot save invalid session')
    }
    const persistable = session.toPersistable()
    localStorage.setItem(STORAGE_KEY, JSON.stringify(persistable))
  }

  load() {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (!raw) return null

    try {
      const data = JSON.parse(raw)
      return AuthSession.fromPersistable(data)
    } catch (e) {
      console.log(e)
      this.clear()
      return null
    }
  }

  clear() {
    localStorage.removeItem(STORAGE_KEY)
  }

  exists() {
    return localStorage.getItem(STORAGE_KEY) !== null
  }
}

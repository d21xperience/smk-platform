// domain/shared/Entity.js
export class Entity {
  constructor({ id, sourceSystem = 'LOCAL', metadata = {} }) {
    this.id = id // UUID internal
    this.sourceSystem = sourceSystem // 'DAPODIK', 'LOCAL', 'IMPORT'
    this.metadata = metadata
    this.createdAt = new Date().toISOString()
  }

  equals(other) {
    if (!other) return false
    return other.id === this.id
  }
}

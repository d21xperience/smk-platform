import uuid from 'uuid'
const { v4: uuidv4 } = uuid
export class BaseEntity {
  constructor(props) {
    this.id = props?.id || uuidv4()
    this.createdAt = props?.createdAt || new Date().toISOString()
    this.updatedAt = props?.updatedAt || new Date().toISOString()
    this.version = props?.version || 1
  }

  updateTimestamps() {
    this.updatedAt = new Date().toISOString()
    this.version += 1
  }

  equals(other) {
    if (!other) return false
    return this.id === other.id
  }

  toJSON() {
    return {
      id: this.id,
      createdAt: this.createdAt,
      updatedAt: this.updatedAt,
      version: this.version,
    }
  }
}

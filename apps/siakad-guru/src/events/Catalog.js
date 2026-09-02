import * as EventTypes from './Type.js'

export class EventCatalog {
  static getValidTypes() {
    return Object.values(EventTypes)
  }

  static isValidType(type) {
    return this.getValidTypes().includes(type)
  }
}

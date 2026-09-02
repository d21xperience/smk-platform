export const CONTEXT_SELECTED = 'CONTEXT_SELECTED'

export class ContextSelectedEvent {
  constructor(context) {
    this.type = CONTEXT_SELECTED
    this.payload = context
    this.timestamp = new Date().toISOString()
  }
}

// PURE JS
export class Projection {
  constructor({ name, initialState, updateFn }) {
    this.name = name
    this.state = initialState
    this.updateFn = updateFn
    this.listeners = [] // bisa ditambahkan subscriber
  }

  getState() {
    return this.state
  }

  applyEvent(event) {
    const newState = this.updateFn(this.state, event)
    this.state = newState
    return newState
  }

  // Untuk testing/reset
  reset(initialState) {
    this.state = initialState
  }
}

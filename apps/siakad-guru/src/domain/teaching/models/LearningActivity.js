// PURE JS - Value Object
export class LearningActivity {
  constructor(props) {
    this.step = props.step || '' // 'opening', 'main', 'closing', atau custom
    this.description = props.description || ''
    this.duration = props.duration || 0 // dalam menit (opsional)
    this.order = props.order || 0
  }

  isEmpty() {
    return !this.description || this.description.trim() === ''
  }

  equals(other) {
    if (!other) return false
    return (
      this.step === other.step &&
      this.description === other.description &&
      this.duration === other.duration
    )
  }

  toJSON() {
    return {
      step: this.step,
      description: this.description,
      duration: this.duration,
      order: this.order,
    }
  }

  static fromJSON(json) {
    return new LearningActivity(json)
  }
}
